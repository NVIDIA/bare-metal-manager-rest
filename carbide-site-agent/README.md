[![pipeline status](https://gitlab-master.nvidia.com/nvmetal/elektra-site-agent/badges/main/pipeline.svg)](https://gitlab-master.nvidia.com/nvmetal/elektra-site-agent/-/commits/main)
[![coverage report](https://gitlab-master.nvidia.com/nvmetal/elektra-site-agent/badges/main/coverage.svg)](https://gitlab-master.nvidia.com/nvmetal/elektra-site-agent/-/commits/main)

# Elektra Site Agent

## Contents
- [Contents](#contents)
- [Intro](#intro)
- [Installation](#installation)
    - [Development Prerequisites](#development-prerequisites)
    - [Quick Start](#quick-start)
- [Components](#components)
- [Testing](#testing)
    - [Local Build Testing](#local-build-testing)

# Intro

Elektra Site Agent is the site manager who manages Forge Data center sites and provisions DPU by interfacing with Site Controller

Elektra Site Agent Components
* Site Agent
* Carbide gRPC Library
* Site Agent Data Store

Elektra Site Agent interfaces
- Northbound - Temporal interface with Forge Cloud
- Southbound - Site Controller

We need the Elektra Site Agent Framework to be:
* Extensible
* Highly-Available
* Scalable

Elektra Site Agent manages Site workflows and store Site workflow information in DB

# Installation

### Development Prerequisites

* [docker-desktop](https://docs.docker.com/desktop/mac/install/)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) (or)
* [minikube](https://minikube.sigs.k8s.io/docs/start/)
* stern `brew install stern`
* make `brew install make`

### Quick start

#### Build Elektra Site Agent
`make build`


#### Deploy Elektra to a k8s cluster

Make sure have username and token configured for nvce registry access in env variable

`export GIT_CREDENTIALS={{username}}:{{token}}`

Build and push the docker image:

`make docker-build && docker-push`

K8s deployment is maintained via kustomize

Create mTLS certs for communicating securely with Temporal (this will be added later)

Deploy elektra site agent on k8s:

`make deploy`

This will start Site Agent

# Components

* Elektra Site agent
* NoSql Database
* Local Temporal

# Testing

You'll want to have three shells for this. One for temporal, one for elektra logs, and one to execute commands.

### Start Temporal

The following steps will run a local instance of the Temporal Server using the default configuration file

- `git clone https://github.com/temporalio/docker-compose.git`
- `cd  docker-compose`
- `docker-compose up`
- `http://localhost:8088`

### Start Local k8s cluster

- `kind create cluster`
- `kubectl cluster-info --context kind-kind`

## Local Build Testing

- `make  build`
- `make  run`

### Build the docker images &  push

- `make docker-build`
- `make docker-push`

## Deploy Elektra to local cluster

- `minikube load image <desired_image>`
- `or if using kind; load image onto local cluster`
- `make deploy`

## Teardown

- `make undeploy`
