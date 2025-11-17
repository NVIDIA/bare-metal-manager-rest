[![pipeline status](https://gitlab-master.nvidia.com/nvmetal/cloud-api/badges/main/pipeline.svg)](https://gitlab-master.nvidia.com/nvmetal/cloud-api/-/commits/main) [![coverage report](https://gitlab-master.nvidia.com/nvmetal/cloud-api/badges/main/coverage.svg)](https://gitlab-master.nvidia.com/nvmetal/cloud-api/-/commits/main)
# Forge Cloud API

F orge Cloud API is the RESTful API interface for Forge

## Setup Postgres/Temporal

Please follow instructions at [Forge Cloud Local](https://gitlab-master.nvidia.com/nvmetal/cloud-local) repo to get Postgres and Temporal set up for your local env.

## Run API Server

Get all modules:

    go mod download

Build the API executable:

    make build

If you wan to run the API server directly from your shell, then you have to copy over the Temporal certs:

    mkdir -p certs/client-cloud
    kubectl get secrets -n cloud-api temporal-client-cloud-certs -o json | jq -r '.data."ca.crt"' | base64 -d > certs/client-cloud/ca.crt
    kubectl get secrets -n cloud-api temporal-client-cloud-certs -o json | jq -r '.data."tls.crt"' | base64 -d > certs/client-cloud/tls.crt
    kubectl get secrets -n cloud-api temporal-client-cloud-certs -o json | jq -r '.data."tls.key"' | base64 -d > certs/client-cloud/tls.key

Run the API server:

    make run

This will start the API server at: [localhost:8388](http://localhost:8388). It will also expose the metrics server at [localhost:9360](http://localhost:9360)
