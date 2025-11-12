package controllers_model

type AmountProdduct struct {
	Product_sku_id string `form:"product_sku_id" json:"product_sku_id" binding:"required"`
	Amount         int    `form:"amount" json:"amount" binding:"required"`
}

type OrderParams struct {
	NumOfProducts []AmountProdduct `form:"num_of_products" json:"num_of_products" binding:"required"`
	Discount_Id   string           `form:"discount_id" json:"discount_id"`
	Address_id    string           `form:"address_id" json:"address_id" binding:"required"`
	Payment_id    string           `form:"payment_id" json:"payment_id" binding:"required"`
}

type OrderIDParams struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required"`
}
