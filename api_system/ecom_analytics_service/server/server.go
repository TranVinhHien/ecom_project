package server

import (
	"mime/multipart"
	"time"

	config_assets "github.com/TranVinhHien/ecom_analytics_service/assets/config"
	server_media "github.com/TranVinhHien/ecom_analytics_service/server/media"
)

// apiClient chứa tất cả các client kết nối tới các service khác
// Pattern: Dependency Injection + Constructor Pattern
type apiClient struct {
	media server_media.MediaServer
}

type ApiServer interface {
	UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error)
	UploadSingleImage(token string, file *multipart.FileHeader) (string, error)
	RemoveImage(token string, imageURLs []string) error
}

func NewAPIServices(jwt config_assets.ReadENV, timeout time.Duration) ApiServer {
	return &apiClient{
		media: server_media.NewMediaServer("", timeout), // unused host for media service
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
