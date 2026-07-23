//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func main() {
	// Initialize client with configuration
	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     os.Getenv("GHION_API_KEY"),
		APISecret:  os.Getenv("GHION_API_SECRET"),
		Passphrase: os.Getenv("GHION_API_PASSPHRASE"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Initialize a payment
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      100,
		Currency:    "ETB",
		Reference:   fmt.Sprintf("order_%d", os.Getpid()),
		Description: "Test payment",
		WebhookURL:  "https://your-domain.com/webhook",
	})
	if err != nil {
		log.Fatalf("Failed to initialize payment: %v", err)
	}

	fmt.Printf("Payment initialized successfully!\n")
	fmt.Printf("Payment ID: %s\n", payment.ID)
	fmt.Printf("Amount: %.2f %s\n", payment.Amount, payment.Currency)
	fmt.Printf("Status: %s\n", payment.Status)
	fmt.Printf("Available channels:\n")
	for _, channel := range payment.Channels {
		fmt.Printf("  - %s (ID: %s)\n", channel.Name, channel.Code)
	}

	// Get payment status
	status, err := client.GetPaymentStatus(payment.ID)
	if err != nil {
		log.Printf("Warning: Failed to get payment status: %v", err)
	} else {
		fmt.Printf("\nPayment status: %s\n", status.Status)
	}
}
