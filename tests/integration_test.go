// +build integration

package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

// Integration tests require real API credentials
// Run with: go test ./... -tags=integration

func init() {
	// Load .env file from project root (parent directory)
	if err := godotenv.Load("../.env"); err != nil {
		// Don't fail if .env file doesn't exist
		// Environment variables might be set manually
		fmt.Printf("Warning: Failed to load .env file: %v\n", err)
	} else {
		fmt.Println("Successfully loaded .env file")
	}
}

func getTestClient(t *testing.T) *ghion.Client {
	apiKey := os.Getenv("GHION_API_KEY")
	apiSecret := os.Getenv("GHION_API_SECRET")
	passphrase := os.Getenv("GHION_API_PASSPHRASE")

	if apiKey == "" || apiSecret == "" || passphrase == "" {
		t.Skip("Skipping integration tests: GHION_API_KEY, GHION_API_SECRET, and GHION_API_PASSPHRASE must be set")
	}

	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		Passphrase: passphrase,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func TestIntegration_InitializePayment(t *testing.T) {
	client := getTestClient(t)

	reference := fmt.Sprintf("test_order_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "Integration test payment",
		WebhookURL:  "https://example.com/webhook",
	})

	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	if payment.ID == "" {
		t.Error("Payment ID should not be empty")
	}

	if payment.Amount != 100 {
		t.Errorf("Expected amount 100, got %f", payment.Amount)
	}

	if payment.Reference != reference {
		t.Errorf("Expected reference %s, got %s", reference, payment.Reference)
	}

	t.Logf("Payment initialized: %s", payment.ID)
	t.Logf("Available channels in response: %d", len(payment.Channels))
	
	// If channels are not in the initial response, get them from checkout
	if len(payment.Channels) == 0 {
		t.Log("No channels in initial response, fetching from checkout...")
		checkout, err := client.GetCheckout(payment.ID)
		if err != nil {
			t.Logf("Warning: Failed to get checkout: %v", err)
		} else {
			t.Logf("Available channels from checkout: %d", len(checkout.AvailableChannels))
			for _, channel := range checkout.AvailableChannels {
				t.Logf("  - %s (%s)", channel.Name, channel.Code)
			}
		}
	} else {
		for _, channel := range payment.Channels {
			t.Logf("  - %s (%s)", channel.Name, channel.Code)
		}
	}
}

func TestIntegration_QRPayment(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("qr_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "QR payment integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized for QR: %s", payment.ID)

	// Step 2: Generate QR code
	qrResult, err := client.PayWithQR(payment.ID)
	if err != nil {
		t.Fatalf("Failed to generate QR: %v", err)
	}

	if qrResult.QRImageURL == "" {
		t.Error("QR Image URL should not be empty")
	}

	if qrResult.QRPayload == "" {
		t.Error("QR Payload should not be empty")
	}

	t.Logf("QR Payment Type: %s", qrResult.Type)
	t.Logf("QR Transaction ID: %s", qrResult.TransactionID)
	t.Logf("QR Status: %s", qrResult.Status)
	t.Logf("QR Image URL: %s", qrResult.QRImageURL)

	// Step 3: Get checkout details
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		t.Fatalf("Failed to get checkout: %v", err)
	}

	if checkout.QR == nil {
		t.Error("Checkout should contain QR information")
	} else {
		t.Logf("Checkout QR Image URL: %s", checkout.QR.QRImageURL)
	}
}

func TestIntegration_OTPPayment(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("otp_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "OTP payment integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized for OTP: %s", payment.ID)

	// Step 2: Send OTP (requires a real phone number)
	phoneNumber := os.Getenv("TEST_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("Skipping OTP send: TEST_PHONE_NUMBER must be set")
	}

	otpSent, err := client.SendOTP(payment.ID, phoneNumber)
	if err != nil {
		t.Logf("Warning: Failed to send OTP (this may be expected): %v", err)
		// Don't fail the test, as OTP sending might fail for various reasons
		return
	}

	t.Logf("OTP sent successfully")
	t.Logf("OTP Status: %s", otpSent.Status)
	t.Logf("OTP Transaction ID: %s", otpSent.TransactionID)
	t.Logf("OTP Message: %s", otpSent.Message)

	// Step 3: Validate OTP (requires the actual OTP code sent to the phone)
	otpCode := os.Getenv("TEST_OTP_CODE")
	if otpCode == "" {
		t.Log("Skipping OTP validation: TEST_OTP_CODE must be set")
		t.Log("To complete OTP validation, you need to:")
		t.Log("1. Check your phone for the OTP code")
		t.Log("2. Set TEST_OTP_CODE environment variable")
		t.Log("3. Run the test again")
		return
	}

	validateResult, err := client.ValidateOTP(payment.ID, otpCode, phoneNumber)
	if err != nil {
		t.Logf("Warning: Failed to validate OTP (this may be expected): %v", err)
		return
	}

	t.Logf("OTP validated successfully")
	t.Logf("Validation Status: %s", validateResult.Status)
	t.Logf("Validation Transaction ID: %s", validateResult.TransactionID)
}

func TestIntegration_CreateBill(t *testing.T) {
	client := getTestClient(t)

	billID := fmt.Sprintf("INV-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		CustomerPhone: "+251911234567",
		CustomerEmail: "test@example.com",
		Description:  "Integration test bill",
		BillCode:     "UTIL",
		Cluster:      "ADDIS ABABA",
	})

	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	if bill.ID == "" {
		t.Error("Bill ID should not be empty")
	}

	if bill.BillID != billID {
		t.Errorf("Expected bill ID %s, got %s", billID, bill.BillID)
	}

	if bill.Amount != 500 {
		t.Errorf("Expected amount 500, got %f", bill.Amount)
	}

	t.Logf("Bill created successfully: %s", bill.ID)
	t.Logf("Bill ID: %s", bill.BillID)
	t.Logf("Status: %s", bill.Status)
}

func TestIntegration_CreateBulkBills(t *testing.T) {
	client := getTestClient(t)

	bulkRequest := &types.BulkCreateBillsRequest{
		Bills: []types.CreateBillRequest{
			{
				BillID:       fmt.Sprintf("BULK-1-%d", time.Now().Unix()),
				Amount:       300,
				Currency:     "ETB",
				DueDate:      "2026-09-01",
				CustomerName: "Customer 1",
				Description:  "Bulk bill 1",
			},
			{
				BillID:       fmt.Sprintf("BULK-2-%d", time.Now().Unix()),
				Amount:       400,
				Currency:     "ETB",
				DueDate:      "2026-09-01",
				CustomerName: "Customer 2",
				Description:  "Bulk bill 2",
			},
		},
	}

	response, err := client.CreateBulkBills(bulkRequest)
	if err != nil {
		t.Fatalf("Failed to create bulk bills: %v", err)
	}

	if response.CreatedCount == 0 {
		t.Error("Expected at least one bill to be created")
	}

	t.Logf("Bulk bills created: %d", response.CreatedCount)
	t.Logf("Errors: %d", response.ErrorCount)
}

func TestIntegration_ListBills(t *testing.T) {
	client := getTestClient(t)

	response, err := client.ListBills(&types.ListBillsRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("Failed to list bills: %v", err)
	}

	t.Logf("Total bills: %d", response.Total)
	t.Logf("Page: %d", response.Page)
	t.Logf("Items on this page: %d", len(response.Items))
}

func TestIntegration_GetBillStatistics(t *testing.T) {
	client := getTestClient(t)

	stats, err := client.GetBillStatistics()
	if err != nil {
		t.Fatalf("Failed to get bill statistics: %v", err)
	}

	t.Logf("Pending: %d", stats.Pending)
	t.Logf("Paid: %d", stats.Paid)
	t.Logf("Forwarded: %d", stats.Forwarded)
	t.Logf("Overdue: %d", stats.Overdue)
	t.Logf("Total Amount: %.2f", stats.TotalAmount)
	t.Logf("Total Paid: %.2f", stats.TotalPaid)
}

func TestIntegration_GetBillDashboard(t *testing.T) {
	client := getTestClient(t)

	dashboard, err := client.GetBillDashboard("2026-08-01", "2026-08-31")
	if err != nil {
		t.Fatalf("Failed to get bill dashboard: %v", err)
	}

	t.Logf("Dashboard Summary - Total Bills: %d", dashboard.Summary.TotalBills)
	t.Logf("Dashboard Summary - Pending: %d", dashboard.Summary.Pending)
	t.Logf("Dashboard Summary - Paid: %d", dashboard.Summary.Paid)
	t.Logf("Dashboard Summary - Total Amount: %.2f", dashboard.Summary.TotalAmount)
	t.Logf("Clusters: %d", len(dashboard.ByCluster))
	t.Logf("Bill Codes: %d", len(dashboard.ByBillCode))
	t.Logf("Trend: %d data points", len(dashboard.Trend))
}

func TestIntegration_GetBillDetail(t *testing.T) {
	client := getTestClient(t)

	// First create a bill
	billID := fmt.Sprintf("INV-DETAIL-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Integration test bill for detail",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Get bill detail
	detail, err := client.GetBillDetail(bill.ID)
	if err != nil {
		t.Fatalf("Failed to get bill detail: %v", err)
	}

	if detail.ID != bill.ID {
		t.Errorf("Expected bill ID %s, got %s", bill.ID, detail.ID)
	}

	t.Logf("Bill detail retrieved: %s", detail.ID)
	t.Logf("Customer Name: %s", detail.CustomerName)
	t.Logf("Amount: %.2f", detail.Amount)
	t.Logf("Payments: %d", len(detail.Payments))
}

func TestIntegration_UpdateBill(t *testing.T) {
	client := getTestClient(t)

	// First create a bill
	billID := fmt.Sprintf("INV-UPDATE-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Original description",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Update bill
	updatedBill, err := client.UpdateBill(bill.ID, &types.UpdateBillRequest{
		Amount:      bill.Amount,
		Description: "Updated description",
		DueDate:     "2026-09-15",
	})
	if err != nil {
		t.Fatalf("Failed to update bill: %v", err)
	}

	if updatedBill.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got %s", updatedBill.Description)
	}

	t.Logf("Bill updated successfully: %s", updatedBill.ID)
}

func TestIntegration_RecordManualPayment(t *testing.T) {
	client := getTestClient(t)

	// First create a bill
	billID := fmt.Sprintf("INV-PAY-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Integration test bill for manual payment",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Record manual payment
	payment, err := client.RecordManualPayment(bill.ID, &types.RecordManualPaymentRequest{
		Amount:        100,
		Source:        "manual",
		PaymentMethod: "cash",
		Reference:     "RECEIPT-001",
		Note:          "Paid at counter",
	})
	if err != nil {
		t.Fatalf("Failed to record manual payment: %v", err)
	}

	t.Logf("Manual payment recorded: %s", payment.PaymentID)
	t.Logf("Amount: %.2f", payment.Amount)
	t.Logf("Bill Status: %s", payment.BillStatus)
	t.Logf("Balance Due: %.2f", payment.BalanceDue)
}

func TestIntegration_DeleteBill(t *testing.T) {
	client := getTestClient(t)

	// First create a bill
	billID := fmt.Sprintf("INV-DELETE-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Integration test bill for deletion",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Delete bill
	deleteResult, err := client.DeleteBill(bill.ID)
	if err != nil {
		t.Fatalf("Failed to delete bill: %v", err)
	}

	t.Logf("Bill deleted successfully: %s", deleteResult.Message)
}

func TestIntegration_GetBillerSettings(t *testing.T) {
	client := getTestClient(t)

	settings, err := client.GetBillerSettings()
	if err != nil {
		t.Fatalf("Failed to get biller settings: %v", err)
	}

	t.Logf("Biller configured: %v", settings.Configured)
	if settings.Settings != nil {
		t.Logf("Biller Code: %s", settings.Settings.BillerCode)
		t.Logf("Biller Name: %s", settings.Settings.BillerName)
		t.Logf("Clusters: %d", len(settings.Settings.Clusters))
		t.Logf("Bill Codes: %d", len(settings.Settings.BillCodes))
	}
}

func TestIntegration_PublicBillLookup(t *testing.T) {
	client := getTestClient(t)

	// First get biller settings to get biller code
	settings, err := client.GetBillerSettings()
	if err != nil {
		t.Fatalf("Failed to get biller settings: %v", err)
	}

	if !settings.Configured || settings.Settings == nil {
		t.Skip("Skipping public bill lookup: Biller not configured")
	}

	// Create a bill for lookup
	billID := fmt.Sprintf("INV-LOOKUP-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Integration test bill for public lookup",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Public bill lookup
	publicBill, err := client.PublicBillLookup(&types.PublicBillLookupRequest{
		BillerCode: settings.Settings.BillerCode,
		BillID:     bill.BillID,
	})
	if err != nil {
		t.Fatalf("Failed to perform public bill lookup: %v", err)
	}

	t.Logf("Public bill lookup successful")
	t.Logf("Bill ID: %s", publicBill.BillID)
	t.Logf("Customer Name: %s", publicBill.CustomerName)
	t.Logf("Amount: %.2f", publicBill.Amount)
	t.Logf("Payment Status: %s", publicBill.PaymentStatus)
	t.Logf("Amount Due: %.2f", publicBill.AmountDue)
}

func TestIntegration_GetBillPaymentLink(t *testing.T) {
	client := getTestClient(t)

	// First create a bill
	billID := fmt.Sprintf("INV-LINK-%d", time.Now().Unix())
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       billID,
		Amount:       500,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Integration test bill for payment link",
	})
	if err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Get payment link
	link, err := client.GetBillPaymentLink(bill.ID)
	if err != nil {
		t.Fatalf("Failed to get bill payment link: %v", err)
	}

	t.Logf("Payment link generated: %s", link.CheckoutURL)
}

func TestIntegration_YaYaWalletPayment(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("yayawallet_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "YaYa Wallet payment integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized for YaYa Wallet: %s", payment.ID)

	// Step 2: Submit YaYa Wallet payment
	phoneNumber := os.Getenv("TEST_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("Skipping YaYa Wallet payment: TEST_PHONE_NUMBER must be set")
	}

	submitResult, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
		Channel:       "yayawallet",
		PaymentMethod: "ussd",
		PhoneNumber:   phoneNumber,
	})
	if err != nil {
		t.Logf("Warning: Failed to submit YaYa Wallet payment (this may be expected): %v", err)
		// Don't fail the test, as payment submission might fail for various reasons
		return
	}

	t.Logf("YaYa Wallet Payment submitted successfully")
	t.Logf("Payment ID: %s", submitResult.ID)
	t.Logf("Status: %s", submitResult.Status)
	t.Logf("Transaction ID: %s", submitResult.TransactionID)
	t.Logf("Message: %s", submitResult.Message)
}

func TestIntegration_TelebirrPayment(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("telebirr_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "Telebirr payment integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized for Telebirr: %s", payment.ID)

	// Step 2: Submit Telebirr payment
	phoneNumber := os.Getenv("TEST_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("Skipping Telebirr payment: TEST_PHONE_NUMBER must be set")
	}

	submitResult, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
		Channel:       "telebirr",
		PaymentMethod: "ussd",
		PhoneNumber:   phoneNumber,
	})
	if err != nil {
		t.Logf("Warning: Failed to submit Telebirr payment (this may be expected): %v", err)
		// Don't fail the test, as payment submission might fail for various reasons
		return
	}

	t.Logf("Telebirr Payment submitted successfully")
	t.Logf("Payment ID: %s", submitResult.ID)
	t.Logf("Status: %s", submitResult.Status)
	t.Logf("Transaction ID: %s", submitResult.TransactionID)
	t.Logf("Message: %s", submitResult.Message)
}

func TestIntegration_USSDPayment(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("ussd_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "USSD payment integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized for USSD: %s", payment.ID)

	// Step 2: Get checkout to find available channels
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		t.Fatalf("Failed to get checkout: %v", err)
	}

	// Step 3: Submit USSD payment
	// Find a USSD channel from available channels
	var ussdChannel string
	for _, channel := range checkout.AvailableChannels {
		if channel.Type == "ussd" || channel.Code == "ussd" {
			ussdChannel = channel.Code
			break
		}
	}

	if ussdChannel == "" {
		// Use first available channel if no USSD channel found
		if len(checkout.AvailableChannels) > 0 {
			ussdChannel = checkout.AvailableChannels[0].Code
			t.Logf("No USSD channel found, using first available: %s", ussdChannel)
		} else {
			t.Fatal("No payment channels available")
		}
	}

	phoneNumber := os.Getenv("TEST_PHONE_NUMBER")
	if phoneNumber == "" {
		t.Skip("Skipping USSD payment: TEST_PHONE_NUMBER must be set")
	}

	submitResult, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
		Channel:     ussdChannel,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		t.Logf("Warning: Failed to submit USSD payment (this may be expected): %v", err)
		// Don't fail the test, as payment submission might fail for various reasons
		return
	}

	t.Logf("USSD Payment submitted successfully")
	t.Logf("Payment ID: %s", submitResult.ID)
	t.Logf("Status: %s", submitResult.Status)
	t.Logf("Transaction ID: %s", submitResult.TransactionID)
	t.Logf("Message: %s", submitResult.Message)
}

func TestIntegration_PaymentStatus(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("status_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "Payment status integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized: %s", payment.ID)

	// Step 2: Get payment status
	status, err := client.GetPaymentStatus(payment.ID)
	if err != nil {
		t.Fatalf("Failed to get payment status: %v", err)
	}

	if status.ID != payment.ID {
		t.Errorf("Expected payment ID %s, got %s", payment.ID, status.ID)
	}

	if status.Amount != 100 {
		t.Errorf("Expected amount 100, got %f", status.Amount)
	}

	t.Logf("Payment Status: %s", status.Status)
	t.Logf("Created At: %s", status.CreatedAt)
	t.Logf("Updated At: %s", status.UpdatedAt)
}

func TestIntegration_GetCheckout(t *testing.T) {
	client := getTestClient(t)

	// Step 1: Initialize payment
	reference := fmt.Sprintf("checkout_test_%d", time.Now().Unix())
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   reference,
		Description: "Checkout integration test",
		WebhookURL:  "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("Failed to initialize payment: %v", err)
	}

	t.Logf("Payment initialized: %s", payment.ID)

	// Step 2: Get checkout details
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		t.Fatalf("Failed to get checkout: %v", err)
	}

	if checkout.ID != payment.ID {
		t.Errorf("Expected payment ID %s, got %s", payment.ID, checkout.ID)
	}

	if checkout.Amount != 100 {
		t.Errorf("Expected amount 100, got %f", checkout.Amount)
	}

	t.Logf("Checkout Status: %s", checkout.Status)
	t.Logf("Is Expired: %v", checkout.IsExpired)
	t.Logf("Expires At: %s", checkout.ExpiresAt)

	if checkout.Merchant != nil {
		t.Logf("Merchant: %s", checkout.Merchant.Name)
	}

	if len(checkout.AvailableChannels) > 0 {
		t.Logf("Available Channels: %d", len(checkout.AvailableChannels))
		for _, channel := range checkout.AvailableChannels {
			t.Logf("  - %s (%s)", channel.Name, channel.Code)
		}
	}
}
