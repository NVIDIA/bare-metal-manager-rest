package rla

import (
	"context"

	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"
)

// CreateGRPCClientActivity - Create GRPC client Activity
func (RLA *API) CreateGrpcClientActivity(ctx context.Context, ResourceID string) (client *client.RlaClient, err error) {
	// Create the VPC
	ManagerAccess.Data.EB.Log.Info().Interface("Request", ResourceID).Msg("RLA: Starting  the gRPC connection Activity")

	// Use temporal logger for temporal logs
	logger := activity.GetLogger(ctx)
	withLogger := log.With(logger, "Activity", "CreateGrpcClientActivity", "ResourceReq", ResourceID)
	withLogger.Info("RLA: Starting  the gRPC connection Activity")

	// Create the client
	ManagerAccess.Data.EB.Log.Info().Interface("Request", ResourceID).Msg("RLA: Creating  grpc client")

	err = RLA.CreateGrpcClient()
	if err != nil {
		return nil, err
	}
	return RLA.GetGrpcClient(), nil
}

// RegisterGRPC - Register GRPC
func (RLA *API) RegisterGRPC() {
	// Register activity
	activityRegisterOptions := activity.RegisterOptions{
		Name: "CreateRlaGrpcClientActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		ManagerAccess.API.RLA.CreateGrpcClientActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GRPC client activity")
}
