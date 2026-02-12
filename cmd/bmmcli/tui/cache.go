// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"time"
)

// NamedItem is a generic resource with a human-readable name and UUID
type NamedItem struct {
	Name   string
	ID     string
	Status string
	Extra  map[string]string // parent refs like "siteId", "vpcId"
	Raw    interface{}       // original API object
}

// Cache holds lazy-loaded resources for the interactive session
type Cache struct {
	items   map[string][]NamedItem
	fetched map[string]time.Time
	ttl     time.Duration
}

// NewCache creates a new session cache with default 30s TTL
func NewCache() *Cache {
	return &Cache{
		items:   make(map[string][]NamedItem),
		fetched: make(map[string]time.Time),
		ttl:     30 * time.Second,
	}
}

// Get returns cached items for a resource type, or nil if stale/missing
func (c *Cache) Get(resourceType string) []NamedItem {
	fetched, ok := c.fetched[resourceType]
	if !ok || time.Since(fetched) > c.ttl {
		return nil
	}
	return c.items[resourceType]
}

// Set stores items for a resource type
func (c *Cache) Set(resourceType string, items []NamedItem) {
	c.items[resourceType] = items
	c.fetched[resourceType] = time.Now()
}

// Invalidate clears a resource type from the cache
func (c *Cache) Invalidate(resourceType string) {
	delete(c.items, resourceType)
	delete(c.fetched, resourceType)
}

// LookupByName finds a cached item by name (case-insensitive)
func (c *Cache) LookupByName(resourceType, name string) *NamedItem {
	items := c.Get(resourceType)
	if items == nil {
		return nil
	}
	for _, item := range items {
		if equalFold(item.Name, name) {
			return &item
		}
	}
	return nil
}

// LookupByID finds a cached item by UUID
func (c *Cache) LookupByID(resourceType, id string) *NamedItem {
	items := c.Get(resourceType)
	if items == nil {
		return nil
	}
	for _, item := range items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
