package instancetype

import (
	"github.com/google/uuid"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the InstanceType Workflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here

	// Collect and Publish InstanceType Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverInstanceTypeInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered the Discover InstanceType Inventory workflow")

	inventoryManager := swa.NewManageInstanceTypeInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverInstanceTypeInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered the Discover InstanceType Inventory activity")

	api.RegisterCron()

	return nil
}
