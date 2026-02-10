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
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
)

type Dumper struct {
	conf       Config
	mtx        sync.RWMutex
	racks      map[string]bool
	nvlDomains map[string]bool
}

func New(conf Config) (*Dumper, error) {
	log.Info().Msgf("Create a new dumper with configuration [%v]", conf)

	racks, nvlDomains, err := conf.ValidateOrLoad()
	if err != nil {
		return nil, err
	}

	return &Dumper{conf: conf, racks: racks, nvlDomains: nvlDomains}, nil
}

func (d *Dumper) Done() error {
	log.Info().Msg("closing the dumper")

	if d.conf.ReadOnly {
		log.Debug().Msg("dumper is read only")
		return nil
	}

	log.Info().Msg("writing the summary file")

	summary := d.Summary()

	bytes, err := summary.ToJSON()
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal the rack summary")
		return err
	}

	if err := os.WriteFile(d.conf.SummaryFilePath(), bytes, 0644); err != nil {
		log.Error().Err(err).Msg("failed to write summary file")
		return err
	}

	return nil
}

func (d *Dumper) Dump(data any) error {
	if d.conf.ReadOnly {
		return fmt.Errorf("dumper is read only")
	}

	var filePath string

	switch obj := data.(type) {
	case *rack.Rack:
		if !d.addRackName(obj.Info.Name) {
			return fmt.Errorf("rack %s already exists", obj.Info.Name)
		}
		filePath = d.conf.RackFilePath(obj.Info.Name)
	case *nvldomain.NVLDomain:
		if !d.addNVLDomainName(obj.Identifier.Name) {
			return fmt.Errorf(
				"nvl domain %s already exists",
				obj.Identifier.Name,
			)
		}
		filePath = d.conf.NVLDomainFilePath(obj.Identifier.Name)
	default:
		return fmt.Errorf("unknown data type: %T", data)
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

func (d *Dumper) Load(name string, to any) error {
	if !d.conf.ReadOnly {
		return fmt.Errorf("dumper is not read only")
	}

	var filePath string

	switch to.(type) {
	case *rack.Rack:
		filePath = d.conf.RackFilePath(name)
	case *nvldomain.NVLDomain:
		filePath = d.conf.NVLDomainFilePath(name)
	default:
		return fmt.Errorf("unknown data type: %T", to)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read from file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		return fmt.Errorf("failed to unmarshal from file %s: %w", filePath, err)
	}

	return nil
}

func (d *Dumper) Summary() *Summary {
	return &Summary{
		SourceType:     d.conf.SourceType,
		RackNames:      d.allRackNames(),
		NVLDomainNames: d.allNVLDomainNames(),
	}
}

func (d *Dumper) allRackNames() []string {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	all := make([]string, 0, len(d.racks))
	for name := range d.racks {
		all = append(all, name)
	}

	return all
}

func (d *Dumper) allNVLDomainNames() []string {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	all := make([]string, 0, len(d.nvlDomains))
	for name := range d.nvlDomains {
		all = append(all, name)
	}

	return all
}

func (d *Dumper) addRackName(name string) bool {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.racks[name]; ok {
		return false
	}

	d.racks[name] = true
	return true
}

func (d *Dumper) addNVLDomainName(name string) bool {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.nvlDomains[name]; ok {
		return false
	}

	d.nvlDomains[name] = true
	return true
}
