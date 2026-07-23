package utils

import (
	"testing"
)

func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		expectError bool
	}{
		{
			name:        "valid API key",
			apiKey:      "gw_test_1234567890abcdef",
			expectError: false,
		},
		{
			name:        "empty API key",
			apiKey:      "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			apiKey:      "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKey(tt.apiKey)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateAPISecret(t *testing.T) {
	tests := []struct {
		name        string
		apiSecret   string
		expectError bool
	}{
		{
			name:        "valid API secret",
			apiSecret:   "test-secret-key-12345",
			expectError: false,
		},
		{
			name:        "empty API secret",
			apiSecret:   "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			apiSecret:   "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPISecret(tt.apiSecret)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePassphrase(t *testing.T) {
	tests := []struct {
		name        string
		passphrase  string
		expectError bool
	}{
		{
			name:        "valid passphrase",
			passphrase:  "test-passphrase",
			expectError: false,
		},
		{
			name:        "empty passphrase",
			passphrase:  "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			passphrase:  "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassphrase(tt.passphrase)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateInitializePaymentRequest(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		reference   string
		currency    string
		webhookURL  string
		returnURL   string
		cancelURL   string
		expectError bool
	}{
		{
			name:        "valid request",
			amount:      100,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "",
			returnURL:   "",
			cancelURL:   "",
			expectError: false,
		},
		{
			name:        "valid request with URLs",
			amount:      100,
			reference:   "order_12345",
			currency:    "ETB",
			webhookURL:  "https://example.com/webhook",
			returnURL:   "https://example.com/success",
			cancelURL:   "",
			expectError: false,
		},
		{
			name:        "zero amount",
			amount:      0,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "",
			returnURL:   "",
			cancelURL:   "",
			expectError: true,
		},
		{
			name:        "negative amount",
			amount:      -100,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "",
			returnURL:   "",
			cancelURL:   "",
			expectError: true,
		},
		{
			name:        "empty reference",
			amount:      100,
			reference:   "",
			currency:    "",
			webhookURL:  "",
			returnURL:   "",
			cancelURL:   "",
			expectError: true,
		},
		{
			name:        "invalid webhook URL",
			amount:      100,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "not-a-valid-url",
			returnURL:   "",
			cancelURL:   "",
			expectError: true,
		},
		{
			name:        "invalid return URL",
			amount:      100,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "",
			returnURL:   "not-a-valid-url",
			cancelURL:   "",
			expectError: true,
		},
		{
			name:        "invalid cancel URL",
			amount:      100,
			reference:   "order_12345",
			currency:    "",
			webhookURL:  "",
			returnURL:   "",
			cancelURL:   "not-a-valid-url",
			expectError: true,
		},
		{
			name:        "valid with all URLs",
			amount:      100,
			reference:   "order_12345",
			currency:    "ETB",
			webhookURL:  "https://example.com/webhook",
			returnURL:   "https://example.com/success",
			cancelURL:   "https://example.com/cancel",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInitializePaymentRequest(tt.amount, tt.reference, tt.currency, tt.webhookURL, tt.returnURL, tt.cancelURL)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateSubmitPaymentRequest(t *testing.T) {
	tests := []struct {
		name        string
		channel     string
		phoneNumber string
		expectError bool
	}{
		{
			name:        "valid request",
			channel:     "telebirr",
			phoneNumber: "",
			expectError: false,
		},
		{
			name:        "valid request with phone",
			channel:     "telebirr",
			phoneNumber: "+251911234567",
			expectError: false,
		},
		{
			name:        "empty channel",
			channel:     "",
			phoneNumber: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSubmitPaymentRequest(tt.channel, tt.phoneNumber, "")
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePaymentID(t *testing.T) {
	tests := []struct {
		name        string
		paymentID   string
		expectError bool
	}{
		{
			name:        "valid payment ID",
			paymentID:   "payment_12345",
			expectError: false,
		},
		{
			name:        "empty payment ID",
			paymentID:   "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			paymentID:   "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentID(tt.paymentID)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expectError bool
	}{
		{
			name:        "valid phone number",
			phoneNumber: "+251911234567",
			expectError: false,
		},
		{
			name:        "empty phone number",
			phoneNumber: "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			phoneNumber: "   ",
			expectError: true,
		},
		{
			name:        "valid phone with spaces",
			phoneNumber: "+251 911 234 567",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.phoneNumber)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateOTPCode(t *testing.T) {
	tests := []struct {
		name        string
		otpCode     string
		expectError bool
	}{
		{
			name:        "valid OTP code",
			otpCode:     "123456",
			expectError: false,
		},
		{
			name:        "empty OTP code",
			otpCode:     "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			otpCode:     "   ",
			expectError: true,
		},
		{
			name:        "OTP code with letters",
			otpCode:     "12345a",
			expectError: false,
		},
		{
			name:        "short OTP code",
			otpCode:     "123",
			expectError: false,
		},
		{
			name:        "long OTP code",
			otpCode:     "1234567890",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOTPCode(tt.otpCode)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateEthiopianPhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expectError bool
	}{
		{
			name:        "valid with country code",
			phoneNumber: "+251911234567",
			expectError: false,
		},
		{
			name:        "valid with leading zero",
			phoneNumber: "0911234567",
			expectError: false,
		},
		{
			name:        "valid without prefix",
			phoneNumber: "911234567",
			expectError: false,
		},
		{
			name:        "invalid - too short",
			phoneNumber: "91123456",
			expectError: true,
		},
		{
			name:        "invalid - too long",
			phoneNumber: "091123456789",
			expectError: true,
		},
		{
			name:        "invalid - wrong prefix",
			phoneNumber: "0811234567",
			expectError: true,
		},
		{
			name:        "invalid - letters",
			phoneNumber: "09123456ab",
			expectError: true,
		},
		{
			name:        "empty",
			phoneNumber: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEthiopianPhoneNumber(tt.phoneNumber)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
