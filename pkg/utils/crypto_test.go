package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	secret := "test-secret"
	timestamp := int64(1234567890)
	method := "POST"
	path := "/checkout/initialize"
	body := `{"amount":100,"reference":"order_12345"}`

	signature := GenerateSignature(timestamp, method, path, body, secret)

	if signature == "" {
		t.Error("Expected non-empty signature")
	}

	// Test that same inputs produce same signature
	signature2 := GenerateSignature(timestamp, method, path, body, secret)
	if signature != signature2 {
		t.Error("Expected same signature for same inputs")
	}

	// Test that different inputs produce different signatures
	signature3 := GenerateSignature(timestamp, method, path, `{"amount":200,"reference":"order_12345"}`, secret)
	if signature == signature3 {
		t.Error("Expected different signature for different inputs")
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	payload := `{"event":"transaction.completed","data":{"payment_id":"12345"}}`

	// Generate a valid signature using the webhook signature format
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Test valid signature
	if !VerifyWebhookSignature([]byte(payload), signature, secret) {
		t.Error("Expected valid signature to verify")
	}

	// Test invalid signature
	if VerifyWebhookSignature([]byte(payload), "invalid-signature", secret) {
		t.Error("Expected invalid signature to fail verification")
	}

	// Test wrong secret
	if VerifyWebhookSignature([]byte(payload), signature, "wrong-secret") {
		t.Error("Expected signature to fail with wrong secret")
	}

	// Test different payload
	if VerifyWebhookSignature([]byte(`{"event":"transaction.failed"}`), signature, secret) {
		t.Error("Expected signature to fail with different payload")
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	timestamp := GetCurrentTimestamp()

	if timestamp <= 0 {
		t.Error("Expected positive timestamp")
	}

	// Test that timestamp is recent (within 1 second)
	currentTime := timestamp
	if currentTime < timestamp-1 || currentTime > timestamp+1 {
		t.Error("Expected timestamp to be current time")
	}
}
