#!/bin/bash
#

docker run -d --rm --name project-test -p 30432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=project postgres:14.4-alpine
