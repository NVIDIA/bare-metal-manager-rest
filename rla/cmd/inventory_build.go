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
package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/builder"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/client"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the rack inventory",
		Long:  `Build the rack inventory from the files which are fetched from various sources`,

		Run: func(cmd *cobra.Command, args []string) {
			doBuild()
		},
	}

	sourceDumperDirs              []string
	buildInventoryDryRun          bool
	buildInventoryRLAHost         string
	buildInventoryRLAPort         int
	buildInventoryOutputDumperDir string
)

func init() {
	inventoryCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringArrayVarP(
		&sourceDumperDirs,
		"source",
		"s",
		[]string{},
		"source dumper directories",
	)

	buildCmd.Flags().BoolVarP(&buildInventoryDryRun, "dry-run", "d", false, "dry-run")                        //nolint
	buildCmd.Flags().StringVarP(&buildInventoryRLAHost, "host", "u", "localhost", "RLA service host")         //nolint
	buildCmd.Flags().IntVarP(&buildInventoryRLAPort, "port", "p", defaultServicePort, "RLA service port")     //nolint
	buildCmd.Flags().StringVarP(&buildInventoryOutputDumperDir, "output", "o", "", "output dumper directory") //nolint
}

func doBuild() {
	builder, err := builder.New(
		builder.Config{
			SourceDumperDirs: sourceDumperDirs,
			OutputDumperDir:  buildInventoryOutputDumperDir,
			DryRun:           buildInventoryDryRun,
			RLAClientConf: client.Config{
				Host: buildInventoryRLAHost,
				Port: buildInventoryRLAPort,
			},
		},
	)

	if err != nil {
		log.Fatal().Msgf("failed to build the builder: %v", err)
	}

	defer builder.Done()

	processed := builder.Build(context.Background())
	log.Info().Msgf("built %d racks, %d NVL domains", processed[0], processed[1])
}
