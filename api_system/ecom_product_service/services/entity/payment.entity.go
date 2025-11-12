package services

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
	Email       string `json:"email"`       // Email của người dùng
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
