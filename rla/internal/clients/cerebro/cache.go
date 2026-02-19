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
	"sync"
)

const (
	// Locations types supported by Cerebro
	LocationTypeUnknown  CachedLocationType = "Unknown"
	LocationTypeProvider CachedLocationType = "Provider"
	LocationTypeRegion   CachedLocationType = "Region"
	LocationTypeSite     CachedLocationType = "Site"
	LocationTypeModule   CachedLocationType = "Module"
)

var (
	locationTypeToParentType = map[CachedLocationType]CachedLocationType{
		LocationTypeUnknown: LocationTypeUnknown,
		LocationTypeModule:  LocationTypeSite,
		LocationTypeSite:    LocationTypeRegion,
		LocationTypeRegion:  LocationTypeProvider,
	}

	locationTypeRoot = LocationTypeProvider
)

type CachedLocationType string

// Cache is a cache for some data fetched from Cerebro.
type Cache struct {
	mtx       sync.RWMutex
	locations map[string]*CachedLocation
}

type CachedLocation struct {
	Name   string
	Type   CachedLocationType
	Parent string
}

func NewCachedLocation(loc *Location) *CachedLocation {
	if loc == nil {
		return nil
	}

	return &CachedLocation{
		Name:   loc.Name,
		Type:   NewCachedLocationTypeFromString(loc.Type.Name),
		Parent: loc.Parent.Name,
	}
}

func (cl *CachedLocation) validateParent(parent *CachedLocation) bool {
	if cl.Type == locationTypeRoot {
		// The location with root type is the root of the chain, it should
		// not have a parent.
		return cl.Parent == "" && parent == nil
	}

	// The location should have a parent.
	if parent == nil {
		return false
	}

	if cl.Parent != parent.Name {
		// The parent should have the expected name.
		return false
	}

	return parent.Type == locationTypeToParentType[cl.Type]
}

func NewCache() *Cache {
	return &Cache{
		locations: make(map[string]*CachedLocation),
	}
}

func (c *Cache) GetLocation(name string) *CachedLocation {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if loc := c.locations[name]; loc != nil {
		cacheLoc := *loc
		return &cacheLoc
	}

	return nil
}

func (c *Cache) GetFullLocation(name string) map[CachedLocationType]CachedLocation {
	chain := make(map[CachedLocationType]CachedLocation)

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for loc := c.locations[name]; loc != nil; loc = c.locations[loc.Parent] {
		chain[loc.Type] = *loc
	}

	return chain
}

func (c *Cache) GetAllLocations() []CachedLocation {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	locations := make([]CachedLocation, 0, len(c.locations))
	for _, loc := range c.locations {
		locations = append(locations, *loc)
	}

	return locations
}

func (c *Cache) AddLocation(name string, loc *Location) bool {
	cacheLoc := NewCachedLocation(loc)
	if cacheLoc == nil {
		return false
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if _, ok := c.locations[name]; ok {
		return false
	}

	c.locations[name] = cacheLoc
	return true
}

func (c *Cache) RemoveLocation(name string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if _, ok := c.locations[name]; !ok {
		return false
	}

	delete(c.locations, name)
	return true
}

// ValidateLocationChain validates the location chain. It returns the list of
// cached locations which do not have valid parent informatino.
func (c *Cache) ValidateLocationChain() []CachedLocation {
	invalidLocations := make([]CachedLocation, 0)

	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for _, loc := range c.locations {
		if !loc.validateParent(c.locations[loc.Parent]) {
			invalidLocations = append(invalidLocations, *loc)
		}
	}

	return invalidLocations
}

func NewCachedLocationTypeFromString(t string) CachedLocationType {
	switch t {
	case string(LocationTypeProvider):
		return LocationTypeProvider
	case string(LocationTypeRegion):
		return LocationTypeRegion
	case string(LocationTypeSite):
		return LocationTypeSite
	case string(LocationTypeModule):
		return LocationTypeModule
	default:
		return LocationTypeUnknown
	}
}
