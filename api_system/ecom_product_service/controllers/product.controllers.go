package controllers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	assets_api "github.com/TranVinhHien/ecom_product_service/assets/api"
	"github.com/TranVinhHien/ecom_product_service/assets/token"
	controllers_model "github.com/TranVinhHien/ecom_product_service/controllers/models"
	services "github.com/TranVinhHien/ecom_product_service/services/entity"
	"github.com/jinzhu/copier"

	"github.com/gin-gonic/gin"
)

func (api *apiController) getAllProductSimple() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		page := ctx.DefaultQuery("page", "1")
		limit := ctx.DefaultQuery("limit", "10")
		brand := ctx.DefaultQuery("brand", "")
		cate_path := ctx.DefaultQuery("cate_path", "")
		shop_id := ctx.DefaultQuery("shop_id", "")
		price_min_s := ctx.DefaultQuery("price_min", "-1")
		price_max_s := ctx.DefaultQuery("price_max", "-1")
		keywords := ctx.DefaultQuery("keywords", "")
		sort := ctx.DefaultQuery("sort", "")
		status := ctx.DefaultQuery("status", "")
		sort_order := []string{"price_asc", "price_desc", "name_asc", "name_desc", "best_sell"}
		// check if sort not in sort_order
		if sort != "" {
			check := false

			for _, v := range sort_order {
				if sort == v {
					sort = v
					check = true
					break
				}
			}
			if !check {
				ctx.JSON(402, assets_api.ResponseError(402, "sort must be one of "+strconv.Quote(strings.Join(sort_order, ", "))))
				return
			}
		}

		deleteStatus := []string{"Pending", "Deleted", "Active"}
		// check if sort not in DeleteStatus
		if status != "" {
			check := false

			for _, v := range deleteStatus {
				if status == v {
					status = v
					check = true
					break
				}
			}
			if !check {
				ctx.JSON(402, assets_api.ResponseError(402, "status must be one of "+strconv.Quote(strings.Join(deleteStatus, ", "))))
				return
			}
		}

		pageInt, errors := strconv.Atoi(page)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}
		pageSizeInt, errors := strconv.Atoi(limit)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}

		price_min, errors := strconv.Atoi(price_min_s)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}
		price_max, errors := strconv.Atoi(price_max_s)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}

		orders, err := api.service.GetAllProductSimple(ctx, services.NewQueryFilter(pageInt, pageSizeInt, nil, nil), cate_path, brand, shop_id, keywords, sort, float64(price_min), float64(price_max), status)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", orders))
	}
}
func (api *apiController) getDetailProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		spu := ctx.Param("id")

		if spu == "" {
			ctx.JSON(402, assets_api.ResponseError(402, "must provide spu_id"))
			return
		}

		product_spu, err := api.service.GetDetailProduct(ctx, spu)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", product_spu))
	}
}
func (api *apiController) createProduct() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		// 1. Parse JSON phần "product"
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		token := ctx.MustGet("token").(string)

		// 2. Xử lý ảnh chính
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mediaFiles := form.File["media"]
		var option_images []struct {
			OptionName string
			Value      string
			Image      *multipart.FileHeader
		}
		images := form.File["image"]
		if len(images) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
			return
		}
		image := images[0]

		var req controllers_model.Product
		productJSON := ctx.PostForm("product")
		if err := json.Unmarshal([]byte(productJSON), &req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Map using copier
		var productParams services.ProductParams
		if err := copier.Copy(&productParams, &req); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 3. Xử lý ảnh media
		// lấy ảnh của option value
		optionValueImages := form.File
		for i := range req.OptionValue {
			key := fmt.Sprintf("option_value_images[%d]", i)
			if files, ok := optionValueImages[key]; ok && len(files) > 0 {
				file := files[0]
				option_images = append(option_images, struct {
					OptionName string
					Value      string
					Image      *multipart.FileHeader
				}{
					OptionName: req.OptionValue[i].OptionName,
					Value:      req.OptionValue[i].Value,
					Image:      file,
				})
			}
		}
		// 5. Gọi hàm xử lý logic
		errors := api.service.CreateProduct(ctx, token, authPayload.Sub,
			productParams, image, mediaFiles, option_images)

		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create product successfull", nil))

	}
}

func (api *apiController) updateProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 1. Lấy product ID từ URL params
		productID := ctx.Param("id")
		if productID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
			return
		}

		// 2. Lấy thông tin user từ token
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		// 3. Parse multipart form
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 4. Xử lý ảnh chính (optional)
		var image *multipart.FileHeader
		images := form.File["image"]
		if len(images) > 0 {
			image = images[0]
		}

		// 5. Xử lý media files (optional)
		mediaFiles := form.File["media"]

		// 6. Parse JSON product data
		var req controllers_model.ProductUpdate
		productJSON := ctx.PostForm("product")
		if productJSON != "" {
			if err := json.Unmarshal([]byte(productJSON), &req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		// 7. Map to service params
		var productParams services.ProductUpdateParams
		if err := copier.Copy(&productParams, &req); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 8. Xử lý ảnh option value (optional)
		var option_images []struct {
			OptionName string
			Value      string
			Image      *multipart.FileHeader
		}
		optionValueImages := form.File
		for i := range req.OptionValue {
			key := fmt.Sprintf("option_value_images[%d]", i)
			if files, ok := optionValueImages[key]; ok && len(files) > 0 {
				file := files[0]
				option_images = append(option_images, struct {
					OptionName string
					Value      string
					Image      *multipart.FileHeader
				}{
					OptionName: req.OptionValue[i].OptionName,
					Value:      req.OptionValue[i].Value,
					Image:      file,
				})
			}
		}

		// 9. Gọi service để update product
		errors := api.service.UpdateProduct(ctx, authPayload.Scope, authPayload.Sub, productID,
			productParams, image, mediaFiles, nil)

		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update product successfully", nil))
	}
}
func (api *apiController) updateSKUReserverProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// _ = ctx.MustGet(authorizationPayload).(*token.Payload)

		var req controllers_model.UpdateSKUReserverRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(req.Data) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_sku is required"})
			return
		}
		var productSKU []services.ProductUpdateSKUReserver
		if err := copier.Copy(&productSKU, &req.Data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		errors := api.service.UpdateSKUReserverProduct(ctx, productSKU, services.ProductUpdateType(req.Status))
		if errors != nil {
			ctx.JSON(errors.Code, assets_api.ResponseError(errors.Code, errors.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update sku reserver product successfully", nil))
	}
}
func (api *apiController) getSKUProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sku_id := ctx.Param("id")
		if sku_id == "" {
			ctx.JSON(402, assets_api.ResponseError(402, "must provide sku_id"))
			return
		}
		product_sku, err := api.service.GetSKUProduct(ctx, sku_id)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get sku product successfull", product_sku))
	}
}
func (api *apiController) getProductWithID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sku_id := ctx.Param("id")
		if sku_id == "" {
			ctx.JSON(402, assets_api.ResponseError(402, "must provide sku_id"))
			return
		}
		product, err := api.service.GetProductWithID(ctx, sku_id)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get product successfull", product))
	}
}
func (api *apiController) buildProductSearchString() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		productID := ctx.Param("id")
		if productID == "" {
			ctx.JSON(402, assets_api.ResponseError(402, "must provide product_id"))
			return
		}
		searchString, err := api.service.BuildProductSearchString(ctx, productID)
		if err != nil {
			ctx.JSON(500, assets_api.ResponseError(500, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("build product search string successfully", gin.H{"search_string": searchString}))
	}
}
func (api *apiController) getAllProductID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		productIDs, err := api.service.GetALLProductID(ctx)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get all product IDs successfully", gin.H{"product_ids": productIDs}))
	}
}
func (api *apiController) getListProductWithIDs() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		productIDs := ctx.QueryArray("product_ids")
		if len(productIDs) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_ids is required"})
			return
		}

		products, err := api.service.GetListProductWithIDs(ctx, productIDs)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get list product with IDs successfully", products))
	}
}
