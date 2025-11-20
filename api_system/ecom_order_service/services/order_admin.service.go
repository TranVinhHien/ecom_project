package services

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	server_product "github.com/TranVinhHien/ecom_order_service/server/product"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
)

// ListShopOrders lấy danh sách đơn hàng của shop với filter và phân trang
func (s *service) ListShopOrders(ctx context.Context, shopID string, status string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError) {

	// Lấy orders từ DB
	orders, err := s.repository.ListShopOrdersSHOP(ctx, db.ListShopOrdersSHOPParams{
		ShopID: shopID,
		Status: db.NullShopOrdersStatus{
			ShopOrdersStatus: db.ShopOrdersStatus(status),
			Valid:            status != "",
		},
		Limit:  int32(query.PageSize),
		Offset: int32(query.PageSize * (query.Page - 1)),
	})
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lấy đơn hàng: %w", err))
	}
	//
	totalElements, err := s.repository.ListShopOrdersSHOPCount(ctx, db.ListShopOrdersSHOPCountParams{
		ShopID: shopID,
		Status: db.NullShopOrdersStatus{
			ShopOrdersStatus: db.ShopOrdersStatus(status),
			Valid:            status != "",
		},
	})
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi đếm đơn hàng: %w", err))
	}
	// Build order summaries
	orderSummaries := make([]services.ShopOrderDetail, len(orders))
	for i, order := range orders {
		// Đếm số items (lấy từ shop_orders và order_items)
		OrderItems, err := s.repository.ListOrderItemsByShopOrderID(ctx, order.ID)
		if err != nil {
			return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lấy order items: %w", err))
		}
		orderSummaries[i] = s.convertDBShopOrderToService(order, OrderItems)
	}
	totalPage := totalElements / int64(query.PageSize)

	result := map[string]interface{}{}
	result["data"] = orderSummaries
	result["currentPage"] = query.Page
	result["totalPages"] = totalPage
	result["totalElements"] = totalElements
	result["limit"] = query.PageSize
	return result, nil
}

func (s *service) GetProductTotalSold(ctx context.Context, productID []string) (map[string]interface{}, *assets_services.ServiceError) {
	totalSold, err := s.repository.GetProductTotalSold(ctx, productID)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lấy tổng số lượng đã bán: %w", err))
	}
	return map[string]interface{}{"data": totalSold}, nil
}

// ShipShopOrder đánh dấu shop order đã được ship
func (s *service) shipShopOrder(ctx context.Context, shopID, user_role, shopOrderID string) *assets_services.ServiceError {
	// Lấy shop order
	shopOrder, err := s.repository.GetShopOrderByID(ctx, shopOrderID)
	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lấy shop orders: %w", err))
	}
	// Verify shop order belongs to this shop
	if shopOrder.ShopID != shopID && user_role != "ROLE_ADMIN" {
		return assets_services.NewError(403, fmt.Errorf("bạn không có quyền cập nhật đơn hàng này"))
	}

	// Verify status is PROCESSING
	if shopOrder.Status != db.ShopOrdersStatusPROCESSING {
		return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái PROCESSING để được giao hàng"))
	}
	shopOrderItems, err := s.repository.ListOrderItemsByShopOrderID(ctx, shopOrder.ID)
	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lấy order items: %w", err))
	}

	err = s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// Update to SHIPPED
		if err := s.repository.UpdateShopOrderStatusToShipped(ctx, db.UpdateShopOrderStatusToShippedParams{
			ID: shopOrder.ID,
		}); err != nil {
			return fmt.Errorf("lỗi khi cập nhật trạng thái giao hàng: %w", err)
		}
		// gửi event tới service product và service vận chuyển để cập nhật kho và tạo đơn vận chuyển
		// gửi tạm trước cho service product vậy
		// hiện tại làm nhanh thì sử lý gọi API luôn cho nó đần

		itemsByShopOrder := []server_product.UpdateProductSKUParams{}
		for _, item := range shopOrderItems {
			itemsByShopOrder = append(itemsByShopOrder, server_product.UpdateProductSKUParams{
				Sku_ID:           item.SkuID,
				QuantityReserved: int(item.Quantity),
			})
		}

		_, err = s.apiServer.UpdateProductSKU(
			s.env.TokenSystem,       // Tên exchange
			string(services.COMMIT), // Routing key
			itemsByShopOrder,
		)
		if err != nil {
			return fmt.Errorf("lỗi khi cập nhật kho sản phẩm: %w", err)
		}

		return nil
	})

	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("lỗi khi cập nhật trạng thái giao hàng: %w", err))
	}
	return nil
}

// UpdateShopOrderStatus cập nhật status của shop order (generic method)
func (s *service) UpdateShopOrderStatus(ctx context.Context, shopID, user_role, shopOrderID, status string, reason string) *assets_services.ServiceError {
	// Lấy shop order
	shopOrder, err := s.repository.GetShopOrderByID(ctx, shopOrderID)

	if err != nil {
		return assets_services.NewError(500, fmt.Errorf("lỗi khi lấy shop order: %w", err))
	}

	// Update status based on target status
	switch status {
	case "PROCESSING":
		// thêm sử lý ở đây
		if shopOrder.Status != db.ShopOrdersStatusAWAITINGPAYMENT {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái AWAITING_PAYMENT để chuyển sang PROCESSING"))
		}
		if err := s.repository.UpdateShopOrderStatusToProcessing(ctx, shopOrder.ID); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
		// thêm sử lý ở đây
	case "CANCELLED":
		if shopOrder.Status != db.ShopOrdersStatusAWAITINGPAYMENT && shopOrder.Status != db.ShopOrdersStatusPROCESSING {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái AWAITING_PAYMENT hoặc PROCESSING để chuyển sang CANCELLED"))
		}
		if err := s.repository.UpdateShopOrderStatusToCancelled(ctx, db.UpdateShopOrderStatusToCancelledParams{
			ID:                 shopOrder.ID,
			CancellationReason: sql.NullString{String: reason, Valid: true},
		}); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
		// đã sử lys song
	case "SHIPPED":
		if shopOrder.Status != db.ShopOrdersStatusPROCESSING {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái PROCESSING để chuyển sang SHIPPED"))
		}
		shipShopOrderErr := s.shipShopOrder(ctx, shopID, user_role, shopOrderID)
		// if err := s.repository.UpdateShopOrderStatusToShipped(ctx, db.UpdateShopOrderStatusToShippedParams{
		// 	ID: shopOrder.ID,
		// });
		if shipShopOrderErr != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", shipShopOrderErr))
		}
		// thêm sử lý ở đây
	case "COMPLETED":
		if shopOrder.Status != db.ShopOrdersStatusSHIPPED {
			return assets_services.NewError(400, fmt.Errorf("đơn hàng phải ở trạng thái SHIPPED để chuyển sang COMPLETED"))
		}
		if err := s.repository.UpdateShopOrderStatusToCompleted(ctx, shopOrder.ID); err != nil {
			return assets_services.NewError(500, fmt.Errorf("lỗi khi cập nhật trạng thái: %w", err))
		}
		// thêm sử lý ở đây
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
