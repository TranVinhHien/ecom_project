package controllers

import (
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_analytics_service/assets/api"
	"github.com/TranVinhHien/ecom_analytics_service/assets/token"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"
	"github.com/gin-gonic/gin"
)

// =================================================================
// PUBLIC ENDPOINTS - Không cần authentication
// =================================================================

// submitChatboxReview: POST /api/v1/public/chatbox/review
// User đánh giá message của chatbot (Like/Dislike)
func (api apiController) submitChatboxReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req entity.SubmitMessageRatingRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "dữ liệu không hợp lệ: "+err.Error()))
			return
		}
		req.UserID = &authPayload.UserId
		serviceErr := api.service.SubmitMessageRating(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"message": "Cảm ơn bạn đã đánh giá!",
		}))
	}
}

// submitCustomerSupportComplaint: POST /api/v1/public/customer-support/complaint
// User gửi phản hồi/khiếu nại đến admin
func (api apiController) submitCustomerSupportComplaint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req entity.SubmitCustomerFeedbackRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "dữ liệu không hợp lệ: "+err.Error()))
			return
		}
		req.UserID = &authPayload.UserId
		req.Email = &authPayload.Email
		feedbackID, serviceErr := api.service.SubmitCustomerFeedback(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusCreated, assets_api.SimpSuccessResponse("success", gin.H{
			"message":     "Phản hồi của bạn đã được gửi thành công. Chúng tôi sẽ xem xét và phản hồi sớm nhất.",
			"feedback_id": feedbackID,
		}))
	}
}

// =================================================================
// ADMIN ENDPOINTS - Chatbox Statistics
// =================================================================

// getChatboxStatistics: GET /api/v1/platform/chatbox/statistics
// Admin xem thống kê tổng quan về message ratings
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getChatboxStatistics() gin.HandlerFunc {
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

		// Get overall stats
		stats, serviceErr := api.service.GetMessageRatingStats(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		// Get time series data
		timeSeries, serviceErr := api.service.GetMessageRatingsTimeSeries(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"overview":    stats,
			"time_series": timeSeries,
		}))
	}
}

// getChatboxReviews: GET /api/v1/platform/chatbox/reviews
// Admin xem danh sách chi tiết các ratings
// Query params: session_id, user_id, rating, start_date, end_date, page, page_size
func (api apiController) getChatboxReviews() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req entity.GetMessageRatingsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "tham số không hợp lệ: "+err.Error()))
			return
		}

		ratings, serviceErr := api.service.GetMessageRatingsList(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"data":      ratings,
			"page":      req.Page,
			"page_size": req.PageSize,
		}))
	}
}

// =================================================================
// ADMIN ENDPOINTS - Customer Support
// =================================================================

// listCustomerFeedbacks: GET /api/v1/platform/customer-feedback
// Admin xem danh sách customer feedbacks
// Query params: feedback_type, status, start_date, end_date, page, page_size
func (api apiController) listCustomerFeedbacks() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req entity.GetCustomerFeedbacksRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "tham số không hợp lệ: "+err.Error()))
			return
		}

		feedbacks, serviceErr := api.service.GetCustomerFeedbacks(ctx.Request.Context(), &req)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"data":      feedbacks,
			"page":      req.Page,
			"page_size": req.PageSize,
		}))
	}
}

// getCustomerFeedbackDetail: GET /api/v1/platform/customer-feedback/:id
// Admin xem chi tiết một feedback
func (api apiController) getCustomerFeedbackDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		if idStr == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "thiếu feedback ID"))
			return
		}

		feedback, serviceErr := api.service.GetCustomerFeedbackByID(ctx.Request.Context(), idStr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", feedback))
	}
}

// getCustomerSupportStatistics: GET /api/v1/platform/customer-feedback/statistics
// Admin xem thống kê về customer support
// Query params: start_date, end_date (YYYY-MM-DD)
func (api apiController) getCustomerSupportStatistics() gin.HandlerFunc {
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

		// Get overall stats
		stats, serviceErr := api.service.GetCustomerFeedbackStats(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		// Get category breakdown
		categoryStats, serviceErr := api.service.GetCustomerFeedbacksByCategory(ctx.Request.Context(), startPtr, endPtr)
		if serviceErr != nil {
			ctx.JSON(serviceErr.Code, assets_api.ResponseError(serviceErr.Code, serviceErr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("success", gin.H{
			"overview":           stats,
			"category_breakdown": categoryStats,
		}))
	}
}
