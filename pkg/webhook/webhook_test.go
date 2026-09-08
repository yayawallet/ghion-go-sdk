package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func TestWebhookHandler(t *testing.T) {
	handler := NewWebhookHandler("test-secret")

	t.Run("verify valid signature", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)

		// Generate a valid signature
		h := hmac.New(sha256.New, []byte("test-secret"))
		h.Write(payload)
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		result := handler.VerifyWebhook(payload, signature)
		if !result {
			t.Error("Expected true for valid signature")
		}
	})

	t.Run("verify invalid signature", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)
		signature := "invalid-signature"

		result := handler.VerifyWebhook(payload, signature)
		if result {
			t.Error("Expected false for invalid signature")
		}
	})

	t.Run("verify empty signature", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)
		signature := ""

		result := handler.VerifyWebhook(payload, signature)
		if result {
			t.Error("Expected false for empty signature")
		}
	})

	t.Run("parse valid webhook", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345","amount":100,"currency":"ETB","reference":"order_12345","status":"completed","timestamp":"2024-01-01T00:00:00Z"}}`)

		// Generate a valid signature
		h := hmac.New(sha256.New, []byte("test-secret"))
		h.Write(payload)
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		event, err := handler.ParseWebhook(payload, signature)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if event == nil {
			t.Error("Expected event but got nil")
		}
		if event.Data.PaymentID != "12345" {
			t.Errorf("Expected payment ID 12345 but got %s", event.Data.PaymentID)
		}
	})

	t.Run("parse webhook with invalid signature", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)
		signature := "invalid-signature"

		_, err := handler.ParseWebhook(payload, signature)
		if err == nil {
			t.Error("Expected error for invalid signature")
		}
	})

	t.Run("parse webhook with empty signature", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)
		signature := ""

		_, err := handler.ParseWebhook(payload, signature)
		if err == nil {
			t.Error("Expected error for empty signature")
		}
	})

	t.Run("parse webhook with invalid JSON", func(t *testing.T) {
		payload := []byte(`invalid json`)
		signature := "signature"

		_, err := handler.ParseWebhook(payload, signature)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

func TestWebhookEventHelpers(t *testing.T) {
	t.Run("IsTransactionCompleted", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionCompleted(event) {
			t.Error("Expected true for transaction completed event")
		}
	})

	t.Run("IsTransactionFailed", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionFailed,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionFailed(event) {
			t.Error("Expected true for transaction failed event")
		}
	})

	t.Run("IsTransactionRefunded", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionRefunded,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionRefunded(event) {
			t.Error("Expected true for transaction refunded event")
		}
	})

	t.Run("IsTransactionExpired", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionExpired,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionExpired(event) {
			t.Error("Expected true for transaction expired event")
		}
	})

	t.Run("IsTransactionPartiallyRefunded", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionPartiallyRefunded,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionPartiallyRefunded(event) {
			t.Error("Expected true for transaction partially refunded event")
		}
	})

	t.Run("IsTransactionDisputed", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionDisputed,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionDisputed(event) {
			t.Error("Expected true for transaction disputed event")
		}
	})

	t.Run("IsTransactionUpdated", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionUpdated,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		if !handler.IsTransactionUpdated(event) {
			t.Error("Expected true for transaction updated event")
		}
	})

	t.Run("GetEventType", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		eventType := handler.GetEventType(event)
		if string(eventType) != string(types.EventTransactionCompleted) {
			t.Errorf("Expected %s but got %s", types.EventTransactionCompleted, eventType)
		}
	})

	t.Run("GetPaymentID", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
			},
		}
		paymentID := handler.GetPaymentID(event)
		if paymentID != "12345" {
			t.Errorf("Expected 12345 but got %s", paymentID)
		}
	})

	t.Run("GetTransactionID", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID:     "12345",
				TransactionID: "txn_12345",
			},
		}
		transactionID := handler.GetTransactionID(event)
		if transactionID != "txn_12345" {
			t.Errorf("Expected txn_12345 but got %s", transactionID)
		}
	})

	t.Run("GetAmount", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
			},
		}
		amount := handler.GetAmount(event)
		if amount != 100 {
			t.Errorf("Expected 100 but got %f", amount)
		}
	})

	t.Run("GetStatus", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Status:    types.PaymentStatusCompleted,
			},
		}
		status := handler.GetStatus(event)
		if status != types.PaymentStatusCompleted {
			t.Errorf("Expected %s but got %s", types.PaymentStatusCompleted, status)
		}
	})

	t.Run("GetCurrency", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Currency:  "ETB",
			},
		}
		currency := handler.GetCurrency(event)
		if currency != "ETB" {
			t.Errorf("Expected ETB but got %s", currency)
		}
	})

	t.Run("GetReference", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Reference: "order_12345",
			},
		}
		reference := handler.GetReference(event)
		if reference != "order_12345" {
			t.Errorf("Expected order_12345 but got %s", reference)
		}
	})

	t.Run("GetTimestamp", func(t *testing.T) {
		handler := NewWebhookHandler("test-secret")
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		timestamp := handler.GetTimestamp(event)
		if timestamp != "2024-01-01T00:00:00Z" {
			t.Errorf("Expected 2024-01-01T00:00:00Z but got %s", timestamp)
		}
	})
}

func TestValidateWebhookEvent(t *testing.T) {
	handler := NewWebhookHandler("test-secret")

	t.Run("valid event", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("missing event type", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: "",
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for missing event type")
		}
	})

	t.Run("missing payment ID", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "",
				Amount:    100,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for missing payment ID")
		}
	})

	t.Run("zero amount", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    0,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for zero amount")
		}
	})

	t.Run("negative amount", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    -100,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for negative amount")
		}
	})

	t.Run("missing currency", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
				Currency:  "",
				Reference: "order_12345",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for missing currency")
		}
	})

	t.Run("missing reference", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
				Currency:  "ETB",
				Reference: "",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for missing reference")
		}
	})

	t.Run("missing timestamp", func(t *testing.T) {
		event := &types.WebhookEvent{
			Event: types.EventTransactionCompleted,
			Data: types.WebhookEventData{
				PaymentID: "12345",
				Amount:    100,
				Currency:  "ETB",
				Reference: "order_12345",
				Timestamp: "",
			},
		}
		err := handler.ValidateWebhookEvent(event)
		if err == nil {
			t.Error("Expected error for missing timestamp")
		}
	})
}

func TestParseWebhookUnsafe(t *testing.T) {
	handler := NewWebhookHandler("test-secret")

	t.Run("parse valid webhook unsafe", func(t *testing.T) {
		payload := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345","amount":100,"currency":"ETB","reference":"order_12345","status":"completed","timestamp":"2024-01-01T00:00:00Z"}}`)

		event, err := handler.ParseWebhookUnsafe(payload)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if event == nil {
			t.Error("Expected event but got nil")
		}
		if event.Data.PaymentID != "12345" {
			t.Errorf("Expected payment ID 12345 but got %s", event.Data.PaymentID)
		}
	})

	t.Run("parse webhook with invalid JSON unsafe", func(t *testing.T) {
		payload := []byte(`invalid json`)

		_, err := handler.ParseWebhookUnsafe(payload)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}
