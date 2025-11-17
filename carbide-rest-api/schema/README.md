# Forge Cloud API Schema

This repo contains OpenAPI schema for Forge Cloud API endpoints.

To view a rendered/browsable version of the schema, make sure you have docker and run the following:

    ./serve.sh

If `serve.sh` gives you error due to volume mount permission issues, then checkout the repo and run the following from project root:

    docker build -f schema/Dockerfile . -t forge-api-docs:latest
    docker run --rm -it -p 8090:80 forge-api-docs:latest

Then access the schema at:

    http://127.0.0.1:8090
