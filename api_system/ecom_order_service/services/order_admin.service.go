package services

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
)

// ListShopOrders lấy danh sách đơn hàng của shop với filter và phân trang
func (s *service) ListShopOrders(ctx context.Context, shopID string, status *string, page, limit int, dateFrom, dateTo *string) (map[string]interface{}, *assets_services.ServiceError) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Lấy shop orders từ DB
	shopOrders, err := s.repository.ListShopOrdersByShopIDPaged(ctx, db.ListShopOrdersByShopIDPagedParams{
		ShopID: shopID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy đơn hàng: %w", err))
	}

	// Filter by status if provided
	filteredOrders := make([]db.ShopOrders, 0)
	for _, order := range shopOrders {
		if status != nil && string(order.Status) != *status {
			continue
		}
		// TODO: Filter by date range if dateFrom/dateTo provided
		filteredOrders = append(filteredOrders, order)
	}

	// Build response
	shopOrderSummaries := make([]map[string]interface{}, len(filteredOrders))
	for i, shopOrder := range filteredOrders {
		// Lấy items count
		items, _ := s.repository.ListOrderItemsByShopOrderID(ctx, shopOrder.ID)

		subtotal, _ := parseFloat(shopOrder.Subtotal)
		totalAmount, _ := parseFloat(shopOrder.TotalAmount)

		shopOrderSummaries[i] = map[string]interface{}{
			"shopOrderId":   shopOrder.ID,
			"shopOrderCode": shopOrder.ShopOrderCode,
			"orderId":       shopOrder.OrderID,
			"status":        string(shopOrder.Status),
			"subtotal":      subtotal,
			"totalAmount":   totalAmount,
			"itemCount":     len(items),
			"createdAt":     shopOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	response := map[string]interface{}{
		"shopOrders": shopOrderSummaries,
		"totalCount": len(filteredOrders),
		"page":       page,
		"limit":      limit,
	}

	return response, nil
}

// ShipShopOrder đánh dấu shop order đã được ship
func (s *service) ShipShopOrder(ctx context.Context, shopID, shopOrderCode string, req services.ShipOrderRequest) *assets_services.ServiceError {
	// Lấy shop order
	shopOrders, err := s.repository.ListShopOrdersByShopIDPaged(ctx, db.ListShopOrdersByShopIDPagedParams{
		ShopID: shopID,
		Limit:  1000, // Get all to find by code
		Offset: 0,
	})
	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lấy shop orders: %w", err))
	}

	var shopOrder *db.ShopOrders
	for _, so := range shopOrders {
		if so.ShopOrderCode == shopOrderCode {
			shopOrder = &so
			break
		}
	}

	if shopOrder == nil {
		return assets_services.NewError(404, fmt.Errorf("không tìm thấy đơn hàng với mã %s", shopOrderCode))
	}

	// Verify shop order belongs to this shop
	if shopOrder.ShopID != shopID {
		return assets_services.NewError(403, fmt.Errorf("bạn không có quyền cập nhật đơn hàng này"))
	}

	// Verify status is PROCESSING
	if shopOrder.Status != db.ShopOrdersStatusPROCESSING {
		return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái PROCESSING để được giao hàng"))
	}

	// Update to SHIPPED
	if err := s.repository.UpdateShopOrderStatusToShipped(ctx, db.UpdateShopOrderStatusToShippedParams{
		ID: shopOrder.ID,
		// TrackingCode:   req.TrackingCode,
		// ShippingMethod: req.ShippingMethod,
	}); err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái đơn hàng thành SHIPPED: %w", err))
	}

	return nil
}

// UpdateShopOrderStatus cập nhật status của shop order (generic method)
func (s *service) UpdateShopOrderStatus(ctx context.Context, shopOrderID, status string) *assets_services.ServiceError {
	// Lấy shop order
	shopOrder, err := s.repository.GetShopOrderByID(ctx, shopOrderID)

	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lấy shop order: %w", err))
	}

	// Update status based on target status
	switch status {
	case "PROCESSING":
		if shopOrder.Status != db.ShopOrdersStatusAWAITINGPAYMENT {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái AWAITING_PAYMENT để chuyển sang PROCESSING"))
		}
		if err := s.repository.UpdateShopOrderStatusToProcessing(ctx, shopOrder.ID); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
	case "CANCELLED":
		if shopOrder.Status != db.ShopOrdersStatusAWAITINGPAYMENT && shopOrder.Status != db.ShopOrdersStatusPROCESSING {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái AWAITING_PAYMENT hoặc PROCESSING để chuyển sang CANCELLED"))
		}
		if err := s.repository.UpdateShopOrderStatusToCancelled(ctx, db.UpdateShopOrderStatusToCancelledParams{
			ID:                 shopOrder.ID,
			CancellationReason: sql.NullString{String: "Cancelled by shop", Valid: true},
		}); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
	case "SHIPPED":
		if shopOrder.Status != db.ShopOrdersStatusPROCESSING {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái PROCESSING để chuyển sang SHIPPED"))
		}
		if err := s.repository.UpdateShopOrderStatusToShipped(ctx, db.UpdateShopOrderStatusToShippedParams{
			ID: shopOrder.ID,
			// TrackingCode:   req.TrackingCode,
			// ShippingMethod: req.ShippingMethod,
		}); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
	case "COMPLETED":
		if shopOrder.Status != db.ShopOrdersStatusSHIPPED {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái SHIPPED để chuyển sang COMPLETED"))
		}
		if err := s.repository.UpdateShopOrderStatusToCompleted(ctx, shopOrder.ID); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
	case "REFUNDED":
		if shopOrder.Status != db.ShopOrdersStatusCOMPLETED {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái COMPLETED để chuyển sang REFUNDED"))
		}
		if err := s.repository.UpdateShopOrderStatusToRefunded(ctx, shopOrder.ID); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}

	default:
		return assets_services.NewError(400, fmt.Errorf("trạng thái không hợp lệ: %s", status))
	}
	return nil
}

// UpdateShopOrderStatus cập nhật status của shop order (generic method)
func (s *service) CallbackPaymentOnline(ctx context.Context, OrderID string) *assets_services.ServiceError {
	// Lấy shop order
	shopOrders, err := s.repository.ListShopOrdersByOrderID(ctx, db.ListShopOrdersByOrderIDParams{
		OrderID: OrderID,
		// Status:  sql.NullString{String: string(db.ShopOrdersStatusAWAITINGPAYMENT), Valid: true},
	})
	// s.repository.CancelShopOrdersByIDs()
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("lỗi khi lấy shop order: %w", err))
	}
	for _, shopOrder := range shopOrders {
		// Update status based on target status
		if shopOrder.Status == db.ShopOrdersStatusAWAITINGPAYMENT {
			if err := s.repository.UpdateShopOrderStatusToProcessing(ctx, shopOrder.ID); err != nil {
				return assets_services.NewError(400, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
			}
		}
	}
	return nil
}

// // Helper: cập nhật status của parent order dựa trên các shop orders
// func (s *service) updateParentOrderStatus(ctx context.Context, orderID string) error {
// 	// Lấy parent order
// 	order, err := s.repository.GetOrderByID(ctx, orderID)
// 	if err != nil {
// 		return err
// 	}

// 	// Lấy tất cả shop orders
// 	shopOrders, err := s.repository.ListShopOrdersByOrderID(ctx, orderID)
// 	if err != nil {
// 		return err
// 	}

// 	if len(shopOrders) == 0 {
// 		return nil
// 	}

// 	// Calculate new status
// 	allCompleted := true
// 	allCancelled := true
// 	anyShipped := false
// 	anyCancelled := false

// 	for _, so := range shopOrders {
// 		if so.Status != db.ShopOrdersStatusCOMPLETED {
// 			allCompleted = false
// 		}
// 		if so.Status != db.ShopOrdersStatusCANCELLED {
// 			allCancelled = false
// 		}
// 		if so.Status == db.ShopOrdersStatusSHIPPED || so.Status == db.ShopOrdersStatusCOMPLETED {
// 			anyShipped = true
// 		}
// 		if so.Status == db.ShopOrdersStatusCANCELLED {
// 			anyCancelled = true
// 		}
// 	}

// 	// var newStatus db.OrdersStatus
// 	// if allCompleted {
// 	// 	newStatus = db.OrdersStatusCOMPLETED
// 	// } else if allCancelled {
// 	// 	newStatus = db.OrdersStatusCANCELLED
// 	// } else if anyCancelled {
// 	// 	// newStatus = db.OrdersStatusPARTIALLY_CANCELLED
// 	// } else if anyShipped {
// 	// 	// newStatus = db.OrdersStatusPARTIALLY_SHIPPED
// 	// } else {
// 	// 	newStatus = db.OrdersStatusPROCESSING
// 	// }

// 	// // Update if changed
// 	// if order.Status != newStatus {
// 	// 	return s.repository.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
// 	// 		ID:     orderID,
// 	// 		Status: newStatus,
// 	// 	})
// 	// }

// 	return nil
// }
