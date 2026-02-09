package vpcprefix

import (
	"github.com/google/uuid"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the VpcPrefix Workflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here

	// Collect and Publish VpcPrefix Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverVpcPrefixInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered the Discover VpcPrefix Inventory workflow")

	inventoryManager := swa.NewManageVpcPrefixInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverVpcPrefixInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered the Discover VpcPrefix Inventory activity")

	api.RegisterCron()

	return nil
}
