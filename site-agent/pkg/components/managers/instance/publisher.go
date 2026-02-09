package instance

import (
	"github.com/google/uuid"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	"go.temporal.io/sdk/activity"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
)

// RegisterPublisher registers the InstanceWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {

	// Register the publishers here
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: Registering the publishers")

	// Get Instance workflow interface
	Instanceinterface := NewInstanceWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)
	activityRegisterOptions := activity.RegisterOptions{
		Name: "PublishInstanceActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		Instanceinterface.PublishInstanceActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered the Publish Instance activity")

	activityRegisterOptions = activity.RegisterOptions{
		Name: "PublishInstancePowerStatus",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		Instanceinterface.PublishInstancePowerStatus, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered the PublishInstancePowerStatus activity")

	// Instance Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverInstanceInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered the Discover Instance Inventory workflow")

	instanceInventoryManager := swa.NewManageInstanceInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceInventoryManager.DiscoverInstanceInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered the Discover Instance Inventory activity")

	api.RegisterCron()
	return nil
}
