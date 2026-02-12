// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Key constants for special keys
const (
	KeyEnter     = '\r'
	KeyNewline   = '\n'
	KeyEscape    = 27
	KeyBackspace = 127
	KeyCtrlC     = 3
	KeyCtrlD     = 4
)

// Special multi-byte key sequences
type SpecialKey int

const (
	KeyNone SpecialKey = iota
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
)

// KeyEvent represents a single keypress
type KeyEvent struct {
	Char    byte
	Special SpecialKey
}

// RawMode enters raw terminal mode and returns a restore function
func RawMode() (restore func(), err error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("entering raw mode: %w", err)
	}
	return func() {
		term.Restore(fd, oldState)
	}, nil
}

// ReadKey reads a single key event from stdin (must be in raw mode)
func ReadKey() (KeyEvent, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return KeyEvent{}, err
	}
	if n == 0 {
		return KeyEvent{}, nil
	}

	// Check for escape sequences (arrow keys)
	if n == 3 && buf[0] == KeyEscape && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return KeyEvent{Special: KeyUp}, nil
		case 'B':
			return KeyEvent{Special: KeyDown}, nil
		case 'C':
			return KeyEvent{Special: KeyRight}, nil
		case 'D':
			return KeyEvent{Special: KeyLeft}, nil
		}
	}

	return KeyEvent{Char: buf[0]}, nil
}

// ANSI escape code helpers

// ClearLine clears the current terminal line
func ClearLine() {
	fmt.Print("\033[2K\r")
}

// ClearDown clears from cursor to end of screen
func ClearDown() {
	fmt.Print("\033[J")
}

// MoveUp moves the cursor up n lines
func MoveUp(n int) {
	if n > 0 {
		fmt.Printf("\033[%dA", n)
	}
}

// MoveDown moves the cursor down n lines
func MoveDown(n int) {
	if n > 0 {
		fmt.Printf("\033[%dB", n)
	}
}

// MoveToColumn moves the cursor to column n (1-based)
func MoveToColumn(n int) {
	fmt.Printf("\033[%dG", n)
}

// HideCursor hides the terminal cursor
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// Bold returns text wrapped in bold ANSI codes
func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

// Dim returns text wrapped in dim ANSI codes
func Dim(s string) string {
	return "\033[2m" + s + "\033[0m"
}

// Reverse returns text with reversed foreground/background
func Reverse(s string) string {
	return "\033[7m" + s + "\033[0m"
}

// Cyan returns text in cyan
func Cyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

// Green returns text in green
func Green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

// Red returns text in red
func Red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

// Yellow returns text in yellow
func Yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}
