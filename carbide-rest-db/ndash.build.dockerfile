FROM golang:1.24

# Install base libraries
RUN apt-get update && \
    apt-get install jq -y && \
    apt install ruby-full -y

# Install dependencies
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2 && \
    go install github.com/mgechev/revive@v1.4.0 && \
    go install github.com/boumenot/gocover-cobertura@latest && \
    go install github.com/jstemmer/go-junit-report@latest

# Install Ruby dependencies
RUN gem install license_finder
