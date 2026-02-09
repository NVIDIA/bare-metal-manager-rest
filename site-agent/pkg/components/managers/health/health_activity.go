package health

import (
	"time"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// GetHealthActivity - get the health status
func (ac *HealthWorkflow) GetHealthActivity() (*wflows.HealthStatus, error) {
	status := &wflows.HealthStatus{
		Timestamp: timestamppb.New(time.Now()),
		SiteInventoryCollection: &wflows.HealthStatusMsg{
			State: ManagerAccess.Data.EB.Managers.Health.Inventory.State,
		},
		SiteControllerConnection: &wflows.HealthStatusMsg{
			State: ManagerAccess.Data.EB.Managers.Health.CarbideInterface.State,
		},
		SiteAgentHighAvailability: &wflows.HealthStatusMsg{
			State: ManagerAccess.Data.EB.Managers.Health.Availabilty.State,
		},
	}

	return status, nil
}
