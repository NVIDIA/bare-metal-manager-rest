// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

//Package config implements vault configuration controller
package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

var (
	maxStateChangesPerCycle = 15
	refreshRateDefault      = 5 * time.Second
	failMetric              prometheus.Counter
	succMetric              prometheus.Counter
)

// NewController returns a new vault controller
func NewController(vaultURL string, svcURL string, secretsRootPath string) (*VaultConfigController, error) {
	v := &VaultConfigController{
		tokenCh:         make(chan string),
		vaultURL:        vaultURL,
		refreshRate:     refreshRateDefault,
		secretsRootPath: secretsRootPath,
		svcURL:          svcURL,
	}
	if v.client == nil {
		vaultClient, err := vault.NewClient(&vault.Config{
			Address: v.vaultURL,
		})
		if err != nil {
			return nil, err
		}
		v.client = &vaultClientInternal{client: vaultClient}
	}
	return v, nil
}

// Start starts the controller
func (v *VaultConfigController) Start(ctx context.Context) {
	go v.run(ctx)
}

// TokenChan returns the channel
func (v *VaultConfigController) TokenChan() <-chan string {
	return v.tokenCh
}

// Token returns the token
func (v *VaultConfigController) Token() string {
	v.credsMutex.RLock()
	defer v.credsMutex.RUnlock()
	if v.creds == nil {
		return ""
	}

	return v.creds.RootToken
}

// CertManagerToken returns the vault token for cert manager
func (v *VaultConfigController) CertManagerToken() string {
	return v.certMgrToken
}

func (v *VaultConfigController) unsealShard() string {
	v.credsMutex.RLock()
	defer v.credsMutex.RUnlock()
	if v.creds == nil {
		return ""
	}

	return v.creds.Keys[0]
}

func (v *VaultConfigController) setToken(creds *vault.InitResponse) {
	v.credsMutex.Lock()
	defer v.credsMutex.Unlock()

	v.creds = creds
	v.client.SetToken(creds.RootToken)
}

// CurrentState gets the current fsm state as seen by the controller
func (v *VaultConfigController) CurrentState(ctx context.Context) (VaultState, error) {
	log := core.GetLogger(ctx)
	current, err := v.currentState(ctx)
	if err != nil {
		log.Errorf("Failed to serve current state with error: %v", err)
		return Unitialized, err
	}

	return current, nil
}

func (v *VaultConfigController) run(ctx context.Context) {
	log := core.GetLogger(ctx)
	succMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "vault_token_renewals",
			Help: "Total number of token renewals",
		},
	)
	failMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "vault_token_renewal_failures",
			Help: "Total number of token renewal_failures",
		},
	)

	if err := v.recoverState(); err != nil {
		log.Infof("Could not recover existing vault state, error: %v. Continuing...", err)
	}

	cmTokenPeriod, err := time.ParseDuration(certManagerTokenPeriod)
	if err != nil {
		log.Fatalf("ParseDuration %s - %v", certManagerTokenPeriod, err)
	}
	cmTokenRenewTime := cmTokenPeriod / 4
	log.Infof("cmTokenRenewTime: %s", cmTokenRenewTime.String())
	renewTicker := time.NewTicker(cmTokenRenewTime)
	defer renewTicker.Stop()
	prevState := Unitialized
loop:
	for {
		select {
		case <-ctx.Done():
			close(v.tokenCh)
			return
		case <-time.After(v.refreshRate):
			break
		case <-renewTicker.C:
			v.renewCMToken(ctx)
			continue loop
		}

		current, err := v.currentState(ctx)
		desired := Done
		log.Infof("Current state: %s, Desired state: %s", vaultStateToString[current], vaultStateToString[desired])

		if err != nil {
			log.Error("Failed to get current state with error:", err)
			continue
		}

		if current == desired {
			if prevState != desired {
				log.Infof("Desired state %s reached.", vaultStateToString[desired])
				prevState = desired
				// Send the token to the channel for initialization to continue
				log.Infof("Current state: %s, Desired state: %s", vaultStateToString[current], vaultStateToString[desired])
				log.Debug("sending vault token to channel to signal vault init complete")
				v.tokenCh <- v.Token()
			} else {
				log.Debugf("Current and Desired states equal: %s", vaultStateToString[desired])
			}
			continue
		}
		prevState = current
		log.Infof("Current state: %s, Desired state: %s", vaultStateToString[current], vaultStateToString[desired])

		err = v.reconcile(ctx, current, desired)
		if err != nil {
			log.Error("Failed to reconcile state with error:", err)
		}
	}
}

func (v *VaultConfigController) renewCMToken(ctx context.Context) {
	log := core.GetLogger(ctx)
	state, _ := v.currentState(ctx)
	if state != Done {
		log.Infof("Skipping renew - current state: %v", state)
		return
	}

	_, err := v.client.AuthToken().RenewWithContext(ctx, v.certMgrToken, 0)
	if err == nil {
		log.Infof("Renewed cert manager token")
		succMetric.Inc()
	} else {
		log.Errorf("renewCMToken - %v", err)
		failMetric.Inc()
	}
}

func (v *VaultConfigController) recoverState() error {
	token, err := v.getSecret(SecretTokenPath, SecretTokenName)
	if err != nil {
		return err
	}

	initResp := &vault.InitResponse{}
	err = json.NewDecoder(bytes.NewReader([]byte(token))).Decode(initResp)
	if err != nil {
		return err
	}

	v.setToken(initResp)

	return nil
}

// We iterate over the different possible states
// A state discovery function returns true if Vault is in this state
func (v *VaultConfigController) currentState(ctx context.Context) (VaultState, error) {
	discovery := []struct {
		state        VaultState
		discoverFunc ControllerStateDiscovery
	}{
		{state: Unitialized, discoverFunc: v.isUninitialized},
		{state: Sealed, discoverFunc: v.isSealed},
		{state: PKINotEnabled, discoverFunc: v.isPKINotEnabled},
		{state: PKICACertNotConfigured, discoverFunc: v.isPKICACertNotConfigured},
		{state: PKIURLsNotConfigured, discoverFunc: v.isPKIURLsNotConfigured},
		{state: PKIRoleNotConfigured, discoverFunc: v.isPKIRoleNotConfigured},
		{state: PolicyNotConfigured, discoverFunc: v.isPolicyNotConfigured},
		{state: TokenAuthNotConfigured, discoverFunc: v.isTokenAuthNotConfigured},

		//{state: AppRoleAuthNotEnabled, discoverFunc: v.isAppRoleAuthNotEnabled},

		//{state: IdentityBackendNotConfigured, discoverFunc: v.isIdentityBackendNotConfigured},
		//{state: IdentityKeyNotConfigured, discoverFunc: v.isIdentityKeyNotConfigured},
	}

	for _, v := range discovery {
		isInthisState, err := v.discoverFunc(ctx)
		if err != nil {
			core.GetLogger(ctx).Errorf("failure evaluating state %s: %v", vaultStateToString[v.state], err)
			return Unitialized, err
		}

		if isInthisState {
			return v.state, nil
		}
	}

	return Done, nil
}

func (v *VaultConfigController) reconcile(ctx context.Context, current, desired VaultState) error {
	log := core.GetLogger(ctx)
	stateMap := map[VaultState]ControllerStateHandler{
		Unitialized:            v.Initialize,
		Sealed:                 v.Unseal,
		PolicyNotConfigured:    v.ConfigurePolicy,
		PKINotEnabled:          v.EnablePKI,
		PKICACertNotConfigured: v.ConfigurePKICACert,
		PKIURLsNotConfigured:   v.ConfigurePKIURLs,
		PKIRoleNotConfigured:   v.ConfigurePKIRole,
		TokenAuthNotConfigured: v.ConfigureTokenAuth,
		// No further state transition needed,
		// previous function transitions to the "Done" state
	}

	var i int
	for {
		if current == desired {
			break
		}

		log.Debugf("Executing stateHandler for state, error: %v", vaultStateToString[current])

		stateHandler, ok := stateMap[current]
		if !ok {
			return fmt.Errorf("No state handler for state %s",
				vaultStateToString[current])
		}

		next, err := stateHandler(ctx)
		if err != nil {
			return fmt.Errorf("StateHandler for state %s exited with error %v",
				vaultStateToString[current],
				err)
		}

		current = next
		i++

		if i >= maxStateChangesPerCycle {
			return fmt.Errorf("maximum state change per cycle reached." +
				" This is likely to happen if the state machine has a sink state or is cyclic")
		}
	}

	return nil
}

func (v *VaultConfigController) desiredPKICACertConfig(bundle string) *vaultPKICACertConfig {
	return &vaultPKICACertConfig{
		PEMBundle: bundle,
	}
}

func (v *VaultConfigController) desiredPKIURLsConfig() *vaultPKIURLsConfig {
	return &vaultPKIURLsConfig{
		IssuingCertificates:   []string{fmt.Sprintf("%s/v1/pki/ca", v.svcURL)},
		CRLDistributionPoints: []string{fmt.Sprintf("%s/v1/pki/crl", v.svcURL)},
	}
}

func (v *VaultConfigController) desiredPKIRoleConfig() *vaultPKIRoleConfig {
	return &vaultPKIRoleConfig{
		AllowAnyName:      true,
		Organization:      []string{PKIRoleName},
		MaxTTL:            desiredPKICertTTL,
		NotBeforeDuration: desiredNotBeforeDuration,
	}
}
