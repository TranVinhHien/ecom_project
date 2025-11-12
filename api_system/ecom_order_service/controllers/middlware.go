package controllers

import (
	"strings"

	assets_api "github.com/TranVinhHien/ecom_order_service/assets/api"
	"github.com/TranVinhHien/ecom_order_service/assets/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationKey     = "authorization"
	authorizationType    = "bearer"
	authorizationPayload = "authorization_payload"
)

func authorization(jwt token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(401, assets_api.ResponseError(401, "authorization header is not provided"))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(400, assets_api.ResponseError(400, "invalid authorization header format"))
			return
		}
		if authType := strings.ToLower(fields[0]); authType != authorizationType {
			ctx.AbortWithStatusJSON(400, assets_api.ResponseError(400, "not support type :%s", authorizationType))
			return
		}
		accessToken := fields[1]
		payload, err := jwt.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, assets_api.ResponseError(401, "invalid access token: "+err.Error()))
			return
		}
		ctx.Set(authorizationPayload, payload)
		ctx.Set("token", accessToken)
		ctx.Next()
	}
}
