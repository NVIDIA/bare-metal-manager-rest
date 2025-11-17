# Copyright 2024 Nvidia
# Forge Cloud Workflow Schema

CLOUD_COMMON_SCRIPTS_PATH?=../cloud-common/scripts
-include $(CLOUD_COMMON_SCRIPTS_PATH)/versioning/versioning.include

.PHONY: env run

all: test vet fmt lint build-local run

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -race ./...

test-ci:
	go test -race ./... 2>&1 | go-junit-report > report.xml

coverage:
	go test $(go list ./... | grep -v /docs) -covermode=count -coverprofile coverage.txt
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	go list ./... | xargs -L1 revive -set_exit_status

clean:
	rm -rf workflows/proto/*.proto
clean-server:
	rm -rf schema/site-agent/workflows/v1

protogen:
	# disabling lint... forge proto definitions are not passing the lint due to too many non-standard definitions
	# echo "Checking validity of proto files"
	# buf lint
	echo "Generating go proto files now"
	buf generate

