.PHONY: fmt lint test coverage build build-local run license clean

all: fmt lint test build run

fmt:
	go fmt ./...

lint:
	go vet ./...
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	go list ./... | xargs -L1 revive -config .revive.toml -set_exit_status

test:
	go test -race ./...

coverage:
	mkdir -p ci-reports
	go test -race ./... -coverpkg=./... -coverprofile ./ci-reports/coverage.txt
	go tool cover -func ./ci-reports/coverage.txt

build:
	CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" ./cmd/credsmgr

build-local:
	go build -ldflags "-extldflags '-static'" ./cmd/credsmgr

run:
	./credsmgr

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

generate-cert-manager:
	@echo "Generating cert-manager manifests from Helm chart..."
	@./scripts/generate-cert-manager-manifests.sh

clean:
	rm -rf ./credsmgr
