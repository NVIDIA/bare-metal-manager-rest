package health

import (
	"errors"
	"time"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/site-agent/pkg/conftypes"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// HealthWorkflow struct to hold the Temporal clients and configuration
type HealthWorkflow struct {
	tcPublish   client.Client
	tcSubscribe client.Client
	cfg         *conftypes.Config
}

const (
	// RetryInterval denotes the initial interval for the retry policy
	RetryInterval = 2
	// RetryCount denotes the maximum number of retries
	RetryCount = 10
	// MaxTemporalActivityRetryCount denotes the maximum number of retries for the Temporal activity
	MaxTemporalActivityRetryCount = 3
)

// NewHealthWorkflows creates an instance for HealthWorkflows
func NewHealthWorkflows(tcPublish client.Client, tcSubscribe client.Client, cfg *conftypes.Config) HealthWorkflow {
	return HealthWorkflow{
		tcPublish:   tcPublish,
		tcSubscribe: tcSubscribe,
		cfg:         cfg,
	}
}

// GetHealth Workflow to get the health status
func GetHealth(ctx workflow.Context, TransactionID *wflows.TransactionID) (*wflows.HealthStatus, error) {
	logger := workflow.GetLogger(ctx)

	slogger := log.With(logger, "Workflow", "CreateHealthWorkflow", "ResourceRequest", TransactionID)
	slogger.Info("Health: Starting  the Health Workflow")

	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Health: Starting  the Health Workflow")

	var status wflows.HealthStatus

	// Validations
	if TransactionID == nil {
		slogger.Error("Health: TransactionID is nil")
		ManagerAccess.Data.EB.Log.Error().Msg("Health: TransactionID is nil")
		return nil, errors.New("Health: TransactionID is nil")
	}
	if TransactionID.ResourceId == "" {
		slogger.Error("Health: TransactionID.ResourceId is empty")
		ManagerAccess.Data.EB.Log.Error().Msg("Health: TransactionID.ResourceId is empty")
		return nil, errors.New("Health: TransactionID.ResourceId is empty")
	}

	// Use default retry interval
	RetryInterval := 1 * time.Second

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    RetryInterval,
		BackoffCoefficient: 2.0,
		MaximumInterval:    1 * time.Minute,
		MaximumAttempts:    MaxTemporalActivityRetryCount,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 20 * time.Second,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	healthWorkflow := HealthWorkflow{}

	err := workflow.ExecuteActivity(ctx, healthWorkflow.GetHealthActivity).Get(ctx, &status)
	if err != nil {
		slogger.Error("Health: Failed to get Health workflow", "Error", err)
		ManagerAccess.Data.EB.Log.Error().Interface("Error", err).Msg("Health: Failed to get Health")
		return &status, err
	}

	slogger.Info("Health: Successfully updated Health")
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Health: Successfully updated Health")

	return &status, err
}
