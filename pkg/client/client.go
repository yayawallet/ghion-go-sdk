package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/yayawallet/ghion-go-sdk/pkg/errors"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
	"github.com/yayawallet/ghion-go-sdk/pkg/utils"
)

const (
	defaultBaseURL         = "https://ghion.financial/api/v1"
	defaultCheckoutBaseURL = "https://app.ghion.financial/api/v1"
	defaultTimeout         = 30 * time.Second
	maxRetries             = 3
)

// GhionClient is the main SDK client for Ghion Finances payment gateway
type GhionClient struct {
	apiKey          string
	apiSecret       string
	passphrase      string
	baseURL         string
	checkoutBaseURL string
	timeout         time.Duration
	httpClient      *http.Client
	webhookHandler  *WebhookHandler
}

// NewGhionClient creates a new Ghion client with the given configuration
func NewGhionClient(config *types.GhionConfig) (*GhionClient, error) {
	// Validate configuration
	if err := utils.ValidateAPIKey(config.APIKey); err != nil {
		return nil, err
	}
	if err := utils.ValidateAPISecret(config.APISecret); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassphrase(config.Passphrase); err != nil {
		return nil, err
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	checkoutBaseURL := config.CheckoutBaseURL
	if checkoutBaseURL == "" {
		checkoutBaseURL = defaultCheckoutBaseURL
	}

	timeout := defaultTimeout
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Millisecond
	}

	client := &GhionClient{
		apiKey:          config.APIKey,
		apiSecret:       config.APISecret,
		passphrase:      config.Passphrase,
		baseURL:         baseURL,
		checkoutBaseURL: checkoutBaseURL,
		timeout:         timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		webhookHandler: NewWebhookHandler(config.APISecret),
	}

	return client, nil
}

// InitializePayment initializes a new payment session
// request: Payment initialization parameters
// Returns: Payment initialization response with available channels
func (c *GhionClient) InitializePayment(request *types.InitializePaymentRequest) (*types.InitializePaymentResponse, error) {
	if err := utils.ValidateInitializePaymentRequest(
		request.Amount,
		request.Reference,
		request.Currency,
		request.WebhookURL,
		request.ReturnURL,
		request.CancelURL,
	); err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"amount":      request.Amount,
		"currency":    request.Currency,
		"reference":   request.Reference,
		"description": request.Description,
		"webhook_url": request.WebhookURL,
		"return_url":  request.ReturnURL,
		"cancel_url":  request.CancelURL,
		"metadata":    request.Metadata,
	}

	// Remove empty values
	body = removeEmptyValues(body)

	var response types.InitializePaymentResponse
	if err := c.apiRequest("POST", "/checkout/initialize", body, &response, "", false); err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitPayment submits payment with chosen channel
// paymentID: Payment session ID
// request: Payment submission parameters
// Returns: Payment submission response
func (c *GhionClient) SubmitPayment(paymentID string, request *types.SubmitPaymentRequest) (*types.SubmitPaymentResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}
	if err := utils.ValidateSubmitPaymentRequest(request.Channel, request.PhoneNumber, request.AccountNumber); err != nil {
		return nil, err
	}

	body := make(map[string]interface{})
	if request.PhoneNumber != "" {
		body["phone_number"] = request.PhoneNumber
	}
	if request.AccountNumber != "" {
		body["account_number"] = request.AccountNumber
	}

	var response types.SubmitPaymentResponse
	path := fmt.Sprintf("/checkout/%s/pay/%s", paymentID, request.Channel)
	if err := c.apiRequest("POST", path, body, &response, "", false); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPaymentStatus gets the current status of a payment
// paymentID: Payment session ID
// Returns: Payment status response
func (c *GhionClient) GetPaymentStatus(paymentID string) (*types.PaymentStatusResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}

	var response types.PaymentStatusResponse
	path := fmt.Sprintf("/checkout/%s", paymentID)
	if err := c.apiRequest("GET", path, nil, &response, "", false); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCheckout gets checkout information including QR, merchant, and provider info
// paymentID: Payment session ID
// Returns: Checkout response
func (c *GhionClient) GetCheckout(paymentID string) (*types.CheckoutResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}

	var response types.CheckoutResponse
	// Checkout endpoint with verify=1 is publicly accessible
	path := fmt.Sprintf("/checkout/%s?verify=1", paymentID)
	if err := c.apiRequest("GET", path, nil, &response, c.checkoutBaseURL, true); err != nil {
		return nil, err
	}

	return &response, nil
}

// PayWithQR pays with QR code
// paymentID: Payment session ID
// Returns: QR payment response with QR image and payload
func (c *GhionClient) PayWithQR(paymentID string) (*types.QRPaymentResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}

	var response types.QRPaymentResponse
	path := fmt.Sprintf("/checkout/%s/pay/other", paymentID)
	if err := c.apiRequest("POST", path, nil, &response, c.checkoutBaseURL, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendOTP sends OTP to user's phone for YaYa Wallet payment
// paymentID: Payment session ID
// phoneNumber: User's phone number
// Returns: OTP send response
func (c *GhionClient) SendOTP(paymentID, phoneNumber string) (*types.OTPSendResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}
	if err := utils.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"phone_number":   phoneNumber,
		"payment_method": "otp",
	}

	var response types.OTPSendResponse
	path := fmt.Sprintf("/checkout/%s/pay/yayawallet", paymentID)
	if err := c.apiRequest("POST", path, body, &response, c.checkoutBaseURL, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// ValidateOTP validates OTP code for payment completion
// paymentID: Payment session ID
// otpCode: OTP code received by user
// phoneNumber: Phone number used to send OTP
// Returns: OTP validation response
func (c *GhionClient) ValidateOTP(paymentID, otpCode, phoneNumber string) (*types.OTPValidateResponse, error) {
	if err := utils.ValidatePaymentID(paymentID); err != nil {
		return nil, err
	}
	if err := utils.ValidateOTPCode(otpCode); err != nil {
		return nil, err
	}
	if err := utils.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"otp_code":     otpCode,
		"phone_number": phoneNumber,
	}

	var response types.OTPValidateResponse
	path := fmt.Sprintf("/checkout/%s/otp-validate", paymentID)
	if err := c.apiRequest("POST", path, body, &response, c.checkoutBaseURL, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// VerifyWebhook verifies webhook signature
// rawBody: Raw request body from webhook
// signature: Signature from X-Ghion-Signature header
// Returns: True if signature is valid
func (c *GhionClient) VerifyWebhook(rawBody []byte, signature string) bool {
	return c.webhookHandler.VerifyWebhook(rawBody, signature)
}

// ParseWebhook parses and verifies webhook event
// rawBody: Raw request body from webhook
// signature: Signature from X-Ghion-Signature header
// Returns: Parsed webhook event
func (c *GhionClient) ParseWebhook(rawBody []byte, signature string) (*types.WebhookEvent, error) {
	return c.webhookHandler.ParseWebhook(rawBody, signature)
}

// GetWebhookHandler returns the webhook handler for advanced webhook processing
func (c *GhionClient) GetWebhookHandler() *WebhookHandler {
	return c.webhookHandler
}

// apiRequest makes an authenticated API request
// method: HTTP method
// path: API path
// data: Request body data
// response: Response struct to unmarshal into
// customBaseURL: Custom base URL (optional)
// skipAuth: Skip authentication (for public endpoints)
func (c *GhionClient) apiRequest(method, path string, data map[string]interface{}, response interface{}, customBaseURL string, skipAuth bool) error {
	return c.apiRequestWithRetry(method, path, data, response, customBaseURL, skipAuth, 1)
}

// apiRequestWithRetry makes an authenticated API request with retry logic
func (c *GhionClient) apiRequestWithRetry(method, path string, data map[string]interface{}, response interface{}, customBaseURL string, skipAuth bool, attempt int) error {
	baseURL := customBaseURL
	if baseURL == "" {
		baseURL = c.baseURL
	}

	fullURL := baseURL + path

	var bodyReader io.Reader
	var bodyString string

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return errors.NewNetworkError("Failed to marshal request body", map[string]interface{}{
				"error": err.Error(),
			})
		}
		bodyString = string(jsonData)
		bodyReader = bytes.NewReader(jsonData)
	}

	// Parse URL to get the path for signature
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return errors.NewNetworkError("Failed to parse URL", map[string]interface{}{
			"error": err.Error(),
		})
	}
	fullPath := parsedURL.Path + parsedURL.RawQuery

	// Create request
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return errors.NewNetworkError("Failed to create request", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	if !skipAuth {
		timestamp := utils.GetCurrentTimestamp()
		signature := utils.GenerateSignature(timestamp, method, fullPath, bodyString, c.apiSecret)

		req.Header.Set("X-Ghion-Key", c.apiKey)
		req.Header.Set("X-Ghion-Timestamp", strconv.FormatInt(timestamp, 10))
		req.Header.Set("X-Ghion-Signature", signature)
		req.Header.Set("X-Ghion-Passphrase", c.passphrase)
	}

	// Execute request with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Retry network errors for safe GET methods
		if c.isSafeEndpoint(path) && c.isRetryableError(err) && attempt < maxRetries {
			waitTime := time.Duration(1000*attempt) * time.Millisecond
			time.Sleep(waitTime)
			return c.apiRequestWithRetry(method, path, data, response, customBaseURL, skipAuth, attempt+1)
		}

		if err == context.DeadlineExceeded {
			return errors.NewNetworkError(fmt.Sprintf("Request timeout after %v", c.timeout), nil)
		}
		return errors.NewNetworkError(fmt.Sprintf("Network error: %v", err), nil)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("Failed to read response body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Handle rate limiting with automatic retry
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		waitMs := 1000 * attempt
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				waitMs = seconds * 1000
			}
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(waitMs) * time.Millisecond)
			return c.apiRequestWithRetry(method, path, data, response, customBaseURL, skipAuth, attempt+1)
		}

		retryAfterInt := 0
		if retryAfter != "" {
			retryAfterInt, _ = strconv.Atoi(retryAfter)
		}
		return errors.NewRateLimitError("API rate limit exceeded", retryAfterInt)
	}

	// Retry transient failures for safe GET methods
	if c.isSafeEndpoint(path) && (resp.StatusCode >= 500 || resp.StatusCode == 0) && attempt < maxRetries {
		waitTime := time.Duration(1000*(1<<(attempt-1))) * time.Millisecond // Exponential backoff
		time.Sleep(waitTime)
		return c.apiRequestWithRetry(method, path, data, response, customBaseURL, skipAuth, attempt+1)
	}

	// Try to parse as JSON
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		// If it's HTML, it's likely an error page
		if strings.Contains(string(respBody), "<") {
			return errors.NewAPIError(
				fmt.Sprintf("API returned HTML error page (status %d). This usually means the payment session has expired or is invalid.", resp.StatusCode),
				resp.StatusCode,
				nil,
			)
		}
		return errors.NewAPIError(
			fmt.Sprintf("Invalid response from API: %s", string(respBody)[:min(100, len(respBody))]),
			resp.StatusCode,
			nil,
		)
	}

	// Handle API errors
	if resp.StatusCode >= 400 {
		errorMessage := "API error"
		if errorData, ok := jsonResp["error"].(map[string]interface{}); ok {
			if msg, ok := errorData["message"].(string); ok {
				errorMessage = msg
			}
		}
		return errors.NewAPIError(errorMessage, resp.StatusCode, c.redactSensitiveData(jsonResp, 0))
	}

	// Unmarshal response
	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return errors.NewNetworkError("Failed to unmarshal response", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return nil
}

// isSafeEndpoint checks if endpoint is safe for retry (idempotent GET operations)
func (c *GhionClient) isSafeEndpoint(path string) bool {
	return strings.Contains(path, "/checkout/") && 
		!strings.Contains(path, "/pay/") && 
		!strings.Contains(path, "/otp-validate")
}

// isRetryableError checks if an error is retryable
func (c *GhionClient) isRetryableError(err error) bool {
	return err != nil && 
		(strings.Contains(err.Error(), "timeout") || 
		 strings.Contains(err.Error(), "connection refused") ||
		 strings.Contains(err.Error(), "EOF"))
}

// redactSensitiveData redacts sensitive data from error responses
func (c *GhionClient) redactSensitiveData(data map[string]interface{}, depth int) map[string]interface{} {
	if depth > 20 {
		return data
	}

	sensitiveKeys := []string{
		"api_key", "api_secret", "passphrase", "signature", 
		"password", "otp_code", "token", "phone_number", "account_number",
	}

	redacted := make(map[string]interface{})
	for key, value := range data {
		keyLower := strings.ToLower(key)
		isSensitive := false
		for _, sk := range sensitiveKeys {
			if strings.Contains(keyLower, sk) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			redacted[key] = "[REDACTED]"
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			redacted[key] = c.redactSensitiveData(nestedMap, depth+1)
		} else if nestedSlice, ok := value.([]interface{}); ok {
			slice := make([]interface{}, len(nestedSlice))
			for i, item := range nestedSlice {
				if nestedItem, ok := item.(map[string]interface{}); ok {
					slice[i] = c.redactSensitiveData(nestedItem, depth+1)
				} else {
					slice[i] = item
				}
			}
			redacted[key] = slice
		} else {
			redacted[key] = value
		}
	}

	return redacted
}

// removeEmptyValues removes empty values from a map
func removeEmptyValues(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if v != "" {
				result[key] = v
			}
		case map[string]interface{}:
			if len(v) > 0 {
				result[key] = v
			}
		default:
			if value != nil {
				result[key] = value
			}
		}
	}
	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
