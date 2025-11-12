package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
)

// ListUserOrders lấy danh sách đơn hàng của user với phân trang
func (s *service) ListUserOrders(ctx context.Context, userID string, query services.QueryFilter, status string) (map[string]interface{}, *assets_services.ServiceError) {

	// Lấy orders từ DB
	orders, err := s.repository.ListShopOrdersByStatus(ctx, db.ListShopOrdersByStatusParams{
		UserID: userID,
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
	totalElements, err := s.repository.ListShopOrdersByStatusCount(ctx, db.ListShopOrdersByStatusCountParams{
		UserID: userID,
		Status: db.NullShopOrdersStatus{
			ShopOrdersStatus: db.ShopOrdersStatus(status),
			Valid:            status != "",
		},
	})
	// Build order summaries
	orderSummaries := make([]services.ShopOrderDetail, len(orders))
	for i, order := range orders {
		// Đếm số items (lấy từ shop_orders và order_items)
		OrderItems, err := s.repository.ListOrderItemsByShopOrderID(ctx, order.ID)
		if err != nil {
			return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lấy order items: %w", err))
		}
		orderSummaries[i] = convertDBShopOrderToService(order, OrderItems)
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

// GetOrderDetail lấy chi tiết đầy đủ của một đơn hàng
func (s *service) GetOrderDetail(ctx context.Context, userID, orderCode string) (map[string]interface{}, *assets_services.ServiceError) {
	// Lấy main order
	order_shop, err := s.repository.GetShopOrderByID(ctx, orderCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, assets_services.NewError(404, fmt.Errorf("không tìm thấy đơn hàng shop"))
		}
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy đơn hàng shop: %w", err))
	}
	OrderItems, err := s.repository.ListOrderItemsByShopOrderID(ctx, order_shop.ID)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lấy order items: %w", err))
	}
	orderSummary := convertDBShopOrderToService(order_shop, OrderItems)

	order, err := s.repository.GetOrderByID(ctx, order_shop.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, assets_services.NewError(404, fmt.Errorf("không tìm thấy đơn hàng"))
		}
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy đơn hàng: %w", err))
	}
	// Verify order thuộc về user
	if order.UserID != userID {
		return nil, assets_services.NewError(403, fmt.Errorf("bạn không có quyền xem đơn hàng này"))
	}

	// Unmarshal shipping address và payment method
	var shippingAddress services.ShippingAddress
	var paymentMethod services.PaymentMethod
	_ = json.Unmarshal([]byte(order.ShippingAddressSnapshot), &shippingAddress)
	_ = json.Unmarshal([]byte(order.PaymentMethodSnapshot), &paymentMethod)

	grandTotal, _ := parseFloat(order.GrandTotal)
	subtotal, _ := parseFloat(order.Subtotal)
	totalShippingFee, _ := parseFloat(order.TotalShippingFee)
	totalDiscount, _ := parseFloat(order.TotalDiscount)
	siteOrderDiscount, _ := parseFloat(order.SiteOrderVoucherDiscount.String)
	siteShippingDiscount, _ := parseFloat(order.SiteShippingVoucherDiscount.String)

	var note *string
	if order.Note.Valid {
		note = &order.Note.String
	}

	var siteOrderVoucherCode, siteShippingVoucherCode *string
	if order.SiteOrderVoucherCode.Valid {
		siteOrderVoucherCode = &order.SiteOrderVoucherCode.String
	}
	if order.SiteShippingVoucherCode.Valid {
		siteShippingVoucherCode = &order.SiteShippingVoucherCode.String
	}

	orderDetail := services.OrderDetail{
		OrderID:   order.ID,
		OrderCode: order.OrderCode,
		UserID:    order.UserID,
		// Status:               string(order.Status),
		GrandTotal:                  grandTotal,
		Subtotal:                    subtotal,
		TotalShippingFee:            totalShippingFee,
		TotalDiscount:               totalDiscount,
		SiteOrderVoucherCode:        siteOrderVoucherCode,
		SiteOrderVoucherDiscount:    siteOrderDiscount,
		SiteShippingVoucherCode:     siteShippingVoucherCode,
		SiteShippingVoucherDiscount: siteShippingDiscount,
		ShippingAddress:             shippingAddress,
		PaymentMethod:               paymentMethod,
		Note:                        note,
		CreatedAt:                   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	response := map[string]interface{}{
		"order":      orderDetail,
		"order_shop": orderSummary,
	}
	return response, nil
}

// SearchOrdersDetail lấy danh sách đơn hàng chi tiết với các bộ lọc
func (s *service) SearchOrdersDetail(ctx context.Context, user_id string, filter services.ShopOrderSearchFilter) (map[string]interface{}, *assets_services.ServiceError) {

	// Chuyển đổi filter sang params cho sqlc
	params := db.SearchShopOrdersParams{
		UserID: user_id,
		Limit:  int32(filter.PageSize),
		Offset: int32(filter.PageSize * (filter.Page - 1)),
	}

	// Set status filter
	if filter.Status != nil {
		params.Status = db.NullShopOrdersStatus{
			ShopOrdersStatus: db.ShopOrdersStatus(*filter.Status),
			Valid:            true,
		}
	}

	// Set shop_id filter
	if filter.ShopID != nil {
		params.ShopID = sql.NullString{String: *filter.ShopID, Valid: true}
	}

	// Set amount filters
	if filter.MinAmount != nil {
		params.MinAmount = sql.NullString{String: fmt.Sprintf("%.2f", *filter.MinAmount), Valid: true}
	}
	if filter.MaxAmount != nil {
		params.MaxAmount = sql.NullString{String: fmt.Sprintf("%.2f", *filter.MaxAmount), Valid: true}
	}

	// Set time filters
	if filter.CreatedFrom != nil {
		params.CreatedFrom = sql.NullTime{Time: *filter.CreatedFrom, Valid: true}
	}
	if filter.CreatedTo != nil {
		params.CreatedTo = sql.NullTime{Time: *filter.CreatedTo, Valid: true}
	}
	if filter.PaidFrom != nil {
		params.PaidFrom = sql.NullTime{Time: *filter.PaidFrom, Valid: true}
	}
	if filter.PaidTo != nil {
		params.PaidTo = sql.NullTime{Time: *filter.PaidTo, Valid: true}
	}
	if filter.ProcessingFrom != nil {
		params.ProcessingFrom = sql.NullTime{Time: *filter.ProcessingFrom, Valid: true}
	}
	if filter.ProcessingTo != nil {
		params.ProcessingTo = sql.NullTime{Time: *filter.ProcessingTo, Valid: true}
	}
	if filter.ShippedFrom != nil {
		params.ShippedFrom = sql.NullTime{Time: *filter.ShippedFrom, Valid: true}
	}
	if filter.ShippedTo != nil {
		params.ShippedTo = sql.NullTime{Time: *filter.ShippedTo, Valid: true}
	}
	if filter.CompletedFrom != nil {
		params.CompletedFrom = sql.NullTime{Time: *filter.CompletedFrom, Valid: true}
	}
	if filter.CompletedTo != nil {
		params.CompletedTo = sql.NullTime{Time: *filter.CompletedTo, Valid: true}
	}
	if filter.CancelledFrom != nil {
		params.CancelledFrom = sql.NullTime{Time: *filter.CancelledFrom, Valid: true}
	}
	if filter.CancelledTo != nil {
		params.CancelledTo = sql.NullTime{Time: *filter.CancelledTo, Valid: true}
	}

	// Set sort_by
	if filter.SortBy != "" {
		params.SortBy = sql.NullString{String: filter.SortBy, Valid: true}
	} else {
		params.SortBy = sql.NullString{String: "created_at", Valid: true}
	}

	// Set OrderBy (default DESC nếu không được cung cấp)
	if filter.SortBy != "" {
		params.SortBy = sql.NullString{String: filter.SortBy, Valid: true}
	} else {
		params.SortBy = sql.NullString{String: "DESC", Valid: true}
	}

	// Count params - same filters but without limit/offset/sort
	countParams := db.SearchShopOrdersCountParams{
		UserID:         params.UserID,
		Status:         params.Status,
		ShopID:         params.ShopID,
		MinAmount:      params.MinAmount,
		MaxAmount:      params.MaxAmount,
		CreatedFrom:    params.CreatedFrom,
		CreatedTo:      params.CreatedTo,
		PaidFrom:       params.PaidFrom,
		PaidTo:         params.PaidTo,
		ProcessingFrom: params.ProcessingFrom,
		ProcessingTo:   params.ProcessingTo,
		ShippedFrom:    params.ShippedFrom,
		ShippedTo:      params.ShippedTo,
		CompletedFrom:  params.CompletedFrom,
		CompletedTo:    params.CompletedTo,
		CancelledFrom:  params.CancelledFrom,
		CancelledTo:    params.CancelledTo,
	}

	// Get total count
	totalElements, err := s.repository.SearchShopOrdersCount(ctx, countParams)
	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi đếm đơn hàng: %w", err))
	}

	// Lấy danh sách shop_orders
	shopOrders, err := s.repository.SearchShopOrders(ctx, params)
	if err != nil {
		return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi tìm kiếm đơn hàng: %w", err))
	}

	// Build kết quả chi tiết cho từng shop_order
	results := make([]map[string]interface{}, 0, len(shopOrders))

	for _, shopOrder := range shopOrders {
		// Lấy order items
		orderItems, err := s.repository.ListOrderItemsByShopOrderID(ctx, shopOrder.ID)
		if err != nil {
			return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy order items: %w", err))
		}

		// Convert shop order detail
		shopOrderDetail := convertDBShopOrderToService(shopOrder, orderItems)

		// Lấy main order
		mainOrder, err := s.repository.GetOrderByID(ctx, shopOrder.OrderID)
		if err != nil {
			return nil, assets_services.NewError(500, fmt.Errorf("lỗi khi lấy main order: %w", err))
		}

		// Parse shipping address và payment method
		var shippingAddress services.ShippingAddress
		var paymentMethod services.PaymentMethod
		_ = json.Unmarshal([]byte(mainOrder.ShippingAddressSnapshot), &shippingAddress)
		_ = json.Unmarshal([]byte(mainOrder.PaymentMethodSnapshot), &paymentMethod)

		// Parse amounts
		grandTotal, _ := parseFloat(mainOrder.GrandTotal)
		subtotal, _ := parseFloat(mainOrder.Subtotal)
		totalShippingFee, _ := parseFloat(mainOrder.TotalShippingFee)
		totalDiscount, _ := parseFloat(mainOrder.TotalDiscount)
		siteOrderDiscount, _ := parseFloat(mainOrder.SiteOrderVoucherDiscount.String)
		siteShippingDiscount, _ := parseFloat(mainOrder.SiteShippingVoucherDiscount.String)

		// Handle nullable fields
		var note *string
		if mainOrder.Note.Valid {
			note = &mainOrder.Note.String
		}

		var siteOrderVoucherCode, siteShippingVoucherCode *string
		if mainOrder.SiteOrderVoucherCode.Valid {
			siteOrderVoucherCode = &mainOrder.SiteOrderVoucherCode.String
		}
		if mainOrder.SiteShippingVoucherCode.Valid {
			siteShippingVoucherCode = &mainOrder.SiteShippingVoucherCode.String
		}

		// Build order detail
		orderDetail := services.OrderDetail{
			OrderID:                     mainOrder.ID,
			OrderCode:                   mainOrder.OrderCode,
			UserID:                      mainOrder.UserID,
			GrandTotal:                  grandTotal,
			Subtotal:                    subtotal,
			TotalShippingFee:            totalShippingFee,
			TotalDiscount:               totalDiscount,
			SiteOrderVoucherCode:        siteOrderVoucherCode,
			SiteOrderVoucherDiscount:    siteOrderDiscount,
			SiteShippingVoucherCode:     siteShippingVoucherCode,
			SiteShippingVoucherDiscount: siteShippingDiscount,
			ShippingAddress:             shippingAddress,
			PaymentMethod:               paymentMethod,
			Note:                        note,
			CreatedAt:                   mainOrder.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:                   mainOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Combine vào response
		result := map[string]interface{}{
			"order":      orderDetail,
			"order_shop": shopOrderDetail,
		}

		results = append(results, result)
	}

	// Calculate total pages
	totalPages := int(totalElements) / filter.PageSize
	if int(totalElements)%filter.PageSize > 0 {
		totalPages++
	}

	// Build response with pagination info
	response := map[string]interface{}{
		"data":          results,
		"currentPage":   filter.Page,
		"pageSize":      filter.PageSize,
		"totalPages":    totalPages,
		"totalElements": totalElements,
	}

	return response, nil
}

func convertDBShopOrderToService(shopOrder db.ShopOrders, Item []db.OrderItems) services.ShopOrderDetail {

	shopSubtotal, _ := parseFloat(shopOrder.Subtotal)
	shippingFee, _ := parseFloat(shopOrder.ShippingFee)
	// shopShippingFee, _ := parseFloat(shopOrder.ShippingFee)
	shopTotalDiscount, _ := parseFloat(shopOrder.TotalDiscount)
	shopTotalAmount, _ := parseFloat(shopOrder.TotalAmount)
	shopVoucherDiscount, _ := parseFloat(shopOrder.ShopVoucherDiscount.String)
	// shopVoucherDiscount, _ := parseFloat(shopOrder.ShopVoucherDiscount)
	// if PaidAt not isZero
	var PaidAt, ProcessingAt, ShippedAt, CompletedAt, CancelledAt *string
	if shopOrder.PaidAt.Valid {
		t := shopOrder.PaidAt.Time.Format("2006-01-02 15:04:05")
		PaidAt = &t
	}
	if shopOrder.ProcessingAt.Valid {
		t := shopOrder.ProcessingAt.Time.Format("2006-01-02 15:04:05")
		ProcessingAt = &t
	}
	if shopOrder.ShippedAt.Valid {
		t := shopOrder.ShippedAt.Time.Format("2006-01-02 15:04:05")
		ShippedAt = &t
	}
	if shopOrder.CompletedAt.Valid {
		t := shopOrder.CompletedAt.Time.Format("2006-01-02 15:04:05")
		CompletedAt = &t
	}
	if shopOrder.CancelledAt.Valid {
		t := shopOrder.CancelledAt.Time.Format("2006-01-02 15:04:05")
		CancelledAt = &t
	}
	// handle item
	var itemDetails []services.OrderItemDetail
	for _, item := range Item {

		temp := services.OrderItemDetail{
			ItemID:    item.ID,
			ProductID: item.ProductID,
			SkuID:     item.SkuID,
			Quantity:  int(item.Quantity),
		}
		originalPrice, _ := parseFloat(item.OriginalUnitPrice)
		finalPrice, _ := parseFloat(item.FinalUnitPrice)
		totalPrice, _ := parseFloat(item.TotalPrice)
		temp.OriginalUnitPrice = originalPrice
		temp.FinalUnitPrice = finalPrice
		temp.TotalPrice = totalPrice
		temp.ProductName = item.ProductNameSnapshot
		if item.ProductImageSnapshot.Valid {
			temp.ProductImage = &item.ProductImageSnapshot.String
		}
		temp.SkuAttributes = item.SkuAttributesSnapshot.String
		promotions, _ := item.PromotionsSnapshot.MarshalJSON()
		_ = json.Unmarshal(promotions, &temp.PromotionsSnapshot)
		itemDetails = append(itemDetails, temp)
	}
	return services.ShopOrderDetail{
		ShopOrderID:         shopOrder.ID,
		ShopOrderCode:       shopOrder.ShopOrderCode,
		ShopID:              shopOrder.ShopID,
		Status:              string(shopOrder.Status),
		Subtotal:            shopSubtotal,
		ShippingFee:         shippingFee,
		TotalDiscount:       shopTotalDiscount,
		TotalAmount:         shopTotalAmount,
		ShopVoucherCode:     &shopOrder.ShopVoucherCode.String,
		ShopVoucherDiscount: shopVoucherDiscount,
		ShippingMethod:      &shopOrder.ShippingMethod.String,
		TrackingCode:        &shopOrder.TrackingCode.String,
		Items:               itemDetails,
		CreatedAt:           shopOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           shopOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
		PaidAt:              PaidAt,
		ProcessingAt:        ProcessingAt,
		ShippedAt:           ShippedAt,
		CompletedAt:         CompletedAt,
		CancelledAt:         CancelledAt,
	}
}
