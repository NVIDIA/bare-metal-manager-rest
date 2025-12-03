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
RUN go env -w GOPRIVATE="gitlab-master.nvidia.com"

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

# Add source code
COPY Makefile .
COPY cmd ./cmd
COPY pkg ./pkg

# Copy files not related to compilation
COPY VERSION ./VERSION

# Download dependencies
RUN --mount=type=secret,id=gitcreds,dst=/root/.netrc go mod download

# Build the binary
RUN CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" ./cmd/migrations

# Multi-Stage production build
FROM nvcr.io/nvidia/distroless/go:v3.2.1 AS final

# Copy the binary from the builder stage
COPY --from=builder /go/src/migrations /usr/local/bin/migrations

# Set the migrations command as the entrypoint of the container
CMD ["/usr/local/bin/migrations", "db", "init_migrate"]
