//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func main() {
	// Initialize client with environment variables
	apiKey := os.Getenv("GHION_API_KEY")
	apiSecret := os.Getenv("GHION_API_SECRET")
	passphrase := os.Getenv("GHION_API_PASSPHRASE")

	if apiKey == "" || apiSecret == "" || passphrase == "" {
		log.Fatal("Please set GHION_API_KEY, GHION_API_SECRET, and GHION_API_PASSPHRASE environment variables")
	}

	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		Passphrase: passphrase,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Test phone number from environment or use default
	testPhone := os.Getenv("TEST_PHONE_NUMBER")
	if testPhone == "" {
		testPhone = "+251991599696"
	}

	log.Println("=== Ghion Go SDK - Comprehensive Payment Lifecycle Test ===")
	log.Printf("Test Phone Number: %s\n", testPhone)
	log.Println("")

	// Test 1: Initialize Payment
	log.Println("TEST 1: Initialize Payment")
	log.Println("----------------------------------------")
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   fmt.Sprintf("lifecycle_test_%d", time.Now().Unix()),
		Description: "Comprehensive lifecycle test",
		WebhookURL:  "https://your-domain.com/webhook",
	})
	if err != nil {
		log.Fatalf("Failed to initialize payment: %v", err)
	}

	log.Printf("✓ Payment initialized successfully")
	log.Printf("  Payment ID: %s", payment.ID)
	log.Printf("  Amount: %.2f %s", payment.Amount, payment.Currency)
	log.Printf("  Reference: %s", payment.Reference)
	log.Printf("  Status: %s", payment.Status)
	log.Printf("  Created At: %s", payment.CreatedAt)
	log.Printf("  Expires At: %s", payment.ExpiresAt)
	log.Printf("  Available Channels: %d", len(payment.Channels))
	for i, channel := range payment.Channels {
		log.Printf("    [%d] %s (Code: %s, Type: %s)", i+1, channel.Name, channel.Code, channel.Type)
	}
	log.Println("")

	// Test 2: Get Checkout Details
	log.Println("TEST 2: Get Checkout Details")
	log.Println("----------------------------------------")
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		log.Fatalf("Failed to get checkout: %v", err)
	}

	log.Printf("✓ Checkout details retrieved")
	log.Printf("  Payment ID: %s", checkout.ID)
	log.Printf("  Status: %s", checkout.Status)
	log.Printf("  Is Expired: %v", checkout.IsExpired)
	log.Printf("  Expires At: %s", checkout.ExpiresAt)

	if checkout.Merchant != nil {
		log.Printf("  Merchant: %s", checkout.Merchant.Name)
		log.Printf("  Business: %s", checkout.Merchant.BusinessName)
	}

	log.Printf("  Available Channels: %d", len(checkout.AvailableChannels))
	for i, channel := range checkout.AvailableChannels {
		log.Printf("    [%d] %s (Code: %s)", i+1, channel.Name, channel.Code)
		if channel.SupportsOTP {
			log.Printf("       - Supports OTP: Yes")
		}
	}

	if checkout.QR != nil {
		log.Printf("  QR Available: Yes")
		log.Printf("    QR Image URL: %s", checkout.QR.QRImageURL)
		log.Printf("    QR Payload: %s", checkout.QR.QRPayload)
	}
	log.Println("")

	// Test 3: Get Payment Status
	log.Println("TEST 3: Get Payment Status")
	log.Println("----------------------------------------")
	status, err := client.GetPaymentStatus(payment.ID)
	if err != nil {
		log.Fatalf("Failed to get payment status: %v", err)
	}

	log.Printf("✓ Payment status retrieved")
	log.Printf("  Payment ID: %s", status.ID)
	log.Printf("  Status: %s", status.Status)
	log.Printf("  Amount: %.2f %s", status.Amount, status.Currency)
	log.Printf("  Reference: %s", status.Reference)
	log.Printf("  Created At: %s", status.CreatedAt)
	log.Printf("  Updated At: %s", status.UpdatedAt)

	if status.Channel != "" {
		log.Printf("  Channel: %s", status.Channel)
	}
	if status.TransactionID != "" {
		log.Printf("  Transaction ID: %s", status.TransactionID)
	}
	if status.Customer != nil {
		log.Printf("  Customer:")
		if status.Customer.PhoneNumber != "" {
			log.Printf("    Phone: %s", status.Customer.PhoneNumber)
		}
		if status.Customer.Name != "" {
			log.Printf("    Name: %s", status.Customer.Name)
		}
	}
	log.Println("")

	// Test 4: QR Payment
	log.Println("TEST 4: QR Payment")
	log.Println("----------------------------------------")
	qrResult, err := client.PayWithQR(payment.ID)
	if err != nil {
		log.Printf("✗ QR payment failed: %v", err)
	} else {
		log.Printf("✓ QR payment generated")
		log.Printf("  Type: %s", qrResult.Type)
		log.Printf("  Transaction ID: %s", qrResult.TransactionID)
		log.Printf("  Status: %s", qrResult.Status)
		log.Printf("  QR Image URL: %s", qrResult.QRImageURL)
		log.Printf("  QR Payload: %s", qrResult.QRPayload)
	}
	log.Println("")

	// Test 5: OTP Flow (YaYa Wallet)
	log.Println("TEST 5: OTP Flow (YaYa Wallet)")
	log.Println("----------------------------------------")

	// Find YaYa Wallet channel
	var yayawalletChannel string
	for _, channel := range checkout.AvailableChannels {
		if channel.Code == "yayawallet" || channel.Name == "YaYa Wallet" {
			yayawalletChannel = channel.Code
			break
		}
	}

	if yayawalletChannel != "" {
		log.Printf("Found YaYa Wallet channel: %s", yayawalletChannel)

		// Send OTP
		log.Printf("Sending OTP to %s...", testPhone)
		otpSent, err := client.SendOTP(payment.ID, testPhone)
		if err != nil {
			log.Printf("✗ OTP send failed: %v", err)
		} else {
			log.Printf("✓ OTP sent successfully")
			log.Printf("  Status: %s", otpSent.Status)
			log.Printf("  Transaction ID: %s", otpSent.TransactionID)
			log.Printf("  Message: %s", otpSent.Message)

			// Note: We cannot validate OTP without the actual code
			// In a real scenario, the user would receive the code and enter it
			log.Printf("  Note: OTP validation requires the actual code sent to the phone")
		}
	} else {
		log.Printf("YaYa Wallet channel not found in available channels")
	}
	log.Println("")

	// Test 6: USSD Payment (YaYa Wallet)
	log.Println("TEST 6: USSD Payment (YaYa Wallet)")
	log.Println("----------------------------------------")
	if yayawalletChannel != "" {
		log.Printf("Submitting USSD payment to %s...", testPhone)
		ussdResult, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
			Channel:     yayawalletChannel,
			PhoneNumber: testPhone,
		})
		if err != nil {
			log.Printf("✗ USSD payment failed: %v", err)
		} else {
			log.Printf("✓ USSD payment submitted")
			log.Printf("  Payment ID: %s", ussdResult.ID)
			log.Printf("  Status: %s", ussdResult.Status)
			log.Printf("  Transaction ID: %s", ussdResult.TransactionID)
			log.Printf("  Message: %s", ussdResult.Message)
			if ussdResult.RedirectURL != "" {
				log.Printf("  Redirect URL: %s", ussdResult.RedirectURL)
			}
		}
	} else {
		log.Printf("YaYa Wallet channel not found, skipping USSD test")
	}
	log.Println("")

	// Test 7: Telebirr USSD Payment
	log.Println("TEST 7: Telebirr USSD Payment")
	log.Println("----------------------------------------")

	// Find Telebirr channel
	var telebirrChannel string
	for _, channel := range checkout.AvailableChannels {
		if channel.Code == "telebirr" || channel.Name == "Telebirr" {
			telebirrChannel = channel.Code
			break
		}
	}

	if telebirrChannel != "" {
		log.Printf("Found Telebirr channel: %s", telebirrChannel)
		log.Printf("Submitting USSD payment to %s...", testPhone)
		telebirrResult, err := client.SubmitPayment(payment.ID, &types.SubmitPaymentRequest{
			Channel:     telebirrChannel,
			PhoneNumber: testPhone,
		})
		if err != nil {
			log.Printf("✗ Telebirr payment failed: %v", err)
		} else {
			log.Printf("✓ Telebirr payment submitted")
			log.Printf("  Payment ID: %s", telebirrResult.ID)
			log.Printf("  Status: %s", telebirrResult.Status)
			log.Printf("  Transaction ID: %s", telebirrResult.TransactionID)
			log.Printf("  Message: %s", telebirrResult.Message)
			if telebirrResult.RedirectURL != "" {
				log.Printf("  Redirect URL: %s", telebirrResult.RedirectURL)
			}
		}
	} else {
		log.Printf("Telebirr channel not found in available channels")
	}
	log.Println("")

	// Test 8: Final Status Check
	log.Println("TEST 8: Final Status Check")
	log.Println("----------------------------------------")
	time.Sleep(2 * time.Second) // Wait a moment for any processing
	finalStatus, err := client.GetPaymentStatus(payment.ID)
	if err != nil {
		log.Printf("✗ Final status check failed: %v", err)
	} else {
		log.Printf("✓ Final status retrieved")
		log.Printf("  Payment ID: %s", finalStatus.ID)
		log.Printf("  Final Status: %s", finalStatus.Status)
		log.Printf("  Updated At: %s", finalStatus.UpdatedAt)

		if finalStatus.CompletedAt != "" {
			log.Printf("  Completed At: %s", finalStatus.CompletedAt)
		}
		if finalStatus.FailedAt != "" {
			log.Printf("  Failed At: %s", finalStatus.FailedAt)
			log.Printf("  Failure Reason: %s", finalStatus.FailureReason)
		}
	}
	log.Println("")

	// Summary
	log.Println("=== Test Summary ===")
	log.Printf("Payment ID: %s", payment.ID)
	log.Printf("Initial Status: %s", payment.Status)
	log.Printf("Final Status: %s", finalStatus.Status)
	log.Printf("Total Channels Available: %d", len(checkout.AvailableChannels))

	// Print all channels with their capabilities
	log.Println("\nAvailable Channels:")
	for i, channel := range checkout.AvailableChannels {
		log.Printf("  [%d] %s (%s)", i+1, channel.Name, channel.Code)
		log.Printf("      Type: %s", channel.Type)
		log.Printf("      Supports OTP: %v", channel.SupportsOTP)
	}

	// Print full checkout JSON for debugging
	log.Println("\nFull Checkout Response:")
	checkoutJSON, _ := json.MarshalIndent(checkout, "", "  ")
	log.Printf("%s\n", string(checkoutJSON))

	log.Println("=== Test Complete ===")
}
