package instance

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	workflow "go.temporal.io/sdk/workflow"
)

// RegisterSubscriber registers the InstanceWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: Registering the subscribers")

	// Get Instance workflow interface
	Instanceinterface := NewInstanceWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	instanceManager := swa.NewManageInstance(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register workflows

	// Sync workflows
	// Register CreateInstance worfklow
	// TODO: Once all Site Agents are updated, remove the legacy CreateInstance workflow, duplicate register Site Workflow as CreateInstance
	// Once all Site Agents are updated with duplicate workflow, switch Cloud API to use call CreateInstance, then de-register CreateInstanceV2 workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateInstanceV2)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Create Instance v2 workflow")

	// Register CreateInstances workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateInstances)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Create Instances workflow")

	// Register DeleteInstanceV2 worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteInstanceV2)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Delete Instance v2 workflow")

	// Register UpdateInstance workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateInstance)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Update Instance workflow")

	// Register RebootInstance workflow
	// TODO: Same as above
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(sww.RebootInstance, workflow.RegisterOptions{
		Name: "RebootInstanceV2",
	})
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Reboot Instance v2 workflow")

	// Legacy workflows
	// Register CreateInstance worfklow (deprecated)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.Instance.CreateInstance)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered deprecated Create Instance workflow")

	// Register DeleteInstance worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.Instance.DeleteInstance)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Delete Instance workflow")

	// Register RebootInstance worfklow (deprecated)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.Instance.RebootInstance)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Reboot Instance workflow")

	// Register activities

	// Sync workflow activities
	// Register CreateInstanceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceManager.CreateInstanceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Create Instance activity")

	// Register CreateInstancesOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceManager.CreateInstancesOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Create Instances activity")

	// Register DeleteInstanceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceManager.DeleteInstanceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Delete Instance activity")

	// Register UpdateInstanceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceManager.UpdateInstanceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Update Instance activity")

	// Register RebootInstanceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceManager.RebootInstanceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Reboot Instance activity")

	// Legacy workflow activities
	// Register CreateInstanceActivity (deprecated)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(Instanceinterface.CreateInstanceActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered deprecated Create Instance activity")

	// Register DeleteInstanceActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(Instanceinterface.DeleteInstanceActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered delete Instance activity")

	// Register RebootInstanceActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(Instanceinterface.RebootInstanceActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: successfully registered Reboot Instance activity")

	return nil
}

// RegisterSubscribers - this is method 2 of registering the subscriber
func RegisterSubscribers() {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("Instance: Registering the subscribers")
	ManagerAccess.API.Orchestrator.AddWorkflow(ManagerAccess.API.Instance.CreateInstance)
}
