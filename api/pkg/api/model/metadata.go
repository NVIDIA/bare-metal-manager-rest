package model

import (
	"github.com/nvidia/carbide-rest/api/pkg/metadata"
)

// APIMetadata is a data structure to capture Forge API system information
type APIMetadata struct {
	// Version contains the API version
	Version string `json:"version"`
	// BuildTime contains the time the binary was built
	BuildTime string `json:"buildTime"`
}

// NewAPIMetadata creates and returns a new APISystemInfo object
func NewAPIMetadata() *APIMetadata {
	amd := &APIMetadata{
		Version:   metadata.Version,
		BuildTime: metadata.BuildTime,
	}

	return amd
}
