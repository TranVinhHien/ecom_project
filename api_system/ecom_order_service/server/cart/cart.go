package server_cart

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"
// )

// type CartServer struct {
// 	baseURL    string
// 	httpClient *http.Client
// }

// // =================================================================
// // Constructor
// // =================================================================

// // NewCartServer tạo mới product client với dependency injection
// func NewCartServer(baseURL string, timeout time.Duration) CartServer {
// 	return CartServer{
// 		baseURL: baseURL,
// 		httpClient: &http.Client{
// 			Timeout: timeout,
// 		},
// 	}
// }

// type RemoveIteminCartResponse struct {
// 	Name        string `json:"name"`        // Tên của người dùng
// 	PhoneNumber string `json:"phoneNumber"` // Số điện thoại của người dùng
// 	Address     string `json:"address"`     // address của người dùng
// }

// func (c CartServer) RemoveItemInCart(token string, items []string) (*GetCartsResponse, error) {

// 	// Tạo request
// 	url := fmt.Sprintf("%s/v1/Cart/payment_method/%s", c.baseURL, payment_method_id)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

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
// 	var result GetCartsResponse
// 	err = json.Unmarshal(responseBody, &result)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}
// 	if result.Code == 400 {
// 		return nil, fmt.Errorf("product service returned error: %s", result.Message)
// 	}
// 	return &result, nil
// }
// func (c CartServer) CreateCart(token string, params InitPaymentParams) (*InitCartResponse, error) {
// 	// Tạo request
// 	url := fmt.Sprintf("%s/v1/Cart/init", c.baseURL)
// 	body, err := json.Marshal(params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal request body: %w", err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")
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
// 	var initCart InitCartResponse
// 	err = json.Unmarshal(responseBody, &initCart)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w", err)
// 	}
// 	if initCart.Code != 200 {
// 		return nil, fmt.Errorf("Cart service returned error: %s", initCart.Message)
// 	}
// 	return &initCart, nil
// }
