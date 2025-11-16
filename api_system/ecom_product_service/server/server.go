package server

import (
	"mime/multipart"
	"time"

	config_assets "github.com/TranVinhHien/ecom_product_service/assets/config"
	server_media "github.com/TranVinhHien/ecom_product_service/server/media"
	server_order "github.com/TranVinhHien/ecom_product_service/server/order"
)

// apiClient chứa tất cả các client kết nối tới các service khác
// Pattern: Dependency Injection + Constructor Pattern
type apiClient struct {
	media server_media.MediaServer
	order server_order.OrderServer
}

type ApiServer interface {
	UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error)
	UploadSingleImage(token string, file *multipart.FileHeader) (string, error)
	RemoveImage(token string, imageURLs []string) error
	GetBulkProductRatingStats(productIDs []string) (map[string]server_order.ProductRatingStatsItem, error)
}

func NewAPIServices(config config_assets.ReadENV, timeout time.Duration) ApiServer {
	return &apiClient{
		media: server_media.NewMediaServer("", timeout), // unused host for media service
		order: server_order.NewOrderServer(config.OrderServiceURL, timeout),
	}
}

// UploadImages upload nhiều ảnh lên media service
func (c apiClient) UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error) {
	return c.media.UploadMultipleImages(token, files)
}
func (c apiClient) UploadSingleImage(token string, file *multipart.FileHeader) (string, error) {
	return c.media.UploadSingleImage(token, file)
}
func (c apiClient) RemoveImage(token string, imageURLs []string) error {
	return nil
	// return c.media.UploadMultipleImages(token, files)
}

// GetBulkProductRatingStats lấy thống kê đánh giá cho nhiều sản phẩm từ order service
func (c apiClient) GetBulkProductRatingStats(productIDs []string) (map[string]server_order.ProductRatingStatsItem, error) {
	return c.order.GetBulkProductRatingStats(productIDs)
}
