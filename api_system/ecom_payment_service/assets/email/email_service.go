package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BrevoEmailService là một Adapter triển khai việc gửi email qua Brevo API.
// Nó tuân thủ SSoT (tệp Postman).
type BrevoEmailService struct {
	httpClient  *http.Client
	apiKey      string
	senderEmail string
	senderName  string
	apiEndpoint string
}

// NewBrevoEmailService là hàm khởi tạo (constructor) cho service.
// Dependencies (apiKey, sender...) nên được inject từ config.
func NewBrevoEmailService(apiKey, senderEmail, senderName string) *BrevoEmailService {
	return &BrevoEmailService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Đặt timeout cho các request HTTP
		},
		apiKey:      apiKey,
		senderEmail: senderEmail,
		senderName:  senderName,
		apiEndpoint: "https://api.brevo.com/v3/smtp/email", //
	}
}

// === CÁC STRUCT NỘI BỘ CHO BREVO API ===
// Dựa trên cấu trúc body của Postman

type brevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type brevoRecipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type brevoRequestBody struct {
	Sender      brevoSender      `json:"sender"`
	To          []brevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

// SendEmail thực hiện việc gửi email.
// Hàm này nhận vào các tham số bạn yêu cầu và thêm một số tham số
// cần thiết cho một email chuyên nghiệp (như Subject).
func (s *BrevoEmailService) SendEmail(
	ctx context.Context,
	recipientEmail string,
	recipientName string, // Tên người nhận (ví dụ: "Nguyễn Văn A")
	subject string,
	htmlBody string, // Đây là string từ GeneratePaymentSuccessEmail
) error {

	// 1. Chuẩn bị Request Body theo chuẩn Brevo
	requestBody := brevoRequestBody{
		Sender: brevoSender{
			Name:  s.senderName,
			Email: s.senderEmail,
		},
		To: []brevoRecipient{
			{
				Name:  recipientName,
				Email: recipientEmail,
			},
		},
		Subject:     subject,
		HTMLContent: htmlBody,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		// Lỗi này là lỗi hệ thống (code), không phải lỗi người dùng
		return fmt.Errorf("brevo: không thể marshal request body: %w", err)
	}

	// 2. Tạo HTTP Request
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("brevo: không thể tạo http request: %w", err)
	}

	// 3. Set Headers theo yêu cầu của Brevo
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	// 4. Gửi Request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("brevo: lỗi khi gửi request: %w", err)
	}
	defer resp.Body.Close()

	// 5. Kiểm tra Response
	// Brevo trả về 201 Created khi thành công
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Đọc body lỗi để debug (nếu có)
		var errResp bytes.Buffer
		errResp.ReadFrom(resp.Body)
		return fmt.Errorf("brevo: API trả về lỗi status %d: %s", resp.StatusCode, errResp.String())
	}

	// Gửi email thành công
	return nil
}
