package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	assets_api "github.com/TranVinhHien/ecom_analytics_service/assets/api"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"

	"github.com/gin-gonic/gin"
)

// === NHÓM I: API CHO SHOP ===

// === Nhóm 1: Tổng quan Shop ===

// getShopOverview: GET /api/v1/shop/overview
// Query params: start_date (YYYY-MM-DD), end_date (YYYY-MM-DD)
func (api apiController) getShopOverview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		// Parse query params
		startDateStr := ctx.Query("start_date")
		endDateStr := ctx.Query("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD"))
				return
			}
		} else {
			// Default: 30 days ago
			startDate = time.Now().AddDate(0, 0, -30)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD"))
				return
			}
		} else {
			// Default: today
			endDate = time.Now()
		}

		result, errors := api.service.GetShopOverview(ctx, shopID.(string), startDate, endDate)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getShopWalletSummary: GET /api/v1/shop/wallet/summary
func (api apiController) getShopWalletSummary() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		result, errors := api.service.GetShopWalletSummary(ctx, shopID.(string))
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 2: Phân tích Đơn hàng ===

// listShopOrders: GET /api/v1/shop/orders
// Query params: status, start_date, end_date, limit, offset
func (api apiController) listShopOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		params := entity.ListShopOrdersParams{
			ShopID: shopID.(string),
			Limit:  20,
			Offset: 0,
		}

		// Parse optional filters
		if status := ctx.Query("status"); status != "" {
			params.Status = sql.NullString{String: status, Valid: true}
		}

		if startDateStr := ctx.Query("start_date"); startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
			params.StartDate = sql.NullTime{Time: startDate, Valid: true}
		}

		if endDateStr := ctx.Query("end_date"); endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
			params.EndDate = sql.NullTime{Time: endDate, Valid: true}
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.ListShopOrders(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getEnrichedShopOrder: GET /api/v1/shop/orders/:shop_order_id
func (api apiController) getEnrichedShopOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		shopOrderID := ctx.Param("shop_order_id")
		if shopOrderID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "shop_order_id is required"))
			return
		}

		result, errors := api.service.GetEnrichedShopOrder(ctx, shopID.(string), shopOrderID)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listShopOrderItems: GET /api/v1/shop/order-items
// Query params: product_id, start_date, end_date, limit, offset
func (api apiController) listShopOrderItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		params := entity.ListShopOrderItemsParams{
			ShopID: shopID.(string),
			Limit:  20,
			Offset: 0,
		}

		// Parse optional filters
		if productID := ctx.Query("product_id"); productID != "" {
			params.ProductID = sql.NullString{String: productID, Valid: true}
		}

		if startDateStr := ctx.Query("start_date"); startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
			params.StartDate = sql.NullTime{Time: startDate, Valid: true}
		}

		if endDateStr := ctx.Query("end_date"); endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
			params.EndDate = sql.NullTime{Time: endDate, Valid: true}
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.ListShopOrderItems(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 3: Phân tích Doanh thu & Dòng tiền ===

// getShopRevenueTimeseries: GET /api/v1/shop/revenue/timeseries
// Query params: start_date, end_date
func (api apiController) getShopRevenueTimeseries() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		startDateStr := ctx.Query("start_date")
		endDateStr := ctx.Query("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
		} else {
			startDate = time.Now().AddDate(0, 0, -30)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
		} else {
			endDate = time.Now()
		}

		result, errors := api.service.GetShopRevenueTimeseries(ctx, shopID.(string), startDate, endDate)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listShopWalletLedgerEntries: GET /api/v1/shop/wallet/ledger-entries
// Query params: limit, offset
func (api apiController) listShopWalletLedgerEntries() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		params := entity.ListWalletLedgerEntriesParams{
			ShopID: shopID.(string),
			Limit:  20,
			Offset: 0,
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.ListShopWalletLedgerEntries(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listShopSettlements: GET /api/v1/shop/settlements
// Query params: status, start_date, end_date, limit, offset
func (api apiController) listShopSettlements() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		params := entity.ListShopSettlementsParams{
			ShopID: shopID.(string),
			Limit:  20,
			Offset: 0,
		}

		if status := ctx.Query("status"); status != "" {
			params.Status = sql.NullString{String: status, Valid: true}
		}

		if startDateStr := ctx.Query("start_date"); startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
			params.StartDate = sql.NullTime{Time: startDate, Valid: true}
		}

		if endDateStr := ctx.Query("end_date"); endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
			params.EndDate = sql.NullTime{Time: endDate, Valid: true}
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.ListShopSettlements(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 4: Phân tích Voucher ===

// listShopVouchers: GET /api/v1/shop/vouchers
// Query params: is_active, limit, offset
func (api apiController) listShopVouchers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		params := entity.ListShopVouchersParams{
			ShopID: shopID.(string),
			Limit:  20,
			Offset: 0,
		}

		if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
			isActive, err := strconv.ParseBool(isActiveStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid is_active, use true or false"))
				return
			}
			params.IsActive = sql.NullBool{Bool: isActive, Valid: true}
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.ListShopVouchers(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getShopVoucherPerformance: GET /api/v1/shop/vouchers/performance
// Query params: start_date, end_date
func (api apiController) getShopVoucherPerformance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		startDateStr := ctx.Query("start_date")
		endDateStr := ctx.Query("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
		} else {
			startDate = time.Now().AddDate(0, 0, -30)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
		} else {
			endDate = time.Now()
		}

		result, errors := api.service.GetShopVoucherPerformance(ctx, shopID.(string), startDate, endDate)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getShopVoucherUsageDetails: GET /api/v1/shop/vouchers/:voucher_id/details
// Query params: limit, offset
func (api apiController) getShopVoucherUsageDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		voucherID := ctx.Param("voucher_id")
		if voucherID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "voucher_id is required"))
			return
		}

		params := entity.ListVoucherUsageDetailsParams{
			VoucherID: voucherID,
			ShopID:    shopID.(string),
			Limit:     20,
			Offset:    0,
		}

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limit, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			params.Limit = int32(limit)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offset, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			params.Offset = int32(offset)
		}

		result, errors := api.service.GetShopVoucherUsageDetails(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 5: Xếp hạng ===

// getShopRankingProducts: GET /api/v1/shop/ranking/products
// Query params: start_date, end_date, sort_by (revenue|quantity), limit
func (api apiController) getShopRankingProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID, exists := ctx.Get("shop_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, assets_api.ResponseError(http.StatusUnauthorized, "shop_id not found in context"))
			return
		}

		startDateStr := ctx.Query("start_date")
		endDateStr := ctx.Query("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid start_date format"))
				return
			}
		} else {
			startDate = time.Now().AddDate(0, 0, -30)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format"))
				return
			}
		} else {
			endDate = time.Now()
		}

		sortBy := ctx.DefaultQuery("sort_by", "revenue")
		if sortBy != "revenue" && sortBy != "quantity" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "sort_by must be 'revenue' or 'quantity'"))
			return
		}

		limit := int32(10)
		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		params := entity.ShopRankingProductsParams{
			ShopID:    shopID.(string),
			StartDate: startDate,
			EndDate:   endDate,
			SortBy:    sortBy,
			Limit:     limit,
		}

		result, errors := api.service.GetShopRankingProducts(ctx, params)
		if errors != nil {
			ctx.JSON(http.StatusInternalServerError, assets_api.ResponseError(http.StatusInternalServerError, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}
