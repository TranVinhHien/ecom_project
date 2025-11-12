package controllers

import (
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	"github.com/TranVinhHien/ecom_order_service/services"

	"github.com/gin-gonic/gin"
)

type apiController struct {
	service services.ServiceUseCase
	jwt     token.Maker
}

func NewAPIController(s services.ServiceUseCase, jwt token.Maker) apiController {
	return apiController{service: s, jwt: jwt}
}

func (api apiController) SetUpRoute(group *gin.RouterGroup) {
	group.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	// =================================================================
	// CUSTOMER ENDPOINTS - Endpoints cho khách hàng
	// =================================================================
	orders := group.Group("/orders")
	{
		orders.PUT("/callback_payment_online/:order_id", api.callbackPaymentOnline())
	}

	orders_auth := orders.Use(authorization(api.jwt))
	{
		// POST /api/v1/orders - Tạo đơn hàng mới
		orders_auth.POST("", api.createOrder())

		// GET /api/v1/orders - Lấy danh sách đơn hàng của user
		orders_auth.GET("", api.listUserOrders())

		// GET /api/v1/orders/{orderCode} - Lấy chi tiết đơn hàng
		orders_auth.GET("/:orderCode", api.getOrderDetail())

		// GET /api/v1/orders/search - Tìm kiếm đơn hàng chi tiết với bộ lọc
		orders_auth.GET("/search/detail", api.searchOrdersDetail())
	}

	// =================================================================
	// ADMIN/SHOP ENDPOINTS - Endpoints cho shop owner/admin
	// =================================================================
	//

	admin := orders.Group("/admin").Use(authorization(api.jwt))
	{
		// // GET /api/v1/admin/shop-orders - Lấy danh sách đơn hàng của shop
		// admin.GET("/shop-orders", api.listShopOrders())

		// // POST /api/v1/admin/shop-orders/{shopOrderCode}/ship - Đánh dấu đơn hàng đã ship
		// admin.POST("/shop-orders/:shopOrderCode/ship", api.shipShopOrder())

		// PUT /api/v1/admin/shop-orders/{shopOrderCode}/status - Cập nhật trạng thái đơn hàng
		admin.PUT("/update_status", api.updateShopOrderStatus())
	}

	// =================================================================
	// PRODUCT ENDPOINTS (các endpoints product đã có từ trước)
	// =================================================================
	// product := group.Group("/product")
	// {
	// 	product.GET("/getall", api.getAllProductSimple())
	// 	product.GET("/getdetail/:id", api.getDetailProduct())
	// 	product_auth := product.Group("").Use(authorization(api.jwt))
	// 	{
	// 		product_auth.POST("/create", api.createProduct())
	// 		product_auth.PUT("/update/:id", api.updateProduct())
	// 		product_auth.POST("/update_sku_reserver", api.updateSKUReserverProduct())
	// 	}
	// }

	// =================================================================
	// VOUCHER ENDPOINTS
	// =================================================================
	vouchers := group.Group("/vouchers")
	{
		vouchers_auth := vouchers.Use(authorization(api.jwt))
		{
			// GET /api/v1/vouchers - list vouchers available to current user
			vouchers_auth.GET("", api.listVouchersForUser())
			// POST /api/v1/vouchers - create a voucher (admin/owner)
			vouchers_auth.POST("", api.createVoucher())
			// PUT /api/v1/vouchers/:voucherID - update a voucher
			vouchers_auth.PUT("/:voucherID", api.updateVoucher())
		}
	}
}
