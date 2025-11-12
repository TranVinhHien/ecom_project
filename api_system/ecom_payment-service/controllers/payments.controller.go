package controllers

import (
	"fmt"
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_payment-service/assets/api"
	"github.com/TranVinhHien/ecom_payment-service/assets/token"
	controllers_model "github.com/TranVinhHien/ecom_payment-service/controllers/models"
	services "github.com/TranVinhHien/ecom_payment-service/services/entity"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

func (api *apiController) ListPayment() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		addrs, err := api.service.ListPaymentMethod(ctx)
		if err != nil {
			fmt.Print(err)
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", addrs))
	}
}

func (api *apiController) PaymentDetail() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "id is required"))
			return
		}
		addrs, err := api.service.PaymentMethodDetail(ctx, id)
		if err != nil {
			fmt.Print(err)
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get payment method detail success", addrs))
	}
}
func (api *apiController) initPayment() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.InitPaymentParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		var order services.InitPaymentParams
		if err := copier.Copy(&order, &req); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		payment, errorr := api.service.InitPayment(ctx, authPayload.Sub, authPayload.Email, order)

		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create order success", payment))
	}
}

func (api *apiController) callbackMoMo() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req services.TransactionMoMO
		if err := ctx.ShouldBindJSON(&req); err != nil {
			// ctx.JSON(http.StatusNoContent, assets_api.SimpSuccessResponse("Hello", nil))
			// return
		}
		api.service.CallBackMoMo(ctx, req)
		// fmt.Printf("momo tra ve : %s", req)
		ctx.JSON(http.StatusNoContent, assets_api.SimpSuccessResponse("Hello", nil))
	}
}
