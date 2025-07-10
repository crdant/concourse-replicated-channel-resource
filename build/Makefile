.PHONY: build test clean docker-build docker-push

BINARY_NAME=concourse-replicated-channel-resource
VERSION?=0.1.0
DOCKER_REPO=ghcr.io/replicatedhq/concourse-replicated-channel-resource

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o $(BINARY_NAME) ./cmd/concourse-replicated-channel-resource

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)

docker-build:
	docker build -f build/Dockerfile -t $(DOCKER_REPO):$(VERSION) .
	docker tag $(DOCKER_REPO):$(VERSION) $(DOCKER_REPO):latest

docker-push:
	docker push $(DOCKER_REPO):$(VERSION)
	docker push $(DOCKER_REPO):latest

install-deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...

all: clean fmt lint test build