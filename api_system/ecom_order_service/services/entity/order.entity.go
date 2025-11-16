package services

import (
	"time"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
)

// ShopOrderSearchFilter đại diện cho các bộ lọc tìm kiếm shop orders
type ShopOrderSearchFilter struct {
	Status         *string    `json:"status,omitempty" form:"status"`
	ShopID         *string    `json:"shop_id,omitempty" form:"shop_id"`
	MinAmount      *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount      *float64   `json:"max_amount,omitempty" form:"max_amount"`
	CreatedFrom    *time.Time `json:"created_from,omitempty" form:"created_from"`
	CreatedTo      *time.Time `json:"created_to,omitempty" form:"created_to"`
	PaidFrom       *time.Time `json:"paid_from,omitempty" form:"paid_from"`
	PaidTo         *time.Time `json:"paid_to,omitempty" form:"paid_to"`
	ProcessingFrom *time.Time `json:"processing_from,omitempty" form:"processing_from"`
	ProcessingTo   *time.Time `json:"processing_to,omitempty" form:"processing_to"`
	ShippedFrom    *time.Time `json:"shipped_from,omitempty" form:"shipped_from"`
	ShippedTo      *time.Time `json:"shipped_to,omitempty" form:"shipped_to"`
	CompletedFrom  *time.Time `json:"completed_from,omitempty" form:"completed_from"`
	CompletedTo    *time.Time `json:"completed_to,omitempty" form:"completed_to"`
	CancelledFrom  *time.Time `json:"cancelled_from,omitempty" form:"cancelled_from"`
	CancelledTo    *time.Time `json:"cancelled_to,omitempty" form:"cancelled_to"`

	// Pagination (optional, có default)
	Page     int `json:"page,omitempty" form:"page" binding:"omitempty,min=1"`
	PageSize int `json:"page_size,omitempty" form:"page_size" binding:"omitempty,min=1,max=100"`

	// Sorting (optional, có default)
	SortBy string `json:"sort_by,omitempty" form:"sort_by"` // created_at, total_amount, paid_at, processing_at, shipped_at, completed_at
}

// CreateOrderRequest đại diện cho request tạo đơn hàng mới
type CreateOrderRequest struct {
	ShippingAddress  ShippingAddress    `json:"shippingAddress" binding:"required"`
	PaymentMethod_ID string             `json:"paymentMethod" binding:"required"`
	Items            []OrderItemRequest `json:"items" binding:"required,min=1"`
	// Vouchers         []VoucherRequest   `json:"vouchers"`
	Note        *string              `json:"note"`
	VoucherShop []VoucherShopRequest `json:"voucher_shop"`
	// Email             string               `json:"email" binding:"required,email"`
	VoucherSiteID     *string `json:"voucher_site_id"` // Có thể là
	VoucherShippingID *string `json:"voucher_shipping_id"`
}

// OrderItemRequest đại diện cho một item trong đơn hàng
type OrderItemRequest struct {
	SkuID    string `json:"sku_id" binding:"required"`
	ShopID   string `json:"shop_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}
type VoucherType string

const (
	VoucherTypeShop  = "SHOP_VOUCHER"
	VoucherTypeShip  = "SHIPPING_VOUCHER"
	VoucherTypeOrder = "ORDER_VOUCHER"
)

type AppliedVoucherInfo struct {
	Voucher        db.Vouchers
	DiscountAmount float64
	ShopOrderID    string // Rỗng nếu là voucher Sàn/Ship
}

// VoucherRequest đại diện cho voucher trong đơn hàng
type VoucherShopRequest struct {
	VoucherID string `json:"voucher_id" binding:"required"`
	// Type      VoucherType `json:"type" binding:"required,oneof=SHOP_VOUCHER SHIPPING_VOUCHER ORDER_VOUCHER"`
	ShopID string `json:"shop_id" binding:"required"` // Chỉ cần thiết khi type là SHOP_VOUCHER
}

// ShippingAddress đại diện cho địa chỉ giao hàng
type ShippingAddress struct {
	FullName   string  `json:"fullName" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	District   *string `json:"district"`
	City       *string `json:"city"`
	PostalCode *string `json:"postalCode"`
}

// PaymentMethod đại diện cho phương thức thanh toán

type PaymentMethodsType string

const (
	PaymentMethodsTypeONLINE  PaymentMethodsType = "ONLINE"
	PaymentMethodsTypeOFFLINE PaymentMethodsType = "OFFLINE"
)

// Quản lý các phương thức thanh toán
type PaymentMethod struct {
	// UUID, Khóa chính
	ID string `json:"id"`
	// Tên hiển thị cho người dùng (Ví dụ: "Ví MoMo")
	Name string `json:"name"`
	// Mã định danh (Ví dụ: "MOMO", "COD")
	Code string `json:"code"`
	// Loại hình thanh toán
	Type PaymentMethodsType `json:"type"`
	// Cho phép bật/tắt phương thức này
	IsActive bool `json:"is_active"`
}

// ShipOrderRequest đại diện cho request đánh dấu đơn hàng đã ship
type ShipOrderRequest struct {
	ShippingMethod string `json:"shipping_method" binding:"required"`
	TrackingCode   string `json:"tracking_code" binding:"required"`
}

// Order status enums
type OrderStatus string

const (
	OrderStatusPending            OrderStatus = "PENDING"
	OrderStatusProcessing         OrderStatus = "PROCESSING"
	OrderStatusPartiallyShipped   OrderStatus = "PARTIALLY_SHIPPED"
	OrderStatusCompleted          OrderStatus = "COMPLETED"
	OrderStatusCancelled          OrderStatus = "CANCELLED"
	OrderStatusPartiallyCancelled OrderStatus = "PARTIALLY_CANCELLED"
)

// ShopOrderStatus enums
type ShopOrderStatus string

const (
	ShopOrderStatusPending              ShopOrderStatus = "PENDING"
	ShopOrderStatusAwaitingPayment      ShopOrderStatus = "AWAITING_PAYMENT"
	ShopOrderStatusAwaitingConfirmation ShopOrderStatus = "AWAITING_CONFIRMATION"
	ShopOrderStatusProcessing           ShopOrderStatus = "PROCESSING"
	ShopOrderStatusShipped              ShopOrderStatus = "SHIPPED"
	ShopOrderStatusCompleted            ShopOrderStatus = "COMPLETED"
	ShopOrderStatusCancelled            ShopOrderStatus = "CANCELLED"
	ShopOrderStatusRefunded             ShopOrderStatus = "REFUNDED"
)

// CreateOrderResponse đại diện cho response sau khi tạo đơn hàng
type CreateOrderResponse struct {
	OrderID    string   `json:"order_id"`
	OrderCode  string   `json:"order_code"`
	GrandTotal float64  `json:"grand_total"`
	Status     string   `json:"status"`
	PaymentURL *string  `json:"payment_url"`
	ShopOrders []string `json:"shop_orders"` // Danh sách shop_order_code
}

// OrderSummary đại diện cho thông tin tóm tắt của đơn hàng
type OrderSummary struct {
	OrderID    string  `json:"order_id"`
	OrderCode  string  `json:"order_code"`
	Status     string  `json:"status"`
	GrandTotal float64 `json:"grand_total"`
	ItemCount  int     `json:"item_count"`
	CreatedAt  string  `json:"created_at"`
}

// OrderDetail đại diện cho thông tin chi tiết đơn hàng
type OrderDetail struct {
	OrderID                     string          `json:"order_id"`
	OrderCode                   string          `json:"order_code"`
	UserID                      string          `json:"user_id"`
	Status                      string          `json:"status"`
	GrandTotal                  float64         `json:"grand_total"`
	Subtotal                    float64         `json:"subtotal"`
	TotalShippingFee            float64         `json:"total_shipping_fee"`
	TotalDiscount               float64         `json:"total_discount"`
	SiteOrderVoucherCode        *string         `json:"site_order_voucher_code"`
	SiteOrderVoucherDiscount    float64         `json:"site_order_voucher_discount"`
	SiteShippingVoucherCode     *string         `json:"site_shipping_voucher_code"`
	SiteShippingVoucherDiscount float64         `json:"site_shipping_voucher_discount"`
	ShippingAddress             ShippingAddress `json:"shipping_address"`
	PaymentMethod               PaymentMethod   `json:"payment_method"`
	Note                        *string         `json:"note"`
	CreatedAt                   string          `json:"created_at"`
	UpdatedAt                   string          `json:"updated_at"`
}

// ShopOrderDetail đại diện cho thông tin chi tiết shop order
type ShopOrderDetail struct {
	ShopOrderID          string            `json:"shop_order_id"`
	ShopOrderCode        string            `json:"shop_order_code"`
	ShopID               string            `json:"shop_id"`
	Status               string            `json:"status"`
	Subtotal             float64           `json:"subtotal"`
	ShippingFee          float64           `json:"shipping_fee"`
	TotalDiscount        float64           `json:"total_discount"`
	TotalAmount          float64           `json:"total_amount"`
	ShopVoucherCode      *string           `json:"shop_voucher_code"`
	ShopVoucherDiscount  float64           `json:"shop_voucher_discount"`
	ShippingMethod       *string           `json:"shipping_method"`
	TrackingCode         *string           `json:"tracking_code"`
	SiteOrderDiscount    float64           `json:"site_order_discount"`
	SiteShippingDiscount float64           `json:"site_shipping_discount"`
	Items                []OrderItemDetail `json:"items"`
	CreatedAt            string            `json:"created_at"`
	UpdatedAt            string            `json:"updated_at"`
	PaidAt               *string           `json:"paid_at"`
	ProcessingAt         *string           `json:"processing_at"`
	ShippedAt            *string           `json:"shipped_at"`
	CompletedAt          *string           `json:"completed_at"`
	CancelledAt          *string           `json:"cancelled_at"`
}

// OrderItemDetail đại diện cho thông tin chi tiết item trong đơn hàng
type OrderItemDetail struct {
	ItemID             string                 `json:"item_id"`
	ProductID          string                 `json:"product_id"`
	SkuID              string                 `json:"sku_id"`
	Quantity           int                    `json:"quantity"`
	OriginalUnitPrice  float64                `json:"original_unit_price"`
	FinalUnitPrice     float64                `json:"final_unit_price"`
	TotalPrice         float64                `json:"total_price"`
	Reviewed           bool                   `json:"reviewed"`
	ProductName        string                 `json:"product_name"`
	ProductImage       *string                `json:"product_image"`
	SkuAttributes      string                 `json:"sku_attributes"`
	PromotionsSnapshot map[string]interface{} `json:"promotions_snapshot"`
}
