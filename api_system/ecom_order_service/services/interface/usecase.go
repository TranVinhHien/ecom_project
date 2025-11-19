package iservices

import (
	"context"

	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
)

type Orders interface {
	// Customer endpoints
	CreateOrder(ctx context.Context, userID string, token string, req services.CreateOrderRequest) (map[string]interface{}, *assets_services.ServiceError)
	ListUserOrders(ctx context.Context, userID string, query services.QueryFilter, status string) (map[string]interface{}, *assets_services.ServiceError)
	GetOrderDetail(ctx context.Context, userID, user_role, orderCode string) (map[string]interface{}, *assets_services.ServiceError)
	SearchOrdersDetail(ctx context.Context, userID string, filter services.ShopOrderSearchFilter) (map[string]interface{}, *assets_services.ServiceError)

	// Admin/Shop endpoints
	ListShopOrders(ctx context.Context, shopID string, status string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError)
	ShipShopOrder(ctx context.Context, shopID, shopOrderCode string) *assets_services.ServiceError
	UpdateShopOrderStatus(ctx context.Context, shopOrderCode, status string) *assets_services.ServiceError
	CallbackPaymentOnline(ctx context.Context, OrderID string) *assets_services.ServiceError

	// Get total sold quantity for given product IDs
	GetProductTotalSold(ctx context.Context, productID []string) (map[string]interface{}, *assets_services.ServiceError)

	// Kafka event handlers
	HandlePaymentSucceededEvent(ctx context.Context, body services.PaymentSucceededEvent) error
	HandlePaymentFailedEvent(ctx context.Context, body services.PaymentFailedEvent) error
}

// Vouchers defines voucher-related use cases
type Vouchers interface {
	// Admin/Seller endpoints
	CreateVoucher(ctx context.Context, req services.CreateVoucherRequest, user_id, user_type string) *assets_services.ServiceError
	UpdateVoucher(ctx context.Context, voucherID string, user_id string, req services.UpdateVoucherRequest) *assets_services.ServiceError
	ListVouchersForManagement(ctx context.Context, ownerID string, ownerType string, filter services.VoucherManagementFilterRequest) (map[string]interface{}, *assets_services.ServiceError)

	// Customer endpoint
	ListVouchersForUser(ctx context.Context, userID string, filter services.VoucherFilterRequest) (map[string]interface{}, *assets_services.ServiceError)
}

// Comments defines comment-related use cases
type Comments interface {
	// Create a new comment/review
	CreateComment(ctx context.Context, userID string, req services.CreateCommentRequest) *assets_services.ServiceError

	// List comments for a product with pagination
	ListComments(ctx context.Context, req services.ListCommentsRequest) (map[string]interface{}, *assets_services.ServiceError)

	// Check which order items have been reviewed
	CheckReviewedItems(ctx context.Context, req services.CheckReviewedItemsRequest) (*services.CheckReviewedItemsResponse, *assets_services.ServiceError)

	// Get bulk product rating stats for multiple products
	GetBulkProductRatingStats(ctx context.Context, req services.GetBulkProductRatingStatsRequest) (map[string]interface{}, *assets_services.ServiceError)
}
