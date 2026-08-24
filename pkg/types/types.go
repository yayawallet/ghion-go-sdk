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

// BillStatus represents the status of a bill
type BillStatus string

const (
	BillStatusPending   BillStatus = "pending"
	BillStatusPaid      BillStatus = "paid"
	BillStatusForwarded BillStatus = "forwarded"
	BillStatusCancelled BillStatus = "cancelled"
	BillStatusExpired   BillStatus = "expired"
	BillStatusOverdue   BillStatus = "overdue"
)

// Penalty represents late payment penalty configuration
type Penalty struct {
	Type        string  `json:"type"`         // 'fixed' or 'percentage'
	Fee         float64 `json:"fee"`          // Penalty fee or percentage rate
	MaxAmount   float64 `json:"max_amount"`   // Maximum total penalty cap
	Recurring   string  `json:"recurring"`    // 'daily', 'weekly', 'monthly', or 'once'
}

// CreateBillRequest represents a bill creation request
type CreateBillRequest struct {
	BillID        string                 `json:"bill_id,omitempty"`        // Optional: Your unique bill identifier (max 100 chars). If omitted, one will be auto-generated
	Amount        float64                `json:"amount"`                    // Required: Bill amount
	Currency      string                 `json:"currency,omitempty"`        // Optional: Currency code (default: ETB)
	DueDate       string                 `json:"due_date"`                  // Required: Due date in Y-m-d format
	StartDate     string                 `json:"start_date,omitempty"`      // Optional: Start date in Y-m-d format
	ExpiresDate   string                 `json:"expires_date,omitempty"`    // Optional: Expiry date in Y-m-d format
	CustomerName  string                 `json:"customer_name,omitempty"`   // Optional: Customer's full name
	CustomerPhone string                 `json:"customer_phone,omitempty"`  // Optional: Customer phone in international format
	CustomerEmail string                 `json:"customer_email,omitempty"`  // Optional: Customer email
	CustomerID    string                 `json:"customer_id,omitempty"`     // Optional: Your internal customer ID
	Description   string                 `json:"description,omitempty"`     // Optional: Bill description
	BillCode      string                 `json:"bill_code,omitempty"`       // Optional: Category code for filtering/routing
	Cluster       string                 `json:"cluster,omitempty"`         // Optional: Geographic/organizational cluster
	Penalty       *Penalty               `json:"penalty,omitempty"`         // Optional: Late payment penalty configuration
	Metadata      map[string]interface{} `json:"metadata,omitempty"`        // Optional: Arbitrary key-value pairs
}

// BillResponse represents a bill response
type BillResponse struct {
	ID              string                 `json:"id"`
	BillID          string                 `json:"bill_id"`
	BillCode        string                 `json:"bill_code,omitempty"`
	BillSeason      string                 `json:"bill_season,omitempty"`
	Cluster         string                 `json:"cluster,omitempty"`
	ExtCustomerID   string                 `json:"ext_customer_id,omitempty"`
	CustomerName    string                 `json:"customer_name,omitempty"`
	CustomerPhone   string                 `json:"customer_phone,omitempty"`
	CustomerEmail   string                 `json:"customer_email,omitempty"`
	CustomerID      string                 `json:"customer_id,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Amount          float64                `json:"amount"`
	ServiceCharge   float64                `json:"service_charge"`
	PenaltyAmount   float64                `json:"penalty_amount"`
	TotalDue        float64                `json:"total_due"`
	Paid            float64                `json:"paid"`
	BalanceDue      float64                `json:"balance_due"`
	Currency        string                 `json:"currency"`
	Status          BillStatus             `json:"status"`
	DueDate         string                 `json:"due_date"`
	StartDate       string                 `json:"start_date,omitempty"`
	ExpiresDate     string                 `json:"expires_date,omitempty"`
	IsOverdue       bool                   `json:"is_overdue"`
	DaysOverdue     int                    `json:"days_overdue"`
	PenaltyType     string                 `json:"penalty_type,omitempty"`
	PenaltyFee      float64                `json:"penalty_fee,omitempty"`
	MaxPenaltyAmount float64               `json:"max_penalty_amount,omitempty"`
	PenaltyRecurring string               `json:"penalty_recurring,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
	ShareToken      string                 `json:"share_token,omitempty"`
	PenaltyConfig   *Penalty               `json:"penalty,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Payments        []PaymentRecord        `json:"payments,omitempty"`
}

// PaymentRecord represents a payment record on a bill
type PaymentRecord struct {
	ID            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Source        string  `json:"source"`
	PaymentMethod string  `json:"payment_method"`
	Reference     string  `json:"reference,omitempty"`
	Note          string  `json:"note,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// BulkCreateBillsRequest represents a bulk bill creation request
type BulkCreateBillsRequest struct {
	Bills []CreateBillRequest `json:"bills"`
}

// BulkCreateBillsResponse represents a bulk bill creation response
type BulkCreateBillsResponse struct {
	CreatedCount int                   `json:"created_count"`
	ErrorCount   int                   `json:"error_count"`
	Created      []BillResponse        `json:"created"`
	Errors       []BulkBillError       `json:"errors"`
}

// BulkBillError represents an error in bulk bill creation
type BulkBillError struct {
	BillID string `json:"bill_id"`
	Error  string `json:"error"`
}

// ListBillsRequest represents a list bills request
type ListBillsRequest struct {
	Status    string  `json:"status,omitempty"`    // Filter by status
	Search    string  `json:"search,omitempty"`    // Search by bill_id, customer_name, customer_phone, or customer_id
	Cluster   string  `json:"cluster,omitempty"`   // Filter by cluster
	BillCode  string  `json:"bill_code,omitempty"` // Filter by bill code
	From      string  `json:"from,omitempty"`      // Filter by creation date (from) in Y-m-d format
	To        string  `json:"to,omitempty"`        // Filter by creation date (to) in Y-m-d format
	Page      int     `json:"page,omitempty"`      // Page number (default: 1)
	Limit     int     `json:"limit,omitempty"`     // Items per page (max: 100, default: 25)
}

// ListBillsResponse represents a list bills response
type ListBillsResponse struct {
	Items []BillResponse `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// BillStatistics represents bill statistics
type BillStatistics struct {
	Pending      int     `json:"pending"`
	Paid         int     `json:"paid"`
	Forwarded    int     `json:"forwarded"`
	Cancelled    int     `json:"cancelled"`
	Expired      int     `json:"expired"`
	Overdue      int     `json:"overdue"`
	TotalAmount  float64 `json:"total_amount"`
	TotalPaid    float64 `json:"total_paid"`
}

// BillDashboardSummary represents bill dashboard summary
type BillDashboardSummary struct {
	TotalBills       int     `json:"total_bills"`
	Pending          int     `json:"pending"`
	Paid             int     `json:"paid"`
	Forwarded        int     `json:"forwarded"`
	Cancelled        int     `json:"cancelled"`
	Expired          int     `json:"expired"`
	Overdue          int     `json:"overdue"`
	TotalAmount      float64 `json:"total_amount"`
	TotalPaid        float64 `json:"total_paid"`
	TotalOutstanding float64 `json:"total_outstanding"`
	AvgBillAmount    float64 `json:"avg_bill_amount"`
	UniqueCustomers  int     `json:"unique_customers"`
}

// BillDashboardByCluster represents bills grouped by cluster
type BillDashboardByCluster struct {
	Cluster string  `json:"cluster"`
	Count   int     `json:"count"`
	Amount  float64 `json:"amount"`
	Paid    float64 `json:"paid"`
}

// BillDashboardByBillCode represents bills grouped by bill code
type BillDashboardByBillCode struct {
	BillCode string  `json:"bill_code"`
	Count    int     `json:"count"`
	Amount   float64 `json:"amount"`
	Paid     float64 `json:"paid"`
}

// BillDashboardTrend represents bill trend data
type BillDashboardTrend struct {
	Date           string  `json:"date"`
	BillsCreated   int     `json:"bills_created"`
	AmountCreated  float64 `json:"amount_created"`
	BillsPaid      int     `json:"bills_paid"`
	AmountPaid     float64 `json:"amount_paid"`
}

// BillDashboard represents bill dashboard analytics
type BillDashboard struct {
	From        string                     `json:"from,omitempty"`
	To          string                     `json:"to,omitempty"`
	Summary     BillDashboardSummary       `json:"summary"`
	ByCluster   []BillDashboardByCluster   `json:"by_cluster,omitempty"`
	ByBillCode  []BillDashboardByBillCode  `json:"by_bill_code,omitempty"`
	Trend       []BillDashboardTrend       `json:"trend,omitempty"`
}

// PaymentLinkResponse represents a payment link response
type PaymentLinkResponse struct {
	CheckoutURL string `json:"checkout_url"`
}

// PublicBillLookupRequest represents a public bill lookup request
type PublicBillLookupRequest struct {
	BillerCode string `json:"biller_code"`
	BillID     string `json:"bill_id"`
}

// PublicBillLookupResponse represents a public bill lookup response
type PublicBillLookupResponse struct {
	ID              string  `json:"id"`
	BillID          string  `json:"bill_id"`
	BillCode        string  `json:"bill_code,omitempty"`
	BillSeason      string  `json:"bill_season,omitempty"`
	Cluster         string  `json:"cluster,omitempty"`
	ExtCustomerID   string  `json:"ext_customer_id,omitempty"`
	CustomerName    string  `json:"customer_name,omitempty"`
	Description     string  `json:"description,omitempty"`
	Amount          float64 `json:"amount"`
	ServiceCharge   float64 `json:"service_charge"`
	PenaltyAmount   float64 `json:"penalty_amount"`
	TotalDue        float64 `json:"total_due"`
	AmountDue       float64 `json:"amount_due"`
	Paid            float64 `json:"paid"`
	Currency        string  `json:"currency"`
	Client          *PublicBillClient `json:"client,omitempty"`
	StartAt         string  `json:"start_at,omitempty"`
	DueAt           string  `json:"due_at,omitempty"`
	PaymentStatus   string  `json:"payment_status"` // Uppercase: 'PENDING', 'PAID', etc.
	PenaltyType     string  `json:"penalty_type,omitempty"`
	PenaltyFee      float64 `json:"penalty_fee,omitempty"`
	MaxPenaltyAmount float64 `json:"max_penalty_amount,omitempty"`
	PenaltyRecurring string `json:"penalty_recurring,omitempty"`
}

// PublicBillClient represents client information in public bill lookup
type PublicBillClient struct {
	UniqueName string `json:"uniqueName"`
	Name       string `json:"name"`
}

// BillerSettingsRequest represents a biller settings update request
type BillerSettingsRequest struct {
	BillerName              string                 `json:"biller_name,omitempty"`
	BillerCategory          string                 `json:"biller_category,omitempty"`
	BillerDescription       string                 `json:"biller_description,omitempty"`
	IconURL                 string                 `json:"icon_url,omitempty"`
	ServiceChargeRate       float64                `json:"service_charge_rate,omitempty"`
	ServiceChargeType       string                 `json:"service_charge_type,omitempty"`
	MinServiceCharge        float64                `json:"min_service_charge,omitempty"`
	MaxServiceCharge        float64                `json:"max_service_charge,omitempty"`
	ServiceChargeRanges     []ServiceChargeRange    `json:"service_charge_ranges,omitempty"`
	Clusters                []string               `json:"clusters,omitempty"`
	BillCodes               []BillCode              `json:"bill_codes,omitempty"`
	WebhookURL              string                 `json:"webhook_url,omitempty"`
	WebhookSecret           string                 `json:"webhook_secret,omitempty"`
	SettlementBankCode      string                 `json:"settlement_bank_code,omitempty"`
	SettlementAccountNumber string                 `json:"settlement_account_number,omitempty"`
	SettlementAccountName   string                 `json:"settlement_account_name,omitempty"`
	SettlementAccounts      []SettlementAccount     `json:"settlement_accounts,omitempty"`
}

// ServiceChargeRange represents a service charge range
type ServiceChargeRange struct {
	MinAmount float64 `json:"min_amount"`
	MaxAmount float64 `json:"max_amount"`
	Rate      float64 `json:"rate"`
}

// BillCode represents a bill code configuration
type BillCode struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// SettlementAccount represents a settlement account
type SettlementAccount struct {
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// BillerSettingsResponse represents a biller settings response
type BillerSettingsResponse struct {
	Configured bool                        `json:"configured"`
	Settings   *BillerSettingsDetail       `json:"settings,omitempty"`
}

// BillerSettingsDetail represents detailed biller settings
type BillerSettingsDetail struct {
	ID                       string                 `json:"id"`
	BillerCode               string                 `json:"biller_code"`
	BillerName               string                 `json:"biller_name"`
	BillerCategory           string                 `json:"biller_category,omitempty"`
	BillerDescription        string                 `json:"biller_description,omitempty"`
	IconURL                  string                 `json:"icon_url,omitempty"`
	ServiceChargeRate        float64                `json:"service_charge_rate"`
	ServiceChargeType        string                 `json:"service_charge_type"`
	MinServiceCharge         float64                `json:"min_service_charge,omitempty"`
	MaxServiceCharge         float64                `json:"max_service_charge,omitempty"`
	ServiceChargeRanges      []ServiceChargeRange    `json:"service_charge_ranges,omitempty"`
	Clusters                 []string               `json:"clusters"`
	BillCodes                []BillCode              `json:"bill_codes"`
	WebhookURL               string                 `json:"webhook_url,omitempty"`
	WebhookSecretConfigured  bool                   `json:"webhook_secret_configured"`
	SettlementBankCode       string                 `json:"settlement_bank_code,omitempty"`
	SettlementAccountNumber  string                 `json:"settlement_account_number,omitempty"`
	SettlementAccountName    string                 `json:"settlement_account_name,omitempty"`
	SettlementAccounts       []SettlementAccount     `json:"settlement_accounts,omitempty"`
	ShortCode                string                 `json:"short_code,omitempty"`
	BillerPrefix             string                 `json:"biller_prefix,omitempty"`
	IsActive                 bool                   `json:"is_active"`
	RequiresExternalSettlement bool                  `json:"requires_external_settlement"`
	Config                   map[string]interface{} `json:"config,omitempty"`
	CreatedAt                string                 `json:"created_at"`
	UpdatedAt                string                 `json:"updated_at"`
}

// UpdateBillRequest represents a bill update request (partial update)
type UpdateBillRequest struct {
	Amount          float64                `json:"amount,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	DueDate         string                 `json:"due_date,omitempty"`
	StartDate       string                 `json:"start_date,omitempty"`
	ExpiresDate     string                 `json:"expires_date,omitempty"`
	CustomerName    string                 `json:"customer_name,omitempty"`
	CustomerPhone   string                 `json:"customer_phone,omitempty"`
	CustomerEmail   string                 `json:"customer_email,omitempty"`
	CustomerID      string                 `json:"customer_id,omitempty"`
	Description     string                 `json:"description,omitempty"`
	BillCode        string                 `json:"bill_code,omitempty"`
	Cluster         string                 `json:"cluster,omitempty"`
	Penalty         *Penalty               `json:"penalty,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// DeleteBillResponse represents a bill deletion response
type DeleteBillResponse struct {
	Message string `json:"message"`
}

// RecordManualPaymentRequest represents a manual payment recording request
type RecordManualPaymentRequest struct {
	Amount        float64 `json:"amount"`
	Source        string  `json:"source,omitempty"`         // 'manual', 'bank_transfer', 'cash', 'yaya_wallet', 'checkout'
	PaymentMethod string  `json:"payment_method,omitempty"` // 'cash', 'bank_transfer', etc.
	Reference     string  `json:"reference,omitempty"`
	Note          string  `json:"note,omitempty"`
}

// RecordManualPaymentResponse represents a manual payment recording response
type RecordManualPaymentResponse struct {
	PaymentID   string  `json:"payment_id"`
	Amount      float64 `json:"amount"`
	BillStatus  string  `json:"bill_status"`
	BalanceDue  float64 `json:"balance_due"`
}

// SendPaymentReminderRequest represents a send payment reminder request
type SendPaymentReminderRequest struct {
	Message string `json:"message,omitempty"`
}

// SendPaymentReminderResponse represents a send payment reminder response
type SendPaymentReminderResponse struct {
	Sent                bool   `json:"sent"`
	ReminderCount       int    `json:"reminder_count"`
	LastReminderSentAt  string `json:"last_reminder_sent_at"`
}

// GenerateBillIDResponse represents a generate bill ID response
type GenerateBillIDResponse struct {
	BillID string `json:"bill_id"`
}

// InitiateCheckoutResponse represents an initiate checkout response
type InitiateCheckoutResponse struct {
	PaymentLinkSlug string  `json:"payment_link_slug"`
	BalanceDue      float64 `json:"balance_due"`
	Currency        string  `json:"currency"`
	CheckoutURL     string  `json:"checkout_url"`
}
