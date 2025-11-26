package controllers

import (
	"net/http"

	assets_api "github.com/TranVinhHien/ecom_payment_service/assets/api"
	"github.com/TranVinhHien/ecom_payment_service/assets/token"

	"github.com/gin-gonic/gin"
)

// create a json to test this controller

func (api *apiController) checkAuth() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("check token ok", authPayload))
	}
}
