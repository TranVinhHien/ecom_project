package controllers

import (
	"net/http"
	"os"

	assets_api "github.com/TranVinhHien/ecom_product_service/assets/api"
	"github.com/TranVinhHien/ecom_product_service/assets/token"
	controllers_model "github.com/TranVinhHien/ecom_product_service/controllers/models"
	"github.com/gin-gonic/gin"
)

func (api *apiController) renderURLLocal() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		filename := ctx.Param("id")

		filePath := api.service.RenderImage(ctx, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Nếu file không tồn tại, trả về lỗi 404
			ctx.JSON(404, gin.H{"error": "File not found", "filePath": filePath})
			return
		}
		ctx.File(filePath)
	}
}
func (api *apiController) uploadMultiMedia() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid form data file"})
			return
		}
		files := form.File["media"]
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		// result, errorr := api.service.UploadMultiMedia(ctx, authPayload.Sub, files)
		result, errorr := api.service.UploadMultiMedia(ctx, authPayload.Sub, files)
		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("upload image success", result))
	}
}

func (api *apiController) deleteMultiImage() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		var req controllers_model.DeleteMediaParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		// errorr := api.service.DeleteMultiImage(ctx, authPayload.Sub, req.ListID)
		errorr := api.service.DeleteMultiImage(ctx, "123", req.ListID)
		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("delete image success", nil))
	}
}
