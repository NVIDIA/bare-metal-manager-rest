export IMAGE_TAG ?= test
export IMAGE_REPO ?= "nvcr.io/nvidian"
export IMAGE_NAME ?= "cloud-site-manager"

.PHONY: fmt lint test coverage build build-local run image image-push chart license clean generate-cert-manager

all: test fmt lint build run

fmt:
	go fmt ./...

lint:
	go vet $$(go list ./... | grep -v '/pkg/client/')
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	go list ./... | grep -v '/pkg/client/' | xargs -L1 revive -config .revive.toml -set_exit_status

test:
	go test -race ./...

coverage:
	mkdir -p ci-reports
	go test -race ./... -coverpkg=./... -coverprofile ./ci-reports/coverage.txt
	go tool cover -func ./ci-reports/coverage.txt

codegen-update:
	docker run --rm -v $(PWD):/go/src/gitlab-master.nvidia.com/nvmetal/cloud-site-manager golang:1.23-alpine sh -c "cd /go/src/gitlab-master.nvidia.com/nvmetal/cloud-site-manager/scripts && ./codegen_update"

build:
	go build -ldflags "-extldflags '-static'" ./cmd/sitemgr

build-local:
	go build ./cmd/sitemgr

run:
	./sitemgr

image:
	docker build --secret id=gitcreds,src=netrc.txt . -t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG} -f Dockerfile

image-push:
	docker push ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_TAG}

chart:
	helm template --output-dir=./test/manifests --values ./test/values.yaml charts/csm

generate-cert-manager:
	@echo "Generating cert-manager manifests..."
	@./scripts/generate-cert-manager-manifests.sh

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

clean:
	rm -rf ./sitemgr
