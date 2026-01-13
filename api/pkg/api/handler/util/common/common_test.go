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

package common

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUniqueChecker_Add_Basic tests the Add method
func TestUniqueChecker_Add_Basic(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	// Add first entry with unique value
	checker.Add(id1, "00:11:22:33:44:55")

	// Add second entry with different unique value
	checker.Add(id2, "AA:BB:CC:DD:EE:FF")

	// Add third entry with duplicate unique value
	checker.Add(id3, "00:11:22:33:44:55")

	// Should have 1 duplicate
	duplicates := checker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "00:11:22:33:44:55")

	// Should detect duplicates
	assert.True(t, checker.HasDuplicates())
}

// TestUniqueChecker_Add_NoDuplicates tests that no duplicates are detected when all values are unique
func TestUniqueChecker_Add_NoDuplicates(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	// Add entries with all unique values
	checker.Add(id1, "00:11:22:33:44:55")
	checker.Add(id2, "AA:BB:CC:DD:EE:FF")
	checker.Add(id3, "FF:FF:FF:FF:FF:FF")

	// Should have no duplicates
	duplicates := checker.GetDuplicates()
	assert.Empty(t, duplicates)
	assert.False(t, checker.HasDuplicates())
}

// TestUniqueChecker_Update tests the Update method
func TestUniqueChecker_Update(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	id2 := uuid.New()

	// Initial add
	oldMac := "00:11:22:33:44:55"
	newMac := "AA:BB:CC:DD:EE:FF"

	checker.Add(id1, oldMac)
	checker.Add(id2, "11:22:33:44:55:66")

	// Initially no duplicates
	assert.False(t, checker.HasDuplicates())

	// Update id1 to a new unique value
	checker.Update(id1, newMac)

	// Old value should no longer be counted
	assert.False(t, checker.HasDuplicates())

	// Update id1 to same value as id2 - should create duplicate
	checker.Update(id1, "11:22:33:44:55:66")

	// Should now have duplicates
	duplicates := checker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "11:22:33:44:55:66")
	assert.True(t, checker.HasDuplicates())
}

// TestUniqueChecker_Update_NoChange tests updating with same value
func TestUniqueChecker_Update_NoChange(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	mac := "00:11:22:33:44:55"

	checker.Add(id1, mac)

	// Update with same value should be no-op
	checker.Update(id1, mac)

	// Should still have no duplicates
	assert.False(t, checker.HasDuplicates())
	assert.Empty(t, checker.GetDuplicates())
}

// TestUniqueChecker_Update_NewID tests updating a new ID that doesn't exist yet
func TestUniqueChecker_Update_NewID(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	id2 := uuid.New()

	checker.Add(id1, "00:11:22:33:44:55")

	// Update a new ID that hasn't been added yet
	checker.Update(id2, "AA:BB:CC:DD:EE:FF")

	// Should have no duplicates
	assert.False(t, checker.HasDuplicates())
	
	// Verify id2 was added
	assert.Equal(t, "AA:BB:CC:DD:EE:FF", checker.IdToUniqueValue[id2])
}

// TestUniqueChecker_GetDuplicates tests the GetDuplicates method with multiple duplicates
func TestUniqueChecker_GetDuplicates(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	id4 := uuid.New()
	id5 := uuid.New()

	// Add values where two unique values are duplicated
	checker.Add(id1, "MAC-A")
	checker.Add(id2, "MAC-B")
	checker.Add(id3, "MAC-A") // duplicate of id1
	checker.Add(id4, "MAC-C")
	checker.Add(id5, "MAC-A") // another duplicate of id1

	// Should have 1 duplicate value (MAC-A appears 3 times)
	duplicates := checker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "MAC-A")
	assert.True(t, checker.HasDuplicates())
}

// TestUniqueChecker_GetDuplicates_MultipleDuplicates tests multiple different duplicate values
func TestUniqueChecker_GetDuplicates_MultipleDuplicates(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	// Create 6 IDs
	ids := make([]uuid.UUID, 6)
	for i := range ids {
		ids[i] = uuid.New()
	}

	// Add values where two different unique values are duplicated
	checker.Add(ids[0], "MAC-A")
	checker.Add(ids[1], "MAC-B")
	checker.Add(ids[2], "MAC-A") // duplicate MAC-A
	checker.Add(ids[3], "MAC-C")
	checker.Add(ids[4], "MAC-B") // duplicate MAC-B
	checker.Add(ids[5], "MAC-D")

	// Should have 2 duplicate values
	duplicates := checker.GetDuplicates()
	assert.Len(t, duplicates, 2)
	assert.Contains(t, duplicates, "MAC-A")
	assert.Contains(t, duplicates, "MAC-B")
	assert.True(t, checker.HasDuplicates())
}

// TestUniqueChecker_BatchOperationExample demonstrates usage in a batch operation
func TestUniqueChecker_BatchOperationExample(t *testing.T) {
	// Simulate batch create operation with MAC addresses
	type MachineRequest struct {
		BmcMacAddress       string
		ChassisSerialNumber string
	}

	requests := []MachineRequest{
		{BmcMacAddress: "00:11:22:33:44:55", ChassisSerialNumber: "SN001"},
		{BmcMacAddress: "AA:BB:CC:DD:EE:FF", ChassisSerialNumber: "SN002"},
		{BmcMacAddress: "00:11:22:33:44:55", ChassisSerialNumber: "SN003"}, // Duplicate MAC
		{BmcMacAddress: "FF:FF:FF:FF:FF:FF", ChassisSerialNumber: "SN001"}, // Duplicate Serial
	}

	macChecker := NewUniqueChecker[int]()
	serialChecker := NewUniqueChecker[int]()

	// Add all requests to checkers
	for i, req := range requests {
		macChecker.Add(i, req.BmcMacAddress)
		serialChecker.Add(i, req.ChassisSerialNumber)
	}

	// Check for duplicates
	macDuplicates := macChecker.GetDuplicates()
	serialDuplicates := serialChecker.GetDuplicates()

	// Should have detected 1 MAC duplicate and 1 Serial duplicate
	assert.Len(t, macDuplicates, 1)
	assert.Contains(t, macDuplicates, "00:11:22:33:44:55")

	assert.Len(t, serialDuplicates, 1)
	assert.Contains(t, serialDuplicates, "SN001")

	// Can build error messages
	var validationErrors []string
	for _, dup := range macDuplicates {
		validationErrors = append(validationErrors, fmt.Sprintf("duplicate BMC MAC address '%s'", dup))
	}
	for _, dup := range serialDuplicates {
		validationErrors = append(validationErrors, fmt.Sprintf("duplicate chassis serial number '%s'", dup))
	}

	assert.Len(t, validationErrors, 2)
}

// TestUniqueChecker_UpdateScenario demonstrates usage in a batch update operation
func TestUniqueChecker_UpdateScenario(t *testing.T) {
	// Simulate existing machines in database
	id1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	id2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	id3 := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	existingMachines := map[uuid.UUID]string{
		id1: "MAC-001",
		id2: "MAC-002",
		id3: "MAC-003",
	}

	// Simulate update requests
	type UpdateRequest struct {
		ID            uuid.UUID
		NewBmcAddress string
	}

	requests := []UpdateRequest{
		{
			ID:            id1,
			NewBmcAddress: "MAC-NEW-1", // Change to unique value - OK
		},
		{
			ID:            id3,
			NewBmcAddress: "MAC-002", // Change to same value as id2 - DUPLICATE
		},
	}

	macChecker := NewUniqueChecker[uuid.UUID]()

	// First, populate checker with existing machines
	for machineID, mac := range existingMachines {
		macChecker.Add(machineID, mac)
	}

	// Initially no duplicates
	assert.False(t, macChecker.HasDuplicates())

	// Now process update requests
	for _, req := range requests {
		macChecker.Update(req.ID, req.NewBmcAddress)
	}

	// Should have detected 1 conflict (MAC-002 now used by both id2 and id3)
	duplicates := macChecker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "MAC-002")
	assert.True(t, macChecker.HasDuplicates())
}

// TestUniqueChecker_ComplexScenario tests a complex scenario with adds and updates
func TestUniqueChecker_ComplexScenario(t *testing.T) {
	checker := NewUniqueChecker[uuid.UUID]()

	// Create some IDs
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	id4 := uuid.New()

	// Initial state: 4 machines with unique MACs
	checker.Add(id1, "MAC-A")
	checker.Add(id2, "MAC-B")
	checker.Add(id3, "MAC-C")
	checker.Add(id4, "MAC-D")

	// No duplicates
	assert.False(t, checker.HasDuplicates())

	// Update id3 to use same MAC as id1
	checker.Update(id3, "MAC-A")

	// Should have 1 duplicate
	duplicates := checker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "MAC-A")

	// Update id3 back to unique value
	checker.Update(id3, "MAC-E")

	// No duplicates again
	assert.False(t, checker.HasDuplicates())

	// Update id1 and id2 to same value
	checker.Update(id1, "MAC-X")
	checker.Update(id2, "MAC-X")

	// Should have 1 duplicate
	duplicates = checker.GetDuplicates()
	assert.Len(t, duplicates, 1)
	assert.Contains(t, duplicates, "MAC-X")
}
