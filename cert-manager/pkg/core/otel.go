// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package core

import (
	"context"

	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// tracer is the OTel tracer to use for the core package
var tracer oteltrace.Tracer

func init() {
	tracer = otel.Tracer("nvmetal/cloud-cert-manager/pkg/core")
}

// StartOTELDaemon starts a go routine that waits on the provided context to quit and then shuts down the daemon
func StartOTELDaemon(ctx context.Context) {
	log := GetLogger(ctx)

	// Ignore this is most likely disabled
	log.Infof("Skipping OTEL startup - not supported")
}
