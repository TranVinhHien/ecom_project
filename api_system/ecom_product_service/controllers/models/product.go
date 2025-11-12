package controllers_model

type Product struct {
	Name                      string        `form:"name" json:"name" binding:"required,min=1,max=255,required"`
	Key                       string        `form:"key" json:"key" binding:"required"`
	Description               string        `form:"description" json:"description" binding:"required,required"`
	ShortDescription          string        `form:"short_description" json:"short_description" binding:"required"`
	BrandID                   string        `form:"brand_id" json:"brand_id" binding:"omitempty,uuid"`
	CategoryID                string        `form:"category_id" json:"category_id" binding:"required,uuid"`
	ShopID                    string        `form:"shop_id" json:"shop_id" binding:"required,uuid"`
	ProductIsPermissionReturn bool          `form:"product_is_permission_return" json:"product_is_permission_return" binding:"omitempty"`
	ProductIsPermissionCheck  bool          `form:"product_is_permission_check" json:"product_is_permission_check" binding:"omitempty"`
	ProductSKU                []ProductSku  `form:"product_sku" json:"product_sku" binding:"dive,required"`
	OptionValue               []OptionValue `form:"option_value" json:"option_value" binding:"dive,required"`
}

type ProductUpdate struct {
	Name                      *string       `form:"name" json:"name" binding:"omitempty,min=1,max=255"`
	Key                       *string       `form:"key" json:"key" binding:"omitempty,min=1"`
	Description               *string       `form:"description" json:"description"`
	ShortDescription          *string       `form:"short_description" json:"short_description"`
	ProductIsPermissionReturn *bool         `form:"product_is_permission_return" json:"product_is_permission_return"`
	ProductIsPermissionCheck  *bool         `form:"product_is_permission_check" json:"product_is_permission_check"`
	DeleteStatus              *bool         `form:"delete_status" json:"delete_status" `
	ProductSKU                []ProductSku  `form:"product_sku" json:"product_sku" binding:"omitempty,dive"`
	OptionValue               []OptionValue `form:"option_value" json:"option_value" binding:"omitempty,dive"`
}

type ProductSku struct {
	ID          string        `form:"id" json:"id" binding:"required"`
	SkuCode     string        `form:"sku_code" json:"sku_code" binding:"required"`
	Price       float64       `form:"price" json:"price" binding:"required,min=0"`
	Quantity    int32         `form:"quantity" json:"quantity" binding:"required,min=0"`
	Weight      float64       `form:"weight" json:"weight" binding:"required,min=0"`
	OptionValue []OptionValue `form:"option_value" json:"option_value" binding:"required"`
}
type OptionValue struct {
	ID         string `form:"id" json:"id" binding:"required"`
	OptionName string `form:"option_name" json:"option_name" binding:"required,min=1,max=100,required"`
	Value      string `form:"value" json:"value" binding:"required,min=1,max=255,required"`
}

type SkuAttr struct {
	SkuID         string `form:"sku_id" json:"sku_id" binding:"required,uuid"`
	OptionValueID string `form:"option_value_id" json:"option_value_id" binding:"required,uuid"`
}

type ProductUpdateSKUReserver struct {
	SkuID            string `json:"sku_id" binding:"required,uuid"`
	QuantityReserver int32  `json:"quantity_reserver" binding:"required,gt=0"`
}

type UpdateSKUReserverRequest struct {
	Data   []ProductUpdateSKUReserver `json:"data" binding:"required,dive"`
	Status string                     `json:"status" binding:"required,oneof=commit hold rollback"`
}
