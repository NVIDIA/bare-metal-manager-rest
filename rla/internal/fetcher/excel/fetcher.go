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
package excel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/dumper"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/fetcher/processor"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/client"
)

type Info struct {
	File  string
	Sheet string
}

type Config struct {
	RackExcelInfo Info
	CompExcelInfo Info
	DumperDir     string
	RLAClientConf client.Config
	DryRun        bool
}

type Fetcher struct {
	conf      Config
	processor processor.Processor
}

func New(conf Config) (*Fetcher, error) {
	log.Info().Msgf("Create a new fetcher with configuration [%v]", conf)

	processor, err := processor.New(
		processor.Config{
			DryRun: conf.DryRun,
			DumperConf: dumper.Config{
				DirName:    conf.DumperDir,
				ReadOnly:   false,
				SourceType: dumper.SourceTypeExcel,
			},
			RLAClientConf: conf.RLAClientConf,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create processor")
		return nil, err
	}

	return &Fetcher{conf: conf, processor: processor}, nil
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	racks, err := f.buildRackInventory()
	if err != nil {
		return fmt.Errorf("failed to build rack invenotry: %v", err)
	}

	components, err := f.buildComponentInventory()
	if err != nil {
		return fmt.Errorf("failed to build comp invenotry: %v", err)
	}

	for sn, comps := range components {
		if rack, ok := racks[sn]; ok {
			patched := rack.PatchComponents(comps)
			log.Info().Msgf(
				"Patched %d components for rack with serial number %s",
				patched,
				sn,
			)
		} else {
			log.Info().Msgf("No rack with serial number %s\n", sn)
		}
	}

	return f.processResult(ctx, racks)
}

func (f *Fetcher) Done() {
	if f.processor == nil {
		return
	}

	if err := f.processor.Done(); err != nil {
		log.Error().Err(err).Msgf("failed to close processor %s", f.processor.Type()) //nolint
	}
}

func (f *Fetcher) buildRackInventory() (map[string]*rack.Rack, error) {
	racks := make(map[string]*rack.Rack)

	rows, err := f.conf.RackExcelInfo.read()
	if err != nil {
		return racks, err
	}

	if len(rows) == 0 {
		return racks, nil
	}

	indexes := buildMapFromStructToExcel(rows[0], RackInvetory{})
	for _, row := range rows[1:] {
		var inv RackInvetory
		buildInventory(indexes, &inv, row)

		c := inv.componentInfo()
		if r, ok := racks[inv.RackSerialNumber]; !ok {
			rname := positionToRackName(inv.RackPosition, 5)
			rack := rack.New(inv.rackInfo(rname), inv.locationInfo())
			rack.AddComponent(c)
			racks[inv.RackSerialNumber] = rack
		} else {
			r.AddComponent(c)
		}
	}

	return racks, nil
}

func (f *Fetcher) buildComponentInventory() (map[string][]component.Component, error) {
	components := make(map[string][]component.Component)

	rows, err := f.conf.CompExcelInfo.read()
	if err != nil {
		return components, err
	}

	if len(rows) == 0 {
		return components, nil
	}

	indexes := buildMapFromStructToExcel(rows[0], ComponentInvetory{})
	for _, row := range rows[1:] {
		var inv ComponentInvetory
		buildInventory(indexes, &inv, row)

		components[inv.ParentTag] = append(
			components[inv.ParentTag],
			inv.componentInfo(),
		)
	}

	return components, nil
}

func (f *Fetcher) processResult(
	ctx context.Context,
	racks map[string]*rack.Rack,
) error {
	for _, r := range racks {
		r.Seal()
		fetchedItem := processor.NewProcessableItem(nil, []*rack.Rack{r})
		if err := f.processor.Process(ctx, fetchedItem); err != nil {
			log.Error().Err(err).Msgf("failed to process fetched item %+v", fetchedItem) //nolint
			return err
		}
	}

	if f.processor.Type() == processor.TypeDryRun {
		fmt.Printf("=========================================")
		fmt.Printf("Total racks: %d\n", len(racks))
	}

	return nil
}

func (ei *Info) read() ([][]string, error) {
	file, err := excelize.OpenFile(ei.File)
	if err != nil {
		return nil, err
	}

	return file.GetRows(ei.Sheet)
}

func buildMapFromStructToExcel(titles []string, inv any) []int {
	m := make(map[string]int)
	for idx, title := range titles {
		m[title] = idx
	}

	t := reflect.TypeOf(inv)
	indexes := make([]int, t.NumField())
	for i := range indexes {
		if idx, ok := m[t.Field(i).Tag.Get("title")]; ok {
			indexes[i] = idx
		} else {
			indexes[i] = -1
		}
	}

	return indexes
}

func buildInventory(indexes []int, s any, info []string) {
	v := reflect.ValueOf(s).Elem()
	if len(indexes) != v.NumField() {
		// indexes is not correct, do nothing
		return
	}

	for i, idx := range indexes {
		field := v.Field(i)
		if field.CanSet() {
			value := retrieveInfoWithIndex(idx, info)
			field.Set(reflect.ValueOf(value))
		}
	}
}

func retrieveInfoWithIndex(idx int, info []string) string {
	if idx < 0 || len(info) < idx+1 {
		return ""
	}

	return info[idx]
}
