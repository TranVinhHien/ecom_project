package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"strings"

	db "github.com/TranVinhHien/ecom_product_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_product_service/services/assets"
	services "github.com/TranVinhHien/ecom_product_service/services/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func (s *service) GetSKUProduct(ctx context.Context, product_sku_id string) (map[string]interface{}, *assets_services.ServiceError) {
	//log.Printf("[GetSKUProduct] Bắt đầu lấy thông tin SKU với ID: %s", product_sku_id)

	product_sku, err := s.repository.GetProductSKU(ctx, product_sku_id)
	if err != nil {
		//log.Printf("[GetSKUProduct] LỖI: Không thể lấy thông tin SKU với ID %s. Chi tiết: %v", product_sku_id, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không tìm thấy SKU với ID: %s. Lỗi: %v", product_sku_id, err))
	}

	result := assets_services.NormalizeSQLNulls(product_sku, "data")
	//log.Printf("[GetSKUProduct] Thành công lấy thông tin SKU với ID: %s", product_sku_id)
	return result, nil
}
func (s *service) GetAllProductSimple(ctx context.Context, query services.QueryFilter, category_path, brand_code, shop_id, keywords, sort string, min_price, max_price float64, status string) (map[string]interface{}, *assets_services.ServiceError) {
	//log.Printf("[GetAllProductSimple] Bắt đầu lấy danh sách sản phẩm - Trang: %d, Kích thước: %d, Danh mục: %s, Thương hiệu: %s, Shop: %s, Từ khóa: %s",
	//	query.Page, query.PageSize, category_path, brand_code, shop_id, keywords)
	cate_id := ""
	brand_id := ""
	if category_path != "" {
		category, err := s.repository.GetCategoryByPath(ctx, sql.NullString{String: category_path, Valid: true})
		if err != nil {
			return nil, assets_services.NewError(400, fmt.Errorf("không tìm thấy danh mục với đường dẫn: %s. Lỗi: %v", category_path, err))
		}
		cate_id = category.CategoryID
	}
	if brand_code != "" {
		brand, err := s.repository.GetBrandByCode(ctx, brand_code)
		if err != nil {
			return nil, assets_services.NewError(400, fmt.Errorf("không tìm thấy thương hiệu với mã: %s. Lỗi: %v", brand_code, err))
		}
		brand_id = brand.BrandID
	}
	var deleteStatus db.ProductDeleteStatus
	switch status {
	case "Pending":
		deleteStatus = db.ProductDeleteStatusPending
	case "Deleted":
		deleteStatus = db.ProductDeleteStatusDeleted
	default:
		deleteStatus = db.ProductDeleteStatusActive
	}
	product_spu, err := s.repository.ListProductsDynamic(ctx, db.ListProductsAdvancedParams{
		Limit:        int32(query.PageSize),
		Offset:       int32((query.Page - 1) * query.PageSize),
		BrandID:      sql.NullString{String: brand_id, Valid: brand_id != ""},
		DeleteStatus: db.NullProductDeleteStatus{ProductDeleteStatus: deleteStatus, Valid: true},
		CategoryID:   sql.NullString{String: cate_id, Valid: cate_id != ""},
		ShopID:       sql.NullString{String: shop_id, Valid: shop_id != ""},
		PriceMin:     sql.NullFloat64{Float64: min_price, Valid: min_price >= 0},
		PriceMax:     sql.NullFloat64{Float64: max_price, Valid: max_price >= 0},
		Keyword:      sql.NullString{String: keywords, Valid: keywords != ""},
		Sort:         sql.NullString{String: strings.ToLower(sort), Valid: sort != ""},
	})
	if err != nil {
		//log.Printf("[GetAllProductSimple] LỖI: Không thể lấy danh sách sản phẩm từ database. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách sản phẩm. Lỗi: %v", err))
	}

	totalElements, err := s.repository.CountProductsDynamic(ctx, db.ListProductsAdvancedParams{
		BrandID:      sql.NullString{String: brand_id, Valid: brand_id != ""},
		CategoryID:   sql.NullString{String: cate_id, Valid: cate_id != ""},
		ShopID:       sql.NullString{String: shop_id, Valid: shop_id != ""},
		PriceMin:     sql.NullFloat64{Float64: min_price, Valid: min_price >= 0},
		PriceMax:     sql.NullFloat64{Float64: max_price, Valid: max_price >= 0},
		Keyword:      sql.NullString{String: keywords, Valid: keywords != ""},
		DeleteStatus: db.NullProductDeleteStatus{ProductDeleteStatus: deleteStatus, Valid: true},
	})
	if err != nil {
		//log.Printf("[GetAllProductSimple] LỖI: Không thể đếm tổng số sản phẩm. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể đếm tổng số sản phẩm. Lỗi: %v", err))
	}

	// Lấy danh sách product_id từ kết quả
	productIDs := make([]string, len(product_spu))
	for i, product := range product_spu {
		productIDs[i] = product.ID
	}

	// Gọi API để lấy thông tin đánh giá cho tất cả sản phẩm
	ratingStats := make(map[string]services.ProductRating)
	if len(productIDs) > 0 {
		stats, err := s.apiServer.GetBulkProductRatingStats(productIDs)
		if err != nil {
			//log.Printf("[GetAllProductSimple] CẢNH BÁO: Không thể lấy thông tin đánh giá. Chi tiết: %v", err)
			// Không return error, chỉ log warning và tiếp tục với dữ liệu rỗng
		} else {
			// Convert sang ProductRating
			for productID, stat := range stats {
				ratingStats[productID] = services.ProductRating{
					ProductID:     stat.ProductID,
					TotalReviews:  stat.TotalReviews,
					AverageRating: stat.AverageRating,
				}
			}
		}
	}
	// Tạo slice mới kết hợp product và rating
	productsWithRating := make([]interface{}, len(product_spu))
	for i, product := range product_spu {
		// Convert product sang map
		productMap := assets_services.NormalizeToInterface(product)

		// Thêm rating vào product
		if rating, exists := ratingStats[product.ID]; exists {
			productMap.(map[string]interface{})["rating"] = rating
		} else {
			// Nếu không có rating, set giá trị mặc định
			productMap.(map[string]interface{})["rating"] = services.ProductRating{
				ProductID:     product.ID,
				TotalReviews:  0,
				AverageRating: 0.0,
			}
		}

		productsWithRating[i] = productMap
	}

	// Tạo result với data đã có rating
	result := make(map[string]interface{})
	result["data"] = productsWithRating
	totalPage := int64(math.Ceil(float64(totalElements) / float64(query.PageSize)))
	result["currentPage"] = query.Page
	result["totalPages"] = totalPage
	result["totalElements"] = totalElements
	result["limit"] = query.PageSize

	//log.Printf("[GetAllProductSimple] Thành công lấy %d sản phẩm - Trang %d/%d", len(product_spu), query.Page, totalPage)
	return result, nil
}

func (s *service) GetProductWithID(ctx context.Context, product_id string) (map[string]interface{}, *assets_services.ServiceError) {
	//log.Printf("[GetProductWithID] Bắt đầu lấy chi tiết sản phẩm với ID: %s", product_id)

	product_spu_detail, err := s.repository.GetProduct(ctx, product_id)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không tìm thấy sản phẩm với ID: %s. Chi tiết: %v", product_id, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không tìm thấy sản phẩm với ID: %s. Lỗi: %v", product_id, err))
	}

	// call sku
	sku, err := s.repository.ListSKUsByProduct(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể lấy danh sách SKU cho sản phẩm %s. Chi tiết: %v", product_id, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách SKU. Lỗi: %s", err.Error()))
	}
	sku_res := make([]services.ProductSku, len(sku))
	if err := copier.Copy(&sku_res, &sku); err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể sao chép dữ liệu SKU. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu SKU: %s", err.Error()))
	}

	// call option value
	option, err := s.repository.ListOptionValuesByProductID(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể lấy danh sách Option Values cho sản phẩm %s. Chi tiết: %v", product_id, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách thuộc tính sản phẩm. Lỗi: %s", err.Error()))
	}
	option_res := make([]services.OptionValue, len(option))
	if err := copier.Copy(&option_res, &option); err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể sao chép dữ liệu Option Values. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu thuộc tính: %s", err.Error()))
	}

	// call sku attr
	sku_attr, err := s.repository.ListSKUOptionValuesByProductID(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể lấy thông tin liên kết SKU-Option cho sản phẩm %s. Chi tiết: %v", product_id, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin liên kết SKU. Lỗi: %s", err.Error()))
	}
	sku_attr_res := make([]services.SkuAttr, len(sku_attr))
	if err := copier.Copy(&sku_attr_res, &sku_attr); err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể sao chép dữ liệu SKU Attributes. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu liên kết: %s", err.Error()))
	}

	// call brand name
	brand, err := s.repository.GetBrand(ctx, product_spu_detail.BrandID.String)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể lấy thông tin thương hiệu %s. Chi tiết: %v", product_spu_detail.BrandID.String, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin thương hiệu. Lỗi: %s", err.Error()))
	}

	// call category name
	category, err := s.repository.GetCategory(ctx, product_spu_detail.CategoryID)
	if err != nil {
		//log.Printf("[GetProductWithID] LỖI: Không thể lấy thông tin danh mục %s. Chi tiết: %v", product_spu_detail.CategoryID, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin danh mục. Lỗi: %s", err.Error()))
	}
	detail := buildProductDetail(option_res, sku_res, sku_attr_res)
	result_summary := struct {
		Product  db.GetProductRow          `json:"product"`
		Brand    db.Brand                  `json:"brand"`
		Category db.Category               `json:"category"`
		Option   []services.OptionResponse `json:"option"`
		SKU      []services.SkuResponse    `json:"sku"`
	}{
		Product:  product_spu_detail,
		Brand:    brand,
		Category: category,
		Option:   detail.OptionMap,
		SKU:      detail.SKUs,
	}

	result := assets_services.NormalizeSQLNulls(result_summary, "data")

	//log.Printf("[GetProductWithID] Thành công lấy chi tiết sản phẩm '%s' (ID: %s) với %d SKU", product_spu_detail.Name, product_id, len(sku))
	return result, nil
}

func (s *service) GetDetailProduct(ctx context.Context, key string) (map[string]interface{}, *assets_services.ServiceError) {
	//log.Printf("[GetDetailProduct] Bắt đầu lấy chi tiết sản phẩm với key: %s", key)

	product_spu_detail, err := s.repository.GetProductByKey(ctx, key)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không tìm thấy sản phẩm với key: %s. Chi tiết: %v", key, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không tìm thấy sản phẩm với key: %s. Lỗi: %v", key, err))
	}

	// call sku
	sku, err := s.repository.ListSKUsByProduct(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể lấy danh sách SKU cho sản phẩm %s. Chi tiết: %v", key, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách SKU. Lỗi: %s", err.Error()))
	}
	sku_res := make([]services.ProductSku, len(sku))
	if err := copier.Copy(&sku_res, &sku); err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể sao chép dữ liệu SKU. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu SKU: %s", err.Error()))
	}

	// call option value
	option, err := s.repository.ListOptionValuesByProductID(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể lấy danh sách Option Values cho sản phẩm %s. Chi tiết: %v", key, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách thuộc tính sản phẩm. Lỗi: %s", err.Error()))
	}
	option_res := make([]services.OptionValue, len(option))
	for i, opt := range option {
		option_res[i] = services.OptionValue{
			ID:         opt.ID,
			OptionName: opt.OptionName,
			Value:      opt.Value,
			ProductID:  opt.ProductID,
			Image:      services.Narg[string]{Data: opt.Image.String, Valid: opt.Image.Valid},
		}
	}

	// call sku attr
	sku_attr, err := s.repository.ListSKUOptionValuesByProductID(ctx, product_spu_detail.ID)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể lấy thông tin liên kết SKU-Option cho sản phẩm %s. Chi tiết: %v", key, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin liên kết SKU. Lỗi: %s", err.Error()))
	}
	sku_attr_res := make([]services.SkuAttr, len(sku_attr))
	if err := copier.Copy(&sku_attr_res, &sku_attr); err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể sao chép dữ liệu SKU Attributes. Chi tiết: %v", err)
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu liên kết: %s", err.Error()))
	}

	// call brand name
	brand, err := s.repository.GetBrand(ctx, product_spu_detail.BrandID.String)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể lấy thông tin thương hiệu %s. Chi tiết: %v", product_spu_detail.BrandID.String, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin thương hiệu. Lỗi: %s", err.Error()))
	}

	// call category name
	category, err := s.repository.GetCategory(ctx, product_spu_detail.CategoryID)
	if err != nil {
		//log.Printf("[GetDetailProduct] LỖI: Không thể lấy thông tin danh mục %s. Chi tiết: %v", product_spu_detail.CategoryID, err)
		return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin danh mục. Lỗi: %s", err.Error()))
	}

	detail := buildProductDetail(option_res, sku_res, sku_attr_res)
	result_summary := struct {
		Product  db.GetProductByKeyRow     `json:"product"`
		Brand    db.Brand                  `json:"brand"`
		Category db.Category               `json:"category"`
		Option   []services.OptionResponse `json:"option"`
		SKU      []services.SkuResponse    `json:"sku"`
	}{
		Product:  product_spu_detail,
		Brand:    brand,
		Category: category,
		Option:   detail.OptionMap,
		SKU:      detail.SKUs,
	}

	result := assets_services.NormalizeSQLNulls(result_summary, "data")

	//log.Printf("[GetDetailProduct] Thành công lấy chi tiết sản phẩm '%s' (key: %s) với %d SKU", product_spu_detail.Name, key, len(sku))
	return result, nil
}
func (s *service) CreateProduct(ctx context.Context, token, userName string, product services.ProductParams, image *multipart.FileHeader, mediaFiles []*multipart.FileHeader, optionImages []struct {
	OptionName string
	Value      string
	Image      *multipart.FileHeader
}) *assets_services.ServiceError {
	//log.Printf("[CreateProduct] Bắt đầu tạo sản phẩm '%s' (key: %s) bởi người dùng: %s", product.Name, product.Key, userName)
	//log.Printf("[CreateProduct] Thông tin: %d Option Values, %d SKUs, %d Option Images", len(product.OptionValue), len(product.ProductSKU), len(optionImages))

	// Upload main image
	//log.Printf("[CreateProduct] Bắt đầu upload ảnh chính...")
	url_image, err := s.UploadMultiMedia(ctx, userName, []*multipart.FileHeader{image})
	if err != nil {
		//log.Printf("[CreateProduct] LỖI: Không thể upload ảnh chính. Chi tiết: %v", err)
		return assets_services.NewError(500, fmt.Errorf("không thể upload ảnh chính. Lỗi: %v", err))
	}
	if len(url_image) == 0 {
		//log.Printf("[CreateProduct] LỖI: Không có ảnh nào được upload")
		return assets_services.NewError(500, fmt.Errorf("không có ảnh nào được upload"))
	}
	//log.Printf("[CreateProduct] Upload ảnh chính thành công: %s", url_image[0])

	// Upload media images
	//log.Printf("[CreateProduct] Bắt đầu upload %d ảnh media...", len(mediaFiles))
	url_media, err := s.UploadMultiMedia(ctx, userName, mediaFiles)
	if err != nil {
		//log.Printf("[CreateProduct] LỖI: Không thể upload ảnh media. Chi tiết: %v", err)
		return assets_services.NewError(500, fmt.Errorf("không thể upload ảnh media. Lỗi: %v", err))
	}
	//log.Printf("[CreateProduct] Upload %d ảnh media thành công", len(url_media))

	url_media_json, errorsJson := json.Marshal(url_media)
	if errorsJson != nil {
		//log.Printf("[CreateProduct] LỖI: Không thể chuyển đổi danh sách ảnh media sang JSON. Chi tiết: %v", errorsJson)
		return assets_services.NewError(500, fmt.Errorf("lỗi xử lý dữ liệu ảnh media. Lỗi: %v", errorsJson))
	}

	// Upload option images
	option_image_path := make(map[string]map[string]string)
	option_image := []string{}
	//log.Printf("[CreateProduct] Bắt đầu upload %d ảnh cho options...", len(optionImages))

	for _, option := range optionImages {
		//log.Printf("[CreateProduct] Upload ảnh option %d/%d: %s = %s", i+1, len(optionImages), option.OptionName, option.Value)
		url_option, err := s.UploadMultiMedia(ctx, userName, []*multipart.FileHeader{option.Image})
		if err != nil {
			//log.Printf("[CreateProduct] LỖI: Không thể upload ảnh cho option '%s - %s'. Chi tiết: %v", option.OptionName, option.Value, err)
			return assets_services.NewError(500, fmt.Errorf("không thể upload ảnh cho option '%s - %s'. Lỗi: %v", option.OptionName, option.Value, err))
		}
		if option_image_path[option.OptionName] == nil {
			option_image_path[option.OptionName] = make(map[string]string)
		}
		option_image_path[option.OptionName][option.Value] = url_option[0]
		option_image = append(option_image, url_option[0])
		//log.Printf("[CreateProduct] Upload ảnh option thành công: %s", url_option[0])
	}

	//log.Printf("[CreateProduct] Bắt đầu transaction tạo sản phẩm trong database...")
	errors := s.repository.ExecTS(ctx, func(tx db.Querier) error {
		product_id := uuid.New().String()
		//log.Printf("[CreateProduct] Tạo product ID: %s", product_id)

		//log.Printf("[CreateProduct] Tạo bản ghi sản phẩm chính...")
		err := tx.CreateProduct(ctx, db.CreateProductParams{
			ID:                        product_id,
			Name:                      product.Name,
			Key:                       product.Key,
			Description:               sql.NullString{String: product.Description, Valid: true},
			ShortDescription:          sql.NullString{String: product.ShortDescription, Valid: true},
			BrandID:                   sql.NullString{String: product.BrandID, Valid: true},
			CategoryID:                product.CategoryID,
			ShopID:                    product.ShopID,
			Image:                     url_image[0],
			Media:                     sql.NullString{String: string(url_media_json), Valid: true},
			ProductIsPermissionReturn: product.ProductIsPermissionReturn,
			ProductIsPermissionCheck:  product.ProductIsPermissionCheck,
			CreateBy:                  sql.NullString{String: userName, Valid: true},
		})
		if err != nil {
			//log.Printf("[CreateProduct] LỖI: Không thể tạo bản ghi sản phẩm trong database. Chi tiết: %v", err)
			return fmt.Errorf("không thể tạo sản phẩm trong database: %w", err)
		}
		//log.Printf("[CreateProduct] Tạo bản ghi sản phẩm thành công")
		// create option value
		//log.Printf("[CreateProduct] Tạo %d Option Values...", len(product.OptionValue))
		createOptionValueParams := make(map[string]map[string]string)
		createSKUAtrrParams := []db.CreateSKUAttrParams{}

		for _, option := range product.OptionValue {
			optionValueID := uuid.New().String()
			option_value := db.CreateOptionValueParams{
				ID:         optionValueID,
				ProductID:  product_id,
				OptionName: option.OptionName,
				Value:      option.Value,
				Image:      sql.NullString{String: option_image_path[option.OptionName][option.Value], Valid: option_image_path[option.OptionName][option.Value] != ""},
			}
			err := tx.CreateOptionValue(ctx, option_value)
			if err != nil {
				//log.Printf("[CreateProduct] LỖI: Không thể tạo Option Value '%s - %s'. Chi tiết: %v", option.OptionName, option.Value, err)
				return fmt.Errorf("không thể tạo thuộc tính '%s - %s': %w", option.OptionName, option.Value, err)
			}
			if createOptionValueParams[option.OptionName] == nil {
				createOptionValueParams[option.OptionName] = make(map[string]string)
			}
			createOptionValueParams[option.OptionName][option.Value] = optionValueID
			//log.Printf("[CreateProduct] Tạo Option Value %d/%d: %s - %s (ID: %s)", i+1, len(product.OptionValue), option.OptionName, option.Value, optionValueID)
		}
		//log.Printf("[CreateProduct] Hoàn thành tạo %d Option Values", len(product.OptionValue))

		// create sku
		//log.Printf("[CreateProduct] Tạo %d Product SKUs...", len(product.ProductSKU))
		for _, sku := range product.ProductSKU {
			skuID := uuid.New().String()
			sku_value := db.CreateProductSKUParams{
				ID:        skuID,
				ProductID: product_id,
				SkuCode:   sku.SkuCode,
				Price:     sku.Price,
				Quantity:  sku.Quantity,
				Weight:    sku.Weight,
			}
			err := tx.CreateProductSKU(ctx, sku_value)
			if err != nil {
				//log.Printf("[CreateProduct] LỖI: Không thể tạo SKU '%s'. Chi tiết: %v", sku.SkuCode, err)
				return fmt.Errorf("không thể tạo SKU '%s': %w", sku.SkuCode, err)
			}
			//log.Printf("[CreateProduct] Tạo SKU %d/%d: %s (ID: %s, Giá: %.0f, Số lượng: %d)", i+1, len(product.ProductSKU), sku.SkuCode, skuID, sku.Price, sku.Quantity)

			// Link SKU với Option Values
			for _, skuOption := range sku.OptionValue {
				OptionValueID := createOptionValueParams[skuOption.OptionName][skuOption.Value]
				if OptionValueID == "" {
					//log.Printf("[CreateProduct] LỖI: Không tìm thấy Option Value '%s - %s' cho SKU '%s'", skuOption.OptionName, skuOption.Value, sku.SkuCode)
					return fmt.Errorf("không tìm thấy thuộc tính '%s - %s' cho SKU '%s'", skuOption.OptionName, skuOption.Value, sku.SkuCode)
				}
				createSKUAtrrParams = append(createSKUAtrrParams, db.CreateSKUAttrParams{
					SkuID:         skuID,
					ProductID:     product_id,
					OptionValueID: OptionValueID,
				})
			}
		}
		//log.Printf("[CreateProduct] Hoàn thành tạo %d SKUs", len(product.ProductSKU))
		// create sku option
		//log.Printf("[CreateProduct] Tạo %d liên kết SKU-Option...", len(createSKUAtrrParams))
		for i, skuAttr := range createSKUAtrrParams {
			err := tx.CreateSKUAttr(ctx, skuAttr)
			if err != nil {
				//log.Printf("[CreateProduct] LỖI: Không thể tạo liên kết SKU-Option (SKU: %s, Option: %s). Chi tiết: %v", skuAttr.SkuID, skuAttr.OptionValueID, err)
				return fmt.Errorf("không thể tạo liên kết SKU-Option: %w", err)
			}
			if (i+1)%10 == 0 || (i+1) == len(createSKUAtrrParams) {
				//log.Printf("[CreateProduct] Đã tạo %d/%d liên kết SKU-Option", i+1, len(createSKUAtrrParams))
			}
		}
		//log.Printf("[CreateProduct] Hoàn thành tạo tất cả liên kết SKU-Option")
		return nil
	})
	// check if user have permission to create product for this shop
	if errors != nil {
		//log.Printf("[CreateProduct] LỖI: Tạo sản phẩm thất bại. Chi tiết: %v", errors)
		// rm all image
		allImages := append(url_image, option_image...)
		allImages = append(allImages, url_media...)
		s.DeleteMultiImage(ctx, userName, allImages)
		return assets_services.NewError(400, fmt.Errorf("không thể tạo sản phẩm. Lỗi: %v", errors))
	}
	return nil
}

func (s *service) UpdateProduct(
	ctx context.Context,
	role_user, userName, productID string,

	product services.ProductUpdateParams, // Struct chứa dữ liệu JSON
	mainImage *multipart.FileHeader, // Ảnh chính mới (nếu có)
	newMediaFiles []*multipart.FileHeader, // Ảnh media mới (nếu có)
	optionImageUpdates []services.OptionImageUpdate, // Cập nhật ảnh option (nếu có)
) *assets_services.ServiceError {
	currentProduct, err := s.repository.GetProduct(ctx, productID) // Giả định hàm này trả về struct có Image, Media (string JSON), etc.
	if err != nil {
		if err == sql.ErrNoRows {
			return &assets_services.ServiceError{Code: 404, Err: fmt.Errorf("sản phẩm không tồn tại")} // Lỗi 404 nên được xử lý bên ngoài transaction
		}
		return assets_services.NewError(400, fmt.Errorf("lỗi khi lấy thông tin sản phẩm: %w", err))
	}
	// kiểm tra ngoại lệ nếu là delete status thì sẽ cập nhật lại sản phẩm trạng thái là xóa
	if product.DeleteStatus != nil && *product.DeleteStatus {
		err := s.repository.UpdateProduct(ctx, db.UpdateProductParams{
			ID:           productID,
			DeleteStatus: db.NullProductDeleteStatus{ProductDeleteStatus: db.ProductDeleteStatusDeleted, Valid: true},
			UpdateBy:     sql.NullString{String: userName, Valid: true},
		})
		if err != nil {
			return assets_services.NewError(400, fmt.Errorf("lỗi khi xóa sản phẩm: %w", err))
		}
		return assets_services.NewError(200, fmt.Errorf("xóa sản phẩm thành công"))
	}
	if product.ApprovalProduct != nil {
		if role_user != "ROLE_ADMIN" {
			return assets_services.NewError(403, fmt.Errorf("bạn không có quyền kiểm duyệt sản phẩm"))
		}
		ProductDeleteStatus := db.ProductDeleteStatusDeleted
		if *product.ApprovalProduct {
			if currentProduct.DeleteStatus.ProductDeleteStatus == db.ProductDeleteStatusActive {
				return assets_services.NewError(400, fmt.Errorf("sản phẩm đã được duyệt trước đó"))
			}
			if currentProduct.DeleteStatus.ProductDeleteStatus == db.ProductDeleteStatusDeleted {
				return assets_services.NewError(400, fmt.Errorf("sản phẩm đã bị xóa không thể duyệt"))
			}

			ProductDeleteStatus = db.ProductDeleteStatusActive

		}
		err := s.repository.UpdateProduct(ctx, db.UpdateProductParams{
			ID:           productID,
			DeleteStatus: db.NullProductDeleteStatus{ProductDeleteStatus: ProductDeleteStatus, Valid: true},
			UpdateBy:     sql.NullString{String: userName, Valid: true},
		})
		if err != nil {
			return assets_services.NewError(400, fmt.Errorf("lỗi khi xóa sản phẩm: %w", err))
		}
	}

	// ----- Bước 1: Upload tất cả ảnh mới LÊN TRƯỚC -----
	var newMainImageUrl string
	var newMediaUrls []string
	newOptionImageUrls := make(map[string]string) // map[optionValueID]newImageUrl
	imagesToDelete := make([]string, 0)           // Thu thập ID ảnh cũ cần xóa
	imagesToDeleteWhenFail := make([]string, 0)   // Thu thập ID ảnh mới đã upload để xóa khi thất bại
	// 1.1 Upload ảnh chính mới
	if mainImage != nil {
		uploadedImage, err := s.UploadMultiMedia(ctx, userName, []*multipart.FileHeader{mainImage}) // Giả định trả về URL
		if err != nil {
			return assets_services.NewError(400, fmt.Errorf("lỗi khi tải ảnh chính lên: %w", err))
		}
		newMainImageUrl = uploadedImage[0]
		imagesToDeleteWhenFail = append(imagesToDeleteWhenFail, newMainImageUrl)
		//log.Printf("Uploaded new main image: %s", newMainImageUrl)
	}

	// 1.2 Upload ảnh media mới
	if len(newMediaFiles) > 0 {
		uploadedMedia, err := s.UploadMultiMedia(ctx, userName, newMediaFiles) // Giả định trả về []string URLs
		if err != nil {
			return assets_services.NewError(400, fmt.Errorf("lỗi khi tải ảnh media lên: %w", err))
		}
		newMediaUrls = uploadedMedia
		imagesToDeleteWhenFail = append(imagesToDeleteWhenFail, newMediaUrls...)
		//log.Printf("Uploaded %d new media images.", len(newMediaUrls))
	}

	// 1.3 Upload ảnh option mới
	for _, optUpdate := range optionImageUpdates {
		if optUpdate.Image != nil && optUpdate.OptionValueID != "" {
			uploadedOptImage, err := s.UploadMultiMedia(ctx, userName, []*multipart.FileHeader{optUpdate.Image})
			if err != nil {
				return assets_services.NewError(400, fmt.Errorf("lỗi khi tải ảnh option lên cho %s: %w", optUpdate.OptionValueID, err))
			}
			newOptionImageUrls[optUpdate.OptionValueID] = uploadedOptImage[0]
			imagesToDeleteWhenFail = append(imagesToDeleteWhenFail, uploadedOptImage[0])
			//log.Printf("Uploaded new image for option value %s: %s", optUpdate.OptionValueID, uploadedOptImage)
		}
	}

	// ----- Bước 2: Thực hiện update trong transaction -----
	txErr := s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// 2.1 Lấy thông tin sản phẩm hiện tại (bao gồm ảnh)

		updateProductParams := db.UpdateProductParams{
			ID:       productID,
			UpdateBy: sql.NullString{String: userName, Valid: true},
		}
		productChanged := false // Cờ để kiểm tra xem có cần update product không

		// --- 2.2 Xử lý cập nhật ảnh chính ---
		currentMainImage := currentProduct.Image // Lấy URL ảnh chính hiện tại
		if newMainImageUrl != "" {               // Có ảnh mới được upload
			if currentMainImage != "" {
				imagesToDelete = append(imagesToDelete, currentMainImage) // Thêm ảnh cũ vào danh sách xóa
				//log.Printf("Marked old main image for deletion: %s", currentMainImage)
			}
			updateProductParams.Image = sql.NullString{String: newMainImageUrl, Valid: true}
			productChanged = true
		} else if product.RemoveMainImage != nil && *product.RemoveMainImage { // Yêu cầu xóa ảnh chính
			if currentMainImage != "" {
				imagesToDelete = append(imagesToDelete, currentMainImage) // Thêm ảnh cũ vào danh sách xóa
				//log.Printf("Marked main image for deletion: %s", currentMainImage)
			}
			updateProductParams.Image = sql.NullString{String: "", Valid: true} // Set rỗng/NULL trong DB
			productChanged = true
		} // Không có ảnh mới và không yêu cầu xóa -> giữ nguyên

		// --- 2.3 Xử lý cập nhật ảnh media ---
		currentMediaUrls := []string{}
		if currentProduct.Media.Valid && currentProduct.Media.String != "" && currentProduct.Media.String != "[]" {
			if err := json.Unmarshal([]byte(currentProduct.Media.String), &currentMediaUrls); err != nil {
				//log.Printf("Warning: Could not unmarshal current media JSON for product %s: %v", productID, err)
				// Có thể quyết định ghi đè hoặc giữ nguyên tùy vào logic của bạn
			}
		}

		finalMediaUrls := []string{}
		mediaChanged := false

		// Giữ lại các ảnh được yêu cầu
		if len(product.KeepMediaURLs) > 0 {
			keepMap := make(map[string]bool)
			for _, url := range product.KeepMediaURLs {
				keepMap[url] = true
			}
			for _, url := range currentMediaUrls {
				if keepMap[url] {
					finalMediaUrls = append(finalMediaUrls, url)
				} else {
					// Nếu không giữ và không nằm trong danh sách xóa tường minh, vẫn xóa
					if !containsString(product.RemoveMediaURLs, url) { // Giả định có hàm containsString
						imagesToDelete = append(imagesToDelete, url)
						//log.Printf("Marked media image for deletion (not kept): %s", url)
						mediaChanged = true
					}
				}
			}
		} else if len(product.RemoveMediaURLs) > 0 {
			// Nếu không có KeepMediaURLs nhưng có RemoveMediaURLs
			removeMap := make(map[string]bool)
			for _, url := range product.RemoveMediaURLs {
				removeMap[url] = true
			}
			for _, url := range currentMediaUrls {
				if !removeMap[url] {
					finalMediaUrls = append(finalMediaUrls, url) // Giữ lại những cái không bị yêu cầu xóa
				} else {
					imagesToDelete = append(imagesToDelete, url)
					//log.Printf("Marked media image for deletion (explicitly removed): %s", url)
					mediaChanged = true
				}
			}
		} else if len(newMediaUrls) > 0 {
			// Chỉ có ảnh mới, không có keep/remove -> Xóa hết ảnh cũ
			// for _, url := range currentMediaUrls {
			imagesToDelete = append(imagesToDelete, currentMediaUrls...)
			//log.Printf("Marked old media image for deletion (replaced by new): %s", url)
			// }
			mediaChanged = true // Vì ảnh cũ bị xóa
		} else {
			// Không có ảnh mới, không keep, không remove -> Giữ nguyên ảnh cũ
			finalMediaUrls = currentMediaUrls
		}

		// Thêm ảnh mới vào cuối
		if len(newMediaUrls) > 0 {
			finalMediaUrls = append(finalMediaUrls, newMediaUrls...)
			mediaChanged = true
		}

		if mediaChanged {
			finalMediaJsonBytes, _ := json.Marshal(finalMediaUrls)
			finalMediaJson := string(finalMediaJsonBytes)
			// Tránh lưu "null" nếu finalMediaUrls rỗng
			if finalMediaJson == "null" {
				finalMediaJson = "[]"
			}
			updateProductParams.Media = sql.NullString{String: finalMediaJson, Valid: true}
			productChanged = true
			//log.Printf("Final media JSON for product %s: %s", productID, finalMediaJson)
		}

		// --- 2.4 Cập nhật các trường thông tin cơ bản khác ---
		if product.Name != nil {
			updateProductParams.Name = sql.NullString{String: *product.Name, Valid: true}
			productChanged = true
		}
		if product.Key != nil {
			updateProductParams.Key = sql.NullString{String: *product.Key, Valid: true}
			productChanged = true
		}
		if product.Description != nil {
			updateProductParams.Description = sql.NullString{String: *product.Description, Valid: true}
			productChanged = true
		}
		if product.ShortDescription != nil {
			updateProductParams.ShortDescription = sql.NullString{String: *product.ShortDescription, Valid: true}
			productChanged = true
		}
		if product.ProductIsPermissionReturn != nil {
			updateProductParams.ProductIsPermissionReturn = sql.NullBool{Bool: *product.ProductIsPermissionReturn, Valid: true}
			productChanged = true
		}
		if product.ProductIsPermissionCheck != nil {
			updateProductParams.ProductIsPermissionCheck = sql.NullBool{Bool: *product.ProductIsPermissionCheck, Valid: true}
			productChanged = true
		}

		// Thực hiện update product nếu có thay đổi
		if productChanged {
			//log.Printf("Updating product basic info for %s", productID)
			err := tx.UpdateProduct(ctx, updateProductParams)
			if err != nil {
				return fmt.Errorf("lỗi khi cập nhật sản phẩm: %w", err)
			}
		}

		// --- 2.5 Cập nhật option_value và ảnh option ---
		// Lấy danh sách option values hiện tại của sản phẩm để lấy ảnh cũ
		currentOptionValues, err := tx.ListOptionValuesByProductID(ctx, productID) // Giả định có hàm này
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin option values hiện tại: %w", err)
		}
		currentOptionImages := make(map[string]string) // map[optionValueID]imageUrl
		for _, ov := range currentOptionValues {
			if ov.Image.Valid && ov.Image.String != "" {
				currentOptionImages[ov.ID] = ov.Image.String
			}
		}
		optionUpdateMap := make(map[string]services.OptionImageUpdate) // Để dễ tra cứu yêu cầu update ảnh option
		for _, ou := range optionImageUpdates {
			optionUpdateMap[ou.OptionValueID] = ou
		}

		for _, option := range product.OptionValue { // Lặp qua các option value cần update thông tin
			if option.ID == "" {
				// Bỏ qua nếu là tạo mới (logic tạo mới nên ở hàm khác)
				//log.Printf("Skipping option value update due to missing ID for %s - %s", option.OptionName, option.Value)
				continue
			}

			updateOptionParams := db.UpdateOptionValueParams{
				ID: option.ID,
			}
			optionValueChanged := false

			// Cập nhật tên và giá trị nếu có
			// không cho cập nhật optionname vì sẽ lỗi logic nhé.
			// if option.OptionName != "" { // Giả định chỉ update nếu có giá trị mới
			// 	updateOptionParams.OptionName = sql.NullString{String: option.OptionName, Valid: true}
			// 	optionValueChanged = true
			// }
			if option.Value != "" {
				updateOptionParams.Value = sql.NullString{String: option.Value, Valid: true}
				optionValueChanged = true
			}

			// Xử lý ảnh option
			currentOptImg := currentOptionImages[option.ID]
			imgUpdateReq, hasImgUpdateReq := optionUpdateMap[option.ID]
			newOptImgUrl := newOptionImageUrls[option.ID]

			if newOptImgUrl != "" { // Có ảnh mới
				if currentOptImg != "" {
					imagesToDelete = append(imagesToDelete, currentOptImg)
					//log.Printf("Marked old option image for deletion: %s (OptionValueID: %s)", currentOptImg, option.ID)
				}
				updateOptionParams.Image = sql.NullString{String: newOptImgUrl, Valid: true}
				optionValueChanged = true
			} else if hasImgUpdateReq && imgUpdateReq.Remove { // Yêu cầu xóa ảnh
				if currentOptImg != "" {
					imagesToDelete = append(imagesToDelete, currentOptImg)
					//log.Printf("Marked option image for deletion: %s (OptionValueID: %s)", currentOptImg, option.ID)
				}
				updateOptionParams.Image = sql.NullString{String: "", Valid: true} // Set rỗng/NULL
				optionValueChanged = true
			} // Không có ảnh mới và không yêu cầu xóa -> giữ nguyên ảnh option

			// Thực hiện update option value nếu có thay đổi
			if optionValueChanged {
				//log.Printf("Updating option value ID: %s", option.ID)
				err := tx.UpdateOptionValue(ctx, updateOptionParams)
				if err != nil {
					return fmt.Errorf("lỗi khi cập nhật option value ID: %s : %w", option.ID, err)
				}
			}
		}

		// --- 2.6 Cập nhật product_sku ---
		for _, sku := range product.ProductSKU {
			if sku.ID == "" {
				// Bỏ qua nếu là tạo mới
				//log.Printf("Skipping SKU update due to missing ID for SKU Code: %s", sku.SkuCode)
				continue
			}

			// Chỉ update các trường được cung cấp (dùng COALESCE hoặc kiểm tra nil)
			updateSkuParams := db.UpdateProductSKUParams{
				ID:       sku.ID,
				SkuCode:  sql.NullString{String: sku.SkuCode, Valid: sku.SkuCode != ""}, // Chỉ update nếu SkuCode không rỗng
				Price:    sql.NullFloat64{Float64: sku.Price, Valid: sku.Price != 0},    // Giả định luôn update Price, Quantity, Weight
				Quantity: sql.NullInt32{Int32: sku.Quantity, Valid: sku.Quantity != 0},
				Weight:   sql.NullFloat64{Float64: sku.Weight, Valid: sku.Weight != 0},
			}
			//log.Printf("Updating SKU ID: %s", sku.ID)
			err := tx.UpdateProductSKU(ctx, updateSkuParams)
			if err != nil {
				return fmt.Errorf("lỗi khi cập nhật SKU ID: %s : %w", sku.ID, err)
			}

			// TODO: Cập nhật bảng liên kết SKU và Option Values nếu cần (product_sku_attributes)
			// Logic này phụ thuộc vào thiết kế CSDL của bạn cho việc liên kết này.
			// Có thể cần xóa các liên kết cũ và tạo lại các liên kết mới dựa trên sku.OptionValue
		}

		return nil // Commit transaction
	})

	if txErr != nil {
		// Nếu lỗi xảy ra trong transaction, không cần gọi xóa ảnh
		//log.Printf("Transaction failed for UpdateProduct %s: %v", productID, txErr)
		// Xử lý lỗi 404 riêng nếu cần
		if strings.Contains(txErr.Error(), "product not found") {
			return assets_services.NewError(404, txErr)
		}
		// Xóa tất cả ảnh mới đã upload do transaction thất bại
		if len(imagesToDeleteWhenFail) > 0 {
			//log.Printf("Transaction failed, deleting %d newly uploaded images...", len(imagesToDeleteWhenFail))
			deleteFailErr := s.DeleteMultiImage(ctx, userName, imagesToDeleteWhenFail)
			if deleteFailErr != nil {
				//log.Printf("Error deleting newly uploaded images after transaction failure: %v", deleteFailErr)
			}
		}
		return assets_services.NewError(400, txErr)
	}

	// ----- Bước 3: Xóa ảnh cũ SAU KHI commit thành công -----
	if len(imagesToDelete) > 0 {
		//log.Printf("Attempting to delete %d old images for product %s...", len(imagesToDelete), productID)
		// Chuyển đổi URLs thành IDs nếu API yêu cầu ID
		// deleteImageIDs := convertUrlsToIDs(imagesToDelete) // Giả định bạn có hàm này
		deleteErr := s.DeleteMultiImage(ctx, userName, imagesToDelete) // Gọi API xóa ảnh
		if deleteErr != nil {
			// Ghi log lỗi xóa ảnh nhưng không trả lỗi cho client vì update CSDL đã thành công
			//log.Printf("Error deleting old images for product %s: %v. Images to delete: %v", productID, deleteErr, imagesToDelete)
		} else {
			//log.Printf("Successfully deleted %d old images.", len(imagesToDelete))
		}
	}

	return nil // Thành công
}

// containsString kiểm tra xem slice có chứa một string không
func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func (s *service) UpdateSKUReserverProduct(ctx context.Context, productSKU []services.ProductUpdateSKUReserver, type_req services.ProductUpdateType) *assets_services.ServiceError {
	err := s.repository.ExecTS(ctx, func(tx db.Querier) error {
		for _, sku := range productSKU {
			sku_db, err := tx.GetProductSKU(ctx, sku.SkuID)
			if err != nil {
				return fmt.Errorf("không tìm thấy SKU với ID: %s. Lỗi: %w", sku.SkuID, err)
			}

			switch {
			case type_req == services.HOLD:
				sku_db.QuantityReserver += sku.QuantityReserver
				if sku_db.Quantity-sku_db.QuantityReserver < 0 {
					return fmt.Errorf("không đủ số lượng tồn kho cho SKU %s (Còn: %d, Yêu cầu: %d)", sku.SkuID, sku_db.Quantity-sku_db.QuantityReserver+sku.QuantityReserver, sku.QuantityReserver)
				}

				err = tx.UpdateProductSKU(ctx, db.UpdateProductSKUParams{
					ID:               sku.SkuID,
					QuantityReserver: sql.NullInt32{Int32: sku_db.QuantityReserver, Valid: true},
				})
				if err != nil {
					return fmt.Errorf("không thể cập nhật số lượng đặt trước cho SKU: %w", err)
				}

			case type_req == services.COMMIT:
				new_QuantityReserver := sku_db.QuantityReserver - sku.QuantityReserver
				new_Quantity := sku_db.Quantity - sku.QuantityReserver
				if new_QuantityReserver < 0 || new_Quantity < 0 {
					return fmt.Errorf("dữ liệu không hợp lệ khi xác nhận đơn hàng cho SKU %s", sku.SkuID)
				}
				err = tx.UpdateProductSKU(ctx, db.UpdateProductSKUParams{
					ID:               sku.SkuID,
					QuantityReserver: sql.NullInt32{Int32: new_QuantityReserver, Valid: true},
					Quantity:         sql.NullInt32{Int32: new_Quantity, Valid: true},
				})
				if err != nil {
					return fmt.Errorf("không thể xác nhận đơn hàng cho SKU: %w", err)
				}
				// trường hợp cập nhật xác nhận sản phẩm thì sẽ cộng số lượng mua nó vào trường số lượng đã bán
				err = tx.IncrementProductTotalSold(ctx, db.IncrementProductTotalSoldParams{
					Quantity: int64(sku.QuantityReserver),
					ID:       sku_db.ProductID,
				})
				if err != nil {
					return fmt.Errorf("không thể cập nhật thêm vào số lượng bán hàng cho shop.: %w", err)
				}

			case type_req == services.ROLLBACK:
				sku_db.QuantityReserver -= sku.QuantityReserver
				if sku_db.QuantityReserver < 0 {
					return fmt.Errorf("dữ liệu không hợp lệ khi hoàn tác đơn hàng cho SKU %s", sku.SkuID)
				}
				err = tx.UpdateProductSKU(ctx, db.UpdateProductSKUParams{
					ID:               sku.SkuID,
					QuantityReserver: sql.NullInt32{Int32: sku_db.QuantityReserver, Valid: true},
				})
				if err != nil {
					return fmt.Errorf("không thể hoàn tác đơn hàng cho SKU: %w", err)
				}
			default:
				return fmt.Errorf("loại cập nhật không hợp lệ: %v", type_req)
			}
		}
		return nil
	})

	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("không thể cập nhật số lượng đặt trước: %w", err))
	}

	return nil
}

func buildProductDetail(options []services.OptionValue, skus []services.ProductSku, attrs []services.SkuAttr) services.ProductDetailResponse {
	result := services.ProductDetailResponse{}

	// 🔹 1. Gom nhóm option_value theo OptionName
	optionMap := make(map[string][]services.OptionValueItem)
	for _, opt := range options {
		var img *string
		if opt.Image.Valid {
			img = &opt.Image.Data
		}
		optionMap[opt.OptionName] = append(optionMap[opt.OptionName], services.OptionValueItem{
			Value:         opt.Value,
			Image:         img,
			OptionValueID: opt.ID,
		})
	}

	for name, values := range optionMap {
		result.OptionMap = append(result.OptionMap, services.OptionResponse{
			OptionName: name,
			Values:     values,
		})
	}

	// 🔹 2. Gom nhóm option_value_id cho từng SKU
	skuMap := make(map[string][]string)
	for _, a := range attrs {
		skuMap[a.SkuID] = append(skuMap[a.SkuID], a.OptionValueID)
	}

	// 🔹 3. Ghép thông tin SKU với danh sách option_value_id
	for _, sku := range skus {
		result.SKUs = append(result.SKUs, services.SkuResponse{
			ID:             sku.ID,
			SkuCode:        sku.SkuCode,
			Price:          sku.Price,
			Quantity:       sku.Quantity - sku.QuantityReserver,
			Weight:         sku.Weight,
			OptionValueIDs: skuMap[sku.ID],
			SkuName:        sku.SkuName,
		})
	}

	return result
}

// BuildProductSearchString tạo chuỗi text đầy đủ từ thông tin sản phẩm để sử dụng cho semantic search
// Hàm này gom tất cả thông tin quan trọng của sản phẩm thành 1 chuỗi có cấu trúc
func (s *service) BuildProductSearchString(ctx context.Context, productID string) (string, error) {
	//log.Printf("[BuildProductSearchString] Bắt đầu build search string cho product ID: %s", productID)

	// 1. Lấy thông tin sản phẩm chính
	product, err := s.repository.GetProduct(ctx, productID)
	if err != nil {
		//log.Printf("[BuildProductSearchString] LỖI: Không tìm thấy sản phẩm ID: %s. Chi tiết: %v", productID, err)
		return "", fmt.Errorf("không tìm thấy sản phẩm: %w", err)
	}

	var searchParts []string

	// 2. Thêm tên sản phẩm (trọng số cao nhất)
	searchParts = append(searchParts, fmt.Sprintf("Tên: %s", product.Name))

	// 3. Thêm key sản phẩm
	// searchParts = append(searchParts, fmt.Sprintf("Key: %s", product.Key))

	// 4. Thêm mô tả ngắn
	if product.ShortDescription.Valid && product.ShortDescription.String != "" {
		searchParts = append(searchParts, fmt.Sprintf("Mô tả ngắn: %s", product.ShortDescription.String))
	}

	// // 5. Thêm mô tả chi tiết
	// if product.Description.Valid && product.Description.String != "" {
	// 	searchParts = append(searchParts, fmt.Sprintf("Mô tả: %s", product.Description.String))
	// }

	// 6. Lấy và thêm thông tin thương hiệu
	if product.BrandID.Valid && product.BrandID.String != "" {
		brand, err := s.repository.GetBrand(ctx, product.BrandID.String)
		if err == nil {
			searchParts = append(searchParts, fmt.Sprintf("Thương hiệu: %s (Mã: %s)", brand.Name, brand.Code))
		}
	}

	// 7. Lấy và thêm thông tin danh mục
	category, err := s.repository.GetCategory(ctx, product.CategoryID)
	if err == nil {
		searchParts = append(searchParts, fmt.Sprintf("Danh mục: %s (Path: %s)", category.Name, category.Path))
	}

	// 8. Lấy và thêm thông tin Option Values
	options, err := s.repository.ListOptionValuesByProductID(ctx, productID)
	if err == nil && len(options) > 0 {
		// Gom nhóm options theo tên
		optionGroups := make(map[string][]string)
		for _, opt := range options {
			optionGroups[opt.OptionName] = append(optionGroups[opt.OptionName], opt.Value)
		}

		// Build chuỗi options
		for optionName, values := range optionGroups {
			searchParts = append(searchParts, fmt.Sprintf("%s: %s", optionName, strings.Join(values, ", ")))
		}
	}

	// 9. Lấy và thêm thông tin SKU
	skus, err := s.repository.ListSKUsByProduct(ctx, productID)
	if err == nil && len(skus) > 0 {
		// Thông tin giá
		var minPrice, maxPrice float64
		minPrice = skus[0].Price
		maxPrice = skus[0].Price

		// skuCodes := make([]string, 0, len(skus))
		// skuNames := make([]string, 0, len(skus))

		for _, sku := range skus {
			if sku.Price < minPrice {
				minPrice = sku.Price
			}
			if sku.Price > maxPrice {
				maxPrice = sku.Price
			}
			// skuCodes = append(skuCodes, sku.SkuCode)
			// if sku.SkuName.Valid && sku.SkuName.String != "" {
			// 	skuNames = append(skuNames, sku.SkuName.String)
			// }
		}

		// Thêm thông tin giá
		if minPrice == maxPrice {
			searchParts = append(searchParts, fmt.Sprintf("Giá: %.0f VNĐ", minPrice))
		} else {
			searchParts = append(searchParts, fmt.Sprintf("Giá: %.0f - %.0f VNĐ", minPrice, maxPrice))
		}

		// Thêm mã SKU
		// searchParts = append(searchParts, fmt.Sprintf("Mã SKU: %s", strings.Join(skuCodes, ", ")))

		// Thêm tên SKU (nếu có)
		// if len(skuNames) > 0 {
		// searchParts = append(searchParts, fmt.Sprintf("Phân loại bao gồm: %s", strings.Join(skuNames, ", ")))
		// }

		// Thêm số lượng biến thể
		// searchParts = append(searchParts, fmt.Sprintf("Số lượng biến thể: %d", len(skus)))
	}

	// 10. Ghép tất cả thành chuỗi cuối cùng
	searchString := strings.Join(searchParts, ". ")

	//log.Printf("[BuildProductSearchString] Hoàn thành build search string cho product ID: %s (Length: %d)", productID, len(searchString))
	return searchString, nil
}
func (s *service) GetALLProductID(ctx context.Context) ([]string, *assets_services.ServiceError) {
	product, err := s.repository.GetAllProductID(ctx)
	if err != nil {
		return nil, assets_services.NewError(400, err)
	}
	return product, nil
}

func (s *service) GetListProductWithIDs(ctx context.Context, productID []string) (map[string]interface{}, *assets_services.ServiceError) {

	products, err := s.repository.GetProductIDs(ctx, productID)
	if err != nil {
		return nil, assets_services.NewError(400, err)
	}
	product_detail := []interface{}{}
	for _, product_spu_detail := range products {
		// call sku
		sku, err := s.repository.ListSKUsByProduct(ctx, product_spu_detail.ID)
		if err != nil {
			//log.Printf("[GetDetailProduct] LỖI: Không thể lấy danh sách SKU cho sản phẩm %s. Chi tiết: %v", key, err)
			return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy danh sách SKU. Lỗi: %s", err.Error()))
		}
		sku_res := make([]services.ProductSkuSearch, len(sku))
		if err := copier.Copy(&sku_res, &sku); err != nil {
			//log.Printf("[GetDetailProduct] LỖI: Không thể sao chép dữ liệu SKU. Chi tiết: %v", err)
			return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu SKU: %s", err.Error()))
		}
		// call brand name
		brand, err := s.repository.GetBrand(ctx, product_spu_detail.BrandID.String)
		if err != nil {
			//log.Printf("[GetDetailProduct] LỖI: Không thể lấy thông tin thương hiệu %s. Chi tiết: %v", product_spu_detail.BrandID.String, err)
			return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin thương hiệu. Lỗi: %s", err.Error()))
		}

		// call category name
		category, err := s.repository.GetCategory(ctx, product_spu_detail.CategoryID)
		if err != nil {
			//log.Printf("[GetDetailProduct] LỖI: Không thể lấy thông tin danh mục %s. Chi tiết: %v", product_spu_detail.CategoryID, err)
			return nil, assets_services.NewError(400, fmt.Errorf("không thể lấy thông tin danh mục. Lỗi: %s", err.Error()))
		}

		// remove field detail
		product_search := services.ProductForSearch{}
		if err := copier.Copy(&product_search, &product_spu_detail); err != nil {
			//log.Printf("[GetDetailProduct] LỖI: Không thể sao chép dữ liệu sản phẩm. Chi tiết: %v", err)
			return nil, assets_services.NewError(400, fmt.Errorf("lỗi xử lý dữ liệu sản phẩm: %s", err.Error()))
		}
		result_summary := struct {
			Product  services.ProductForSearch `json:"product"`
			Brand    string                    `json:"brand"`
			Category string                    `json:"category"`
			// Option   []services.OptionResponse `json:"option"`
			SKU []services.ProductSkuSearch `json:"sku"`
		}{
			Product:  product_search,
			Brand:    brand.Name,
			Category: category.Name,
			// Option:   detail.OptionMap,
			SKU: sku_res,
		}
		product_detail = append(product_detail, result_summary)
	}
	result := assets_services.NormalizeListSQLNulls(product_detail, "data")
	return result, nil
}

// BuildProductSearchStringFromParams tạo search string trực tiếp từ params khi tạo sản phẩm
// Sử dụng hàm này ngay sau khi CreateProduct thành công để không cần query lại DB
func BuildProductSearchStringFromParams(
	product services.ProductParams,
	brandName, brandCode string,
	categoryName, categoryPath string,
	skuInfos []struct {
		SkuCode string
		SkuName string
		Price   float64
	},
) string {
	var searchParts []string

	// 1. Tên sản phẩm
	searchParts = append(searchParts, fmt.Sprintf("Tên: %s", product.Name))

	// 2. Key sản phẩm
	searchParts = append(searchParts, fmt.Sprintf("Key: %s", product.Key))

	// 3. Mô tả ngắn
	if product.ShortDescription != "" {
		searchParts = append(searchParts, fmt.Sprintf("Mô tả ngắn: %s", product.ShortDescription))
	}

	// // 4. Mô tả chi tiết
	// if product.Description != "" {
	// 	searchParts = append(searchParts, fmt.Sprintf("Mô tả: %s", product.Description))
	// }

	// 5. Thương hiệu
	if brandName != "" {
		searchParts = append(searchParts, fmt.Sprintf("Thương hiệu: %s (Mã: %s)", brandName, brandCode))
	}

	// 6. Danh mục
	if categoryName != "" {
		searchParts = append(searchParts, fmt.Sprintf("Danh mục: %s (Path: %s)", categoryName, categoryPath))
	}

	// 7. Option Values
	if len(product.OptionValue) > 0 {
		optionGroups := make(map[string][]string)
		for _, opt := range product.OptionValue {
			optionGroups[opt.OptionName] = append(optionGroups[opt.OptionName], opt.Value)
		}

		for optionName, values := range optionGroups {
			searchParts = append(searchParts, fmt.Sprintf("%s: %s", optionName, strings.Join(values, ", ")))
		}
	}

	// 8. Thông tin SKU
	if len(skuInfos) > 0 {
		var minPrice, maxPrice float64
		minPrice = skuInfos[0].Price
		maxPrice = skuInfos[0].Price

		skuCodes := make([]string, 0, len(skuInfos))
		skuNames := make([]string, 0, len(skuInfos))

		for _, sku := range skuInfos {
			if sku.Price < minPrice {
				minPrice = sku.Price
			}
			if sku.Price > maxPrice {
				maxPrice = sku.Price
			}
			skuCodes = append(skuCodes, sku.SkuCode)
			if sku.SkuName != "" {
				skuNames = append(skuNames, sku.SkuName)
			}
		}

		// Giá
		if minPrice == maxPrice {
			searchParts = append(searchParts, fmt.Sprintf("Giá: %.0f VNĐ", minPrice))
		} else {
			searchParts = append(searchParts, fmt.Sprintf("Giá: %.0f - %.0f VNĐ", minPrice, maxPrice))
		}

		// Mã SKU
		searchParts = append(searchParts, fmt.Sprintf("Mã SKU: %s", strings.Join(skuCodes, ", ")))

		// Tên SKU
		if len(skuNames) > 0 {
			searchParts = append(searchParts, fmt.Sprintf("Phân loại: %s", strings.Join(skuNames, ", ")))
		}

		// Số lượng biến thể
		searchParts = append(searchParts, fmt.Sprintf("Số lượng biến thể: %d", len(skuInfos)))
	}

	return strings.Join(searchParts, ". ")
}
