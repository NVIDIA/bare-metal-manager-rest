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
package dumper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	summaryFileName   = "rack_summary.json"
	rackFileName      = "rack_%s.json"
	nvlDomainFileName = "nvl_domain_%s.json"
)

const (
	SourceTypeUnknown SourceType = "Unknown"
	SourceTypeExcel   SourceType = "Excel"
	SourceTypeCerebro SourceType = "Cerebro"
	SourceTypeRLA     SourceType = "RLA"
)

type SourceType string

type Config struct {
	DirName    string
	ReadOnly   bool
	SourceType SourceType
}

func (c *Config) ValidateOrLoad() (map[string]bool, map[string]bool, error) {
	if c.DirName == "" {
		return nil, nil, errors.New("directory name is required")
	}

	if !c.ReadOnly {
		if c.SourceType == "" || c.SourceType == SourceTypeUnknown {
			// For the writeable dumper, the source type is required.
			return nil, nil, errors.New("source type is required")
		}

		// Make sure the directory gets created.
		if err := os.MkdirAll(c.DirName, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create directory %s: %w", c.DirName, err)
		}

		return make(map[string]bool), make(map[string]bool), nil
	}

	if _, err := os.Stat(c.DirName); err != nil {
		return nil, nil, fmt.Errorf("directory %s is not accessible", c.DirName)
	}

	data, err := os.ReadFile(c.SummaryFilePath())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read summary file: %w", err)
	}

	summary, err := NewSummary(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build summary: %w", err)
	}

	racks := make(map[string]bool)
	for _, name := range summary.RackNames {
		if _, err := os.Stat(c.RackFilePath(name)); err != nil {
			return nil, nil, fmt.Errorf("rack %s is not accessible", name)
		}
		racks[name] = true
	}

	nvlDomains := make(map[string]bool)
	for _, name := range summary.NVLDomainNames {
		if _, err := os.Stat(c.NVLDomainFilePath(name)); err != nil {
			return nil, nil, fmt.Errorf("nvl domain %s is not accessible", name)
		}
		nvlDomains[name] = true
	}

	c.SourceType = summary.SourceType

	return racks, nvlDomains, nil
}

func (c *Config) RackFilePath(name string) string {
	return filepath.Join(c.DirName, fmt.Sprintf(rackFileName, name))
}

func (c *Config) NVLDomainFilePath(name string) string {
	return filepath.Join(c.DirName, fmt.Sprintf(nvlDomainFileName, name))
}

func (c *Config) SummaryFilePath() string {
	return filepath.Join(c.DirName, summaryFileName)
}
