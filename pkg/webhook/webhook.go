package webhook

import (
	"encoding/json"

	ghionerrors "github.com/yayawallet/ghion-go-sdk/pkg/errors"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
	"github.com/yayawallet/ghion-go-sdk/pkg/utils"
)

// WebhookHandler handles webhook signature verification and parsing
type WebhookHandler struct {
	apiSecret string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(apiSecret string) *WebhookHandler {
	return &WebhookHandler{
		apiSecret: apiSecret,
	}
}

// VerifyWebhook verifies the webhook signature
// rawBody: Raw request body as byte slice
// signature: Signature from X-Ghion-Signature header
// Returns: True if signature is valid
func (wh *WebhookHandler) VerifyWebhook(rawBody []byte, signature string) bool {
	return utils.VerifyWebhookSignature(rawBody, signature, wh.apiSecret)
}

// ParseWebhook parses and verifies a webhook event
// rawBody: Raw request body as byte slice
// signature: Signature from X-Ghion-Signature header
// Returns: Parsed webhook event or error if signature is invalid
func (wh *WebhookHandler) ParseWebhook(rawBody []byte, signature string) (*types.WebhookEvent, error) {
	if !wh.VerifyWebhook(rawBody, signature) {
		return nil, ghionerrors.NewAuthenticationError("Invalid webhook signature", nil)
	}

	var event types.WebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return nil, ghionerrors.NewWebhookError("Failed to parse webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return &event, nil
}

// ParseWebhookUnsafe parses a webhook event without signature verification
// WARNING: Only use this for testing purposes
// rawBody: Raw request body as byte slice
// Returns: Parsed webhook event or error if parsing fails
func (wh *WebhookHandler) ParseWebhookUnsafe(rawBody []byte) (*types.WebhookEvent, error) {
	var event types.WebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return nil, ghionerrors.NewWebhookError("Failed to parse webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return &event, nil
}

// ValidateWebhookEvent validates that a webhook event has all required fields
func (wh *WebhookHandler) ValidateWebhookEvent(event *types.WebhookEvent) error {
	if event.Event == "" {
		return ghionerrors.NewValidationError("Webhook event type is required", "event", nil)
	}

	if event.Data.PaymentID == "" {
		return ghionerrors.NewValidationError("Webhook event data payment_id is required", "data.payment_id", nil)
	}

	if event.Data.Amount <= 0 {
		return ghionerrors.NewValidationError("Webhook event data amount must be positive", "data.amount", event.Data.Amount)
	}

	if event.Data.Currency == "" {
		return ghionerrors.NewValidationError("Webhook event data currency is required", "data.currency", nil)
	}

	if event.Data.Reference == "" {
		return ghionerrors.NewValidationError("Webhook event data reference is required", "data.reference", nil)
	}

	if event.Data.Timestamp == "" {
		return ghionerrors.NewValidationError("Webhook event data timestamp is required", "data.timestamp", nil)
	}

	return nil
}

// IsTransactionCompleted checks if the event is a transaction completed event
func (wh *WebhookHandler) IsTransactionCompleted(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionCompleted
}

// IsTransactionFailed checks if the event is a transaction failed event
func (wh *WebhookHandler) IsTransactionFailed(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionFailed
}

// IsTransactionRefunded checks if the event is a transaction refunded event
func (wh *WebhookHandler) IsTransactionRefunded(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionRefunded
}

// IsTransactionExpired checks if the event is a transaction expired event
func (wh *WebhookHandler) IsTransactionExpired(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionExpired
}

// IsTransactionPartiallyRefunded checks if the event is a transaction partially refunded event
func (wh *WebhookHandler) IsTransactionPartiallyRefunded(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionPartiallyRefunded
}

// IsTransactionDisputed checks if the event is a transaction disputed event
func (wh *WebhookHandler) IsTransactionDisputed(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionDisputed
}

// IsTransactionUpdated checks if the event is a transaction updated event
func (wh *WebhookHandler) IsTransactionUpdated(event *types.WebhookEvent) bool {
	return event.Event == types.EventTransactionUpdated
}

// GetEventType returns the event type as a string
func (wh *WebhookHandler) GetEventType(event *types.WebhookEvent) string {
	return string(event.Event)
}

// GetPaymentID returns the payment ID from the webhook event
func (wh *WebhookHandler) GetPaymentID(event *types.WebhookEvent) string {
	return event.Data.PaymentID
}

// GetTransactionID returns the transaction ID from the webhook event
func (wh *WebhookHandler) GetTransactionID(event *types.WebhookEvent) string {
	return event.Data.TransactionID
}

// GetAmount returns the amount from the webhook event
func (wh *WebhookHandler) GetAmount(event *types.WebhookEvent) float64 {
	return event.Data.Amount
}

// GetStatus returns the payment status from the webhook event
func (wh *WebhookHandler) GetStatus(event *types.WebhookEvent) types.PaymentStatus {
	return event.Data.Status
}

// GetCurrency returns the currency from the webhook event
func (wh *WebhookHandler) GetCurrency(event *types.WebhookEvent) string {
	return event.Data.Currency
}

// GetReference returns the reference from the webhook event
func (wh *WebhookHandler) GetReference(event *types.WebhookEvent) string {
	return event.Data.Reference
}

// GetTimestamp returns the timestamp from the webhook event
func (wh *WebhookHandler) GetTimestamp(event *types.WebhookEvent) string {
	return event.Data.Timestamp
}
