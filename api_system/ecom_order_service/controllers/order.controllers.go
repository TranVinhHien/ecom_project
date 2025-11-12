package controllers

import (
	"net/http"
	"strconv"
	"time"

	assets_api "github.com/TranVinhHien/ecom_order_service/assets/api"
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"

	"github.com/gin-gonic/gin"
)

// =================================================================
// CUSTOMER ENDPOINTS - /api/v1/orders
// =================================================================

// createOrder tạo đơn hàng mới từ giỏ hàng
func (api *apiController) createOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		token := ctx.MustGet("token").(string)
		// code chả có logic gì ở đây
		// phải đổi thành tham số ở phần controllers
		var req services.CreateOrderRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		// Gọi service để tạo đơn hàng
		result, err := api.service.CreateOrder(ctx, authPayload.Sub, token, req)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusCreated, assets_api.SimpSuccessResponse("Order created successfully", result))
	}
}

// listUserOrders lấy danh sách đơn hàng của user hiện tại
func (api *apiController) listUserOrders() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		status := ctx.Query("status") //
		// if status == "" {
		// 	ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Status query parameter is required"))
		// 	return
		// }
		if status != "" {
			if status != "AWAITING_PAYMENT" && status != "PROCESSING" && status != "SHIPPED" && status != "COMPLETED" && status != "CANCELED" && status != "REFUNDED" {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid status value"))
				return
			}
		}

		result, err := api.service.ListUserOrders(ctx, authPayload.Sub, services.NewQueryFilter(page, limit, nil, nil), status)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Get orders successfully", result))
	}
}

// getOrderDetail lấy chi tiết đầy đủ của một đơn hàng
func (api *apiController) getOrderDetail() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		orderCode := ctx.Param("orderCode")

		if orderCode == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Order code is required"))
			return
		}

		result, err := api.service.GetOrderDetail(ctx, authPayload.Sub, orderCode)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Get order detail successfully", result))
	}
}

// searchOrdersDetail tìm kiếm danh sách đơn hàng chi tiết với các bộ lọc
func (api *apiController) searchOrdersDetail() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		// Helper function để parse time từ string
		parseTime := func(timeStr string) (*time.Time, error) {
			if timeStr == "" {
				return nil, nil
			}
			// Try format YYYY-MM-DD first
			t, err := time.Parse("2006-01-02", timeStr)
			if err != nil {
				// Try format YYYY-MM-DD HH:MM:SS
				t, err = time.Parse("2006-01-02 15:04:05", timeStr)
				if err != nil {
					return nil, err
				}
			}
			return &t, nil
		}

		// Parse basic fields (không bao gồm time)
		var filter services.ShopOrderSearchFilter
		
		// Parse basic parameters
		if status := ctx.Query("status"); status != "" {
			filter.Status = &status
		}
		if shopID := ctx.Query("shop_id"); shopID != "" {
			filter.ShopID = &shopID
		}
		if minAmount := ctx.Query("min_amount"); minAmount != "" {
			if val, err := strconv.ParseFloat(minAmount, 64); err == nil {
				filter.MinAmount = &val
			}
		}
		if maxAmount := ctx.Query("max_amount"); maxAmount != "" {
			if val, err := strconv.ParseFloat(maxAmount, 64); err == nil {
				filter.MaxAmount = &val
			}
		}

		// Parse pagination
		if page := ctx.Query("page"); page != "" {
			if val, err := strconv.Atoi(page); err == nil {
				filter.Page = val
			}
		}
		if pageSize := ctx.Query("page_size"); pageSize != "" {
			if val, err := strconv.Atoi(pageSize); err == nil {
				filter.PageSize = val
			}
		}
		
		filter.SortBy = ctx.Query("sort_by")

		// Parse time fields
		var err error
		if filter.CreatedFrom, err = parseTime(ctx.Query("created_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid created_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.CreatedTo, err = parseTime(ctx.Query("created_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid created_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.PaidFrom, err = parseTime(ctx.Query("paid_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid paid_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.PaidTo, err = parseTime(ctx.Query("paid_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid paid_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.ProcessingFrom, err = parseTime(ctx.Query("processing_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid processing_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.ProcessingTo, err = parseTime(ctx.Query("processing_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid processing_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.ShippedFrom, err = parseTime(ctx.Query("shipped_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid shipped_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.ShippedTo, err = parseTime(ctx.Query("shipped_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid shipped_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.CompletedFrom, err = parseTime(ctx.Query("completed_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid completed_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.CompletedTo, err = parseTime(ctx.Query("completed_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid completed_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.CancelledFrom, err = parseTime(ctx.Query("cancelled_from")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid cancelled_from format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}
		if filter.CancelledTo, err = parseTime(ctx.Query("cancelled_to")); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid cancelled_to format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS"))
			return
		}

		// Set default values nếu không được truyền vào
		if filter.Page == 0 {
			filter.Page = 1
		}
		if filter.PageSize == 0 {
			filter.PageSize = 10
		}
		if filter.SortBy == "" {
			filter.SortBy = "created_at"
		}

		// Validate sort_by values nếu có truyền vào
		validSortBy := map[string]bool{
			"created_at":    true,
			"total_amount":  true,
			"paid_at":       true,
			"processing_at": true,
			"shipped_at":    true,
			"completed_at":  true,
		}
		if !validSortBy[filter.SortBy] {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid sort_by value. Allowed: created_at, total_amount, paid_at, processing_at, shipped_at, completed_at"))
			return
		}

		// Validate status nếu có truyền vào
		if filter.Status != nil {
			validStatus := map[string]bool{
				"PENDING":               true,
				"AWAITING_PAYMENT":      true,
				"AWAITING_CONFIRMATION": true,
				"PROCESSING":            true,
				"SHIPPED":               true,
				"COMPLETED":             true,
				"CANCELLED":             true,
				"REFUNDED":              true,
			}
			if !validStatus[*filter.Status] {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid status value"))
				return
			}
		}

		// Call service
		result, errSvc := api.service.SearchOrdersDetail(ctx, authPayload.Sub, filter)
		if errSvc != nil {
			ctx.JSON(errSvc.Code, assets_api.ResponseError(errSvc.Code, errSvc.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Search orders successfully", result))
	}
}

// =================================================================
// ADMIN/SHOP ENDPOINTS - /api/v1/admin/shop-orders
// =================================================================

// listShopOrders lấy danh sách đơn hàng cho shop owner
func (api *apiController) listShopOrders() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		// Giả định shopID lấy từ token hoặc query param
		// Trong production, bạn có thể lấy từ user profile
		shopID := ctx.Query("shop_id")
		if shopID == "" {
			// Fallback: sử dụng userID làm shopID nếu không có
			shopID = authPayload.Sub
		}

		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

		var status *string
		if statusParam := ctx.Query("status"); statusParam != "" {
			status = &statusParam
		}

		var dateFrom, dateTo *string
		if df := ctx.Query("date_from"); df != "" {
			dateFrom = &df
		}
		if dt := ctx.Query("date_to"); dt != "" {
			dateTo = &dt
		}

		result, err := api.service.ListShopOrders(ctx, shopID, status, page, limit, dateFrom, dateTo)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Get shop orders successfully", result))
	}
}

// shipShopOrder đánh dấu đơn hàng của shop đã được ship
func (api *apiController) shipShopOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		shopOrderCode := ctx.Param("shopOrderCode")

		if shopOrderCode == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Shop order code is required"))
			return
		}

		var req services.ShipOrderRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		// Giả định shopID lấy từ token
		shopID := ctx.Query("shop_id")
		if shopID == "" {
			shopID = authPayload.Sub
		}

		err := api.service.ShipShopOrder(ctx, shopID, shopOrderCode, req)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Order marked as shipped successfully", nil))
	}
}

// updateShopOrderStatus cập nhật trạng thái đơn hàng của shop
func (api *apiController) updateShopOrderStatus() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		// shopOrderCode := ctx.Param("shopOrderCode")

		// if shopOrderCode == "" {
		// 	ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Shop order code is required"))
		// 	return
		// }

		var req struct {
			Status      string `json:"status" binding:"required"`
			ShopOrderID string `json:"shop_order_id" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		err := api.service.UpdateShopOrderStatus(ctx, req.ShopOrderID, req.Status)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Order status updated successfully", nil))
	}
}

// updateShopOrderStatus cập nhật trạng thái đơn hàng thành sử lý khi thanh toán online
func (api *apiController) callbackPaymentOnline() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		orderID := ctx.Param("order_id")

		if orderID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Order ID is required"))
			return
		}

		err := api.service.CallbackPaymentOnline(ctx, orderID)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Order status updated successfully", nil))
	}
}
