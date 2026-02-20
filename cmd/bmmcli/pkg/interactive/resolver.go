// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"context"
	"fmt"
	"strings"
)

// FetchFunc fetches resources from the API and returns them as NamedItems.
type FetchFunc func(ctx context.Context) ([]NamedItem, error)

// Resolver resolves resource names to UUIDs using the cache and interactive select.
type Resolver struct {
	cache    *Cache
	fetchers map[string]FetchFunc
}

// NewResolver creates a new resolver with the given cache.
func NewResolver(cache *Cache) *Resolver {
	return &Resolver{
		cache:    cache,
		fetchers: make(map[string]FetchFunc),
	}
}

// RegisterFetcher registers a fetch function for a resource type.
func (r *Resolver) RegisterFetcher(resourceType string, fn FetchFunc) {
	r.fetchers[resourceType] = fn
}

// Fetch fetches (or returns cached) items for a resource type.
func (r *Resolver) Fetch(ctx context.Context, resourceType string) ([]NamedItem, error) {
	if items := r.cache.Get(resourceType); items != nil {
		return items, nil
	}

	fn, ok := r.fetchers[resourceType]
	if !ok {
		return nil, fmt.Errorf("no fetcher registered for resource type %q", resourceType)
	}

	items, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	r.cache.Set(resourceType, items)
	return items, nil
}

// FetchFiltered fetches items and filters by a parent key/value in Extra.
func (r *Resolver) FetchFiltered(ctx context.Context, resourceType, filterKey, filterValue string) ([]NamedItem, error) {
	items, err := r.Fetch(ctx, resourceType)
	if err != nil {
		return nil, err
	}

	var filtered []NamedItem
	for _, item := range items {
		if item.Extra != nil && item.Extra[filterKey] == filterValue {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

// Resolve interactively selects a resource by showing a TUI select menu.
func (r *Resolver) Resolve(ctx context.Context, resourceType, label string) (*NamedItem, error) {
	items, err := r.Fetch(ctx, resourceType)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", resourceType, err)
	}

	return r.SelectFromItems(label, items)
}

// ResolveFiltered interactively selects from a filtered set of resources.
func (r *Resolver) ResolveFiltered(ctx context.Context, resourceType, label, filterKey, filterValue string) (*NamedItem, error) {
	items, err := r.FetchFiltered(ctx, resourceType, filterKey, filterValue)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", resourceType, err)
	}

	return r.SelectFromItems(label, items)
}

// SelectFromItems shows an interactive select menu for the given items.
func (r *Resolver) SelectFromItems(label string, items []NamedItem) (*NamedItem, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no %s available", label)
	}

	if len(items) == 1 {
		fmt.Printf("%s %s %s\n", Bold(label+":"), Green(items[0].Name), Dim("(auto-selected)"))
		return &items[0], nil
	}

	selectItems := make([]SelectItem, len(items))
	for i, item := range items {
		lbl := item.Name
		if item.Status != "" {
			lbl += "  " + Dim(item.Status)
		}
		selectItems[i] = SelectItem{
			Label: lbl,
			ID:    item.ID,
		}
	}

	selected, err := Select(label+":", selectItems)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.ID == selected.ID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("selected item not found")
}

// ResolveWithArgs resolves a resource either from the provided args (name or ID)
// or by showing the interactive select widget.
func (r *Resolver) ResolveWithArgs(ctx context.Context, resourceType, label string, args []string) (*NamedItem, error) {
	items, err := r.Fetch(ctx, resourceType)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", resourceType, err)
	}

	if len(args) > 0 && args[0] != "" {
		query := strings.ToLower(args[0])
		for _, item := range items {
			if strings.ToLower(item.Name) == query || strings.ToLower(item.ID) == query {
				fmt.Printf("%s %s %s\n", Bold(label+":"), Green(item.Name), Dim("(matched)"))
				return &item, nil
			}
		}
		return nil, fmt.Errorf("no %s matching %q found", resourceType, args[0])
	}

	return r.SelectFromItems(label, items)
}

// ResolveName looks up a name in the cache and returns the UUID, or empty string.
func (r *Resolver) ResolveName(resourceType, name string) string {
	item := r.cache.LookupByName(resourceType, name)
	if item != nil {
		return item.ID
	}
	return ""
}

// ResolveID looks up an ID in the cache and returns the name, or the ID itself.
func (r *Resolver) ResolveID(resourceType, id string) string {
	item := r.cache.LookupByID(resourceType, id)
	if item != nil {
		return item.Name
	}
	return id
}
