# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-08-22

### Added
- `SendPaymentReminder` method to send payment reminders to customers via email and SMS
- `GenerateBillID` method to generate suggested auto-generated bill IDs before creating bills
- `InitiateCheckout` method to initiate checkout for bills and ensure payment links exist
- Made `BillID` field optional in `CreateBillRequest` - if omitted, the API will auto-generate one
- `SendPaymentReminderRequest` and `SendPaymentReminderResponse` types
- `GenerateBillIDResponse` type
- `InitiateCheckoutResponse` type
- Validator for `SendPaymentReminderRequest`

### Changed
- Updated bill creation validator to accept optional BillID
- Updated documentation with examples for new endpoints
- Added integration tests for all new bill payment endpoints
- Enhanced developer guide with usage examples for GenerateBillID and InitiateCheckout

## [1.1.0] - 2026-08-07

### Added
- Complete Bill Payment API support for creating and managing bills
- `CreateBill` method to create single bills with customer information, due dates, and penalty configurations
- `CreateBulkBills` method to create multiple bills in a single request for batch billing cycles
- `ListBills` method to list bills with filters (status, search, cluster, bill_code, date range, pagination)
- `GetBillStatistics` method to retrieve aggregate counts and amounts for bills
- `GetBillDashboard` method to get detailed analytics including summary statistics, trends, and breakdowns by cluster/bill code
- `GetBillDetail` method to retrieve full bill details including penalty configuration, metadata, and payment history
- `UpdateBill` method to update bill properties (partial update - only provided fields are updated)
- `DeleteBill` method to delete bills (only bills with no payments can be deleted)
- `RecordManualPayment` method to record manual payments (cash, bank transfer, etc.) for reconciliation purposes
- `GetBillPaymentLink` method to generate shareable payment links for bills
- `GetBillerSettings` method to retrieve biller configuration including biller_code, clusters, bill codes, and webhook settings
- `UpdateBillerSettings` method to update biller configuration
- `PublicBillLookup` method for public bill lookup (no authentication) used by bank branches and mobile banking apps
- `BillStatus` enum (pending, paid, forwarded, cancelled, expired, overdue)
- `Penalty` interface for late payment penalty configuration
- Comprehensive bill payment types: CreateBillRequest, BillResponse, BulkCreateBillsRequest, BulkCreateBillsResponse, ListBillsRequest, ListBillsResponse, BillStatistics, BillDashboard, PaymentLinkResponse, PublicBillLookupRequest, PublicBillLookupResponse, BillerSettingsRequest, BillerSettingsResponse, UpdateBillRequest, DeleteBillResponse, RecordManualPaymentRequest, RecordManualPaymentResponse
- Validators for all bill payment requests with date format validation (Y-m-d), email validation, and penalty configuration validation
- Bill payment examples and integration tests

### Fixed
- Fixed signature generation for GET requests to exclude query parameters (matching API specification)
- Simplified webhook handling by removing redundant wrapper and using webhook package directly

### Security
- HMAC-SHA256 signature generation and verification
- Timing-safe comparison for webhook signatures
- Input validation for all API requests
- Automatic sensitive data redaction

## [1.0.0] - 2026-07-23

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

---

## Versioning

For the versions available, see the [tags on this repository](https://github.com/yayawallet/ghion-go-sdk/tags).
