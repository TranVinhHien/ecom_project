package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
	assets_services "github.com/TranVinhHien/ecom_analytics_service/services/assets"
	entity "github.com/TranVinhHien/ecom_analytics_service/services/entity"

	"golang.org/x/sync/errgroup"
)

// === NHÓM 1: Tổng quan ===

// GetPlatformOverview xử lý API: GET /api/v1/platform/overview
func (s *service) GetPlatformOverview(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformOverviewResponse, *assets_services.ServiceError) {

	g, gCtx := errgroup.WithContext(ctx)

	var orderOverview db_order.GetPlatformOrderOverviewRow
	var revenueSummary db_transaction.GetPlatformRevenueSummaryRow
	var costSummary db_transaction.GetPlatformCostSummaryRow

	// Tác vụ 1: Lấy tổng quan GMV, Order, User, Shop (từ order_db)
	g.Go(func() error {
		var err error
		orderOverview, err = s.order.GetPlatformOrderOverview(gCtx, db_order.GetPlatformOrderOverviewParams{
			FromCompletedAt: sql.NullTime{Time: startDate, Valid: true},
			ToCompletedAt:   sql.NullTime{Time: endDate, Valid: true},
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy tổng quan đơn hàng: %w", err)
		}
		return nil
	})

	// Tác vụ 2: Lấy Doanh thu Sàn (từ transaction_db)
	g.Go(func() error {
		var err error
		revenueSummary, err = s.transaction.GetPlatformRevenueSummary(gCtx, db_transaction.GetPlatformRevenueSummaryParams{
			FromSettledAt: sql.NullTime{Time: startDate, Valid: true},
			ToSettledAt:   sql.NullTime{Time: endDate, Valid: true},
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy doanh thu sàn: %w", err)
		}
		return nil
	})

	// Tác vụ 3: Lấy Chi phí Sàn (từ transaction_db)
	g.Go(func() error {
		var err error
		costSummary, err = s.transaction.GetPlatformCostSummary(gCtx, db_transaction.GetPlatformCostSummaryParams{
			FromCreatedAt: startDate,
			ToCreatedAt:   endDate,
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy chi phí sàn: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("lỗi khi chờ các tác vụ: %w", err)}
	}

	// Tổng hợp dữ liệu
	resp := &entity.PlatformOverviewResponse{
		TotalOrders: orderOverview.TotalOrders,
		TotalShops:  orderOverview.TotalShops,
	}

	// (Parse từ sql.NullString sang float64. Cần hàm helper 'entity.ParseAmountToFloat64' đã viết)
	totalGMV, _ := entity.ParseAmountToFloat64(orderOverview.TotalGmv)
	totalCommission, _ := entity.ParseAmountToFloat64(revenueSummary.TotalCommission)
	totalShippingRevenue, _ := entity.ParseAmountToFloat64(revenueSummary.TotalShippingRevenue)
	totalOrderVoucherCost, _ := entity.ParseAmountToFloat64(costSummary.TotalOrderVoucherCost)
	totalPromoCost, _ := entity.ParseAmountToFloat64(costSummary.TotalPromotionCost)
	totalShippingDiscount, _ := entity.ParseAmountToFloat64(costSummary.TotalShippingDiscountCost)
	totalSubsidy, _ := entity.ParseAmountToFloat64(costSummary.TotalProductSubsidyCost)

	resp.TotalGMV = totalGMV
	resp.TotalPlatformRevenue = totalCommission + totalShippingRevenue
	resp.TotalPlatformCost = totalOrderVoucherCost + totalPromoCost + totalShippingDiscount + totalSubsidy
	resp.PlatformProfit = resp.TotalPlatformRevenue - resp.TotalPlatformCost

	return resp, nil
}

// === NHÓM 2: Quản lý Đơn hàng ===

// ListPlatformOrders xử lý API: GET /api/v1/platform/orders
func (s *service) ListPlatformOrders(ctx context.Context, params entity.ListPlatformOrdersParams) (*entity.ListPlatformOrdersResponse, *assets_services.ServiceError) {

	// Bước 1: Lấy danh sách đơn hàng từ order_db
	shopOrders, err := s.order.ListPlatformOrders(ctx, db_order.ListPlatformOrdersParams{
		ShopIDFilter: params.ShopID,
		StatusFilter: db_order.NullShopOrdersStatus{ShopOrdersStatus: db_order.ShopOrdersStatus(params.Status.String), Valid: params.Status.Valid},
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.ListPlatformOrdersResponse{Orders: []entity.ShopOrderWithPaymentStatus{}}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy danh sách đơn hàng: %w", err)}
	}

	if len(shopOrders) == 0 {
		return &entity.ListPlatformOrdersResponse{Orders: []entity.ShopOrderWithPaymentStatus{}}, nil
	}

	// Bước 2: Thu thập order_id để lấy trạng thái thanh toán
	orderIDs := make([]sql.NullString, 0, len(shopOrders))
	for _, so := range shopOrders {
		orderIDs = append(orderIDs, sql.NullString{String: so.OrderID, Valid: true})
	}

	// Bước 3: Lấy trạng thái thanh toán từ transaction_db
	txnStatuses, err := s.transaction.GetTransactionStatusesByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy trạng thái giao dịch: %w", err)}
	}

	// Tạo map tra cứu
	paymentStatusMap := make(map[string]string)
	for _, ts := range txnStatuses {
		paymentStatusMap[ts.OrderID.String] = string(ts.Status)
	}

	// Bước 4: Tổng hợp dữ liệu
	resp := &entity.ListPlatformOrdersResponse{
		Orders: make([]entity.ShopOrderWithPaymentStatus, len(shopOrders)),
	}
	for i, so := range shopOrders {
		// status := "UNKNOWN"
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

// GetEnrichedPlatformOrder xử lý API: GET /api/v1/platform/orders/{order_id}/detail
func (s *service) GetEnrichedPlatformOrder(ctx context.Context, orderID string) (*entity.EnrichedPlatformOrderResponse, *assets_services.ServiceError) {

	g, gCtx := errgroup.WithContext(ctx)
	resp := &entity.EnrichedPlatformOrderResponse{}

	// Tác vụ 1: Lấy đơn hàng TỔNG
	g.Go(func() error {
		var err error
		resp.ParentOrder, err = s.order.GetOrderByID(gCtx, orderID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin đơn hàng: %w", err)
		}
		return nil
	})

	// Tác vụ 2: Lấy các đơn hàng SHOP con
	g.Go(func() error {
		var err error
		resp.ShopOrders, err = s.order.GetShopOrdersByOrderID(gCtx, orderID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy các đơn hàng shop: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		// Nếu 1 trong 2 lỗi (ví dụ: không tìm thấy)
		if err == sql.ErrNoRows {
			return nil, &assets_services.ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("đơn hàng không tồn tại")}
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: err}
	}

	return resp, nil
}

// === NHÓM 3: Quản lý Tài chính ===

// GetPlatformRevenueTimeseries xử lý API: GET /api/v1/platform/finance/revenue-timeseries
func (s *service) GetPlatformRevenueTimeseries(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformRevenueTimeseriesResponse, *assets_services.ServiceError) {

	// Phân tích: Chúng ta cần 3 dòng dữ liệu theo thời gian:
	// 1. GMV (Từ order_db)
	// 2. Doanh thu Sàn (Từ transaction_db, bảng settlements)
	// 3. Chi phí Sàn (Từ transaction_db, bảng platform_costs)

	g, gCtx := errgroup.WithContext(ctx)

	var gmvData []db_order.GetPlatformGMVTimeSeriesRow // (Sử dụng lại struct của Shop)
	var revenueData []db_transaction.GetPlatformRevenueTimeSeriesRow
	var costData []db_transaction.GetPlatformCostTimeSeriesRow

	// Tác vụ 1: Lấy GMV
	g.Go(func() error {
		var err error
		gmvData, err = s.order.GetPlatformGMVTimeSeries(gCtx, db_order.GetPlatformGMVTimeSeriesParams{
			StartDate: sql.NullTime{Time: startDate, Valid: true},
			EndDate:   sql.NullTime{Time: endDate, Valid: true},
		})
		if err != nil {
			return err
		}
		// Tạm thời bỏ qua Tác vụ 1 vì sqlc query chưa tồn tại
		return nil
	})

	// Tác vụ 2: Lấy Doanh thu Sàn
	g.Go(func() error {
		var err error
		revenueData, err = s.transaction.GetPlatformRevenueTimeSeries(gCtx, db_transaction.GetPlatformRevenueTimeSeriesParams{
			FromSettledAt: sql.NullTime{Time: startDate, Valid: true},
			ToSettledAt:   sql.NullTime{Time: endDate, Valid: true},
		})
		if err != nil {
			return err
		}
		return nil
	})

	// Tác vụ 3: Lấy Chi phí Sàn
	g.Go(func() error {
		var err error
		costData, err = s.transaction.GetPlatformCostTimeSeries(gCtx, db_transaction.GetPlatformCostTimeSeriesParams{
			FromCreatedAt: startDate,
			ToCreatedAt:   endDate,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy dữ liệu theo thời gian: %w", err)}
	}

	// Tổng hợp 3 luồng dữ liệu (Logic Join trong Go)
	dataMap := make(map[string]*entity.PlatformRevenueDatapoint)

	// (Nếu gmvData có dữ liệu)
	for _, row := range gmvData {
		dateStr := row.ReportDate.Format("2006-01-02")
		gmv, _ := entity.ParseAmountToFloat64(row.Gmv)
		dataMap[dateStr] = &entity.PlatformRevenueDatapoint{
			Date:     dateStr,
			TotalGMV: gmv,
		}
	}

	for _, row := range revenueData {
		dateStr := row.ReportDate.Format("2006-01-02")
		revenue, _ := entity.ParseAmountToFloat64(row.PlatformRevenue)
		if dp, ok := dataMap[dateStr]; ok {
			dp.PlatformRevenue = revenue
		} else {
			dataMap[dateStr] = &entity.PlatformRevenueDatapoint{Date: dateStr, PlatformRevenue: revenue}
		}
	}

	for _, row := range costData {
		dateStr := row.ReportDate.Format("2006-01-02")
		cost, _ := entity.ParseAmountToFloat64(row.TotalCost)
		if dp, ok := dataMap[dateStr]; ok {
			dp.PlatformCost = cost
		} else {
			dataMap[dateStr] = &entity.PlatformRevenueDatapoint{Date: dateStr, PlatformCost: cost}
		}
	}

	// Tính lợi nhuận và chuyển sang slice
	resp := &entity.PlatformRevenueTimeseriesResponse{
		Data: make([]entity.PlatformRevenueDatapoint, 0, len(dataMap)),
	}
	for _, dp := range dataMap {
		dp.PlatformProfit = dp.PlatformRevenue - dp.PlatformCost
		resp.Data = append(resp.Data, *dp)
	}
	// (Nên sắp xếp resp.Data theo ngày)

	return resp, nil
}

// ListPlatformTransactions (Wrapper)
func (s *service) ListPlatformTransactions(ctx context.Context, params entity.ListPlatformTransactionsParams) ([]db_transaction.Transactions, *assets_services.ServiceError) {
	txns, err := s.transaction.ListPlatformTransactions(ctx, db_transaction.ListPlatformTransactionsParams{
		TypeFilter:   db_transaction.NullTransactionsType{Valid: params.Type.Valid, TransactionsType: db_transaction.TransactionsType(params.Type.String)},
		StatusFilter: db_transaction.NullTransactionsStatus{Valid: params.Status.Valid, TransactionsStatus: db_transaction.TransactionsStatus(params.Status.String)},
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.Transactions{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin giao dịch: %w", err)}
	}
	return txns, nil
}

// ListPlatformSettlements (Wrapper)
func (s *service) ListPlatformSettlements(ctx context.Context, params entity.ListPlatformSettlementsParams) ([]db_transaction.ShopOrderSettlements, *assets_services.ServiceError) {
	settlements, err := s.transaction.ListPlatformSettlements(ctx, db_transaction.ListPlatformSettlementsParams{
		StatusFilter: db_transaction.NullShopOrderSettlementsStatus{Valid: params.Status.Valid, ShopOrderSettlementsStatus: db_transaction.ShopOrderSettlementsStatus(params.Status.String)},
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.ShopOrderSettlements{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin thanh toán shop: %w", err)}
	}
	return settlements, nil
}

// ListPlatformLedgers (Wrapper)
func (s *service) ListPlatformLedgers(ctx context.Context, params entity.ListPlatformLedgersParams) ([]db_transaction.AccountLedgers, *assets_services.ServiceError) {
	ledgers, err := s.transaction.ListPlatformLedgers(ctx, db_transaction.ListPlatformLedgersParams{
		OwnerTypeFilter: db_transaction.NullAccountLedgersOwnerType{Valid: params.OwnerType.Valid, AccountLedgersOwnerType: db_transaction.AccountLedgersOwnerType(params.OwnerType.String)},
		Limit:           params.Limit,
		Offset:          params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.AccountLedgers{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin sổ cái: %w", err)}
	}
	return ledgers, nil
}

// ListLedgerEntries (Wrapper)
func (s *service) ListLedgerEntries(ctx context.Context, ledgerID string, limit, offset int32) ([]db_transaction.LedgerEntries, *assets_services.ServiceError) {
	entries, err := s.transaction.ListLedgerEntriesByLedgerID(ctx, db_transaction.ListLedgerEntriesByLedgerIDParams{
		LedgerID: ledgerID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_transaction.LedgerEntries{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy các mục sổ cái: %w", err)}
	}
	return entries, nil
}

// === NHÓM 4: Phân tích Voucher ===

// ListPlatformVouchers (Wrapper)
func (s *service) ListPlatformVouchers(ctx context.Context, params entity.ListPlatformVouchersParams) ([]db_order.Vouchers, *assets_services.ServiceError) {
	vouchers, err := s.order.ListPlatformVouchers(ctx, db_order.ListPlatformVouchersParams{
		OwnerTypeFilter: db_order.NullVouchersOwnerType{Valid: params.OwnerType.Valid, VouchersOwnerType: db_order.VouchersOwnerType(params.OwnerType.String)},
		IsActiveFilter:  params.IsActive,
		Limit:           params.Limit,
		Offset:          params.Offset,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.Vouchers{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy danh sách voucher: %w", err)}
	}
	return vouchers, nil
}

// GetPlatformVoucherPerformance xử lý API: GET /api/v1/platform/vouchers/performance/platform
func (s *service) GetPlatformVoucherPerformance(ctx context.Context, startDate, endDate time.Time) (*entity.PlatformVoucherPerformanceResponse, *assets_services.ServiceError) {

	g, gCtx := errgroup.WithContext(ctx)

	var usageStats db_order.GetPlatformVoucherPerformanceRow
	var costStats db_transaction.GetPlatformCostSummaryRow

	// Tác vụ 1: Lấy từ voucher_db (lịch sử sử dụng)
	g.Go(func() error {
		var err error
		usageStats, err = s.order.GetPlatformVoucherPerformance(gCtx, db_order.GetPlatformVoucherPerformanceParams{
			FromUsedAt: startDate,
			ToUsedAt:   endDate,
		})
		if err != nil && err != sql.ErrNoRows {
			return &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin hiệu suất voucher: %w", err)}
		}
		return nil
	})

	// Tác vụ 2: Lấy từ transaction_db (chi phí thực tế đã hạch toán)
	g.Go(func() error {
		var err error
		costStats, err = s.transaction.GetPlatformCostSummary(gCtx, db_transaction.GetPlatformCostSummaryParams{
			FromCreatedAt: startDate,
			ToCreatedAt:   endDate,
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("lỗi khi lấy thông tin chi phí voucher: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy hiệu suất voucher toàn nền tảng: %w", err)}
	}

	// (Chi phí voucher Sàn chịu = voucher giảm giá đơn hàng + voucher giảm ship)
	costVoucher, _ := entity.ParseAmountToFloat64(costStats.TotalOrderVoucherCost)
	costShip, _ := entity.ParseAmountToFloat64(costStats.TotalShippingDiscountCost)
	// usageStatsValue ,_ = entity.ParseAmountToFloat64(usageStats.)
	usageStats.TotalDiscountValue, _ = entity.ParseAmountToFloat64(usageStats.TotalDiscountValue)

	costStats.TotalOrderVoucherCost, _ = entity.ParseAmountToFloat64(costStats.TotalOrderVoucherCost)
	costStats.TotalPromotionCost, _ = entity.ParseAmountToFloat64(costStats.TotalPromotionCost)
	costStats.TotalShippingDiscountCost, _ = entity.ParseAmountToFloat64(costStats.TotalShippingDiscountCost)
	costStats.TotalProductSubsidyCost, _ = entity.ParseAmountToFloat64(costStats.TotalProductSubsidyCost)
	resp := &entity.PlatformVoucherPerformanceResponse{
		UsageHistoryStats: usageStats,
		PlatformCostStats: costStats,
		TotalVoucherCost:  costVoucher + costShip,
	}

	return resp, nil
}

// === NHÓM 5: Phân tích Shop ===

// ListPlatformShops xử lý API: GET /api/v1/platform/shops
func (s *service) ListPlatformShops(ctx context.Context, params entity.ListPlatformShopsParams) ([]entity.PlatformShopRow, *assets_services.ServiceError) {

	// (Như đã phân tích, chúng ta không có kết nối đến user_service,
	// nên chúng ta lấy danh sách shop từ order_db)

	// (Chúng ta cần một hàm sqlc mới `GetAllShopsStats` giống `GetPlatformTopShopsByGMV`
	// nhưng có phân trang và không chỉ lấy TOP)
	// Tạm thời dùng:
	shopStats, err := s.order.GetPlatformTopShopsByGMV(ctx, db_order.GetPlatformTopShopsByGMVParams{
		// (Bỏ qua date range để lấy TẤT CẢ)
		Limit: params.Limit,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.PlatformShopRow{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin shop: %w", err)}
	}

	resp := make([]entity.PlatformShopRow, len(shopStats))
	for i, stat := range shopStats {
		stat.TotalGmv, _ = entity.ParseAmountToFloat64(stat.TotalGmv)
		resp[i] = entity.PlatformShopRow{
			GetPlatformTopShopsByGMVRow: stat,
			// (Thiếu ShopName, ShopAvatar... vì không có user_service)
		}
	}

	return resp, nil
}

// GetPlatformShopDetail xử lý API: GET /api/v1/platform/shops/{shop_id}/detail
func (s *service) GetPlatformShopDetail(ctx context.Context, shopID string, startDate, endDate time.Time) (*entity.PlatformShopDetailResponse, *assets_services.ServiceError) {

	// API này gọi lại các hàm của Shop (Nhóm 1)

	g, gCtx := errgroup.WithContext(ctx)
	resp := &entity.PlatformShopDetailResponse{}

	// Tác vụ 1: Lấy Overview
	g.Go(func() error {
		var err *assets_services.ServiceError
		resp.Overview, err = s.GetShopOverview(gCtx, shopID, startDate, endDate)
		if err != nil {
			return fmt.Errorf("lỗi khi lấy shop overview: %w", err.Err) // Trả về lỗi (nếu có)
		}
		return nil
	})

	// Tác vụ 2: Lấy Wallet
	g.Go(func() error {
		var err *assets_services.ServiceError
		resp.Wallet, err = s.GetShopWalletSummary(gCtx, shopID)
		if err != nil {
			return fmt.Errorf("lỗi khi lấy shop wallet: %w", err.Err) // Trả về lỗi (nếu có)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy chi tiết shop: %w", err)}
	}

	return resp, nil
}

// === NHÓM 6: Xếp hạng ===

// GetPlatformRankingShops (Wrapper)
func (s *service) GetPlatformRankingShops(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopShopsByGMVRow, *assets_services.ServiceError) {
	shops, err := s.order.GetPlatformTopShopsByGMV(ctx, db_order.GetPlatformTopShopsByGMVParams{
		StartDate: sql.NullTime{Time: params.StartDate, Valid: true},
		EndDate:   sql.NullTime{Time: params.EndDate, Valid: true},
		Limit:     params.Limit,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.GetPlatformTopShopsByGMVRow{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin shop: %w", err)}
	}
	return shops, nil
}

// GetPlatformRankingProducts (Wrapper)
func (s *service) GetPlatformRankingProducts(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopProductsByQuantityRow, *assets_services.ServiceError) {
	products, err := s.order.GetPlatformTopProductsByQuantity(ctx, db_order.GetPlatformTopProductsByQuantityParams{
		StartDate: sql.NullTime{Time: params.StartDate, Valid: true},
		EndDate:   sql.NullTime{Time: params.EndDate, Valid: true},
		Limit:     params.Limit,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.GetPlatformTopProductsByQuantityRow{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy sản phẩm xếp hạng theo số lượng: %w", err)}
	}
	return products, nil
}

// GetPlatformRankingUsers (Wrapper)
func (s *service) GetPlatformRankingUsers(ctx context.Context, params entity.PlatformRankingParams) ([]db_order.GetPlatformTopUsersBySpendRow, *assets_services.ServiceError) {
	users, err := s.order.GetPlatformTopUsersBySpend(ctx, db_order.GetPlatformTopUsersBySpendParams{
		StartDate: sql.NullTime{Time: params.StartDate, Valid: params.StartDate.IsZero() != true},
		EndDate:   sql.NullTime{Time: params.EndDate, Valid: params.EndDate.IsZero() != true},
		Limit:     params.Limit,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []db_order.GetPlatformTopUsersBySpendRow{}, nil
		}
		return nil, &assets_services.ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("lỗi khi lấy thông tin người dùng: %w", err)}
	}
	return users, nil
}

// GetPlatformRankingCategories
func (s *service) GetPlatformRankingCategories(ctx context.Context, params entity.PlatformRankingParams) (*entity.PlatformRankingCategoriesResponse, *assets_services.ServiceError) {
	// Ghi chú Kiến trúc:
	// Như đã phân tích, để thực hiện API này, chúng ta cần:
	// 1. Lấy (order_items) từ `order_db`.
	// 2. Lấy (product_id, category_id) từ `product_db`.
	//
	// Cấu trúc `service` của bạn KHÔNG có kết nối đến `product_db`.
	// Vì vậy, API này KHÔNG THỂ thực hiện được với kiến trúc hiện tại.

	log.Println("Architectural Warning: GetPlatformRankingCategories requires a database connection to 'product_service' which is not provided in the current service struct.")

	return &entity.PlatformRankingCategoriesResponse{
		Message: "API Not Implemented: This feature requires access to Product Database (product_db) to map products to categories. Please update service architecture.",
	}, nil
}
