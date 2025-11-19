package server_order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OrderServer struct {
	baseURL    string
	httpClient *http.Client
}

// ProductRatingStatsItem represents rating statistics for a single product
type ProductRatingStatsItem struct {
	ProductID     string  `json:"product_id"`
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

// GetProductTotalSold
type GetProductTotalSold struct {
	ProductID string `json:"product_id"`
	TotalSold int64  `json:"total_sold"`
}

// GetBulkProductRatingStatsRequest represents the request to get rating stats
type GetBulkProductRatingStatsRequest struct {
	ProductIDs []string `json:"product_ids"`
}

// GetBulkProductRatingStatsResponse represents the response from order service
type GetBulkProductRatingStatsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Result  struct {
		Data []ProductRatingStatsItem `json:"data"`
	} `json:"result"`
}

// GetBulkProductRatingStatsResponse represents the response from order service
type GetProductTotalSoldResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Result  struct {
		Data []GetProductTotalSold `json:"data"`
	} `json:"result"`
}

// NewOrderServer tạo mới order client với dependency injection
func NewOrderServer(baseURL string, timeout time.Duration) OrderServer {
	return OrderServer{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetBulkProductRatingStats gọi API để lấy thống kê đánh giá cho nhiều sản phẩm
func (c OrderServer) GetBulkProductRatingStats(productIDs []string) (map[string]ProductRatingStatsItem, error) {
	if len(productIDs) == 0 {
		return make(map[string]ProductRatingStatsItem), nil
	}

	// Tạo request body
	reqBody := GetBulkProductRatingStatsRequest{
		ProductIDs: productIDs,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Tạo HTTP request
	url := fmt.Sprintf("%s/v1/comments/bulk-stats", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var apiResp GetBulkProductRatingStatsResponse
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Kiểm tra response status
	if apiResp.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", apiResp.Message)
	}
	// Chuyển đổi slice thành map
	result := make(map[string]ProductRatingStatsItem)
	for _, item := range apiResp.Result.Data {
		result[item.ProductID] = item
	}

	return result, nil
}

// GetProductTotalSold gọi API Order Service để lấy tổng số sản phẩm đã bán
func (c OrderServer) GetProductTotalSold(ctx context.Context, productIDs []string) (map[string]GetProductTotalSold, error) {
	if len(productIDs) == 0 {
		return make(map[string]GetProductTotalSold), nil
	}

	// Tạo request body
	reqBody := GetBulkProductRatingStatsRequest{
		ProductIDs: productIDs,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Tạo HTTP request
	url := fmt.Sprintf("%s/v1/orders/get_product_total_sold", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 4. Đọc response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 5. Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// 6. Parse response (Dùng struct ĐÚNG)
	var apiResp GetProductTotalSoldResponse // Dùng struct cho response tổng
	err = json.Unmarshal(responseBody, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 7. Kiểm tra response status
	if apiResp.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", apiResp.Message)
	}

	// 8. Chuyển đổi slice thành map (để tra cứu O(1) hiệu quả)
	result := make(map[string]GetProductTotalSold)
	for _, item := range apiResp.Result.Data {
		result[item.ProductID] = item
	}

	return result, nil
}
