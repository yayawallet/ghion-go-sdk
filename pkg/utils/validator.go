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
