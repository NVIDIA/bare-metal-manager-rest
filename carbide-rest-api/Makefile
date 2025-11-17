CLOUD_COMMON_SCRIPTS_PATH?=../cloud-common/scripts
include $(CLOUD_COMMON_SCRIPTS_PATH)/versioning/versioning.include

.PHONY: fmt lint test coverage build run swagger license clean

all: fmt lint coverage build run

fmt:
	go fmt ./...

lint:
	go vet ./...
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	go list ./... | xargs -L1 revive -config .revive.toml -set_exit_status

test:
	go test -race -p 1 ./... -count=1

coverage:
	mkdir -p ci-reports
	go test -race -p 1 ./... -coverpkg=./... -coverprofile ./ci-reports/coverage.txt
	go tool cover -func ./ci-reports/coverage.txt

build-local:
	go build -ldflags "-X 'gitlab-master.nvidia.com/nvmetal/cloud-api/pkg/metadata.Version=$(cat VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/cloud-api/pkg/metadata.BuildTime=$(date +"%Y-%m-%d %H:%M:%S")'" ./cmd/api

build:
	go build -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/cloud-api/pkg/metadata.Version=$(cat VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/cloud-api/pkg/metadata.BuildTime=$(date +"%Y-%m-%d %H:%M:%S")'" ./cmd/api

run:
	./api

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

clean:
	rm api
