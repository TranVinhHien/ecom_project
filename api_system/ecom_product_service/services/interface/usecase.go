package iservices

import (
	"context"
	"mime/multipart"

	assets_services "github.com/TranVinhHien/ecom_product_service/services/assets"
	services "github.com/TranVinhHien/ecom_product_service/services/entity"
)

type Categories interface {
	GetCategoris(ctx context.Context, cate_id string) (map[string]interface{}, *assets_services.ServiceError)
	AddCategory(ctx context.Context, userName string, cat services.Categorys, file *multipart.FileHeader) *assets_services.ServiceError
	UpdateCategory(ctx context.Context, userName string, cat services.Categorys, file *multipart.FileHeader) *assets_services.ServiceError
	DeleteCategory(ctx context.Context, userName, categoryID string) *assets_services.ServiceError
}
type Products interface {
	UpdateSKUReserverProduct(ctx context.Context, productSKU []services.ProductUpdateSKUReserver, type_req services.ProductUpdateType) *assets_services.ServiceError
	GetAllProductSimple(ctx context.Context, query services.QueryFilter, category_path, brand_code, shop_id, keywords, sort string, min_price, max_price float64) (map[string]interface{}, *assets_services.ServiceError)
	GetDetailProduct(ctx context.Context, productSpuID string) (map[string]interface{}, *assets_services.ServiceError)
	CreateProduct(ctx context.Context, token, userName string, product services.ProductParams, image *multipart.FileHeader, mediaFiles []*multipart.FileHeader, optionImages []struct {
		OptionName string
		Value      string
		Image      *multipart.FileHeader
	}) *assets_services.ServiceError
	UpdateProduct(
		ctx context.Context,
		token, userName, productID string,
		product services.ProductUpdateParams, // Struct chứa dữ liệu JSON
		mainImage *multipart.FileHeader, // Ảnh chính mới (nếu có)
		newMediaFiles []*multipart.FileHeader, // Ảnh media mới (nếu có)
		optionImageUpdates []services.OptionImageUpdate, // Cập nhật ảnh option (nếu có)

	) *assets_services.ServiceError
	GetSKUProduct(ctx context.Context, product_sku_id string) (map[string]interface{}, *assets_services.ServiceError)
	GetProductWithID(ctx context.Context, product_id string) (map[string]interface{}, *assets_services.ServiceError)
	BuildProductSearchString(ctx context.Context, productID string) (string, error)
	GetALLProductID(ctx context.Context) ([]string, *assets_services.ServiceError)
	GetListProductWithIDs(ctx context.Context, productID []string) (map[string]interface{}, *assets_services.ServiceError)
}
type Media interface {
	RenderImage(ctx context.Context, id string) string
}
