package controllers

// func (api *apiController) listCategories() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
// 		id := ctx.Query("cate_id")
// 		addrs, err := api.service.GetCategoris(ctx, id)
// 		if err != nil {
// 			fmt.Print(err)
// 			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", addrs))
// 	}
// }
// func (api *apiController) createCategories() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

// 		form, err := ctx.MultipartForm()
// 		if err != nil {
// 			ctx.JSON(400, gin.H{"error": "Invalid form data"})
// 			return
// 		}
// 		name := form.Value["name"][0]
// 		parent := form.Value["parent"][0]
// 		media := form.File["media"]
// 		var file *multipart.FileHeader
// 		if len(media) > 0 {
// 			file = media[0]
// 		}
// 		cat := services.Categorys{
// 			Name:   name,
// 			Parent: services.Narg[string]{Valid: parent != "", Data: parent},
// 		}
// 		errors := api.service.AddCategory(ctx, authPayload.Sub, cat, file)
// 		if errors != nil {
// 			fmt.Print(errors)
// 			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
// 			return
// 		}
// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create category successful", nil))
// 	}
// }
// func (api *apiController) updateCategories() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

// 		form, err := ctx.MultipartForm()
// 		if err != nil {
// 			ctx.JSON(400, gin.H{"error": "Invalid form data"})
// 			return
// 		}
// 		categoryID := form.Value["cate_id"][0]
// 		if categoryID == "" {
// 			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "cate_id is required"))
// 			return
// 		}
// 		name := form.Value["name"][0]
// 		parent := form.Value["parent"][0]
// 		media := form.File["media"]
// 		var file *multipart.FileHeader
// 		if len(media) > 0 {
// 			file = media[0]
// 		}
// 		cat := services.Categorys{
// 			CategoryID: categoryID,
// 			Name:       name,
// 			Parent:     services.Narg[string]{Valid: parent != "", Data: parent},
// 		}
// 		errors := api.service.UpdateCategory(ctx, authPayload.Sub, cat, file)
// 		if errors != nil {
// 			fmt.Print(errors)
// 			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update category successful", nil))
// 	}
// }
// func (api *apiController) deleteCategories() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
// 		categoryID := ctx.Param("cate_id") // or ctx.Query("cate_id")
// 		if categoryID == "" {
// 			ctx.JSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, "cate_id is required"))
// 			return
// 		}
// 		err := api.service.DeleteCategory(ctx, authPayload.Sub, categoryID)
// 		if err != nil {
// 			fmt.Print(err)
// 			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("delete category successful", nil))
// 	}
// }
