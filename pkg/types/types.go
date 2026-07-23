package types

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusExpired    PaymentStatus = "expired"
)

// WebhookEventType represents the type of webhook event
type WebhookEventType string

const (
	EventTransactionCompleted         WebhookEventType = "transaction.completed"
	EventTransactionFailed            WebhookEventType = "transaction.failed"
	EventTransactionRefunded          WebhookEventType = "transaction.refunded"
	EventTransactionPartiallyRefunded WebhookEventType = "transaction.partially_refunded"
	EventTransactionExpired           WebhookEventType = "transaction.expired"
	EventTransactionDisputed          WebhookEventType = "transaction.disputed"
	EventTransactionUpdated           WebhookEventType = "transaction.updated"
)

// GhionConfig represents the configuration for the Ghion client
type GhionConfig struct {
	APIKey          string
	APISecret       string
	Passphrase      string
	BaseURL         string // Optional, defaults to production
	CheckoutBaseURL string // Optional, for checkout endpoints
	Timeout         int    // Optional, request timeout in milliseconds
}

// InitializePaymentRequest represents a payment initialization request
type InitializePaymentRequest struct {
	Amount      float64              `json:"amount"`
	Currency    string               `json:"currency,omitempty"`
	Reference   string               `json:"reference"`
	Description string               `json:"description,omitempty"`
	WebhookURL  string               `json:"webhook_url,omitempty"`
	ReturnURL   string               `json:"return_url,omitempty"`
	CancelURL   string               `json:"cancel_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// InitializePaymentResponse represents a payment initialization response
type InitializePaymentResponse struct {
	ID               string           `json:"id"`
	Amount           float64          `json:"amount"`
	Currency         string           `json:"currency"`
	Reference        string           `json:"reference"`
	Description      string           `json:"description"`
	Status           string           `json:"status"`
	Channels         []PaymentChannel `json:"channels"`
	AvailableChannels []PaymentChannel `json:"available_channels,omitempty"`
	ExpiresAt        string           `json:"expires_at"`
	CreatedAt        string           `json:"created_at"`
	Providers        []Provider       `json:"providers,omitempty"`
	CardEnabled      bool             `json:"card_enabled,omitempty"`
	OtherEnabled     bool             `json:"other_enabled,omitempty"`
	QR               *QRInfo          `json:"qr,omitempty"`
	YayaUniqueRef    string           `json:"yaya_unique_reference,omitempty"`
	CheckoutURL      string           `json:"checkout_url,omitempty"`
	AllowAmountEdit  bool             `json:"allow_amount_edit,omitempty"`
	CallbackURL      string           `json:"callback_url,omitempty"`
	ReturnURL        string           `json:"return_url,omitempty"`
	CancelURL        string           `json:"cancel_url,omitempty"`
}

// PaymentChannel represents an available payment channel
type PaymentChannel struct {
	ID             string `json:"id,omitempty"`
	Code           string `json:"code,omitempty"`
	Name           string `json:"name"`
	Icon           string `json:"icon,omitempty"`
	RequiresPhone  bool   `json:"requires_phone,omitempty"`
	RequiresAccount bool  `json:"requires_account,omitempty"`
	Type           string `json:"type,omitempty"`
	Logo           string `json:"logo,omitempty"`
	SupportsOTP    bool   `json:"supports_otp,omitempty"`
}

// SubmitPaymentRequest represents a payment submission request
type SubmitPaymentRequest struct {
	Channel        string `json:"channel"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	AccountNumber  string `json:"account_number,omitempty"`
	PaymentMethod  string `json:"payment_method,omitempty"` // 'ussd', 'otp', or 'qr'
}

// SubmitPaymentResponse represents a payment submission response
type SubmitPaymentResponse struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	TransactionID string `json:"transaction_id,omitempty"`
	Message      string  `json:"message"`
	RedirectURL  string  `json:"redirect_url,omitempty"`
}

// PaymentStatusResponse represents a payment status response
type PaymentStatusResponse struct {
	ID           string       `json:"id"`
	Amount       float64      `json:"amount"`
	Currency     string       `json:"currency"`
	Reference    string       `json:"reference"`
	Description  string       `json:"description"`
	Status       PaymentStatus `json:"status"`
	Channel      string       `json:"channel,omitempty"`
	TransactionID string      `json:"transaction_id,omitempty"`
	Customer     *Customer    `json:"customer,omitempty"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	CompletedAt  string       `json:"completed_at,omitempty"`
	FailedAt     string       `json:"failed_at,omitempty"`
	FailureReason string      `json:"failure_reason,omitempty"`
}

// Customer represents customer information
type Customer struct {
	PhoneNumber   string `json:"phone_number,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
}

// CheckoutResponse represents checkout information
type CheckoutResponse struct {
	ID               string           `json:"id"`
	Reference        string           `json:"reference"`
	Amount           float64          `json:"amount"`
	Currency         string           `json:"currency"`
	Status           string           `json:"status"`
	Mode             string           `json:"mode,omitempty"`
	Description      string           `json:"description"`
	Channel          string           `json:"channel,omitempty"`
	IsExpired        bool             `json:"is_expired"`
	ExpiresAt        string           `json:"expires_at"`
	Providers        []Provider       `json:"providers,omitempty"`
	CardEnabled      bool             `json:"card_enabled,omitempty"`
	OtherEnabled     bool             `json:"other_enabled,omitempty"`
	AvailableChannels []PaymentChannel `json:"available_channels,omitempty"`
	CreatedAt        string           `json:"created_at"`
	CallbackURL      string           `json:"callback_url,omitempty"`
	ReturnURL        string           `json:"return_url,omitempty"`
	CancelURL        string           `json:"cancel_url,omitempty"`
	YayaUniqueRef    string           `json:"yaya_unique_reference,omitempty"`
	QR               *QRInfo          `json:"qr,omitempty"`
	Merchant         *Merchant        `json:"merchant,omitempty"`
	CollectPhone     bool             `json:"collect_phone,omitempty"`
	CollectEmail     bool             `json:"collect_email,omitempty"`
	AllowAmountEdit  bool             `json:"allow_amount_edit,omitempty"`
	PayerPhone       string           `json:"payer_phone,omitempty"`
	FeeOnMerchant    bool             `json:"fee_on_merchant,omitempty"`
	GatewayFee       float64          `json:"gateway_fee,omitempty"`
	TotalAmount      float64          `json:"total_amount,omitempty"`
	CheckoutURL      string           `json:"checkout_url,omitempty"`
}

// QRPaymentResponse represents a QR payment response
type QRPaymentResponse struct {
	Type         string `json:"type"`
	TransactionID string `json:"transaction_id"`
	Status       string `json:"status"`
	QRImageURL   string `json:"qr_image_url"`
	QRPayload    string `json:"qr_payload"`
}

// OTPSendResponse represents an OTP send response
type OTPSendResponse struct {
	Type         string `json:"type"`
	TransactionID string `json:"transaction_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

// OTPValidateResponse represents an OTP validation response
type OTPValidateResponse struct {
	Status       string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// Provider represents payment provider information
type Provider struct {
	Code   string   `json:"code"`
	Name   string   `json:"name"`
	Logo   string   `json:"logo,omitempty"`
	Methods []string `json:"methods"`
}

// QRInfo represents QR code information
type QRInfo struct {
	QRImageURL string `json:"qr_image_url"`
	QRPayload  string `json:"qr_payload"`
}

// Merchant represents merchant information
type Merchant struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	BusinessName string `json:"business_name"`
	LogoURL      string `json:"logo_url,omitempty"`
}

// WebhookEvent represents a webhook event payload
type WebhookEvent struct {
	Event     WebhookEventType `json:"event"`
	Data      WebhookEventData `json:"data"`
	Signature string           `json:"signature"`
}

// WebhookEventData represents the data within a webhook event
type WebhookEventData struct {
	PaymentID     string       `json:"payment_id"`
	TransactionID string       `json:"transaction_id,omitempty"`
	Amount        float64      `json:"amount"`
	Currency      string       `json:"currency"`
	Reference     string       `json:"reference"`
	Status        PaymentStatus `json:"status"`
	Timestamp     string       `json:"timestamp"`
}

// APIErrorResponse represents an API error response
type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

// APIErrorDetail represents the error detail
type APIErrorDetail struct {
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}
