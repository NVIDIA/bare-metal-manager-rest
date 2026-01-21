// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vault

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/pkg/errors"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/pki"
	vaultcfg "github.com/nvidia/carbide-rest/cert-manager/pkg/vault/config"
	kapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	vaultSecretsRootPathDefault = "/vault/secrets"
	caBaseDNSDefault            = "temporal.nvidia.com"
)

// NewCommand creates a new cli command
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "Forge Credentials Service",
		Usage: "Forge Credentials Service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Log debug message to stderr",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "tls-port",
				Value: "8000",
				Usage: "TLS port to listen to",
			},
			&cli.StringFlag{
				Name:  "insecure-port",
				Value: "8001",
				Usage: "http port to listen to",
			},
			&cli.StringFlag{
				Name:  "dns-name",
				Value: "credsmgr.csm",
				Usage: "DNS name for incluster tls access",
			},
			&cli.StringFlag{
				Name:  "vault-endpoint",
				Value: "http://localhost:8200",
				Usage: "Vault service endpoint used by backend",
			},
			&cli.StringFlag{
				Name:  "vault-secrets-root-path",
				Value: vaultSecretsRootPathDefault,
				Usage: "The root path of the vault secrets and configuration written to disk",
			},
			&cli.StringFlag{
				Name:  "ca-base-dns",
				Value: caBaseDNSDefault,
				Usage: "Base dns appended to common names",
			},
		&cli.StringFlag{
			Name:    "vault-ingress-url",
			Usage:   "ingress URL for this vault/csm pod for setting up CA (required when not using --use-native-pki)",
			EnvVars: []string{"FORGE_VAULT_INGRESS_URL"},
		},
			&cli.StringFlag{
				Name:    "cert-manager-secret-name",
				Usage:   "name of secret to create for cert manager",
				EnvVars: []string{"FORGE_CERT_MANAGER_SECRET_NAME"},
			},
			&cli.StringFlag{
				Name:    "cert-manager-namespace",
				Usage:   "namespace of the cert-maanger",
				EnvVars: []string{"CERT_MANAGER_NS"},
			},
			&cli.StringFlag{
				Name:    "sentry-dsn",
				Value:   "",
				EnvVars: []string{"SENTRY_DSN"},
				Usage:   "DSN for sentry/glitchtip",
			},
			&cli.BoolFlag{
				Name:    "use-native-pki",
				Value:   false,
				EnvVars: []string{"USE_NATIVE_PKI"},
				Usage:   "Use native Go PKI instead of Vault for certificate generation",
			},
			&cli.StringFlag{
				Name:    "ca-common-name",
				Value:   "Carbide Local CA",
				EnvVars: []string{"CA_COMMON_NAME"},
				Usage:   "Common name for the CA certificate (only used with --use-native-pki)",
			},
			&cli.StringFlag{
				Name:    "ca-organization",
				Value:   "NVIDIA",
				EnvVars: []string{"CA_ORGANIZATION"},
				Usage:   "Organization for the CA certificate (only used with --use-native-pki)",
			},
			&cli.StringFlag{
				Name:    "ca-cert-file",
				Value:   "/vault/secrets/vault-root-ca-certificate/certificate",
				EnvVars: []string{"CA_CERT_FILE"},
				Usage:   "Primary path to CA certificate file (vault-style path)",
			},
			&cli.StringFlag{
				Name:    "ca-key-file",
				Value:   "/vault/secrets/vault-root-ca-private-key/privatekey",
				EnvVars: []string{"CA_KEY_FILE"},
				Usage:   "Primary path to CA private key file (vault-style path)",
			},
			&cli.StringFlag{
				Name:    "alt-ca-cert-file",
				Value:   "/etc/pki/ca/tls.crt",
				EnvVars: []string{"ALT_CA_CERT_FILE"},
				Usage:   "Alternate path to CA certificate file (for migration)",
			},
			&cli.StringFlag{
				Name:    "alt-ca-key-file",
				Value:   "/etc/pki/ca/tls.key",
				EnvVars: []string{"ALT_CA_KEY_FILE"},
				Usage:   "Alternate path to CA private key file (for migration)",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				core.GetLogger(c.Context).Logger.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			ctx := c.Context
			log := core.GetLogger(ctx)
			useNativePKI := c.Bool("use-native-pki")

			o := Options{
				Addr:         ":" + c.String("tls-port"),
				InsecureAddr: ":" + c.String("insecure-port"),
				DNSName:      c.String("dns-name"),
				CABaseDNS:    c.String("ca-base-dns"),
				sentryDSN:    c.String("sentry-dsn"),
			}

			var certIssuer CertificateIssuer

			if useNativePKI {
				// Use native Go PKI instead of Vault
				log.Info("Using native Go PKI for certificate generation")
				caCertFile := c.String("ca-cert-file")
				caKeyFile := c.String("ca-key-file")
				altCACertFile := c.String("alt-ca-cert-file")
				altCAKeyFile := c.String("alt-ca-key-file")
				log.Infof("CA paths - primary: %s, alternate: %s", caCertFile, altCACertFile)
				issuer, err := pki.NewNativeCertificateIssuer(pki.NativeCertificateIssuerOptions{
					BaseDNS:       c.String("ca-base-dns"),
					CACommonName:  c.String("ca-common-name"),
					CAOrganization: c.String("ca-organization"),
					CACertFile:    caCertFile,
					CAKeyFile:     caKeyFile,
					AltCACertFile: altCACertFile,
					AltCAKeyFile:  altCAKeyFile,
				})
				if err != nil {
					log.Errorf("Failed to create native PKI issuer: %v", err)
					return err
				}
				certIssuer = issuer
				log.Info("Native PKI issuer initialized successfully")
		} else {
			// Use Vault for certificate generation
			vaultIngressURL := c.String("vault-ingress-url")
			if vaultIngressURL == "" {
				return errors.New("--vault-ingress-url is required when not using --use-native-pki")
			}

			vaultEndpoint := c.String("vault-endpoint")
			o.VaultEndpoint = vaultEndpoint

			log.Infof("Configuring vault at %s", vaultEndpoint)
				vaultCfgCtrl, err := vaultcfg.NewController(vaultEndpoint,
					c.String("vault-ingress-url"),
					c.String("vault-secrets-root-path"))
				if err != nil {
					log.Errorf("Vault initialized failed, error: %v", err)
					return err
				}
				vaultCfgCtrl.Start(ctx)

				// Wait for vault token to be sent from the channel once vault is initialized
				log.Info("Waiting for vault to be initialized...")
				vaultToken, ok := <-vaultCfgCtrl.TokenChan()
				if !ok {
					return errors.New("vault failed to initialize")
				}
				o.VaultToken = vaultToken
				log.Info("Vault successfully initialized!")

				certMgrSecret := c.String("cert-manager-secret-name")
				if certMgrSecret != "" {
					createCertMgrSecret(ctx, c.String("cert-manager-namespace"), certMgrSecret, vaultCfgCtrl.CertManagerToken())
				}
			}

			log.Info("Configuring Certificate server")
			s, err := NewServerWithIssuer(ctx, o, certIssuer)
			if err != nil {
				return err
			}
			log.Info("Starting Certificate server")
			s.Start(ctx)

			<-ctx.Done()

			gracePeriod := 5 * time.Second
			log.Infof("Shut down requested, wait for %v grace period ...", gracePeriod)
			time.Sleep(gracePeriod)
			log.Infof("Server terminated.")
			return nil
		},
	}
}

func createCertMgrSecret(ctx context.Context, namespace, name, token string) {
	secretData := map[string][]byte{
		"token": []byte(token),
	}
	log := core.GetLogger(ctx)
	log.Infof("Creating secret %s/%s", namespace, name)
	if namespace == "" {
		log.Fatalf("namespace is required")
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "rest.InClusterConfig()"))
	}
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "kubernetes.NewForConfig"))
	}

	secObj, err := k8sClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		secObj.Data = secretData
		_, err = k8sClient.CoreV1().Secrets(namespace).Update(ctx, secObj, metav1.UpdateOptions{})
	} else {
		secObj = &kapi.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: secretData,
		}
		_, err = k8sClient.CoreV1().Secrets(namespace).Create(ctx, secObj, metav1.CreateOptions{})
	}

	if err != nil {
		log.Fatal(errors.Wrap(err, "secret create/update"))
	}

	log.Infof("Created secret %s/%s", namespace, name)
}
