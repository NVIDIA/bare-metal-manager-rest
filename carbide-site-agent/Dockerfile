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
RUN go env -w GOPRIVATE="gitlab-master.nvidia.com"

WORKDIR /go/src

# Fetch and cache dependencies
COPY ./go.mod ./
COPY ./go.sum ./
RUN --mount=type=secret,id=gitcreds,dst=/root/.netrc go mod download

# Add source code
COPY Makefile .
# Copy the framework
COPY pkg ./pkg
# Copy the executable
COPY cmd ./cmd
COPY VERSION .

#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on make build
RUN export VERSION=$(cat ./VERSION)
RUN export BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S")
RUN mkdir -p /usr/elektra
RUN cd cmd/elektra && GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o elektra && mv elektra /usr/elektra/elektra
RUN cd cmd/elektractl && GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o elektractl && mv elektractl /usr/elektra/elektractl


# Multi-Stage final build
FROM alpine

#RUN apk update && apk add bash && apk add curl

RUN apk add --update --no-cache \
      ca-certificates \
      bash \
      curl \
      gcompat
#RUN apk add sed attr dialog dialog-doc bash bash-doc bash-completion grep grep-doc
#RUN apk add util-linux util-linux-doc pciutils usbutils binutils findutils readline &&


# Retrieve the binary from the previous stage
COPY --from=builder /usr/elektra/elektra /usr/local/bin/elektra
COPY --from=builder /usr/elektra/elektractl /usr/local/bin/elektractl

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY cmd/elektractl/elektractl.env /usr/local/bin/elektractl.env


# Install kubectl
RUN curl -LO https://dl.k8s.io/release/v1.23.0/bin/linux/amd64/kubectl
RUN chmod u+x kubectl && mv kubectl /bin/kubectl && chmod +x /usr/local/bin/elektra && chmod +x /usr/local/bin/elektractl

RUN chmod +x /usr/local/bin/entrypoint.sh

# Start the entrypoint script
# Set the elektra worker binary as the entrypoint of the container
CMD ["/usr/local/bin/entrypoint.sh", "elektra"]
