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

FROM golang:1.23 AS builder

ENV GIT_TERMINAL_PROMPT=1
RUN go env -w GOPRIVATE="github.com/NVIDIA"

WORKDIR /go/src

# Fetch and cache dependencies
COPY ./go.mod ./
COPY ./go.sum ./

RUN --mount=type=secret,id=gitcreds,dst=/root/.netrc go mod download

# Add source code
COPY Makefile .
COPY pkg ./pkg
COPY cmd ./cmd
COPY VERSION .

# Build the elektra binary
RUN export VERSION=$(cat ./VERSION)
RUN export BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S")
RUN mkdir -p /usr/elektra
RUN CGO_ENABLED=0 go build -ldflags "-extldflags '-static' -X 'github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/metadata.Version=$(VERSION)' -X 'github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" ./cmd/elektra

# Multi-Stage production build
FROM nvcr.io/nvidia/distroless/go:v3.2.1 AS final

# Retrieve the binary from the previous stage
COPY --from=builder /go/src/elektra /usr/local/bin/elektra

# Set the elektra worker binary as the entrypoint of the container
CMD ["/usr/local/bin/elektra"]
