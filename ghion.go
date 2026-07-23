// Package ghion is the Go SDK for Ghion Finances payment gateway
//
// This SDK provides a simple and secure way to integrate Ghion Finances payment gateway
// into your Go applications. It supports multiple payment channels including USSD, QR, OTP,
// and provides comprehensive webhook handling.
//
// Quick Start:
//
//	client, err := ghion.NewClient(&ghion.Config{
//	    APIKey:     "your-api-key",
//	    APISecret:  "your-api-secret",
//	    Passphrase: "your-passphrase",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	payment, err := client.InitializePayment(&ghion.InitializePaymentRequest{
//	    Amount:    100,
//	    Reference: "order_12345",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Payment initialized: %s\n", payment.ID)
package ghion

import (
	"github.com/yayawallet/ghion-go-sdk/pkg/client"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

// Config represents the configuration for the Ghion client
type Config = types.GhionConfig

// Client is the main SDK client for Ghion Finances payment gateway
type Client = client.GhionClient

// PaymentStatus represents the status of a payment
type PaymentStatus = types.PaymentStatus

// WebhookEventType represents the type of webhook event
type WebhookEventType = types.WebhookEventType

// Request and Response types
type (
	// InitializePaymentRequest represents a payment initialization request
	InitializePaymentRequest = types.InitializePaymentRequest
	// InitializePaymentResponse represents a payment initialization response
	InitializePaymentResponse = types.InitializePaymentResponse
	// SubmitPaymentRequest represents a payment submission request
	SubmitPaymentRequest = types.SubmitPaymentRequest
	// SubmitPaymentResponse represents a payment submission response
	SubmitPaymentResponse = types.SubmitPaymentResponse
	// PaymentStatusResponse represents a payment status response
	PaymentStatusResponse = types.PaymentStatusResponse
	// CheckoutResponse represents checkout information
	CheckoutResponse = types.CheckoutResponse
	// QRPaymentResponse represents a QR payment response
	QRPaymentResponse = types.QRPaymentResponse
	// OTPSendResponse represents an OTP send response
	OTPSendResponse = types.OTPSendResponse
	// OTPValidateResponse represents an OTP validation response
	OTPValidateResponse = types.OTPValidateResponse
	// PaymentChannel represents an available payment channel
	PaymentChannel = types.PaymentChannel
	// WebhookEvent represents a webhook event payload
	WebhookEvent = types.WebhookEvent
	// Customer represents customer information
	Customer = types.Customer
	// Provider represents payment provider information
	Provider = types.Provider
	// QRInfo represents QR code information
	QRInfo = types.QRInfo
	// Merchant represents merchant information
	Merchant = types.Merchant
)

const (
	// PaymentStatusPending represents a pending payment
	PaymentStatusPending = types.PaymentStatusPending
	// PaymentStatusProcessing represents a processing payment
	PaymentStatusProcessing = types.PaymentStatusProcessing
	// PaymentStatusCompleted represents a completed payment
	PaymentStatusCompleted = types.PaymentStatusCompleted
	// PaymentStatusFailed represents a failed payment
	PaymentStatusFailed = types.PaymentStatusFailed
	// PaymentStatusCancelled represents a cancelled payment
	PaymentStatusCancelled = types.PaymentStatusCancelled
	// PaymentStatusExpired represents an expired payment
	PaymentStatusExpired = types.PaymentStatusExpired
)

// Shorter aliases for payment status (status.something format)
const (
	// StatusPending is a shorter alias for PaymentStatusPending
	StatusPending = types.PaymentStatusPending
	// StatusProcessing is a shorter alias for PaymentStatusProcessing
	StatusProcessing = types.PaymentStatusProcessing
	// StatusCompleted is a shorter alias for PaymentStatusCompleted
	StatusCompleted = types.PaymentStatusCompleted
	// StatusFailed is a shorter alias for PaymentStatusFailed
	StatusFailed = types.PaymentStatusFailed
	// StatusCancelled is a shorter alias for PaymentStatusCancelled
	StatusCancelled = types.PaymentStatusCancelled
	// StatusExpired is a shorter alias for PaymentStatusExpired
	StatusExpired = types.PaymentStatusExpired
)

const (
	// EventTransactionCompleted represents a transaction completed event
	EventTransactionCompleted = types.EventTransactionCompleted
	// EventTransactionFailed represents a transaction failed event
	EventTransactionFailed = types.EventTransactionFailed
	// EventTransactionRefunded represents a transaction refunded event
	EventTransactionRefunded = types.EventTransactionRefunded
	// EventTransactionPartiallyRefunded represents a transaction partially refunded event
	EventTransactionPartiallyRefunded = types.EventTransactionPartiallyRefunded
	// EventTransactionExpired represents a transaction expired event
	EventTransactionExpired = types.EventTransactionExpired
	// EventTransactionDisputed represents a transaction disputed event
	EventTransactionDisputed = types.EventTransactionDisputed
	// EventTransactionUpdated represents a transaction updated event
	EventTransactionUpdated = types.EventTransactionUpdated
)

// NewClient creates a new Ghion client with the given configuration
func NewClient(config *Config) (*Client, error) {
	return client.NewGhionClient(config)
}

