// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"fmt"
	"strings"
)

// SelectItem represents one option in a select menu
type SelectItem struct {
	Label string
	ID    string
	Extra map[string]string
}

// Select displays an interactive arrow-key menu and returns the selected item.
// The user navigates with up/down arrows, type to filter, Enter to select.
// Returns nil and an error on Ctrl+C / Ctrl+D.
func Select(label string, items []SelectItem) (*SelectItem, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no items to select from")
	}

	restore, err := RawMode()
	if err != nil {
		return nil, err
	}
	defer restore()

	cursor := 0
	filter := ""
	filtered := filterItems(items, filter)

	// Initial render
	render(label, filtered, cursor, filter)

	for {
		key, err := ReadKey()
		if err != nil {
			ShowCursor()
			return nil, err
		}

		switch {
		case key.Char == KeyCtrlC || key.Char == KeyCtrlD:
			// Clean up display and return
			clearRendered(len(filtered))
			ShowCursor()
			return nil, fmt.Errorf("selection cancelled")

		case key.Special == KeyUp:
			if cursor > 0 {
				cursor--
			}

		case key.Special == KeyDown:
			if cursor < len(filtered)-1 {
				cursor++
			}

		case key.Char == KeyEnter || key.Char == KeyNewline:
			if len(filtered) > 0 {
				selected := filtered[cursor]
				clearRendered(len(filtered))
				ShowCursor()
				// Print the selection on one line
				fmt.Printf("%s %s\r\n", label, Green(selected.Label))
				return &selected, nil
			}

		case key.Char == KeyBackspace:
			if len(filter) > 0 {
				filter = filter[:len(filter)-1]
				filtered = filterItems(items, filter)
				if cursor >= len(filtered) {
					cursor = max(0, len(filtered)-1)
				}
			}

		case key.Char >= 32 && key.Char < 127:
			// Printable character -- add to filter
			filter += string(key.Char)
			filtered = filterItems(items, filter)
			if cursor >= len(filtered) {
				cursor = max(0, len(filtered)-1)
			}

		default:
			continue
		}

		// Re-render
		clearRendered(len(filtered))
		filtered = filterItems(items, filter)
		if cursor >= len(filtered) {
			cursor = max(0, len(filtered)-1)
		}
		render(label, filtered, cursor, filter)
	}
}

func render(label string, items []SelectItem, cursor int, filter string) {
	HideCursor()

	// Header line with filter
	if filter != "" {
		fmt.Printf("%s %s\r\n", Bold(label), Dim("(filter: "+filter+")"))
	} else {
		fmt.Printf("%s %s\r\n", Bold(label), Dim("(type to filter, arrows to move, enter to select)"))
	}

	// Menu items
	for i, item := range items {
		if i == cursor {
			fmt.Printf("  %s %s\r\n", Cyan(">"), Reverse(" "+item.Label+" "))
		} else {
			fmt.Printf("    %s\r\n", item.Label)
		}
	}

	if len(items) == 0 {
		fmt.Printf("    %s\r\n", Dim("(no matches)"))
	}
}

func clearRendered(itemCount int) {
	// Move up past all rendered lines (header + items or "no matches")
	lines := 1 + itemCount
	if itemCount == 0 {
		lines = 2 // header + "(no matches)"
	}
	MoveUp(lines)
	MoveToColumn(1)
	ClearDown()
}

func filterItems(items []SelectItem, filter string) []SelectItem {
	if filter == "" {
		return items
	}
	lowerFilter := strings.ToLower(filter)
	var result []SelectItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Label), lowerFilter) {
			result = append(result, item)
		}
	}
	return result
}
