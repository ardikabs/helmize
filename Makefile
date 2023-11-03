PLUGIN_NAME := helmize
Version 	:= $(shell git describe --tags --dirty --always)
GitCommit 	:= $(shell git rev-parse HEAD)
LDFLAGS 	:= "-s -w -X github.com/ardikabs/helmize/cmd.Version=$(Version) -X github.com/ardikabs/helmize/cmd.GitCommit=$(GitCommit)"
OUTDIR 		:= bin

GOLANGCI_VERSION = 1.53.3

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(OUTDIR) v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint "$@"

## audit: format, tidying, vet, lint and test all code
.PHONY: audit
audit: fmt mod lint vet test

## fmt: formatting the code
fmt:
	@echo 'Formatting code...'
	@go fmt $(shell go list ./... | grep -v /vendor/|xargs echo)

## mod: tidying module dependencies
.PHONY: mod
mod:
	@echo 'Tidying and verifying module dependencies...'
	@go mod tidy
	@go mod verify

## lint: linting the code
.PHONY: lint
lint: bin/golangci-lint
	@echo 'Linting code...'
	bin/golangci-lint run

## test: running unit test
.PHONY: test
test:
	@echo 'Running tests...'
	@mkdir -p output
	@go test $(shell go list ./... | grep -v /vendor/|xargs echo) -v -covermode=atomic -cover -coverprofile=./output/coverage.out -race && \
		go tool cover -html=./output/coverage.out -o ./output/coverage.html && \
		go tool cover -func=./output/coverage.out

## vet: vetting the code
.PHONY: vet
vet:
	@echo 'Vetting code...'
	@go vet $(shell go list ./... | grep -v /vendor/|xargs echo)

.PHONY: build
build:
	@mkdir -p $(OUTDIR)
	go build -o $(OUTDIR)/$(PLUGIN_NAME)