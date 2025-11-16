package server_transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
	// OrderInfo         string             `json:"order_isnfo" binding:"required"`
}

type UserInfo struct {
	Name        string `json:"name"`        // Tên của người dùng
	PhoneNumber string `json:"phoneNumber"` // Số điện thoại của người dùng
	Address     string `json:"address"`     // address của người dùng
}

type TransactionServer struct {
	baseURL    string
	httpClient *http.Client
}
type GetTransactionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Code     string `json:"code"`
			Type     string `json:"type"`
			IsActive bool   `json:"is_active"`
		} `json:"data"`
	} `json:"result"`
}
type InitTransactionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Data struct {
			Amount       float64 `json:"amount"`
			Message      string  `json:"message"`
			OrderID      string  `json:"orderId"`
			PartnerCode  string  `json:"partnerCode"`
			PayURL       string  `json:"payUrl"`
			RequestID    string  `json:"requestId"`
			ResponseTime int64   `json:"responseTime"`
			ResultCode   int     `json:"resultCode"`
			ShortLink    string  `json:"shortLink"`
		} `json:"data"`
	} `json:"result"`
}

// =================================================================
// Constructor
// =================================================================

// NewTransactionServer tạo mới product client với dependency injection
func NewTransactionServer(baseURL string, timeout time.Duration) TransactionServer {
	return TransactionServer{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c TransactionServer) GetTransaction(payment_method_id string) (*GetTransactionsResponse, error) {

	// Tạo request
	url := fmt.Sprintf("%s/v1/transaction/payment_method/%s", c.baseURL, payment_method_id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var result GetTransactionsResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code == 400 {
		return nil, fmt.Errorf("product service returned error: %s", result.Message)
	}
	return &result, nil
}
func (c TransactionServer) CreateTransaction(token string, params InitPaymentParams) (*InitTransactionResponse, error) {
	// Tạo request
	url := fmt.Sprintf("%s/v1/transaction/init", c.baseURL)
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var initTransaction InitTransactionResponse
	err = json.Unmarshal(responseBody, &initTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if initTransaction.Code != 200 {
		return nil, fmt.Errorf("transaction service returned error: %s", initTransaction.Message)
	}
	return &initTransaction, nil
}

// func (c TransactionServer) UpdateProductSKU(token, status string, params []UpdateProductSKUParams) (*GetProductDetailResponse, error) {

// 	// Tạo request
// 	if status != "commit" && status != "hold" && status != "rollback" {
// 		return nil, fmt.Errorf("status must be 'commit', 'hold' or 'rollback'")
// 	}

// 	url := fmt.Sprintf("%s/v1/product/update_sku_reserver", c.baseURL)
// 	data := struct {
// 		Status string                   `json:"status"`
// 		Params []UpdateProductSKUParams `json:"data"`
// 	}{
// 		Status: status,
// 		Params: params,
// 	}
// 	body, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal request body: %w", err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

// 	// Gửi request
// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// Đọc response body
// 	responseBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response: %w", err)
// 	}

// 	// Kiểm tra status code
// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(responseBody))
// 	}

// 	// Parse response
// 	var result GetProductDetailResponse
// 	err = json.Unmarshal(responseBody, &result)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}
// 	if result.Code == 400 {
// 		return nil, fmt.Errorf("product service returned error: %s", result.Message)
// 	}
// 	return &result, nil
// }
