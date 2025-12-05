// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package config

import (
	"context"
	"sync"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// VaultConfigController defines a controller
type VaultConfigController struct {
	tokenCh         chan string
	client          vaultClient
	vaultURL        string
	refreshRate     time.Duration
	secretsRootPath string
	svcURL          string
	creds           *vault.InitResponse
	credsMutex      sync.RWMutex
	certMgrToken    string
}

// Note this module expects the following environment variable to be filled
// with the pod namespace
const (
	SecretTokenPath      = "token"
	SecretTokenName      = "vault-token"
	CertManagerTokenName = "certmgr-token"
	// #nosec G101
	SecretCACertPath = "vault-root-ca-certificate"
	// #nosec G101
	SecretCAPrivateKeyPath   = "vault-root-ca-private-key"
	PolicyName               = "root-identity"
	AppRolePolicyName        = "approle-policy"
	CertManagerPolicyName    = "certmgr-policy"
	IdentityKeyName          = "apiserver"
	IdentityRoleName         = "apiserver"
	PKIRoleName              = "cloud-cert"
	desiredRotationPeriod    = 3600 * 24 // vault default
	desiredVerificationTTL   = 3600 * 24 // vault default
	desiredPKICertTTL        = 3600 * 24 * 180
	desiredPKICACertTTL      = "87600h"
	desiredNotBeforeDuration = 120 // 2 min slack
	certManagerTokenPeriod   = "24h"
)

// Definitions of vault config
const (
	RootIdentityPolicy = `path "identity/*" {
		capabilities = ["create", "read", "update", "delete", "list"]
	}
	path "auth/approle/role/*" {
		capabilities = ["create", "read", "update", "delete", "list"]
	}
	path "pki/issue/*" {
		capabilities = ["create", "read", "update"]
	}`

	AppRolePolicy = `path "identity/oidc/token/*" {
		capabilities = ["read"]
	}`

	CertManagerPolicy = `path "pki/sign/*" {
		capabilities = ["create", "read", "update"]
	}`
)

// VaultState defines the fsm state
type VaultState int

// State definitions for vault
const (
	Unitialized VaultState = iota
	Sealed
	PolicyNotConfigured

	AppRoleAuthNotEnabled

	IdentityBackendNotConfigured
	IdentityKeyNotConfigured

	PKINotEnabled
	PKICACertNotConfigured
	PKIURLsNotConfigured
	PKIRoleNotConfigured
	TokenAuthNotConfigured

	Done
)

var (
	vaultNextState = map[VaultState]VaultState{
		Unitialized:            Sealed,
		Sealed:                 PKINotEnabled,
		PKINotEnabled:          PKICACertNotConfigured,
		PKICACertNotConfigured: PKIURLsNotConfigured,
		PKIURLsNotConfigured:   PKIRoleNotConfigured,
		PKIRoleNotConfigured:   PolicyNotConfigured,
		PolicyNotConfigured:    TokenAuthNotConfigured,
		TokenAuthNotConfigured: Done,
		//		AppRoleAuthNotEnabled:        IdentityBackendNotConfigured,
		//		IdentityBackendNotConfigured: IdentityKeyNotConfigured,
		//		IdentityKeyNotConfigured:     PKINotEnabled,
	}
	vaultStateToString = map[VaultState]string{
		Unitialized:                  "Unitialized",
		Sealed:                       "Sealed",
		PolicyNotConfigured:          "PolicyNotConfigured",
		AppRoleAuthNotEnabled:        "AppRoleAuthNotEnabled",
		IdentityBackendNotConfigured: "IdentityBackendNotConfigured",
		IdentityKeyNotConfigured:     "IdentityKeyNotConfigured",
		PKINotEnabled:                "PKINotEnabled",
		PKICACertNotConfigured:       "PKICACertNotConfigured",
		PKIURLsNotConfigured:         "PKIURLsNotConfigured",
		PKIRoleNotConfigured:         "PKIRoleNotConfigured",
		TokenAuthNotConfigured:       "TokenAuthNotConfigured",
		Done:                         "Done",
	}
)

// ControllerStateHandler is the state handler type
type ControllerStateHandler func(ctx context.Context) (VaultState, error)

// ControllerStateDiscovery is the discovery handler type
type ControllerStateDiscovery func(ctx context.Context) (bool, error)

type vaultIdentityBackendConfig struct {
	Issuer string `json:"issuer"`
}

type vaultIdentityKey struct {
	Algorithm        string      `json:"algorithm"`
	RotationPeriod   interface{} `json:"rotation_period"`
	AllowedClientIDs []string    `json:"allowed_client_ids"`
	VerificationTTL  interface{} `json:"verification_ttl"`
}

type vaultPKICACertConfig struct {
	PEMBundle string `json:"pem_bundle"`
}

type vaultPKIRoleConfig struct {
	AllowAnyName      bool          `json:"allow_any_name"`
	Organization      []string      `json:"organization"`
	MaxTTL            time.Duration `json:"max_ttl"`
	NotBeforeDuration time.Duration `json:"not_before_duration"`
}

type vaultPKIURLsConfig struct {
	IssuingCertificates   []string `json:"issuing_certificates"`
	CRLDistributionPoints []string `json:"crl_distribution_points"`
}

type vaultClient interface {
	Sys() vaultClientSys
	Logical() vaultClientLogical
	AuthToken() vaultClientToken
	Address() string
	SetToken(token string)
}

type vaultClientSys interface {
	Health() (*vault.HealthResponse, error)

	Init(opts *vault.InitRequest) (*vault.InitResponse, error)
	Unseal(shard string) (*vault.SealStatusResponse, error)

	GetPolicy(name string) (string, error)
	PutPolicy(name, rules string) error

	ListAuth() (map[string]*vault.AuthMount, error)
	EnableAuth(path, authType, desc string) error

	ListMounts() (map[string]*vault.MountOutput, error)
	Mount(path string, mountInfo *vault.MountInput) error
	TuneMount(path string, config vault.MountConfigInput) error
}

type vaultClientToken interface {
	Create(opts *vault.TokenCreateRequest) (*vault.Secret, error)
	Lookup(token string) (*vault.Secret, error)
	RenewWithContext(context.Context, string, int) (*vault.Secret, error)
}

type vaultClientLogical interface {
	Read(path string) (*vault.Secret, error)
	Write(path string, data map[string]interface{}) (*vault.Secret, error)
}

type vaultClientInternal struct {
	client *vault.Client
}

func (vci *vaultClientInternal) Sys() vaultClientSys {
	return vci.client.Sys()
}

func (vci *vaultClientInternal) AuthToken() vaultClientToken {
	return vci.client.Auth().Token()
}

func (vci *vaultClientInternal) Logical() vaultClientLogical {
	return vci.client.Logical()
}

func (vci *vaultClientInternal) Address() string {
	return vci.client.Address()
}

func (vci *vaultClientInternal) SetToken(token string) {
	vci.client.SetToken(token)
}
