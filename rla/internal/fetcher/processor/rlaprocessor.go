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

	"github.com/rs/zerolog/log"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
)

type RLAProcessor struct {
	rlaClient *client.Client
}

func newRLAProcessor(c client.Config) (*RLAProcessor, error) {
	rlaClient, err := client.New(c)
	if err != nil {
		return nil, err
	}

	return &RLAProcessor{
		rlaClient: rlaClient,
	}, nil
}

func (rp *RLAProcessor) Process(ctx context.Context, item *ProcessableItem) error {
	for _, r := range item.Racks {
		if r == nil {
			continue
		}

		_, err := rp.rlaClient.CreateExpectedRack(ctx, r)
		if err != nil {
			log.Error().Err(err).Msgf("failed to create expected rack %s", r.Info.Name)
			return err
		}
	}

	if item.NVLDomain != nil {
		_, err := rp.rlaClient.CreateNVLDomain(ctx, item.NVLDomain)
		if err != nil {
			log.Error().Err(err).Msg("failed to create NVL domain")
			return err
		}

		err = rp.rlaClient.AttachRacksToNVLDomain(
			ctx,
			item.NVLDomain.Identifier,
			item.NVLDomain.RackIdentifiers,
		)

		if err != nil {
			log.Error().Err(err).Msg("failed to attach racks to NVL domain")
			return err
		}
	}

	return nil
}

func (rp *RLAProcessor) Type() Type {
	return TypeRLA
}

func (rp *RLAProcessor) Done() error {
	return rp.rlaClient.Close()
}
