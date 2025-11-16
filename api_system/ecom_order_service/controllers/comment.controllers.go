package controllers

import (
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_order_service/assets/api"
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"

	"github.com/gin-gonic/gin"
)

// createComment handles POST /api/v1/comments
// Tạo đánh giá/bình luận cho sản phẩm đã mua
func (api *apiController) createComment() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Require auth - lấy user_id từ token
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		var req services.CreateCommentRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		// Call service để tạo comment
		if err := api.service.CreateComment(ctx, authPayload.Sub, req); err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusCreated, assets_api.SimpSuccessResponse("Đánh giá sản phẩm thành công", nil))
	}
}

// listComments handles GET /api/v1/comments
// Lấy danh sách bình luận cho một sản phẩm (có phân trang)
func (api *apiController) listComments() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Parse query parameters
		var req services.ListCommentsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid query parameters: "+err.Error()))
			return
		}

		// Set default limit nếu không được cung cấp
		if req.PageSize == 0 {
			req.PageSize = 20 // Default 20 comments per page
		}
		if req.Page <= 0 {
			req.Page = 1
		}
		// Call service để lấy comments
		result, err := api.service.ListComments(ctx, req)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Lấy danh sách đánh giá thành công", result))
	}
}

// checkReviewedItems handles POST /api/v1/comments/check-reviewed
// Kiểm tra danh sách order_item_id nào đã được đánh giá
// API này dành cho service khác (như order service) gọi để check trạng thái review
func (api *apiController) checkReviewedItems() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req services.CheckReviewedItemsRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		// Call service để check reviewed items
		result, err := api.service.CheckReviewedItems(ctx, req)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Kiểm tra trạng thái đánh giá thành công", result))
	}
}

// getBulkProductRatingStats handles POST /api/v1/comments/bulk-stats
// Lấy thống kê đánh giá (điểm trung bình và tổng số lượt đánh giá) cho nhiều sản phẩm
// Trả về map với product_id là key
func (api *apiController) getBulkProductRatingStats() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req services.GetBulkProductRatingStatsRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		// Call service để lấy bulk stats
		result, err := api.service.GetBulkProductRatingStats(ctx, req)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Lấy thống kê đánh giá sản phẩm thành công", result))
	}
}
