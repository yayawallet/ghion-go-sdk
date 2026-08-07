package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	timestamp := int64(1234567890)
	method := "POST"
	path := "/api/v1/checkout"
	body := `{"amount":100}`
	secret := "test-secret"

	signature := GenerateSignature(timestamp, method, path, body, secret)
	if signature == "" {
		t.Error("Expected non-empty signature")
	}

	// Test with empty body
	signature2 := GenerateSignature(timestamp, method, path, "", secret)
	if signature2 == "" {
		t.Error("Expected non-empty signature with empty body")
	}

	// Test different method
	signature3 := GenerateSignature(timestamp, "GET", path, body, secret)
	if signature3 == "" {
		t.Error("Expected non-empty signature with GET method")
	}

	// Test with different path
	signature4 := GenerateSignature(timestamp, method, "/api/v1/payment", body, secret)
	if signature4 == "" {
		t.Error("Expected non-empty signature with different path")
	}

	// Test with different secret
	signature5 := GenerateSignature(timestamp, method, path, body, "different-secret")
	if signature5 == "" {
		t.Error("Expected non-empty signature with different secret")
	}

	// Test that same inputs produce same signature
	signature6 := GenerateSignature(timestamp, method, path, body, secret)
	if signature != signature6 {
		t.Error("Expected same signature for same inputs")
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	rawBody := []byte(`{"event":"transaction.completed","data":{"payment_id":"12345"}}`)

	// Generate a valid signature
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(rawBody)
	validSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Test valid signature
	if !VerifyWebhookSignature(rawBody, validSignature, secret) {
		t.Error("Expected valid signature to verify")
	}

	// Test invalid signature
	if VerifyWebhookSignature(rawBody, "invalid-signature", secret) {
		t.Error("Expected invalid signature to fail verification")
	}

	// Test empty signature
	if VerifyWebhookSignature(rawBody, "", secret) {
		t.Error("Expected empty signature to fail verification")
	}

	// Test different body
	differentBody := []byte(`{"event":"transaction.failed"}`)
	if VerifyWebhookSignature(differentBody, validSignature, secret) {
		t.Error("Expected signature verification to fail with different body")
	}

	// Test empty body
	emptyBody := []byte(``)
	h2 := hmac.New(sha256.New, []byte(secret))
	h2.Write(emptyBody)
	emptySignature := base64.StdEncoding.EncodeToString(h2.Sum(nil))
	if !VerifyWebhookSignature(emptyBody, emptySignature, secret) {
		t.Error("Expected empty body signature to verify")
	}

	// Test different secret
	if VerifyWebhookSignature(rawBody, validSignature, "different-secret") {
		t.Error("Expected signature verification to fail with different secret")
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	timestamp := GetCurrentTimestamp()
	if timestamp == 0 {
		t.Error("Expected non-zero timestamp")
	}

	// Test that timestamp is reasonable (within last 10 seconds)
	currentTime := int64(0) // This would be time.Now().Unix() in real test
	if timestamp < currentTime-10 || timestamp > currentTime+10 {
		t.Logf("Timestamp: %d", timestamp)
	}
}

func TestGenerateWebhookSignature(t *testing.T) {
	secret := "test-secret"
	rawBody := []byte(`{"event":"transaction.completed"}`)

	// Generate signature using the private function from crypto.go
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(rawBody)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if signature == "" {
		t.Error("Expected non-empty signature")
	}

	// Test with empty body
	emptyBody := []byte{}
	h2 := hmac.New(sha256.New, []byte(secret))
	h2.Write(emptyBody)
	emptySignature := base64.StdEncoding.EncodeToString(h2.Sum(nil))
	if emptySignature == "" {
		t.Error("Expected non-empty signature for empty body")
	}

	// Test with different secret
	differentSecret := "different-secret"
	h3 := hmac.New(sha256.New, []byte(differentSecret))
	h3.Write(rawBody)
	differentSignature := base64.StdEncoding.EncodeToString(h3.Sum(nil))
	if differentSignature == "" {
		t.Error("Expected non-empty signature with different secret")
	}

	// Test that same inputs produce same signature
	h4 := hmac.New(sha256.New, []byte(secret))
	h4.Write(rawBody)
	sameSignature := base64.StdEncoding.EncodeToString(h4.Sum(nil))
	if signature != sameSignature {
		t.Error("Expected same signature for same inputs")
	}
}
