// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptText displays a label and reads a line of text input.
// If required is true and the user enters empty text, it re-prompts.
func PromptText(label string, required bool) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", Bold(label))
		if !scanner.Scan() {
			return "", fmt.Errorf("input cancelled")
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" && required {
			fmt.Println(Red("  (required)"))
			continue
		}
		return text, nil
	}
}

// PromptConfirm displays a y/N confirmation prompt.
func PromptConfirm(label string) (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s [y/N] ", Bold(label))
	if !scanner.Scan() {
		return false, fmt.Errorf("input cancelled")
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
