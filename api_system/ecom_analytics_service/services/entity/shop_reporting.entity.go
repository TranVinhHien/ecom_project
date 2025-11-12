package services

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	// "bytes"
	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
)

// ShopOverviewResponse là DTO cho API GET /api/v1/shop/overview
type ShopOverviewResponse struct {
	TotalGMV         float64 `json:"total_gmv"`         // Doanh thu gộp (từ shop_orders)
	TotalNetRevenue  float64 `json:"total_net_revenue"` // Doanh thu thuần (đã đối soát)
	TotalOrders      int64   `json:"total_orders"`
	ProcessingOrders int64   `json:"processing_orders"`
	WalletBalance    float64 `json:"wallet_balance"`  // Số dư khả dụng
	PendingBalance   float64 `json:"pending_balance"` // Số dư chờ
}

// ShopWalletSummaryResponse là DTO cho API GET /api/v1/shop/wallet/summary
type ShopWalletSummaryResponse struct {
	Balance             float64 `json:"balance"`               // Số dư khả dụng
	PendingBalance      float64 `json:"pending_balance"`       // Số dư chờ (từ ví)
	TotalSettledRevenue float64 `json:"total_settled_revenue"` // Tổng tiền đã đối soát (từ settlements)
	TotalFundsHeld      float64 `json:"total_funds_held"`      // Tổng tiền đang tạm giữ (từ settlements)
	TotalWithdrawn      float64 `json:"total_withdrawn"`       // Tổng tiền đã rút (từ ledger_entries)
}

// === Nhóm 2 ===

type ListShopOrdersParams struct {
	ShopID    string
	Status    sql.NullString // 'PROCESSING', 'COMPLETED', v.v.
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

type ShopOrderWithPaymentStatus struct {
	ShopOrder db_order.ShopOrders
	// PaymentStatus string `json:"payment_status"` // 'PENDING', 'SUCCESS', 'FAILED'
}

type ListShopOrdersResponse struct {
	Orders []ShopOrderWithPaymentStatus `json:"orders"`
	// (Có thể thêm metadata phân trang tại đây)
}

type EnrichedShopOrderResponse struct {
	OrderDetails    db_order.ShopOrders                  `json:"order_details"`
	OrderItems      []db_order.OrderItems                `json:"order_items"`
	SettlementInfo  *db_transaction.ShopOrderSettlements `json:"settlement_info"` // Dùng con trỏ
	PaymentStatus   string                               `json:"payment_status"`
	ParentOrderInfo db_order.Orders                      `json:"parent_order_info"`
}

type ListShopOrderItemsParams struct {
	ShopID    string
	ProductID sql.NullString // Lọc theo product_id (tùy chọn)
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

// === Nhóm 3 ===

type RevenueDatapoint struct {
	Date       string  `json:"date"`
	GMV        float64 `json:"gmv"`         // Doanh thu gộp
	NetRevenue float64 `json:"net_revenue"` // Doanh thu thuần (đã đối soát)
}

type RevenueTimeseriesResponse struct {
	Data []RevenueDatapoint `json:"data"`
}

type ListWalletLedgerEntriesParams struct {
	ShopID string
	Limit  int32
	Offset int32
}

type ListShopSettlementsParams struct {
	ShopID    string
	Status    sql.NullString // 'PENDING_SETTLEMENT', 'SETTLED', v.v.
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

// === Nhóm 4 ===

type ListShopVouchersParams struct {
	ShopID   string
	IsActive sql.NullBool // Lọc theo trạng thái active
	Limit    int32
	Offset   int32
}

type VoucherPerformanceResponse struct {
	TotalUsageCount    int64   `json:"total_usage_count"`
	TotalDiscountValue float64 `json:"total_discount_value"`
}

type ListVoucherUsageDetailsParams struct {
	VoucherID string
	ShopID    string // Để đảm bảo shop chỉ xem voucher của mình
	Limit     int32
	Offset    int32
}

// === Nhóm 5 ===

type ShopRankingProductsParams struct {
	ShopID    string
	StartDate time.Time
	EndDate   time.Time
	SortBy    string // "revenue" hoặc "quantity"
	Limit     int32
}

type ProductRankingRow struct {
	ProductID             string  `json:"product_id"`
	SkuID                 string  `json:"sku_id"`
	ProductNameSnapshot   string  `json:"product_name_snapshot"`
	SkuAttributesSnapshot string  `json:"sku_attributes_snapshot"`
	TotalRevenue          float64 `json:"total_revenue"`  // SUM(total_price)
	TotalQuantity         int64   `json:"total_quantity"` // SUM(quantity)
}

// --- Helpers ---
// parseAmountToFloat64 chuyển đổi kiểu DECIMAL(15,2) (thường là sql.NullString) sang float64

// ParseAmountToFloat64 chuyển đổi mọi kiểu dữ liệu từ database sang float64
// Hỗ trợ: sql.NullString, string, []byte, []uint8, float64, float32, int*, nil
func ParseAmountToFloat64(amount interface{}) (float64, error) {
	// Xử lý nil ngay từ đầu
	if amount == nil {
		return 0.0, nil
	}

	// Sử dụng reflect để xử lý động
	val := reflect.ValueOf(amount)

	// Xử lý theo Kind (kiểu cơ bản)
	switch val.Kind() {
	case reflect.String:
		// Trường hợp: string
		str := val.String()
		if str == "" {
			return 0.0, nil
		}
		result, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Printf("[ParseAmountToFloat64] Failed to parse string '%s': %v", str, err)
			return 0.0, fmt.Errorf("failed to parse string '%s': %w", str, err)
		}
		return result, nil

	case reflect.Slice:
		// Trường hợp: []byte hoặc []uint8
		if val.Type().Elem().Kind() == reflect.Uint8 {
			bytes := val.Bytes()
			if len(bytes) == 0 {
				return 0.0, nil
			}
			str := string(bytes)
			result, err := strconv.ParseFloat(str, 64)
			if err != nil {
				log.Printf("[ParseAmountToFloat64] Failed to parse []byte '%s' (raw: %v): %v", str, bytes, err)
				return 0.0, fmt.Errorf("failed to parse []byte '%s': %w", str, err)
			}
			return result, nil
		}

	case reflect.Float64:
		return val.Float(), nil

	case reflect.Float32:
		return float64(val.Float()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	}

	// Xử lý các kiểu đặc biệt từ database/sql
	switch v := amount.(type) {
	case sql.NullString:
		if !v.Valid {
			return 0.0, nil
		}
		result, err := strconv.ParseFloat(v.String, 64)
		if err != nil {
			log.Printf("[ParseAmountToFloat64] Failed to parse sql.NullString '%s': %v", v.String, err)
			return 0.0, fmt.Errorf("failed to parse sql.NullString '%s': %w", v.String, err)
		}
		return result, nil

	case sql.NullFloat64:
		if !v.Valid {
			return 0.0, nil
		}
		return v.Float64, nil

	case sql.NullInt64:
		if !v.Valid {
			return 0.0, nil
		}
		return float64(v.Int64), nil
	}

	// Nếu không match case nào
	log.Printf("[ParseAmountToFloat64] Unsupported type: %T (value: %v, kind: %v)", amount, amount, val.Kind())
	return 0.0, fmt.Errorf("unsupported type: %T (value: %v)", amount, amount)
}

// ParseInterfaceToInt tương tự, xử lý chuyển đổi sang int64
func ParseInterfaceToInt(amount interface{}) (int64, error) {
	// Xử lý nil ngay từ đầu
	if amount == nil {
		return 0, nil
	}

	// Sử dụng reflect để xử lý động
	val := reflect.ValueOf(amount)

	// Xử lý theo Kind
	switch val.Kind() {
	case reflect.String:
		str := val.String()
		if str == "" {
			return 0, nil
		}
		result, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			log.Printf("[ParseInterfaceToInt] Failed to parse string '%s': %v", str, err)
			return 0, fmt.Errorf("failed to parse string '%s': %w", str, err)
		}
		return result, nil

	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			bytes := val.Bytes()
			if len(bytes) == 0 {
				return 0, nil
			}
			str := string(bytes)
			result, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Printf("[ParseInterfaceToInt] Failed to parse []byte '%s' (raw: %v): %v", str, bytes, err)
				return 0, fmt.Errorf("failed to parse []byte '%s': %w", str, err)
			}
			return result, nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint()), nil

	case reflect.Float64, reflect.Float32:
		return int64(val.Float()), nil
	}

	// Xử lý các kiểu đặc biệt
	switch v := amount.(type) {
	case sql.NullString:
		if !v.Valid {
			return 0, nil
		}
		result, err := strconv.ParseInt(v.String, 10, 64)
		if err != nil {
			log.Printf("[ParseInterfaceToInt] Failed to parse sql.NullString '%s': %v", v.String, err)
			return 0, fmt.Errorf("failed to parse sql.NullString '%s': %w", v.String, err)
		}
		return result, nil

	case sql.NullInt64:
		if !v.Valid {
			return 0, nil
		}
		return v.Int64, nil

	case sql.NullFloat64:
		if !v.Valid {
			return 0, nil
		}
		return int64(v.Float64), nil
	}

	// Nếu không match case nào
	log.Printf("[ParseInterfaceToInt] Unsupported type: %T (value: %v, kind: %v)", amount, amount, val.Kind())
	return 0, fmt.Errorf("unsupported type: %T (value: %v)", amount, amount)
}
