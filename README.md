# github.com/yayawallet/ghion-go-sdk

A type-safe Go SDK for the Ghion Finances payment gateway. Built with security, maintainability, and scalability in mind for developers maintaining or contributing to this repository.

## Features

- **Type-Safe**: Full Go type definitions with comprehensive structs
- **Secure**: HMAC-SHA256 authentication with timing-safe signature verification
- **Robust**: Comprehensive error handling with custom error types
- **Validated**: Built-in input validation for all API requests
- **Webhook Support**: Secure webhook signature verification and parsing
- **Modern**: Built with modern Go best practices
- **Multi-Channel**: Support for USSD, QR, and OTP payment methods
- **Retry Logic**: Automatic retry for transient failures and rate limits

## Repository Structure

```text
ghion-go-sdk/
├── ghion.go              # Main entry point and public API aliases
├── pkg/
│   ├── client/           # HTTP client and API methods
│   ├── types/            # Type definitions and structs
│   ├── errors/           # Custom error types
│   ├── utils/            # Utility functions (crypto, validation)
│   └── webhook/          # Webhook handling
├── examples/             # Usage examples (HTTP server, OTP, QR, Payment Lifecycle)
└── tests/                # Test files (unit and integration)
```

## Setup for Development

### 1. Prerequisites
- Go 1.19 or higher
- Git

### 2. Clone the Repository
```bash
git clone https://github.com/yayawallet/ghion-go-sdk.git
cd ghion-go-sdk
```

### 3. Install Dependencies
```bash
go mod tidy
```

## Adding New Features

### 1. Adding a New API Endpoint
1. **Define the types** in `pkg/types/types.go` (Requests, Responses, Enums).
2. **Add validation** in `pkg/utils/validator.go` if necessary.
3. **Implement the method** in `pkg/client/client.go` using the generic `apiRequest` method.
4. **Export the types/methods** in `ghion.go` via aliases.
5. **Add unit tests** in `pkg/client/client_test.go`.
6. **Add integration tests** in `tests/integration_test.go` (if applicable).

### 2. Adding a New Error Type
1. **Define the error** in `pkg/errors/errors.go`.
2. **Export the error** in `ghion.go`.
3. **Add tests** in `pkg/errors/errors_test.go`.

## Testing the SDK

### Running All Tests
Run all tests across the SDK (unit tests only):
```bash
go test ./... -v
```

### Running Package-Specific Tests
Run tests for specific packages:
```bash
# Client tests
go test ./pkg/client -v

# Utils tests (crypto, validation)
go test ./pkg/utils -v

# Error handling tests
go test ./pkg/errors -v

# Webhook tests
go test ./pkg/webhook -v
```

### Test Coverage
Generate a test coverage report:
```bash
go test ./pkg/... -coverprofile=coverage.out -covermode=atomic
```

View coverage summary:
```bash
go tool cover -func=coverage.out
```

Generate HTML coverage report:
```bash
go tool cover -html=coverage.out
```

**Current Coverage:**
- `pkg/client`: 68.5% (validation and configuration logic)
- `pkg/errors`: 95.0% (error type constructors)
- `pkg/utils`: 92.1% (crypto and validation functions)
- `pkg/webhook`: 97.5% (signature verification and parsing)

### Unit Tests (No Credentials Required)
Unit tests verify SDK logic without making API calls:
- Client configuration and instantiation
- Input validation functions
- Cryptographic signature generation
- Webhook signature verification
- Error type constructors
- Helper functions

Run unit tests:
```bash
go test ./pkg/... -v
```

### Integration Tests (Credentials Required)
Integration tests make real API calls to verify SDK functionality with the Ghion API.

**Setup:**
1. Create a `.env` file in the project root:
```env
GHION_API_KEY=your_api_key
GHION_API_SECRET=your_api_secret
GHION_API_PASSPHRASE=your_passphrase
TEST_PHONE_NUMBER=+251911234567
TEST_OTP_CODE=123456  # For OTP validation test
```

2. Install the `godotenv` package for loading environment variables:
```bash
go get github.com/joho/godotenv
```

3. Run integration tests:
```bash
go test ./tests/ -v -tags=integration
```

**Integration Test Coverage:**
- **Initialize Payment**: Tests payment initialization and channel availability
- **QR Payment**: Tests QR code generation and checkout retrieval
- **OTP Payment**: Tests OTP sending and validation (requires phone number)
- **YaYa Wallet Payment**: Tests YaYa Wallet USSD payment flow
- **Telebirr Payment**: Tests Telebirr USSD payment flow
- **USSD Payment**: Tests generic USSD payment flow
- **Payment Status**: Tests payment status retrieval
- **Get Checkout**: Tests checkout details and available channels

**Available Payment Channels:**
- YaYa Wallet (yayawallet)
- Card (card)
- Telebirr (telebirr)
- CBE Birr (cbebirr)
- M-PESA (mpesa)
- Kacha (kacha)
- CBE Mobile (cbemobile)
- Awash Birr (awashbirr)
- Hijra Bank (hijirabank)
- Ahadu (ahadu)
- NIB Bank (NIB)

**Test Organization:**
- Unit tests are located in `pkg/*/` directories alongside source code
- Integration tests are located in `tests/` directory with `// +build integration` tag
- Tests are organized by package for accurate coverage reporting

### Writing Tests
Tests are organized by package in the `pkg/` directory:

**Test Structure:**
```go
package client

import (
    "testing"
    "github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func TestNewGhionClient(t *testing.T) {
    tests := []struct {
        name        string
        config      *types.GhionConfig
        expectError bool
    }{
        {
            name: "valid config",
            config: &types.GhionConfig{
                APIKey:     "test-key",
                APISecret:  "test-secret",
                Passphrase: "test-passphrase",
            },
            expectError: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := NewGhionClient(tt.config)
            if tt.expectError && err == nil {
                t.Error("Expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
            if !tt.expectError && client == nil {
                t.Error("Expected client but got nil")
            }
        })
    }
}
```

**Testing Guidelines:**
- Use table-driven tests for multiple test cases
- Test both success and error paths
- Validate error types when expecting errors
- Test edge cases (empty strings, zero values, invalid inputs)
- Mock external dependencies when possible
- Keep tests independent and deterministic

### Running Examples
See the `examples/` directory for complete working scripts:

```bash
# Basic quick start example
go run examples/quick-start.go

# Comprehensive payment lifecycle test
go run examples/payment-lifecycle.go

# HTTP server with webhook handling
go run examples/http-server.go

# OTP flow example
go run examples/otp-flow.go
```

### CI/CD Testing
For continuous integration, run:
```bash
# Run all tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic

# Check coverage threshold (e.g., 70%)
go tool cover -func=coverage.out | grep total
```

## Code Quality & Guidelines

- **Style**: Use `gofmt` to format code before committing (`go fmt ./...`).
- **Linting**: We recommend using `golangci-lint` to ensure code quality.
- **Documentation**: All exported functions, types, and constants must have proper Go doc comments.
- **Error Handling**: Use the custom error types in `pkg/errors` instead of generic Go errors. Never expose sensitive information in error messages. Redact sensitive data from logs.

## Security Considerations

- **HMAC-SHA256 Authentication**: All API requests are signed using HMAC-SHA256.
- **Timing-Safe Comparison**: Webhook signatures use constant-time comparison (`subtle.ConstantTimeCompare`) to prevent timing attacks.
- **Input Validation**: All inputs are validated before being sent over the network.
- **Sensitive Data Redaction**: Error responses automatically redact sensitive information (API keys, passphrases, signatures).

## Contributing

Contributions are welcome! Please ensure:
1. Code adheres to existing style (`go fmt`).
2. All tests pass (`go test ./...`).
3. Documentation is updated.
4. Changes are backwards compatible when possible.

Please open an issue to discuss proposed changes before creating a pull request.

For release guidelines, see [RELEASE.md](RELEASE.md).

## License

See [LICENSE](LICENSE) file for details.
