// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const maxSuggestions = 6
const maxHistory = 100

// argResourceMap maps command prefixes to the resource type whose names should
// be offered as argument completions.
var argResourceMap = map[string]string{
	"site get":                     "site",
	"site delete":                  "site",
	"vpc get":                      "vpc",
	"vpc delete":                   "vpc",
	"subnet get":                   "subnet",
	"subnet delete":                "subnet",
	"instance get":                 "instance",
	"instance update":              "instance",
	"instance delete":              "instance",
	"allocation get":               "allocation",
	"allocation delete":            "allocation",
	"allocation-constraint list":   "allocation",
	"allocation-constraint get":    "allocation",
	"audit get":                    "audit",
	"machine get":                  "machine",
	"ip-block get":                 "ip-block",
	"ip-block delete":              "ip-block",
	"operating-system get":         "operating-system",
	"operating-system update":      "operating-system",
	"operating-system delete":      "operating-system",
	"ssh-key-group get":            "ssh-key-group",
	"ssh-key get":                  "ssh-key",
	"sku get":                      "sku",
	"rack get":                     "rack",
	"vpc-prefix get":               "vpc-prefix",
	"tenant-account get":           "tenant-account",
	"expected-machine get":         "expected-machine",
	"dpu-extension-service get":    "dpu-extension-service",
	"infiniband-partition get":     "infiniband-partition",
	"nvlink-logical-partition get": "nvlink-logical-partition",
	"network-security-group get":   "network-security-group",
}

var history []string
var historyPos int

// RunREPL starts the interactive REPL loop with inline autocomplete.
func RunREPL(s *Session) error {
	commands := AllCommands()
	cmdNames := make([]string, len(commands))
	cmdMap := make(map[string]Command, len(commands))
	for i, cmd := range commands {
		cmdNames[i] = cmd.Name
		cmdMap[cmd.Name] = cmd
	}
	cmdNames = append(cmdNames, "org", "org list", "org set",
		"scope", "scope site", "scope vpc", "scope clear",
		"exit", "quit")

	fmt.Printf("\n%s\n", Bold("BMM Interactive Mode"))
	fmt.Printf("Org: %s\n", Cyan(s.Org))
	fmt.Printf("Start typing a command. %s to quit.\n\n", Bold("Ctrl+D"))

	for {
		line, err := readLineWithSuggestions(s, cmdNames)
		if err != nil {
			fmt.Println("\nGoodbye.")
			return nil
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if len(history) == 0 || history[len(history)-1] != line {
			history = append(history, line)
			if len(history) > maxHistory {
				history = history[1:]
			}
		}

		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye.")
			return nil
		}

		if line == "org" {
			fmt.Printf("Current org: %s\n\n", Cyan(s.Org))
			continue
		}
		if line == "org list" {
			runOrgList(s)
			fmt.Println()
			continue
		}
		if strings.HasPrefix(line, "org set ") {
			newOrg := strings.TrimSpace(line[len("org set "):])
			if newOrg != "" {
				s.Org = newOrg
				s.Cache.InvalidateAll()
				fmt.Printf("Org set to: %s\n\n", Cyan(s.Org))
			} else {
				fmt.Fprintf(os.Stderr, "%s org name required\n\n", Red("Error:"))
			}
			continue
		}

		if line == "scope" {
			if s.Scope.SiteID == "" && s.Scope.VpcID == "" {
				fmt.Println("No scope set. All list commands return unfiltered results.")
			} else {
				if s.Scope.SiteName != "" {
					fmt.Printf("  site: %s (%s)\n", Cyan(s.Scope.SiteName), s.Scope.SiteID)
				}
				if s.Scope.VpcName != "" {
					fmt.Printf("  vpc:  %s (%s)\n", Cyan(s.Scope.VpcName), s.Scope.VpcID)
				}
			}
			fmt.Println()
			continue
		}
		if line == "scope clear" {
			s.Scope = Scope{}
			s.Cache.InvalidateFiltered()
			fmt.Println("Scope cleared.")
			fmt.Println()
			continue
		}
		if line == "scope site" || strings.HasPrefix(line, "scope site ") {
			runScopeSet(s, "site", strings.TrimSpace(strings.TrimPrefix(line, "scope site")))
			continue
		}
		if line == "scope vpc" || strings.HasPrefix(line, "scope vpc ") {
			runScopeSet(s, "vpc", strings.TrimSpace(strings.TrimPrefix(line, "scope vpc")))
			continue
		}

		if cmd, ok := cmdMap[line]; ok {
			if err := cmd.Run(s, nil); err != nil {
				fmt.Fprintf(os.Stderr, "%s %v\n", Red("Error:"), err)
			}
			fmt.Println()
			continue
		}

		matched := false
		for _, cmd := range commands {
			if strings.HasPrefix(line, cmd.Name) {
				rest := strings.TrimSpace(line[len(cmd.Name):])
				var args []string
				if rest != "" {
					args = strings.Fields(rest)
				}
				if err := cmd.Run(s, args); err != nil {
					fmt.Fprintf(os.Stderr, "%s %v\n", Red("Error:"), err)
				}
				fmt.Println()
				matched = true
				break
			}
		}

		if !matched {
			fmt.Fprintf(os.Stderr, "%s unknown command: %s\n", Red("Error:"), line)
			fmt.Println()
		}
	}
}

func readLineWithSuggestions(s *Session, cmdNames []string) (string, error) {
	restore, err := RawMode()
	if err != nil {
		return "", err
	}
	defer func() {
		restore()
		ShowCursor()
	}()

	prompt := s.PromptString()
	line := ""
	savedLine := ""
	historyPos = -1
	selectedSuggestion := -1
	prevSuggestionCount := 0

	allSuggestions := func() []string {
		return getAllSuggestions(s, line, cmdNames)
	}

	renderInput := func() {
		suggestions := allSuggestions()
		if len(suggestions) > maxSuggestions {
			suggestions = suggestions[:maxSuggestions]
		}

		clearSuggestionLines(prevSuggestionCount)

		ClearLine()
		fmt.Print("\r" + prompt + line)

		if selectedSuggestion >= len(suggestions) {
			selectedSuggestion = len(suggestions) - 1
		}

		if len(line) > 0 && len(suggestions) > 0 {
			for i, sg := range suggestions {
				fmt.Print("\r\n")
				ClearLine()
				if i == selectedSuggestion {
					fmt.Print("  " + Reverse(" "+sg+" "))
				} else {
					fmt.Print("  " + Dim(sg))
				}
			}
			MoveUp(len(suggestions))
			MoveToColumn(len(stripAnsi(prompt)) + len(line) + 1)
		}

		prevSuggestionCount = len(suggestions)
		if len(line) == 0 {
			prevSuggestionCount = 0
		}
	}

	ShowCursor()
	renderInput()

	for {
		key, err := ReadKey()
		if err != nil {
			return "", err
		}

		switch {
		case key.Char == KeyCtrlC:
			line = ""
			selectedSuggestion = -1
			historyPos = -1
			clearSuggestionLines(prevSuggestionCount)
			prevSuggestionCount = 0
			renderInput()

		case key.Char == KeyCtrlD:
			clearSuggestionLines(prevSuggestionCount)
			return "", fmt.Errorf("EOF")

		case key.Char == KeyEnter || key.Char == KeyNewline:
			suggestions := allSuggestions()
			if len(suggestions) > maxSuggestions {
				suggestions = suggestions[:maxSuggestions]
			}

			if selectedSuggestion >= 0 && selectedSuggestion < len(suggestions) {
				line = suggestions[selectedSuggestion]
				selectedSuggestion = -1
				historyPos = -1
				clearSuggestionLines(prevSuggestionCount)
				prevSuggestionCount = 0
				renderInput()
				continue
			}

			clearSuggestionLines(prevSuggestionCount)
			ClearLine()
			fmt.Print("\r" + prompt + line + "\r\n")
			historyPos = -1
			return line, nil

		case key.Char == '\t':
			suggestions := allSuggestions()
			if len(suggestions) > 0 {
				idx := selectedSuggestion
				if idx < 0 {
					idx = 0
				}
				if idx < len(suggestions) {
					line = suggestions[idx]
					selectedSuggestion = -1
				}
			}
			renderInput()

		case key.Special == KeyUp:
			suggestions := allSuggestions()
			if len(suggestions) > maxSuggestions {
				suggestions = suggestions[:maxSuggestions]
			}

			if len(line) > 0 && len(suggestions) > 0 {
				selectedSuggestion--
				if selectedSuggestion < 0 {
					selectedSuggestion = len(suggestions) - 1
				}
				renderInput()
				continue
			}

			if len(history) > 0 {
				if historyPos == -1 {
					savedLine = line
					historyPos = len(history) - 1
				} else if historyPos > 0 {
					historyPos--
				}
				line = history[historyPos]
				selectedSuggestion = -1
			}
			renderInput()

		case key.Special == KeyDown:
			suggestions := allSuggestions()
			if len(suggestions) > maxSuggestions {
				suggestions = suggestions[:maxSuggestions]
			}

			if len(line) > 0 && len(suggestions) > 0 {
				selectedSuggestion++
				if selectedSuggestion >= len(suggestions) {
					selectedSuggestion = 0
				}
				renderInput()
				continue
			}

			if historyPos >= 0 {
				historyPos++
				if historyPos >= len(history) {
					historyPos = -1
					line = savedLine
				} else {
					line = history[historyPos]
				}
				selectedSuggestion = -1
			}
			renderInput()

		case key.Char == KeyBackspace:
			if len(line) > 0 {
				line = line[:len(line)-1]
				selectedSuggestion = -1
				historyPos = -1
			}
			renderInput()

		case key.Char >= 32 && key.Char < 127:
			line += string(key.Char)
			selectedSuggestion = -1
			historyPos = -1
			renderInput()

		default:
			continue
		}
	}
}

func getAllSuggestions(s *Session, input string, cmdNames []string) []string {
	if input == "" {
		return nil
	}

	for cmdPrefix, resourceType := range argResourceMap {
		withSpace := cmdPrefix + " "
		if strings.HasPrefix(strings.ToLower(input), strings.ToLower(withSpace)) {
			argPart := input[len(withSpace):]
			return getResourceSuggestions(s, cmdPrefix, resourceType, argPart)
		}
	}

	return getCommandSuggestions(input, cmdNames)
}

func getResourceSuggestions(s *Session, cmdPrefix, resourceType, argFilter string) []string {
	items := s.Cache.Get(resourceType)
	if items == nil {
		fetched, err := s.Resolver.Fetch(s.Ctx, resourceType)
		if err != nil {
			return nil
		}
		items = fetched
	}

	lowerFilter := strings.ToLower(argFilter)
	var matches []string
	for _, item := range items {
		name := item.Name
		if name == "" {
			name = item.ID
		}
		if lowerFilter == "" || strings.Contains(strings.ToLower(name), lowerFilter) {
			matches = append(matches, cmdPrefix+" "+name)
		}
	}
	return matches
}

func getCommandSuggestions(input string, cmdNames []string) []string {
	lower := strings.ToLower(input)
	var matches []string
	for _, name := range cmdNames {
		if strings.HasPrefix(strings.ToLower(name), lower) {
			matches = append(matches, name)
		}
	}
	return matches
}

func clearSuggestionLines(count int) {
	if count == 0 {
		return
	}
	for i := 0; i < count; i++ {
		fmt.Print("\r\n")
		ClearLine()
	}
	MoveUp(count)
}

func stripAnsi(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func runScopeSet(s *Session, resourceType, nameOrID string) {
	currentSiteID := s.Scope.SiteID
	currentSiteName := s.Scope.SiteName
	if currentSiteName == "" {
		currentSiteName = currentSiteID
	}

	if resourceType == "vpc" {
		s.Cache.Invalidate("vpc")
	}

	var item *NamedItem
	var err error

	if nameOrID != "" {
		item, err = resolveByNameOrID(s, resourceType, nameOrID)
	} else {
		item, err = s.Resolver.Resolve(s.Ctx, resourceType, strings.Title(resourceType))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n\n", Red("Error:"), err)
		return
	}

	switch resourceType {
	case "site":
		s.Scope.SiteID = item.ID
		s.Scope.SiteName = item.Name
		s.Scope.VpcID = ""
		s.Scope.VpcName = ""
		s.Cache.InvalidateFiltered()
	case "vpc":
		if currentSiteID != "" && item.Extra != nil && item.Extra["siteId"] != "" && item.Extra["siteId"] != currentSiteID {
			fmt.Fprintf(
				os.Stderr,
				"%s selected VPC belongs to a different site (%s). Current site scope is %s. Use %s first or choose a VPC in the current site.\n\n",
				Red("Error:"),
				Cyan(s.Resolver.ResolveID("site", item.Extra["siteId"])),
				Cyan(currentSiteName),
				Bold("scope clear"),
			)
			return
		}

		s.Scope.VpcID = item.ID
		s.Scope.VpcName = item.Name
		if item.Extra != nil && item.Extra["siteId"] != "" {
			siteID := item.Extra["siteId"]
			if currentSiteID == "" && s.Scope.SiteID != siteID {
				siteName := s.Resolver.ResolveID("site", siteID)
				s.Scope.SiteID = siteID
				s.Scope.SiteName = siteName
				fmt.Printf("Scope set: site = %s (from VPC)\n", Cyan(siteName))
			}
		}
		s.Cache.InvalidateFiltered()
	}

	fmt.Printf("Scope set: %s = %s\n\n", resourceType, Cyan(item.Name))
}

func resolveByNameOrID(s *Session, resourceType, query string) (*NamedItem, error) {
	items, err := s.Resolver.Fetch(s.Ctx, resourceType)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(query)
	for _, item := range items {
		if strings.ToLower(item.Name) == lower || strings.ToLower(item.ID) == lower {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no %s matching %q", resourceType, query)
}

func runOrgList(s *Session) {
	if s.Token == "" {
		fmt.Printf("Current org: %s\n", Cyan(s.Org))
		fmt.Printf("%s No token available. Run %s first.\n", Yellow("Note:"), Bold("login"))
		return
	}

	orgs := extractOrgsFromJWT(s.Token)
	if len(orgs) == 0 {
		fmt.Printf("Current org: %s\n", Cyan(s.Org))
		fmt.Println("Could not extract orgs from token.")
		fmt.Printf("Switch manually: %s\n", Bold("org set <org-name>"))
		return
	}

	fmt.Printf("Current org: %s\n\n", Cyan(s.Org))
	for _, org := range orgs {
		marker := "  "
		if org == s.Org {
			marker = Cyan("> ")
		}
		fmt.Printf("%s%s\n", marker, org)
	}
	fmt.Printf("\nSwitch with: %s\n", Bold("org set <org-name>"))
}

type jwtAccessClaim struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func extractOrgsFromJWT(tokenStr string) []string {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil
	}

	payload := parts[1]
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil
	}

	var claims struct {
		Access []jwtAccessClaim `json:"access"`
	}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil
	}

	var orgs []string
	seen := map[string]bool{}
	for _, c := range claims.Access {
		if strings.HasPrefix(c.Type, "group/ngc") && c.Name != "" && !seen[c.Name] {
			orgs = append(orgs, c.Name)
			seen[c.Name] = true
		}
	}
	return orgs
}
