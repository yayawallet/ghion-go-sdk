package errors

import "fmt"

// GhionError is the base error class for all SDK errors
type GhionError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *GhionError) Error() string {
	if e.Details != nil && len(e.Details) > 0 {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewGhionError creates a new base GhionError
func NewGhionError(message, code string, details map[string]interface{}) *GhionError {
	return &GhionError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ConfigurationError represents invalid SDK configuration
type ConfigurationError struct {
	*GhionError
}

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(message string, details map[string]interface{}) *ConfigurationError {
	return &ConfigurationError{
		GhionError: NewGhionError(message, "CONFIGURATION_ERROR", details),
	}
}

// AuthenticationError represents invalid credentials or signature
type AuthenticationError struct {
	*GhionError
}

// NewAuthenticationError creates a new AuthenticationError
func NewAuthenticationError(message string, details map[string]interface{}) *AuthenticationError {
	return &AuthenticationError{
		GhionError: NewGhionError(message, "AUTHENTICATION_ERROR", details),
	}
}

// APIError represents a failed HTTP request
type APIError struct {
	*GhionError
	StatusCode int
	Response   map[string]interface{}
}

// NewAPIError creates a new APIError
func NewAPIError(message string, statusCode int, response map[string]interface{}) *APIError {
	details := map[string]interface{}{
		"status_code": statusCode,
	}
	if response != nil {
		details["response"] = response
	}
	return &APIError{
		GhionError: NewGhionError(message, "API_ERROR", details),
		StatusCode: statusCode,
		Response:   response,
	}
}

// ValidationError represents invalid input parameters
type ValidationError struct {
	*GhionError
	Field string
	Value interface{}
}

// NewValidationError creates a new ValidationError
func NewValidationError(message, field string, value interface{}) *ValidationError {
	details := map[string]interface{}{
		"field": field,
		"value": value,
	}
	return &ValidationError{
		GhionError: NewGhionError(message, "VALIDATION_ERROR", details),
		Field:      field,
		Value:      value,
	}
}

// NetworkError represents connection or timeout issues
type NetworkError struct {
	*GhionError
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(message string, details map[string]interface{}) *NetworkError {
	return &NetworkError{
		GhionError: NewGhionError(message, "NETWORK_ERROR", details),
	}
}

// PaymentError represents payment processing failures
type PaymentError struct {
	*GhionError
	PaymentID string
}

// NewPaymentError creates a new PaymentError
func NewPaymentError(message, paymentID string, details map[string]interface{}) *PaymentError {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["payment_id"] = paymentID
	return &PaymentError{
		GhionError: NewGhionError(message, "PAYMENT_ERROR", details),
		PaymentID:  paymentID,
	}
}

// BillError represents bill processing failures
type BillError struct {
	*GhionError
	BillID string
}

// NewBillError creates a new BillError
func NewBillError(message, billID string, details map[string]interface{}) *BillError {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["bill_id"] = billID
	return &BillError{
		GhionError: NewGhionError(message, "BILL_ERROR", details),
		BillID:     billID,
	}
}

// WebhookError represents webhook signature verification or processing failures
type WebhookError struct {
	*GhionError
}

// NewWebhookError creates a new WebhookError
func NewWebhookError(message string, details map[string]interface{}) *WebhookError {
	return &WebhookError{
		GhionError: NewGhionError(message, "WEBHOOK_ERROR", details),
	}
}

// RateLimitError represents API rate limit exceeded
type RateLimitError struct {
	*GhionError
	RetryAfter int
}

// NewRateLimitError creates a new RateLimitError
func NewRateLimitError(message string, retryAfter int) *RateLimitError {
	details := map[string]interface{}{
		"retry_after": retryAfter,
	}
	return &RateLimitError{
		GhionError: NewGhionError(message, "RATE_LIMIT_ERROR", details),
		RetryAfter: retryAfter,
	}
}
