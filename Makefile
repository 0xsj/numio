# Makefile for numio - Natural Math Calculator

# ════════════════════════════════════════════════════════════════
# VARIABLES
# ════════════════════════════════════════════════════════════════

# Binary names
APP_NAME := numio
TUI_BIN := numio-tui
CLI_BIN := numio-cli

# Directories
CMD_DIR := ./cmd
BUILD_DIR := ./bin
TMP_DIR := ./tmp

# Go settings
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Version (can be overridden: make build VERSION=1.0.0)
VERSION ?= 0.1.0
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Linker flags with version info
LDFLAGS_VERSION := -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)

# ════════════════════════════════════════════════════════════════
# DEFAULT TARGET
# ════════════════════════════════════════════════════════════════

.PHONY: all
all: build

# ════════════════════════════════════════════════════════════════
# DEVELOPMENT
# ════════════════════════════════════════════════════════════════

## dev: Run TUI with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Starting hot reload..."
	@air -c .air.toml

## dev-cli: Run CLI with hot reload
.PHONY: dev-cli
dev-cli:
	@air -c .air.cli.toml

## run: Run TUI directly (no hot reload)
.PHONY: run
run:
	@$(GO) run $(CMD_DIR)/$(TUI_BIN)/main.go

## run-cli: Run CLI directly
.PHONY: run-cli
run-cli:
	@$(GO) run $(CMD_DIR)/$(CLI_BIN)/main.go

## repl: Run CLI in REPL mode
.PHONY: repl
repl:
	@$(GO) run $(CMD_DIR)/$(CLI_BIN)/main.go

# ════════════════════════════════════════════════════════════════
# BUILD
# ════════════════════════════════════════════════════════════════

## build: Build both TUI and CLI binaries
.PHONY: build
build: build-tui build-cli

## build-tui: Build TUI binary
.PHONY: build-tui
build-tui:
	@echo "Building $(TUI_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(BUILD_DIR)/$(TUI_BIN) $(CMD_DIR)/$(TUI_BIN)

## build-cli: Build CLI binary
.PHONY: build-cli
build-cli:
	@echo "Building $(CLI_BIN)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(BUILD_DIR)/$(CLI_BIN) $(CMD_DIR)/$(CLI_BIN)

## build-release: Build optimized release binaries
.PHONY: build-release
build-release:
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(BUILD_DIR)/$(TUI_BIN) $(CMD_DIR)/$(TUI_BIN)
	@CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(BUILD_DIR)/$(CLI_BIN) $(CMD_DIR)/$(CLI_BIN)

# ════════════════════════════════════════════════════════════════
# INSTALL
# ════════════════════════════════════════════════════════════════

## install: Install binaries to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(TUI_BIN) and $(CLI_BIN)..."
	@$(GO) install $(CMD_DIR)/$(TUI_BIN)
	@$(GO) install $(CMD_DIR)/$(CLI_BIN)

## uninstall: Remove installed binaries
.PHONY: uninstall
uninstall:
	@echo "Uninstalling..."
	@rm -f $(shell go env GOPATH)/bin/$(TUI_BIN)
	@rm -f $(shell go env GOPATH)/bin/$(CLI_BIN)

# ════════════════════════════════════════════════════════════════
# TEST
# ════════════════════════════════════════════════════════════════

## test: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test ./... -v

## test-short: Run tests without verbose output
.PHONY: test-short
test-short:
	@$(GO) test ./...

## test-cover: Run tests with coverage
.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	@$(GO) test ./... -coverprofile=coverage.out
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## test-race: Run tests with race detector
.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	@$(GO) test ./... -race

## bench: Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	@$(GO) test ./... -bench=. -benchmem

# ════════════════════════════════════════════════════════════════
# LINT & FORMAT
# ════════════════════════════════════════════════════════════════

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GO) vet ./...

## lint: Run golangci-lint (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## check: Run fmt, vet, and test
.PHONY: check
check: fmt vet test-short

# ════════════════════════════════════════════════════════════════
# DEPENDENCIES
# ════════════════════════════════════════════════════════════════

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download

## deps-update: Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	@$(GO) get -u ./...
	@$(GO) mod tidy

## deps-tidy: Tidy go.mod
.PHONY: deps-tidy
deps-tidy:
	@echo "Tidying dependencies..."
	@$(GO) mod tidy

## deps-verify: Verify dependencies
.PHONY: deps-verify
deps-verify:
	@echo "Verifying dependencies..."
	@$(GO) mod verify

# ════════════════════════════════════════════════════════════════
# TOOLS
# ════════════════════════════════════════════════════════════════

## tools: Install development tools
.PHONY: tools
tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed."

## tools-check: Check if required tools are installed
.PHONY: tools-check
tools-check:
	@echo "Checking tools..."
	@which air > /dev/null || (echo "air not found. Run: make tools" && exit 1)
	@echo "✓ air"
	@which golangci-lint > /dev/null || echo "⚠ golangci-lint not found (optional)"
	@echo "All required tools installed."

# ════════════════════════════════════════════════════════════════
# CLEAN
# ════════════════════════════════════════════════════════════════

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(TMP_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete."

## clean-cache: Clean Go build cache
.PHONY: clean-cache
clean-cache:
	@echo "Cleaning Go cache..."
	@$(GO) clean -cache -testcache

# ════════════════════════════════════════════════════════════════
# HELP
# ════════════════════════════════════════════════════════════════

## help: Show this help message
.PHONY: help
help:
	@echo ""
	@echo "numio - Natural Math Calculator"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## /  /' | column -t -s ':'
	@echo ""