package services

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

type ProductStatus string

const (
	ProductDeleteStatusActive  ProductStatus = "Active"
	ProductDeleteStatusDeleted ProductStatus = "Deleted"
)

type Condition struct {
	Field    string      // Tên cột
	Operator string      // Dấu so sánh (>, <, =, >=, <=, !=)
	Value    interface{} // Giá trị so sánh
}
type OrderBy struct {
	Field string         // Tên cột
	Value OrderDirection // Giá trị so sánh
}
type QueryFilter struct {
	Conditions []Condition
	OrderBy    *OrderBy // Trường để sắp xếp
	Page       int      // Trang hiện tại
	PageSize   int      // Số lượng kết quả mỗi trang
}

func NewQueryFilter(page int, pageSize int, conditions []Condition, orderBy *OrderBy) QueryFilter {
	if pageSize > 100 {
		pageSize = 100
	}

	return QueryFilter{
		Page:       page,
		PageSize:   pageSize,
		Conditions: conditions,
		OrderBy:    orderBy,
	}
}

type Categorys struct {
	CategoryID string            `json:"category_id"`
	Name       string            `json:"name"`
	Key        string            `json:"key"`
	Path       string            `json:"path"`
	Childs     Narg[[]Categorys] `json:"child"`
	Parent     Narg[string]      `json:"parent"`
	Image      Narg[string]      `json:"image"`
}

type Product struct {
	ID                        string          `json:"id"`
	Name                      string          `json:"name"`
	Key                       string          `json:"key"`
	Description               string          `json:"description"`
	ShortDescription          string          `json:"short_description"`
	BrandID                   string          `json:"brand_id"`
	CategoryID                string          `json:"category_id"`
	ShopID                    string          `json:"shop_id"`
	Image                     string          `json:"image"`
	Media                     []string        `json:"media"`
	DeleteStatus              ProductStatus   `json:"delete_status"`
	ProductIsPermissionReturn bool            `json:"product_is_permission_return"`
	ProductIsPermissionCheck  bool            `json:"product_is_permission_check"`
	CreateDate                time.Time       `json:"create_date"`
	UpdateDate                Narg[time.Time] `json:"update_date"`
	CreateBy                  string          `json:"create_by"`
	UpdateBy                  Narg[string]    `json:"update_by"`
	ProductSKU                []ProductSku    `json:"product_sku"`
	SKUAttr                   []SkuAttr       `json:"sku_attr"`
	OptionValue               []OptionValue   `json:"option_value"`
}
type ProductSku struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	SkuCode          string    `json:"sku_code"`
	Price            float64   `json:"price"`
	Quantity         int32     `json:"quantity"`
	QuantityReserver int32     `json:"quantity_reserver"`
	Weight           float64   `json:"weight"`
	CreateDate       time.Time `json:"create_date"`
	UpdateDate       time.Time `json:"update_date"`
	SkuName          string    `json:"sku_name"`
}

type ProductSkuSearch struct {
	ID string `json:"id"`
	// ProductID        string    `json:"product_id"`
	// SkuCode          string    `json:"sku_code"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
	// QuantityReserver int32   `json:"quantity_reserver"`
	// Weight float64 `json:"weight"`
	// CreateDate       time.Time `json:"create_date"`
	// UpdateDate       time.Time `json:"update_date"`
	SkuName string `json:"sku_name"`
}
type SkuAttr struct {
	SkuID         string `json:"sku_id"`
	OptionValueID string `json:"option_value_id"`
}

type OptionValue struct {
	ID         string       `json:"id"`
	OptionName string       `json:"option_name"`
	Value      string       `json:"value"`
	ProductID  string       `json:"product_id"`
	Image      Narg[string] `json:"image"`
}

type ProductParams struct {
	Name                      string                `json:"name" `
	Key                       string                `json:"key" `
	Description               string                `json:"description" `
	ShortDescription          string                `json:"short_description" `
	BrandID                   string                `json:"brand_id" `
	CategoryID                string                `json:"category_id" `
	ShopID                    string                `json:"shop_id" `
	ProductIsPermissionReturn bool                  `json:"product_is_permission_return" `
	ProductIsPermissionCheck  bool                  `json:"product_is_permission_check" `
	ProductSKU                []ProductSKUParams    `json:"product_sku" `
	OptionValue               []ProductOptionParams `json:"option_value" `
}

type ProductSKUAttrParams struct {
	SkuID         string `json:"sku_id" `
	OptionValueID string `json:"option_value_id" `
}

// product res
type OptionResponse struct {
	OptionName string            `json:"option_name"`
	Values     []OptionValueItem `json:"values"`
}

type OptionValueItem struct {
	Value         string  `json:"value"`
	Image         *string `json:"image"`
	OptionValueID string  `json:"option_value_id"`
}

type SkuResponse struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	SkuCode          string    `json:"sku_code"`
	Price            float64   `json:"price"`
	Quantity         int32     `json:"quantity"`
	QuantityReserver int32     `json:"quantity_reserver"`
	Weight           float64   `json:"weight"`
	CreateDate       time.Time `json:"create_date"`
	UpdateDate       time.Time `json:"update_date"`
	SkuName          string    `json:"sku_name"`
	OptionValueIDs   []string  `json:"option_value_ids"`
}

type ProductDetailResponse struct {
	OptionMap []OptionResponse `json:"option_map"`
	SKUs      []SkuResponse    `json:"skus"`
}

// ///////////////////////
// service update sku reserver product
// ///////////////////////
type ProductUpdateSKUReserver struct {
	SkuID            string `json:"sku_id" `
	QuantityReserver int32  `json:"quantity_reserver"`
}

type ProductUpdateType string

const (
	COMMIT   ProductUpdateType = "commit"
	HOLD     ProductUpdateType = "hold"
	ROLLBACK ProductUpdateType = "rollback"
)

// ProductUpdateParams cần thêm các trường quản lý ảnh media và ảnh chính
type ProductUpdateParams struct {
	Name                      *string               `json:"name,omitempty"`
	Key                       *string               `json:"key,omitempty"`
	Description               *string               `json:"description,omitempty"`
	ShortDescription          *string               `json:"short_description,omitempty"`
	ProductIsPermissionReturn *bool                 `json:"product_is_permission_return,omitempty"`
	ProductIsPermissionCheck  *bool                 `json:"product_is_permission_check,omitempty"`
	ProductSKU                []ProductSKUParams    `json:"product_sku,omitempty"`
	OptionValue               []ProductOptionParams `json:"option_value,omitempty"`
	DeleteStatus              *bool                 `json:"delete_status" `

	// --- Cập nhật quản lý ảnh ---
	RemoveMainImage *bool    `json:"remove_main_image,omitempty"` // Cờ để xóa ảnh chính
	KeepMediaURLs   []string `json:"keep_media_urls,omitempty"`   // Giữ lại media URLs này
	RemoveMediaURLs []string `json:"remove_media_urls,omitempty"` // Xóa các media URLs này
	// NewMediaFiles []*multipart.FileHeader `json:"-"` // Sẽ nhận từ form-data, không phải JSON
}

type ProductSKUParams struct {
	ID          string                `json:"id,omitempty"` // Bắt buộc khi update
	SkuCode     string                `json:"sku_code"`
	Price       float64               `json:"price"`
	Quantity    int32                 `json:"quantity"`
	Weight      float64               `json:"weight"`
	OptionValue []ProductOptionParams `json:"option_value"` // Thường không cần ở đây khi update SKU
}

type ProductOptionParams struct {
	ID         string `json:"id,omitempty"` // Bắt buộc khi update option value
	OptionName string `json:"option_name"`
	Value      string `json:"value"`
	// Image string `json:"image,omitempty"` // Không truyền image URL qua đây
}

// OptionImageUpdate dùng cho input form-data
type OptionImageUpdate struct {
	OptionValueID string                `form:"option_value_id"` // ID của OptionValue cần cập nhật ảnh
	Image         *multipart.FileHeader `form:"image"`           // Ảnh mới nếu có
	Remove        bool                  `form:"remove"`          // Đánh dấu nếu muốn xóa ảnh hiện tại
}

type ProductForSearch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	// Description               sql.NullString          `json:"description"`
	ShortDescription sql.NullString `json:"short_description"`
	// BrandID                   sql.NullString          `json:"brand_id"`
	// CategoryID                string                  `json:"category_id"`
	// ShopID                    string                  `json:"shop_id"`
	Image string `json:"image"`
	// Media                     sql.NullString          `json:"media"`
	// DeleteStatus              NullProductDeleteStatus `json:"delete_status"`
	ProductIsPermissionReturn sql.NullBool `json:"product_is_permission_return"`
	ProductIsPermissionCheck  sql.NullBool `json:"product_is_permission_check"`
	// CreateDate                sql.NullTime            `json:"create_date"`
	// UpdateDate                sql.NullTime            `json:"update_date"`
	// CreateBy                  sql.NullString          `json:"create_by"`
	// UpdateBy                  sql.NullString          `json:"update_by"`
}

// ProductRating represents rating information for a product
type ProductRating struct {
	ProductID     string  `json:"product_id"`
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

// ProductRatingStatsItem represents rating statistics for a single product
type ProductTotalSold struct {
	ProductID string `json:"product_id"`
	TotalSold int64  `json:"total_sold"`
}
