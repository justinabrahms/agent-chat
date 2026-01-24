# Makefile for agent-chat
# Run 'make' or 'make help' to see available targets

BINARY_NAME := agent-chat
BINARY_DIR := ./bin
CMD_PATH := ./cmd/agent-chat
GO := go

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION_PKG := github.com/justinabrahms/agent-chat/internal/version
LDFLAGS := -ldflags "-X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) -X $(VERSION_PKG).BuildDate=$(BUILD_DATE)"

# Detect OS and architecture
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

.PHONY: all help build build-linux build-darwin test lint clean install

all: help

##@ General

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

build: ## Build the binary for current platform
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Built $(BINARY_DIR)/$(BINARY_NAME) for $(GOOS)/$(GOARCH)"

test: ## Run all tests with verbose output
	$(GO) test -v ./...

lint: ## Run linters (go vet and staticcheck if available)
	$(GO) vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed, skipping (install with: go install honnef.co/go/tools/cmd/staticcheck@latest)"; \
	fi

##@ Build Variants

build-linux: ## Cross-compile for Linux amd64
	@mkdir -p $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@echo "Built $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64"

build-darwin: ## Cross-compile for macOS (darwin) amd64
	@mkdir -p $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@echo "Built $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64"

build-all: build-linux build-darwin ## Build for all supported platforms
	@echo "Built all platform binaries"

##@ Installation

install: build ## Install binary to GOPATH/bin
	@install -d $(shell $(GO) env GOPATH)/bin
	@install -m 755 $(BINARY_DIR)/$(BINARY_NAME) $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(shell $(GO) env GOPATH)/bin"

##@ Cleanup

clean: ## Remove build artifacts
	@rm -rf $(BINARY_DIR)
	@echo "Cleaned $(BINARY_DIR)"
