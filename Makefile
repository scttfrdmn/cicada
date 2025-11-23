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
test-integration: ## Run integration tests with AWS (set AWS_PROFILE=aws AWS_REGION=us-west-2)
	@echo "Running integration tests..."
	@echo "Using AWS_PROFILE=$${AWS_PROFILE:-aws}, AWS_REGION=$${AWS_REGION:-us-west-2}"
	@AWS_PROFILE=$${AWS_PROFILE:-aws} AWS_REGION=$${AWS_REGION:-us-west-2} \
		go test -v -tags=integration -timeout=10m ./internal/integration/...
	@echo "✅ Integration tests passed"

.PHONY: test-integration-setup
test-integration-setup: ## Set up integration test S3 bucket (cicada-integration-test)
	@echo "Setting up integration test infrastructure..."
	@AWS_PROFILE=$${AWS_PROFILE:-aws} AWS_REGION=$${AWS_REGION:-us-west-2} \
		aws s3 mb s3://cicada-integration-test --region $${AWS_REGION:-us-west-2} 2>/dev/null || true
	@echo "✅ Test bucket ready: s3://cicada-integration-test"

.PHONY: test-all
test-all: test test-integration ## Run all tests (unit + integration)

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
