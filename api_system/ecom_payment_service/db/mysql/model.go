package db

import (
	"database/sql"
)

type ProductSkusDetailss struct {
	ProductSkuID     string        `json:"product_sku_id"`
	Value            string        `json:"value"`
	SkuStock         sql.NullInt32 `json:"sku_stock"`
	Price            float64       `json:"price"`
	Sort             sql.NullInt32 `json:"sort"`
	CreateDate       sql.NullTime  `json:"create_date"`
	UpdateDate       sql.NullTime  `json:"update_date"`
	ProductsSpuID    string        `json:"products_spu_id"`
	Name             string        `json:"name"`
	ShortDescription string        `json:"short_description"`
	Image            string        `json:"image"`
	InfoProduct      string        `json:"info_sku_attr"`
}
