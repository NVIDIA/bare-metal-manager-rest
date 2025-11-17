/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package conftypes

import (
	"encoding/json"
	"time"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/grpc/client"
)

// RunInEnvironment provides a strongly-typed indicator for the environment
// in which the app is running.
type RunInEnvironment int

const (
	// RunningInUnknown - Running In Unknown Env
	RunningInUnknown RunInEnvironment = iota
	// RunningInDocker - Running In Docker
	RunningInDocker
	// RunningInK8s - Running In K8s
	RunningInK8s
)

// TemporalConfig holds configurations for connecting to Temporal server
type TemporalConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	ClusterID                  string `json:"clusterID"`
	TemporalServer             string `json:"temporalServer"`
	TemporalPublishNamespace   string `json:"temporalPublishNamespace"`
	TemporalSubscribeNamespace string `json:"temporalSubscribeNamespace"`
	TemporalPublishQueue       string `json:"temporalPublishQueue"`
	TemporalSubscribeQueue     string `json:"temporalSubscribeQueue"`
	TemporalInventorySchedule  string `json:"temporalInventorySchedule"`
	TemporalCertPath           string `json:"temporalCertPath"`
}

// GetTemporalCertOTPFullPath - Get Temporal Cert OTP path
func (tc *TemporalConfig) GetTemporalCertOTPFullPath() string {
	return tc.TemporalCertPath + "/otp"
}

// GetTemporalCACertFilePath - Get Temporal CA Cert File n Path
func (tc *TemporalConfig) GetTemporalCACertFilePath() (string, string) {
	file := "ca.crt"
	return file, tc.TemporalCertPath + "/ca/"
}

// GetTemporalCACertFullPath - Get Temporal CA Cert Full Path
func (tc *TemporalConfig) GetTemporalCACertFullPath() string {
	file := "tls.crt"
	return tc.TemporalCertPath + "/ca/" + file
}

// GetTemporalClientCertFilePath - Get Temporal client Cert File n Path
func (tc *TemporalConfig) GetTemporalClientCertFilePath() ([]string, string) {
	file := []string{"tls.crt", "tls.key"}
	return file, tc.TemporalCertPath + "/client/"
}

// GetTemporalClientCertFullPath - Get Temporal client Cert Full Path
func (tc *TemporalConfig) GetTemporalClientCertFullPath() string {
	file := "tls.crt"
	return tc.TemporalCertPath + "/client/" + file
}

// GetTemporalClientKeyFullPath - Get Temporal client Key Full Path
func (tc *TemporalConfig) GetTemporalClientKeyFullPath() string {
	file := "tls.key"
	return tc.TemporalCertPath + "/client/" + file
}

// CarbideConfig holds configurations for connecting to Carbide server
type CarbideConfig struct {
	Address        string               `json:"carbideAddress"`
	Secure         client.SecureOptions `json:"carbideSecureOptions"`
	SkipServerAuth bool                 `json:"carbideSkipServerAuth"`
	ServerCAPath   string               `json:"carbideCertPath"`
	ClientCertPath string               `json:"carbideClientCertPath"`
	ClientKeyPath  string               `json:"carbideClientKeyPath"`
}

// DBConfig is the Elektra data store
type DBConfig struct {
	Server     string `json:"server"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	CACertPath string `json:"cacertpath"`
}

// Config for Elektra
type Config struct {
	Temporal         TemporalConfig
	Carbide          CarbideConfig
	DB               DBConfig
	IsMasterPod      bool          `json:"isMasterPod"`
	EnableDebug      bool          `json:"enableDebug"`
	DevMode          bool          `json:"devMode"`
	EnableTLS        bool          `json:"enableTLS"`
	DisableBootstrap bool          `json:"disableBootstrap"`
	BootstrapSecret  string        `json:"bootstrapSecret"`
	WatcherInterval  time.Duration `json:"watcherInterval"`
	PodNamespace     string        `json:"podNamespace"`
	TemporalSecret   string        `json:"temporalSecret"`
	MetricsPort      string        `json:"metricsPort"`
	SiteVersion      string        `json:"siteVersion"`
	CloudVersion     string        `json:"cloudVersion"`
	RunningIn        RunInEnvironment
	UtMode           bool
}

// String - json string
func (c *Config) String() string {
	str, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(str)
}

// NewConfType - new config
func NewConfType() *Config {
	// We can set the default config here
	return &Config{
		Temporal: TemporalConfig{},
		Carbide:  CarbideConfig{},
		DB:       DBConfig{},
	}
}
