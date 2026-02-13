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
package processor

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/dumper"
)

type Dumper struct {
	dumper *dumper.Dumper
}

func newDumper(c dumper.Config) (*Dumper, error) {
	if c.DirName == "" {
		return nil, errors.New("dirName is required")
	}

	if c.SourceType == "" || c.SourceType == dumper.SourceTypeUnknown {
		return nil, errors.New("source type is required")
	}

	// For the dumper for processing the fetced information, it should not be
	// read-only.
	c.ReadOnly = false

	dumper, err := dumper.New(c)
	if err != nil {
		return nil, err
	}

	return &Dumper{dumper: dumper}, nil
}

func (d *Dumper) Process(ctx context.Context, item *ProcessableItem) error {
	for _, r := range item.Racks {
		if r == nil {
			continue
		}

		if err := d.dumper.Dump(r); err != nil {
			log.Error().Err(err).Msgf("failed to dump rack %s", r.Info.Name)
			return err
		}
	}

	if item.NVLDomain != nil {
		if err := d.dumper.Dump(item.NVLDomain); err != nil {
			log.Error().Err(err).Msgf("failed to dump NVL domain %s", item.NVLDomain.Identifier.Name)
			return err
		}
	}

	return nil
}

func (d *Dumper) Type() Type {
	return TypeDumper
}

func (d *Dumper) Done() error {
	return d.dumper.Done()
}
