package controllers

import (
	"github.com/TranVinhHien/ecom_payment-service/assets/token"
	"github.com/TranVinhHien/ecom_payment-service/services"

	"github.com/gin-gonic/gin"
)

type apiController struct {
	service services.ServiceUseCase
	jwt     token.Maker
}

func NewAPIController(s services.ServiceUseCase, jwt token.Maker) apiController {
	return apiController{service: s, jwt: jwt}
}

func (api apiController) SetUpRoute(group *gin.RouterGroup) {
	group.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(200)
	})

	payment := group.Group("/transaction") //.Use(authorization(api.jwt))
	{
		payment_auth := payment.Group("").Use(authorization(api.jwt))
		{
			payment_auth.POST("/init", api.initPayment())
		}
		payment.GET("/payment_method", api.ListPayment())
		payment.GET("/payment_method/:id", api.PaymentDetail())
		payment.POST("/callback", api.callbackMoMo())

	}
}
