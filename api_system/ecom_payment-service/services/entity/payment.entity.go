package services

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
	OrderSubtotal             float64 `json:"order_subtotal"`               // Tiền hàng gốc (để tính hoa hồng)
	ShopFundedProductDiscount float64 `json:"shop_funded_product_discount"` // Giảm giá SP Shop chịu
	SiteFundedProductDiscount float64 `json:"site_funded_product_discount"` // Giảm giá SP Sàn trợ giá (Sàn bù cho Shop)
	ShopVoucherDiscount       float64 `json:"shop_voucher_discount"`        // Voucher Shop chịu
	ShippingFee               float64 `json:"shipping_fee"`                 // Phí ship khách trả cho đơn này
	ShopShippingDiscount      float64 `json:"shop_shipping_discount"`       // Ship Shop tài trợ
	// --- LOẠI BỎ site_shipping_discount và site_order_discount ---
	CommissionFee    float64 `json:"commission_fee"`     // Phí hoa hồng Sàn tính (Order Service tính sẵn)
	NetSettledAmount float64 `json:"net_settled_amount"` // Tiền Shop thực nhận (Order Service tính sẵn)
}
type InitPaymentParams struct {
	OrderID         string       `json:"order_id" binding:"required"`
	Amount          float64      `json:"amount" binding:"required"`            // Tổng tiền khách phải trả cuối cùng
	PaymentMethodID string       `json:"payment_method_id" binding:"required"` // UUID của phương thức thanh toán
	Items           []DetailItem `json:"items"`                                // Giữ lại nếu cổng TT yêu cầu
	// --- THÊM CÁC TRƯỜNG TỔNG CHI PHÍ SÀN ---
	SiteOrderVoucherDiscountAmount float64 `json:"site_order_voucher_discount_amount"` // Tổng voucher Sàn giảm giá đơn hàng
	SitePromotionDiscountAmount    float64 `json:"site_promotion_discount_amount"`     // Tổng KM Sàn giảm giá đơn hàng
	SiteShippingDiscountAmount     float64 `json:"site_shipping_discount_amount"`      // Tổng Sàn giảm giá Ship
	TotalSiteFundedProductDiscount float64 `json:"total_site_funded_product_discount"` // Tổng Sàn trợ giá SP

	SettlementDetails []SettlementDetail `json:"settlement_details" binding:"required,min=1"` // *Thông tin tài chính chi tiết cho từng shop*
	UserInfo          User_MOMO          `json:"user_info" binding:"required"`                // Thông tin người dùng

	// --- THÊM CÁC URL CHO THANH TOÁN ONLINE ---
	// ReturnURL string `json:"return_url" binding:"required_if=PaymentMethod.Type ONLINE"` // URL trả về trình duyệt
	// NotifyURL string `json:"notify_url" binding:"required_if=PaymentMethod.Type ONLINE"` // URL callback server-to-server
	// OrderInfo string `json:"order_info" binding:"required"`
}

type Product_MOMO struct {
	ID          string `json:"id"`          // SKU number
	Name        string `json:"name"`        // Tên sản phẩm
	Description string `json:"description"` // Miêu tả sản phẩm
	ImageURL    string `json:"imageUrl"`    // Link hình ảnh của sản phẩm
	// Manufacturer string `json:"manufacturer"` // Tên nhà sản xuất
	Price    int64  `json:"price"`    // Đơn giá (Long)
	Currency string `json:"currency"` // Đơn vị tiền tệ (VND)
	Quantity int    `json:"quantity"` // Số lượng sản phẩm (lớn hơn 0)
	// Unit         string `json:"unit"`         // Đơn vị đo lường của sản phẩm
	TotalPrice int64 `json:"totalPrice"` // Tổng giá (Đơn giá x Số lượng)
	// TaxAmount    int64  `json:"taxAmount"`    // Tổng thuế
}
type User_MOMO struct {
	Name        string `json:"name"`        // Tên của người dùng
	PhoneNumber string `json:"phoneNumber"` // Số điện thoại của người dùng
	Address     string `json:"address"`     // address của người dùng
}

type Payload_MOMO struct {
	PartnerCode  string         `json:"partnerCode"`
	AccessKey    string         `json:"accessKey"`
	RequestID    string         `json:"requestId"`
	Amount       string         `json:"amount"`
	OrderID      string         `json:"orderId"`
	OrderInfo    string         `json:"orderInfo"`
	PartnerName  string         `json:"partnerName"`
	Items        []Product_MOMO `json:"items"`
	UserInfo     User_MOMO      `json:"userInfo"`
	StoreId      string         `json:"storeId"`
	OrderGroupId string         `json:"orderGroupId"`
	Lang         string         `json:"lang"`
	AutoCapture  bool           `json:"autoCapture"`
	RedirectUrl  string         `json:"redirectUrl"`
	IpnUrl       string         `json:"ipnUrl"`
	ExtraData    string         `json:"extraData"`
	RequestType  string         `json:"requestType"`
	Signature    string         `json:"signature"`
}

type TransactionMoMO struct {
	Amount        float64 `json:"amount"`       // Số tiền thanh toán
	ExtraData     string  `json:"extraData"`    // Dữ liệu bổ sung (nếu có)
	Message       string  `json:"message"`      // Thông báo trạng thái
	OrderID       string  `json:"orderId"`      // ID của đơn hàng
	OrderInfo     string  `json:"orderInfo"`    // Thông tin đơn hàng
	OrderType     string  `json:"orderType"`    // Loại đơn hàng (ví dụ: momo_wallet)
	PartnerCode   string  `json:"partnerCode"`  // Mã đối tác
	PayType       string  `json:"payType"`      // Loại thanh toán (qr, etc.)
	RequestID     string  `json:"requestId"`    // ID của yêu cầu thanh toán
	ResponseTime  float64 `json:"responseTime"` // Thời gian phản hồi (dạng timestamp)
	ResultCode    int     `json:"resultCode"`   // Mã kết quả
	Signature     string  `json:"signature"`    // Chuỗi chữ ký xác thực
	TransactionID float64 `json:"transId"`      // ID giao dịch

}
type OrderValue struct {
	OrderID     string  `json:"orderId"`
	TotalAmount float64 `json:"amount"`
}
type CombinedDataPayLoadMoMo struct {
	Info  User_MOMO      `json:"info"`
	Items []Product_MOMO `json:"items"`
	// OrderTX *Orders        `json:"orderTX"`
	Order         OrderValue `json:"order"`
	TransactionID string     `json:"transaction_id"`
	Email         string     `json:"email"`
}
