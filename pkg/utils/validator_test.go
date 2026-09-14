package utils

import (
	"testing"

	"github.com/yayawallet/ghion-go-sdk/pkg/types"
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

func TestValidateBillID(t *testing.T) {
	tests := []struct {
		name        string
		billID      string
		expectError bool
	}{
		{
			name:        "valid bill ID",
			billID:      "INV-12345",
			expectError: false,
		},
		{
			name:        "empty bill ID",
			billID:      "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			billID:      "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBillID(tt.billID)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateCreateBillRequest(t *testing.T) {
	tests := []struct {
		name          string
		billID        string
		amount        float64
		dueDate       string
		customerEmail string
		expectError   bool
	}{
		{
			name:          "valid request",
			billID:        "INV-12345",
			amount:        500,
			dueDate:       "2026-09-01",
			customerEmail: "",
			expectError:   false,
		},
		{
			name:          "valid with email",
			billID:        "INV-12345",
			amount:        500,
			dueDate:       "2026-09-01",
			customerEmail: "test@example.com",
			expectError:   false,
		},
		{
			name:          "zero amount",
			billID:        "INV-12345",
			amount:        0,
			dueDate:       "2026-09-01",
			customerEmail: "",
			expectError:   true,
		},
		{
			name:          "negative amount",
			billID:        "INV-12345",
			amount:        -100,
			dueDate:       "2026-09-01",
			customerEmail: "",
			expectError:   true,
		},
		{
			name:          "empty due date",
			billID:        "INV-12345",
			amount:        500,
			dueDate:       "",
			customerEmail: "",
			expectError:   true,
		},
		{
			name:          "invalid date format",
			billID:        "INV-12345",
			amount:        500,
			dueDate:       "2026/09/01",
			customerEmail: "",
			expectError:   true,
		},
		{
			name:          "invalid email",
			billID:        "INV-12345",
			amount:        500,
			dueDate:       "2026-09-01",
			customerEmail: "invalid-email",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateBillRequest(tt.billID, tt.amount, tt.dueDate, tt.customerEmail)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateBulkCreateBillsRequest(t *testing.T) {
	tests := []struct {
		name        string
		bills       []interface{}
		expectError bool
	}{
		{
			name:        "valid request",
			bills:       []interface{}{"bill1", "bill2"},
			expectError: false,
		},
		{
			name:        "empty bills array",
			bills:       []interface{}{},
			expectError: true,
		},
		{
			name:        "nil bills array",
			bills:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBulkCreateBillsRequest(tt.bills)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateListBillsRequest(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		search      string
		cluster     string
		billCode    string
		from        string
		to          string
		page        int
		limit       int
		expectError bool
	}{
		{
			name:        "valid request",
			status:      "pending",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "",
			to:          "",
			page:        1,
			limit:       10,
			expectError: false,
		},
		{
			name:        "valid with dates",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "2026-08-01",
			to:          "2026-08-31",
			page:        1,
			limit:       10,
			expectError: false,
		},
		{
			name:        "invalid from date",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "2026/08/01",
			to:          "2026-08-31",
			page:        1,
			limit:       10,
			expectError: true,
		},
		{
			name:        "invalid to date",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "2026-08-01",
			to:          "2026/08/31",
			page:        1,
			limit:       10,
			expectError: true,
		},
		{
			name:        "negative page",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "",
			to:          "",
			page:        -1,
			limit:       10,
			expectError: true,
		},
		{
			name:        "limit too high",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "",
			to:          "",
			page:        1,
			limit:       101,
			expectError: true,
		},
		{
			name:        "negative limit",
			status:      "",
			search:      "",
			cluster:     "",
			billCode:    "",
			from:        "",
			to:          "",
			page:        1,
			limit:       -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListBillsRequest(tt.status, tt.search, tt.cluster, tt.billCode, tt.from, tt.to, tt.page, tt.limit)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePublicBillLookupRequest(t *testing.T) {
	tests := []struct {
		name        string
		billerCode  string
		billID      string
		expectError bool
	}{
		{
			name:        "valid request",
			billerCode:  "BILLER001",
			billID:      "INV-12345",
			expectError: false,
		},
		{
			name:        "empty biller code",
			billerCode:  "",
			billID:      "INV-12345",
			expectError: true,
		},
		{
			name:        "empty bill ID",
			billerCode:  "BILLER001",
			billID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePublicBillLookupRequest(tt.billerCode, tt.billID)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateRecordManualPaymentRequest(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		expectError bool
	}{
		{
			name:        "valid amount",
			amount:      100,
			expectError: false,
		},
		{
			name:        "zero amount",
			amount:      0,
			expectError: true,
		},
		{
			name:        "negative amount",
			amount:      -100,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRecordManualPaymentRequest(tt.amount)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateBillerSettingsRequest(t *testing.T) {
	tests := []struct {
		name              string
		serviceChargeRate float64
		clusters          []interface{}
		billCodes         []interface{}
		expectError       bool
	}{
		{
			name:              "valid request",
			serviceChargeRate: 0.05,
			clusters:          []interface{}{"ADDIS ABABA", "HAWASSA"},
			billCodes:         []interface{}{map[string]interface{}{"code": "UTIL", "name": "Utilities"}},
			expectError:       false,
		},
		{
			name:              "negative service charge rate",
			serviceChargeRate: -0.05,
			clusters:          nil,
			billCodes:         nil,
			expectError:       true,
		},
		{
			name:              "invalid cluster type",
			serviceChargeRate: 0.05,
			clusters:          []interface{}{123},
			billCodes:         nil,
			expectError:       true,
		},
		{
			name:              "invalid bill code type",
			serviceChargeRate: 0.05,
			clusters:          nil,
			billCodes:         []interface{}{"invalid"},
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBillerSettingsRequest(tt.serviceChargeRate, tt.clusters, tt.billCodes)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateBillDashboardRequest(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		expectError bool
	}{
		{
			name:        "valid request",
			from:        "2026-08-01",
			to:          "2026-08-31",
			expectError: false,
		},
		{
			name:        "empty dates",
			from:        "",
			to:          "",
			expectError: false,
		},
		{
			name:        "invalid from date",
			from:        "2026/08/01",
			to:          "2026-08-31",
			expectError: true,
		},
		{
			name:        "invalid to date",
			from:        "2026-08-01",
			to:          "2026/08/31",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBillDashboardRequest(tt.from, tt.to)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected bool
	}{
		{
			name:     "valid date",
			dateStr:  "2026-08-07",
			expected: true,
		},
		{
			name:     "invalid format with slashes",
			dateStr:  "2026/08/07",
			expected: false,
		},
		{
			name:     "invalid format without dashes",
			dateStr:  "20260807",
			expected: false,
		},
		{
			name:     "empty string",
			dateStr:  "",
			expected: false,
		},
		{
			name:     "invalid month",
			dateStr:  "2026-13-01",
			expected: true, // Pattern matches, even if invalid date
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidDate(tt.dateStr)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "valid email with subdomain",
			email:    "test@mail.example.com",
			expected: true,
		},
		{
			name:     "invalid - no @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "invalid - no domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "invalid - no local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "invalid - spaces",
			email:    "test @example.com",
			expected: false,
		},
		{
			name:     "empty string",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestValidateEscrowID(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name:        "valid escrow ID",
			id:          "escrow-123",
			expectError: false,
		},
		{
			name:        "empty escrow ID",
			id:          "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			id:          "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEscrowID(tt.id)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateListEscrowsRequest(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectError bool
	}{
		{
			name:        "valid status",
			status:      "funded",
			expectError: false,
		},
		{
			name:        "empty status",
			status:      "",
			expectError: false,
		},
		{
			name:        "invalid status",
			status:      "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListEscrowsRequest(tt.status)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateUpdateDirectPaySettingsRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *types.UpdateDirectPaySettingsRequest
		expectError bool
	}{
		{
			name: "valid request",
			request: &types.UpdateDirectPaySettingsRequest{
				ValidationTimeout: func() *int { i := 10; return &i }(),
			},
			expectError: false,
		},
		{
			name: "validation timeout too low",
			request: &types.UpdateDirectPaySettingsRequest{
				ValidationTimeout: func() *int { i := 0; return &i }(),
			},
			expectError: true,
		},
		{
			name: "validation timeout too high",
			request: &types.UpdateDirectPaySettingsRequest{
				ValidationTimeout: func() *int { i := 61; return &i }(),
			},
			expectError: true,
		},
		{
			name: "validation adapter http without URL",
			request: &types.UpdateDirectPaySettingsRequest{
				ValidationAdapter: func() *string { s := "http"; return &s }(),
				ValidationURL:     func() *string { s := ""; return &s }(),
			},
			expectError: true,
		},
		{
			name: "validation adapter http with URL",
			request: &types.UpdateDirectPaySettingsRequest{
				ValidationAdapter: func() *string { s := "http"; return &s }(),
				ValidationURL:     func() *string { s := "https://api.example.com"; return &s }(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateDirectPaySettingsRequest(tt.request)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateTestDirectPaySettingsRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *types.TestDirectPaySettingsRequest
		expectError bool
	}{
		{
			name: "valid request",
			request: &types.TestDirectPaySettingsRequest{
				CustomerID: "C-123",
				Reference:  "INV-001",
			},
			expectError: false,
		},
		{
			name: "empty customer ID",
			request: &types.TestDirectPaySettingsRequest{
				CustomerID: "",
			},
			expectError: true,
		},
		{
			name: "valid without reference",
			request: &types.TestDirectPaySettingsRequest{
				CustomerID: "C-123",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTestDirectPaySettingsRequest(tt.request)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
