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

	"github.com/nvidia/carbide-rest/rla/internal/dumper"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
)

type Type string

const (
	TypeDryRun Type = "DryRun"
	TypeDumper Type = "Dumper"
	TypeRLA    Type = "RLA"
)

type Config struct {
	DryRun        bool
	DumperConf    dumper.Config
	RLAClientConf client.Config
}

type Processor interface {
	Process(ctx context.Context, item *ProcessableItem) error
	Type() Type
	Done() error
}

func New(c Config) (Processor, error) {
	if c.DryRun {
		// Create a dry run processor.
		return newDryruner()
	}

	if c.DumperConf.DirName != "" {
		// Create a dumper processor
		return newDumper(c.DumperConf)
	}

	return newRLAProcessor(c.RLAClientConf)
}
