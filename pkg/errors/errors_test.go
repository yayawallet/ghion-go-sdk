package errors

import (
	"testing"
)

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name       string
		errorType  interface{}
		errorCode  string
		errorMsg   string
	}{
		{
			name:      "ValidationError",
			errorType: &ValidationError{},
			errorCode: "VALIDATION_ERROR",
			errorMsg:  "Validation failed",
		},
		{
			name:      "ConfigurationError",
			errorType: &ConfigurationError{},
			errorCode: "CONFIGURATION_ERROR",
			errorMsg:  "Configuration error",
		},
		{
			name:      "AuthenticationError",
			errorType: &AuthenticationError{},
			errorCode: "AUTHENTICATION_ERROR",
			errorMsg:  "Authentication failed",
		},
		{
			name:      "APIError",
			errorType: &APIError{},
			errorCode: "API_ERROR",
			errorMsg:  "API error",
		},
		{
			name:      "NetworkError",
			errorType: &NetworkError{},
			errorCode: "NETWORK_ERROR",
			errorMsg:  "Network error",
		},
		{
			name:      "PaymentError",
			errorType: &PaymentError{},
			errorCode: "PAYMENT_ERROR",
			errorMsg:  "Payment error",
		},
		{
			name:      "BillError",
			errorType: &BillError{},
			errorCode: "BILL_ERROR",
			errorMsg:  "Bill error",
		},
		{
			name:      "WebhookError",
			errorType: &WebhookError{},
			errorCode: "WEBHOOK_ERROR",
			errorMsg:  "Webhook error",
		},
		{
			name:      "RateLimitError",
			errorType: &RateLimitError{},
			errorCode: "RATE_LIMIT_ERROR",
			errorMsg:  "Rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			switch tt.errorType.(type) {
			case *ValidationError:
				err = NewValidationError(tt.errorMsg, "test_field", "test_value")
			case *ConfigurationError:
				err = NewConfigurationError(tt.errorMsg, nil)
			case *AuthenticationError:
				err = NewAuthenticationError(tt.errorMsg, nil)
			case *APIError:
				err = NewAPIError(tt.errorMsg, 400, nil)
			case *NetworkError:
				err = NewNetworkError(tt.errorMsg, nil)
			case *PaymentError:
				err = NewPaymentError(tt.errorMsg, "payment_12345", nil)
			case *BillError:
				err = NewBillError(tt.errorMsg, "bill_12345", nil)
			case *WebhookError:
				err = NewWebhookError(tt.errorMsg, nil)
			case *RateLimitError:
				err = NewRateLimitError(tt.errorMsg, 60)
			}

			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			// Check error code via the embedded GhionError
			var ghionErr *GhionError
			switch e := err.(type) {
			case *ValidationError:
				ghionErr = e.GhionError
			case *ConfigurationError:
				ghionErr = e.GhionError
			case *AuthenticationError:
				ghionErr = e.GhionError
			case *APIError:
				ghionErr = e.GhionError
			case *NetworkError:
				ghionErr = e.GhionError
			case *PaymentError:
				ghionErr = e.GhionError
			case *BillError:
				ghionErr = e.GhionError
			case *WebhookError:
				ghionErr = e.GhionError
			case *RateLimitError:
				ghionErr = e.GhionError
			}

			if ghionErr == nil {
				t.Errorf("Expected GhionError but got %T", err)
				return
			}

			if ghionErr.Code != tt.errorCode {
				t.Errorf("Expected error code %s but got %s", tt.errorCode, ghionErr.Code)
			}

			if ghionErr.Message != tt.errorMsg {
				t.Errorf("Expected error message %s but got %s", tt.errorMsg, ghionErr.Message)
			}

			// Test Error() method
			errorStr := err.Error()
			if errorStr == "" {
				t.Error("Expected non-empty error string")
			}
		})
	}
}

func TestAPIErrorWithStatusCode(t *testing.T) {
	err := NewAPIError("API error", 404, nil)
	
	apiErr, ok := err.GhionError.Details["status_code"]
	if !ok {
		t.Errorf("Expected status_code in details")
		return
	}

	if apiErr != 404 {
		t.Errorf("Expected status code 404 but got %v", apiErr)
	}
}

func TestRateLimitErrorWithRetryAfter(t *testing.T) {
	err := NewRateLimitError("Rate limit exceeded", 60)

	retryAfter, ok := err.GhionError.Details["retry_after"]
	if !ok {
		t.Errorf("Expected retry_after in details")
		return
	}

	if retryAfter != 60 {
		t.Errorf("Expected retry after 60 but got %v", retryAfter)
	}
}

func TestBillErrorWithBillID(t *testing.T) {
	err := NewBillError("Bill not found", "bill_12345", nil)

	billID, ok := err.GhionError.Details["bill_id"]
	if !ok {
		t.Errorf("Expected bill_id in details")
		return
	}

	if billID != "bill_12345" {
		t.Errorf("Expected bill ID bill_12345 but got %v", billID)
	}

	if err.BillID != "bill_12345" {
		t.Errorf("Expected BillID field to be bill_12345 but got %s", err.BillID)
	}
}
