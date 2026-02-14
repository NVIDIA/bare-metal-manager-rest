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
package cerebro

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/clients/cerebro"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/dumper"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/fetcher/processor"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/client"
	identifier "github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/Identifier"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/workerpool"
)

type Config struct {
	URL               string
	APIKey            string
	RLAClientConf     client.Config
	CerebroClientConf cerebro.Config
	DryRun            bool
	DumperDir         string
}

type Fetcher struct {
	conf          Config
	cerebroClient *cerebro.Client
	processor     processor.Processor
}

func New(ctx context.Context, conf Config) (*Fetcher, error) {
	cerebroClient, err := cerebro.New(ctx, conf.CerebroClientConf)
	if err != nil {
		return nil, err
	}

	processor, err := processor.New(
		processor.Config{
			DryRun: conf.DryRun,
			DumperConf: dumper.Config{
				DirName:    conf.DumperDir,
				ReadOnly:   false,
				SourceType: dumper.SourceTypeCerebro,
			},
			RLAClientConf: conf.RLAClientConf,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		conf:          conf,
		processor:     processor,
		cerebroClient: cerebroClient,
	}, nil
}

func (f *Fetcher) Fetch(ctx context.Context, names []string) int {
	fetchByRackName := true

	ids := make([]identifier.Identifier, 0, len(names))
	if len(names) == 0 {
		// Fetch all the NVL domains, then fetch the rack for each NVL domain.
		nvlDomainIDs, err := f.cerebroClient.GetAllNVLDomainIdentifiers(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get all NVL domain IDs")
			return 0
		}

		if len(nvlDomainIDs) == 0 {
			log.Info().Msg("no NVL domain IDs found")
			return 0
		}

		ids = append(ids, nvlDomainIDs...)

		fetchByRackName = false
	} else {
		for _, name := range names {
			ids = append(ids, identifier.Identifier{
				ID:   uuid.Nil,
				Name: name,
			})
		}
	}

	start := time.Now()

	// Experiments show that the performance is better with 12 workers.
	wp := workerpool.New(
		&workerpool.Config{
			MaxWorkers:    12,
			QueueSize:     len(ids),
			WorkerTimeout: 5 * time.Second,
			TaskTimeout:   60 * time.Second,
			EnableMetrics: false,
		},
	)

	if err := wp.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start worker pool")
		return 0
	}
	defer wp.Stop()

	var dryRunProcessorWG sync.WaitGroup
	var fetchTaskWG sync.WaitGroup
	dryRunProcessChan := make(chan *processor.ProcessableItem)
	var processed int32 = 0

	dryRunProcessorWG.Add(1)
	go func() {
		log.Debug().Msg("starting dry run processor")
		defer dryRunProcessorWG.Done()

		for fetched := range dryRunProcessChan {
			f.processor.Process(ctx, fetched)
		}

		log.Debug().Msg("dry run processor finished")
	}()

	for _, id := range ids {
		task := &fetchTask{
			id:                id,
			fetchByRackName:   fetchByRackName,
			fetcher:           f,
			wg:                &fetchTaskWG,
			dryRunProcessChan: dryRunProcessChan,
			processed:         &processed,
		}
		task.wg.Add(1)
		if err := wp.Submit(task); err != nil {
			log.Error().Err(err).Msg("failed to submit task")
			task.wg.Done()
		}
	}

	fetchTaskWG.Wait()
	close(dryRunProcessChan)
	dryRunProcessorWG.Wait()

	log.Info().Msgf("fetched %d racks in %s", processed, time.Since(start))

	return int(processed)
}

func (f *Fetcher) Done() {
	if f.processor == nil {
		return
	}

	if err := f.processor.Done(); err != nil {
		log.Error().Err(err).Msgf("failed to close processor %s", f.processor.Type())
	}
}

func fixRackInfo(rack *rack.Rack) {
	if rack.Info.Manufacturer == "" {
		// Try to use the manufacturer of the compute component.
		for _, c := range rack.Components {
			if c.Type == devicetypes.ComponentTypeCompute && c.Info.Manufacturer != "" {
				rack.Info.Manufacturer = c.Info.Manufacturer
				break
			}
		}

		if rack.Info.Manufacturer == "" {
			rack.Info.Manufacturer = "Unknown"
		}
	}

	if rack.Info.SerialNumber == "" {
		// Use the name as the serial number.
		rack.Info.SerialNumber = rack.Info.Name
	}
}

type fetchTask struct {
	id                identifier.Identifier
	fetchByRackName   bool
	fetcher           *Fetcher
	wg                *sync.WaitGroup
	dryRunProcessChan chan *processor.ProcessableItem
	processed         *int32
}

func (ft *fetchTask) ID() string {
	return ft.id.Name
}

func (ft *fetchTask) Execute(ctx context.Context) error {
	defer ft.wg.Done()

	var fetchedItem *processor.ProcessableItem

	if ft.fetchByRackName {
		log.Debug().Msgf("fetching rack %s", ft.id.Name)
		r, err := ft.fetcher.cerebroClient.GetRackByName(ctx, ft.id.Name)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get rack %s", ft.id.Name)
			return err
		}

		fetchedItem = processor.NewProcessableItem(nil, []*rack.Rack{r})
	} else {
		log.Debug().Msgf("fetching rack by NVL domain ID %s", ft.id.ID)
		racks, err := ft.fetcher.cerebroClient.GetRackByNVLDomainID(
			ctx,
			ft.id.ID,
		)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get rack by NVL domain ID %s", ft.id.ID)
			return err
		}

		rackIdentifiers := make([]identifier.Identifier, 0, len(racks))
		for _, r := range racks {
			rackIdentifiers = append(
				rackIdentifiers,
				identifier.Identifier{
					ID:   r.Info.ID,
					Name: r.Info.Name,
				},
			)
		}

		fetchedItem = processor.NewProcessableItem(
			&nvldomain.NVLDomain{
				Identifier:      ft.id,
				RackIdentifiers: rackIdentifiers,
			},
			racks,
		)
	}

	for _, r := range fetchedItem.Racks {
		fixRackInfo(r)
	}

	if ft.fetcher.processor.Type() == processor.TypeDryRun {
		ft.dryRunProcessChan <- fetchedItem
		atomic.AddInt32(ft.processed, 1)
		return nil
	}

	if err := ft.fetcher.processor.Process(ctx, fetchedItem); err != nil {
		log.Error().Err(err).Msgf("failed to process fectched item %+v", fetchedItem) //nolint
		return err
	}

	atomic.AddInt32(ft.processed, 1)

	log.Debug().Msgf("fetched item %v", fetchedItem)

	return nil
}
