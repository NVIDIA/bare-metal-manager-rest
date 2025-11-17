.PHONY: fmt lint test coverage build run license clean

all: fmt lint test build run

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

build:
	go build -race ./cmd/migrations

run:
	./migrations

license:
	license_finder report --decisions_file .license-config.yml --format html --save=./ci-reports/license-report.html

clean:
	rm -f migrations
