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

	"github.com/rs/zerolog/log"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	swe "github.com/nvidia/carbide-rest/site-workflow/pkg/error"
	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	cClient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ManageInfiniBandPartitionInventory is an activity wrapper for InfiniBand Partition inventory collection and publishing
type ManageInfiniBandPartitionInventory struct {
	config ManageInventoryConfig
}

// DiscoverInfiniBandPartitionInventory is an activity to collect InfiniBand Partition inventory and publish to Temporal queue
func (mmi *ManageInfiniBandPartitionInventory) DiscoverInfiniBandPartitionInventory(ctx context.Context) error {
	logger := log.With().Str("Activity", "DiscoverIBPartitionInventory").Logger()
	logger.Info().Msg("Starting activity")
	inventoryImpl := manageInventoryImpl[*cwssaws.IBPartitionId, *cwssaws.IBPartition, *cwssaws.InfiniBandPartitionInventory]{
		itemType:               "InfiniBandPartition",
		config:                 mmi.config,
		internalFindIDs:        ibpFindIDs,
		internalFindByIDs:      ibpFindByIDs,
		internalPagedInventory: ibpPagedInventory,
		internalFindFallback:   ibpFindFallback,
	}
	return inventoryImpl.CollectAndPublishInventory(ctx, &logger)
}

// NewManageInfiniBandPartitionInventory returns a ManageInventory implementation for InfiniBand Partition activity
func NewManageInfiniBandPartitionInventory(config ManageInventoryConfig) ManageInfiniBandPartitionInventory {
	return ManageInfiniBandPartitionInventory{
		config: config,
	}
}

func ibpFindIDs(ctx context.Context, carbideClient *cClient.CarbideClient) ([]*cwssaws.IBPartitionId, error) {
	idList, err := carbideClient.Networks().FindInfinibandPartitionIDs(ctx, &cwssaws.IBPartitionSearchFilter{})
	if err != nil {
		return nil, err
	}
	return idList.GetIbPartitionIds(), nil
}

func ibpFindByIDs(ctx context.Context, carbideClient *cClient.CarbideClient, ids []*cwssaws.IBPartitionId) ([]*cwssaws.IBPartition, error) {
	list, err := carbideClient.Networks().FindInfinibandPartitionsByIDs(ctx, &cwssaws.IBPartitionsByIdsRequest{
		IbPartitionIds: ids,
	})
	if err != nil {
		return nil, err
	}
	return list.GetIbPartitions(), nil
}

func ibpPagedInventory(allItemIDs []*cwssaws.IBPartitionId, pagedItems []*cwssaws.IBPartition, input *pagedInventoryInput) *cwssaws.InfiniBandPartitionInventory {
	itemIDs := []string{}
	for _, id := range allItemIDs {
		itemIDs = append(itemIDs, id.GetValue())
	}

	// Create an inventory page with the subset of VPCs
	inventory := &cwssaws.InfiniBandPartitionInventory{
		IbPartitions: pagedItems,
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

func ibpFindFallback(ctx context.Context, carbideClient *cClient.CarbideClient) ([]*cwssaws.IBPartitionId, []*cwssaws.IBPartition, error) {
	items, err := carbideClient.Networks().ListInfiniBandPartition(ctx)
	if err != nil {
		return nil, nil, err
	}
	var ids []*cwssaws.IBPartitionId
	for _, it := range items.GetIbPartitions() {
		ids = append(ids, it.Id)
	}
	return ids, items.GetIbPartitions(), nil
}

// ManageInfiniBandPartition is an activity wrapper for InfiniBand Partition management
type ManageInfiniBandPartition struct {
	CarbideAtomicClient *client.CarbideAtomicClient
}

// NewManageInfiniBandPartition returns a new ManageInfiniBandPartition client
func NewManageInfiniBandPartition(carbideClient *client.CarbideAtomicClient) ManageInfiniBandPartition {
	return ManageInfiniBandPartition{
		CarbideAtomicClient: carbideClient,
	}
}

// Function to create InfiniBand Partition with Carbide
func (mibp *ManageInfiniBandPartition) CreateInfiniBandPartitionOnSite(ctx context.Context, request *cwssaws.IBPartitionCreationRequest) error {
	logger := log.With().Str("Activity", "CreateInfiniBandPartitionOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty create InfiniBand Partition request")
	} else if request.Id == nil || request.GetId().GetValue() == "" {
		err = errors.New("received create InfiniBand Partition request missing ID")
	} else if request.Config == nil {
		err = errors.New("received create InfiniBand Partition request missing Config")
	} else if request.Config.Name == "" {
		err = errors.New("received create InfiniBand Partition request missing Name")
	} else if request.Config.TenantOrganizationId == "" {
		err = errors.New("received create InfiniBand Partition request missing TenantOrganizationId")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mibp.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	// Call Forge gRPC endpoint
	_, err = forgeClient.CreateIBPartition(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create InfiniBand Partition using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// Function to delete InfiniBand Partition on Carbide
func (mipb *ManageInfiniBandPartition) DeleteInfiniBandPartitionOnSite(ctx context.Context, request *cwssaws.IBPartitionDeletionRequest) error {
	logger := log.With().Str("Activity", "DeleteInfiniBandPartitionOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty delete InfiniBand Partition request")
	} else if request.Id == nil || request.Id.GetValue() == "" {
		err = errors.New("received delete InfiniBand Partition request without ID")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mipb.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.DeleteIBPartition(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to delete InfiniBand Partition using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}
