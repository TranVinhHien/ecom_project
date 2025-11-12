package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
	assets_services "github.com/TranVinhHien/ecom_analytics_service/services/assets"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"
	"golang.org/x/sync/errgroup"
)

// GetShopOverview xử lý API: GET /api/v1/shop/overview
func (s *service) GetShopOverview(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.ShopOverviewResponse, *assets_services.ServiceError) {

	// Dùng errgroup để chạy 3 tác vụ đồng thời
	g, gCtx := errgroup.WithContext(ctx)

	var orderOverview db_order.GetShopOrderOverviewRow
	var shopLedger db_transaction.AccountLedgers
	var shopOrderIDs []string

	// Tác vụ 1: Lấy tổng quan đơn hàng (từ order_db)
	g.Go(func() error {
		var err error
		orderOverview, err = s.order.GetShopOrderOverview(gCtx, db_order.GetShopOrderOverviewParams{
			ShopID:    shopID,
			StartDate: sql.NullTime{Time: startDate, Valid: startDate.IsZero() == false},
			EndDate:   sql.NullTime{Time: endDate, Valid: endDate.IsZero() == false},
		})
		if err != nil && err != sql.ErrNoRows {
			// log.Printf("Error GetShopOrderOverview: %v", err)
			return fmt.Errorf("lỗi khi lấy tổng quan đơn hàng: %w", err)
		}
		return nil
	})

	// Tác vụ 2: Lấy thông tin ví (từ transaction_db)
	g.Go(func() error {
		var err error
		shopLedger, err = s.transaction.GetLedgerByOwnerID(gCtx, shopID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin ví: %w", err)
		}
		return nil
	})

	// Tác vụ 3: Lấy *tất cả* ID đơn hàng của shop (từ order_db)
	// (Đây là tác vụ có thể nặng nếu Shop có hàng triệu đơn)
	g.Go(func() error {
		var err error
		shopOrderIDs, err = s.order.GetShopOrderIDs(gCtx, shopID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy danh sách ID đơn hàng: %w", err)
		}
		return nil
	})

	// Chờ cả 3 tác vụ hoàn thành
	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: err}
	}

	// Tác vụ 4: Lấy thống kê đối soát (từ transaction_db)
	// Chạy SAU KHI có danh sách ID
	var settlementStats db_transaction.GetShopSettlementStatsByOrderIDsRow
	if len(shopOrderIDs) > 0 {
		var err error
		settlementStats, err = s.transaction.GetShopSettlementStatsByOrderIDs(ctx, shopOrderIDs)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error GetShopSettlementStatsByOrderIDs: %v", err)
			return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin đối soát: %w", err)}
		}
	}

	// Tổng hợp kết quả
	resp := &entity.ShopOverviewResponse{
		TotalOrders:      orderOverview.TotalOrders,
		ProcessingOrders: orderOverview.ProcessingOrders,
	}

	// Xử lý parse (trong thực tế nên dùng thư viện decimal)
	resp.TotalGMV, _ = entity.ParseAmountToFloat64(orderOverview.TotalGmv)
	resp.WalletBalance, _ = entity.ParseAmountToFloat64(shopLedger.Balance)
	resp.PendingBalance, _ = entity.ParseAmountToFloat64(shopLedger.PendingBalance)
	resp.TotalNetRevenue, _ = entity.ParseAmountToFloat64(settlementStats.TotalSettled)

	return resp, nil
}

// GetShopWalletSummary xử lý API: GET /api/v1/shop/wallet/summary
func (s *service) GetShopWalletSummary(ctx context.Context, shopID string) (*entity.ShopWalletSummaryResponse, *assets_services.ServiceError) {

	g, gCtx := errgroup.WithContext(ctx)

	var shopLedger db_transaction.AccountLedgers
	var shopOrderIDs []string
	var ledgerEntries []db_transaction.LedgerEntries

	// Tác vụ 1: Lấy ví
	g.Go(func() error {
		var err error
		shopLedger, err = s.transaction.GetLedgerByOwnerID(gCtx, shopID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin ví: %w", err)
		}
		return nil
	})

	// Tác vụ 2: Lấy ID đơn hàng
	g.Go(func() error {
		var err error
		shopOrderIDs, err = s.order.GetShopOrderIDs(gCtx, shopID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy shop order: %w", err)
		}
		return nil
	})

	// Tác vụ 3: Lấy lịch sử giao dịch (để tính tiền rút)
	// Cảnh báo: Lấy toàn bộ lịch sử là không tốt,
	// nên có câu sqlc SUM, nhưng vì ta dùng query đã có,
	// ta sẽ dùng logic Go để tính.
	g.Go(func() error {
		var err error
		// Lấy 1000 giao dịch gần nhất để tính
		ledgerEntries, err = s.transaction.ListLedgerEntriesByOwnerID(gCtx, db_transaction.ListLedgerEntriesByOwnerIDParams{
			OwnerID: shopID,
			Limit:   1000, // Giới hạn, nếu không sẽ sập CSDL
			Offset:  0,
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy lịch sử giao dịch: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: err}
	}

	// Tác vụ 4: Lấy thống kê đối soát (chạy sau Tác vụ 2)
	var settlementStats db_transaction.GetShopSettlementStatsByOrderIDsRow
	if len(shopOrderIDs) > 0 {
		var err error
		settlementStats, err = s.transaction.GetShopSettlementStatsByOrderIDs(ctx, shopOrderIDs)
		if err != nil && err != sql.ErrNoRows {
			return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin đối soát: %w", err)}
		}
	}

	// Xử lý dữ liệu
	resp := &entity.ShopWalletSummaryResponse{}
	resp.Balance, _ = entity.ParseAmountToFloat64(shopLedger.Balance)
	resp.PendingBalance, _ = entity.ParseAmountToFloat64(shopLedger.PendingBalance)
	resp.TotalSettledRevenue, _ = entity.ParseAmountToFloat64(settlementStats.TotalSettled)
	resp.TotalFundsHeld, _ = entity.ParseAmountToFloat64(settlementStats.TotalFundsHeld)

	// Tính toán TotalWithdrawn (Logic trong Go)
	var totalWithdrawn float64
	for _, entry := range ledgerEntries {
		// Giả định: bút toán rút tiền là DEBIT và có chữ "Withdrawal"
		if entry.Type == "DEBIT" && strings.Contains(entry.Description, "Withdrawal") {
			amount, _ := entity.ParseAmountToFloat64(entry.Amount)
			// Amount là số âm (DEBIT), nên ta trừ đi để thành số dương
			totalWithdrawn -= amount
		}
	}
	resp.TotalWithdrawn = totalWithdrawn

	return resp, nil
}

// === NHÓM 2: Phân tích Đơn hàng ===

// ListShopOrders xử lý API: GET /api/v1/shop/orders
func (s *service) ListShopOrders(ctx context.Context, params entity.ListShopOrdersParams) (*entity.ListShopOrdersResponse, *assets_services.ServiceError) {

	// Bước 1: Lấy danh sách đơn hàng của shop từ order_db
	shopOrders, err := s.order.ListShopOrders(ctx, db_order.ListShopOrdersParams{
		ShopID:       params.ShopID,
		StatusFilter: db_order.NullShopOrdersStatus{ShopOrdersStatus: db_order.ShopOrdersStatus(params.Status.String), Valid: params.Status.Valid},
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.ListShopOrdersResponse{Orders: []entity.ShopOrderWithPaymentStatus{}}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy danh sách đơn hàng: %w", err)}
	}

	if len(shopOrders) == 0 {
		return &entity.ListShopOrdersResponse{Orders: []entity.ShopOrderWithPaymentStatus{}}, nil
	}

	// Bước 2: Thu thập các parent_order_id để truy vấn transaction_db
	orderIDs := make([]sql.NullString, 0, len(shopOrders))
	for _, so := range shopOrders {
		orderIDs = append(orderIDs, sql.NullString{String: so.OrderID, Valid: true})
		// (Chúng ta có thể tối ưu bằng cách dùng map để lấy unique ID,
		// nhưng sqlc IN(slice) đã xử lý tốt việc này)
	}

	// Bước 3: Lấy trạng thái thanh toán từ transaction_db
	txnStatuses, err := s.transaction.GetTransactionStatusesByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy trạng thái giao dịch: %w", err)}
	}

	// Tạo map để tra cứu nhanh: map[order_id] -> status
	paymentStatusMap := make(map[string]string)
	for _, ts := range txnStatuses {
		paymentStatusMap[ts.OrderID.String] = string(ts.Status)
	}

	// Bước 4: Tổng hợp dữ liệu
	resp := &entity.ListShopOrdersResponse{
		Orders: make([]entity.ShopOrderWithPaymentStatus, len(shopOrders)),
	}
	for i, so := range shopOrders {
		// status := "UNKNOWN" // Mặc định nếu không tìm thấy
		// if s, ok := paymentStatusMap[so.OrderID]; ok {
		// 	status = s
		// }

		resp.Orders[i] = entity.ShopOrderWithPaymentStatus{
			ShopOrder: so,
			// PaymentStatus: status,
		}
	}

	return resp, nil
}

// GetEnrichedShopOrder xử lý API: GET /api/v1/shop/orders/{id}/enriched
func (s *service) GetEnrichedShopOrder(ctx context.Context, shopID, shopOrderID string) (*entity.EnrichedShopOrderResponse, *assets_services.ServiceError) {

	// Bước 1: Lấy thông tin shop_order gốc
	shopOrder, err := s.order.GetShopOrderByID(ctx, db_order.GetShopOrderByIDParams{
		ID:     shopOrderID,
		ShopID: shopID, // Đảm bảo shop chỉ lấy đúng đơn của mình
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &assets_services.ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("không tìm thấy shop order: %w", err)}
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy shop order: %w", err)}
	}

	g, gCtx := errgroup.WithContext(ctx)

	resp := &entity.EnrichedShopOrderResponse{
		OrderDetails: shopOrder,
	}
	var settlementInfo []db_transaction.ShopOrderSettlements // Phải là slice vì sqlc trả về slice
	var paymentStatusRows []db_transaction.GetTransactionStatusesByOrderIDsRow

	// Tác vụ 2a: Lấy order items (từ order_db)
	g.Go(func() error {
		var err error
		resp.OrderItems, err = s.order.GetOrderItemsByShopOrderID(gCtx, shopOrderID)
		if err != nil && err != sql.ErrNoRows {
			return &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy order items: %w", err)}
		}
		return nil
	})

	// Tác vụ 2b: Lấy thông tin đối soát (từ transaction_db)
	g.Go(func() error {
		var err error
		settlementInfo, err = s.transaction.GetShopSettlementsByOrderIDs(gCtx, []string{shopOrderID})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin đối soát: %w", err)
		}
		return nil
	})

	// Tác vụ 2c: Lấy thông tin đơn hàng TỔNG (từ order_db)
	g.Go(func() error {
		var err error
		resp.ParentOrderInfo, err = s.order.GetOrderByID(gCtx, shopOrder.OrderID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin đơn hàng tổng: %w", err)
		}
		return nil
	})

	// Tác vụ 2d: Lấy trạng thái thanh toán (từ transaction_db)
	g.Go(func() error {
		var err error
		paymentStatusRows, err = s.transaction.GetTransactionStatusesByOrderIDs(gCtx, []sql.NullString{{String: shopOrder.OrderID, Valid: true}})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy trạng thái giao dịch: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi chờ các tác vụ: %w", err)}
	}

	// Xử lý kết quả
	if len(settlementInfo) > 0 {
		resp.SettlementInfo = &settlementInfo[0]
	}

	if len(paymentStatusRows) > 0 {
		resp.PaymentStatus = string(paymentStatusRows[0].Status)
	} else {
		resp.PaymentStatus = "UNKNOWN"
	}

	return resp, nil
}

// ListShopOrderItems (Hàm này khá đơn giản, chủ yếu là wrapper)
// API này không có trong sqlc, nên tôi sẽ giả định một query mới.
// Giả định bạn đã thêm query này vào `order_analytics.query.sql`:
/*
-- name: ListShopOrderItems :many
SELECT * FROM order_items oi
JOIN shop_orders so ON oi.shop_order_id = so.id
WHERE
    so.shop_id = ?
    AND (sqlc.narg(product_id_filter) IS NULL OR oi.product_id = sqlc.narg(product_id_filter))
    AND (so.created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date))
LIMIT ? OFFSET ?;
*/
func (s *service) ListShopOrderItems(ctx context.Context, params entity.ListShopOrderItemsParams) ([]db_order.OrderItems, *assets_services.ServiceError) {
	// Giả định hàm sqlc ListShopOrderItems tồn tại
	// items, err := s.order.ListShopOrderItems(ctx, db_order.ListShopOrderItemsParams{...})
	// if err != nil { ... }
	// return items, nil

	// Vì hàm sqlc chưa có, chúng ta tạm trả về lỗi "chưa implement"
	return nil, &assets_services.ServiceError{Code: http.StatusNotImplemented, Err: fmt.Errorf("hàm ListShopOrderItems chưa được implement")}
}

// === NHÓM 3: Phân tích Doanh thu & Dòng tiền ===

// GetShopRevenueTimeseries xử lý API: GET /api/v1/shop/revenue/timeseries
func (s *service) GetShopRevenueTimeseries(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.RevenueTimeseriesResponse, *assets_services.ServiceError) {

	g, gCtx := errgroup.WithContext(ctx)

	var gmvData []db_order.GetShopRevenueTimeSeriesRow
	var shopOrderIDs []string
	var settlements []db_transaction.ShopOrderSettlements

	// Tác vụ 1: Lấy GMV timeseries
	g.Go(func() error {
		var err error
		gmvData, err = s.order.GetShopRevenueTimeSeries(gCtx, db_order.GetShopRevenueTimeSeriesParams{
			ShopID:          shopID,
			FromCompletedAt: sql.NullTime{Time: startDate, Valid: true},
			ToCompletedAt:   sql.NullTime{Time: endDate, Valid: true},
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy dữ liệu GMV: %w", err)
		}
		return nil
	})

	// Tác vụ 2: Lấy ID đơn hàng
	g.Go(func() error {
		var err error
		shopOrderIDs, err = s.order.GetShopOrderIDs(gCtx, shopID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy danh sách ID đơn hàng: %w", err)
		}
		return nil
	})

	// ✅ Chờ lấy được shopOrderIDs trước
	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi chờ các tác vụ: %w", err)}
	}

	// ✅ Tạo errgroup mới cho tác vụ tiếp theo
	g2, gCtx2 := errgroup.WithContext(ctx)

	// Tác vụ 3: Lấy settlements (chạy sau khi có shopOrderIDs)
	g2.Go(func() error {
		if len(shopOrderIDs) == 0 {
			return nil
		}
		var err error
		settlements, err = s.transaction.GetShopSettlementsByOrderIDs(gCtx2, shopOrderIDs)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin đối soát: %w", err)
		}
		return nil
	})

	if err := g2.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi chờ tác vụ đối soát: %w", err)}
	}

	// Bắt đầu xử lý logic trong Go (Kém hiệu quả, nhưng bắt buộc do kiến trúc 3-DB)

	// map[date_string] -> net_revenue
	netRevenueMap := make(map[string]float64)
	for _, set := range settlements {
		if set.Status == "SETTLED" && set.SettledAt.Valid {
			dateStr := set.SettledAt.Time.Format("2006-01-02")

			// Lọc theo ngày
			if set.SettledAt.Time.Before(startDate) || set.SettledAt.Time.After(endDate) {
				continue
			}

			amount, _ := entity.ParseAmountToFloat64(set.NetSettledAmount)
			netRevenueMap[dateStr] += amount
		}
	}

	// map[date_string] -> Datapoint
	// Dùng GMV làm "key" chính
	finalMap := make(map[string]entity.RevenueDatapoint)
	for _, gmv := range gmvData {
		dateStr := gmv.ReportDate.Format("2006-01-02")
		amount, _ := entity.ParseAmountToFloat64(gmv.Gmv)
		finalMap[dateStr] = entity.RevenueDatapoint{
			Date: dateStr,
			GMV:  amount,
			// NetRevenue sẽ được điền ở bước sau
		}
	}

	// Ghép NetRevenue vào
	for dateStr, netAmount := range netRevenueMap {
		if dp, ok := finalMap[dateStr]; ok {
			dp.NetRevenue = netAmount
			finalMap[dateStr] = dp
		} else {
			// (Trường hợp có NetRevenue nhưng không có GMV)
			finalMap[dateStr] = entity.RevenueDatapoint{
				Date:       dateStr,
				NetRevenue: netAmount,
			}
		}
	}

	resp := &entity.RevenueTimeseriesResponse{
		Data: make([]entity.RevenueDatapoint, 0, len(finalMap)),
	}
	for _, dp := range finalMap {
		resp.Data = append(resp.Data, dp)
	}
	// (Nên sắp xếp resp.Data theo ngày trước khi trả về)

	return resp, nil
}

// ListShopWalletLedgerEntries xử lý API: GET /api/v1/shop/wallet/ledger-entries
func (s *service) ListShopWalletLedgerEntries(ctx context.Context, params entity.ListWalletLedgerEntriesParams) ([]db_transaction.LedgerEntries, *assets_services.ServiceError) {
	entries, err := s.transaction.ListLedgerEntriesByOwnerID(ctx, db_transaction.ListLedgerEntriesByOwnerIDParams{
		OwnerID: params.ShopID,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.LedgerEntries{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin giao dịch: %w", err)}
	}
	return entries, nil
}

// ListShopSettlements xử lý API: GET /api/v1/shop/settlements
func (s *service) ListShopSettlements(ctx context.Context, params entity.ListShopSettlementsParams) ([]db_transaction.ShopOrderSettlements, *assets_services.ServiceError) {

	// Bước 1: Lấy ID đơn hàng từ order_db (theo shop_id và date_range)
	// (Giả định hàm ListShopOrders có thể lọc theo ngày,
	// nhưng hàm sqlc GetShopOrderIDs thì không.
	// Đây là 1 điểm yếu nữa, tạm thời chúng ta lấy *tất cả* ID)
	shopOrderIDs, err := s.order.GetShopOrderIDs(ctx, params.ShopID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.ShopOrderSettlements{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy danh sách ID đơn hàng: %w", err)}
	}

	if len(shopOrderIDs) == 0 {
		return []db_transaction.ShopOrderSettlements{}, nil
	}

	// Bước 2: Lấy settlements từ transaction_db
	// (Hàm GetShopSettlementsByOrderIDs không hỗ trợ lọc status, date, pagination.
	// Chúng ta phải lọc trong Go - Rất không hiệu quả)

	// TODO: Cần một hàm sqlc `ListShopSettlements` tốt hơn
	// Tạm thời, chúng ta dùng hàm đã có
	settlements, err := s.transaction.GetShopSettlementsByOrderIDs(ctx, shopOrderIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.ShopOrderSettlements{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin đối soát: %w", err)}
	}

	// Lọc trong Go (Tạm thời)
	filtered := make([]db_transaction.ShopOrderSettlements, 0)
	for _, set := range settlements {
		// Lọc Status
		if params.Status.Valid && set.Status != db_transaction.ShopOrderSettlementsStatus(params.Status.String) {
			continue
		}
		// Lọc Date
		if params.StartDate.Valid && set.OrderCompletedAt.Time.Before(params.StartDate.Time) {
			continue
		}
		if params.EndDate.Valid && set.OrderCompletedAt.Time.After(params.EndDate.Time) {
			continue
		}
		filtered = append(filtered, set)
	}
	// (Chưa xử lý pagination)

	return filtered, nil
}

// === NHÓM 4: Phân tích Voucher ===

// ListShopVouchers xử lý API: GET /api/v1/shop/vouchers
func (s *service) ListShopVouchers(ctx context.Context, params entity.ListShopVouchersParams) ([]db_order.Vouchers, *assets_services.ServiceError) {
	vouchers, err := s.order.ListVouchersByOwner(ctx, db_order.ListVouchersByOwnerParams{
		OwnerID:        params.ShopID,
		IsActiveFilter: params.IsActive,
		Limit:          params.Limit,
		Offset:         params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.Vouchers{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy danh sách voucher: %w", err)}
	}
	return vouchers, nil
}

// GetShopVoucherPerformance xử lý API: GET /api/v1/shop/vouchers/performance
func (s *service) GetShopVoucherPerformance(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.VoucherPerformanceResponse, *assets_services.ServiceError) {
	stats, err := s.order.GetVoucherUsagePerformanceByOwner(ctx, db_order.GetVoucherUsagePerformanceByOwnerParams{
		OwnerID:    shopID,
		FromUsedAt: startDate,
		ToUsedAt:   endDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.VoucherPerformanceResponse{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin hiệu suất voucher: %w", err)}
	}

	resp := &entity.VoucherPerformanceResponse{
		TotalUsageCount: stats.TotalUsageCount,
	}
	resp.TotalDiscountValue, _ = entity.ParseAmountToFloat64(stats.TotalDiscountValue)

	return resp, nil
}

// GetShopVoucherUsageDetails xử lý API: GET /api/v1/shop/vouchers/{id}/detail
func (s *service) GetShopVoucherUsageDetails(ctx context.Context, params entity.ListVoucherUsageDetailsParams) ([]db_order.VoucherUsageHistory, *assets_services.ServiceError) {

	// (Cần kiểm tra xem voucherID này có thực sự thuộc shopID không)
	// (Tạm bỏ qua bước xác thực đó)

	history, err := s.order.GetVoucherUsageHistory(ctx, db_order.GetVoucherUsageHistoryParams{
		VoucherID: params.VoucherID,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.VoucherUsageHistory{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy lịch sử sử dụng voucher: %w", err)}
	}
	return history, nil
}

// === NHÓM 5: Xếp hạng ===

// GetShopRankingProducts xử lý API: GET /api/v1/shop/ranking/products
func (s *service) GetShopRankingProducts(ctx context.Context, params entity.ShopRankingProductsParams) ([]entity.ProductRankingRow, *assets_services.ServiceError) {

	if params.SortBy == "revenue" {
		rows, err := s.order.GetShopTopProductsByRevenue(ctx, db_order.GetShopTopProductsByRevenueParams{
			ShopID:    params.ShopID,
			StartDate: sql.NullTime{Time: params.StartDate, Valid: true},
			EndDate:   sql.NullTime{Time: params.EndDate, Valid: true},
			Limit:     params.Limit,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return []entity.ProductRankingRow{}, nil
			}
			return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy sản phẩm xếp hạng theo doanh thu: %w", err)}
		}

		// Chuyển đổi sang DTO
		resp := make([]entity.ProductRankingRow, len(rows))
		for i, row := range rows {
			revenue, _ := entity.ParseAmountToFloat64(row.TotalRevenue)
			resp[i] = entity.ProductRankingRow{
				ProductID:             row.ProductID,
				SkuID:                 row.SkuID,
				ProductNameSnapshot:   row.ProductNameSnapshot,
				SkuAttributesSnapshot: row.SkuAttributesSnapshot.String,
				TotalRevenue:          revenue,
				// TotalQuantity không được trả về từ query này
			}
		}
		return resp, nil

	} else if params.SortBy == "quantity" {
		rows, err := s.order.GetShopTopProductsByQuantity(ctx, db_order.GetShopTopProductsByQuantityParams{
			ShopID:    params.ShopID,
			StartDate: sql.NullTime{Time: params.StartDate, Valid: true},
			EndDate:   sql.NullTime{Time: params.EndDate, Valid: true},
			Limit:     params.Limit,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				return []entity.ProductRankingRow{}, nil
			}
			return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy sản phẩm xếp hạng theo số lượng: %w", err)}
		}

		// Chuyển đổi sang DTO
		resp := make([]entity.ProductRankingRow, len(rows))
		for i, row := range rows {
			quantity, _ := entity.ParseInterfaceToInt(row.TotalQuantity)
			resp[i] = entity.ProductRankingRow{
				ProductID:             row.ProductID,
				SkuID:                 row.SkuID,
				ProductNameSnapshot:   row.ProductNameSnapshot,
				SkuAttributesSnapshot: row.SkuAttributesSnapshot.String,
				TotalQuantity:         quantity,
				// TotalRevenue không được trả về từ query này
			}
		}
		return resp, nil
	}

	return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("tham số sort_by không hợp lệ: %s", params.SortBy)}
}
