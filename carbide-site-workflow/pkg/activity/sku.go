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

package activity

import (
	"context"

	cwssaws "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	cclient "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/grpc/client"
	"github.com/rs/zerolog/log"
)

// ManageSkuInventory is an activity wrapper for Sku inventory collection and publishing
type ManageSkuInventory struct {
	config ManageInventoryConfig
}

// DiscoverSkuInventory is an activity to collect Sku inventory and publish to Temporal queue
func (msi *ManageSkuInventory) DiscoverSkuInventory(ctx context.Context) error {
	logger := log.With().Str("Activity", "DiscoverSkuInventory").Logger()
	logger.Info().Msg("Starting activity")
	inventoryImpl := manageInventoryImpl[string, *cwssaws.Sku, *cwssaws.SkuInventory]{
		itemType:               "Sku",
		config:                 msi.config,
		internalFindIDs:        skuFindIDs,
		internalFindByIDs:      skuFindByIDs,
		internalPagedInventory: skuPagedInventory,
	}
	return inventoryImpl.CollectAndPublishInventory(ctx, &logger)
}

// NewManageSkuInventory returns a ManageInventory implementation for Sku activity
func NewManageSkuInventory(config ManageInventoryConfig) ManageSkuInventory {
	return ManageSkuInventory{
		config: config,
	}
}

func skuFindIDs(ctx context.Context, carbideClient *cclient.CarbideClient) ([]string, error) {
	forgeClient := carbideClient.Carbide()
	result, err := forgeClient.GetAllSkuIds(ctx, nil)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(result.Ids))
	for i, id := range result.Ids {
		ids[i] = id
	}

	return ids, nil
}

func skuFindByIDs(ctx context.Context, carbideClient *cclient.CarbideClient, ids []string) ([]*cwssaws.Sku, error) {
	forgeClient := carbideClient.Carbide()
	result, err := forgeClient.FindSkusByIds(ctx, &cwssaws.SkusByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	return result.Skus, nil
}

func skuPagedInventory(ids []string, skus []*cwssaws.Sku, input *pagedInventoryInput) *cwssaws.SkuInventory {
	// Build inventory page
	page := input.buildPage()

	// Copy IDs to page
	page.ItemIds = make([]string, len(ids))
	copy(page.ItemIds, ids)

	return &cwssaws.SkuInventory{
		InventoryStatus: input.status,
		StatusMsg:       input.statusMessage,
		Skus:            skus,
		InventoryPage:   page,
	}
}
