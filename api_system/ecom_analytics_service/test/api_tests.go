package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// CONFIGURATION
const (
	BaseURL = "http://localhost:9001/v1"

	ShopToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJTRUxMRVIiLCJpc3MiOiJsZW1hcmNoZW5vYmxlLmlkLnZuIiwiZXhwIjo0OTE3NTExMjUyLCJpYXQiOjE3NjE3NTEyNTIsInVzZXJJZCI6IjE2NzQwOGUzLWFmZWYtNDhiOS04ZTRmLTZkZDQxZWJmMzQ2NCIsImp0aSI6ImU2YzgyN2E2LTIyOTYtNGNlOC1iMjQ1LWM3MDIxNWM4MGJjNyIsImVtYWlsIjoidmluaGhpZW4xMnpAZ21haWwuY29tIn0.2N5eJfouPOkZ47Nh2PLdrtJ_Md7zoDuvCfPn7XqncuQ"
	AdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJoaWVubGF6YWRhIiwic2NvcGUiOiJBRE1JTiIsImlzcyI6ImxlbWFyY2hlbm9ibGUuaWQudm4iLCJleHAiOjQ5MTc1MTEyNTIsImlhdCI6MTc2MTc1MTI1MiwidXNlcklkIjoiMTY3NDA4ZTMtYWZlZi00OGI5LThlNGYtNmRkNDFlYmYzNDY0IiwianRpIjoiZTZjODI3YTYtMjI5Ni00Y2U4LWIyNDUtYzcwMjE1YzgwYmM3IiwiZW1haWwiOiJ2aW5oaGllbjEyekBnbWFpbC5jb20ifQ.RnuEug1ThzMTTejZOBhq0zHOiEkL-4lrwYdsMp0mADM"

	TestShopID      = "uuid-shop"
	TestOrderID     = "31456046-7055-462b-9233-129aa3e749c5"
	TestShopOrderID = "0c1bd408-b9d5-4640-b686-02e1c12a8c32"
	TestProductID   = "69310fca-3af6-4a53-88b0-d7114aeecef0"
	TestVoucherID   = "v-expired-03"
	TestLedgerID    = "111111111111111111111111111111111111"
	TestUserID      = "user123"

	TestStartDate  = "2025-10-01"
	TestEndDate    = "2025-11-30"
	TestSingleDate = "2025-10-18"
)

// TEST RESULT TRACKING
var (
	totalTests  = 0
	passedTests = 0
	failedTests = 0
)

// Colors
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

type APIResponse struct {
	Code    int             `json:"code"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func printSection(title string) {
	fmt.Printf("\n%s========================================%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s%s%s\n", ColorBlue, title, ColorReset)
	fmt.Printf("%s========================================%s\n\n", ColorBlue, ColorReset)
}

func printTest(description string) {
	totalTests++
	fmt.Printf("%s[TEST %d] %s%s\n", ColorYellow, totalTests, description, ColorReset)
}

func printPass(message string) {
	passedTests++
	fmt.Printf("%s✓ PASS: %s%s\n", ColorGreen, message, ColorReset)
}

func printFail(message string) {
	failedTests++
	fmt.Printf("%s✗ FAIL: %s%s\n", ColorRed, message, ColorReset)
}

func printSummary() {
	printSection("TEST SUMMARY")
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("%sPassed: %d%s\n", ColorGreen, passedTests, ColorReset)
	fmt.Printf("%sFailed: %d%s\n", ColorRed, failedTests, ColorReset)

	if failedTests == 0 {
		fmt.Printf("\n%s✓ All tests passed! 🎉%s\n\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("\n%s✗ Some tests failed.%s\n\n", ColorRed, ColorReset)
	}
}

func TestAPI(description, method, endpoint, token string, expectedStatus int) {
	printTest(description)

	url := BaseURL + endpoint
	fmt.Printf("  → %s %s\n", method, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		printFail(fmt.Sprintf("%s - Error creating request: %v", description, err))
		return
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		printFail(fmt.Sprintf("%s - Error executing request: %v", description, err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		printFail(fmt.Sprintf("%s - Error reading response: %v", description, err))
		return
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		printFail(fmt.Sprintf("%s - Invalid JSON response: %v", description, err))
		fmt.Printf("  Response body: %s\n", string(body))
		return
	}

	if resp.StatusCode != expectedStatus {
		printFail(fmt.Sprintf("%s - Expected status %d, got %d", description, expectedStatus, resp.StatusCode))
		fmt.Printf("  Response: %s\n", string(body))
		return
	}

	if expectedStatus == 200 {
		if apiResp.Status != "success" {
			printFail(fmt.Sprintf("%s - Expected status 'success', got '%s'", description, apiResp.Status))
			fmt.Printf("  Response: %s\n", string(body))
			return
		}
		if len(apiResp.Result) == 0 {
			printFail(fmt.Sprintf("%s - Missing 'result' field in response", description))
			return
		}
		printPass(fmt.Sprintf("%s (Status: %d)", description, resp.StatusCode))
	} else {
		if apiResp.Status == "" {
			printFail(fmt.Sprintf("%s - Invalid error response structure", description))
			return
		}
		printPass(fmt.Sprintf("%s (Expected error: %d)", description, resp.StatusCode))
	}

	time.Sleep(300 * time.Millisecond)
}

func TestShopAPIs() {
	printSection("TESTING SHOP APIs")

	printSection("Shop - Nhóm 1: Tổng quan")

	TestAPI("Shop Overview - Full params", "GET",
		"/shop/overview?start_date="+TestStartDate+"&end_date="+TestEndDate, ShopToken, 200)
	TestAPI("Shop Overview - No params", "GET", "/shop/overview", ShopToken, 200)
	TestAPI("Shop Overview - Invalid date", "GET", "/shop/overview?start_date=2025/10/01", ShopToken, 400)
	TestAPI("Shop Overview - No token", "GET", "/shop/overview", "", 401)
	TestAPI("Shop Overview - Wrong role", "GET", "/shop/overview", AdminToken, 403)
	TestAPI("Shop Wallet Summary", "GET", "/shop/wallet/summary", ShopToken, 200)

	printSection("Shop - Nhóm 2: Phân tích Đơn hàng")

	TestAPI("List Shop Orders - Full params", "GET",
		"/shop/orders?status=COMPLETED&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=10&offset=0", ShopToken, 200)
	TestAPI("List Shop Orders - No filters", "GET", "/shop/orders", ShopToken, 200)

	statuses := []string{"AWAITING_PAYMENT", "PROCESSING", "SHIPPED", "COMPLETED", "CANCELLED", "REFUNDED"}
	for _, status := range statuses {
		TestAPI(fmt.Sprintf("List Shop Orders - Status: %s", status), "GET", "/shop/orders?status="+status, ShopToken, 200)
	}

	TestAPI("List Shop Orders - Invalid limit", "GET", "/shop/orders?limit=abc", ShopToken, 400)
	TestAPI("Get Enriched Shop Order", "GET", "/shop/orders/"+TestShopOrderID, ShopToken, 200)
	TestAPI("List Shop Order Items - Full", "GET",
		"/shop/order-items?product_id="+TestProductID+"&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20&offset=0", ShopToken, 200)
	TestAPI("List Shop Order Items - No filters", "GET", "/shop/order-items", ShopToken, 200)

	printSection("Shop - Nhóm 3: Doanh thu & Dòng tiền")

	TestAPI("Shop Revenue Timeseries - Full", "GET",
		"/shop/revenue/timeseries?start_date="+TestStartDate+"&end_date="+TestEndDate, ShopToken, 200)
	TestAPI("Shop Revenue Timeseries - No params", "GET", "/shop/revenue/timeseries", ShopToken, 200)
	TestAPI("Shop Wallet Ledger Entries - With pagination", "GET",
		"/shop/wallet/ledger-entries?limit=50&offset=0", ShopToken, 200)
	TestAPI("Shop Wallet Ledger Entries - No params", "GET", "/shop/wallet/ledger-entries", ShopToken, 200)
	TestAPI("List Shop Settlements - Full", "GET",
		"/shop/settlements?status=SETTLED&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20&offset=0", ShopToken, 200)

	settlementStatuses := []string{"PENDING_SETTLEMENT", "FUNDS_HELD", "SETTLED", "FAILED"}
	for _, status := range settlementStatuses {
		TestAPI(fmt.Sprintf("List Shop Settlements - Status: %s", status), "GET",
			"/shop/settlements?status="+status, ShopToken, 200)
	}

	printSection("Shop - Nhóm 4: Voucher")

	TestAPI("List Shop Vouchers - Active", "GET", "/shop/vouchers?is_active=true&limit=20&offset=0", ShopToken, 200)
	TestAPI("List Shop Vouchers - Inactive", "GET", "/shop/vouchers?is_active=false", ShopToken, 200)
	TestAPI("List Shop Vouchers - All", "GET", "/shop/vouchers", ShopToken, 200)
	TestAPI("List Shop Vouchers - Invalid is_active", "GET", "/shop/vouchers?is_active=maybe", ShopToken, 400)
	TestAPI("Shop Voucher Performance - Full", "GET",
		"/shop/vouchers/performance?start_date="+TestStartDate+"&end_date="+TestEndDate, ShopToken, 200)
	TestAPI("Shop Voucher Performance - No params", "GET", "/shop/vouchers/performance", ShopToken, 200)
	TestAPI("Shop Voucher Usage Details", "GET",
		"/shop/vouchers/"+TestVoucherID+"/details?limit=20&offset=0", ShopToken, 200)

	printSection("Shop - Nhóm 5: Xếp hạng")

	TestAPI("Shop Ranking Products - Sort by revenue", "GET",
		"/shop/ranking/products?start_date="+TestStartDate+"&end_date="+TestEndDate+"&sort_by=revenue&limit=10", ShopToken, 200)
	TestAPI("Shop Ranking Products - Sort by quantity", "GET",
		"/shop/ranking/products?start_date="+TestStartDate+"&end_date="+TestEndDate+"&sort_by=quantity&limit=10", ShopToken, 200)
	TestAPI("Shop Ranking Products - Invalid sort_by", "GET", "/shop/ranking/products?sort_by=invalid", ShopToken, 400)
	TestAPI("Shop Ranking Products - No params", "GET", "/shop/ranking/products", ShopToken, 200)
}

func TestPlatformAPIs() {
	printSection("TESTING PLATFORM APIs")

	printSection("Platform - Nhóm 1: Tổng quan")

	TestAPI("Platform Overview - Full params", "GET",
		"/platform/overview?start_date="+TestStartDate+"&end_date="+TestEndDate, AdminToken, 200)
	TestAPI("Platform Overview - No params", "GET", "/platform/overview", AdminToken, 200)
	TestAPI("Platform Overview - Invalid date", "GET", "/platform/overview?start_date=invalid", AdminToken, 400)
	TestAPI("Platform Overview - Wrong role", "GET", "/platform/overview", ShopToken, 403)

	printSection("Platform - Nhóm 2: Quản lý Đơn hàng")

	TestAPI("List Platform Orders - Full params", "GET",
		"/platform/orders?shop_id="+TestShopID+"&user_id="+TestUserID+"&status=COMPLETED&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20&offset=0", AdminToken, 200)
	TestAPI("List Platform Orders - Filter shop_id", "GET", "/platform/orders?shop_id="+TestShopID, AdminToken, 200)
	TestAPI("List Platform Orders - Filter user_id", "GET", "/platform/orders?user_id="+TestUserID, AdminToken, 200)
	TestAPI("List Platform Orders - Filter status", "GET", "/platform/orders?status=PROCESSING", AdminToken, 200)
	TestAPI("List Platform Orders - No filters", "GET", "/platform/orders", AdminToken, 200)

	statuses := []string{"AWAITING_PAYMENT", "PROCESSING", "SHIPPED", "COMPLETED", "CANCELLED", "REFUNDED"}
	for _, status := range statuses {
		TestAPI(fmt.Sprintf("List Platform Orders - Status: %s", status), "GET",
			"/platform/orders?status="+status, AdminToken, 200)
	}

	TestAPI("Get Enriched Platform Order", "GET", "/platform/orders/"+TestOrderID, AdminToken, 200)

	printSection("Platform - Nhóm 3: Quản lý Tài chính")

	TestAPI("Platform Revenue Timeseries - Full", "GET",
		"/platform/finance/revenue-timeseries?start_date="+TestStartDate+"&end_date="+TestEndDate, AdminToken, 200)
	TestAPI("Platform Revenue Timeseries - No params", "GET", "/platform/finance/revenue-timeseries", AdminToken, 200)
	TestAPI("List Platform Transactions - Full", "GET",
		"/platform/finance/transactions?type=PAYMENT&status=SUCCESS&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20&offset=0", AdminToken, 200)

	txnTypes := []string{"PAYMENT", "REFUND", "PAYOUT", "DEPOSIT"}
	for _, txnType := range txnTypes {
		TestAPI(fmt.Sprintf("List Platform Transactions - Type: %s", txnType), "GET",
			"/platform/finance/transactions?type="+txnType, AdminToken, 200)
	}

	txnStatuses := []string{"PENDING", "SUCCESS", "FAILED"}
	for _, status := range txnStatuses {
		TestAPI(fmt.Sprintf("List Platform Transactions - Status: %s", status), "GET",
			"/platform/finance/transactions?status="+status, AdminToken, 200)
	}

	TestAPI("List Platform Settlements - Full", "GET",
		"/platform/finance/settlements?status=SETTLED&start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20&offset=0", AdminToken, 200)

	settlementStatuses := []string{"PENDING_SETTLEMENT", "FUNDS_HELD", "SETTLED", "FAILED"}
	for _, status := range settlementStatuses {
		TestAPI(fmt.Sprintf("List Platform Settlements - Status: %s", status), "GET",
			"/platform/finance/settlements?status="+status, AdminToken, 200)
	}

	TestAPI("List Platform Ledgers - SHOP", "GET", "/platform/finance/ledgers?owner_type=SHOP&limit=20&offset=0", AdminToken, 200)
	TestAPI("List Platform Ledgers - PLATFORM", "GET", "/platform/finance/ledgers?owner_type=PLATFORM", AdminToken, 200)
	TestAPI("List Platform Ledgers - No filter", "GET", "/platform/finance/ledgers", AdminToken, 200)
	TestAPI("List Ledger Entries", "GET", "/platform/finance/ledgers/"+TestLedgerID+"/entries?limit=50&offset=0", AdminToken, 200)

	printSection("Platform - Nhóm 4: Voucher")

	TestAPI("List Platform Vouchers - SHOP", "GET", "/platform/vouchers?owner_type=SHOP&is_active=true&limit=20&offset=0", AdminToken, 200)
	TestAPI("List Platform Vouchers - PLATFORM", "GET", "/platform/vouchers?owner_type=PLATFORM&is_active=true", AdminToken, 200)
	TestAPI("List Platform Vouchers - Active", "GET", "/platform/vouchers?is_active=true", AdminToken, 200)
	TestAPI("List Platform Vouchers - Inactive", "GET", "/platform/vouchers?is_active=false", AdminToken, 200)
	TestAPI("List Platform Vouchers - No filter", "GET", "/platform/vouchers", AdminToken, 200)
	TestAPI("Platform Voucher Performance - Full", "GET",
		"/platform/vouchers/performance?start_date="+TestStartDate+"&end_date="+TestEndDate, AdminToken, 200)
	TestAPI("Platform Voucher Performance - No params", "GET", "/platform/vouchers/performance", AdminToken, 200)

	printSection("Platform - Nhóm 5: Phân tích Shop")

	TestAPI("List Platform Shops - With pagination", "GET", "/platform/shops?limit=20&offset=0", AdminToken, 200)
	TestAPI("List Platform Shops - No pagination", "GET", "/platform/shops", AdminToken, 200)
	TestAPI("Platform Shop Detail - Full", "GET",
		"/platform/shops/"+TestShopID+"/detail?start_date="+TestStartDate+"&end_date="+TestEndDate, AdminToken, 200)
	TestAPI("Platform Shop Detail - No params", "GET", "/platform/shops/"+TestShopID+"/detail", AdminToken, 200)

	printSection("Platform - Nhóm 6: Xếp hạng")

	TestAPI("Platform Ranking Shops - Full", "GET",
		"/platform/ranking/shops?start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=10", AdminToken, 200)
	TestAPI("Platform Ranking Shops - No params", "GET", "/platform/ranking/shops", AdminToken, 200)
	TestAPI("Platform Ranking Products - Full", "GET",
		"/platform/ranking/products?start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=20", AdminToken, 200)
	TestAPI("Platform Ranking Products - No params", "GET", "/platform/ranking/products", AdminToken, 200)
	TestAPI("Platform Ranking Users - Full", "GET",
		"/platform/ranking/users?start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=15", AdminToken, 200)
	TestAPI("Platform Ranking Users - No params", "GET", "/platform/ranking/users", AdminToken, 200)
	TestAPI("Platform Ranking Categories - Full", "GET",
		"/platform/ranking/categories?start_date="+TestStartDate+"&end_date="+TestEndDate+"&limit=10", AdminToken, 200)
	TestAPI("Platform Ranking Categories - No params", "GET", "/platform/ranking/categories", AdminToken, 200)
}

func TestEdgeCases() {
	printSection("TESTING EDGE CASES")

	TestAPI("Edge Case - Same start and end date", "GET",
		"/shop/overview?start_date="+TestSingleDate+"&end_date="+TestSingleDate, ShopToken, 200)
	TestAPI("Edge Case - End before start", "GET",
		"/shop/overview?start_date="+TestEndDate+"&end_date="+TestStartDate, ShopToken, 200)
	TestAPI("Edge Case - Large date range", "GET",
		"/shop/overview?start_date=2024-01-01&end_date=2026-12-31", ShopToken, 200)
	TestAPI("Edge Case - Zero limit", "GET", "/shop/orders?limit=0", ShopToken, 200)
	TestAPI("Edge Case - Large limit", "GET", "/shop/orders?limit=10000", ShopToken, 200)
	TestAPI("Edge Case - Large offset", "GET", "/shop/orders?offset=999999", ShopToken, 200)
}

func main() {
	printSection("E-COMMERCE ANALYTICS API TEST SUITE")
	fmt.Printf("Base URL: %s\n", BaseURL)
	fmt.Printf("Started at: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	if strings.Contains(ShopToken, "YOUR_") || strings.Contains(AdminToken, "YOUR_") {
		fmt.Printf("%sWarning: Using placeholder tokens. Update tokens in api_test.go%s\n", ColorYellow, ColorReset)
	}

	testSuite := "all"
	if len(os.Args) > 1 {
		testSuite = os.Args[1]
	}

	switch testSuite {
	case "shop":
		TestShopAPIs()
	case "platform":
		TestPlatformAPIs()
	case "edge":
		TestEdgeCases()
	case "all":
		TestShopAPIs()
		TestPlatformAPIs()
		TestEdgeCases()
	default:
		fmt.Printf("%sInvalid test suite: %s%s\n", ColorRed, testSuite, ColorReset)
		fmt.Println("Usage: go run api_test.go [all|shop|platform|edge]")
		os.Exit(1)
	}

	printSummary()

	if failedTests > 0 {
		os.Exit(1)
	}
}
