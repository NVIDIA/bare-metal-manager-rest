/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cerebro

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewCachedLocation(t *testing.T) {
	tests := map[string]struct {
		location *Location
		expected *CachedLocation
	}{
		"nil location": {
			location: nil,
			expected: nil,
		},
		"valid location": {
			location: &Location{
				Name:   "test-module",
				Type:   LocationType{Name: "Module"},
				Parent: ParentLocation{Name: "test-site"},
			},
			expected: &CachedLocation{
				Name:   "test-module",
				Type:   LocationTypeModule,
				Parent: "test-site",
			},
		},
		"unknown location type": {
			location: &Location{
				Name:   "test-unknown",
				Type:   LocationType{Name: "UnknownType"},
				Parent: ParentLocation{Name: "test-parent"},
			},
			expected: &CachedLocation{
				Name:   "test-unknown",
				Type:   LocationTypeUnknown,
				Parent: "test-parent",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := NewCachedLocation(tt.location)
			if result == nil && tt.expected == nil {
				return
			}
			if result == nil || tt.expected == nil {
				t.Errorf("NewCachedLocation() = %v, expected %v", result, tt.expected)
				return
			}
			if result.Name != tt.expected.Name || result.Type != tt.expected.Type || result.Parent != tt.expected.Parent {
				t.Errorf("NewCachedLocation() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCachedLocation_validateParent(t *testing.T) {
	tests := map[string]struct {
		location *CachedLocation
		parent   *CachedLocation
		expected bool
	}{
		"root location with no parent": {
			location: &CachedLocation{
				Name:   "test-provider",
				Type:   LocationTypeProvider,
				Parent: "",
			},
			parent:   nil,
			expected: true,
		},
		"root location with parent (invalid)": {
			location: &CachedLocation{
				Name:   "test-provider",
				Type:   LocationTypeProvider,
				Parent: "some-parent",
			},
			parent:   &CachedLocation{Name: "some-parent", Type: LocationTypeProvider},
			expected: false,
		},
		"non-root location with nil parent (invalid)": {
			location: &CachedLocation{
				Name:   "test-site",
				Type:   LocationTypeSite,
				Parent: "test-region",
			},
			parent:   nil,
			expected: false,
		},
		"valid parent-child relationship": {
			location: &CachedLocation{
				Name:   "test-site",
				Type:   LocationTypeSite,
				Parent: "test-region",
			},
			parent: &CachedLocation{
				Name: "test-region",
				Type: LocationTypeRegion,
			},
			expected: true,
		},
		"wrong parent name": {
			location: &CachedLocation{
				Name:   "test-site",
				Type:   LocationTypeSite,
				Parent: "expected-region",
			},
			parent: &CachedLocation{
				Name: "actual-region",
				Type: LocationTypeRegion,
			},
			expected: false,
		},
		"wrong parent type": {
			location: &CachedLocation{
				Name:   "test-site",
				Type:   LocationTypeSite,
				Parent: "test-parent",
			},
			parent: &CachedLocation{
				Name: "test-parent",
				Type: LocationTypeSite, // Should be Region
			},
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.location.validateParent(tt.parent)
			if result != tt.expected {
				t.Errorf("validateParent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewCache(t *testing.T) {
	cache := NewCache()
	if cache == nil {
		t.Error("NewCache() should not return nil")
	}
	if cache.locations == nil {
		t.Error("NewCache() should initialize locations map")
	}
	if len(cache.locations) != 0 {
		t.Error("NewCache() should initialize empty locations map")
	}
}

func TestCache_GetLocation(t *testing.T) {
	cache := NewCache()

	// Test getting non-existent location
	result := cache.GetLocation("non-existent")
	if result != nil {
		t.Error("GetLocation() should return nil for non-existent location")
	}

	// Add a location
	testLoc := &CachedLocation{
		Name:   "test-location",
		Type:   LocationTypeModule,
		Parent: "test-parent",
	}
	cache.locations["test-location"] = testLoc

	// Test getting existing location
	result = cache.GetLocation("test-location")
	if result == nil {
		t.Error("GetLocation() should return location for existing key")
	}
	if result.Name != testLoc.Name || result.Type != testLoc.Type || result.Parent != testLoc.Parent {
		t.Errorf("GetLocation() = %v, expected %v", result, testLoc)
	}

	// Verify it returns a copy (not the same pointer)
	if result == testLoc {
		t.Error("GetLocation() should return a copy, not the original pointer")
	}
}

func TestCache_AddLocation(t *testing.T) {
	cache := NewCache()

	// Test adding valid location
	testLocation := &Location{
		Name:   "test-module",
		Type:   LocationType{Name: "Module"},
		Parent: ParentLocation{Name: "test-site"},
	}

	success := cache.AddLocation("test-module", testLocation)
	if !success {
		t.Error("AddLocation() should return true for valid location")
	}

	// Verify location was added
	result := cache.GetLocation("test-module")
	if result == nil {
		t.Error("Location should be retrievable after adding")
	}

	// Test adding duplicate location (should fail)
	success = cache.AddLocation("test-module", testLocation)
	if success {
		t.Error("AddLocation() should return false for duplicate location")
	}

	// Test adding nil location
	success = cache.AddLocation("nil-test", nil)
	if success {
		t.Error("AddLocation() should return false for nil location")
	}
}

func TestCache_RemoveLocation(t *testing.T) {
	cache := NewCache()

	// Test removing non-existent location
	success := cache.RemoveLocation("non-existent")
	if success {
		t.Error("RemoveLocation() should return false for non-existent location")
	}

	// Add a location
	testLoc := &CachedLocation{
		Name:   "test-location",
		Type:   LocationTypeModule,
		Parent: "test-parent",
	}
	cache.locations["test-location"] = testLoc

	// Test removing existing location
	success = cache.RemoveLocation("test-location")
	if !success {
		t.Error("RemoveLocation() should return true for existing location")
	}

	// Verify location was removed
	result := cache.GetLocation("test-location")
	if result != nil {
		t.Error("Location should not be retrievable after removal")
	}
}

func TestCache_GetAllLocations(t *testing.T) {
	cache := NewCache()

	// Test empty cache
	locations := cache.GetAllLocations()
	if len(locations) != 0 {
		t.Error("GetAllLocations() should return empty slice for empty cache")
	}

	// Add some locations
	testLocs := []*CachedLocation{
		{Name: "loc1", Type: LocationTypeProvider, Parent: ""},
		{Name: "loc2", Type: LocationTypeRegion, Parent: "loc1"},
		{Name: "loc3", Type: LocationTypeSite, Parent: "loc2"},
	}

	for _, loc := range testLocs {
		cache.locations[loc.Name] = loc
	}

	// Test getting all locations
	locations = cache.GetAllLocations()
	if len(locations) != len(testLocs) {
		t.Errorf("GetAllLocations() returned %d locations, expected %d", len(locations), len(testLocs))
	}

	// Verify all locations are present (order doesn't matter)
	found := make(map[string]bool)
	for _, loc := range locations {
		found[loc.Name] = true
	}

	for _, expected := range testLocs {
		if !found[expected.Name] {
			t.Errorf("GetAllLocations() missing location: %s", expected.Name)
		}
	}
}

func TestCache_GetFullLocation(t *testing.T) {
	cache := NewCache()

	// Create a location chain: provider -> region -> site -> module
	provider := &CachedLocation{Name: "provider1", Type: LocationTypeProvider, Parent: ""}
	region := &CachedLocation{Name: "region1", Type: LocationTypeRegion, Parent: "provider1"}
	site := &CachedLocation{Name: "site1", Type: LocationTypeSite, Parent: "region1"}
	module := &CachedLocation{Name: "module1", Type: LocationTypeModule, Parent: "site1"}

	cache.locations["provider1"] = provider
	cache.locations["region1"] = region
	cache.locations["site1"] = site
	cache.locations["module1"] = module

	// Test getting full location chain from module
	chain := cache.GetFullLocation("module1")

	expectedTypes := []CachedLocationType{LocationTypeModule, LocationTypeSite, LocationTypeRegion, LocationTypeProvider}
	if len(chain) != len(expectedTypes) {
		t.Errorf("GetFullLocation() returned chain of length %d, expected %d", len(chain), len(expectedTypes))
	}

	for _, expectedType := range expectedTypes {
		if _, exists := chain[expectedType]; !exists {
			t.Errorf("GetFullLocation() missing location type: %s", expectedType)
		}
	}

	// Test getting full location for non-existent location
	emptyChain := cache.GetFullLocation("non-existent")
	if len(emptyChain) != 0 {
		t.Error("GetFullLocation() should return empty chain for non-existent location")
	}
}

func TestCache_ValidateLocationChain(t *testing.T) {
	cache := NewCache()

	// Create valid chain
	provider := &CachedLocation{Name: "provider1", Type: LocationTypeProvider, Parent: ""}
	region := &CachedLocation{Name: "region1", Type: LocationTypeRegion, Parent: "provider1"}
	site := &CachedLocation{Name: "site1", Type: LocationTypeSite, Parent: "region1"}

	// Create invalid location (wrong parent type)
	invalidLoc := &CachedLocation{Name: "invalid1", Type: LocationTypeSite, Parent: "provider1"}

	cache.locations["provider1"] = provider
	cache.locations["region1"] = region
	cache.locations["site1"] = site
	cache.locations["invalid1"] = invalidLoc

	// Test validation
	invalidLocations := cache.ValidateLocationChain()

	if len(invalidLocations) != 1 {
		t.Errorf("ValidateLocationChain() found %d invalid locations, expected 1", len(invalidLocations))
	}

	if len(invalidLocations) > 0 && invalidLocations[0].Name != "invalid1" {
		t.Errorf("ValidateLocationChain() found invalid location %s, expected invalid1", invalidLocations[0].Name)
	}
}

func TestNewCachedLocationTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected CachedLocationType
	}{
		{"Provider", LocationTypeProvider},
		{"Region", LocationTypeRegion},
		{"Site", LocationTypeSite},
		{"Module", LocationTypeModule},
		{"UnknownType", LocationTypeUnknown},
		{"", LocationTypeUnknown},
		{"provider", LocationTypeUnknown}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NewCachedLocationTypeFromString(tt.input)
			if result != tt.expected {
				t.Errorf("NewCachedLocationTypeFromString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCache_ConcurrentOperations(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup

	// Test concurrent adds
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			location := &Location{
				Name:   fmt.Sprintf("location-%d", id),
				Type:   LocationType{Name: "Module"},
				Parent: ParentLocation{Name: "parent"},
			}

			cache.AddLocation(fmt.Sprintf("location-%d", id), location)
		}(i)
	}

	wg.Wait()

	// Verify all locations were added
	locations := cache.GetAllLocations()
	if len(locations) != numGoroutines {
		t.Errorf("Expected %d locations after concurrent adds, got %d", numGoroutines, len(locations))
	}

	// Test concurrent reads and removes
	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		// Concurrent reads
		go func(id int) {
			defer wg.Done()
			cache.GetLocation(fmt.Sprintf("location-%d", id))
		}(i)

		// Concurrent removes
		go func(id int) {
			defer wg.Done()
			cache.RemoveLocation(fmt.Sprintf("location-%d", id))
		}(i)
	}

	wg.Wait()

	// After removes, cache should be empty
	locations = cache.GetAllLocations()
	if len(locations) != 0 {
		t.Errorf("Expected 0 locations after concurrent removes, got %d", len(locations))
	}
}

func TestLocationTypeToParentType(t *testing.T) {
	// Test the global mapping variable
	expectedMappings := map[CachedLocationType]CachedLocationType{
		LocationTypeUnknown: LocationTypeUnknown,
		LocationTypeModule:  LocationTypeSite,
		LocationTypeSite:    LocationTypeRegion,
		LocationTypeRegion:  LocationTypeProvider,
	}

	for childType, expectedParentType := range expectedMappings {
		actualParentType, exists := locationTypeToParentType[childType]
		if !exists {
			t.Errorf("locationTypeToParentType missing entry for %s", childType)
		}
		if actualParentType != expectedParentType {
			t.Errorf("locationTypeToParentType[%s] = %s, expected %s", childType, actualParentType, expectedParentType)
		}
	}

	// Test that locationTypeRoot is set correctly
	if locationTypeRoot != LocationTypeProvider {
		t.Errorf("locationTypeRoot = %s, expected %s", locationTypeRoot, LocationTypeProvider)
	}
}
