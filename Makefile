.PHONY: help build test test-verbose test-coverage lint fmt clean install run dev

# Variables
BINARY_NAME=gomodsync
BIN_DIR=bin
COVERAGE_FILE=coverage.out

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "Coverage report generated: $(COVERAGE_FILE)"
	@$(GOCMD) tool cover -func=$(COVERAGE_FILE)

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	@$(GOCMD) tool cover -html=$(COVERAGE_FILE)

lint: ## Run linters
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run --timeout=5m

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@which goimports > /dev/null && goimports -w . || echo "goimports not found, skipping..."

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BIN_DIR)
	@rm -f $(COVERAGE_FILE)

install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BIN_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

run: build ## Build and run with example arguments
	@echo "Running $(BINARY_NAME)..."
	@$(BIN_DIR)/$(BINARY_NAME)

dev: ## Run tests, lint, and build
	@$(MAKE) test
	@$(MAKE) lint
	@$(MAKE) build

ci: ## Run CI checks (test, lint, build)
	@$(MAKE) test-coverage
	@$(MAKE) lint
	@$(MAKE) build

all: clean dev ## Clean, test, lint, and build

.DEFAULT_GOAL := help
