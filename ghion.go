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
	"github.com/yayawallet/ghion-go-sdk/pkg/errors"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

// Config represents the configuration for the Ghion client
type Config = types.GhionConfig

// Client is the main SDK client for Ghion Finances payment gateway
type Client = client.GhionClient

// Error types
type (
	// GhionError is the base error class for all SDK errors
	GhionError = errors.GhionError
	// ConfigurationError represents invalid SDK configuration
	ConfigurationError = errors.ConfigurationError
	// AuthenticationError represents invalid credentials or signature
	AuthenticationError = errors.AuthenticationError
	// APIError represents a failed HTTP request
	APIError = errors.APIError
	// ValidationError represents invalid input parameters
	ValidationError = errors.ValidationError
	// NetworkError represents connection or timeout issues
	NetworkError = errors.NetworkError
	// PaymentError represents payment processing failures
	PaymentError = errors.PaymentError
	// BillError represents bill processing failures
	BillError = errors.BillError
	// WebhookError represents webhook signature verification or processing failures
	WebhookError = errors.WebhookError
	// RateLimitError represents API rate limit exceeded
	RateLimitError = errors.RateLimitError
)

// PaymentStatus represents the status of a payment
type PaymentStatus = types.PaymentStatus

// WebhookEventType represents the type of webhook event
type WebhookEventType = types.WebhookEventType

// BillStatus represents the status of a bill
type BillStatus = types.BillStatus

// EscrowStatus represents the status of an escrow
type EscrowStatus = types.EscrowStatus

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
	// Escrow represents a hold payment (escrow)
	Escrow = types.Escrow
	// ListEscrowsRequest represents a request to list escrows
	ListEscrowsRequest = types.ListEscrowsRequest
	// ListEscrowsResponse represents a list of escrows
	ListEscrowsResponse = types.ListEscrowsResponse
	// PullEscrowFundsResponse represents the response when pulling escrow funds
	PullEscrowFundsResponse = types.PullEscrowFundsResponse
	// DirectPaySettings represents the Direct Pay settings
	DirectPaySettings = types.DirectPaySettings
	// GetDirectPaySettingsResponse represents the response when getting Direct Pay settings
	GetDirectPaySettingsResponse = types.GetDirectPaySettingsResponse
	// UpdateDirectPaySettingsRequest represents a request to update Direct Pay settings
	UpdateDirectPaySettingsRequest = types.UpdateDirectPaySettingsRequest
	// TestDirectPaySettingsRequest represents a request to test Direct Pay settings
	TestDirectPaySettingsRequest = types.TestDirectPaySettingsRequest
	// TestDirectPaySettingsResponse represents the response when testing Direct Pay settings
	TestDirectPaySettingsResponse = types.TestDirectPaySettingsResponse
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

const (
	// BillStatusPending represents a pending bill
	BillStatusPending = types.BillStatusPending
	// BillStatusPaid represents a paid bill
	BillStatusPaid = types.BillStatusPaid
	// BillStatusOverdue represents an overdue bill
	BillStatusOverdue = types.BillStatusOverdue
	// BillStatusCancelled represents a cancelled bill
	BillStatusCancelled = types.BillStatusCancelled
)

const (
	// EscrowStatusFunded represents a funded escrow
	EscrowStatusFunded = types.EscrowStatusFunded
	// EscrowStatusWithdrawing represents a withdrawing escrow
	EscrowStatusWithdrawing = types.EscrowStatusWithdrawing
	// EscrowStatusWithdrawn represents a withdrawn escrow
	EscrowStatusWithdrawn = types.EscrowStatusWithdrawn
	// EscrowStatusReleased represents a released escrow
	EscrowStatusReleased = types.EscrowStatusReleased
	// EscrowStatusCancelled represents a cancelled escrow
	EscrowStatusCancelled = types.EscrowStatusCancelled
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

// Error constructors
var (
	// NewConfigurationError creates a new ConfigurationError
	NewConfigurationError = errors.NewConfigurationError
	// NewAuthenticationError creates a new AuthenticationError
	NewAuthenticationError = errors.NewAuthenticationError
	// NewAPIError creates a new APIError
	NewAPIError = errors.NewAPIError
	// NewValidationError creates a new ValidationError
	NewValidationError = errors.NewValidationError
	// NewNetworkError creates a new NetworkError
	NewNetworkError = errors.NewNetworkError
	// NewPaymentError creates a new PaymentError
	NewPaymentError = errors.NewPaymentError
	// NewBillError creates a new BillError
	NewBillError = errors.NewBillError
	// NewWebhookError creates a new WebhookError
	NewWebhookError = errors.NewWebhookError
	// NewRateLimitError creates a new RateLimitError
	NewRateLimitError = errors.NewRateLimitError
)
