# Copyright 2021 Nvidia
# Elektra Makefile

# Image URL to use all building/pushing image targets
IMG ?= nvcr.io/nvidian/nvpubcloud/elektra:latest
M1IMG ?= nvcr.io/nvidian/nvpubcloud/elektra_mac:vEtcd
HELMIFY = $(shell pwd)/bin/helmify
HELMIFY_VERSION = "v0.0.1"
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")

# Versioning
CLOUD_COMMON_SCRIPTS_PATH?=../cloud-common/scripts
-include $(CLOUD_COMMON_SCRIPTS_PATH)/versioning/versioning.include

.PHONY: env run install-kustomize docker-build-linux docker-build-local docker-build

all: test vet fmt lint build-local run

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	cd ./cmd/elektraserver; go build -race -o ../bin/elektraserver; cd ../..
	./cmd/bin/elektraserver -tout=100 & echo $$! > elektra_server.pid
	if CARBIDE_ADDRESS=localhost:11079 CARBIDE_SEC_OPT=0 go test -v -race ./pkg/...; then echo "Tests passed"; else echo "Tests failed"; fi
	kill `cat elektra_server.pid`; rm elektra_server.pid

test-ci:
	go install github.com/jstemmer/go-junit-report/v2@latest
	go install github.com/boumenot/gocover-cobertura@latest
	cd ./cmd/elektraserver; go build -race -o ../bin/elektraserver; cd ../..
	./cmd/bin/elektraserver -tout=330 & echo $$! > elektra_server.pid
	if CARBIDE_ADDRESS=localhost:11079 CARBIDE_SEC_OPT=0 go test -v -race -p 1 ./pkg/... -coverpkg=./... -coverprofile=coverage.txt 2>&1 | tee report.tmp; then cat report.tmp | go-junit-report > report.xml; go tool cover -func coverage.txt; gocover-cobertura < coverage.txt > coverage.xml; kill `cat elektra_server.pid`; rm elektra_server.pid; else kill `cat elektra_server.pid`; rm elektra_server.pid; exit 1; fi

coverage:
	cd ./cmd/elektraserver; go build -race -o ../bin/elektraserver; cd ../..
	./cmd/bin/elektraserver -tout=40 & echo $$! > elektra_server.pid
	if CARBIDE_ADDRESS=localhost:11079 CARBIDE_SEC_OPT=0 go test -v ./pkg/... -coverpkg=./pkg/... -coverprofile=coverage.txt -covermode=count; then go tool cover -html=coverage.txt -o coverage.html; kill `cat elektra_server.pid`; rm elektra_server.pid; else kill `cat elektra_server.pid`; rm elektra_server.pid; exit 1; fi

lint:
	golangci-lint run --issues-exit-code 0 --out-format code-climate | jq .
	go list ./... | xargs -L1 revive -set_exit_status

build:
	cd ./cmd/elektractl; GOOS=linux GOARCH=amd64 go build -race -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektractl
	cd ./cmd/elektra; GOOS=linux GOARCH=amd64 go build -race -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektra
	cd ./cmd/elektraserver; GOOS=linux GOARCH=amd64 go build -race -ldflags "-extldflags '-static' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektraserver

build-local:
	cd ./cmd/elektractl; go build -ldflags "-X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektractl
	cd ./cmd/elektra; go build -ldflags "-X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektra
	cd ./cmd/elektraserver; go build -ldflags "-X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.Version=$(VERSION)' -X 'gitlab-master.nvidia.com/nvmetal/elektra-site-agent/pkg/metadata.BuildTime=$(BUILD_TIME)'" -o ../bin/elektraserver

clean-local:
	cd ./cmd/elektractl; rm -f ../bin/elektractl
	cd ./cmd/elektra; rm -f ../bin/elektra
	cd ./cmd/elektraserver; rm -f ../bin/elektraserver

clean-test:
	git checkout HEAD config/secret/ca/tls.crt config/secret/client/tls.crt config/secret/client/tls.key

run:
	./cmd/bin/elektra

build-clean:
	rm -f cmd/bin/*

clean:
	rm -f ./pkg/components/elektractl.env ./cmd/bin/elektra ./cmd/bin/elektraserver ./coverage.html
	rm -f ./coverage.txt ./report.xml ./report.tmp ./coverage.xml
	rm -f ./cmd/elektractl/cacert ./cmd/elektractl/creds-url ./cmd/elektractl/otp ./cmd/elektractl/site-uuid
	rm -f elektra_server.pid
	git restore config/secret/ca/tls.crt config/secret/client/tls.crt config/secret/client/tls.key code-quality-report.json

docker-build-linux:  ## Build docker image for linux env
	docker build -f Dockerfile.local --platform linux/amd64 -t ${IMG} --build-arg GIT_CREDENTIALS=${GIT_CREDENTIALS} .
	docker push ${IMG}

docker-build-local:  ## Build docker image for linux env
	docker build -f Dockerfile.local -t ${IMG} --build-arg GIT_CREDENTIALS=${GIT_CREDENTIALS} .
	docker push ${IMG}

docker-build:  ## Build docker image for linux env
	docker build -t ${IMG} .

docker-push:  ## Build docker image for linux env
	docker build -t ${IMG} .

docker-mac:  ## Build docker image for linux env
	docker buildx build --tag ${M1IMG} --platform linux/arm64 . ; docker push ${M1IMG}

kustomize:
	kustomize build k8s/overlays/dev

load-minikube:
	minikube image load ${M1IMG}

update-temporal:
	cd temporal
	helm dependencies update

temporal:
	cd temporal; helm install temporaltest . --timeout 900s

deploy-all:
	kustomize build k8s/overlays/dev | kubectl apply -f -
	kubectl create -f ./k8s/overlay/dev/elektra-data.yaml
	kubectl create -f ./k8s/overlays/dev/pull.yaml
	kubectl create secret -n elektra-site-agent generic temporal-cert --from-file=cacertificate=./config/secret/ca/tls.crt --from-file=certificate=./config/secret/client/tls.crt --from-file=key=./config/secret/client/tls.key
	kubectl create secret -n elektra-site-agent generic bootstrap-info --from-file=cacert=./config/secret/bootstrap/cacert --from-file=otp=./config/secret/bootstrap/otp --from-file=site-uuid=./config/secret/bootstrap/site-uuid --from-file=creds-url=./config/secret/bootstrap/creds-url
	kubectl create secret -n elektra-site-agent tls elektra-tls --key ./config/secret/client/tls.key --cert ./config/secret/client/tls.crt
	kubectl create secret -n elektra-site-agent generic cloud-tls  --from-file=ca.crt=./config/secret/ca/tls.crt
	kubectl rollout restart statefulset elektra -n elektra-site-agent

deploy:
	kustomize build k8s/overlays/dev | kubectl apply -f -

deploy-pgo:
	kubectl apply -k k8s/pgo/namespace
	kubectl apply --server-side -k k8s/pgo/default
	kubectl apply -k k8s/pgo/postgres
ifndef ignore-not-found
  ignore-not-found = true
endif

undeploy:
	kubectl delete --ignore-not-found=$(ignore-not-found) -k k8s/pgo/default
	kubectl delete --ignore-not-found=$(ignore-not-found) -k k8s/pgo/namespace
	kubectl delete --ignore-not-found=$(ignore-not-found) -k k8s/pgo/postgres
	kustomize build k8s/overlays/dev | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

undeploy-pgo:
	kubectl delete -n elektra-site-agent --ignore-not-found=$(ignore-not-found) -k k8s/pgo/default
	kubectl delete -n elektra-site-agent --ignore-not-found=$(ignore-not-found) -k k8s/pgo/namespace
	kubectl delete -n elektra-site-agent --ignore-not-found=$(ignore-not-found) -k k8s/pgo/postgres


## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize

## Tool Versions
KUSTOMIZE_VERSION ?= v4.5.7

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"

install-kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

helmify:
	GOBIN=$(PROJECT_DIR)/bin go install github.com/vinodchitraliNVIDIA/helmify/cmd/helmify@${HELMIFY_VERSION}

helm: install-kustomize helmify
	$(KUSTOMIZE) build k8s/overlays/helm | $(HELMIFY) elektra-site-agent

elektractl:
	cd cmd/elektractl; GOOS=linux GOARCH=amd64 go build -o elektractl ; chmod 777 elektractl


tls-secret:
	kubectl create secret -n elektra-site-agent generic temporal-cert --from-file=cacertificate=./config/secret/bootstrap/cacert --from-file=certificate=./config/secret/client/tls.crt --from-file=key=./config/secret/client/tls.key
	kubectl create secret -n elektra-site-agent generic bootstrap-info --from-file=cacert=./config/secret/bootstrap/cacert --from-file=otp=./config/secret/bootstrap/otp --from-file=site-uuid=./config/secret/bootstrap/site-uuid --from-file=creds-url=./config/secret/bootstrap/creds-url
	kubectl create secret -n elektra-site-agent tls elektra-tls --key ./config/secret/client/tls.key --cert ./config/secret/client/tls.crt
	kubectl create secret -n elektra-site-agent generic cloud-tls  --from-file=ca.crt=./config/secret/ca/tls.crt

update-tls:
	kubectl delete secrets elektra-tls cloud-tls -n elektra-site-agent
	kubectl create secret -n elektra-site-agent tls elektra-tls --key ./config/secret/client/tls.key --cert ./config/secret/client/tls.crt
	kubectl create secret -n elektra-site-agent generic cloud-tls  --from-file=ca.crt=./config/secret/ca/tls.crt

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

update-schema:
	GOPROXY=direct GOSUMDB=off go get "gitlab-master.nvidia.com/nvmetal/cloud-workflow-schema/schema/site-agent/workflows/v1"
