package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	ghion "github.com/yayawallet/ghion-go-sdk"
	"github.com/yayawallet/ghion-go-sdk/pkg/types"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not load .env file, using environment variables")
	}

	// Initialize the client
	client, err := ghion.NewClient(&ghion.Config{
		APIKey:     getEnv("GHION_API_KEY"),
		APISecret:  getEnv("GHION_API_SECRET"),
		Passphrase: getEnv("GHION_API_PASSPHRASE"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== Bill Payment API Examples ===")
	fmt.Println()

	// Example 1: Create a single bill
	createBillExample(client)

	// Example 2: Create bulk bills
	createBulkBillsExample(client)

	// Example 3: List bills with filters
	listBillsExample(client)

	// Example 4: Get bill statistics
	getBillStatisticsExample(client)

	// Example 5: Get bill dashboard
	getBillDashboardExample(client)

	// Example 6: Get bill details
	getBillDetailExample(client)

	// Example 7: Update a bill
	updateBillExample(client)

	// Example 8: Record a manual payment
	recordManualPaymentExample(client)

	// Example 9: Get bill payment link
	getBillPaymentLinkExample(client)

	// Example 10: Get biller settings
	getBillerSettingsExample(client)

	// Example 11: Public bill lookup
	publicBillLookupExample(client)

	// Example 12: Delete a bill
	deleteBillExample(client)
}

func createBillExample(client *ghion.Client) {
	fmt.Println("--- Example 1: Create a Single Bill ---")

	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:        fmt.Sprintf("INV-%d", time.Now().Unix()),
		Amount:        500.00,
		Currency:      "ETB",
		DueDate:       "2026-09-01",
		CustomerName:  "John Doe",
		CustomerPhone: "+251911234567",
		CustomerEmail: "john@example.com",
		Description:   "Monthly utility bill",
		BillCode:      "UTIL",
		Cluster:       "ADDIS ABABA",
	})
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill created successfully!\n")
	fmt.Printf("  ID: %s\n", bill.ID)
	fmt.Printf("  Bill ID: %s\n", bill.BillID)
	fmt.Printf("  Amount: %.2f %s\n", bill.Amount, bill.Currency)
	fmt.Printf("  Status: %s\n", bill.Status)
	fmt.Printf("  Due Date: %s\n\n", bill.DueDate)
}

func createBulkBillsExample(client *ghion.Client) {
	fmt.Println("--- Example 2: Create Bulk Bills ---")

	request := &types.BulkCreateBillsRequest{
		Bills: []types.CreateBillRequest{
			{
				BillID:       fmt.Sprintf("BULK-1-%d", time.Now().Unix()),
				Amount:       300.00,
				Currency:     "ETB",
				DueDate:      "2026-09-01",
				CustomerName: "Customer 1",
				Description:  "Bulk bill 1",
			},
			{
				BillID:       fmt.Sprintf("BULK-2-%d", time.Now().Unix()),
				Amount:       400.00,
				Currency:     "ETB",
				DueDate:      "2026-09-01",
				CustomerName: "Customer 2",
				Description:  "Bulk bill 2",
			},
		},
	}

	response, err := client.CreateBulkBills(request)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bulk bills created successfully!\n")
	fmt.Printf("  Created: %d\n", response.CreatedCount)
	fmt.Printf("  Errors: %d\n\n", response.ErrorCount)
}

func listBillsExample(client *ghion.Client) {
	fmt.Println("--- Example 3: List Bills with Filters ---")

	response, err := client.ListBills(&types.ListBillsRequest{
		Status: "pending",
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bills listed successfully!\n")
	fmt.Printf("  Total: %d\n", response.Total)
	fmt.Printf("  Page: %d\n", response.Page)
	fmt.Printf("  Items: %d\n", len(response.Items))

	for _, bill := range response.Items {
		fmt.Printf("  - %s: %.2f %s (%s)\n", bill.BillID, bill.Amount, bill.Currency, bill.Status)
	}
	fmt.Println()
}

func getBillStatisticsExample(client *ghion.Client) {
	fmt.Println("--- Example 4: Get Bill Statistics ---")

	stats, err := client.GetBillStatistics()
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill statistics retrieved successfully!\n")
	fmt.Printf("  Pending: %d\n", stats.Pending)
	fmt.Printf("  Paid: %d\n", stats.Paid)
	fmt.Printf("  Forwarded: %d\n", stats.Forwarded)
	fmt.Printf("  Overdue: %d\n", stats.Overdue)
	fmt.Printf("  Total Amount: %.2f\n", stats.TotalAmount)
	fmt.Printf("  Total Paid: %.2f\n\n", stats.TotalPaid)
}

func getBillDashboardExample(client *ghion.Client) {
	fmt.Println("--- Example 5: Get Bill Dashboard ---")

	dashboard, err := client.GetBillDashboard("2026-08-01", "2026-08-31")
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill dashboard retrieved successfully!\n")
	fmt.Printf("  Total Bills: %d\n", dashboard.Summary.TotalBills)
	fmt.Printf("  Pending: %d\n", dashboard.Summary.Pending)
	fmt.Printf("  Paid: %d\n", dashboard.Summary.Paid)
	fmt.Printf("  Total Amount: %.2f\n", dashboard.Summary.TotalAmount)
	fmt.Printf("  Clusters: %d\n", len(dashboard.ByCluster))
	fmt.Printf("  Bill Codes: %d\n", len(dashboard.ByBillCode))
	fmt.Printf("  Trend Data Points: %d\n\n", len(dashboard.Trend))
}

func getBillDetailExample(client *ghion.Client) {
	fmt.Println("--- Example 6: Get Bill Details ---")

	// First create a bill to get its ID
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("DETAIL-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Bill for detail example",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	detail, err := client.GetBillDetail(bill.ID)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill details retrieved successfully!\n")
	fmt.Printf("  Bill ID: %s\n", detail.BillID)
	fmt.Printf("  Customer: %s\n", detail.CustomerName)
	fmt.Printf("  Amount: %.2f %s\n", detail.Amount, detail.Currency)
	fmt.Printf("  Status: %s\n", detail.Status)
	fmt.Printf("  Due Date: %s\n", detail.DueDate)
	fmt.Printf("  Payments: %d\n\n", len(detail.Payments))
}

func updateBillExample(client *ghion.Client) {
	fmt.Println("--- Example 7: Update a Bill ---")

	// First create a bill
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("UPDATE-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Original description",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	// Update the bill
	updatedBill, err := client.UpdateBill(bill.ID, &types.UpdateBillRequest{
		Amount:      bill.Amount,
		Description: "Updated description",
		DueDate:     "2026-09-15",
	})
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill updated successfully!\n")
	fmt.Printf("  ID: %s\n", updatedBill.ID)
	fmt.Printf("  Description: %s\n", updatedBill.Description)
	fmt.Printf("  Due Date: %s\n\n", updatedBill.DueDate)
}

func recordManualPaymentExample(client *ghion.Client) {
	fmt.Println("--- Example 8: Record Manual Payment ---")

	// First create a bill
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("PAY-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Bill for manual payment",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	// Record a manual payment
	payment, err := client.RecordManualPayment(bill.ID, &types.RecordManualPaymentRequest{
		Amount:        100.00,
		Source:        "manual",
		PaymentMethod: "cash",
		Reference:     "RECEIPT-001",
		Note:          "Paid at counter",
	})
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Manual payment recorded successfully!\n")
	fmt.Printf("  Payment ID: %s\n", payment.PaymentID)
	fmt.Printf("  Amount: %.2f\n", payment.Amount)
	fmt.Printf("  Bill Status: %s\n", payment.BillStatus)
	fmt.Printf("  Balance Due: %.2f\n\n", payment.BalanceDue)
}

func getBillPaymentLinkExample(client *ghion.Client) {
	fmt.Println("--- Example 9: Get Bill Payment Link ---")

	// First create a bill
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("LINK-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Bill for payment link",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	// Get payment link
	link, err := client.GetBillPaymentLink(bill.ID)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Payment link generated successfully!\n")
	fmt.Printf("  Checkout URL: %s\n\n", link.CheckoutURL)
}

func getBillerSettingsExample(client *ghion.Client) {
	fmt.Println("--- Example 10: Get Biller Settings ---")

	settings, err := client.GetBillerSettings()
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Biller settings retrieved successfully!\n")
	fmt.Printf("  Configured: %v\n", settings.Configured)
	if settings.Settings != nil {
		fmt.Printf("  Biller Code: %s\n", settings.Settings.BillerCode)
		fmt.Printf("  Biller Name: %s\n", settings.Settings.BillerName)
		fmt.Printf("  Biller Category: %s\n", settings.Settings.BillerCategory)
		fmt.Printf("  Clusters: %d\n", len(settings.Settings.Clusters))
		fmt.Printf("  Bill Codes: %d\n\n", len(settings.Settings.BillCodes))
	}
}

func publicBillLookupExample(client *ghion.Client) {
	fmt.Println("--- Example 11: Public Bill Lookup ---")

	// Get biller settings first to get biller code
	settings, err := client.GetBillerSettings()
	if err != nil {
		fmt.Printf("Error getting biller settings: %v\n\n", err)
		return
	}

	if !settings.Configured || settings.Settings == nil {
		fmt.Println("Biller not configured, skipping public lookup")
		fmt.Println()
		return
	}

	// Create a bill for lookup
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("LOOKUP-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Bill for public lookup",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	// Public bill lookup
	publicBill, err := client.PublicBillLookup(&types.PublicBillLookupRequest{
		BillerCode: settings.Settings.BillerCode,
		BillID:     bill.BillID,
	})
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Public bill lookup successful!\n")
	fmt.Printf("  Bill ID: %s\n", publicBill.BillID)
	fmt.Printf("  Customer: %s\n", publicBill.CustomerName)
	fmt.Printf("  Amount: %.2f\n", publicBill.Amount)
	fmt.Printf("  Payment Status: %s\n", publicBill.PaymentStatus)
	fmt.Printf("  Amount Due: %.2f\n\n", publicBill.AmountDue)
}

func deleteBillExample(client *ghion.Client) {
	fmt.Println("--- Example 12: Delete a Bill ---")

	// First create a bill
	bill, err := client.CreateBill(&types.CreateBillRequest{
		BillID:       fmt.Sprintf("DELETE-%d", time.Now().Unix()),
		Amount:       500.00,
		Currency:     "ETB",
		DueDate:      "2026-09-01",
		CustomerName: "Test Customer",
		Description:  "Bill for deletion",
	})
	if err != nil {
		fmt.Printf("Error creating bill: %v\n\n", err)
		return
	}

	// Delete the bill
	result, err := client.DeleteBill(bill.ID)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bill deleted successfully!\n")
	fmt.Printf("  Message: %s\n\n", result.Message)
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
