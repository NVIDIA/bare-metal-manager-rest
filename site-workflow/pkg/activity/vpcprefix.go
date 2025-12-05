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

// ManageVpcPrefix is an activity wrapper for VpcPrefix management
type ManageVpcPrefix struct {
	CarbideAtomicClient *client.CarbideAtomicClient
}

// NewManageVpcPrefix returns a new ManageVpcPrefix client
func NewManageVpcPrefix(carbideClient *client.CarbideAtomicClient) ManageVpcPrefix {
	return ManageVpcPrefix{
		CarbideAtomicClient: carbideClient,
	}
}

// Function to create VpcPrefix with Carbide
func (mvp *ManageVpcPrefix) CreateVpcPrefixOnSite(ctx context.Context, request *cwssaws.VpcPrefixCreationRequest) error {
	logger := log.With().Str("Activity", "CreateVpcPrefixOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty create VPC prefix request")
	} else if request.Name == "" {
		err = errors.New("received create VPC prefix request missing Name")
	} else if request.Prefix == "" {
		err = errors.New("received create VPC prefix request missing Prefix")
	} else if request.Id == nil || request.Id.Value == "" {
		// Don't let a request come in without a cloud-provided ID
		// or carbide will generate one and cloud won't know the relationship.
		err = errors.New("received create VPC prefix request missing ID")
	} else if request.VpcId == nil || request.VpcId.Value == "" {
		err = errors.New("received create VPC prefix request missing VPC ID")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mvp.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.CreateVpcPrefix(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create VPC prefix using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// Function to update VpcPrefix with Carbide
func (mvp *ManageVpcPrefix) UpdateVpcPrefixOnSite(ctx context.Context, request *cwssaws.VpcPrefixUpdateRequest) error {
	logger := log.With().Str("Activity", "UpdateVpcPrefixOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty update VPC prefix request")
	} else if request.Id == nil || request.Id.Value == "" {
		err = errors.New("received update VPC prefix request missing ID")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mvp.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.UpdateVpcPrefix(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to update VPC prefix using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// Function to delete VpcPrefix on Carbide
func (mvp *ManageVpcPrefix) DeleteVpcPrefixOnSite(ctx context.Context, request *cwssaws.VpcPrefixDeletionRequest) error {
	logger := log.With().Str("Activity", "DeleteVpcPrefixOnSite").Logger()

	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	if request == nil {
		err = errors.New("received empty delete VPC prefix request")
	} else if request.Id == nil || request.Id.Value == "" {
		err = errors.New("reveived delete VPC prefix missing VPC Prefix ID")
	}

	if err != nil {
		return temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call Site Controller gRPC endpoint
	carbideClient := mvp.CarbideAtomicClient.GetClient()
	forgeClient := carbideClient.Carbide()

	_, err = forgeClient.DeleteVpcPrefix(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to delete VPC prefix using Site Controller API")
		return swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return nil
}

// ManageVpcPrefixInventory is an activity wrapper for VpcPrefix inventory collection and publishing
type ManageVpcPrefixInventory struct {
	config ManageInventoryConfig
}

// NewManageVpcPrefixInventory returns a ManageInventory implementation for VpcPrefix
func NewManageVpcPrefixInventory(config ManageInventoryConfig) ManageVpcPrefixInventory {
	return ManageVpcPrefixInventory{
		config: config,
	}
}

// DiscoverVpcPrefixInventory is an activity to collect VpcPrefix inventory and publish to Temporal queue
func (moii *ManageVpcPrefixInventory) DiscoverVpcPrefixInventory(ctx context.Context) error {
	logger := log.With().Str("Activity", "DiscoverVpcPrefixInventory").Logger()
	logger.Info().Msg("Starting activity")

	inventoryImpl := manageInventoryImpl[*cwssaws.VpcPrefixId, *cwssaws.VpcPrefix, *cwssaws.VpcPrefixInventory]{
		itemType:               "VpcPrefix",
		config:                 moii.config,
		internalFindIDs:        VpcPrefixFindIDs,
		internalFindByIDs:      VpcPrefixFindByIDs,
		internalPagedInventory: VpcPrefixPagedInventory,
		internalFindFallback:   VpcPrefixFindFallback,
	}
	return inventoryImpl.CollectAndPublishInventory(ctx, &logger)
}

func VpcPrefixFindIDs(ctx context.Context, carbideClient *cClient.CarbideClient) ([]*cwssaws.VpcPrefixId, error) {
	idList, err := carbideClient.Carbide().SearchVpcPrefixes(ctx, &cwssaws.VpcPrefixSearchQuery{})
	if err != nil {
		return nil, err
	}
	return idList.VpcPrefixIds, nil
}

func VpcPrefixFindByIDs(ctx context.Context, carbideClient *cClient.CarbideClient, ids []*cwssaws.VpcPrefixId) ([]*cwssaws.VpcPrefix, error) {
	list, err := carbideClient.Carbide().GetVpcPrefixes(ctx, &cwssaws.VpcPrefixGetRequest{
		VpcPrefixIds: ids,
	})

	if err != nil {
		return nil, err
	}
	return list.GetVpcPrefixes(), nil
}

func VpcPrefixPagedInventory(allItemIDs []*cwssaws.VpcPrefixId, pagedItems []*cwssaws.VpcPrefix, input *pagedInventoryInput) *cwssaws.VpcPrefixInventory {
	itemIDs := []string{}
	for _, id := range allItemIDs {
		itemIDs = append(itemIDs, id.GetValue())
	}

	// Create an inventory page with the subset of VpcPrefixs
	inventory := &cwssaws.VpcPrefixInventory{
		VpcPrefixes: pagedItems,
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

func VpcPrefixFindFallback(ctx context.Context, carbideClient *cClient.CarbideClient) ([]*cwssaws.VpcPrefixId, []*cwssaws.VpcPrefix, error) {
	request := &cwssaws.VpcPrefixGetRequest{}
	items, err := carbideClient.Carbide().GetVpcPrefixes(ctx, request)
	if err != nil {
		return nil, nil, err
	}

	var ids []*cwssaws.VpcPrefixId
	for _, it := range items.GetVpcPrefixes() {
		ids = append(ids, it.GetId())
	}
	return ids, items.GetVpcPrefixes(), nil
}
