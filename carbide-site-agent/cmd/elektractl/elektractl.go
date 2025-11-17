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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	wflows "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	zlog "github.com/rs/zerolog/log"
	"go.temporal.io/sdk/worker"
)

var (
	// TemporalPublishQueue - Temporal Publish Queue
	TemporalPublishQueue = "forge-publish"
	// TemporalSubscribeQueue - Temporal Subscribe Queue
	TemporalSubscribeQueue = "forge-subscribe"
	// DefaultNamespace - Default Namespace
	DefaultNamespace = "default"
	// DefaultVPCName - Default VPC Name
	DefaultVPCName = "testVPC"
	// DefaultVPCOrg - Default VPC Org
	DefaultVPCOrg = uuid.New().String()
	// ElektractlConfig - Elektractl Config path
	ElektractlConfig = "/usr/local/bin/elektractl.env"
)

func main() {
	// Update ElektraConfig from the env if present
	config, isPresent := os.LookupEnv("ELEKTRACTL_CONFIG")
	if isPresent {
		ElektractlConfig = config
	}
	err := godotenv.Load(ElektractlConfig)
	if err != nil {
		log.Fatalln("Unable to create elektractl client ", err)
	}
	Execute()
}

func epilogue(we client.WorkflowRun, wflow string) {
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	// Synchronously wait for the workflow completion.
	var result string
	err := we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get ", wflow, " workflow result", err)
	}
	log.Println(wflow, "Workflow result:", result)
}

func vpcCreate(name, org string) {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()

	vpcReq := &wflows.CreateVPCRequest{
		Name:                 name,
		TenantOrganizationId: org,
	}
	transaction := &wflows.TransactionID{
		ResourceId: "testVPC-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	wflow := "CreateVPC"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, vpcReq)
	if err != nil {
		log.Fatalln("Unable execute workflow CreateVPCWorkflow result", err)
	}
	epilogue(we, wflow)
}
func vpcUpdate(name, org, id string) {

	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()
	transaction := &wflows.TransactionID{
		ResourceId: "testVPC-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	// Update the workflow here
	vpcUpReq := &wflows.UpdateVPCRequest{
		Name: name,
		Id: &wflows.UUID{
			Value: id,
		},
		TenantOrganizationId: org,
	}
	wflow := "UpdateVPC"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, vpcUpReq)
	if err != nil {
		log.Fatalln("Unable execute workflow UpdateVPCWorkflow result", err)
	}
	epilogue(we, wflow)
}
func vpcDelete(id string) {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()
	transaction := &wflows.TransactionID{
		ResourceId: "testVPC-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	deleteReq := &wflows.DeleteVPCRequest{
		Id: &wflows.UUID{
			Value: id,
		},
	}
	wflow := "DeleteVPC"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, deleteReq)
	if err != nil {
		log.Fatalln("Unable execute workflow DeleteVPC result", err)
	}
	epilogue(we, wflow)
}

func subnetCreate() {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()

	createRequest := &wflows.CreateSubnetRequest{
		Name: "testSubnet",
		VpcId: &wflows.UUID{
			Value: "47a657d6-1fd0-497b-aaa0-b2072887fd17",
		},
		NetworkPrefixes: []*wflows.NetworkPrefixInfo{
			&wflows.NetworkPrefixInfo{
				Prefix:       "10.3.3.3/16",
				ReserveFirst: 2,
			},
		},
	}
	transaction := &wflows.TransactionID{
		ResourceId: "testSubnet-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	wflow := "CreateSubnet"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, createRequest)
	if err != nil {
		log.Fatalln("Unable execute workflow CreateSubnet result", err)
	}
	epilogue(we, wflow)
}
func instanceCreate() {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()
	createRequest := &wflows.CreateInstanceRequest{
		MachineId: &wflows.MachineId{
			Id: uuid.NewString(),
		},
	}
	transaction := &wflows.TransactionID{
		ResourceId: "testInstance-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	wflow := "CreateInstance"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, createRequest)
	if err != nil {
		log.Fatalln("Unable execute workflow CreateInstance result", err)
	}
	epilogue(we, wflow)
}
func instanceDelete(id string) {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c, workflowOptions := temporalclient.Publisher, temporalclient.PublishWflowOpts
	defer c.Close()
	transaction := &wflows.TransactionID{
		ResourceId: "testInstance-" + uuid.New().String(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	resourceID := id
	if id == "" {
		resourceID = transaction.ResourceId
	}
	deleteReq := &wflows.DeleteInstanceRequest{
		InstanceId: &wflows.UUID{
			Value: resourceID,
		},
	}
	wflow := "DeleteInstance"
	we, err := c.ExecuteWorkflow(context.Background(), *workflowOptions, wflow, transaction, deleteReq)
	if err != nil {
		log.Fatalln("Unable execute workflow DeleteInstance result", err)
	}
	epilogue(we, wflow)
}
func proxy() {
	temporalclient, err := createTemporalClient()
	if err != nil {
		log.Fatalln("cannot create temporal client", err)

	}
	c := temporalclient.Publisher
	defer c.Close()

	zlog.Info().Msg("Running as cloud proxy forever")
	// Start the Cloud workflow registration
	go RegisterCloudWorkflows(c)

	// sleep
	// Wait forever
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termChan:
		return
	}
}

// RegisterCloudWorkflows - Register Cloud Workflows
func RegisterCloudWorkflows(tclient client.Client) {

	zlog.Info().Msg("Registering Cloud workflows")

	// Create a temporal worker
	cworker := worker.New(
		tclient,
		TemporalPublishQueue,
		worker.Options{},
	)

	// Register Inventory worfklow
	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "UpdateMachineInventory",
	}
	cworker.RegisterWorkflowWithOptions(
		UpdateMachineInventory, wflowRegisterOptions,
	)

	// Register vpc workflow
	wflowRegisterOptions = workflow.RegisterOptions{
		Name: "UpdateVpcInfo",
	}
	cworker.RegisterWorkflowWithOptions(
		UpdateVpcInfo, wflowRegisterOptions,
	)

	// Register subnet workflow
	wflowRegisterOptions = workflow.RegisterOptions{
		Name: "UpdateSubnetInfo",
	}
	cworker.RegisterWorkflowWithOptions(
		UpdateSubnetInfo, wflowRegisterOptions,
	)

	// Register Instance workflow
	wflowRegisterOptions = workflow.RegisterOptions{
		Name: "UpdateInstanceInfo",
	}
	cworker.RegisterWorkflowWithOptions(
		UpdateInstanceInfo, wflowRegisterOptions,
	)

	// Start listening to the Task Queue
	zlog.Info().Msg("Starting Cloud worker")
	err := cworker.Run(worker.InterruptCh())
	if err != nil {
		zlog.Panic().Err(err).Msg("Failed to start cloud proxy worker")
	}

}

// UpdateMachineInventory - Update Machine Inventory
func UpdateMachineInventory(ctx workflow.Context, SiteID string, Inventory *wflows.MachineInventory) (err error) {
	log.Println("Received Inventory update for site", SiteID)
	log.Println("UpdateMachineInventory", SiteID, Inventory)
	return
}

// UpdateVpcInfo - Update Vpc Info
func UpdateVpcInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, VPCInfo *wflows.VPCInfo) (err error) {
	log.Println("Received VPC update for site", SiteID)
	log.Println("UpdateVpcInfo", SiteID, TransactionID, VPCInfo)
	return
}

// UpdateSubnetInfo - Update Subnet Info
func UpdateSubnetInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, SubnetInfo *wflows.SubnetInfo) (err error) {
	log.Println("Received Subnet update for site", SiteID)
	log.Println("UpdateSubnetInfo", SiteID, TransactionID, SubnetInfo)
	return
}

// UpdateInstanceInfo - Update Instance Info
func UpdateInstanceInfo(ctx workflow.Context, SiteID string, TransactionID *wflows.TransactionID, InstanceInfo *wflows.InstanceInfo) (err error) {
	log.Println("Received Instance update for site", SiteID)
	log.Println("UpdateInstanceInfo", SiteID, TransactionID, InstanceInfo)
	return
}
