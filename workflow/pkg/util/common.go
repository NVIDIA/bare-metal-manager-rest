// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package util

import (
	"slices"
	"time"

	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

var (
	// ErrMsgSiteControllerRowNotFound is returned when an entity is not found in Site Controller
	ErrMsgSiteControllerRowNotFound = "row not found"
	// ErrMsgSiteControllerNoRowsReturned is returned when lookup for entity returns nothing in Site Controller
	ErrMsgSiteControllerNoRowsReturned = "no rows returned"
	// ErrMsgSiteControllerMarkedForDeletion is returned when an entity is marked for deletion in Site Controller
	ErrMsgSiteControllerMarkedForDeletion = "marked for deletion"
	// ErrMsgSiteControllerCouldNotFind is returned when an entity is not found in Site Controller
	ErrMsgSiteControllerCouldNotFind = "could not find"
	// ErrMsgSiteControllerDuplicateEntryFound is returned when an entity is found in Site Controller
	ErrMsgSiteControllerDuplicateEntryFound = "duplicate key value violates unique constraint"
)

const (
	// InventoryReceiptInterval is the interval between 2 subsequent inventory receipts
	// TODO: Move this to cloud-common so it can be used by Site Agent as well
	InventoryReceiptInterval = 3 * time.Minute
)

func PtrsEqual[T comparable](i1 *T, i2 *T) bool {
	// They're either both nil or both non-nil
	// Otherwise, they certainly don't match.
	if (i1 == nil) != (i2 == nil) {
		return false
	}

	// We know their nil-ness is the same,
	// so if one is non-nil, then we can
	// compare the actual values being pointed
	// to by both.
	if i1 != nil && *i1 != *i2 {
		return false
	}

	return true
}

func NetworkSecurityGroupPropagationDetailsEqual(pd1, pd2 *cdbm.NetworkSecurityGroupPropagationDetails) bool {
	if (pd1 == nil) != (pd2 == nil) {
		return false
	}

	// If pd1 was nil but we made it here, then
	// both were nil, so we can return true.
	if pd1 == nil {
		return true
	}

	return pd1.Status.Number() == pd2.Status.Number() &&
		PtrsEqual(pd1.Details, pd2.Details) &&
		slices.Equal(pd1.UnpropagatedInstanceIds, pd2.UnpropagatedInstanceIds) &&
		slices.Equal(pd1.RelatedInstanceIds, pd2.RelatedInstanceIds)

}

func MachineCapabilitiesEqual(cap1 *cdbm.MachineCapability, cap2 *cdbm.MachineCapability) bool {

	return PtrsEqual(cap1.Cores, cap2.Cores) &&
		PtrsEqual(cap1.Threads, cap2.Threads) &&
		PtrsEqual(cap1.Count, cap2.Count) &&
		PtrsEqual(cap1.DeviceType, cap2.DeviceType) &&
		cap1.Name == cap2.Name &&
		cap1.Type == cap2.Type &&
		PtrsEqual(cap1.Capacity, cap2.Capacity) &&
		PtrsEqual(cap1.Frequency, cap2.Frequency) &&
		PtrsEqual(cap1.HardwareRevision, cap2.HardwareRevision) &&
		PtrsEqual(cap1.Vendor, cap2.Vendor) &&
		slices.Equal(cap1.InactiveDevices, cap2.InactiveDevices) &&
		cap1.Index == cap2.Index
}

// IsTimeWithinStaleInventoryThreshold checks if the action time is within the threshold where we could be processing an older inventory
func IsTimeWithinStaleInventoryThreshold(actionTime time.Time) bool {
	return time.Since(actionTime) < InventoryReceiptInterval+(time.Second*10)
}
