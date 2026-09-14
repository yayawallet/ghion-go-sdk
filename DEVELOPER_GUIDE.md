# Ghion Finances Go SDK - Developer Guide

Welcome to the official developer guide for the `ghion-go-sdk`. This document provides comprehensive instructions on setting up, implementing, and troubleshooting the SDK in your Go applications.

## Table of Contents
1. [Installation & Setup](#installation--setup)
2. [Core Concepts](#core-concepts)
3. [Implementation Guide](#implementation-guide)
   - [Initializing a Payment](#1-initializing-a-payment)
   - [Handling Payment Methods](#2-handling-payment-methods)
     - [OTP Flow (YaYa Wallet)](#otp-flow-yaya-wallet)
     - [USSD Flow](#ussd-flow)
     - [QR Code Flow](#qr-code-flow)
   - [Bill Payment API](#4-bill-payment-api)
   - [Hold Payment (Escrow)](#5-hold-payment-escrow)
   - [Pay Merchant (Direct Pay)](#6-pay-merchant-direct-pay)
   - [Webhooks Integration](#3-webhooks-integration)
   - [Error Handling](#7-error-handling)
4. [Best Practices](#best-practices)
5. [Common Issues & Fixes](#common-issues--fixes)

---

## Installation & Setup

### 1. Install the SDK
Install the package via `go get`:

```bash
go get github.com/yayawallet/ghion-go-sdk
```

### 2. Environment Variables
You will need your API credentials from the Ghion Developer Dashboard. Securely store them in your environment (e.g., using a `.env` file):

```env
GHION_API_KEY=your_api_key_here
GHION_API_SECRET=your_api_secret_here
GHION_API_PASSPHRASE=your_passphrase_here
WEBHOOK_URL=https://your-domain.com/webhook
```

### 3. Initialize the Client
Import and initialize the `GhionClient` in your application:

```go
package main

import (
	"log"
	"os"

	"github.com/yayawallet/ghion-go-sdk"
)

func main() {
	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     os.Getenv("GHION_API_KEY"),
		APISecret:  os.Getenv("GHION_API_SECRET"),
		Passphrase: os.Getenv("GHION_API_PASSPHRASE"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}
}
```

---

## Core Concepts

The SDK revolves around a few key resources:
- **Payment Session:** Created when a user initiates a checkout. Represents the transaction lifecycle.
- **Channels:** Different payment methods available (e.g., YaYa Wallet, Card, Telebirr).
- **Webhooks:** The primary, asynchronous mechanism for receiving definitive payment statuses (Success, Failure, Expiry).

---

## Implementation Guide

### 1. Initializing a Payment
When a user clicks "Checkout", you must create a payment session. This returns the available channels, provider information, and a unique `payment.ID`.

```go
import (
    "fmt"
    "time"
    "github.com/yayawallet/ghion-go-sdk/pkg/types"
)

// Inside your payment handler:
payment, err := client.InitializePayment(&types.InitializePaymentRequest{
    Amount:      100,
    Currency:    "ETB", // Default is ETB
    Reference:   fmt.Sprintf("order_%d", time.Now().Unix()), // Your internal unique order ID
    Description: "Purchase of premium coffee",
    WebhookURL:  os.Getenv("WEBHOOK_URL"),
    ReturnURL:   "https://your-domain.com/success",
    CancelURL:   "https://your-domain.com/cancel",
})

if err != nil {
    log.Printf("Error initializing payment: %v", err)
    return
}

// Optionally fetch full checkout details (for QR codes, specific provider rules, etc.)
checkoutInfo, err := client.GetCheckout(payment.ID)

fmt.Printf("Payment Initialized. ID: %s", payment.ID)
// Return payment.ID and checkoutInfo.AvailableChannels to your frontend
```

### 2. Handling Payment Methods

#### OTP Flow (YaYa Wallet)
The OTP flow requires sending an OTP to the user's phone, then validating it.

**Step A: Send OTP**
```go
otpResponse, err := client.SendOTP(payment.ID, phoneNumber)
if err != nil {
    // Handle error
}
// Display an input field to the user to enter the 6-digit OTP
```

**Step B: Validate OTP**
```go
validateResponse, err := client.ValidateOTP(payment.ID, otpCode, phoneNumber)
if err != nil {
    // Handle error
}
if validateResponse.Status == string(ghion.StatusCompleted) {
    // Payment is successful
}
```

#### USSD Flow
For wallets supporting direct push prompts (USSD):

```go
result, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
    Channel:       "yayawallet", // or 'telebirr', etc.
    PhoneNumber:   "0912345678",
})
if err != nil {
    // Handle error
}
// User will receive a prompt on their phone to enter their PIN.
```

#### QR Code Flow
To display a QR code for the user to scan with their banking app:

```go
qrPayment, err := client.PayWithQR(payment.ID)
if err != nil {
    // Handle error
}
// qrPayment.QRImageURL contains the URL to the generated QR code image
// qrPayment.QRPayload contains the raw string payload for the QR code
```

### 4. Bill Payment API

The Bill Payment API allows you to create and manage bills for recurring payments, utility bills, invoices, and other billing scenarios.

#### Create a Bill

Create a single bill with customer information:

```go
bill, err := client.CreateBill(&types.CreateBillRequest{
    BillID:        "INV-12345", // Optional: If omitted, one will be auto-generated
    Amount:        500.00,
    Currency:      "ETB",
    DueDate:       "2026-09-01",
    CustomerName:  "John Doe",
    CustomerPhone: "+251911234567",
    CustomerEmail: "john@example.com",
    Description:   "Monthly utility bill",
    BillCode:      "UTIL",
    Cluster:       "ADDIS ABABA",
})
if err != nil {
    // Handle error
}
```

**Required Parameters:**
- `amount`: Bill amount
- `dueDate`: Due date in Y-m-d format
- `customerName`: Customer name

**Optional Parameters:**
- `billID`: Unique bill identifier (if omitted, one will be auto-generated)
- `currency`: Currency code (default: ETB)
- `customerPhone`: Customer phone number
- `customerEmail`: Customer email
- `description`: Bill description
- `billCode`: Bill code for categorization
- `cluster`: Geographic cluster

**Important Notes:**
- `billID` is now optional. If omitted, the API will auto-generate one
- You can use `GenerateBillID()` to get a suggested auto-generated bill ID before creating a bill

#### Generate Bill ID

Generate a suggested auto-generated bill ID before creating a bill:

```go
generatedID, err := client.GenerateBillID()
if err != nil {
    // Handle error
}

fmt.Printf("Generated Bill ID: %s\n", generatedID.BillID)
```

**Use Case:** Use this to preview what the auto-generated bill ID would be, or to generate IDs in bulk before creating bills.

#### Create Bulk Bills

Create multiple bills in a single request for batch billing cycles:

```go
response, err := client.CreateBulkBills(&types.BulkCreateBillsRequest{
    Bills: []types.CreateBillRequest{
        {
            BillID:       "BULK-1",
            Amount:       300.00,
            DueDate:      "2026-09-01",
            CustomerName: "Customer 1",
        },
        {
            BillID:       "BULK-2",
            Amount:       400.00,
            DueDate:      "2026-09-01",
            CustomerName: "Customer 2",
        },
    },
})
if err != nil {
    // Handle error
}
fmt.Printf("Created: %d, Errors: %d\n", response.CreatedCount, response.ErrorCount)
```

#### List Bills

List bills with filters for pagination and search:

```go
response, err := client.ListBills(&types.ListBillsRequest{
    Status:   "pending",
    Page:     1,
    Limit:    10,
    From:     "2026-08-01",
    To:       "2026-08-31",
})
if err != nil {
    // Handle error
}

for _, bill := range response.Items {
    fmt.Printf("%s: %.2f %s (%s)\n", bill.BillID, bill.Amount, bill.Currency, bill.Status)
}
```

**Filter Parameters:**
- `status`: Filter by status (pending, paid, forwarded, cancelled, expired, overdue)
- `search`: Search by customer name or bill ID
- `cluster`: Filter by cluster
- `billCode`: Filter by bill code
- `from`: Start date (Y-m-d format)
- `to`: End date (Y-m-d format)
- `page`: Page number
- `limit`: Items per page (max 100)

#### Get Bill Statistics

Get aggregate statistics for all bills:

```go
stats, err := client.GetBillStatistics()
if err != nil {
    // Handle error
}

fmt.Printf("Pending: %d\n", stats.Pending)
fmt.Printf("Paid: %d\n", stats.Paid)
fmt.Printf("Total Amount: %.2f\n", stats.TotalAmount)
```

#### Get Bill Dashboard

Get detailed analytics including trends and breakdowns:

```go
dashboard, err := client.GetBillDashboard("2026-08-01", "2026-08-31")
if err != nil {
    // Handle error
}

fmt.Printf("Total Bills: %d\n", dashboard.Summary.TotalBills)
fmt.Printf("Pending: %d\n", dashboard.Summary.Pending)
fmt.Printf("Paid: %d\n", dashboard.Summary.Paid)
```

#### Get Bill Details

Retrieve full bill details including payment history:

```go
detail, err := client.GetBillDetail(bill.ID)
if err != nil {
    // Handle error
}

fmt.Printf("Customer: %s\n", detail.CustomerName)
fmt.Printf("Amount: %.2f\n", detail.Amount)
fmt.Printf("Status: %s\n", detail.Status)
fmt.Printf("Payments: %d\n", len(detail.Payments))
```

#### Update Bill

Update bill properties (partial update - only provided fields are updated):

```go
updatedBill, err := client.UpdateBill(bill.ID, &types.UpdateBillRequest{
    Amount:      600.00,
    Description: "Updated description",
    DueDate:     "2026-09-15",
})
if err != nil {
    // Handle error
}
```

#### Record Manual Payment

Record manual payments (cash, bank transfer, etc.) for reconciliation:

```go
payment, err := client.RecordManualPayment(bill.ID, &types.RecordManualPaymentRequest{
    Amount:        100.00,
    Source:        "manual",
    PaymentMethod: "cash",
    Reference:     "RECEIPT-001",
    Note:          "Paid at counter",
})
if err != nil {
    // Handle error
}

fmt.Printf("Bill Status: %s\n", payment.BillStatus)
fmt.Printf("Balance Due: %.2f\n", payment.BalanceDue)
```

#### Send Payment Reminder

Send payment reminders to customers via email and SMS to encourage timely payments:

```go
reminder, err := client.SendPaymentReminder(bill.ID, &types.SendPaymentReminderRequest{
    Message: "Please pay your bill before the due date.",
})
if err != nil {
    // Handle error
}

fmt.Printf("Reminder Sent: %v\n", reminder.Sent)
fmt.Printf("Reminder Count: %d\n", reminder.ReminderCount)
fmt.Printf("Last Reminder Sent At: %s\n", reminder.LastReminderSentAt)
```

**Important Notes:**
- The message parameter is optional - if not provided, a default reminder message will be sent
- Requires the `bill_payment.reminders.send` permission
- The response includes the total reminder count and timestamp of the last reminder sent
- Reminders are sent via both email and SMS to the customer

#### Get Bill Payment Link

Generate a shareable payment link for bills:

```go
link, err := client.GetBillPaymentLink(bill.ID)
if err != nil {
    // Handle error
}

fmt.Printf("Payment Link: %s\n", link.CheckoutURL)
// Share this link with customers via SMS, email, etc.
```

#### Delete Bill

Delete bills (only bills with no payments can be deleted):

```go
result, err := client.DeleteBill(bill.ID)
if err != nil {
    // Handle error
}
fmt.Printf("Message: %s\n", result.Message)
```

#### Public Bill Lookup

Lookup bills without authentication (for public-facing apps like bank branches or mobile banking):

```go
publicBill, err := client.PublicBillLookup(&types.PublicBillLookupRequest{
    BillerCode: "BILLER001",
    BillID:     "INV-12345",
})
if err != nil {
    // Handle error
}

fmt.Printf("Payment Status: %s\n", publicBill.PaymentStatus)
fmt.Printf("Amount Due: %.2f\n", publicBill.AmountDue)
```

#### Initiate Checkout

Initiate checkout for a bill to ensure a payment link exists. This creates a payment link if one doesn't already exist:

```go
checkout, err := client.InitiateCheckout(bill.ID)
if err != nil {
    // Handle error
}

fmt.Printf("Payment Link Slug: %s\n", checkout.PaymentLinkSlug)
fmt.Printf("Balance Due: %.2f %s\n", checkout.BalanceDue, checkout.Currency)
fmt.Printf("Checkout URL: %s\n", checkout.CheckoutURL)
```

**Important Notes:**
- This endpoint ensures a payment link exists for the bill (creates if needed)
- Returns the checkout URL directly, which can be shared with customers
- Useful when you need to programmatically generate payment links for bills
- The `payment_link_slug` can be used to construct custom URLs

### 5. Hold Payment (Escrow)

Hold Payment (Escrow) allows customers to hold funds until you pull them. This is useful for marketplace scenarios where you want to ensure funds are available before releasing them to sellers.

**Note:** These features require module enablement by Ghion support. Contact Ghion support to enable Hold Payment for your account.

#### Listing Escrows

List all held payments (escrows) for your account, optionally filtered by status:

```go
// List all escrows
escrows, err := client.ListEscrows(nil)
if err != nil {
    // Handle error
}

for _, escrow := range escrows.Escrows {
    fmt.Printf("Escrow ID: %s, Status: %s, Amount: %s %s\n",
        escrow.ID, escrow.Status, escrow.Amount, escrow.Currency)
}

// List escrows with status filter
fundedEscrows, err := client.ListEscrows(&ghion.ListEscrowsRequest{
    Status: ghion.EscrowStatusFunded,
})
if err != nil {
    // Handle error
}
```

**Escrow Status Values:**
- `EscrowStatusFunded`: Funds are held and available to pull
- `EscrowStatusWithdrawing`: Funds are being withdrawn
- `EscrowStatusWithdrawn`: Funds have been withdrawn
- `EscrowStatusReleased`: Funds have been released
- `EscrowStatusCancelled`: Escrow was cancelled

#### Getting Escrow Details

Retrieve a single escrow by its ID:

```go
escrow, err := client.GetEscrow("escrow-123")
if err != nil {
    // Handle error
}

fmt.Printf("Escrow ID: %s\n", escrow.ID)
fmt.Printf("Wallet: %s (%s)\n", escrow.WalletName, escrow.WalletPhone)
fmt.Printf("Amount: %s %s\n", escrow.Amount, escrow.Currency)
fmt.Printf("Status: %s\n", escrow.Status)
fmt.Printf("Purpose: %s\n", escrow.Purpose)
```

#### Pulling Escrow Funds

Pull funds from a funded escrow to your balance:

```go
response, err := client.PullEscrowFunds("escrow-123")
if err != nil {
    // Handle error
}

fmt.Printf("Escrow ID: %s\n", response.ID)
fmt.Printf("New Status: %s\n", response.Status)
fmt.Printf("Amount Released: %s %s\n", response.Amount, response.Currency)
fmt.Printf("Released At: %s\n", response.ReleasedAt)
fmt.Printf("Released By: %s\n", response.ReleasedBy)
```

**Important Notes:**
- Only escrows with status `funded` can have funds pulled
- After pulling, the escrow status changes to `released`
- This operation transfers funds from the escrow to your merchant balance

#### Webhook Events for Escrows

The following webhook events are sent for escrow operations:
- `escrow.funded`: When funds are held in escrow
- `escrow.released`: When funds are pulled from escrow
- `escrow.cancelled`: When an escrow is cancelled

### 6. Pay Merchant (Direct Pay)

Pay Merchant (Direct Pay) allows wallet users to pay merchants directly using their customer ID and reference. You can configure customer validation via HTTP adapter to verify customer details before payment.

**Note:** These features require module enablement by Ghion support. Contact Ghion support to enable Pay Merchant for your account.

#### Getting Direct Pay Settings

Retrieve current Pay Merchant settings:

```go
settings, err := client.GetDirectPaySettings()
if err != nil {
    // Handle error
}

fmt.Printf("Configured: %v\n", settings.Configured)
fmt.Printf("Customer ID Required: %v\n", settings.Settings.CustomerIDRequired)
fmt.Printf("Reference Required: %v\n", settings.Settings.ReferenceRequired)
fmt.Printf("Validation Adapter: %s\n", settings.Settings.ValidationAdapter)
fmt.Printf("Validation URL: %s\n", settings.Settings.ValidationURL)
fmt.Printf("Validation Timeout: %d seconds\n", settings.Settings.ValidationTimeout)
```

#### Updating Direct Pay Settings

Configure Pay Merchant settings including customer validation:

```go
customerIDRequired := true
referenceRequired := true
validationAdapter := "http"
validationURL := "https://api.example.com/validate"
validationTimeout := 10

settings, err := client.UpdateDirectPaySettings(&ghion.UpdateDirectPaySettingsRequest{
    CustomerIDRequired:   &customerIDRequired,
    ReferenceRequired:    &referenceRequired,
    ValidationAdapter:    &validationAdapter,
    ValidationURL:        &validationURL,
    ValidationTimeout:    &validationTimeout,
    ValidationMethod:     func() *string { s := "POST"; return &s }(),
    ValidationAPIKey:     func() *string { s := "your-api-key"; return &s }(),
    ValidationAuthHeader: func() *string { s := "Authorization"; return &s }(),
})
if err != nil {
    // Handle error
}
```

**Validation Configuration Options:**
- `validation_adapter`: "http" for HTTP validation, or "none" to disable
- `validation_url`: Your validation endpoint URL (required when adapter is "http")
- `validation_method`: HTTP method to use (POST, GET)
- `validation_api_key`: API key to send to your validation endpoint
- `validation_auth_header`: Header name for the API key (e.g., "Authorization")
- `validation_timeout`: Request timeout in seconds (1-60)
- `validation_strict`: If true, reject payments when validation fails

#### Testing Direct Pay Settings

Test your validation configuration before going live:

```go
testResult, err := client.TestDirectPaySettings(&ghion.TestDirectPaySettingsRequest{
    CustomerID: "C-123",
    Reference:  "INV-001",
})
if err != nil {
    // Handle error
}

fmt.Printf("Validation Performed: %v\n", testResult.ValidationPerformed)
if testResult.CustomerName != nil {
    fmt.Printf("Customer Name: %s\n", *testResult.CustomerName)
}
if testResult.ReferenceValid != nil {
    fmt.Printf("Reference Valid: %v\n", *testResult.ReferenceValid)
}
if testResult.Amount != nil {
    fmt.Printf("Amount: %.2f\n", *testResult.Amount)
}
if testResult.Error != nil {
    fmt.Printf("Error: %s\n", *testResult.Error)
}
```

#### Customer Validation API

When `validation_adapter` is set to "http", the SDK will call your validation endpoint with the following request:

**Request Body:**
```json
{
  "customer_id": "C-123",
  "reference": "INV-001"
}
```

**Expected Response:**
```json
{
  "customer_name": "John Doe",
  "reference_valid": true,
  "reference": "INV-001",
  "amount": 100.00,
  "error": null
}
```

**Response Path Configuration:**
- `validation_customer_name_path`: JSON path to customer name in response (default: "customer_name")
- `validation_reference_valid_path`: JSON path to reference validity (default: "reference_valid")
- `validation_reference_path`: JSON path to reference (default: "reference")
- `validation_amount_path`: JSON path to amount (default: "amount")
- `validation_error_path`: JSON path to error message (default: "error")

**Important Notes:**
- Validation timeout must be between 1 and 60 seconds
- If `validation_strict` is true, payments will be rejected when validation fails
- If `validation_strict` is false, payments proceed even if validation fails (validation result is informational only)
- The validation endpoint should return HTTP 200 with valid JSON

#### Webhook Events for Direct Pay

The following webhook events are sent for Direct Pay operations:
- `direct_pay.completed`: When a direct payment is completed
- `direct_pay.failed`: When a direct payment fails

### 7. Webhooks Integration
Webhooks are **mandatory** for robust payment verification. Users might close the browser while a payment is processing, so your server must rely on webhooks to fulfill orders.

**Important:** Webhooks must parse the *raw* request body to verify the cryptographic signature.

```go
import (
    "io"
    "net/http"
    "github.com/yayawallet/ghion-go-sdk"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    signature := r.Header.Get("X-Ghion-Signature")
    
    // Read raw body for signature verification
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusInternalServerError)
        return
    }

    event, err := client.ParseWebhook(body, signature)
    if err != nil {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    switch event.Event {
    case ghion.EventTransactionCompleted:
        // Fulfill the order! (e.g., mark DB as paid, send email)
        log.Printf("Payment Success: %s", event.Data.PaymentID)
    case ghion.EventTransactionFailed:
        // Handle failure
        log.Printf("Payment Failed: %s", event.Data.PaymentID)
    // Handle ghion.EventTransactionExpired, etc...
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"received": true}`))
}
```

### 8. Error Handling

The SDK provides detailed error information through custom error types. Always handle errors appropriately to provide good user experience and debugging information.

```go
import "github.com/yayawallet/ghion-go-sdk/pkg/errors"

payment, err := client.InitializePayment(&types.InitializePaymentRequest{
    Amount:    100,
    Reference: "order-123",
})

if err != nil {
    // Check error type for specific handling
    switch e := err.(type) {
    case *errors.ValidationError:
        // Input validation failed (invalid phone number, missing required fields, etc.)
        log.Printf("Validation Error: %s (Field: %s)", e.Message, e.Field)
        // Show user-friendly error message
        return fmt.Errorf("Invalid input: %s", e.Message)
    case *errors.NetworkError:
        // Network connectivity issues
        log.Printf("Network Error: %s", e.Message)
        // Implement retry logic or show connection error
        return fmt.Errorf("Connection error. Please try again.")
    case *errors.APIError:
        // API returned an error (invalid credentials, insufficient funds, etc.)
        log.Printf("API Error: %s (Status: %d)", e.Message, e.StatusCode)
        // Handle specific API errors
        if e.StatusCode == 401 {
            return fmt.Errorf("Authentication failed. Check your API credentials.")
        }
        return fmt.Errorf("Payment error: %s", e.Message)
    case *errors.RateLimitError:
        // Too many requests
        log.Printf("Rate Limit Error: %s", e.Message)
        // Implement exponential backoff
        return fmt.Errorf("Too many requests. Please wait and try again.")
    default:
        // Unknown error
        log.Printf("Unknown Error: %v", err)
        return fmt.Errorf("An unexpected error occurred.")
    }
}
```

**Common Error Scenarios:**

1. **ValidationError**: Invalid input parameters (e.g., invalid phone number format, missing required fields, invalid email format)
2. **NetworkError**: Network connectivity issues, timeout, or DNS resolution failures
3. **APIError**: API returned an error (401 for invalid credentials, 400 for bad request, 500 for server errors)
4. **RateLimitError**: Too many requests within a short time period (implement exponential backoff)

**Best Practices for Error Handling:**
- Log errors with context for debugging
- Show user-friendly error messages to end users
- Implement retry logic for transient errors (network issues, rate limits)
- Validate user input before making API calls to prevent validation errors
- Never expose sensitive information (API secrets, internal details) in error messages

---

## Best Practices

1. **Rely on Webhooks, Not Polling:** 
   Always use Webhooks (`/webhook`) or Server-Sent Events (SSE) driven by webhooks to update the frontend. Avoid setting up frequent polling loops to call `client.GetPaymentStatus()`, as excessive polling can trigger rate limits or interfere with active transactions.
2. **Raw Body for Webhooks:** 
   Always pass the raw byte array of the request body to `client.ParseWebhook()`. If your framework parses the body into a struct before the SDK verifies it, the HMAC signature verification will fail.
3. **Idempotency:** 
   Webhook events can theoretically be delivered more than once. Ensure your database fulfillment logic checks if an order is already marked as "paid" before granting the user access to the product again.
4. **Environment Separation:** 
   Keep a strict separation between Test/Sandbox keys and Production keys.

---

## Common Issues & Fixes

### 1. Error: "Transaction is not in a state awaiting OTP validation"
**Symptom:** You call `SendOTP`, wait for the user to input the code, call `ValidateOTP`, and receive this error.
**Cause:** Calling `client.GetPaymentStatus(paymentID)` via a polling interval *while* the OTP is pending. Querying the Ghion API manually during an active OTP session can forcefully reset the transaction state on the gateway's end.
**Fix:** Remove any background status polling during the OTP flow. Send the OTP, wait for user input, and immediately validate the OTP. Rely on Webhooks for background state changes.

### 2. Webhook Signature Verification Fails
**Symptom:** `client.ParseWebhook()` returns an "Invalid webhook signature" error.
**Cause:** The web framework is parsing the request body before it reaches your handler, modifying the raw byte stream. The signature verification requires the exact raw byte string sent by Ghion.
**Fix:** Ensure you read `r.Body` directly using `io.ReadAll()` before any JSON decoders attempt to parse the request.

### 3. "OTP code is required" error during ValidateOTP
**Symptom:** The SDK returns a validation error or the API rejects the request stating the OTP code is missing, even when provided.
**Cause:** In older versions of the SDK, there was a payload key mismatch (`otp` instead of `otp_code`).
**Fix:** Ensure you are using the latest version of the SDK (`go get -u github.com/yayawallet/ghion-go-sdk`). The SDK internally maps `otpCode` to the correct `otp_code` payload expected by the API.

### 4. Missing Environment Variables
**Symptom:** The SDK fails to initialize with validation errors.
**Cause:** The required `APIKey`, `APISecret`, or `Passphrase` are missing or empty.
**Fix:** Double check that your environment variables are correctly loaded and passed to `ghion.NewClient()`. Use tools like `godotenv` to manage `.env` files in Go.
