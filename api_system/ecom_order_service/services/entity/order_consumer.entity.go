package services

import (
	"context"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
)

// EventPublisher là interface tối thiểu cho RabbitMQ
type EventPublisher interface {
	PublishEvent(ctx context.Context, exchange, routingKey string, body []byte) error
}

// PaymentFailedEvent là cấu trúc message nhận được
type PaymentFailedEvent struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"` // Đây là order_id (cha)
	Reason        string `json:"reason"`
}
type PaymentSucceededEvent struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"` // Đây là order_id (cha)
	// Amount        string `json:"amount"`
}

// OrderItemInfo là thông tin item để gửi cho Product Service
type OrderItemInfo struct {
	SKUID    string `json:"sku_id"`
	Quantity int    `json:"quantity"`
}

// OrderCancelledEvent là message gửi đi cho Product Service
type OrderCancelledEvent struct {
	ShopOrderID string          `json:"shop_order_id"`
	Items       []OrderItemInfo `json:"items"`
}

// PaymentEventHandler chứa các dependency
type PaymentEventHandler struct {
	repo      db.Querier     // Interface sqlc
	publisher EventPublisher // Publisher để gửi event 'order_cancelled'
}
