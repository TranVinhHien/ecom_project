package controllers

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (api *apiController) renderAndroid() func(c *gin.Context) {
	// lệnh lấy sha256 // keytool -list -v -keystore "%USERPROFILE%\.android\debug.keystore" -alias androiddebugkey -storepass android -keypass android
	return func(ctx *gin.Context) {
		// Lấy đường dẫn gốc của project
		rootDir, _ := os.Getwd()
		// Ghép đường dẫn tới file json
		jsonPath := filepath.Join(rootDir, "assets", "setup-mobile", "check_android.json")
		// fmt.Println("Serving file:", jsonPath)
		ctx.Header("Content-Type", "application/json")
		ctx.File(jsonPath)
	}
}
