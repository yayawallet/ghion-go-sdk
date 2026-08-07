package utils

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/yayawallet/ghion-go-sdk/pkg/errors"
)

// ValidateAPIKey validates API key format
func ValidateAPIKey(apiKey string) error {
	if err := requireNonEmptyString(apiKey, "API key"); err != nil {
		return err
	}
	return nil
}

// ValidateAPISecret validates API secret format
func ValidateAPISecret(apiSecret string) error {
	if err := requireNonEmptyString(apiSecret, "API secret"); err != nil {
		return err
	}
	return nil
}

// ValidatePassphrase validates passphrase format
func ValidatePassphrase(passphrase string) error {
	if err := requireNonEmptyString(passphrase, "Passphrase"); err != nil {
		return err
	}
	return nil
}

// ValidateInitializePaymentRequest validates payment initialization request
func ValidateInitializePaymentRequest(amount float64, reference, currency, webhookURL, returnURL, cancelURL string) error {
	if amount <= 0 {
		return errors.NewValidationError("Amount must be a positive number", "amount", amount)
	}

	if err := requireNonEmptyString(reference, "Reference"); err != nil {
		return err
	}

	if currency != "" {
		if err := requireStringIfPresent(currency, "Currency"); err != nil {
			return err
		}
	}

	if webhookURL != "" {
		if !isValidURL(webhookURL) {
			return errors.NewValidationError("Webhook URL must be a valid URL", "webhook_url", webhookURL)
		}
	}

	if returnURL != "" {
		if !isValidURL(returnURL) {
			return errors.NewValidationError("Return URL must be a valid URL", "return_url", returnURL)
		}
	}

	if cancelURL != "" {
		if !isValidURL(cancelURL) {
			return errors.NewValidationError("Cancel URL must be a valid URL", "cancel_url", cancelURL)
		}
	}

	return nil
}

// ValidateSubmitPaymentRequest validates payment submission request
func ValidateSubmitPaymentRequest(channel, phoneNumber, accountNumber string) error {
	if err := requireNonEmptyString(channel, "Channel"); err != nil {
		return err
	}

	if phoneNumber != "" {
		if err := requireStringIfPresent(phoneNumber, "Phone number"); err != nil {
			return err
		}
	}

	if accountNumber != "" {
		if err := requireStringIfPresent(accountNumber, "Account number"); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePaymentID validates payment ID
func ValidatePaymentID(paymentID string) error {
	return requireNonEmptyString(paymentID, "Payment ID")
}

// ValidatePhoneNumber validates phone number for OTP
func ValidatePhoneNumber(phoneNumber string) error {
	return requireNonEmptyString(phoneNumber, "Phone number")
}

// ValidateOTPCode validates OTP code
func ValidateOTPCode(otpCode string) error {
	if otpCode == "" || strings.TrimSpace(otpCode) == "" {
		return errors.NewValidationError("OTP code is required and must be a non-empty string", "otp_code", otpCode)
	}
	return nil
}

// requireNonEmptyString validates that a value is a non-empty string
func requireNonEmptyString(value, fieldName string) error {
	if value == "" || strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fieldName+" is required and must be a non-empty string", fieldName, value)
	}
	return nil
}

// requireStringIfPresent validates that if a value is present, it must be a string
func requireStringIfPresent(value, fieldName string) error {
	if value != "" && strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fieldName+" must be a non-empty string if provided", fieldName, value)
	}
	return nil
}

// isValidURL checks if a string is a valid URL
func isValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// ValidateEthiopianPhoneNumber validates Ethiopian phone number format
func ValidateEthiopianPhoneNumber(phoneNumber string) error {
	// Remove any spaces or special characters
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phoneNumber, "")

	// Ethiopian phone numbers start with +251 followed by 9 digits
	// Or can start with 0 followed by 9 digits
	ethiopianPattern := regexp.MustCompile(`^(\+251|0)?9\d{8}$`)

	if !ethiopianPattern.MatchString(cleaned) {
		return errors.NewValidationError("Invalid Ethiopian phone number format. Expected format: +2519XXXXXXXX or 09XXXXXXXX", "phone_number", phoneNumber)
	}

	return nil
}

// ValidateBillID validates bill ID
func ValidateBillID(billID string) error {
	return requireNonEmptyString(billID, "Bill ID")
}

// ValidateCreateBillRequest validates bill creation request
func ValidateCreateBillRequest(billID string, amount float64, dueDate, customerEmail string) error {
	if err := requireNonEmptyString(billID, "Bill ID"); err != nil {
		return err
	}

	if amount <= 0 {
		return errors.NewValidationError("Amount must be a positive number", "amount", amount)
	}

	if err := requireNonEmptyString(dueDate, "Due date"); err != nil {
		return err
	}

	if !isValidDate(dueDate) {
		return errors.NewValidationError("Due date must be in Y-m-d format", "due_date", dueDate)
	}

	if customerEmail != "" {
		if !isValidEmail(customerEmail) {
			return errors.NewValidationError("Customer email must be a valid email address", "customer_email", customerEmail)
		}
	}

	return nil
}

// ValidateBulkCreateBillsRequest validates bulk bill creation request
func ValidateBulkCreateBillsRequest(bills []interface{}) error {
	if len(bills) == 0 {
		return errors.NewValidationError("Bills array is required and must not be empty", "bills", bills)
	}
	return nil
}

// ValidateListBillsRequest validates list bills request
func ValidateListBillsRequest(status, search, cluster, billCode, from, to string, page, limit int) error {
	if status != "" {
		if err := requireStringIfPresent(status, "Status"); err != nil {
			return err
		}
	}

	if search != "" {
		if err := requireStringIfPresent(search, "Search"); err != nil {
			return err
		}
	}

	if cluster != "" {
		if err := requireStringIfPresent(cluster, "Cluster"); err != nil {
			return err
		}
	}

	if billCode != "" {
		if err := requireStringIfPresent(billCode, "Bill code"); err != nil {
			return err
		}
	}

	if from != "" {
		if !isValidDate(from) {
			return errors.NewValidationError("From date must be in Y-m-d format", "from", from)
		}
	}

	if to != "" {
		if !isValidDate(to) {
			return errors.NewValidationError("To date must be in Y-m-d format", "to", to)
		}
	}

	if page < 0 {
		return errors.NewValidationError("Page must be a non-negative number", "page", page)
	}

	if limit < 0 || limit > 100 {
		return errors.NewValidationError("Limit must be a number between 0 and 100", "limit", limit)
	}

	return nil
}

// ValidatePublicBillLookupRequest validates public bill lookup request
func ValidatePublicBillLookupRequest(billerCode, billID string) error {
	if err := requireNonEmptyString(billerCode, "Biller code"); err != nil {
		return err
	}
	if err := requireNonEmptyString(billID, "Bill ID"); err != nil {
		return err
	}
	return nil
}

// ValidateRecordManualPaymentRequest validates manual payment recording request
func ValidateRecordManualPaymentRequest(amount float64) error {
	if amount <= 0 {
		return errors.NewValidationError("Amount must be a positive number", "amount", amount)
	}
	return nil
}

// ValidateBillerSettingsRequest validates biller settings request
func ValidateBillerSettingsRequest(serviceChargeRate float64, clusters, billCodes []interface{}) error {
	if serviceChargeRate < 0 {
		return errors.NewValidationError("Service charge rate must be a non-negative number", "service_charge_rate", serviceChargeRate)
	}

	if clusters != nil && len(clusters) > 0 {
		// Check if all elements are strings
		for _, cluster := range clusters {
			if _, ok := cluster.(string); !ok {
				return errors.NewValidationError("Clusters must be an array of strings", "clusters", clusters)
			}
		}
	}

	if billCodes != nil && len(billCodes) > 0 {
		// Check if all elements are maps (bill code objects)
		for _, billCode := range billCodes {
			if _, ok := billCode.(map[string]interface{}); !ok {
				return errors.NewValidationError("Bill codes must be an array of bill code objects", "bill_codes", billCodes)
			}
		}
	}

	return nil
}

// ValidateBillDashboardRequest validates bill dashboard date parameters
func ValidateBillDashboardRequest(from, to string) error {
	if from != "" {
		if !isValidDate(from) {
			return errors.NewValidationError("From date must be in Y-m-d format", "from", from)
		}
	}

	if to != "" {
		if !isValidDate(to) {
			return errors.NewValidationError("To date must be in Y-m-d format", "to", to)
		}
	}

	return nil
}

// isValidDate checks if a string is a valid date in Y-m-d format
func isValidDate(dateStr string) bool {
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return datePattern.MatchString(dateStr)
}

// isValidEmail checks if a string is a valid email
func isValidEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return emailPattern.MatchString(email)
}
