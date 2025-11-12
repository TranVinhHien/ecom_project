package services

import (
	"database/sql"
	"time"

	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
)

// === Nhóm 1: Tổng quan ===

type PlatformOverviewResponse struct {
	TotalGMV             float64 `json:"total_gmv"`
	TotalPlatformRevenue float64 `json:"total_platform_revenue"`
	TotalPlatformCost    float64 `json:"total_platform_cost"`
	PlatformProfit       float64 `json:"platform_profit"`
	TotalOrders          int64   `json:"total_orders"`
	TotalShops           int64   `json:"total_shops"`
}

// === Nhóm 2: Quản lý Đơn hàng ===

type ListPlatformOrdersParams struct {
	ShopID    sql.NullString
	UserID    sql.NullString // (Lưu ý: sqlc query hiện tại chưa hỗ trợ filter này)
	Status    sql.NullString
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

// (Sử dụng lại ShopOrderWithPaymentStatus từ dtos.go)
type ListPlatformOrdersResponse struct {
	Orders []ShopOrderWithPaymentStatus `json:"orders"`
	// (Có thể thêm metadata phân trang)
}

type EnrichedPlatformOrderResponse struct {
	ParentOrder db_order.Orders       `json:"parent_order"`
	ShopOrders  []db_order.ShopOrders `json:"shop_orders"`
}

// === Nhóm 3: Quản lý Tài chính ===

type PlatformRevenueDatapoint struct {
	Date            string  `json:"date"`
	TotalGMV        float64 `json:"total_gmv"`
	PlatformRevenue float64 `json:"platform_revenue"`
	PlatformCost    float64 `json:"platform_cost"`
	PlatformProfit  float64 `json:"platform_profit"`
}

type PlatformRevenueTimeseriesResponse struct {
	Data []PlatformRevenueDatapoint `json:"data"`
}

type ListPlatformTransactionsParams struct {
	Type      sql.NullString
	Status    sql.NullString
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

type ListPlatformSettlementsParams struct {
	Status    sql.NullString
	StartDate sql.NullTime
	EndDate   sql.NullTime
	Limit     int32
	Offset    int32
}

type ListPlatformLedgersParams struct {
	OwnerType sql.NullString // 'SHOP' hoặc 'PLATFORM'
	Limit     int32
	Offset    int32
}

// === Nhóm 4: Phân tích Voucher ===

type ListPlatformVouchersParams struct {
	OwnerType sql.NullString // 'SHOP' hoặc 'PLATFORM'
	IsActive  sql.NullBool
	Limit     int32
	Offset    int32
}

type PlatformVoucherPerformanceResponse struct {
	UsageHistoryStats db_order.GetPlatformVoucherPerformanceRow `json:"usage_history_stats"`
	PlatformCostStats db_transaction.GetPlatformCostSummaryRow  `json:"platform_cost_stats"`
	TotalVoucherCost  float64                                   `json:"total_voucher_cost"` // Tổng chi phí voucher Sàn chịu
}

// === Nhóm 5: Phân tích Shop ===

type ListPlatformShopsParams struct {
	Limit  int32
	Offset int32
	// (Có thể thêm filter theo tên shop, v.v. nếu user_service hỗ trợ)
}

// (Dùng lại db_order.GetPlatformTopShopsByGMVRow làm DTO)
type PlatformShopRow struct {
	db_order.GetPlatformTopShopsByGMVRow
	// (Có thể bổ sung ShopName, ShopAvatar... nếu gọi User service)
}

type PlatformShopDetailResponse struct {
	Overview *ShopOverviewResponse      `json:"overview"`
	Wallet   *ShopWalletSummaryResponse `json:"wallet"`
	// (Thêm các thông tin profile shop nếu gọi User service)
}

// === Nhóm 6: Xếp hạng ===

type PlatformRankingParams struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int32
}

type PlatformRankingCategoriesResponse struct {
	Message string `json:"message"`
	// (Dữ liệu sẽ được thêm nếu có kết nối product_db)
}
