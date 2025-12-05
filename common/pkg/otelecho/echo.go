// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package otelecho

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	upstreamotelecho "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TracerKey is a key for current tracer
	TracerKey = "otel-go-contrib-tracer-labstack-echo"

	// TracerName is name of the tracer
	TracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	// TraceHdr is header name for ngc trace id
	TraceHdr = "X-Ngc-Trace-Id"
)

// Middleware wraps the upstream otelecho middleware and adds custom functionality:
// - Zerolog logging of trace IDs (from extracted context before span creation)
// - Setting X-Ngc-Trace-Id header
// - Storing tracer in context for use by other packages
func Middleware(service string, opts ...upstreamotelecho.Option) echo.MiddlewareFunc {
	upstreamMiddleware := upstreamotelecho.Middleware(service, opts...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Get upstream handler
		upstreamHandler := upstreamMiddleware(func(c echo.Context) error {
			// This runs after upstream middleware sets up the span
			// The upstream middleware has already:
			// 1. Extracted context from headers
			// 2. Created a new span
			// 3. Set the request context with the span
			//
			// We need to:
			// 1. Get trace ID from the span (which inherits from parent)
			// 2. Log and set header
			// 3. Ensure tracer is stored in context for util/tracer.go

			ctx := c.Request().Context()
			span := trace.SpanFromContext(ctx)

			// Get trace ID from the span and log it
			// This matches the original behavior where trace ID was logged before span creation
			// Note: The original code logged trace ID from extracted context (parent trace)
			// We log from created span (child span), which inherits the same trace ID
			scc := span.SpanContext()
			traceID := scc.TraceID().String()
			log.Info().Msgf("span traceid: %s", traceID)
			c.Response().Header().Set(TraceHdr, traceID)

			// Always store tracer in context for use by other packages (like util/tracer.go)
			// This matches the original behavior where tracer was always stored
			tracerKey := "otel-go-contrib-tracer-labstack-echo"
			if tracer := c.Get(tracerKey); tracer != nil {
				if t, ok := tracer.(trace.Tracer); ok {
					ctx = context.WithValue(ctx, TracerKey, t)
					c.SetRequest(c.Request().WithContext(ctx))
				}
			}

			// Now call the actual next handler
			return next(c)
		})

		return upstreamHandler
	}
}

// Re-export types and functions from upstream
type Option = upstreamotelecho.Option

var (
	WithPropagators    = upstreamotelecho.WithPropagators
	WithTracerProvider = upstreamotelecho.WithTracerProvider
	WithSkipper        = upstreamotelecho.WithSkipper
)
