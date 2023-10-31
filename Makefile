Version := $(shell git describe --tags --dirty --always)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/ardikabs/kasque/cmd.Version=$(Version) -X github.com/ardikabs/kasque/cmd.GitCommit=$(GitCommit)"
OUTDIR := bin

GOLANGCI_VERSION = 1.31.0

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint "$@"

## audit: tidy and vendor dependencies and format, vet, lint and test all code
.PHONY: audit
audit: fmt vendor lint vet test

## fmt: formatting the code
fmt:
	@echo 'Formatting code...'
	@go fmt $(shell go list ./... | grep -v /vendor/|xargs echo)

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	@go mod tidy
	@go mod verify
	@echo 'Vendoring dependencies...'
	@go mod vendor

## lint: linting the code
.PHONY: lint
lint: bin/golangci-lint
	@echo 'Linting code...'
	bin/golangci-lint run

## test: run unit test
.PHONY: test
test:
	@echo 'Running tests...'
	@mkdir -p output
	@go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover -coverprofile=./output/coverage.out -race && \
		go tool cover -html=./output/coverage.out -o ./output/coverage.html && \
		go tool cover -func=./output/coverage.out

## vet: run vetting the code
.PHONY: vet
vet:
	@echo 'Vetting code...'
	@go vet $(shell go list ./... | grep -v /vendor/|xargs echo)
