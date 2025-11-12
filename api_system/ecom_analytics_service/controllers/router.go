package controllers

import (
	"github.com/TranVinhHien/ecom_analytics_service/assets/token"
	"github.com/TranVinhHien/ecom_analytics_service/services"

	"github.com/gin-gonic/gin"
)

type apiController struct {
	service services.ServiceUseCase
	jwt     token.Maker
}

func NewAPIController(s services.ServiceUseCase, jwt token.Maker) apiController {
	return apiController{
		service: s,
		jwt:     jwt,
	}
}

// SetUpRoute định nghĩa tất cả các đường dẫn cho Service Thống Kê
func (api apiController) SetUpRoute(group *gin.RouterGroup) {
	// === Cấu hình CORS chung ===
	group.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	// === Nhóm I: API cho Nhà bán hàng (SHOP) ===
	// Tất cả API trong nhóm này yêu cầu đăng nhập và có vai trò "SHOP"
	shop := group.Group("/shop").
		Use(authorization(api.jwt)).
		Use(checkRole("SELLER")).
		Use(getShopID()) // Middleware để lấy shop_id và lưu vào context
	{
		// Nhóm 1: Tổng quan
		shop.GET("/overview", api.getShopOverview())
		shop.GET("/wallet/summary", api.getShopWalletSummary())

		// Nhóm 2: Phân tích Đơn hàng
		shop.GET("/orders", api.listShopOrders())
		shop.GET("/orders/:shop_order_id", api.getEnrichedShopOrder()) // Lấy chi tiết đơn
		shop.GET("/order-items", api.listShopOrderItems())

		// Nhóm 3: Phân tích Doanh thu & Dòng tiền
		shop.GET("/revenue/timeseries", api.getShopRevenueTimeseries())
		shop.GET("/wallet/ledger-entries", api.listShopWalletLedgerEntries())
		shop.GET("/settlements", api.listShopSettlements())

		// Nhóm 4: Phân tích Voucher
		shop.GET("/vouchers", api.listShopVouchers())
		shop.GET("/vouchers/performance", api.getShopVoucherPerformance())
		shop.GET("/vouchers/:voucher_id/details", api.getShopVoucherUsageDetails())

		// Nhóm 5: Xếp hạng
		shop.GET("/ranking/products", api.getShopRankingProducts())
	}

	// === Nhóm II: API cho Nền tảng (PLATFORM / ADMIN) ===
	// Tất cả API trong nhóm này yêu cầu đăng nhập và có vai trò "ADMIN"
	platform := group.Group("/platform").
		Use(authorization(api.jwt)).
		Use(checkRole("ADMIN")) // Sử dụng middleware từ ví dụ của bạn
	{
		// Nhóm 1: Tổng quan
		platform.GET("/overview", api.getPlatformOverview())

		// Nhóm 2: Quản lý Đơn hàng
		platform.GET("/orders", api.listPlatformOrders())
		platform.GET("/orders/:order_id", api.getEnrichedPlatformOrder()) // Lấy chi tiết đơn TỔNG

		// Nhóm 3: Quản lý Tài chính
		platform.GET("/finance/revenue-timeseries", api.getPlatformRevenueTimeseries())
		platform.GET("/finance/transactions", api.listPlatformTransactions())
		platform.GET("/finance/settlements", api.listPlatformSettlements())
		platform.GET("/finance/ledgers", api.listPlatformLedgers())
		platform.GET("/finance/ledgers/:ledger_id/entries", api.listLedgerEntries())

		// Nhóm 4: Phân tích Voucher
		platform.GET("/vouchers", api.listPlatformVouchers())
		platform.GET("/vouchers/performance", api.getPlatformVoucherPerformance())

		// Nhóm 5: Phân tích Shop
		platform.GET("/shops", api.listPlatformShops())
		platform.GET("/shops/:shop_id/detail", api.getPlatformShopDetail()) // Xem chi tiết 1 shop

		// Nhóm 6: Xếp hạng Toàn Sàn
		platform.GET("/ranking/shops", api.getPlatformRankingShops())
		platform.GET("/ranking/products", api.getPlatformRankingProducts())
		platform.GET("/ranking/users", api.getPlatformRankingUsers())
		platform.GET("/ranking/categories", api.getPlatformRankingCategories())

		// Nhóm 7: Thống kê Tương tác với Chatbox
		platform.GET("/chatbox/statistics", api.getChatboxStatistics())
		platform.GET("/chatbox/reviews", api.getChatboxReviews())

		// Nhóm 8: Thống kê Hỗ trợ Khách hàng
		platform.GET("/customer-support/statistics", api.getCustomerSupportStatistics())
		platform.GET("/customer-support/feedbacks", api.listCustomerFeedbacks())
		platform.GET("/customer-support/feedbacks/:id", api.getCustomerFeedbackDetail())

	}

	// === Nhóm III: API Công khai (Public) ===
	public := group.Group("/public").Use(authorization(api.jwt))
	{
		// Nhóm 1 : Gửi đánh giá của chatbox
		public.POST("/chatbox/review", api.submitChatboxReview())
		// Nhóm 2 : Gửi Khiếu nại tới Hỗ trợ Khách hàng
		public.POST("/customer-support/complaint", api.submitCustomerSupportComplaint())
	}
}
