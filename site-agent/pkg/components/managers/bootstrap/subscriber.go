package bootstrap

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// RegisterSubscriber registers the Bootstrap workflows and activities with the Temporal client
func (api *BoostrapAPI) RegisterSubscriber() error {
	// Initialize logger
	logger := ManagerAccess.Data.EB.Log

	// Only master pod should watch for the OTP rotation workflow
	if !ManagerAccess.Conf.EB.IsMasterPod {
		return nil
	}

	// Register the workflows
	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "RotateTemporalCertAccessOTP",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(api.RotateTemporalCertAccessOTP, wflowRegisterOptions)
	logger.Info().Msg("Bootstrap: successfully registered the ReceiveAndProcessOTP workflow")

	// Register the activities
	otpHandler := NewOTPHandler(ManagerAccess.Data.EB.Managers.Bootstrap.Secret)

	activityRegisterOptions := activity.RegisterOptions{
		Name: "ReceiveAndSaveOTP",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(otpHandler.ReceiveAndSaveOTP, activityRegisterOptions)
	logger.Info().Msg("Bootstrap: successfully registered the ReceiveAndSaveOTP activity")

	return nil
}
