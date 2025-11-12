package services

import (
	"time"
)

type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
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

type Users struct {
}
type Accounts struct {
	AccountID    string          `json:"account_id"`
	Username     string          `json:"username"`
	Password     string          `json:"password"`
	ActiveStatus string          `json:"active_status"`
	CreateDate   time.Time       `json:"create_date"`
	UpdateDate   Narg[time.Time] `json:"update_date"`
}

type Categorys struct {
	CategoryID string            `json:"category_id"`
	Name       string            `json:"name"`
	Key        string            `json:"key"`
	Path       string            `json:"path"`
	Parent     Narg[string]      `json:"parent"`
	Childs     Narg[[]Categorys] `json:"child"`
}

type CustomerAddress struct {
	IDAddress   string          `json:"id_address"`
	CustomerID  string          `json:"customer_id"`
	Address     string          `json:"address"`
	PhoneNumber string          `json:"phone_number"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
}

type Customers struct {
	CustomerID              string          `json:"customer_id"`
	Name                    string          `json:"name"`
	Email                   string          `json:"email"`
	Image                   Narg[string]    `json:"image"`
	Dob                     time.Time       `json:"dob"`
	Gender                  string          `json:"gender"`
	DeviceRegistrationToken Narg[string]    `json:"device_registration_token"`
	AccountID               string          `json:"account_id"`
	CreateDate              time.Time       `json:"create_date"`
	UpdateDate              Narg[time.Time] `json:"update_date"`
}

type DescriptionAttr struct {
	DescriptionAttrID string `json:"description_attr_id"`
	Name              string `json:"name"`
	Value             string `json:"value"`
	ProductsSpuID     string `json:"products_spu_id"`
}

type Discounts struct {
	DiscountID     string          `json:"discount_id"`
	DiscountCode   string          `json:"discount_code"`
	DiscountValue  float64         `json:"discount_value"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	MinOrderValue  float64         `json:"min_order_value"`
	Amount         int32           `json:"amount"`
	StatusDiscount string          `json:"status_discount"`
	CreateDate     time.Time       `json:"create_date"`
	UpdateDate     Narg[time.Time] `json:"update_date"`
}

type Employees struct {
	EmployeeID  string          `json:"employee_id"`
	Gender      string          `json:"gender"`
	Dob         time.Time       `json:"dob"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phone_number"`
	Address     string          `json:"address"`
	Salary      float64         `json:"salary"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
	AccountID   string          `json:"account_id"`
}

type OrderDetail struct {
	OrderDetailID string                  `json:"order_detail_id"`
	Quantity      int32                   `json:"quantity"`
	UnitPrice     float64                 `json:"unit_price"`
	ProductSkuID  string                  `json:"product_sku_id"`
	OrderID       string                  `json:"order_id"`
	ProductSKU    Narg[ProductSkusDetail] `json:"product_info"`
}

type Orders struct {
	OrderID           string                `json:"order_id"`
	OrderDate         time.Time             `json:"order_date"`
	TotalAmount       float64               `json:"total_amount"`
	CustomerAddressID string                `json:"customer_address_id"`
	DiscountID        Narg[string]          `json:"discount_id"`
	PaymentMethodID   string                `json:"payment_method_id"`
	PaymentStatus     string                `json:"payment_status"`
	OrderStatus       string                `json:"order_status"`
	CreateDate        time.Time             `json:"create_date"`
	UpdateDate        Narg[time.Time]       `json:"update_date"`
	CustomerID        string                `json:"customer_id"`
	OrderDetail       []OrderDetail         `json:"order_detail"`
	Address           Narg[CustomerAddress] `json:"customer_address"`
	PaymentMethod     Narg[PaymentMethods]  `json:"payment_method"`
}

type PaymentMethods struct {
	PaymentMethodID string       `json:"payment_method_id"`
	Name            string       `json:"name"`
	Description     Narg[string] `json:"description"`
}

type ProductSkuAttrs struct {
	ProductSkuAttrID string       `json:"product_sku_attr_id"`
	Name             string       `json:"name"`
	Value            string       `json:"value"`
	Image            Narg[string] `json:"image"`
	ProductsSpuID    string       `json:"products_spu_id"`
}

type ProductSkus struct {
	ProductSkuID  string          `json:"product_sku_id"`
	Value         string          `json:"value"`
	SkuStock      int32           `json:"sku_stock"`
	Price         float64         `json:"price"`
	Sort          int32           `json:"sort"`
	CreateDate    time.Time       `json:"create_date"`
	UpdateDate    Narg[time.Time] `json:"update_date"`
	ProductsSpuID string          `json:"products_spu_id"`
}
type ProductSkusDetail struct {
	ProductSkuID     string          `json:"product_sku_id"`
	Value            string          `json:"value"`
	SkuStock         int32           `json:"sku_stock"`
	Price            float64         `json:"price"`
	Sort             int32           `json:"sort"`
	CreateDate       time.Time       `json:"create_date"`
	UpdateDate       Narg[time.Time] `json:"update_date"`
	ProductsSpuID    string          `json:"products_spu_id"`
	Name             string          `json:"name"`
	ShortDescription string          `json:"short_description"`
	Image            string          `json:"image"`
	InfoProduct      string          `json:"info_sku_attr"`
}
type ProductsSpu struct {
	ProductsSpuID    string          `json:"products_spu_id"`
	Name             string          `json:"name"`
	BrandID          string          `json:"brand_id"`
	Description      string          `json:"description"`
	ShortDescription string          `json:"short_description"`
	StockStatus      string          `json:"stock_status"`
	DeleteStatus     string          `json:"delete_status"`
	Sort             int32           `json:"sort"`
	CreateDate       time.Time       `json:"create_date"`
	UpdateDate       Narg[time.Time] `json:"update_date"`
	Image            string          `json:"image"`
	Media            string          `json:"media"`
	Key              string          `json:"key"`
	CategoryID       string          `json:"category_id"`
}

type PurchaseOrderDetail struct {
	PurchaseOrderDetailID string  `json:"purchase_order_detail_id"`
	Quantity              int32   `json:"quantity"`
	UnitPrice             float64 `json:"unit_price"`
	PurchaseOrderID       string  `json:"purchase_order_id"`
	ProductSkuID          string  `json:"product_sku_id"`
}

type PurchaseOrders struct {
	PurchaseOrderID string          `json:"purchase_order_id"`
	TotalAmount     float64         `json:"total_amount"`
	Status          string          `json:"status"`
	CreateDate      time.Time       `json:"create_date"`
	UpdateDate      Narg[time.Time] `json:"update_date"`
	SupplierID      string          `json:"supplier_id"`
	EmployeeID      string          `json:"employee_id"`
}

type Ratings struct {
	RatingID      string          `json:"rating_id"`
	Comment       Narg[string]    `json:"comment"`
	Star          int32           `json:"star"`
	CreateDate    time.Time       `json:"create_date"`
	UpdateDate    Narg[time.Time] `json:"update_date"`
	CustomerID    string          `json:"custonmer_id"`
	ProductsSpuID string          `json:"products_spu_id"`
	UserInfo      Narg[Customers] `json:"user_info"`
}

type RoleAccount struct {
	RoleAccountID string `json:"role_account_id"`
	AccountID     string `json:"account_id"`
	RoleID        string `json:"role_id"`
}

type Roles struct {
	RoleID      string       `json:"role_id"`
	Name        string       `json:"name"`
	Description Narg[string] `json:"description"`
}

type Suppliers struct {
	SupplierID  string          `json:"supplier_id"`
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Address     Narg[string]    `json:"address"`
	CreateDate  time.Time       `json:"create_date"`
	UpdateDate  Narg[time.Time] `json:"update_date"`
}

// // struct cho những trường hợp sử lý ở db thay vì service
type AmountProdduct struct {
	Product_sku_id string ` json:"product_sku_id"`
	Amount         int    `json:"amount"`
}

type CreateOrderParams struct {
	NumOfProducts []AmountProdduct ` json:"num_of_products" `
	Discount_Id   string           ` json:"discount_id"`
	Address_id    string           `json:"address_id"`
	Payment_id    string           `json:"payment_id"`
}

type ProductSimple struct {
	ProductsSpuID    string `json:"products_spu_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	StockStatus      string `json:"stock_status"`
	DeleteStatus     string `json:"delete_status"`
	Image            string `json:"image"`
	// Media            string  `json:"media"`
	Key         string  `json:"key"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
	Avg_star    float64 `json:"average_star"`
	TotalRating int32   `json:"total_rating"`
}
type ProductDetail struct {
	Spu     ProductsSpu       `json:"spu"`
	Sku     []ProductSkus     `json:"sku"`
	DesAttr []DescriptionAttr `json:"description_attrs"`
	SkuAttr []ProductSkuAttrs `json:"sku_attrs"`
}
