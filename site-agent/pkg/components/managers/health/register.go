package health

import (
	"go.temporal.io/sdk/activity"
	workflow "go.temporal.io/sdk/workflow"
)

// RegisterSubscriber registers the HealthWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("Health: Registering the subscribers")

	// Get Health workflow interface
	Healthinterface := NewHealthWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	// Register worfklow
	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "GetHealth",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(
		GetHealth, wflowRegisterOptions,
	)

	ManagerAccess.Data.EB.Log.Info().Msg("Health: successfully registered the get Health workflow")

	// Register activity
	activityRegisterOptions := activity.RegisterOptions{
		Name: "CreateHealthActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		Healthinterface.GetHealthActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("Health: successfully registered the get Health activity")

	return nil
}
