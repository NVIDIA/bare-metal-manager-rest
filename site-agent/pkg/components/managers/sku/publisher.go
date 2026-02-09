package sku

import (
	"github.com/google/uuid"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the SKU workflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: Registering the publishers")

	// Collect and Publish SKU Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverSkuInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: successfully registered the DiscoverSkuInventory workflow")

	inventoryManager := swa.NewManageSkuInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverSkuInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: successfully registered the DiscoverSkuInventory activity")

	_ = api.RegisterCron()

	return nil
}
