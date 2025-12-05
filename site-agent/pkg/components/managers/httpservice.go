// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package managers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	computils "github.com/nvidia/carbide-rest/site-agent/pkg/components/utils"
	"go.temporal.io/sdk/client"
)

func handleSiteStatusRequest(w http.ResponseWriter, r *http.Request) {
	// Get the status of Bootstrap n write to the HTTP response body.
	siteStatus := ManagerAccess.API.Bootstrap.GetState()
	for _, v := range siteStatus {
		fmt.Fprint(w, v)
	}
	siteStatus = ManagerAccess.API.Orchestrator.GetState()
	for _, v := range siteStatus {
		fmt.Fprint(w, v)
	}
	siteStatus = ManagerAccess.API.Carbide.GetState()
	for _, v := range siteStatus {
		fmt.Fprint(w, v)
	}
	fmt.Fprint(w, fmt.Sprintln(" Site Agent Health: ",
		computils.CompStatus(ManagerAccess.Data.EB.HealthStatus.Load()).String()))
}

func handleVpcStatusRequest(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("GET params were: %v", r.URL.Query())

	vpcName := r.URL.Query().Get(computils.ParamName)
	if vpcName != "" {
		workflowID := "vpc-get-" + vpcName
		log.Info().Msgf("VPC GET : %v", vpcName)

		workflowOptions := client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: ManagerAccess.Conf.EB.Temporal.TemporalSubscribeQueue,
		}

		we, err := ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber.ExecuteWorkflow(
			context.Background(),
			workflowOptions,
			ManagerAccess.API.VPC.GetVPCByName,
		)
		if err != nil {
			log.Info().Msgf("Op error: %v", err.Error())
			fmt.Fprint(w, err.Error())
			return
		}
		ResourceResponse := &wflows.GetVPCResponse{}
		we.Get(context.Background(), ResourceResponse)
		fmt.Fprint(w, ResourceResponse.Status)
		fmt.Fprint(w, ResourceResponse.StatusMsg)
		for _, v := range ResourceResponse.List.Vpcs {
			fmt.Fprint(w, v)
		}
		return
	}

	// Get the status of VPC n write to the HTTP response body.
	vpcStatus := ManagerAccess.API.VPC.GetState()
	for _, v := range vpcStatus {
		fmt.Fprint(w, v)
	}
}

func handleSubnetStatusRequest(w http.ResponseWriter, r *http.Request) {
	// Get the status of Subnet n write to the HTTP response body.
	sStatus := ManagerAccess.API.Subnet.GetState()
	for _, v := range sStatus {
		fmt.Fprint(w, v)
	}
}

func handleInstanceStatusRequest(w http.ResponseWriter, r *http.Request) {
	// Get the status of Instance n write to the HTTP response body.
	sStatus := ManagerAccess.API.Instance.GetState()
	for _, v := range sStatus {
		fmt.Fprint(w, v)
	}
}

func handleMachineStatusRequest(w http.ResponseWriter, r *http.Request) {
	// Get the status of Instance n write to the HTTP response body.
	sStatus := ManagerAccess.API.Machine.GetState()
	for _, v := range sStatus {
		fmt.Fprint(w, v)
	}
}

// StartHTTPServer - start a web server on the specified port.
func StartHTTPServer() {
	port := os.Getenv("ESA_PORT")
	http.HandleFunc(computils.SiteStatus, handleSiteStatusRequest)
	http.HandleFunc(computils.VPCStatus, handleVpcStatusRequest)
	http.HandleFunc(computils.SubnetStatus, handleSubnetStatusRequest)
	http.HandleFunc(computils.InstanceStatus, handleInstanceStatusRequest)
	http.HandleFunc(computils.MachineStatus, handleMachineStatusRequest)
	go http.ListenAndServe(fmt.Sprintf("localhost:%v", port), nil)
}
