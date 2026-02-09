package managertypes

import (
	bootstraptypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/bootstrap"
	carbidetypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/carbide"
	healthtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/health"
	rlatypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/rla"
	workflowtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/workflow"
)

// Managers - manager ds
type Managers struct {
	Version string
	// All the datastructures of Managers below
	Workflow  *workflowtypes.Workflow
	Carbide   *carbidetypes.Carbide
	RLA       *rlatypes.RLA
	Bootstrap *bootstraptypes.Bootstrap
	Health    *healthtypes.HealthCache
}

// NewManagerType - get new type of all managers
func NewManagerType() *Managers {
	return &Managers{
		Version: "0.0.1",
		// All the managers below
		Workflow:  workflowtypes.NewWorkflowInstance(),
		Carbide:   carbidetypes.NewCarbideInstance(),
		RLA:       rlatypes.NewRLAInstance(),
		Bootstrap: bootstraptypes.NewBootstrapInstance(),
		Health:    healthtypes.NewHealthCache(),
	}
}
