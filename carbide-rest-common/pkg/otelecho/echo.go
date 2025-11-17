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
//
// Note: The upstream middleware calls c.Error(err) AND returns err, causing double
// error handling. Our wrapper fixes this by catching the error after upstream processing
// and returning nil if the error was already handled.
func Middleware(service string, opts ...upstreamotelecho.Option) echo.MiddlewareFunc {
	upstreamMiddleware := upstreamotelecho.Middleware(service, opts...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Wrap next with our custom logic
		customNext := func(c echo.Context) error {
			// Add our custom logic before calling the handler
			ctx := c.Request().Context()
			span := trace.SpanFromContext(ctx)
			scc := span.SpanContext()

			// Get trace ID from the span and set header (even for noop/zero spans)
			traceID := scc.TraceID().String()
			c.Response().Header().Set(TraceHdr, traceID)
			
			// Only log if we have a non-zero trace ID
			if traceID != "00000000000000000000000000000000" {
				log.Info().Msgf("span traceid: %s", traceID)
			}

			// Always store tracer in context for use by other packages (like util/tracer.go)
			tracerKey := "otel-go-contrib-tracer-labstack-echo"
			if tracer := c.Get(tracerKey); tracer != nil {
				if t, ok := tracer.(trace.Tracer); ok {
					ctx = context.WithValue(ctx, TracerKey, t)
					c.SetRequest(c.Request().WithContext(ctx))
				}
			}

			// Call the actual handler
			return next(c)
		}

		// Apply upstream middleware to our wrapped handler
		instrumentedHandler := upstreamMiddleware(customNext)

		// Return a handler that prevents double error handling
		return func(c echo.Context) error {
			// The upstream middleware calls c.Error(err) internally when there's an error,
			// which invokes the HTTPErrorHandler. Then it also returns the error, causing
			// Echo to invoke the HTTPErrorHandler a second time.
			// To fix this, we call the instrumented handler and return nil to prevent
			// the second invocation, since the error was already handled.
			err := instrumentedHandler(c)
			if err != nil {
				// Error was already handled by upstream middleware's c.Error(err) call
				// Return nil to prevent Echo from handling it again
				return nil
			}
			return nil
		}
	}
}

// Re-export types and functions from upstream
type Option = upstreamotelecho.Option

var (
	WithPropagators    = upstreamotelecho.WithPropagators
	WithTracerProvider = upstreamotelecho.WithTracerProvider
	WithSkipper        = upstreamotelecho.WithSkipper
)
