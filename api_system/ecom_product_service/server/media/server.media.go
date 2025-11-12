package server_media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type MediaServer struct {
	baseURL    string
	httpClient *http.Client
}

type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Result  struct {
		Urls []string `json:"urls"`
	} `json:"result"` // URLs của các ảnh đã upload
}

// NewMediaServer tạo mới media client với dependency injection
func NewMediaServer(baseURL string, timeout time.Duration) MediaServer {
	return MediaServer{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// UploadImages upload nhiều ảnh lên media service
// files: slice các multipart.FileHeader từ Gin context
func (c MediaServer) uploadImages(token string, files []*multipart.FileHeader) (*UploadResponse, error) {
	// Tạo buffer để chứa multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Thêm các file vào form data
	for _, fileHeader := range files {
		// Tạo form file với field name "media"
		part, err := writer.CreateFormFile("media", fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		// Mở file từ FileHeader
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}

		// Copy nội dung file vào form
		_, err = io.Copy(part, file)
		file.Close() // Đóng file ngay sau khi copy xong

		if err != nil {
			return nil, fmt.Errorf("failed to write file content %s: %w", fileHeader.Filename, err)
		}
	}

	// Đóng writer để hoàn thành multipart form
	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Tạo request
	url := fmt.Sprintf("%s/v1/media/uploads", c.baseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
	var uploadResp UploadResponse
	err = json.Unmarshal(responseBody, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &uploadResp, nil
}
func (c MediaServer) UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error) {
	resp, err := c.uploadImages(token, files)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("upload failed: %s", resp.Error)
	}

	return resp.Result.Urls, nil
}

// UploadSingleImage upload 1 ảnh đơn lẻ
func (c MediaServer) UploadSingleImage(token string, file *multipart.FileHeader) (string, error) {
	// Tạo slice chứa 1 file
	files := []*multipart.FileHeader{file}

	resp, err := c.uploadImages(token, files)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", fmt.Errorf("upload failed: %s", resp.Error)
	}

	return resp.Result.Urls[0], nil
}
