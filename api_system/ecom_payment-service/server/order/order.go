package server_transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OrderServices struct {
	baseURL    string
	httpClient *http.Client
}
type UpdateOrderStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	// Result  struct {
	// 	Data struct {
	// 		ID       string `json:"id"`
	// 		Name     string `json:"name"`
	// 		Code     string `json:"code"`
	// 		Type     string `json:"type"`
	// 		IsActive bool   `json:"is_active"`
	// 	} `json:"data"`
	// } `json:"result"`
	Result interface{} `json:"result"`
}

// =================================================================
// Constructor
// =================================================================

// NewOrderServices tạo mới product client với dependency injection
func NewOrderServices(baseURL string, timeout time.Duration) OrderServices {
	return OrderServices{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}
func (c OrderServices) UpdateOrderCallbackPayment(order_id string) error {

	url := fmt.Sprintf("%s/v1/orders/callback_payment_online/%s", c.baseURL, order_id)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Gửi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var result UpdateOrderStatusResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Code == 400 {
		return fmt.Errorf("product service returned error: %s", result.Message)
	}
	return nil
}

// func (c OrderServices) UpdateProductSKU(token, status string, params []UpdateProductSKUParams) (*GetProductDetailResponse, error) {

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
