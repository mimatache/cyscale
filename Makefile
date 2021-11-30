SHELL:=/bin/bash
APP:=cyscale-cli
BUILD_DIR:=build
BIN_DIR:=$(BUILD_DIR)/$(APP)/_bin

DOCKER_REPO?="matache91mh"
IMAGE?=$(DOCKER_REPO)/$(APP)

VERSION ?= $(shell git describe --tags --dirty --always)
BUILD_DATE ?= $(shell date +%FT%T%z)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)

LDFLAGS += -X 'github.com/mimatache/cyscale/internal/info.appName=${APP}'
LDFLAGS += -X 'github.com/mimatache/cyscale/internal/info.version=${VERSION}'
LDFLAGS += -X 'github.com/mimatache/cyscale/internal/info.commitHash=${COMMIT_HASH}'
LDFLAGS += -X 'github.com/mimatache/cyscale/internal/info.buildDate=${BUILD_DATE}'


all: install-go-tools generate fmt lint test build
	
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -v -o $(BIN_DIR)/$(APP) .

test-ci:
	go test -v  -race -json -coverprofile=coverage.out ./... > unit-test.json
	go tool cover -func=coverage.out

test:
	go test -v -race -cover ./...

install-go-tools:
	GO111MODULE=on CGO_ENABLED=0 go get github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	go vet ./...
	golangci-lint run

fmt:
	go mod tidy
	goimports -w .
	gofmt -s -w .

run:
	GO111MODULE=on CGO_ENABLED=0 go run .

docker-build:
	docker build -t $(IMAGE):$(VERSION) --build-arg VERSION=${VERSION} --build-arg BUILD_DATE=${BUILD_DATE} --build-arg COMMIT_HASH=${COMMIT_HASH} --build-arg APP=${APP} .

docker-push: app-image
	docker push $(IMAGE):$(VERSION)

