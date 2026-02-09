package expectedmachine

import (
	"github.com/google/uuid"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the ExpectedMachineWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here

	// Collect and Publish ExpectedMachine Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverExpectedMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DiscoverExpectedMachineInventory workflow")

	inventoryManager := swa.NewManageExpectedMachineInventory(
		uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		ManagerAccess.Data.EB.Managers.Carbide.Client,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		InventoryCarbidePageSize,
	)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverExpectedMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DiscoverExpectedMachineInventory activity")

	_ = api.RegisterCron()

	return nil
}
