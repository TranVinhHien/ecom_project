package server

import (
	"time"

	config_assets "github.com/TranVinhHien/ecom_payment-service/assets/config"
	server_order "github.com/TranVinhHien/ecom_payment-service/server/order"
	server_product "github.com/TranVinhHien/ecom_payment-service/server/product"
)

// apiClient chứa tất cả các client kết nối tới các service khác
// Pattern: Dependency Injection + Constructor Pattern
type apiClient struct {
	// media   server_media.MediaServer
	product server_product.ProductServer
	order   server_order.OrderServices
}

type ApiServer interface {
	// UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error)
	// UploadSingleImage(token string, file *multipart.FileHeader) (string, error)
	GetSKUs(sku_id string) (*server_product.GetSKUResponse, error)
	GetProductDetail(sku_id string) (*server_product.GetProductDetailResponse, error)
	UpdateProductSKU(token, status string, params []server_product.UpdateProductSKUParams) (*server_product.GetProductDetailResponse, error)
	UpdateOrderCallback(payment_method_id string) error
	// GetTransaction(payment_method_id string) (*server_transaction.GetTransactionsResponse, error)
}

func NewAPIServices(jwt config_assets.ReadENV, timeout time.Duration) ApiServer {
	return &apiClient{
		product: server_product.NewProductServer(jwt.URLProductService, timeout),
		order:   server_order.NewOrderServices(jwt.URLOrderService, timeout),
	}
}

// UploadImages upload nhiều ảnh lên media service
//
//	func (c apiClient) UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error) {
//		return c.media.UploadMultipleImages(token, files)
//	}
//
//	func (c apiClient) UploadSingleImage(token string, file *multipart.FileHeader) (string, error) {
//		return c.media.UploadSingleImage(token, file)
//	}
func (c apiClient) GetSKUs(sku_id string) (*server_product.GetSKUResponse, error) {
	return c.product.GetSKUs(sku_id)
}
func (c apiClient) GetProductDetail(product_id string) (*server_product.GetProductDetailResponse, error) {
	return c.product.GetProductDetail(product_id)
}
func (c apiClient) UpdateProductSKU(token, status string, params []server_product.UpdateProductSKUParams) (*server_product.GetProductDetailResponse, error) {
	return c.product.UpdateProductSKU(token, status, params)
}
func (c apiClient) UpdateOrderCallback(payment_method_id string) error {
	return c.order.UpdateOrderCallbackPayment(payment_method_id)
}
