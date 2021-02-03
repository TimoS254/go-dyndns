PROJECT_NAME=go-dyndns

GIT_VERSION=$(shell git describe --tags --always)
GIT_BRANCH=$(shell git branch --show-current)
GIT_DEFAULT_BRANCH=main

LD_FLAGS = -X main.GitVersion=${GIT_VERSION} -X main.GitBranch=${GIT_BRANCH}

OUTPUT_PREFIX=./bin/${PROJECT_NAME}-${GIT_VERSION}

OUTPUT_SUFFIX=$(shell go env GOEXE)

# tests project with the built-in Golang tool
.PHONY: build
test:
	@go test -timeout 1m ./...

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-${GOOS}-${GOARCH}${OUTPUT_SUFFIX}" ./cmd/${PROJECT_NAME}/*

# build go application for docker usage
.PHONY: build-docker
build-docker:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-docker" ./cmd/${PROJECT_NAME}/*

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/${PROJECT_NAME}/*