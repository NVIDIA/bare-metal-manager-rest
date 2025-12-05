// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package elektra

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"

	"github.com/rs/zerolog/log"
	"github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/components/managers/bootstrap"
	"github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/components/managers/carbide"
	"github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/datatypes/elektratypes"

	computils "github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/components/utils"
	bootstraptypes "github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/datatypes/managertypes/bootstrap"
	workflowtypes "github.com/nvidia/carbide-rest-api/carbide-rest-site-agent/pkg/datatypes/managertypes/workflow"
)

var (
	// NOTE: These values must match values in test Carbide server in elektra-carbide-lib

	// DefaultTestVpcID is the default VPC ID for testing
	DefaultTestVpcID = "00000000-0000-4000-8000-000000000000"
	// DefaultTestNetworkSegmentID is the default NetworkSegment ID for testing
	DefaultTestNetworkSegmentID = "00000000-0000-4000-9000-000000000000"
	// DefaultTestTenantKeysetID is the default TenantKeyset ID for testing
	DefaultTestTenantKeysetID = "00000000-0000-4000-a000-000000000000"
	// DefaultTestIBParitionID is the default IBPartition ID for testing
	DefaultTestIBParitionID = "00000000-0000-4000-b000-000000000000"

	// Counters to track various failure/success counts
	wflowGrpcFail = 0
	wflowGrpcSucc = 0
	wflowActFail  = uint64(0)
	wflowActSucc  = uint64(0)
	wflowPubFail  = uint64(0)
	wflowPubSucc  = uint64(0)
)

// Test Elektra objects
var testElektra *Elektra
var testElektraTypes *elektratypes.Elektra

func checkGrpcState(stats *workflowtypes.MgrState) {
	fail := int(carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcFail.Load())
	if wflowGrpcFail != fail {
		log.Info().Msgf("wflowGrpcFail: %v, state fail: %v ", wflowGrpcFail, fail)
		panic("wflowGrpcFail ctr incorrect")
	}
	succ := int(carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcSucc.Load())
	if wflowGrpcSucc != succ {
		log.Info().Msgf("wflowGrpcSucc: %v, state succ %v", wflowGrpcSucc, succ)
		panic("wflowGrpcSucc ctr incorrect")
	}
	state := uint64(carbide.ManagerAccess.Data.EB.Managers.Carbide.State.HealthStatus.Load())
	if uint64(computils.CompHealthy) != state {
		log.Info().Msgf("state %v ", state)
		panic("Component not in Healthy State")
	}

	if stats.WflowActFail.Load() != wflowActFail {
		log.Info().Msgf("%v != %v", stats.WflowActFail.Load(), wflowActFail)
		panic("wflowActFail")
	}
	if stats.WflowActSucc.Load() != wflowActSucc {
		log.Info().Msgf("%v != %v", stats.WflowActSucc.Load(), wflowActSucc)
		panic("wflowActSucc")
	}
	if stats.WflowPubSucc.Load() != wflowPubSucc {
		log.Info().Msgf("%v != %v", stats.WflowPubSucc.Load(), wflowPubSucc)
		panic("wflowPubSucc")
	}
	if stats.WflowPubFail.Load() != wflowPubFail {
		log.Info().Msgf("%v != %v", stats.WflowPubFail.Load(), wflowPubFail)
		panic("wflowPubFail")
	}
}

// SetupTestCA generates a test CA certificate and key.
func SetupTestCA(t *testing.T) (string, string) {
	caKeyPath := "/tmp/ca.key"
	caCertPath := "/tmp/ca.crt"

	// Generate CA private key
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Create a CA certificate template
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Self-sign the CA certificate
	caCertBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	assert.NoError(t, err)

	// Write the CA private key to file
	caKeyFile, err := os.Create(caKeyPath)
	assert.NoError(t, err)
	defer caKeyFile.Close()
	pem.Encode(caKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey)})

	// Write the CA certificate to file
	caCertFile, err := os.Create(caCertPath)
	assert.NoError(t, err)
	defer caCertFile.Close()
	pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes})

	return caKeyPath, caCertPath
}

// SetupTestCerts generates a test server/client certificate and key signed by the provided CA.
func SetupTestCerts(t *testing.T, caCertPath, caKeyPath string) (string, string) {
	keyPath := "/tmp/tls.key"
	certPath := "/tmp/tls.crt"

	// Load CA
	caCertPEM, err := os.ReadFile(caCertPath)
	assert.NoError(t, err)
	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	assert.NoError(t, err)

	caKeyPEM, err := os.ReadFile(caKeyPath)
	assert.NoError(t, err)
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	assert.NoError(t, err)

	// Generate keypair for test server/client
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publickey := &privatekey.PublicKey

	// Create certificate template to be signed by CA
	certTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"Test Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		BasicConstraintsValid: true,
	}

	// Create the server/client certificate signed by CA
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, publickey, caPrivateKey)
	assert.NoError(t, err)

	// Encode and save the server/client private key
	keyFile, err := os.Create(keyPath)
	assert.NoError(t, err)
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)})

	// Encode and save the server/client certificate
	certFile, err := os.Create(certPath)
	assert.NoError(t, err)
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	return keyPath, certPath
}

// TestInitElektra initializes a test version of the Site Agent
func TestInitElektra(t *testing.T) {
	if testElektra != nil {
		return
	}
	os.Setenv("CARBIDE_CERT_CHECK_INTERVAL", "1") // set this to check if certs were rotated every second to help with unit tests
	defer os.Unsetenv("CARBIDE_CERT_CHECK_INTERVAL")

	// Initialize Elektra microservice
	log.Info().Msg("Elektra: Initializing Elektra microservice")

	// Generate CA and client certs
	caKeyPath, caCertPath := SetupTestCA(t)
	keyPath, certPath := SetupTestCerts(t, caCertPath, caKeyPath)

	// Set environment variables for the test
	if err := os.Setenv("CARBIDE_CA_CERT_PATH", caCertPath); err != nil {
		t.Fatalf("Failed to set CARBIDE_CA_CERT_PATH: %v", err)
	}
	if err := os.Setenv("CARBIDE_CLIENT_CERT_PATH", certPath); err != nil {
		t.Fatalf("Failed to set CARBIDE_CLIENT_CERT_PATH: %v", err)
	}
	if err := os.Setenv("CARBIDE_CLIENT_KEY_PATH", keyPath); err != nil {
		t.Fatalf("Failed to set CARBIDE_CLIENT_KEY_PATH: %v", err)
	}

	// Initialize Elektra Data Structures
	testElektraTypes = elektratypes.NewElektraTypes()
	// Initialize Elektra API
	api, initErr := NewElektraAPI(testElektraTypes, true)
	if initErr != nil {
		log.Fatal().Err(initErr).Msg("Elektra: Failed to initialize Elektra API")
	} else {
		log.Info().Msg("Elektra: Successfully initialized Elektra API")
	}
	err := simulateMountedSecretFile(bootstrap.ManagerAccess.Conf.EB.BootstrapSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Elektra: simulateMountedSecretFile")
	}
	// Initialize Elektra Managers
	api.Init()
	api.Start()
	testElektra = api
}

func simulateMountedSecretFile(dir string) error {
	result := &struct {
		Data bootstraptypes.SecretConfig `yaml:"data"`
	}{}
	cfgFile, err := os.ReadFile(dir + "bootstrapInfo")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(cfgFile, result)
	if err != nil {
		return err
	}
	bCfg := &result.Data

	// Return error if credentials are not available
	if bCfg.CACert == "" || bCfg.CredsURL == "" || bCfg.OTP == "" || bCfg.UUID == "" {
		return fmt.Errorf("Empty secret")
	}
	log.Info().Msgf("Read %v %v %v %v", bCfg.UUID, bCfg.OTP, bCfg.CredsURL, bCfg.CACert)

	err = os.WriteFile(dir+bootstraptypes.TagUUID, []byte(bCfg.UUID), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(dir+bootstraptypes.TagOTP, []byte(bCfg.OTP), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(dir+bootstraptypes.TagCACert, []byte(bCfg.CACert), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(dir+bootstraptypes.TagCredsURL, []byte(bCfg.CredsURL), 0644)
	if err != nil {
		return err
	}

	return nil
}
