package controllers

import (
	"os"

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
