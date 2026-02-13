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
package builder

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/dumper"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/fetcher/processor"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
)

// Builder is responsible for building the rack inventory from the given
// multiple sources. Currently, the builder only supports the Excel and
// Cerebro sources.
type Builder struct {
	conf      Config
	sources   map[dumper.SourceType]*dumper.Dumper
	summaries map[dumper.SourceType]*dumper.Summary
	processor processor.Processor
}

func New(conf Config) (*Builder, error) {
	log.Info().Msgf("Create a new builder with configuration [%v]", conf)

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	builder := &Builder{
		conf:      conf,
		sources:   make(map[dumper.SourceType]*dumper.Dumper),
		summaries: make(map[dumper.SourceType]*dumper.Summary),
	}

	for _, dir := range conf.SourceDumperDirs {
		dumper, err := dumper.New(dumper.Config{DirName: dir, ReadOnly: true})
		if err != nil {
			return nil, err
		}

		summary := dumper.Summary()
		builder.sources[summary.SourceType] = dumper
		builder.summaries[summary.SourceType] = summary
	}

	if !builder.validateSources() {
		return nil, fmt.Errorf("unexpected sources")
	}

	processor, err := processor.New(
		processor.Config{
			DryRun: conf.DryRun,
			DumperConf: dumper.Config{
				DirName:    conf.OutputDumperDir,
				ReadOnly:   false,
				SourceType: dumper.SourceTypeRLA,
			},
			RLAClientConf: conf.RLAClientConf,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create processor")
		return nil, err
	}

	builder.processor = processor

	return builder, nil
}

func (b *Builder) Done() {
	for _, source := range b.sources {
		if err := source.Done(); err != nil {
			log.Error().Err(err).Msg("failed to close source dumper")
		}
	}

	if b.processor != nil {
		if err := b.processor.Done(); err != nil {
			log.Error().Err(err).Msg("failed to close processor")
		}
	}
}

func (b *Builder) Build(ctx context.Context) []int {
	processed := make([]int, 2)

	// Use the rack names from the Excel source as the source of truth.
	// Ideally, we should use the Cerebro as the SOT, unfortunately Cerebro
	// does not provide enough information to serve as the SOT.
	for _, name := range b.summaries[dumper.SourceTypeExcel].RackNames {
		var rackExcel rack.Rack
		var rackCerebro rack.Rack

		err := b.sources[dumper.SourceTypeExcel].Load(name, &rackExcel)
		if err != nil {
			log.Error().Err(err).Msgf(
				"failed to load rack %s from Excel source",
				name,
			)
			continue
		}

		err = b.sources[dumper.SourceTypeCerebro].Load(name, &rackCerebro)
		if err != nil {
			log.Error().Err(err).Msgf(
				"failed to load rack %s from Cerebro source",
				name,
			)
			continue
		}

		builder := newRackBuilder(&rackExcel, &rackCerebro)
		r := builder.build()
		for _, r := range builder.reports {
			log.Debug().Msgf("Rack %s: %s", name, r.String())
		}

		item := processor.NewProcessableItem(nil, []*rack.Rack{r})
		if err := b.processor.Process(ctx, item); err != nil {
			log.Error().Err(err).Msgf(
				"failed to process rack %s",
				name,
			)
			continue
		}

		processed[0]++
	}

	// NVL domain is only available from the Cerebro source.
	for _, name := range b.summaries[dumper.SourceTypeCerebro].NVLDomainNames {
		var nvlDomain nvldomain.NVLDomain

		err := b.sources[dumper.SourceTypeCerebro].Load(name, &nvlDomain)
		if err != nil {
			log.Error().Err(err).Msgf(
				"failed to load NVL domain %s from Cerebro source",
				name,
			)
			continue
		}

		item := processor.NewProcessableItem(&nvlDomain, nil)
		if err := b.processor.Process(ctx, item); err != nil {
			log.Error().Err(err).Msgf(
				"failed to process NVL domain %s",
				name,
			)
			continue
		}
		processed[1]++
	}

	return processed
}

func (b *Builder) validateSources() bool {
	return b.summaries[dumper.SourceTypeExcel] != nil &&
		b.summaries[dumper.SourceTypeCerebro] != nil &&
		b.sources[dumper.SourceTypeExcel] != nil &&
		b.sources[dumper.SourceTypeCerebro] != nil
}
