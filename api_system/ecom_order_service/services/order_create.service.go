package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	db "github.com/TranVinhHien/ecom_order_service/db/sqlc"
	server_product "github.com/TranVinhHien/ecom_order_service/server/product"
	server_transaction "github.com/TranVinhHien/ecom_order_service/server/transaction"
	assets_services "github.com/TranVinhHien/ecom_order_service/services/assets"
	services "github.com/TranVinhHien/ecom_order_service/services/entity"
	"github.com/google/uuid"
)

// CreateOrder tạo đơn hàng mới với logic phức tạp multi-shop
func (s *service) CreateOrder(ctx context.Context, userID string, token string, req services.CreateOrderRequest) (map[string]interface{}, *assets_services.ServiceError) {
	// Bước 1: Validate request
	if err := s.validateCreateOrderRequest(req); err != nil {
		return nil, err
	}

	// Bước 2: Lấy thông tin sản phẩm từ Product Service
	productInfoMap, err := s.fetchProductInfoForOrder(ctx, req.Items)
	if err != nil {
		return nil, err

	}
	// đặt tên hàm sai làm biến sửa lại
	paymentMethod, errors := s.apiServer.GetTransaction(req.PaymentMethod_ID)
	if errors != nil {
		return nil, &assets_services.ServiceError{
			Code: 400,
			Err:  fmt.Errorf("lỗi khi lấy phương thức thanh toán: %s", errors.Error()),
		}
	}
	// Bước 3: Validate tồn kho
	if err := s.validateStock(productInfoMap, req.Items); err != nil {
		return nil, err
	}

	// Bước 4: Nhóm items theo shop
	shopItemsMap := s.groupItemsByShop(req.Items)

	// Bước 5: Tạo order ID và order code
	orderID := uuid.New().String()
	orderCode := s.generateOrderCode()

	// Bước 6: Tạo các shop orders và tính toán tổng tiền
	shopOrders, grandTotal, subtotal, totalShippingFee, totalDiscount, voucherTotalSite, voucherShippingSite, voucherTotalDiscount, voucherShippingDiscount, err := s.createShopOrdersWithItems(
		ctx, userID, orderID, shopItemsMap, productInfoMap, req.VoucherShop, req.VoucherSiteID, req.VoucherShippingID)
	if err != nil {
		return nil, err
	}

	// Bước 7: Reserve stock cho tất cả items
	reservations := s.buildStockReservations(req.Items)
	if err := s.reserveStockForOrder(ctx, token, reservations); err != nil {
		return nil, assets_services.NewError(409, fmt.Errorf("lỗi khi cập nhật số lượng tồn kho: %w", err))
	}

	// Bước 8: Lưu order vào database trong transaction
	var paymentURL *string
	saveErr := s.repository.ExecTS(ctx, func(tx db.Querier) error {
		// Lưu main order
		shippingJSON, _ := json.Marshal(req.ShippingAddress)
		paymentJSON, _ := json.Marshal(paymentMethod.Result.Data)

		var noteSQL sql.NullString
		if req.Note != nil {
			noteSQL = sql.NullString{String: *req.Note, Valid: true}
		}
		discountTotalSite := 0.0
		discountShippingSite := 0.0
		if voucherTotalSite != nil {
			discountTotalSite, _ = countDiscountAmount(*voucherTotalSite, subtotal)
			// trừ voucher site
			err := s.releaseStockForVoucher(ctx, userID, voucherTotalSite.ID, discountTotalSite)
			if err != nil {
				return fmt.Errorf("lỗi khi giải phóng voucher: %w", err)
			}
		}
		if voucherShippingSite != nil {
			discountShippingSite, _ = countDiscountAmount(*voucherShippingSite, subtotal)
			// trừ voucher giao hàng
			err := s.releaseStockForVoucher(ctx, userID, voucherShippingSite.ID, discountShippingSite)
			if err != nil {
				return fmt.Errorf("lỗi khi giải phóng voucher: %w", err)
			}
		}
		var siteOrderVoucherCode, siteOrderVoucherDiscount, siteShippingVoucherCode, siteShippingVoucherDiscount sql.NullString

		if voucherTotalSite != nil {
			siteOrderVoucherCode = sql.NullString{String: voucherTotalSite.VoucherCode, Valid: true}
			siteOrderVoucherDiscount = sql.NullString{String: fmt.Sprintf("%.2f", discountTotalSite), Valid: true}
		}
		if voucherShippingSite != nil {
			siteShippingVoucherCode = sql.NullString{String: voucherShippingSite.VoucherCode, Valid: true}
			siteShippingVoucherDiscount = sql.NullString{String: fmt.Sprintf("%.2f", discountShippingSite), Valid: true}
		}

		// trừ voucher giao hàng
		if err := tx.CreateOrder(ctx, db.CreateOrderParams{
			ID:        orderID,
			OrderCode: orderCode,
			UserID:    userID,
			// Status:                      initialStatus,
			GrandTotal:                  fmt.Sprintf("%.2f", grandTotal),
			Subtotal:                    fmt.Sprintf("%.2f", subtotal),
			TotalShippingFee:            fmt.Sprintf("%.2f", totalShippingFee),
			TotalDiscount:               fmt.Sprintf("%.2f", totalDiscount),
			SiteOrderVoucherCode:        siteOrderVoucherCode,
			SiteOrderVoucherDiscount:    siteOrderVoucherDiscount,
			SiteShippingVoucherCode:     siteShippingVoucherCode,
			SiteShippingVoucherDiscount: siteShippingVoucherDiscount,
			ShippingAddressSnapshot:     json.RawMessage(shippingJSON),
			PaymentMethodSnapshot:       json.RawMessage(paymentJSON),
			Note:                        noteSQL,
		}); err != nil {
			return fmt.Errorf("lỗi khi tạo đơn hàng: %w", err)
		}

		// Lưu shop orders và items
		for _, shopOrder := range shopOrders {
			status := shopOrder.Status
			if paymentMethod.Result.Data.Type == "OFFLINE" {
				status = services.ShopOrderStatusProcessing
			}
			// create shop_order_voucher_usage_history --- IGNORE ---
			if shopOrder.DiscountCode != "" {
				voucherShop, err := tx.GetVoucherForValidation(ctx, shopOrder.DiscountCode)
				if err != nil {
					return fmt.Errorf("lỗi khi lấy voucher: %w", err)
				}
				discountOrderShop, _ := countDiscountAmount(voucherShop, shopOrder.Subtotal)
				// trừ voucher site
				errrors := s.releaseStockForVoucher(ctx, userID, voucherShop.ID, discountOrderShop)
				if errrors != nil {
					return fmt.Errorf("lỗi khi giải phóng voucher: %w", errrors)
				}
			}
			if err := tx.CreateShopOrder(ctx, db.CreateShopOrderParams{
				ID:                  shopOrder.ShopOrderID,
				ShopOrderCode:       shopOrder.ShopOrderCode,
				OrderID:             orderID,
				ShopID:              shopOrder.ShopID,
				Status:              db.ShopOrdersStatus(status),
				Subtotal:            fmt.Sprintf("%.2f", shopOrder.Subtotal),
				TotalDiscount:       fmt.Sprintf("%.2f", shopOrder.TotalDiscount),
				TotalAmount:         fmt.Sprintf("%.2f", shopOrder.TotalAmount),
				ShippingFee:         fmt.Sprintf("%.2f", shopOrder.ShippingFee),
				ShopVoucherCode:     sql.NullString{String: shopOrder.DiscountCode, Valid: shopOrder.DiscountCode != ""},
				ShopVoucherDiscount: sql.NullString{String: fmt.Sprintf("%.2f", shopOrder.TotalDiscount), Valid: shopOrder.DiscountCode != ""},
				ProcessingAt:        sql.NullTime{Time: time.Now(), Valid: status == services.ShopOrderStatusProcessing},
			}); err != nil {
				return fmt.Errorf("lỗi khi tạo shop order: %w", err)
			}

			// Lưu order items
			for _, item := range shopOrder.Items {
				var productImage sql.NullString
				if item.ProductImage != nil {
					productImage = sql.NullString{String: *item.ProductImage, Valid: true}
				}

				// skuAttrsJSON, _ := json.Marshal(item.SkuAttributes)

				if err := tx.CreateOrderItem(ctx, db.CreateOrderItemParams{
					ID:                    item.ItemID,
					ShopOrderID:           shopOrder.ShopOrderID,
					ProductID:             item.ProductID,
					SkuID:                 item.SkuID,
					Quantity:              uint32(item.Quantity),
					OriginalUnitPrice:     fmt.Sprintf("%.2f", item.OriginalUnitPrice),
					FinalUnitPrice:        fmt.Sprintf("%.2f", item.FinalUnitPrice),
					TotalPrice:            fmt.Sprintf("%.2f", item.TotalPrice),
					PromotionsSnapshot:    json.RawMessage("{}"),
					ProductNameSnapshot:   item.ProductName,
					ProductImageSnapshot:  productImage,
					SkuAttributesSnapshot: sql.NullString{String: string(item.SkuAttributes), Valid: true},
				}); err != nil {
					return fmt.Errorf("lỗi khi tạo order item: %w", err)
				}
			}
		}

		// // // Bước 9: Tạo transaction
		// res, err := s.apiServer.CreateTransaction(token, server_transaction.InitPaymentParams{
		params := createInitPaymentParams(orderID, req.PaymentMethod_ID, grandTotal, totalShippingFee, req.ShippingAddress, req.Items, productInfoMap, shopOrders, voucherTotalDiscount, voucherShippingDiscount)
		// Gọi API CreateTransaction của Transaction Service
		res, err := s.apiServer.CreateTransaction(token, params)
		if err != nil {
			// Ghi log chi tiết lỗi gọi Transaction Service
			log.Printf("Error calling CreateTransaction API: %v, params: %+v", err, params)
			return fmt.Errorf("lỗi khi tạo giao dịch: %w", err)
		}

		if res.Result.Data.PayURL != "" {
			paymentURL = &res.Result.Data.PayURL
		}
		return nil
	})

	if saveErr != nil {
		// Rollback stock reservation nếu lưu thất bại
		_ = s.releaseStockForOrder(ctx, token, reservations)

		if voucherTotalSite != nil {
			// trừ voucher site
			err := s.reserveStockForVoucher(ctx, userID, voucherTotalSite.ID)
			if err != nil {
				return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lưu đơn hàng: %w", saveErr))
			}
		}
		if voucherShippingSite != nil {
			// trừ voucher giao hàng
			err := s.reserveStockForVoucher(ctx, userID, voucherShippingSite.ID)
			if err != nil {
				return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lưu đơn hàng: %w", saveErr))
			}
		}
		// trar loi shopVoucher
		for _, shopOrder := range shopOrders {
			if shopOrder.DiscountCode != "" {
				voucherShop, err := s.repository.GetVoucherForValidation(ctx, shopOrder.DiscountCode)
				if err != nil {
					return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi cập nhật lại voucher: %w", saveErr))
				}
				// trừ voucher site
				errrors := s.reserveStockForVoucher(ctx, userID, voucherShop.ID)
				if errrors != nil {
					return nil, errrors
				}
			}
		}
		return nil, assets_services.NewError(400, fmt.Errorf("lỗi khi lưu đơn hàng: %w", saveErr))
	}

	// Bước 10: Build response
	shopOrderCodes := make([]string, len(shopOrders))
	for i, so := range shopOrders {
		shopOrderCodes[i] = so.ShopOrderCode
	}

	response := map[string]interface{}{
		"orderId":    orderID,
		"orderCode":  orderCode,
		"grandTotal": grandTotal,
		// "status":     in,
		"shopOrders": shopOrderCodes,
	}

	if paymentURL != nil {
		response["paymentUrl"] = *paymentURL
	}

	return response, nil
}

func createInitPaymentParams(
	orderID, PaymentMethod_ID string, grandTotal, totalShippingFee float64, ShippingAddress services.ShippingAddress,
	OrderItemRequest []services.OrderItemRequest, productInfoMap map[string]*ProductInfo,
	shopOrders []ShopOrderWithItems,
	voucherTotalDiscount, voucherShippingDiscount float64) server_transaction.InitPaymentParams {
	// })
	addres := ShippingAddress.Address
	if ShippingAddress.District != nil {
		addres += " - " + *ShippingAddress.District
		if ShippingAddress.City != nil {
			addres += " - " + *ShippingAddress.City

		}
	}

	userInfo := server_transaction.UserInfo{
		Name:        ShippingAddress.FullName, // Lấy tạm từ địa chỉ ship
		PhoneNumber: ShippingAddress.Phone,    // Lấy tạm từ địa chỉ ship
		Address:     *ShippingAddress.City + " " + *ShippingAddress.District + " " + ShippingAddress.Address,
	}
	// Xây dựng danh sách DetailItem cho Transaction Service (nếu cổng TT cần)
	detailItemsForTx := make([]server_transaction.DetailItem, 0, len(OrderItemRequest))
	for _, itemReq := range OrderItemRequest {
		if productInfo, ok := productInfoMap[itemReq.SkuID]; ok {
			// TODO: Xác định giá cuối cùng của sản phẩm sau KM (FinalUnitPrice) đã được tính ở đâu?
			// Giả sử nó nằm trong shopOrder.Items tương ứng hoặc productInfo đã cập nhật
			// Ở đây lấy tạm productInfo.Price (cần đảm bảo đây là giá cuối)
			detailItemsForTx = append(detailItemsForTx, server_transaction.DetailItem{
				ProductID: productInfo.ProductID,
				Name:      productInfo.ProductName,
				ImageURL:  *productInfo.Image, // Kiểm tra nil nếu cần
				Quantity:  itemReq.Quantity,
				Price:     productInfo.Price, // Giả định đây là giá cuối đã giảm
			})
		}
	}
	// Xây dựng danh sách SettlementDetail từ shopOrders đã tính toán
	settlementDetails := make([]server_transaction.SettlementDetail, 0, len(shopOrders))
	// Các biến để tổng hợp chi phí Sàn
	var totalSiteOrderVoucherDiscount float64 = voucherTotalDiscount
	var totalSitePromotionDiscount float64 = 0.0
	var totalSiteShippingDiscount float64 = voucherShippingDiscount
	var totalSiteFundedProductDiscount float64 = voucherTotalDiscount + voucherShippingDiscount

	for _, shopOrder := range shopOrders {
		// TODO: Bổ sung logic tính toán các trường sau trong createShopOrdersWithItems
		// và trả về trong struct ShopOrderWithItems

		// hiện tại không có đợt giảm giá
		var shopFundedProductDiscount float64 = 0 // giá tiền giảm giá của sản phẩm trong đợt giảm giá
		var siteFundedProductDiscount float64 = 0 // tiền sàn trong đợt giảm giá

		var shopVoucherDiscount float64 = 0.0 // mã giảm giá shop
		shopVoucherDiscount = shopOrder.TotalDiscount

		var shopShippingDiscount float64 = 0.0 // Voucher ship shop(hiện tại ko hỗ trợ)

		var orderSubtotal float64 = 0.0    // giá gốc của đơn
		orderSubtotal = shopOrder.Subtotal // Tạm coi subtotal đã tính là giá gốc

		var commissionFee float64 = 0.0     // hoa hồng trên giá gốc
		commissionFee = orderSubtotal * 0.1 // Tạm tính hoa hồng 10%

		var siteOrderDiscount float64 = 0.0                                                                                  // voucher của sàn cho order tổng chia riêng cho shop để giảm giá. Tính với công thức  tổng tiền voucher *( tổng tiền đơn hàng shop / tổng tiền đơn hàng tất cả shop)
		var siteShippingDiscount float64 = 0.0                                                                               // voucher ship sàn giảm cho shop
		siteOrderDiscount = totalSiteOrderVoucherDiscount * ((shopOrder.TotalAmount - shopOrder.TotalDiscount) / grandTotal) // phải lấy giá sản phẩm sau khi giảm giá của shop
		siteShippingDiscount = totalSiteShippingDiscount * (shopOrder.ShippingFee / totalShippingFee)                        // phải lấy giá sản phẩm sau khi giảm giá của shop

		var netSettledAmount float64 = 0.0 // tiền gốc người bán nhận đc
		netSettledAmount = orderSubtotal - shopFundedProductDiscount + siteFundedProductDiscount - shopVoucherDiscount - shopShippingDiscount - commissionFee

		settlementDetails = append(settlementDetails, server_transaction.SettlementDetail{
			ShopOrderID:               shopOrder.ShopOrderID,
			OrderSubtotal:             orderSubtotal,             // Cần giá trị gốc thực tế
			ShopFundedProductDiscount: shopFundedProductDiscount, // Cần giá trị thực tế
			SiteFundedProductDiscount: siteFundedProductDiscount, // Cần giá trị thực tế
			ShopVoucherDiscount:       shopVoucherDiscount,       // Cần giá trị thực tế
			ShippingFee:               shopOrder.ShippingFee,
			SiteOrderDiscount:         siteOrderDiscount, // Cần giá trị thực tếq
			SiteShippingDiscount:      siteShippingDiscount,
			ShopShippingDiscount:      shopShippingDiscount, // Cần giá trị thực tế
			CommissionFee:             commissionFee,        // Cần giá trị thực tế
			NetSettledAmount:          netSettledAmount,     // Cần giá trị thực tế
		})

	}

	// Tạo params cho CreateTransaction
	params := server_transaction.InitPaymentParams{
		OrderID:                        orderID,
		Amount:                         grandTotal, // Tổng tiền khách trả
		PaymentMethodID:                PaymentMethod_ID,
		Items:                          detailItemsForTx,               // Danh sách item cho cổng TT
		SiteOrderVoucherDiscountAmount: totalSiteOrderVoucherDiscount,  // TODO: Cần giá trị thực tế
		SitePromotionDiscountAmount:    totalSitePromotionDiscount,     // TODO: Cần giá trị thực tế
		SiteShippingDiscountAmount:     totalSiteShippingDiscount,      // TODO: Cần giá trị thực tế
		TotalSiteFundedProductDiscount: totalSiteFundedProductDiscount, // TODO: Cần giá trị thực tế
		SettlementDetails:              settlementDetails,              // Chi tiết tài chính từng shop
		UserInfo:                       userInfo,                       // Thông tin người dùng
	}
	return params
}

// Helper: validate request
func (s *service) validateCreateOrderRequest(req services.CreateOrderRequest) *assets_services.ServiceError {
	if len(req.Items) == 0 {
		return assets_services.NewError(400, fmt.Errorf("đơn hàng phải chứa ít nhất một sản phẩm"))
	}

	if req.ShippingAddress.FullName == "" || req.ShippingAddress.Phone == "" || req.ShippingAddress.Address == "" {
		return assets_services.NewError(400, fmt.Errorf("địa chỉ giao hàng không hợp lệ"))
	}

	// if req.PaymentMethod.Type != "ONLINE" && req.PaymentMethod.Type != "OFFLINE" {
	// 	return assets_services.NewError(400, fmt.Errorf("invalid payment method type"))
	// }

	return nil
}

// Helper: lấy thông tin sản phẩm (giả định gọi qua repository hoặc service khác)
func (s *service) fetchProductInfoForOrder(ctx context.Context, items []services.OrderItemRequest) (map[string]*ProductInfo, *assets_services.ServiceError) {
	productMap := make(map[string]*ProductInfo)
	// s.apiServer
	for _, item := range items {
		// Trong production, bạn sẽ gọi Product Service hoặc repository
		// Ở đây tôi giả định lấy từ local DB
		sku, err := s.apiServer.GetSKUs(item.SkuID)
		if err != nil {
			return nil, assets_services.NewError(404, fmt.Errorf("lỗi khi lấy product SKU %s: %w", item.SkuID, err))
		}

		product, err := s.apiServer.GetProductDetail(sku.Result.Data.ProductID)
		if err != nil {
			return nil, assets_services.NewError(404, fmt.Errorf("lỗi khi lấy product %s: %w", sku.Result.Data.ProductID, err))
		}

		price := sku.Result.Data.Price
		quantity := sku.Result.Data.Quantity

		productMap[item.SkuID] = &ProductInfo{
			ProductID:   sku.Result.Data.ProductID,
			SkuID:       sku.Result.Data.ID,
			ProductName: product.Result.Data.Product.Name,
			Image:       &product.Result.Data.Product.Image,
			Price:       price,
			Stock:       quantity,
			Attributes:  sku.Result.Data.SkuName,
		}
	}

	return productMap, nil
}

// Helper: validate stock
func (s *service) validateStock(productMap map[string]*ProductInfo, items []services.OrderItemRequest) *assets_services.ServiceError {
	for _, item := range items {
		product := productMap[item.SkuID]
		if product.Stock < item.Quantity {
			return assets_services.NewError(409, fmt.Errorf(
				"không đủ hàng cho %s. Sẵn có: %d, Yêu cầu: %d",
				product.ProductName, product.Stock, item.Quantity,
			))
		}
	}
	return nil
}

// Helper: nhóm items theo shop
func (s *service) groupItemsByShop(items []services.OrderItemRequest) map[string][]services.OrderItemRequest {
	shopMap := make(map[string][]services.OrderItemRequest)
	for _, item := range items {
		shopMap[item.ShopID] = append(shopMap[item.ShopID], item)
	}
	return shopMap
}

// Helper: tạo shop orders với items
func (s *service) createShopOrdersWithItems(
	ctx context.Context,
	userID string,
	orderID string,
	shopItemsMap map[string][]services.OrderItemRequest,
	productMap map[string]*ProductInfo,
	voucherShop []services.VoucherShopRequest,
	voucherTotalSiteID *string,
	voucherShippingSiteID *string,
) (
	orderShopAndItem []ShopOrderWithItems,
	grandTotal float64, subtotal float64, totalShippingFee float64, totalDiscount float64,
	voucherTotalSite *db.Vouchers, voucherShippingSite *db.Vouchers, voucherTotalDiscount, voucherShippingDiscount float64,
	errors *assets_services.ServiceError) {
	shopOrders := make([]ShopOrderWithItems, 0)
	voucherShopInfo := map[string]*db.Vouchers{}
	// kieểm tra voucher shop có hay không để tạo map trước
	if len(voucherShop) > 0 {
		for _, vs := range voucherShop {
			valid, reason, voucherShopData := s.checkSingleVoucher(ctx, userID, vs.VoucherID)
			if !valid {
				return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher shop không hợp lệ: %s", reason))
			}
			voucherShopInfo[vs.ShopID] = &voucherShopData
		}
	}
	// kiểm tra xem
	for shopID, items := range shopItemsMap {
		shopOrderID := uuid.New().String()
		shopOrderCode := s.generateShopOrderCode(shopID)

		// Tính shipping fee (giả định cố định, trong production sẽ gọi Shipping Service)
		shippingFee := 30000.0

		// Xác định status dựa trên payment method
		status := services.ShopOrderStatusAwaitingPayment

		shopOrder := ShopOrderWithItems{
			ShopOrderID:   shopOrderID,
			ShopOrderCode: shopOrderCode,
			ShopID:        shopID,
			Status:        status,
			ShippingFee:   shippingFee,

			Items: make([]OrderItemData, 0),
		}

		var shopSubtotal float64

		for _, itemReq := range items {
			product := productMap[itemReq.SkuID]
			itemID := uuid.New().String()
			itemTotal := product.Price * float64(itemReq.Quantity)

			item := OrderItemData{
				ItemID:            itemID,
				ProductID:         product.ProductID,
				SkuID:             product.SkuID,
				Quantity:          itemReq.Quantity,
				OriginalUnitPrice: product.Price,
				FinalUnitPrice:    product.Price,
				TotalPrice:        itemTotal,
				ProductName:       product.ProductName,
				ProductImage:      product.Image,
				SkuAttributes:     product.Attributes,
			}

			shopOrder.Items = append(shopOrder.Items, item)
			shopSubtotal += itemTotal
		}

		shopOrder.Subtotal = shopSubtotal // tổng tiền chay 1 đơn hàng
		// check voucher shop
		voucher := voucherShopInfo[shopID]
		if voucher != nil {
			min, err := assets_services.ConvertStringToFloat(voucher.MinPurchaseAmount)
			if err != nil {
				return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi định dạng giá trị tối thiểu của voucher: %w", err))
			}
			if shopSubtotal < min {
				return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher shop không áp dụng cho shop %s: giá trị đơn hàng tối thiểu là %.2f", shopID, min))
			}
			// tính toán discount
			discount := 0.0
			// chỉ áp dụng voucher shop là giảm tiền trên đơn
			if voucher.AppliesToType == db.VouchersAppliesToTypeORDERTOTAL {
				discount, err = countDiscountAmount(*voucher, shopSubtotal)
				if err != nil {
					return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi khi tính toán giảm giá voucher shop cho shop %s: %w", shopID, err))
				}
			}
			shopOrder.TotalDiscount = discount
			shopOrder.DiscountCode = voucher.VoucherCode
		}

		shopOrder.TotalAmount = shopSubtotal + shippingFee - shopOrder.TotalDiscount

		shopOrders = append(shopOrders, shopOrder)
		// tong gia tri don hang
		subtotal += shopSubtotal
		// tong tiem ship taon bo don hang
		totalShippingFee += shippingFee
		// tong tien giam gia ma shop ap dung
		totalDiscount += shopOrder.TotalDiscount
	}

	// tính giảm giá cho toan bo don hang
	if voucherTotalSiteID != nil {
		valid, reason, voucherData := s.checkSingleVoucher(ctx, userID, *voucherTotalSiteID)
		if !valid {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher tổng không hợp lệ: %s", reason))
		}
		// tính toán giảm giá
		min, err := assets_services.ConvertStringToFloat(voucherData.MinPurchaseAmount)
		if err != nil {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi định dạng giá trị tối thiểu của voucher: %w", err))
		}
		if subtotal < min {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher tổng không áp dụng: giá trị đơn hàng tối thiểu là %.2f", min))
		}
		// tính toán discount
		discount := 0.0
		if voucherData.AppliesToType == db.VouchersAppliesToTypeORDERTOTAL {
			discount, err = countDiscountAmount(voucherData, subtotal)
			if err != nil {
				return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi khi tính toán giảm giá voucher tổng: %w", err))
			}
			// TÍNH TOÁN TỔNG TIỀN GIẢM TỔNG CỘNG VÀO GIÁ TRỊ ĐƠN HÀNG
			totalDiscount += discount
			voucherTotalSite = &voucherData
			voucherTotalDiscount = discount
		} else {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("loại voucher tổng không được hỗ trợ"))
		}
	}

	// tính giảm giá cho shop và giảm giá giao hang
	if voucherShippingSiteID != nil {
		valid, reason, voucherData := s.checkSingleVoucher(ctx, userID, *voucherShippingSiteID)
		if !valid {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher tổng không hợp lệ: %s", reason))
		}
		// tính toán giảm giá
		min, err := assets_services.ConvertStringToFloat(voucherData.MinPurchaseAmount)
		if err != nil {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi định dạng giá trị tối thiểu của voucher: %w", err))
		}
		if subtotal < min {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("voucher tổng không áp dụng: giá trị đơn hàng tối thiểu là %.2f", min))
		}
		// tính toán discount
		discount := 0.0
		if voucherData.AppliesToType == db.VouchersAppliesToTypeSHIPPINGFEE {
			discount, err = countDiscountAmount(voucherData, subtotal)
			if err != nil {
				return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(500, fmt.Errorf("lỗi khi tính toán giảm giá voucher giao hàng: %w", err))
			}
			// tính giảm giá cho voucher giao hàng cộng vào tổng giảm giá
			totalDiscount += discount
			voucherShippingDiscount = discount
			voucherShippingSite = &voucherData
		} else {
			return nil, 0, 0, 0, 0, nil, nil, 0, 0, assets_services.NewError(400, fmt.Errorf("loại voucher giao hàng không được hỗ trợ"))
		}
	}

	grandTotal = subtotal + totalShippingFee - totalDiscount
	grandTotal = math.Ceil(grandTotal) // Làm tròn lên

	return shopOrders, grandTotal, subtotal, totalShippingFee, totalDiscount, voucherTotalSite, voucherShippingSite, voucherTotalDiscount, voucherShippingDiscount, nil
}

func countDiscountAmount(voucher db.Vouchers, total float64) (discount float64, err error) {
	if voucher.DiscountType == db.VouchersDiscountTypePERCENTAGE {
		percent, err := assets_services.ConvertStringToFloat(voucher.DiscountValue)
		if err != nil {
			return 0, fmt.Errorf("lỗi định dạng giá trị giảm giá voucher: %w", err)
		}
		if !voucher.MaxDiscountAmount.Valid {
			return 0, fmt.Errorf("thiếu giá trị giảm giá tối đa của voucher")
		}
		maxDiscount, err := assets_services.ConvertStringToFloat(voucher.MaxDiscountAmount.String)
		if err != nil {
			return 0, fmt.Errorf("lỗi định dạng giá trị giảm giá tối đa của voucher: %w", err)
		}
		discount = total * (percent / 100.0)
		if discount > maxDiscount {
			discount = maxDiscount
		}
	} else {
		discountValue, err := assets_services.ConvertStringToFloat(voucher.DiscountValue)
		if err != nil {
			return 0, fmt.Errorf("lỗi định dạng giá trị giảm giá voucher: %w", err)
		}
		discount = discountValue
	}

	return discount, nil
}

// Helper: build stock reservations
func (s *service) buildStockReservations(items []services.OrderItemRequest) []services.ProductUpdateSKUReserver {
	reservations := make([]services.ProductUpdateSKUReserver, len(items))
	for i, item := range items {
		reservations[i] = services.ProductUpdateSKUReserver{
			SkuID:            item.SkuID,
			QuantityReserver: int32(item.Quantity),
		}
	}
	return reservations
}

// Helper: reserve stock
func (s *service) reserveStockForOrder(ctx context.Context, token string, reservations []services.ProductUpdateSKUReserver) *assets_services.ServiceError {
	// copy
	// Map using copier
	var skuParams []server_product.UpdateProductSKUParams
	for _, r := range reservations {
		skuParams = append(skuParams, server_product.UpdateProductSKUParams{
			Sku_ID:           r.SkuID,
			QuantityReserved: int(r.QuantityReserver),
		})
	}
	res, err := s.apiServer.UpdateProductSKU(s.env.TokenSystem, string(services.HOLD), skuParams)
	if err != nil {
		return &assets_services.ServiceError{
			Code: 500,
			Err:  err,
		}
	}
	if res.Code != 200 {
		return &assets_services.ServiceError{
			Code: res.Code,
			Err:  errors.New(res.Message),
		}
	}
	return nil
}

// Helper: release stock
func (s *service) releaseStockForOrder(ctx context.Context, token string, reservations []services.ProductUpdateSKUReserver) *assets_services.ServiceError {
	var skuParams []server_product.UpdateProductSKUParams
	for _, r := range reservations {
		skuParams = append(skuParams, server_product.UpdateProductSKUParams{
			Sku_ID:           r.SkuID,
			QuantityReserved: int(r.QuantityReserver),
		})
	}
	res, err := s.apiServer.UpdateProductSKU(token, string(services.ROLLBACK), skuParams)
	if err != nil {
		return &assets_services.ServiceError{
			Code: 500,
			Err:  err,
		}
	}
	if res.Code != 200 {
		return &assets_services.ServiceError{
			Code: res.Code,
			Err:  errors.New(res.Message),
		}
	}
	return nil
}

// Helper: release stock for voucher
func (s *service) releaseStockForVoucher(ctx context.Context, userID, voucherID string, discountAmount float64) *assets_services.ServiceError {

	err := s.UseVoucher(ctx, services.UseVoucherInput{
		UserID:         userID,
		VoucherID:      voucherID,
		DiscountAmount: discountAmount,
	})
	if err != nil {
		return &assets_services.ServiceError{
			Code: 400,
			Err:  err,
		}
	}
	return nil
}

// Helper: reserve stock for voucher
func (s *service) reserveStockForVoucher(ctx context.Context, userID, voucherID string) *assets_services.ServiceError {

	err := s.RollbackVoucher(ctx, services.RollbackVoucherInput{
		UserID:    userID,
		VoucherID: voucherID,
	})
	if err != nil {
		return &assets_services.ServiceError{
			Code: 400,
			Err:  err,
		}
	}
	return nil
}

// Helper: generate order code
func (s *service) generateOrderCode() string {
	timestamp := time.Now().Format("20060102")
	randomPart := uuid.New().String()[:8]
	return fmt.Sprintf("YAN%s%s", timestamp, randomPart)
}

// Helper: generate shop order code
func (s *service) generateShopOrderCode(shopID string) string {
	randomPart := uuid.New().String()[:8]
	shopIDShort := shopID
	if len(shopID) > 8 {
		shopIDShort = shopID[:8]
	}
	return fmt.Sprintf("SHOP-%s-%s", shopIDShort, randomPart)
}

// Struct helpers
type ProductInfo struct {
	ProductID   string
	SkuID       string
	ProductName string
	Image       *string
	Price       float64
	Stock       int
	Attributes  string
}

type ShopOrderWithItems struct {
	ShopOrderID   string
	ShopOrderCode string
	ShopID        string
	Status        services.ShopOrderStatus
	Subtotal      float64
	ShippingFee   float64
	TotalDiscount float64
	DiscountCode  string
	TotalAmount   float64
	Items         []OrderItemData
}

type OrderItemData struct {
	ItemID            string
	ProductID         string
	SkuID             string
	Quantity          int
	OriginalUnitPrice float64
	FinalUnitPrice    float64
	TotalPrice        float64
	ProductName       string
	ProductImage      *string
	SkuAttributes     string
}

// Helper functions
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

func parseInt(i int32) (int, error) {
	return int(i), nil
}
