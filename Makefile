.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the cicada binary
	@echo "Building cicada..."
	@go build -o bin/cicada ./cmd/cicada

.PHONY: install
install: ## Install cicada to $GOPATH/bin
	@echo "Installing cicada..."
	@go install ./cmd/cicada

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Linting targets
.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running linters..."
	@golangci-lint run --config .golangci.yml

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linters with auto-fix..."
	@golangci-lint run --fix --config .golangci.yml

.PHONY: fmt
fmt: ## Run go fmt on all files
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: reportcard
reportcard: ## Check Go Report Card grade (requires goreportcard-cli)
	@echo "Running Go Report Card checks..."
	@goreportcard-cli -v || echo "Install: go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest"

# Testing targets
.PHONY: test
test: ## Run unit tests
	@echo "Running tests..."
	@go test -v -race ./...

.PHONY: test-short
test-short: ## Run unit tests (short mode)
	@go test -v -short ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-coverage
test-coverage: ## Show test coverage percentage
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

.PHONY: test-integration
test-integration: ## Run integration tests with AWS (requires .env.test)
	@echo "Running integration tests..."
	@test -f .env.test || (echo "Error: .env.test not found. Copy .env.test.example" && exit 1)
	@export $$(grep -v '^\#' .env.test | xargs) && \
		go test -v -tags=integration -timeout=10m ./...

.PHONY: test-integration-setup
test-integration-setup: ## Set up integration test infrastructure
	@echo "Setting up integration test infrastructure..."
	@test -f .env.test || (echo "Error: .env.test not found. Copy .env.test.example" && exit 1)
	@export $$(grep -v '^\#' .env.test | xargs) && \
		aws s3 mb s3://$${CICADA_TEST_BUCKET} --region $${AWS_REGION} 2>/dev/null || true
	@echo "✅ Test bucket ready"

# Development targets
.PHONY: dev
dev: ## Run in development mode
	@go run ./cmd/cicada

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	@echo "Tidying dependencies..."
	@go mod tidy

.PHONY: verify
verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

# Quality checks
.PHONY: check
check: lint test ## Run all checks (lint + test)
	@echo "✅ All checks passed"

.PHONY: ci
ci: deps lint test-short ## Run CI checks (fast)
	@echo "✅ CI checks passed"

# Default target
.DEFAULT_GOAL := help
