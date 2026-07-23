.PHONY: help test test-coverage fmt lint clean build install

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	go test ./...

test-coverage: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests (requires credentials)
	go test ./... -tags=integration -v

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter
	golangci-lint run

clean: ## Clean build artifacts
	go clean
	rm -f coverage.out coverage.html
	rm -rf bin/ dist/

build: ## Build the SDK
	go build ./...

install: ## Install the SDK locally
	go install ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

verify: ## Verify dependencies
	go mod verify

examples: ## Run examples
	@echo "Set environment variables before running examples:"
	@echo "export GHION_API_KEY=your_key"
	@echo "export GHION_API_SECRET=your_secret"
	@echo "export GHION_API_PASSPHRASE=your_passphrase"
	@echo ""
	@echo "Then run:"
	@echo "go run examples/quick-start.go"
	@echo "go run examples/http-server.go"
	@echo "go run examples/otp-flow.go"
	@echo "go run examples/qr-payment.go"

all: fmt lint test ## Run format, lint, and tests
