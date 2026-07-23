//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

var client *ghion.Client

func init() {
	// Initialize client
	var err error
	client, err = ghion.NewClient(&ghion.Config{
		APIKey:     os.Getenv("GHION_API_KEY"),
		APISecret:  os.Getenv("GHION_API_SECRET"),
		Passphrase: os.Getenv("GHION_API_PASSPHRASE"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
}

func main() {
	// Setup HTTP routes
	http.HandleFunc("/api/payments/initialize", initializePaymentHandler)
	http.HandleFunc("/api/payments/", paymentStatusHandler)
	http.HandleFunc("/webhook", webhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// initializePaymentHandler handles payment initialization requests
func initializePaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount      float64              `json:"amount"`
		Currency    string               `json:"currency"`
		Reference   string               `json:"reference"`
		Description string               `json:"description"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Initialize payment
	payment, err := client.InitializePayment(&types.InitializePaymentRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Reference:   req.Reference,
		Description: req.Description,
		WebhookURL:  fmt.Sprintf("%s/webhook", getBaseURL(r)),
		Metadata:    req.Metadata,
	})
	if err != nil {
		log.Printf("Failed to initialize payment: %v", err)
		http.Error(w, fmt.Sprintf("Failed to initialize payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    payment,
	})
}

// paymentStatusHandler handles payment status requests
func paymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	paymentID := r.URL.Path[len("/api/payments/"):]
	if paymentID == "" {
		http.Error(w, "Payment ID required", http.StatusBadRequest)
		return
	}

	status, err := client.GetPaymentStatus(paymentID)
	if err != nil {
		log.Printf("Failed to get payment status: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get payment status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    status,
	})
}

// webhookHandler handles webhook events from Ghion
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Ghion-Signature")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusBadRequest)
		return
	}

	// Parse and verify webhook
	event, err := client.ParseWebhook(body, signature)
	if err != nil {
		log.Printf("Webhook verification failed: %v", err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	log.Printf("Webhook received: %s", event.Event)
	log.Printf("Payment ID: %s", event.Data.PaymentID)
	log.Printf("Status: %s", event.Data.Status)

	// Handle different event types
	switch event.Event {
	case ghion.EventTransactionCompleted:
		handleTransactionCompleted(event)
	case ghion.EventTransactionFailed:
		handleTransactionFailed(event)
	case ghion.EventTransactionExpired:
		handleTransactionExpired(event)
	default:
		log.Printf("Unhandled webhook event: %s", event.Event)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"received": true,
	})
}

func handleTransactionCompleted(event *types.WebhookEvent) {
	log.Printf("Transaction completed: %s", event.Data.PaymentID)
	// Update your database, send confirmation email, etc.
}

func handleTransactionFailed(event *types.WebhookEvent) {
	log.Printf("Transaction failed: %s", event.Data.PaymentID)
	// Handle failed transaction
}

func handleTransactionExpired(event *types.WebhookEvent) {
	log.Printf("Transaction expired: %s", event.Data.PaymentID)
	// Handle expired transaction
}

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
