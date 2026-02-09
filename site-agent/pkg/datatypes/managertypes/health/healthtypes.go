package healthtypes

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// We will define our own state later
// type HealthState int

// const (
// 	UNKNOWN HealthState = iota
// 	UP
// 	DOWN
// 	ERROR
// )

type SiteInventoryHealth struct {
	State     wflows.HealthState
	StatusMsg string
	// More fields to be added later
}

type SiteControllerConnection struct {
	State     wflows.HealthState
	StatusMsg string
	// More fields to be added later
}

type HighAvailability struct {
	State     wflows.HealthState
	StatusMsg string
}

// HealthCache Site Agent HealthCache
type HealthCache struct {
	Inventory        SiteInventoryHealth
	CarbideInterface SiteControllerConnection
	Availabilty      HighAvailability
}

// NewHealthCache - initialize site agent health cache
func NewHealthCache() *HealthCache {
	return &HealthCache{}
}
