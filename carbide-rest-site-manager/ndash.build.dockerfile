FROM golang:1.23

# Install base libraries
RUN apt-get update && \
    apt-get install jq -y && \
    apt install ruby-full -y

# Install dependencies
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v1.55.2
RUN go install github.com/mgechev/revive@v1.6.1
RUN go install github.com/boumenot/gocover-cobertura@latest
RUN go install github.com/jstemmer/go-junit-report@latest

# Install Ruby dependencies
RUN gem install license_finder
