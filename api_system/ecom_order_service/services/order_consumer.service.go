package services

import (
	"context"
	"database/sql"
	"log"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	server_product "github.com/TranVinhHien/ecom_order_service/server/product"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
)

func (s *service) HandlePaymentFailedEvent(ctx context.Context, body services.PaymentFailedEvent) error {

	log.Printf("Processing payment_failed event for order %s", body.OrderID)

	// 1. Tìm các shop_orders con đang ở trạng thái AWAITING_PAYMENT
	shopOrderIDs, err := s.repository.ListShopOrdersByOrderID(ctx, db.ListShopOrdersByOrderIDParams{
		OrderID: body.OrderID,
		Status:  db.NullShopOrdersStatus{ShopOrdersStatus: db.ShopOrdersStatusAWAITINGPAYMENT, Valid: true},
	}) //
	if err != nil {
		// log.Printf("Error fetching shop orders: %v", err)
		return err // Thử lại (NACK)
	}

	if len(shopOrderIDs) == 0 {
		log.Printf("No shop orders found in 'AWAITING_PAYMENT' for order %s. Already processed or COD.", body.OrderID)
		return nil // Hoàn tất (ACK)
	}
	itemIDs := make([]string, len(shopOrderIDs))
	for i, so := range shopOrderIDs {
		itemIDs[i] = so.ID
	}
	// 2. Lấy tất cả items thuộc các shop_orders này
	items, err := s.repository.GetOrderItemsByShopOrderIDs(ctx, itemIDs) //
	if err != nil {
		log.Printf("Error fetching order items: %v", err)
		return err // Thử lại (NACK)
	}

	// Code đơn giản (nếu không dùng transaction cho 1 lệnh update)
	params := db.CancelShopOrdersByIDsParams{
		CancellationReason: sql.NullString{String: "Thanh toán thất bại hoặc hết hạn", Valid: true},
		ShopOrderIds:       itemIDs,
	}
	err = s.repository.ExecTS(ctx, func(tx db.Querier) error {
		err = tx.CancelShopOrdersByIDs(ctx, params) //
		if err != nil {
			log.Printf("Error cancelling shop orders: %v", err)
			return err // Thử lại (NACK)
		}

		// 4. Gửi sự kiện 'order_cancelled' cho Product Service (để trả kho)
		// Nhóm các items theo shop_order_id
		itemsByShopOrder := []server_product.UpdateProductSKUParams{}
		for _, item := range items {
			itemsByShopOrder = append(itemsByShopOrder, server_product.UpdateProductSKUParams{
				Sku_ID:           item.SkuID,
				QuantityReserved: int(item.Quantity),
			})
		}
		_, err = s.apiServer.UpdateProductSKU(
			s.env.TokenSystem, // Tên exchange
			"rollback",        // Routing key
			itemsByShopOrder,
		)
		if err != nil {
			log.Printf("Error updating product SKUs for rollback: %v", err)
			return err // Thử lại (NACK)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error in transaction while cancelling shop orders: %v", err)
		return err // Thử lại (NACK)
	}

	log.Printf("Successfully cancelled %d shop orders for main order %s", len(shopOrderIDs), body.OrderID)
	return nil // Hoàn tất (ACK)
}
func (s *service) HandlePaymentSucceededEvent(ctx context.Context, body services.PaymentSucceededEvent) error {
	log.Printf("Processing payment_succeeded event for order %s", body.OrderID)

	// 1. Tìm các shop_orders con đang ở trạng thái AWAITING_PAYMENT
	shopOrderItems, err := s.repository.ListShopOrdersByOrderID(ctx, db.ListShopOrdersByOrderIDParams{
		OrderID: body.OrderID,
		Status:  db.NullShopOrdersStatus{ShopOrdersStatus: db.ShopOrdersStatusAWAITINGPAYMENT, Valid: true},
	}) //
	if err != nil {
		log.Printf("Error fetching shop orders: %v", err)
		return err // Thử lại (NACK)
	}

	if len(shopOrderItems) == 0 {
		log.Printf("No shop orders found in 'AWAITING_PAYMENT' for order %s. Already processed or COD.", body.OrderID)
		return nil // Hoàn tất (ACK)
	}
	shopOrderID := make([]string, len(shopOrderItems))
	for i, so := range shopOrderItems {
		shopOrderID[i] = so.ID
	}

	// 3. Cập nhật CSDL (phải nằm trong 1 DB Transaction, dù sqlc không trực tiếp hỗ trợ)
	// update toàn bộ orderShop
	s.repository.ExecTS(ctx, func(tx db.Querier) error {
		for _, id := range shopOrderID {
			err = s.repository.UpdateShopOrderStatusToProcessing(ctx, id) //
			if err != nil {
				log.Printf("Error updating shop order %s: %v", id, err)
				return err // Thử lại (NACK)
			}
		}
		return nil
	})

	log.Printf("Successfully update %d shop orders for main order %s", len(shopOrderID), body.OrderID)
	return nil
}

// sử lý khi phía admin xác nhận đã hoàn thành đơn hàng . đây sẽ là 1 API gọi từ phía seller gọi tới đây rồi sẽ gọi tới phía vận chuyển để gửi tin nhắn cho nó giao
func (s *service) HandleOrderDeliveredEvent(ctx context.Context, body interface{}) error {
	return nil
}

// sử lý khi phía vận chuyển xác nhận đã giao hàng thì sẽ cập nhật ở hàm này. Hàm này sẽ bắt đầu cộng riêng ra chứ không cộng ở trong kia.
func (s *service) HandleOrderReceivedEvent(ctx context.Context, body interface{}) error {
	return nil
}
