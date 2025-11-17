
MAIN_BRANCH=master

.ONESHELL:
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

CGO_ENABLED := $(or ${CGO_ENABLED},0)
GO := go
GO111MODULE := on
LINKMODE := -extldflags '-static -s -w'


.EXPORT_ALL_VARIABLES:

.PHONY: bench benchstat lint test postgres-up postgres-rm cockroach-up cockroach-up-cluster cockroach-rm license
all: proto server client test fuzz bench
ci: proto test

bench:
	CGO_ENABLED=1 $(GO) test -bench ./... -run=- -benchmem -timeout 20m

benchstat:
	git stash
	CGO_ENABLED=1 $(GO) test -bench . -run=- -count 5 -benchmem > old.txt
	git stash pop
	CGO_ENABLED=1 $(GO) test -bench . -run=- -count 5 -benchmem > new.txt
	benchstat old.txt new.txt

lint:
	CGO_ENABLED=1 $(GO) vet ./...
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	CGO_ENABLED=1 $(GO) list ./... | xargs -L1 revive -config .revive.toml -set_exit_status

test:
	CGO_ENABLED=1 $(GO) test ./...

coverage:
	mkdir ci-reports
	CGO_ENABLED=1 $(GO) test ./... -coverpkg=./... -coverprofile ./ci-reports/coverage.txt
	CGO_ENABLED=1 $(GO) tool cover -func ./ci-reports/coverage.txt

.PHONY: fuzz
fuzz:
	CGO_ENABLED=1 $(GO) test -fuzz=Fuzz -fuzztime 30s -v -run ^Fuzz

.PHONY: golangcicheck
golangcicheck:
	@bash -c "type -P golangci-lint;" 2>/dev/null || (echo "golangci-lint is required but not available in current PATH. Install: https://github.com/golangci/golangci-lint#install"; exit 1)

.PHONY: server
server:
	go build -tags netgo,osusergo,urfave_cli_no_docs \
		 -ldflags "$(LINKMODE) -X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
	   -o bin/server ./cmd/server
	# strip bin/server

.PHONY: client
client:
	go build -tags netgo,osusergo,urfave_cli_no_docs \
		 -ldflags "$(LINKMODE) -X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
	   -o bin/cli ./cmd/client
	# strip bin/cli

.PHONY: proto
proto:
	$(MAKE) -C proto protoc

.PHONY: postgres-up
postgres-up: postgres-rm
	docker run -d --name ipamdb -p 5433:5432 -e POSTGRES_PASSWORD="password" postgres:$(PG_VERSION) -c 'max_connections=200'

postgres-rm:
	docker rm -f ipamdb || true

cockroach-up: cockroach-rm postgres-rm
	# https://www.cockroachlabs.com/docs/v19.2/start-a-local-cluster-in-docker-linux.html#main-content
	docker network create -d bridge roachnet
	docker run -d --name=roach1 --hostname=roach1 --net=roachnet -p 5433:26257 -p 8080:8080 cockroachdb/cockroach:$(COCKROACH_VERSION) start-single-node --insecure --listen-addr=0.0.0.0

cockroach-up-cluster: cockroach-rm
	# https://www.cockroachlabs.com/docs/v19.2/start-a-local-cluster-in-docker-linux.html#main-content
	docker network create -d bridge roachnet
	docker run -d --name=roach1 --hostname=roach1 --net=roachnet -p 5433:26257 -p 8080:8080 cockroachdb/cockroach:$(COCKROACH_VERSION) start --insecure --join=roach1,roach2,roach3
	docker run -d --name=roach2 --hostname=roach2 --net=roachnet cockroachdb/cockroach:$(COCKROACH_VERSION) start --insecure --join=roach1,roach2,roach3
	docker run -d --name=roach3 --hostname=roach3 --net=roachnet cockroachdb/cockroach:$(COCKROACH_VERSION) start --insecure --join=roach1,roach2,roach3
	docker exec -it roach1 ./cockroach init --insecure

cockroach-rm:
	docker rm -f roach1 roach2 roach3 || true
	docker network rm roachnet || true

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

.PHONY: redis-up
redis-up: redis-rm
	docker run -d --name ipamredis -p 6379:6379 redis

.PHONY: redis-rm
redis-rm:
	docker rm -f ipamredis || true
