# Makefile for concourse-replicated-channel-resource

# Variables
BINARY_NAME = concourse-replicated-channel-resource
PACKAGE_NAME = $(BINARY_NAME)
IMAGE_NAME = $(BINARY_NAME)
TTL_REGISTRY = ttl.sh
TTL_IMAGE = $(TTL_REGISTRY)/$(IMAGE_NAME):4h

# Output directories
BIN_DIR = bin
DIST_DIR = dist

# Source files
GO_FILES = $(shell find . -name "*.go" -type f)
GO_MOD_FILES = go.mod go.sum

# Default target
.PHONY: all
all: build

# Test target
.PHONY: test
test: $(GO_FILES) $(GO_MOD_FILES)
	go test -v ./...

# Format target
.PHONY: format
format: $(GO_FILES)
	go fmt ./...

# Lint target
.PHONY: lint
lint: $(GO_FILES)
	golangci-lint run ./...

# Build target
.PHONY: build
build: $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME): $(GO_FILES) $(GO_MOD_FILES)
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Package target using melange
.PHONY: package
package: $(DIST_DIR)/x86_64/$(PACKAGE_NAME)-0.1.0-r0.apk

$(DIST_DIR)/x86_64/$(PACKAGE_NAME)-0.1.0-r0.apk: $(BIN_DIR)/$(BINARY_NAME) build/melange.yaml
	mkdir -p $(DIST_DIR)
	melange build build/melange.yaml --arch x86_64 --out-dir $(DIST_DIR)/ --source-dir .

# Image target using apko
.PHONY: image
image: $(DIST_DIR)/$(IMAGE_NAME).tar

$(DIST_DIR)/$(IMAGE_NAME).tar: $(DIST_DIR)/x86_64/$(PACKAGE_NAME)-0.1.0-r0.apk build/apko.yaml
	mkdir -p $(DIST_DIR)/x86_64/x86_64
	cp $(DIST_DIR)/x86_64/APKINDEX.tar.gz $(DIST_DIR)/x86_64/x86_64/
	cp $(DIST_DIR)/x86_64/*.apk $(DIST_DIR)/x86_64/x86_64/
	apko build build/apko.yaml $(IMAGE_NAME) $(DIST_DIR)/$(IMAGE_NAME).tar --ignore-signatures -r $(DIST_DIR)/x86_64 -p concourse-replicated-channel-resource

# Deploy to ttl.sh
.PHONY: ttl.sh
ttl.sh: $(DIST_DIR)/$(IMAGE_NAME).tar
	docker load < $(DIST_DIR)/$(IMAGE_NAME).tar
	docker tag $(IMAGE_NAME):latest-amd64 $(TTL_IMAGE)
	docker push $(TTL_IMAGE)
	@echo "Image deployed to: $(TTL_IMAGE)"

# Clean target
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/ $(DIST_DIR)/

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Build the binary (default)"
	@echo "  test        - Run tests"
	@echo "  format      - Format Go code"
	@echo "  lint        - Run linter"
	@echo "  build       - Build the binary"
	@echo "  package     - Package using melange"
	@echo "  image       - Build OCI image using apko"
	@echo "  ttl.sh      - Deploy image to ttl.sh"
	@echo "  clean       - Clean build artifacts"
	@echo "  help        - Show this help"