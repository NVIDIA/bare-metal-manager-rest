/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package main

import (
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

var Elog zerolog.Logger

const (
	elektractlLogFile = "elektractl.log"
)

func InitLog() {
	file, err := os.OpenFile(
		"elektractlLogFile",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatalln("Elektractl log cannot be created", err)
	}

	defer file.Close()

	buildInfo, _ := debug.ReadBuildInfo()

	// Detailed Logger can be finetuned anytime later
	Elog = zerolog.New(zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	Elog.Info().Msg("Initializing Elektractl cli log")
}
