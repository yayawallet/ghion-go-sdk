package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

// GenerateSignature generates HMAC-SHA256 signature for API authentication
// timestamp: Unix timestamp in seconds
// method: HTTP method (GET, POST, etc.)
// path: API path including /api/v1 prefix
// body: Request body as string
// secret: API secret key
// Returns: Base64-encoded signature
func GenerateSignature(timestamp int64, method, path, body, secret string) string {
	message := strconv.FormatInt(timestamp, 10) + method + path + body
	
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifyWebhookSignature verifies webhook signature
// rawBody: Raw request body as byte slice
// signature: Signature from X-Ghion-Signature header
// secret: API secret key
// Returns: True if signature is valid
func VerifyWebhookSignature(rawBody []byte, signature, secret string) bool {
	expected := generateWebhookSignature(rawBody, secret)
	
	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(expected), []byte(signature))
}

// generateWebhookSignature generates the expected webhook signature
func generateWebhookSignature(rawBody []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(rawBody)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetCurrentTimestamp returns current Unix timestamp in seconds
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
