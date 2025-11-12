package server

import (
	"mime/multipart"
	"time"

	config_assets "github.com/TranVinhHien/ecom_order_service/assets/config"
	server_media "github.com/TranVinhHien/ecom_order_service/server/media"
	server_product "github.com/TranVinhHien/ecom_order_service/server/product"
	server_transaction "github.com/TranVinhHien/ecom_order_service/server/transaction"
)

// apiClient chứa tất cả các client kết nối tới các service khác
// Pattern: Dependency Injection + Constructor Pattern
type apiClient struct {
	media       server_media.MediaServer
	product     server_product.ProductServer
	transaction server_transaction.TransactionServer
}

type ApiServer interface {
	UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error)
	UploadSingleImage(token string, file *multipart.FileHeader) (string, error)
	GetSKUs(sku_id string) (*server_product.GetSKUResponse, error)
	GetProductDetail(sku_id string) (*server_product.GetProductDetailResponse, error)
	UpdateProductSKU(token, status string, params []server_product.UpdateProductSKUParams) (*server_product.GetProductDetailResponse, error)
	GetTransaction(payment_method_id string) (*server_transaction.GetTransactionsResponse, error)
	CreateTransaction(token string, params server_transaction.InitPaymentParams) (*server_transaction.InitTransactionResponse, error)
}

func NewAPIServices(jwt config_assets.ReadENV, timeout time.Duration) ApiServer {
	return &apiClient{
		// media:       server_media.NewMediaServer(jwt.URLMediaService, timeout),
		product:     server_product.NewProductServer(jwt.URLProductService, timeout),
		transaction: server_transaction.NewTransactionServer(jwt.URLTransactionService, timeout),
	}
}

// UploadImages upload nhiều ảnh lên media service
func (c apiClient) UploadMultipleImages(token string, files []*multipart.FileHeader) ([]string, error) {
	return c.media.UploadMultipleImages(token, files)
}
func (c apiClient) UploadSingleImage(token string, file *multipart.FileHeader) (string, error) {
	return c.media.UploadSingleImage(token, file)
}
func (c apiClient) GetSKUs(sku_id string) (*server_product.GetSKUResponse, error) {
	return c.product.GetSKUs(sku_id)
}
func (c apiClient) GetProductDetail(product_id string) (*server_product.GetProductDetailResponse, error) {
	return c.product.GetProductDetail(product_id)
}
func (c apiClient) UpdateProductSKU(token, status string, params []server_product.UpdateProductSKUParams) (*server_product.GetProductDetailResponse, error) {
	return c.product.UpdateProductSKU(token, status, params)
}
func (c apiClient) GetTransaction(payment_method_id string) (*server_transaction.GetTransactionsResponse, error) {
	return c.transaction.GetTransaction(payment_method_id)
}
func (c apiClient) CreateTransaction(token string, params server_transaction.InitPaymentParams) (*server_transaction.InitTransactionResponse, error) {
	return c.transaction.CreateTransaction(token, params)
}
