/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	zlogadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"go.temporal.io/sdk/client"
	temporalClient "go.temporal.io/sdk/client"
	temporalWorker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/internal/config"
	cwfh "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/health"

	cwssaws "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
)

const (
	// ZerologMessageFieldName specifies the field name for log message
	ZerologMessageFieldName = "msg"
	// ZerologLevelFieldName specifies the field name for log level
	ZerologLevelFieldName = "type"
)

var (
	SiteID              = ""
	TemporalLocalClient temporalClient.Client
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = ZerologLevelFieldName
	zerolog.MessageFieldName = ZerologMessageFieldName

	cfg := config.NewConfig()
	defer cfg.Close()

	// Initializer Temporal client
	// Create the client object just once per process
	log.Info().Msg("creating Temporal client")

	var tc temporalClient.Client
	tLogger := logur.LoggerToKV(zlogadapter.New(zerolog.New(os.Stderr)))

	tcfg, err := cfg.GetTemporalConfig()
	if err != nil {
		log.Panic().Err(err).Msg("unable to get temporal config")
	}

	tc, err = temporalClient.NewLazyClient(temporalClient.Options{
		HostPort:  fmt.Sprintf("%v:%v", tcfg.Host, tcfg.Port),
		Namespace: tcfg.Namespace,
		ConnectionOptions: temporalClient.ConnectionOptions{
			TLS: tcfg.ClientTLSCfg,
		},
		Logger: tLogger,
	})

	if err != nil {
		log.Panic().Err(err).Str("Namespace", tcfg.Namespace).Msg("failed to create Temporal client")
	} else {
		defer tc.Close()
	}

	stc, err := temporalClient.Dial(temporalClient.Options{
		HostPort:  fmt.Sprintf("%v:%v", tcfg.Host, tcfg.Port),
		Namespace: "site",
		ConnectionOptions: temporalClient.ConnectionOptions{
			TLS: tcfg.ClientTLSCfg,
		},
		Logger: tLogger,
	})
	if err != nil {
		log.Panic().Err(err).Str("Namespace", "site").Msg("failed to create Temporal client")
	} else {
		defer tc.Close()
	}

	w := temporalWorker.New(tc, tcfg.Queue, temporalWorker.Options{})

	log.Info().Str("Temporal Namespace", tcfg.Namespace).Msg("registering workflow and activities")

	// Register site agent workflows
	w.RegisterWorkflow(CreateVPC)
	w.RegisterWorkflow(DeleteVPC)

	w.RegisterWorkflow(CreateSubnet)
	w.RegisterWorkflow(DeleteSubnet)

	w.RegisterWorkflow(CreateInstance)
	w.RegisterWorkflow(GetInstance)
	w.RegisterWorkflow(DeleteInstance)
	w.RegisterWorkflow(RebootInstance)

	w.RegisterWorkflow(GetHealth)

	w.RegisterWorkflow(CreateSSHKeyGroup)
	w.RegisterWorkflow(UpdateSSHKeyGroup)
	w.RegisterWorkflow(DeleteSSHKeyGroup)

	// Serve health endpoint
	go func() {
		log.Info().Msg("starting health check API server")
		http.HandleFunc("/healthz", cwfh.StatusHandler)
		http.HandleFunc("/readyz", cwfh.StatusHandler)

		hostPort := ":8899"

		serr := http.ListenAndServe(hostPort, nil)
		if serr != nil {
			log.Panic().Err(serr).Msg("failed to start health check server")
		}
	}()

	// PublishDummyMachineInventory one time publish
	// here tcfg.Namespace == site_id
	SiteID = tcfg.Namespace
	TemporalLocalClient = stc
	go PublishDummyMachineInventory(context.Background(), stc, tcfg.Namespace)

	// Start listening to the Task Queue
	log.Info().Str("Temporal Namespace", tcfg.Namespace).Msg("starting Temporal worker")
	err = w.Run(temporalWorker.InterruptCh())
	if err != nil {
		log.Panic().Err(err).Str("Temporal Namespace", tcfg.Namespace).Msg("failed to start site agent worker")
	}
}

func CreateVPC(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.CreateVPCRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest", ResourceRequest.Name).Msg("Received CreateVPC")
	return nil
}

func DeleteVPC(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.CreateVPCRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest", ResourceRequest.Name).Msg("Received DeleteVPC")
	return nil
}

func CreateSubnet(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.CreateSubnetRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest", ResourceRequest.Name).Msg("Received CreateSubnet")
	return nil
}

func DeleteSubnet(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.DeleteSubnetRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received DeleteSubnet")
	return nil
}

func CreateInstance(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.CreateInstanceRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received CreateInstance")
	if ResourceRequest != nil {
		if ResourceRequest.Metadata != nil {
			log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received CreateInstance with Labels")
		}
	}
	return nil
}

func GetInstance(ctx workflow.Context, TransactionID *cwssaws.TransactionID, InstanceID *cwssaws.UUID) (InstanceInfo *cwssaws.InstanceInfo, err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received GetInstance")
	return nil, nil
}

func DeleteInstance(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.DeleteInstanceRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received DeleteInstance")
	return nil
}

func RebootInstance(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.RebootInstanceRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received RebootInstance")
	return nil
}

func CreateSSHKeyGroup(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.CreateSSHKeyGroupRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest.KeysetId", ResourceRequest.KeysetId).Str("ResourceRequest.Version", ResourceRequest.Version).Msg("Received CreateSSHKeyGroup")

	TransactionID.Timestamp = timestamppb.Now()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "UpdateSSHKeyGroupInfo-" + ResourceRequest.KeysetId + "-" + ResourceRequest.Version,
		TaskQueue: "site",
	}

	sshKeyGroupInfo := &cwssaws.SSHKeyGroupInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
		StatusMsg: "SSH Key Group was successfully synced",
		TenantKeyset: &cwssaws.TenantKeyset{
			Version: ResourceRequest.Version,
		},
		ObjectStatus: cwssaws.ObjectStatus_OBJECT_STATUS_CREATED,
	}

	we, err := TemporalLocalClient.ExecuteWorkflow(context.Background(), workflowOptions, "UpdateSSHKeyGroupInfo", SiteID, TransactionID, sshKeyGroupInfo)
	if err != nil {
		return err
	}

	_ = we.GetID()
	return nil
}

func UpdateSSHKeyGroup(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.UpdateSSHKeyGroupRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest.KeysetId", ResourceRequest.KeysetId).Str("ResourceRequest.Version", ResourceRequest.Version).Msg("Received UpdateSSHKeyGroup")
	TransactionID.Timestamp = timestamppb.Now()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "UpdateSSHKeyGroupInfo-" + ResourceRequest.KeysetId + "-" + ResourceRequest.KeysetId,
		TaskQueue: "site",
	}

	sshKeyGroupInfo := &cwssaws.SSHKeyGroupInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
		StatusMsg: "SSH Key Group was successfully synced",
		TenantKeyset: &cwssaws.TenantKeyset{
			Version: ResourceRequest.Version,
		},
		ObjectStatus: cwssaws.ObjectStatus_OBJECT_STATUS_CREATED,
	}

	we, err := TemporalLocalClient.ExecuteWorkflow(context.Background(), workflowOptions, "UpdateSSHKeyGroupInfo", SiteID, TransactionID, sshKeyGroupInfo)
	if err != nil {
		return err
	}

	_ = we.GetID()
	return nil
}

func DeleteSSHKeyGroup(ctx workflow.Context, TransactionID *cwssaws.TransactionID, ResourceRequest *cwssaws.DeleteSSHKeyGroupRequest) (err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Str("ResourceRequest.KeysetId", ResourceRequest.KeysetId).Msg("Received DeleteSSHKeyGroup")
	return nil
}

func GetHealth(ctx workflow.Context, TransactionID *cwssaws.TransactionID) (HealthStatus *cwssaws.HealthStatus, err error) {
	log.Info().Str("TransactionID", TransactionID.ResourceId).Msg("Received GetHealth")
	HealthStatus = &cwssaws.HealthStatus{
		SiteInventoryCollection: &cwssaws.HealthStatusMsg{
			State: 1,
		},
		SiteControllerConnection: &cwssaws.HealthStatusMsg{
			State: 1,
		},
		SiteAgentHighAvailability: &cwssaws.HealthStatusMsg{
			State: 1,
		},
	}
	return HealthStatus, nil
}

// PublishDummyMachineInventory - Publish Inventory to the Forge cloud
func PublishDummyMachineInventory(ctx context.Context, tc temporalClient.Client, siteID string) (ID string, err error) {
	// Create the Machine
	log.Info().Str("Site ID", siteID).Msg("PublishDummyMachineInventory: Starting the Publish Dummy Machine")

	workflowOptions := client.StartWorkflowOptions{
		ID:        "MachineInventoryUpdate-" + siteID,
		TaskQueue: "site",
	}

	Inventory := &cwssaws.MachineInventory{
		Machines:  []*cwssaws.MachineInfo{},
		Timestamp: timestamppb.Now(),
	}

	// dummy machine ids
	ids := []string{"bd65600d-8669-4903-8a14-af88203add38", "7ee19487-e9f2-4a53-b4f0-26b7865f2f37", "bc75354e-098c-4ef3-82f5-8285f712909ds"}
	for idx, id := range ids {
		Inventory.Machines = append(Inventory.Machines, &cwssaws.MachineInfo{})
		Inventory.Machines[idx].Machine = &cwssaws.Machine{
			Id:    &cwssaws.MachineId{Id: id},
			State: "ready",
		}
	}

	we, err := tc.ExecuteWorkflow(context.Background(), workflowOptions, "UpdateMachineInventory", siteID, Inventory)
	if err != nil {
		return "", err
	}

	wid := we.GetID()
	return wid, nil
}
