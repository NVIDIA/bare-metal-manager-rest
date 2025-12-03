# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
# property and proprietary rights in and to this material, related
# documentation and any modifications thereto. Any use, reproduction,
# disclosure or distribution of this material and related documentation
# without an express license agreement from NVIDIA CORPORATION or
# its affiliates is strictly prohibited.
#

FROM golang:1.24-alpine AS builder

ENV GIT_TERMINAL_PROMPT=1
RUN go env -w GOPRIVATE="github.com/NVIDIA"

RUN apk add --no-cache \
      bash \
      git \
      jq \
      gcc \
      musl-dev \
      gcompat \
      make

WORKDIR /go/src

# Fetch and cache dependencies
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=secret,id=gitcreds,dst=/root/.netrc go mod download

# Add source code
COPY Makefile .
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal
COPY VERSION ./VERSION

# Create Temporal certs directory
RUN mkdir -p /var/secrets/temporal/certs

# Build App
RUN CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" ./cmd/workflow

# Multi-Stage production build
FROM nvcr.io/nvidia/distroless/go:v3.2.1 AS final

# Copy the binary from the builder stage
COPY --from=builder /go/src/workflow /app/workflow

# Copy Temporal certs directory
COPY --from=builder /var/secrets/temporal/certs /var/secrets/temporal/certs

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/app/workflow"]
