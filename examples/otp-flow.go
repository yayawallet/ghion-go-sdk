//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func main() {
	// Initialize client
	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     os.Getenv("GHION_API_KEY"),
		APISecret:  os.Getenv("GHION_API_SECRET"),
		Passphrase: os.Getenv("GHION_API_PASSPHRASE"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Step 1: Initialize payment
	fmt.Println("Step 1: Initializing payment...")
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   fmt.Sprintf("otp_order_%d", time.Now().Unix()),
		Description: "OTP payment test",
		WebhookURL:  "https://your-domain.com/webhook",
	})
	if err != nil {
		log.Fatalf("Failed to initialize payment: %v", err)
	}

	fmt.Printf("Payment initialized: %s\n", payment.ID)
	fmt.Printf("Amount: %.2f %s\n", payment.Amount, payment.Currency)

	// Step 2: Send OTP to customer's phone
	phoneNumber := "+251911234567" // Replace with actual phone number
	fmt.Printf("\nStep 2: Sending OTP to %s...\n", phoneNumber)

	otpSent, err := client.SendOTP(payment.ID, phoneNumber)
	if err != nil {
		log.Fatalf("Failed to send OTP: %v", err)
	}

	fmt.Printf("OTP sent successfully: %s\n", otpSent.Status)
	fmt.Printf("Transaction ID: %s\n", otpSent.TransactionID)

	// Step 3: Prompt for OTP code
	fmt.Println("\nStep 3: Please enter the OTP code sent to your phone:")
	var otpCode string
	fmt.Scanln(&otpCode)

	// Step 4: Validate OTP
	fmt.Printf("\nStep 4: Validating OTP...\n")
	otpValidated, err := client.ValidateOTP(payment.ID, otpCode, phoneNumber)
	if err != nil {
		log.Fatalf("Failed to validate OTP: %v", err)
	}

	fmt.Printf("OTP validation result: %s\n", otpValidated.Status)
	fmt.Printf("Transaction ID: %s\n", otpValidated.TransactionID)

	// Step 5: Get final payment status
	fmt.Println("\nStep 5: Getting final payment status...")
	status, err := client.GetPaymentStatus(payment.ID)
	if err != nil {
		log.Printf("Warning: Failed to get payment status: %v", err)
	} else {
		fmt.Printf("Final payment status: %s\n", status.Status)
		if status.Status == ghion.PaymentStatusCompleted {
			fmt.Printf("Payment completed successfully!\n")
			fmt.Printf("Transaction ID: %s\n", status.TransactionID)
		}
	}

	// Step 6: Get full checkout details
	fmt.Println("\nStep 6: Getting checkout details...")
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		log.Printf("Warning: Failed to get checkout details: %v", err)
	} else {
		fmt.Printf("Checkout status: %s\n", checkout.Status)
		fmt.Printf("Is expired: %v\n", checkout.IsExpired)
		if checkout.Merchant != nil {
			fmt.Printf("Merchant: %s\n", checkout.Merchant.Name)
		}
	}
}
