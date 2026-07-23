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
	fmt.Println("Step 1: Initializing payment for QR payment...")
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   fmt.Sprintf("qr_order_%d", time.Now().Unix()),
		Description: "QR payment test",
		WebhookURL:  "https://your-domain.com/webhook",
	})
	if err != nil {
		log.Fatalf("Failed to initialize payment: %v", err)
	}

	fmt.Printf("Payment initialized: %s\n", payment.ID)
	fmt.Printf("Amount: %.2f %s\n", payment.Amount, payment.Currency)

	// Step 2: Get QR code for payment
	fmt.Println("\nStep 2: Generating QR code...")
	qrResult, err := client.PayWithQR(payment.ID)
	if err != nil {
		log.Fatalf("Failed to generate QR: %v", err)
	}

	fmt.Printf("QR payment type: %s\n", qrResult.Type)
	fmt.Printf("Transaction ID: %s\n", qrResult.TransactionID)
	fmt.Printf("Status: %s\n", qrResult.Status)
	fmt.Printf("QR Image URL: %s\n", qrResult.QRImageURL)
	fmt.Printf("QR Payload: %s\n", qrResult.QRPayload)

	// Step 3: Get checkout details (includes QR info)
	fmt.Println("\nStep 3: Getting checkout details...")
	checkout, err := client.GetCheckout(payment.ID)
	if err != nil {
		log.Fatalf("Failed to get checkout details: %v", err)
	}

	fmt.Printf("Checkout status: %s\n", checkout.Status)
	fmt.Printf("Is expired: %v\n", checkout.IsExpired)
	if checkout.QR != nil {
		fmt.Printf("QR Image URL: %s\n", checkout.QR.QRImageURL)
		fmt.Printf("QR Payload: %s\n", checkout.QR.QRPayload)
	}
	if checkout.Merchant != nil {
		fmt.Printf("Merchant: %s\n", checkout.Merchant.Name)
		fmt.Printf("Business: %s\n", checkout.Merchant.BusinessName)
	}

	// Step 4: Display available payment channels
	fmt.Println("\nStep 4: Available payment channels:")
	for _, channel := range checkout.AvailableChannels {
		fmt.Printf("  - %s (Code: %s)\n", channel.Name, channel.Code)
		if channel.SupportsOTP {
			fmt.Printf("    Supports OTP: Yes\n")
		}
		if channel.RequiresPhone {
			fmt.Printf("    Requires Phone: Yes\n")
		}
	}

	fmt.Println("\nQR Payment Flow:")
	fmt.Println("1. Display the QR code to your customer")
	fmt.Println("2. Customer scans the QR code with their mobile banking app")
	fmt.Println("3. Customer confirms the payment on their phone")
	fmt.Println("4. You'll receive a webhook notification when payment is completed")
	fmt.Println("5. Use the webhook to update your system")
}
