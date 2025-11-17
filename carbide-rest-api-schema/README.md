# Cloud Workflow Schema

This repository defines the Temporal workflow interfaces between Forge Cloud and Forge Site Agent.

It also hosts the copy of the Carbide gRPC API proto that must be used by all Cloud services.

## Overview

Carbide gRPC API proto file is periodically updated to pull in latest changes from Carbide repo. Since proto is generally backwards compatible, this update is usually done when a new feature is being implemented in Cloud and the corresponding proto attributes are missing in the current proto file in Cloud Workflow Schema.

Not all objects/methods changed in proto is actively used by Cloud. But we pull them in to stay in sync.

## Using Workflow Schema

Incorporate the Go module using:

    go get github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema

The Cloud Workflow Schema has semantic versioning e.g. v0.0.33. Specific versions should be pulled in as needed.
