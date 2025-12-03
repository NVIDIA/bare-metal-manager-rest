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

FROM golang:1.23-alpine AS builder

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
COPY test ./test
COPY internal ./internal
COPY VERSION ./VERSION

# Build App
RUN make build-test-site-agent

# Multi-Stage final build
FROM alpine

RUN apk add --update --no-cache \
      ca-certificates

WORKDIR /usr/local

COPY --from=builder /go/src/testsiteagent /usr/local/bin/testsiteagent

# Create Temporal certs directory
RUN mkdir -p /var/secrets/temporal/certs

# Set the system worker binary as the entrypoint of the container
CMD ["testsiteagent"]
