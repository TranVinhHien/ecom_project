package server_product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ProductServer struct {
	baseURL    string
	httpClient *http.Client
}

// =================================================================
// Response wrappers
// =================================================================

type GetSKUResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Data SKUData `json:"data"`
	} `json:"result"`
}

type GetProductListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Data []ProductListItem `json:"data"`
	} `json:"result"`
}

type GetProductDetailResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Data ProductDetail `json:"data"`
	} `json:"result"`
}

// =================================================================
// SKU Data structures
// =================================================================

type SKUData struct {
	CreateDate       time.Time `json:"create_date"`
	ID               string    `json:"id"`
	Price            float64   `json:"price"`
	ProductID        string    `json:"product_id"`
	Quantity         int       `json:"quantity"`
	QuantityReserver int       `json:"quantity_reserver"`
	SkuCode          string    `json:"sku_code"`
	SkuName          string    `json:"sku_name"`
	UpdateDate       time.Time `json:"update_date"`
	Weight           float64   `json:"weight"`
}

// =================================================================
// Product List Item structures
// =================================================================

type ProductListItem struct {
	BrandID                   string    `json:"brand_id"`
	CategoryID                string    `json:"category_id"`
	CreateBy                  string    `json:"create_by"`
	CreateDate                time.Time `json:"create_date"`
	DeleteStatus              string    `json:"delete_status"`
	Description               string    `json:"description"`
	ID                        string    `json:"id"`
	Image                     string    `json:"image"`
	Key                       string    `json:"key"`
	MaxPrice                  float64   `json:"max_price"`
	MaxPriceSkuID             string    `json:"max_price_sku_id"`
	Media                     string    `json:"media"` // JSON array as string
	MinPrice                  float64   `json:"min_price"`
	MinPriceSkuID             string    `json:"min_price_sku_id"`
	Name                      string    `json:"name"`
	ProductIsPermissionCheck  bool      `json:"product_is_permission_check"`
	ProductIsPermissionReturn bool      `json:"product_is_permission_return"`
	ShopID                    string    `json:"shop_id"`
	ShortDescription          string    `json:"short_description"`
	UpdateBy                  *string   `json:"update_by"`
	UpdateDate                time.Time `json:"update_date"`
}

// =================================================================
// Product Detail structures
// =================================================================

type ProductDetail struct {
	Brand    Brand           `json:"brand"`
	Category Category        `json:"category"`
	Option   []ProductOption `json:"option"`
	Product  ProductInfo     `json:"product"`
	SKU      []ProductSKU    `json:"sku"`
}

type Brand struct {
	BrandID    string    `json:"brand_id"`
	Code       string    `json:"code"`
	CreateDate time.Time `json:"create_date"`
	Image      *string   `json:"image"`
	Name       string    `json:"name"`
	UpdateDate time.Time `json:"update_date"`
}

type Category struct {
	CategoryID string  `json:"category_id"`
	Image      *string `json:"image"`
	Key        string  `json:"key"`
	Name       string  `json:"name"`
	Parent     *string `json:"parent"`
	Path       string  `json:"path"`
}

type ProductOption struct {
	OptionName string               `json:"option_name"`
	Values     []ProductOptionValue `json:"values"`
}

type ProductOptionValue struct {
	Image         *string `json:"image,omitempty"`
	OptionValueID string  `json:"option_value_id"`
	Value         string  `json:"value"`
}

type ProductInfo struct {
	BrandID                   string    `json:"brand_id"`
	CategoryID                string    `json:"category_id"`
	CreateBy                  string    `json:"create_by"`
	CreateDate                time.Time `json:"create_date"`
	DeleteStatus              string    `json:"delete_status"`
	Description               string    `json:"description"`
	ID                        string    `json:"id"`
	Image                     string    `json:"image"`
	Key                       string    `json:"key"`
	MaxPrice                  float64   `json:"max_price"`
	MaxPriceSkuID             string    `json:"max_price_sku_id"`
	Media                     string    `json:"media"` // Comma-separated URLs
	MinPrice                  float64   `json:"min_price"`
	MinPriceSkuID             string    `json:"min_price_sku_id"`
	Name                      string    `json:"name"`
	ProductIsPermissionCheck  bool      `json:"product_is_permission_check"`
	ProductIsPermissionReturn bool      `json:"product_is_permission_return"`
	ShopID                    string    `json:"shop_id"`
	ShortDescription          string    `json:"short_description"`
	UpdateBy                  *string   `json:"update_by"`
	UpdateDate                time.Time `json:"update_date"`
}

type ProductSKU struct {
	ID             string   `json:"id"`
	OptionValueIDs []string `json:"option_value_ids"`
	Price          float64  `json:"price"`
	Quantity       int      `json:"quantity"`
	SkuCode        string   `json:"sku_code"`
	Weight         float64  `json:"weight"`
}

type UpdateProductSKUParams struct {
	Sku_ID           string `json:"sku_id"`
	QuantityReserved int    `json:"quantity_reserver"`
}

// =================================================================
// Constructor
// =================================================================

// NewProductServer tạo mới product client với dependency injection
func NewProductServer(baseURL string, timeout time.Duration) ProductServer {
	return ProductServer{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c ProductServer) GetSKUs(sku_id string) (*GetSKUResponse, error) {

	// Tạo request
	url := fmt.Sprintf("%s/v1/product/getsku/%s", c.baseURL, sku_id)
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
	var result GetSKUResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code == 400 {
		return nil, fmt.Errorf("product service returned error: %s", result.Message)
	}
	return &result, nil
}
func (c ProductServer) GetProductDetail(product_id string) (*GetProductDetailResponse, error) {

	// Tạo request
	url := fmt.Sprintf("%s/v1/product/getdetail_with_id/%s", c.baseURL, product_id)
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
	var result GetProductDetailResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code == 400 {
		return nil, fmt.Errorf("product service returned error: %s", result.Message)
	}
	return &result, nil
}
func (c ProductServer) UpdateProductSKU(token, status string, params []UpdateProductSKUParams) (*GetProductDetailResponse, error) {

	// Tạo request
	if status != "commit" && status != "hold" && status != "rollback" {
		return nil, fmt.Errorf("trạng thái phải là 'commit', 'hold' hoặc 'rollback'")
	}

	url := fmt.Sprintf("%s/v1/product/update_sku_reserver", c.baseURL)
	data := struct {
		Status string                   `json:"status"`
		Params []UpdateProductSKUParams `json:"data"`
	}{
		Status: status,
		Params: params,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi marshal dữ liệu: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gửi request: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi đọc response: %w", err)
	}
	// Parse response
	var result GetProductDetailResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf(" lỗi khi parse response: %w", err)
	}
	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lỗi khi upload với status %d: %s", result.Code, result.Message)
	}

	if result.Code == 400 {
		return nil, fmt.Errorf("lỗi: %s", result.Message)
	}
	return &result, nil
}
