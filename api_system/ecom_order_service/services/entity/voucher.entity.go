package services

import "time"

//================================================================
// 1. CÁC STRUCT ĐẦU VÀO (DTOs) CHO LỚP HANDLER
//================================================================

// CreateVoucherRequest là dữ liệu đầu vào khi Handler gọi Usecase
type CreateVoucherRequest struct {
	Name              string    `json:"name"`
	VoucherCode       string    `json:"voucher_code"`
	DiscountType      string    `json:"discount_type"` // "PERCENTAGE" or "FIXED_AMOUNT"
	DiscountValue     float64   `json:"discount_value"`
	MaxDiscountAmount *float64  `json:"max_discount_amount"` // Dùng con trỏ cho giá trị nullable
	AppliesToType     string    `json:"applies_to_type"`     // "ORDER_TOTAL" or "SHIPPING_FEE"
	MinPurchaseAmount float64   `json:"min_purchase_amount"`
	AudienceType      string    `json:"audience_type"` // "PUBLIC" or "ASSIGNED"
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	TotalQuantity     int32     `json:"total_quantity"`
	MaxUsagePerUser   int32     `json:"max_usage_per_user"`
	UserUse           []string  `json:"user_use"`
}

// UpdateVoucherRequest là dữ liệu đầu vào cho việc cập nhật từng phần
// Tất cả các trường đều là con trỏ để chúng ta biết trường nào được gửi lên (không phải nil)
type UpdateVoucherRequest struct {
	Name              *string    `json:"name"`
	VoucherCode       *string    `json:"voucher_code"`
	DiscountType      *string    `json:"discount_type"`
	DiscountValue     *float64   `json:"discount_value"`
	MaxDiscountAmount *float64   `json:"max_discount_amount"`
	AppliesToType     *string    `json:"applies_to_type"`
	MinPurchaseAmount *float64   `json:"min_purchase_amount"`
	AudienceType      *string    `json:"audience_type"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	TotalQuantity     *int32     `json:"total_quantity"`
	MaxUsagePerUser   *int32     `json:"max_usage_per_user"`
	IsActive          *bool      `json:"is_active"`
}

// UseVoucherInput định nghĩa thông tin cần thiết để sử dụng 1 voucher
type UseVoucherInput struct {
	VoucherID      string
	UserID         string
	DiscountAmount float64 // Bắt buộc: để ghi vào history
}

// RollbackVoucherInput định nghĩa thông tin để hoàn trả 1 voucher
type RollbackVoucherInput struct {
	VoucherID string
	UserID    string
}

// VoucherFilterRequest định nghĩa các điều kiện lọc voucher
type VoucherFilterRequest struct {
	OwnerType     *string `form:"owner_type"`      // "PLATFORM" hoặc "SHOP"
	ShopID        *string `form:"shop_id"`         // Lọc theo shop cụ thể (khi owner_type=SHOP)
	AppliesToType *string `form:"applies_to_type"` // "ORDER_TOTAL" hoặc "SHIPPING_FEE"
	SortBy        string  `form:"sort_by"`         // "discount_asc", "discount_desc", "created_at"
}
