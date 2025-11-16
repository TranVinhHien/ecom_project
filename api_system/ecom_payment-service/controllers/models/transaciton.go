package controllers_model

type AmountProdduct struct {
	Product_sku_id string `form:"product_sku_id" json:"product_sku_id" binding:"required"`
	Amount         int    `form:"amount" json:"amount" binding:"required"`
}

type DetailItem struct {
	ProductID string  `json:"product_id"` // Mã sản phẩm
	Name      string  `json:"name"`       // Tên sản phẩm
	ImageURL  string  `json:"image_url"`  // Link hình ảnh của sản phẩm
	Quantity  int     `json:"quantity"`   // Số lượng sản phẩm
	Price     float64 `json:"price"`      // Đơn giá
}

// SettlementDetail: Chỉ chứa các thông tin tài chính ảnh hưởng đến Shop
// và các thông tin cần thiết để Sàn tính toán chi phí riêng
type SettlementDetail struct {
	ShopOrderID               string  `json:"shop_order_id" binding:"required"`
	OrderSubtotal             float64 `json:"order_subtotal" binding:"required"` // Tiền hàng gốc (để tính hoa hồng)
	ShopFundedProductDiscount float64 `json:"shop_funded_product_discount" `     // Giảm giá SP Shop chịu
	SiteFundedProductDiscount float64 `json:"site_funded_product_discount" `     // Giảm giá SP Sàn trợ giá (Sàn bù cho Shop)
	ShopVoucherDiscount       float64 `json:"shop_voucher_discount" `            // Voucher Shop chịu
	ShippingFee               float64 `json:"shipping_fee" binding:"required"`   // Phí ship khách trả cho đơn này
	ShopShippingDiscount      float64 `json:"shop_shipping_discount"`            // Ship Shop tài trợ
	SiteOrderDiscount         float64 `json:"site_order_discount"`               // --- THÊM TRƯỜNG --- Số tiền giảm từ voucher SÀN (tiền hàng) đã được PHÂN BỔ cho đơn hàng shop này
	SiteShippingDiscount      float64 `json:"site_shipping_discount"`            // --- THÊM TRƯỜNG --- Tiền Sàn hỗ trợ ship (voucher ship) đã được PHÂN BỔ cho đơn hàng shop này
	// --- LOẠI BỎ site_shipping_discount và site_order_discount ---
	CommissionFee    float64 `json:"commission_fee" binding:"required"`     // Phí hoa hồng Sàn tính (Order Service tính sẵn)
	NetSettledAmount float64 `json:"net_settled_amount" binding:"required"` // Tiền Shop thực nhận (Order Service tính sẵn)
}
type InitPaymentParams struct {
	OrderID         string       `json:"order_id" binding:"required" `
	Amount          float64      `json:"amount" binding:"required"`            // Tổng tiền khách phải trả cuối cùng
	PaymentMethodID string       `json:"payment_method_id" binding:"required"` // UUID của phương thức thanh toán
	Items           []DetailItem `json:"items" binding:"required"`             // Giữ lại nếu cổng TT yêu cầu
	// --- THÊM CÁC TRƯỜNG TỔNG CHI PHÍ SÀN ---
	SiteOrderVoucherDiscountAmount float64 `json:"site_order_voucher_discount_amount" ` // Tổng voucher Sàn giảm giá đơn hàng
	SitePromotionDiscountAmount    float64 `json:"site_promotion_discount_amount" `     // Tổng KM Sàn giảm giá đơn hàng
	SiteShippingDiscountAmount     float64 `json:"site_shipping_discount_amount" `      // Tổng Sàn giảm giá Ship
	TotalSiteFundedProductDiscount float64 `json:"total_site_funded_product_discount" ` // Tổng Sàn trợ giá SP

	SettlementDetails []SettlementDetail `json:"settlement_details" binding:"required,min=1" ` // *Thông tin tài chính chi tiết cho từng shop*
	UserInfo          UserInfo           `json:"user_info" binding:"required"`                 // Thông tin người dùng

	// --- THÊM CÁC URL CHO THANH TOÁN ONLINE ---
	// ReturnURL string `json:"return_url" binding:"required_if=PaymentMethod.Type ONLINE"` // URL trả về trình duyệt
	// NotifyURL string `json:"notify_url" binding:"required_if=PaymentMethod.Type ONLINE"` // URL callback server-to-server
	// OrderInfo string `json:"order_info" binding:"required"`
}

type UserInfo struct {
	Name        string `json:"name" binding:"required"`        // Tên của người dùng
	PhoneNumber string `json:"phoneNumber" binding:"required"` // Số điện thoại của người dùng
	Address     string `json:"address" binding:"required"`     // address của người dùng
}

type OrderIDParams struct {
	OrderID string `form:"order_id" json:"order_id" binding:"required"`
}
