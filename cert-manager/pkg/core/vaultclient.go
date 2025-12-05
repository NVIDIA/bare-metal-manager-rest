// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	vault "github.com/hashicorp/vault/api"
	otelattr "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	// k8sServiceAccountTokenPath is the K8s service account token path
	// within a pod.
	// #nosec G101
	k8sServiceAccountTokenPath = "/run/secrets/kubernetes.io/serviceaccount/token"
)

const (
	attrVaultRequestIDKey   = "fc.core.vault.req.id"
	attrVaultRequestPathKey = "fc.core.vault.req.path"
	// #nosec G101
	attrVaultTokenExpirationKey = "fc.core.vault.token.expiration"

	// DefaultVaultTokenExpiryMargin is the default margin for expiration of
	// the vault token. This value determines when the vault client
	// should consider the token expired, and go fetch a new one.
	DefaultVaultTokenExpiryMargin = 5 * time.Minute

	// DefaultVaultTokenExpiryMarginEnvKey is the environment variable name
	// for the env var that allows for configuring an expiry margin
	// of the vault token.
	// #nosec G101
	DefaultVaultTokenExpiryMarginEnvKey = "FC_VAULT_TOKEN_EXPIRY_MARGIN"
)

// VaultOptions defines vault options
type VaultOptions struct {
	// Addr is the vault service endpoint
	Addr string

	// VauleToken is the vault token used to authenticate to vault.
	VaultToken string

	// K8sRole is the K8s role name configured in vault's K8s auth
	// backend.
	K8sRole string

	// AppRole and AppRoleSecret are approle role-id and secret-id
	// configured in vault's approle auth backend.
	AppRole, AppRoleSecret string

	// VaultTokenExpiryMargin is the margin for expiration of
	// the vault token. This value determines when the vault client
	// should consider the token expired, and go fetch a new one.
	VaultTokenExpiryMargin time.Duration
}

// VaultClient defines a vault client
type VaultClient struct {
	deps              vaultDeps
	options           VaultOptions
	tokenLock         sync.Mutex
	tokenExpiration   time.Time
	tokenExpiryMargin time.Duration
}

// NewVaultClient returns a vault client
func NewVaultClient(ctx context.Context, o VaultOptions) (*VaultClient, error) {
	client, err := vault.NewClient(&vault.Config{Address: o.Addr})
	if err != nil {
		return nil, fmt.Errorf("failed to create a vault client with error: %v", err)
	}

	deps := &vaultDepsImpl{client: client}

	vc := &VaultClient{deps: deps, options: o, tokenExpiryMargin: DefaultVaultTokenExpiryMargin}

	// Check for expiry margin, and use if set
	if o.VaultTokenExpiryMargin != 0 {
		vc.tokenExpiryMargin = o.VaultTokenExpiryMargin
	}

	if err := vc.vaultLogin(ctx); err != nil {
		return nil, err
	}
	return vc, nil
}

var vaultReqSequenceID = uint64(0)

// nextVaultRequestID returns a 6 digits string as request ID
func nextVaultRequestID(_ context.Context) string {
	seqID := atomic.AddUint64(&vaultReqSequenceID, 1)
	return fmt.Sprintf("%06d", seqID%1000000)
}

func logicalRead(ctx context.Context, deps vaultDeps, path string) (*vault.Secret, error) {
	log := GetLogger(ctx).WithFields(map[string]interface{}{"rpc": "vault.LogicalRead", "path": path})
	reqID := nextVaultRequestID(ctx)
	reqLog := log.WithField("req", reqID)
	reqLog.Infof("")
	_, span := tracer.Start(ctx,
		"vault.LogicalRead",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(otelattr.String(attrVaultRequestIDKey, reqID),
			otelattr.String(attrVaultRequestPathKey, path)))
	defer span.End()

	resp, err := deps.LogicalRead(path)

	respLog := log.WithFields(map[string]interface{}{"resp": reqID, "err": err})
	respLog.Infof("")
	// Only log payload in debug mode
	respLog.WithField("payload", resp).Debugf("")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
	}
	return resp, err
}

// LogicalRead implements the read method
func (c *VaultClient) LogicalRead(ctx context.Context, path string) (*vault.Secret, error) {
	// Ensure we're logged into vault before write
	if err := c.vaultLogin(ctx); err != nil {
		return nil, err
	}

	return logicalRead(ctx, c.deps, path)
}

func logicalWrite(ctx context.Context, deps vaultDeps, path string, data map[string]interface{}) (*vault.Secret, error) {
	log := GetLogger(ctx).WithFields(map[string]interface{}{"rpc": "vault.LogicalWrite", "path": path})

	reqID := nextVaultRequestID(ctx)
	reqLog := log.WithField("req", reqID)
	reqLog.Infof("")
	reqLog.WithField("payload", data).Debugf("")
	_, span := tracer.Start(ctx,
		"vault.LogicalWrite",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(otelattr.String(attrVaultRequestIDKey, reqID),
			otelattr.String(attrVaultRequestPathKey, path)))
	defer span.End()

	resp, err := deps.LogicalWrite(path, data)

	respLog := log.WithFields(map[string]interface{}{"resp": reqID, "err": err})
	respLog.Infof("")
	respLog.WithField("payload", resp).Debugf("")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
	}

	return resp, err
}

// LogicalWrite implements the write method
func (c *VaultClient) LogicalWrite(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error) {
	// Ensure we're logged into vault before write
	if err := c.vaultLogin(ctx); err != nil {
		return nil, err
	}

	return logicalWrite(ctx, c.deps, path, data)
}

// vaultDeps declares the vault API calls we used in order to
// implement core.VaultClient as dependencies.
type vaultDeps interface {
	SetToken(token string)
	LogicalRead(path string) (*vault.Secret, error)
	LogicalWrite(path string, data map[string]interface{}) (*vault.Secret, error)
}

type vaultDepsImpl struct{ client *vault.Client }

func (v *vaultDepsImpl) SetToken(token string) {
	v.client.SetToken(token)
}

func (v *vaultDepsImpl) LogicalWrite(path string, data map[string]interface{}) (*vault.Secret, error) {
	return v.client.Logical().Write(path, data)
}

func (v *vaultDepsImpl) LogicalRead(path string) (*vault.Secret, error) {
	return v.client.Logical().Read(path)
}

func (c *VaultClient) vaultLogin(ctx context.Context) error {
	useVaultToken := c.options.VaultToken != ""
	useKubeAuth := c.options.K8sRole != ""
	useAppRoleAuth := c.options.AppRole != ""
	count := 0
	if useVaultToken {
		count++
	}
	if useKubeAuth {
		count++
	}
	if useAppRoleAuth {
		count++
	}
	if count != 1 {
		return fmt.Errorf("Exactly one of VaultToken, K8sRole and AppRole MUST be specified for vault login")
	}
	if useAppRoleAuth && c.options.AppRoleSecret == "" {
		return fmt.Errorf("AppRoleSecret is empty when login useAppRoleAuth")
	}

	c.tokenLock.Lock()
	defer c.tokenLock.Unlock()

	// If the token is set to zero. Check if the token has expired. If it is not expired we can use the existing token.
	if !c.tokenExpiration.IsZero() && time.Now().Before(c.tokenExpiration.Add(-c.tokenExpiryMargin)) {
		GetLogger(ctx).WithField("useVaultToken", useVaultToken).
			WithField("tokenExpiration", c.tokenExpiration).
			Debug("Skipping login attempt due to token being unexpired or vault token in use.")
		return nil
	}

	// Returns token, tokenExpiration and error if one occurs
	var login func(ctx context.Context, deps vaultDeps, o VaultOptions) (string, time.Time, error)
	if useVaultToken {
		login = vaultLoginWithToken
	} else if useKubeAuth {
		login = vaultLoginWithKubeAuth
	} else if useAppRoleAuth {
		login = vaultLoginWithAppRole
	}
	token, tokenExpiresAt, err := login(ctx, c.deps, c.options)
	if err != nil {
		return err
	}
	// Set the token and its respective expiration for the next run
	c.deps.SetToken(token)
	c.tokenExpiration = tokenExpiresAt
	return nil
}

func vaultLoginWithToken(ctx context.Context, _ vaultDeps, o VaultOptions) (string, time.Time, error) {
	log := GetLogger(ctx).WithField("span", "vaultLoginWithToken")
	log.Debugf("begin")
	defer log.Debugf("end")
	_, span := tracer.Start(ctx,
		"vaultLoginWithToken",
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()
	//span.SetAttributes(otelattr.Any(attrVaultTokenExpirationKey, time.Time{}))
	return o.VaultToken, time.Time{}, nil
}

func vaultLoginWithKubeAuth(ctx context.Context, deps vaultDeps, o VaultOptions) (string, time.Time, error) {
	log := GetLogger(ctx).WithField("span", "vaultLoginWithKubeAuth")
	log.Debugf("begin")
	defer log.Debugf("end")
	ctx, span := tracer.Start(ctx,
		"vaultLoginWithKubeAuth",
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	saToken, err := ioutil.ReadFile(k8sServiceAccountTokenPath)
	if err != nil {
		return "", time.Time{}, err
	}

	tNow := time.Now()
	resp, err := logicalWrite(ctx, deps, "auth/kubernetes/login", map[string]interface{}{
		"role": o.K8sRole,
		"jwt":  string(saToken),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	if resp == nil {
		err = fmt.Errorf("resp is nil from auth/kubernetes/login")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	if resp.Auth == nil {
		err = fmt.Errorf("resp.Auth is nil from auth/kubernetes/login")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	tokenExpiration := tNow.Add(time.Duration(resp.Auth.LeaseDuration) * time.Second)
	//span.SetAttributes(otelattr.Any(attrVaultTokenExpirationKey, tokenExpiration))
	span.SetStatus(codes.Ok, "K8s login successful.")
	return resp.Auth.ClientToken, tokenExpiration, nil
}

func vaultLoginWithAppRole(ctx context.Context, deps vaultDeps, o VaultOptions) (string, time.Time, error) {
	log := GetLogger(ctx).WithField("span", "vaultLoginWithAppRole")
	log.Debugf("begin")
	defer log.Debugf("end")
	ctx, span := tracer.Start(ctx,
		"vaultLoginWithAppRole",
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	tNow := time.Now()
	resp, err := logicalWrite(ctx, deps, "auth/approle/login", map[string]interface{}{
		"role_id":   o.AppRole,
		"secret_id": o.AppRoleSecret,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	if resp == nil {
		err = fmt.Errorf("resp is nil from auth/approle/login")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	if resp.Auth == nil {
		err = fmt.Errorf("resp.Auth is nil from auth/approle/login")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(otelattr.Bool("error", true))
		return "", time.Time{}, err
	}
	tokenExpiration := tNow.Add(time.Duration(resp.Auth.LeaseDuration) * time.Second)
	//span.SetAttributes(otelattr.Any(attrVaultTokenExpirationKey, tokenExpiration))
	span.SetStatus(codes.Ok, "K8s approle login successful.")
	return resp.Auth.ClientToken, tokenExpiration, nil
}
