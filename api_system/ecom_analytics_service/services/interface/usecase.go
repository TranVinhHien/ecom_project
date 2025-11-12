package iservices

import (
	"context"
	"time"

	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
	assets_services "github.com/TranVinhHien/ecom_analytics_service/services/assets"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"
)

// ServiceUseCase định nghĩa tất cả các nghiệp vụ của service
type ServiceSITEUseCase interface {
	// === NHÓM I: API CHO SHOP ===
	// === Nhóm 1: Thống kê Tổng quan Shop ===
	GetShopOverview(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.ShopOverviewResponse, *assets_services.ServiceError)
	GetShopWalletSummary(ctx context.Context, shopID string) (*entity.ShopWalletSummaryResponse, *assets_services.ServiceError)
	ListShopOrders(ctx context.Context, params entity.ListShopOrdersParams) (*entity.ListShopOrdersResponse, *assets_services.ServiceError)
	GetEnrichedShopOrder(ctx context.Context, shopID, shopOrderID string) (*entity.EnrichedShopOrderResponse, *assets_services.ServiceError)

	// === Nhóm 2: Phân tích Đơn hàng ===
	ListShopOrderItems(ctx context.Context, params entity.ListShopOrderItemsParams) ([]db_order.OrderItems, *assets_services.ServiceError)

	// === Nhóm 3: Phân tích Doanh thu & Dòng tiền ===
	GetShopRevenueTimeseries(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.RevenueTimeseriesResponse, *assets_services.ServiceError)
	ListShopWalletLedgerEntries(ctx context.Context, params entity.ListWalletLedgerEntriesParams) ([]db_transaction.LedgerEntries, *assets_services.ServiceError)
	ListShopSettlements(ctx context.Context, params entity.ListShopSettlementsParams) ([]db_transaction.ShopOrderSettlements, *assets_services.ServiceError)

	// === Nhóm 4: Phân tích Voucher ===
	ListShopVouchers(ctx context.Context, params entity.ListShopVouchersParams) ([]db_order.Vouchers, *assets_services.ServiceError)
	GetShopVoucherPerformance(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.VoucherPerformanceResponse, *assets_services.ServiceError)
	GetShopVoucherUsageDetails(ctx context.Context, params entity.ListVoucherUsageDetailsParams) ([]db_order.VoucherUsageHistory, *assets_services.ServiceError)

	// === Nhóm 5: Xếp hạng ===
	GetShopRankingProducts(ctx context.Context, params entity.ShopRankingProductsParams) ([]entity.ProductRankingRow, *assets_services.ServiceError)

	// === NHÓM II: API CHO SÀN ===

	// Nhóm 1: Tổng quan
	GetPlatformOverview(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformOverviewResponse, *assets_services.ServiceError)

	// Nhóm 2: Quản lý Đơn hàng
	ListPlatformOrders(ctx context.Context, params entity.ListPlatformOrdersParams) (*entity.ListPlatformOrdersResponse, *assets_services.ServiceError)
	GetEnrichedPlatformOrder(ctx context.Context, orderID string) (*entity.EnrichedPlatformOrderResponse, *assets_services.ServiceError)

	// Nhóm 3: Quản lý Tài chính
	GetPlatformRevenueTimeseries(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformRevenueTimeseriesResponse, *assets_services.ServiceError)
	ListPlatformTransactions(ctx context.Context, params entity.ListPlatformTransactionsParams) ([]db_transaction.Transactions, *assets_services.ServiceError)
	ListPlatformSettlements(ctx context.Context, params entity.ListPlatformSettlementsParams) ([]db_transaction.ShopOrderSettlements, *assets_services.ServiceError)
	ListPlatformLedgers(ctx context.Context, params entity.ListPlatformLedgersParams) ([]db_transaction.AccountLedgers, *assets_services.ServiceError)
	ListLedgerEntries(ctx context.Context, ledgerID string, limit, offset int32) ([]db_transaction.LedgerEntries, *assets_services.ServiceError)

	// Nhóm 4: Phân tích Voucher
	ListPlatformVouchers(ctx context.Context, params entity.ListPlatformVouchersParams) ([]db_order.Vouchers, *assets_services.ServiceError)
	GetPlatformVoucherPerformance(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformVoucherPerformanceResponse, *assets_services.ServiceError)

	// Nhóm 5: Phân tích Shop
	ListPlatformShops(ctx context.Context, params entity.ListPlatformShopsParams) ([]entity.PlatformShopRow, *assets_services.ServiceError)
	GetPlatformShopDetail(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.PlatformShopDetailResponse, *assets_services.ServiceError)

	// Nhóm 6: Xếp hạng
	GetPlatformRankingShops(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopShopsByGMVRow, *assets_services.ServiceError)
	GetPlatformRankingProducts(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopProductsByQuantityRow, *assets_services.ServiceError)
	GetPlatformRankingUsers(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopUsersBySpendRow, *assets_services.ServiceError)
	GetPlatformRankingCategories(ctx context.Context, params entity.PlatformRankingParams) (*entity.PlatformRankingCategoriesResponse, *assets_services.ServiceError)
}

type FeedbackUseCase interface {
	// Message Ratings
	SubmitMessageRating(ctx context.Context, req *entity.SubmitMessageRatingRequest) *assets_services.ServiceError
	GetMessageRatingStats(ctx context.Context, startDate, endDate *string) (*entity.MessageRatingStatsResponse, *assets_services.ServiceError)
	GetMessageRatingsTimeSeries(ctx context.Context, startDate, endDate *string) ([]entity.MessageRatingTimeSeriesItem, *assets_services.ServiceError)
	GetMessageRatingsList(ctx context.Context, req *entity.GetMessageRatingsRequest) ([]entity.MessageRatingDetailItem, *assets_services.ServiceError)

	// Customer Feedback
	SubmitCustomerFeedback(ctx context.Context, req *entity.SubmitCustomerFeedbackRequest) (string, *assets_services.ServiceError)
	GetCustomerFeedbacks(ctx context.Context, req *entity.GetCustomerFeedbacksRequest) ([]entity.CustomerFeedbackItem, *assets_services.ServiceError)
	GetCustomerFeedbackByID(ctx context.Context, id string) (*entity.CustomerFeedbackItem, *assets_services.ServiceError)
	GetCustomerFeedbackStats(ctx context.Context, startDate, endDate *string) (*entity.CustomerFeedbackStatsResponse, *assets_services.ServiceError)
	GetCustomerFeedbacksByCategory(ctx context.Context, startDate, endDate *string) ([]entity.CustomerFeedbackCategoryStats, *assets_services.ServiceError)
}
