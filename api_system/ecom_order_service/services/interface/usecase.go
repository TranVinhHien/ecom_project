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
	GetOrderDetail(ctx context.Context, userID, orderCode string) (map[string]interface{}, *assets_services.ServiceError)
	SearchOrdersDetail(ctx context.Context, userID string, filter services.ShopOrderSearchFilter) (map[string]interface{}, *assets_services.ServiceError)

	// Admin/Shop endpoints
	ListShopOrders(ctx context.Context, shopID string, status *string, page, limit int, dateFrom, dateTo *string) (map[string]interface{}, *assets_services.ServiceError)
	ShipShopOrder(ctx context.Context, shopID, shopOrderCode string, req services.ShipOrderRequest) *assets_services.ServiceError
	UpdateShopOrderStatus(ctx context.Context, shopOrderCode, status string) *assets_services.ServiceError
	CallbackPaymentOnline(ctx context.Context, OrderID string) *assets_services.ServiceError

	// Kafka event handlers
	HandlePaymentSucceededEvent(ctx context.Context, body services.PaymentSucceededEvent) error
	HandlePaymentFailedEvent(ctx context.Context, body services.PaymentFailedEvent) error
}

// Vouchers defines voucher-related use cases
type Vouchers interface {
	// Admin endpoints
	CreateVoucher(ctx context.Context, req services.CreateVoucherRequest) *assets_services.ServiceError
	UpdateVoucher(ctx context.Context, voucherID string, req services.UpdateVoucherRequest) *assets_services.ServiceError

	// Customer endpoint
	ListVouchersForUser(ctx context.Context, userID string, filter services.VoucherFilterRequest) (map[string]interface{}, *assets_services.ServiceError)
}
