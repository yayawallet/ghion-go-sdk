# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Ghion Go SDK
- Payment initialization API
- Payment submission API
- Payment status checking
- QR code payment support
- OTP payment flow support
- Webhook signature verification and parsing
- Comprehensive input validation
- HMAC-SHA256 authentication
- Automatic retry logic for transient failures
- Rate limit handling with exponential backoff
- Sensitive data redaction in error responses
- Ethiopian phone number validation
- Complete test suite
- Documentation and examples
- Developer guide

### Security
- HMAC-SHA256 signature generation and verification
- Timing-safe comparison for webhook signatures
- Input validation for all API requests
- Automatic sensitive data redaction

## [1.0.0] - 2024-01-XX

### Added
- Initial stable release
- Full API coverage for payment operations
- Webhook handling
- Comprehensive error handling
- Security features
- Documentation and examples

---

## Versioning

For the versions available, see the [tags on this repository](https://github.com/yayawallet/ghion-go-sdk/tags).
