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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gorilla/mux"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

const (
	testToken       = "testRootToken"
	testKey         = "testRootKey"
	testCMToken     = "testCertToken"
	testRefreshRate = 10 * time.Millisecond
	//testRefreshRate = 2*time.Second
)

var (
	testCACert = `-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIRAOSueOOLdD79GWtlYv0Uc9MwDQYJKoZIhvcNAQELBQAw
QjEPMA0GA1UEChMGbnZpZGlhMS8wLQYDVQQDEyZldHMtbnZpZGlhLW1jYW1wZGV2
LmRldi5lZ3gubnZpZGlhLmNvbTAeFw0yMTA5MTQyMDI0MjFaFw0zMTA5MTIyMDI0
MjFaMEIxDzANBgNVBAoTBm52aWRpYTEvMC0GA1UEAxMmZXRzLW52aWRpYS1tY2Ft
cGRldi5kZXYuZWd4Lm52aWRpYS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDVCsUYWHk/qan/7MVaG5sH039jKITm2I5mRpRuprjVrNa1YNWZYHi0
wAgLQPM+ftl54tq1svzHMWUNvOHDZZmnZR72KsFquZItz8PjQb/xZSbG30JS/eJF
JgFvbjASec1gihfKR+x7rviYfPnMwWEmo1VQ/3su9XgOUb9D7Ce8RrGpmPmpV+P8
X031rDjDGVr8gHm/sCBQCZEwgUlswut+6HKasZ1gbS2U4ffZG3crOkEuMMC9E2QU
KXqAmJY9qt8p5DEyQh7/nm/y8Qb0GCeDn65XHDHtHdkPMFUHSG74/blHbRngrVcg
Iv9f7H+TABs3PSaQ43FuNGWZGr/8R6evAgMBAAGjgZUwgZIwDgYDVR0PAQH/BAQD
AgGmMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAPBgNVHRMBAf8EBTAD
AQH/MB0GA1UdDgQWBBToTQsj9vDA1GCCFlzScBoi73+jtDAxBgNVHREEKjAogiZl
dHMtbnZpZGlhLW1jYW1wZGV2LmRldi5lZ3gubnZpZGlhLmNvbTANBgkqhkiG9w0B
AQsFAAOCAQEAwAieHlp1xFseZHIZRcEYmQWhXZtYOA3Qvt3HjXYhZD54UcOOCOD8
gnt/dLNEgcPIoxxk2taxaM8CkiBzQdzTHJxNOWXyHOTf0vko8pidhot0W/UnL0uF
+nY43yKM+LrNYefJTSZBjbcgmkCEDZPR9pHDFLP+DrU4p1SA9fl8eYnHPOF6pYco
LgdhMG8nyL+iawksoIKJorpmmuIyc4qYj5OaNDR3CvJsfy1i4FDiWgnBqZRwZozC
R2i+pzoinIcG0X4K3G1EGEZbJf9tvH8wc07t/lejUYXgtKRKulw9kD1KJeQzzs1X
b4/fjPLgIVzZefZWCG5CAVvFGEvTdMuPkA==
-----END CERTIFICATE-----`
)

type suite struct {
	l           net.Listener
	srv         *httptest.Server
	forceErr    bool
	initialized bool
	sealed      bool
	policies    map[string]string
	tokens      map[string]*vault.TokenCreateRequest
	urlsCfg     map[string]interface{}
	rolesCfg    map[string]interface{}
}

func (s *suite) setup(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	s.l = l
	s.sealed = true
	s.policies = make(map[string]string)
	s.tokens = make(map[string]*vault.TokenCreateRequest)
	s.urlsCfg = make(map[string]interface{})
	s.rolesCfg = map[string]interface{}{
		"allow_any_name":      true,
		"organization":        []string{PKIRoleName},
		"max_ttl":             desiredPKICertTTL,
		"not_before_duration": desiredNotBeforeDuration,
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/v1/sys/init", func(w http.ResponseWriter, r *http.Request) {
		var c []byte
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodGet {
			isr := &vault.InitStatusResponse{Initialized: s.initialized}
			c, err = json.Marshal(isr)
			require.NoError(t, err)
		} else {
			s.initialized = true
			ir := &vault.InitResponse{
				RootToken: testToken,
				Keys:      []string{testKey},
			}
			c, err = json.Marshal(ir)
			require.NoError(t, err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/sys/unseal", func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
		s.sealed = false
		ssr := &vault.SealStatusResponse{
			Initialized: s.initialized,
			Sealed:      s.sealed,
		}
		c, err := json.Marshal(ssr)
		require.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/sys/health", func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
		hr := &vault.HealthResponse{
			Initialized: s.initialized,
			Sealed:      s.sealed,
		}
		c, err := json.Marshal(hr)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/sys/policies/acl/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		req := make(map[string]string)
		if r.Method == http.MethodPut {
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			s.policies[name] = req["policy"]
			return
		}

		resp := &vault.Secret{
			Data: map[string]interface{}{
				"policy": s.policies[name],
			},
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/auth/token/create", func(w http.ResponseWriter, r *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
		req := &vault.TokenCreateRequest{}
		err = json.NewDecoder(r.Body).Decode(req)
		require.NoError(t, err)
		s.tokens[testCMToken] = req
		pol := s.policies[CertManagerPolicyName]
		resp := &vault.Secret{
			Auth: &vault.SecretAuth{
				ClientToken:   testCMToken,
				Policies:      []string{pol},
				TokenPolicies: []string{pol},
			},
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/auth/token/lookup", func(w http.ResponseWriter, r *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		req := make(map[string]string)
		err = json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		token := req["token"]
		if token != "" && s.tokens[token] != nil {
			pol := s.policies[CertManagerPolicyName]
			resp := &vault.Secret{
				Auth: &vault.SecretAuth{
					ClientToken:   testCMToken,
					Policies:      []string{pol},
					TokenPolicies: []string{pol},
				},
			}
			c, err := json.Marshal(resp)
			require.NoError(t, err)

			w.Header().Set("Content-Type", "application/json")
			w.Write(c)
			return
		}
		http.Error(w, "notfound", http.StatusInternalServerError)
	})
	rtr.HandleFunc("/v1/sys/mounts", func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		resp := &vault.Secret{
			Data: map[string]interface{}{
				"pki/": &vault.MountOutput{},
			},
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.PathPrefix("/v1/sys/mounts/").HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
	})
	rtr.HandleFunc("/v1/pki/cert/ca", func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		resp := &vault.Secret{
			Data: map[string]interface{}{
				"certificate": testCACert,
			},
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/pki/config/ca", func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}

		resp := &vault.Secret{
			Data: map[string]interface{}{
				"certificate": testCACert,
			},
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})
	rtr.HandleFunc("/v1/pki/config/urls", func(w http.ResponseWriter, r *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
		if r.Method == http.MethodPut {
			err := json.NewDecoder(r.Body).Decode(&s.urlsCfg)
			require.NoError(t, err)
		}
		resp := &vault.Secret{
			Data: s.urlsCfg,
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})

	rtr.HandleFunc("/v1/pki/roles/"+PKIRoleName, func(w http.ResponseWriter, _ *http.Request) {
		if s.forceErr {
			http.Error(w, "forced error", http.StatusInternalServerError)
			return
		}
		resp := &vault.Secret{
			Data: s.rolesCfg,
		}
		c, err := json.Marshal(resp)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	})

	s.srv = httptest.NewUnstartedServer(rtr)
	s.srv.Listener = l
	s.srv.Start()
}

func (s *suite) teardown() {
	s.srv.Close()
}

func TestVaultConfig_NewController(t *testing.T) {
	vcc, err := NewController("https://localhost", "https://ets.example.com", "testdata/secrets")
	assert.NoError(t, err)
	assert.NotNil(t, vcc)
	assert.NotNil(t, vcc.client)
	assert.Equal(t, "https://localhost", vcc.client.Address())
	assert.Equal(t, "https://localhost", vcc.vaultURL)
	assert.Equal(t, refreshRateDefault, vcc.refreshRate)
	assert.Equal(t, "testdata/secrets", vcc.secretsRootPath)
}

func TestVaultFSM(t *testing.T) {
	ts := &suite{}
	ts.setup(t)
	defer ts.teardown()
	// make a copy of testdata
	dir, err := ioutil.TempDir("", "secrets")
	require.NoError(t, err)
	defer os.RemoveAll(dir) // clean up
	cmd := exec.Command("cp", "-r", "testdata/secrets", dir)
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	vcc, err := NewController(fmt.Sprintf("http://%s", ts.l.Addr().String()), "https://test.example.com", dir+"/secrets")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(core.NewDefaultContext(context.Background()))

	vcc.refreshRate = testRefreshRate
	vcc.Start(ctx)
	select {
	case _, ok := <-vcc.TokenChan():
		assert.Equal(t, true, ok)
	case <-time.After(10 * time.Second):
		t.Error("timeout")
	}
	cancel()
	<-vcc.TokenChan()

	vcc, err = NewController(fmt.Sprintf("http://%s", ts.l.Addr().String()), "https://test.example.com", dir+"/secrets")
	require.NoError(t, err)
	ts.urlsCfg = make(map[string]interface{})
	ts.rolesCfg = make(map[string]interface{})
	status, err := vcc.isPKIURLsNotConfigured(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, true, status)
	status, err = vcc.isPKIRoleNotConfigured(context.Background())
	assert.Equal(t, nil, err)
	assert.Equal(t, true, status)

	ts.forceErr = true
	discoverers := []ControllerStateDiscovery{
		vcc.isUninitialized,
		vcc.isSealed,
		vcc.isPolicyNotConfigured,
		vcc.isTokenAuthNotConfigured,
		vcc.isPKINotEnabled,
		vcc.isPKICACertNotConfigured,
		vcc.isPKIURLsNotConfigured,
		vcc.isPKIRoleNotConfigured,
	}
	for _, d := range discoverers {
		_, err = d(context.Background())
		assert.NotEqual(t, nil, err)
	}

	handlers := []ControllerStateHandler{
		vcc.Initialize,
		vcc.Unseal,
		vcc.ConfigurePolicy,
		vcc.EnablePKI,
		vcc.ConfigurePKICACert,
		vcc.ConfigurePKIURLs,
		vcc.ConfigurePKIRole,
		vcc.ConfigureTokenAuth,
	}

	for _, h := range handlers {
		_, err = h(context.Background())
		assert.NotEqual(t, nil, err)
	}
}
