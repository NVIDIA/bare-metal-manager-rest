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

RUN --mount=type=secret,id=gitcreds,dst=/root/.netrc go mod download
# Build App
RUN make build

# Copy files not related to compilation
COPY scripts/migrate.sh ./migrate.sh
COPY VERSION ./VERSION

# Multi-Stage final build
FROM alpine

RUN apk add --update --no-cache \
      ca-certificates

WORKDIR /usr/local

COPY --from=builder /go/src/migrations /usr/local/bin/migrations
COPY --from=builder /go/src/migrate.sh /usr/local/bin/migrate.sh

# Set the system worker binary as the entrypoint of the container
ENTRYPOINT ["/usr/local/bin/migrate.sh"]
