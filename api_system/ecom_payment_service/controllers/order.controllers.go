package controllers

// func (api *apiController) getMyOrders() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
// 		page := ctx.DefaultQuery("page", "1")
// 		limit := ctx.DefaultQuery("limit", "10")
// 		pageInt, errors := strconv.Atoi(page)
// 		if errors != nil {
// 			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
// 			return
// 		}
// 		pageSizeInt, errors := strconv.Atoi(limit)
// 		if errors != nil {
// 			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
// 			return
// 		}
// 		order_by := ctx.Query("order_by")
// 		order_option := ctx.Query("order_option")
// 		search := ctx.Query("search")
// 		condition, errors := controllers_assets.ParseConditions(search)
// 		if errors != nil {
// 			ctx.JSON(400, assets_api.ResponseError(400, errors.Error()))
// 			return
// 		}
// 		order := &services.OrderBy{}
// 		if order_by != "" {
// 			order.Field = order_by
// 			order.Field = order_option
// 		} else {
// 			order = nil
// 		}
// 		// orders, err := api.service.ListOrderByUserID(ctx, authPayload.Sub, services.NewQueryFilter(pageInt, pageSizeInt, condition, order))
// 		if err != nil {
// 			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", orders))
// 	}
// }
// func (api *apiController) getOrderOnline() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
// 		orders, err := api.service.GetURLOrderMoMOAgain(ctx, authPayload.Sub)
// 		if err != nil {
// 			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", orders))
// 	}
// }
// func (api *apiController) cancelOrder() func(ctx *gin.Context) {
// 	return func(ctx *gin.Context) {

// 		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
// 		var req controllers_model.OrderIDParams
// 		if err := ctx.ShouldBindJSON(&req); err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
// 			return
// 		}
// 		err := api.service.CancelOrder(ctx, authPayload.Sub, req.OrderID)
// 		if err != nil {
// 			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("cancel order successfull", nil))
// 	}
// }
