package client

import (
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
	"github.com/yayawallet/ghion-go-sdk/pkg/webhook"
)

// WebhookHandler wraps the webhook handler for use in the client
type WebhookHandler struct {
	handler *webhook.WebhookHandler
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(apiSecret string) *WebhookHandler {
	return &WebhookHandler{
		handler: webhook.NewWebhookHandler(apiSecret),
	}
}

// VerifyWebhook verifies the webhook signature
func (wh *WebhookHandler) VerifyWebhook(rawBody []byte, signature string) bool {
	return wh.handler.VerifyWebhook(rawBody, signature)
}

// ParseWebhook parses and verifies a webhook event
func (wh *WebhookHandler) ParseWebhook(rawBody []byte, signature string) (*types.WebhookEvent, error) {
	return wh.handler.ParseWebhook(rawBody, signature)
}

// ValidateWebhookEvent validates that a webhook event has all required fields
func (wh *WebhookHandler) ValidateWebhookEvent(event *types.WebhookEvent) error {
	return wh.handler.ValidateWebhookEvent(event)
}

// IsTransactionCompleted checks if the event is a transaction completed event
func (wh *WebhookHandler) IsTransactionCompleted(event *types.WebhookEvent) bool {
	return wh.handler.IsTransactionCompleted(event)
}

// IsTransactionFailed checks if the event is a transaction failed event
func (wh *WebhookHandler) IsTransactionFailed(event *types.WebhookEvent) bool {
	return wh.handler.IsTransactionFailed(event)
}

// IsTransactionRefunded checks if the event is a transaction refunded event
func (wh *WebhookHandler) IsTransactionRefunded(event *types.WebhookEvent) bool {
	return wh.handler.IsTransactionRefunded(event)
}

// IsTransactionExpired checks if the event is a transaction expired event
func (wh *WebhookHandler) IsTransactionExpired(event *types.WebhookEvent) bool {
	return wh.handler.IsTransactionExpired(event)
}

// GetEventType returns the event type as a string
func (wh *WebhookHandler) GetEventType(event *types.WebhookEvent) string {
	return wh.handler.GetEventType(event)
}

// GetPaymentID returns the payment ID from the webhook event
func (wh *WebhookHandler) GetPaymentID(event *types.WebhookEvent) string {
	return wh.handler.GetPaymentID(event)
}

// GetTransactionID returns the transaction ID from the webhook event
func (wh *WebhookHandler) GetTransactionID(event *types.WebhookEvent) string {
	return wh.handler.GetTransactionID(event)
}

// GetAmount returns the amount from the webhook event
func (wh *WebhookHandler) GetAmount(event *types.WebhookEvent) float64 {
	return wh.handler.GetAmount(event)
}

// GetStatus returns the payment status from the webhook event
func (wh *WebhookHandler) GetStatus(event *types.WebhookEvent) types.PaymentStatus {
	return wh.handler.GetStatus(event)
}
