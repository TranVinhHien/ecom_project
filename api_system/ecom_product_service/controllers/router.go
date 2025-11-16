package controllers

import (
	"github.com/TranVinhHien/ecom_product_service/assets/token"
	"github.com/TranVinhHien/ecom_product_service/services"

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
	group.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Status(200)

	})

	categories := group.Group("/categories")
	{
		categories.GET("/get", api.listCategories())
		categories_auth := categories.Group("").Use(authorization(api.jwt)).Use(checkRole("ADMIN"))
		{
			categories_auth.POST("/create", api.createCategories())
			categories_auth.PUT("/update", api.updateCategories())
			categories_auth.DELETE("/delete/:id", api.deleteCategories())
		}
	}
	product := group.Group("/product")
	{
		product.GET("/getall", api.getAllProductSimple())
		product.GET("/getdetail/:id", api.getDetailProduct())
		product.GET("/build_search_string/:id", api.buildProductSearchString())
		product.GET("/getallproductid", api.getAllProductID())
		product.GET("/get_products_detail_for_search", api.getListProductWithIDs())
		// chỉ cho phép shop mới được tạo/sửa sản phẩm
		product_auth := product.Group("").Use(authorization(api.jwt)).Use(checkRole("ROLE_SELLER"))
		{
			product_auth.POST("/create", api.createProduct())
			product_auth.PUT("/update/:id", api.updateProduct())
		}
		// sau này tạo thêm check endpoint chỉ cho phép admin mới được xóa sản phẩm
		product.POST("/update_sku_reserver", api.updateSKUReserverProduct())

		product.GET("/getsku/:id", api.getSKUProduct())
		product.GET("/getdetail_with_id/:id", api.getProductWithID())
	}
	// group.GET("/.well-known/assetlinks.json", api.renderAndroid())
	group.GET("/media/:id", api.renderURLLocal())
}
