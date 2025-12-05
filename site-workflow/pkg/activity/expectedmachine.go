// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package activity

import (
	"context"
	"errors"
	"time"

	"github.com/gogo/status"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	swe "github.com/nvidia/carbide-rest/site-workflow/pkg/error"
	cclient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
)

// ManageExpectedMachineInventory is an activity wrapper for Expected Machine inventory collection and publishing
type ManageExpectedMachineInventory struct {
	config ManageInventoryConfig
}

// DiscoverExpectedMachineInventory is an activity to collect Expected Machine inventory and publish to Temporal queue
func (memi *ManageExpectedMachineInventory) DiscoverExpectedMachineInventory(ctx context.Context) error {
	logger := log.With().Str("Activity", "DiscoverExpectedMachineInventory").Logger()
	logger.Info().Msg("Starting activity")
	inventoryImpl := manageInventoryImpl[string, *cwssaws.ExpectedMachine, *cwssaws.ExpectedMachineInventory]{
		itemType:               "ExpectedMachine",
		config:                 memi.config,
		internalFindIDs:        expectedMachineFindIDs,
		internalFindByIDs:      expectedMachineFindByIDs,
		internalPagedInventory: expectedMachinePagedInventory,
		internalFindFallback:   expectedMachineFindFallback,
	}
	return inventoryImpl.CollectAndPublishInventory(ctx, &logger)
}

// NewManageExpectedMachineInventory returns a ManageInventory implementation for Expected Machine activity
func NewManageExpectedMachineInventory(config ManageInventoryConfig) ManageExpectedMachineInventory {
	return ManageExpectedMachineInventory{
		config: config,
	}
}

func expectedMachineFindIDs(_ context.Context, _ *cclient.CarbideClient) ([]string, error) {
	// TODO: Re-implement this function when pagination support is added in Carbide

	// For now, to trigger fallback return a gRPC Not Implemented error
	return nil, status.Errorf(codes.Unimplemented, "method FindExpectedMachineIDs is not implemented")
}

func expectedMachineFindByIDs(_ context.Context, _ *cclient.CarbideClient, _ []string) ([]*cwssaws.ExpectedMachine, error) {
	// TODO: Implement this function when pagination is support is added in Carbide
	return []*cwssaws.ExpectedMachine{}, nil
}

func expectedMachinePagedInventory(allItemIDs []string, pagedItems []*cwssaws.ExpectedMachine, input *pagedInventoryInput) *cwssaws.ExpectedMachineInventory {
	itemIDs := allItemIDs

	// Create an inventory page with the subset of Expected Machines
	inventory := &cwssaws.ExpectedMachineInventory{
		ExpectedMachines: pagedItems,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
		InventoryStatus: input.status,
		StatusMsg:       input.statusMessage,
		InventoryPage:   input.buildPage(),
	}
	if inventory.InventoryPage != nil {
		inventory.InventoryPage.ItemIds = itemIDs
	}
	return inventory
}

func expectedMachineFindFallback(ctx context.Context, carbideClient *cclient.CarbideClient) ([]string, []*cwssaws.ExpectedMachine, error) {
	// TODO: This is highly inefficient and transferring all data at once. Pagination support is needed in Carbide.
	// Retrieve all Expected Machines from Carbide
	emList, err := carbideClient.Carbide().GetAllExpectedMachines(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, nil, err
	}

	// Get all IDs
	var ids []string
	for _, it := range emList.ExpectedMachines {
		if it.Id != nil {
			ids = append(ids, it.Id.Value)
		}
	}

	return ids, emList.ExpectedMachines, nil
}

// ManageExpectedMachine is an activity wrapper for Expected Machine management
type ManageExpectedMachine struct {
	CarbideAtomicClient *cclient.CarbideAtomicClient
}

// NewManageExpectedMachine returns a new ManageExpectedMachine client
func NewManageExpectedMachine(carbideClient *cclient.CarbideAtomicClient) ManageExpectedMachine {
	return ManageExpectedMachine{
		CarbideAtomicClient: carbideClient,
	}
}

// CreateExpectedMachineOnSite creates Expected Machine with Carbide
func (mem *ManageExpectedMachine) CreateExpectedMachineOnSite(ctx context.Context, request *cwssaws.ExpectedMachine) error {
	logger := log.With().Str("Activity", "CreateExpectedMachineOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty create Expected Machine request")
	} else if id := request.GetId(); id == nil || (*id).String() == "" {
		err = errors.New("received create Expected Machine request without required id field")
	} else if request.GetBmcMacAddress() == "" || request.GetChassisSerialNumber() == "" {
		err = errors.New("received create Expected Machine request with missing MAC or serial")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mem.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	// Call Forge gRPC endpoint
	_, err = forgeClient.AddExpectedMachine(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create Expected Machine using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// UpdateExpectedMachineOnSite updates Expected Machine on Carbide
func (mem *ManageExpectedMachine) UpdateExpectedMachineOnSite(ctx context.Context, request *cwssaws.ExpectedMachine) error {
	logger := log.With().Str("Activity", "UpdateExpectedMachineOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty update Expected Machine request")
	} else if id := request.GetId(); id == nil || (*id).String() == "" {
		err = errors.New("received update Expected Machine request without required id field")
	} else if request.GetBmcMacAddress() == "" || request.GetChassisSerialNumber() == "" {
		err = errors.New("received update Expected Machine request with missing MAC or serial")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mem.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.UpdateExpectedMachine(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to update Expected Machine using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// DeleteExpectedMachineOnSite deletes Expected Machine on Carbide
func (mem *ManageExpectedMachine) DeleteExpectedMachineOnSite(ctx context.Context, request *cwssaws.ExpectedMachineRequest) error {
	logger := log.With().Str("Activity", "DeleteExpectedMachineOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty delete Expected Machine request")
	} else if id := request.GetId(); id == nil || (*id).String() == "" {
		err = errors.New("received delete Expected Machine request without required id field")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mem.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.DeleteExpectedMachine(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to delete Expected Machine using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}
