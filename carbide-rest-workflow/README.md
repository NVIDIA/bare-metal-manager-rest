[![pipeline status](https://gitlab-master.nvidia.com/nvmetal/cloud-workflow/badges/main/pipeline.svg)](https://gitlab-master.nvidia.com/nvmetal/cloud-workflow/-/commits/main)  [![coverage report](https://gitlab-master.nvidia.com/nvmetal/cloud-workflow/badges/main/coverage.svg)](https://gitlab-master.nvidia.com/nvmetal/cloud-workflow/-/commits/main)

## Forge Cloud Workflow

Cloud workflow provides Temporal workflow definitions for asynchronous workers

## Setup Postgres/Temporal

Please follow instructions at [Forge Cloud Local](https://gitlab-master.nvidia.com/nvmetal/cloud-local) repo to get Postgres and Temporal set up for your local env.

Copy the following certs over from Cloud Local directory:

    export FORGE_CLOUD_LOCAL_DIR=<path-to-cloud-local-repo>
    mkdir certs
    cp -R $FORGE_CLOUD_LOCAL_DIR/temporal/tls/certs/client-cloud certs/
    cp -R $FORGE_CLOUD_LOCAL_DIR/temporal/tls/certs/server-intermediate-ca certs/

Create DB for Forge:

    PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -p 30432 < $FORGE_CLOUD_LOCAL_DIR/postgres/scripts/setup.sql

## Run Cloud Worker

Get all modules:

    go mod download

Build the workflow executable:

    make build

Run the API server:

    make run

## Run Site Worker

To run the Site Worker instead, use the following command:

    TEMPORAL_NAMESPACE=site TEMPORAL_QUEUE=site make run

## Run Dummy Site Agent Worker

Get all modules:

    go mod download

Build the workflow executable:

    make build-test-site-agent

Run the API server:

    make run-test-site-agent

To run the Dummy Site Agent Worker instead, use the following command:

    TEMPORAL_NAMESPACE=<<site_id>> TEMPORAL_QUEUE=<<site_id>> make site-agent-run

To run the Dummy Site Agent Worker as k8s local deplopyment.
    
    to use dummy site agent for local temporal call testing, following steps can be done
    
    after deploying local env using cloud-local, run the following steps (note: from cloud-local, cloud-workflow (cloud/site worker) has to be up and running on local dev)
    
    dummy site agent images can be built using the following commands.

    docker build -f test/siteagent/Dockerfile --secret id=gitcreds,src=netrc.txt . -t localhost:5001/test-site-agent:local-dev
    docker push localhost:5001/test-site-agent:local-dev

    k8s apply command.

    kubectl apply -f deploy/kustomize/overlays/local/test-site-agent-deployment.yaml
