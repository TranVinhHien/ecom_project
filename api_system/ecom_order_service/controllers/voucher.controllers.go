package controllers

import (
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_order_service/assets/api"
	"github.com/TranVinhHien/ecom_order_service/assets/token"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"

	"github.com/gin-gonic/gin"
)

// createVoucher handles POST /api/v1/vouchers
func (api *apiController) createVoucher() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Require auth
		tokenPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		var req services.CreateVoucherRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		if err := api.service.CreateVoucher(ctx, req, tokenPayload.UserId, tokenPayload.Scope); err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusCreated, assets_api.SimpSuccessResponse("Voucher created successfully", nil))
	}
}

// updateVoucher handles PUT /api/v1/vouchers/:voucherID
func (api *apiController) updateVoucher() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// Require auth
		tokenPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		voucherID := ctx.Param("voucherID")
		if voucherID == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "voucherID is required"))
			return
		}

		var req services.UpdateVoucherRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid request body: "+err.Error()))
			return
		}

		if err := api.service.UpdateVoucher(ctx, voucherID, tokenPayload.UserId, req); err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Voucher updated successfully", nil))
	}
}

// listVouchersForUser handles GET /api/v1/vouchers
func (api *apiController) listVouchersForUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		// Parse filter từ query parameters
		var filter services.VoucherFilterRequest
		if err := ctx.ShouldBindQuery(&filter); err != nil {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "Invalid query parameters: "+err.Error()))
			return
		}

		result, err := api.service.ListVouchersForUser(ctx, authPayload.Sub, filter)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Get vouchers successfully", result))
	}
}
