.PHONY: fmt lint test coverage license

all: fmt lint coverage

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

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html
