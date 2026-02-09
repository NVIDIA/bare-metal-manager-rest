#!/bin/bash
#

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Starting OpenAPI Redoc Server at http://127.0.0.1:8090";

docker run -it --rm -p 8090:80 -v $SCRIPT_DIR/spec.yaml:/usr/share/nginx/html/openapi.yaml -e SPEC_URL=openapi.yaml redocly/redoc
