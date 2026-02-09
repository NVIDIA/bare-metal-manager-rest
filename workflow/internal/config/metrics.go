package config

import (
	"fmt"
)

// MetricsConfig holds configuration of Metrics
type MetricsConfig struct {
	Enabled bool
	Port    int
}

// GetListenAddr returns the local address for listen socket.
func (mcfg *MetricsConfig) GetListenAddr() string {
	return fmt.Sprintf(":%v", mcfg.Port)
}

// NewMetricsConfig initializes and returns a configuration object for managing Metrics
func NewMetricsConfig(enabled bool, port int) *MetricsConfig {
	return &MetricsConfig{
		Enabled: enabled,
		Port:    port,
	}
}
