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

// === NHÓM II: API CHO PLATFORM ===

// === Nhóm 1: Tổng quan ===

// getPlatformOverview: GET /api/v1/platform/overview
// Query params: start_date, end_date
func (api apiController) getPlatformOverview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
			startDate = time.Now().AddDate(0, 0, -30)
		}

		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD"))
				return
			}
		} else {
			endDate = time.Now()
		}

		result, errors := api.service.GetPlatformOverview(ctx, startDate, endDate)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 2: Quản lý Đơn hàng ===

// listPlatformOrders: GET /api/v1/platform/orders
// Query params: shop_id, user_id, status, start_date, end_date, limit, offset
func (api apiController) listPlatformOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformOrdersParams{
			Limit:  20,
			Offset: 0,
		}

		// Parse optional filters
		if shopID := ctx.Query("shop_id"); shopID != "" {
			params.ShopID = sql.NullString{String: shopID, Valid: true}
		}

		if userID := ctx.Query("user_id"); userID != "" {
			params.UserID = sql.NullString{String: userID, Valid: true}
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

		result, errors := api.service.ListPlatformOrders(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getEnrichedPlatformOrder: GET /api/v1/platform/orders/:order_id
func (api apiController) getEnrichedPlatformOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderID := ctx.Param("order_id")
		if orderID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "order_id is required"))
			return
		}

		result, errors := api.service.GetEnrichedPlatformOrder(ctx, orderID)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 3: Quản lý Tài chính ===

// getPlatformRevenueTimeseries: GET /api/v1/platform/finance/revenue-timeseries
// Query params: start_date, end_date
func (api apiController) getPlatformRevenueTimeseries() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		result, errors := api.service.GetPlatformRevenueTimeseries(ctx, startDate, endDate)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listPlatformTransactions: GET /api/v1/platform/finance/transactions
// Query params: type, status, start_date, end_date, limit, offset
func (api apiController) listPlatformTransactions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformTransactionsParams{
			Limit:  20,
			Offset: 0,
		}

		if typeStr := ctx.Query("type"); typeStr != "" {
			params.Type = sql.NullString{String: typeStr, Valid: true}
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

		result, errors := api.service.ListPlatformTransactions(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listPlatformSettlements: GET /api/v1/platform/finance/settlements
// Query params: status, start_date, end_date, limit, offset
func (api apiController) listPlatformSettlements() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformSettlementsParams{
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

		result, errors := api.service.ListPlatformSettlements(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listPlatformLedgers: GET /api/v1/platform/finance/ledgers
// Query params: owner_type, limit, offset
func (api apiController) listPlatformLedgers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformLedgersParams{
			Limit:  20,
			Offset: 0,
		}

		if ownerType := ctx.Query("owner_type"); ownerType != "" {
			params.OwnerType = sql.NullString{String: ownerType, Valid: true}
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

		result, errors := api.service.ListPlatformLedgers(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// listLedgerEntries: GET /api/v1/platform/finance/ledgers/:ledger_id/entries
// Query params: limit, offset
func (api apiController) listLedgerEntries() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ledgerID := ctx.Param("ledger_id")
		if ledgerID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "ledger_id is required"))
			return
		}

		limit := int32(20)
		offset := int32(0)

		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			offsetInt, err := strconv.ParseInt(offsetStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid offset"))
				return
			}
			offset = int32(offsetInt)
		}

		result, errors := api.service.ListLedgerEntries(ctx, ledgerID, limit, offset)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 4: Phân tích Voucher ===

// listPlatformVouchers: GET /api/v1/platform/vouchers
// Query params: owner_type, is_active, limit, offset
func (api apiController) listPlatformVouchers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformVouchersParams{
			Limit:  20,
			Offset: 0,
		}

		if ownerType := ctx.Query("owner_type"); ownerType != "" {
			params.OwnerType = sql.NullString{String: ownerType, Valid: true}
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

		result, errors := api.service.ListPlatformVouchers(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getPlatformVoucherPerformance: GET /api/v1/platform/vouchers/performance
// Query params: start_date, end_date
func (api apiController) getPlatformVoucherPerformance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		result, errors := api.service.GetPlatformVoucherPerformance(ctx, startDate, endDate)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 5: Phân tích Shop ===

// listPlatformShops: GET /api/v1/platform/shops
// Query params: limit, offset
func (api apiController) listPlatformShops() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := entity.ListPlatformShopsParams{
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

		result, errors := api.service.ListPlatformShops(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getPlatformShopDetail: GET /api/v1/platform/shops/:shop_id/detail
// Query params: start_date, end_date
func (api apiController) getPlatformShopDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shopID := ctx.Param("shop_id")
		if shopID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "shop_id is required"))
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

		result, errors := api.service.GetPlatformShopDetail(ctx, shopID, startDate, endDate)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// === Nhóm 6: Xếp hạng Toàn Sàn ===

// getPlatformRankingShops: GET /api/v1/platform/ranking/shops
// Query params: start_date, end_date, limit
func (api apiController) getPlatformRankingShops() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		limit := int32(10)
		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		params := entity.PlatformRankingParams{
			StartDate: startDate,
			EndDate:   endDate,
			Limit:     limit,
		}

		result, errors := api.service.GetPlatformRankingShops(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getPlatformRankingProducts: GET /api/v1/platform/ranking/products
// Query params: start_date, end_date, limit
func (api apiController) getPlatformRankingProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		limit := int32(10)
		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		params := entity.PlatformRankingParams{
			StartDate: startDate,
			EndDate:   endDate,
			Limit:     limit,
		}

		result, errors := api.service.GetPlatformRankingProducts(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getPlatformRankingUsers: GET /api/v1/platform/ranking/users
// Query params: start_date, end_date, limit
func (api apiController) getPlatformRankingUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		limit := int32(10)
		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		params := entity.PlatformRankingParams{
			StartDate: startDate,
			EndDate:   endDate,
			Limit:     limit,
		}

		result, errors := api.service.GetPlatformRankingUsers(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}

// getPlatformRankingCategories: GET /api/v1/platform/ranking/categories
// Query params: start_date, end_date, limit
func (api apiController) getPlatformRankingCategories() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		limit := int32(10)
		if limitStr := ctx.Query("limit"); limitStr != "" {
			limitInt, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "invalid limit"))
				return
			}
			limit = int32(limitInt)
		}

		params := entity.PlatformRankingParams{
			StartDate: startDate,
			EndDate:   endDate,
			Limit:     limit,
		}

		result, errors := api.service.GetPlatformRankingCategories(ctx, params)
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", result))
	}
}
