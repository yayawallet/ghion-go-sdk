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
   - [Webhooks Integration](#3-webhooks-integration)
   - [Error Handling](#5-error-handling)
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
    BillID:        "INV-12345",
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
- `billID`: Unique bill identifier
- `amount`: Bill amount
- `dueDate`: Due date in Y-m-d format
- `customerName`: Customer name

**Optional Parameters:**
- `currency`: Currency code (default: ETB)
- `customerPhone`: Customer phone number
- `customerEmail`: Customer email
- `description`: Bill description
- `billCode`: Bill code for categorization
- `cluster`: Geographic cluster

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

### 5. Webhooks Integration
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

### 6. Error Handling

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
