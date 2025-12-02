package controllers

import (
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_analytics_service/assets/api"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"
	"github.com/gin-gonic/gin"
)

// =================================================================
// AGENT ANALYTICS ENDPOINTS - For Platform/Admin
// =================================================================

// getDashboardStats: GET /api/v1/platform/agent-analytics/dashboard
// Admin xem thống kê tổng quan về agent interactions
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getDashboardStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		var startPtr, endPtr *string
		if startDate != "" {
			startPtr = &startDate
		}
		if endDate != "" {
			endPtr = &endDate
		}

		stats, serviceErr := api.service.GetDashboardStats(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", stats))
	}
}

// getMessageVolumeByHour: GET /api/v1/platform/agent-analytics/message-volume
// Admin xem mật độ tin nhắn theo giờ (heatmap)
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getMessageVolumeByHour() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		var startPtr, endPtr *string
		if startDate != "" {
			startPtr = &startDate
		}
		if endDate != "" {
			endPtr = &endDate
		}

		data, serviceErr := api.service.GetMessageVolumeByHour(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", data))
	}
}

// getTopActiveUsers: GET /api/v1/platform/agent-analytics/top-users
// Admin xem top người dùng tích cực nhất
// Query params: start_date, end_date (YYYY-MM-DD), limit (default 10)
func (api apiController) getTopActiveUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req entity.GetTopActiveUsersRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "tham số không hợp lệ: "+err.Error()))
			return
		}

		data, serviceErr := api.service.GetTopActiveUsers(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"data":  data,
			"limit": req.Limit,
		}))
	}
}

// getTopicStats: GET /api/v1/platform/agent-analytics/topics
// Admin xem thống kê các chủ đề thường được hỏi
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getTopicStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		var startPtr, endPtr *string
		if startDate != "" {
			startPtr = &startDate
		}
		if endDate != "" {
			endPtr = &endDate
		}

		data, serviceErr := api.service.GetTopicStats(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", data))
	}
}

// getPurchaseIntentStats: GET /api/v1/platform/agent-analytics/purchase-intent
// Admin xem thống kê ý định mua hàng từ hội thoại
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getPurchaseIntentStats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		var startPtr, endPtr *string
		if startDate != "" {
			startPtr = &startDate
		}
		if endDate != "" {
			endPtr = &endDate
		}

		data, serviceErr := api.service.GetPurchaseIntentStats(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", data))
	}
}

// getTopMentionedCategories: GET /api/v1/platform/agent-analytics/top-categories
// Admin xem top danh mục sản phẩm được nhắc đến nhiều nhất
// Query params: start_date, end_date (YYYY-MM-DD), limit (default 10)
func (api apiController) getTopMentionedCategories() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req entity.GetTopMentionedCategoriesRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "tham số không hợp lệ: "+err.Error()))
			return
		}

		data, serviceErr := api.service.GetTopMentionedCategories(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"data":  data,
			"limit": req.Limit,
		}))
	}
}
