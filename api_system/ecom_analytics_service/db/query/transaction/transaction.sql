-- =================================================================
-- I. API CHO SHOP (Shop-facing APIs)
-- =================================================================

-- name: GetLedgerByOwnerID :one
-- Tác dụng: Lấy thông tin ví của Shop (API: GET /shop/overview, /shop/wallet/summary)
SELECT * FROM account_ledgers
WHERE owner_id = ? AND owner_type = 'SHOP';

-- name: ListLedgerEntriesByOwnerID :many
-- Tác dụng: Lấy lịch sử giao dịch ví của Shop (API: GET /shop/wallet/ledger-entries)
SELECT le.*
FROM ledger_entries le
JOIN account_ledgers al ON le.ledger_id = al.id
WHERE
    al.owner_id = ? AND al.owner_type = 'SHOP'
ORDER BY le.created_at DESC
LIMIT ? OFFSET ?;

-- name: GetShopSettlementsByOrderIDs :many
-- Tác dụng: Lấy thông tin đối soát của Shop (API: GET /shop/settlements)
-- ĐÃ SỬA: Sửa 'ANY(sqlc.slice())' thành 'IN (sqlc.slice())'
SELECT * FROM shop_order_settlements
WHERE shop_order_id IN (sqlc.slice(shop_order_ids));

-- name: GetShopSettlementStatsByOrderIDs :one
-- Tác dụng: Thống kê tổng tiền đối soát của Shop (API: GET /shop/wallet/summary)
-- ĐÃ SỬA: Sửa 'ANY(sqlc.slice())' thành 'IN (sqlc.slice())'
SELECT
    COALESCE(SUM(CASE WHEN status = 'SETTLED' THEN net_settled_amount END), 0.00) AS total_settled,
    COALESCE(SUM(CASE WHEN status = 'FUNDS_HELD' THEN net_settled_amount END), 0.00) AS total_funds_held
FROM shop_order_settlements
WHERE shop_order_id IN (sqlc.slice(shop_order_ids));


-- =================================================================
-- II. API CHO SÀN (Platform-facing APIs)
-- =================================================================

-- name: GetPlatformRevenueSummary :one
-- Tác dụng: Lấy Doanh thu & Chi phí của Sàn (API: GET /platform/overview)
SELECT
    COALESCE(SUM(commission_fee), 0.00) AS total_commission,
    COALESCE(SUM(shipping_fee), 0.00) AS total_shipping_revenue
FROM shop_order_settlements
WHERE
    status = 'SETTLED'
    AND settled_at BETWEEN ? AND ?;

-- name: GetPlatformCostSummary :one
-- Tác dụng: Lấy Tổng chi phí Sàn (API: GET /platform/overview)
SELECT
    COALESCE(SUM(site_order_voucher_discount_amount), 0.00) AS total_order_voucher_cost,
    COALESCE(SUM(site_promotion_discount_amount), 0.00) AS total_promotion_cost,
    COALESCE(SUM(site_shipping_discount_amount), 0.00) AS total_shipping_discount_cost,
    COALESCE(SUM(total_site_funded_product_discount), 0.00) AS total_product_subsidy_cost
FROM order_platform_costs
WHERE created_at BETWEEN ? AND ?;

-- name: GetPlatformRevenueTimeSeries :many
-- Tác dụng: Vẽ biểu đồ Doanh thu Sàn (API: GET /platform/finance/revenue-timeseries)
SELECT
    DATE(settled_at) AS report_date,
    COALESCE(SUM(commission_fee), 0.00) AS platform_revenue
FROM shop_order_settlements
WHERE
    status = 'SETTLED'
    AND settled_at BETWEEN ? AND ?
GROUP BY report_date
ORDER BY report_date ASC;

-- name: GetPlatformCostTimeSeries :many
-- Tác dụng: Vẽ biểu đồ Chi phí Sàn (API: GET /platform/finance/revenue-timeseries)
SELECT
    DATE(created_at) AS report_date,
    COALESCE(SUM(
        site_order_voucher_discount_amount +
        site_promotion_discount_amount +
        site_shipping_discount_amount +
        total_site_funded_product_discount
    ), 0.00) AS total_cost
FROM order_platform_costs
WHERE created_at BETWEEN ? AND ?
GROUP BY report_date
ORDER BY report_date ASC;

-- name: ListPlatformTransactions :many
-- Tác dụng: Lấy tất cả giao dịch (API: GET /platform/finance/transactions)
-- ĐÃ SỬA: Bỏ các cast '::text'
SELECT * FROM transactions
WHERE
    (sqlc.narg(type_filter) IS NULL OR type = sqlc.narg(type_filter))
    AND (sqlc.narg(status_filter) IS NULL OR status = sqlc.narg(status_filter))
    AND (created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date))
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListPlatformSettlements :many
-- Tác dụng: Lấy tất cả bản ghi đối soát (API: GET /platform/finance/settlements)
-- ĐÃ SỬA: Bỏ các cast '::text'
SELECT * FROM shop_order_settlements
WHERE
    (sqlc.narg(status_filter) IS NULL OR status = sqlc.narg(status_filter))
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR order_completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
ORDER BY order_completed_at DESC
LIMIT ? OFFSET ?;

-- name: ListPlatformLedgers :many
-- Tác dụng: Lấy tất cả ví trên Sàn (API: GET /platform/finance/ledgers)
-- ĐÃ SỬA: Bỏ các cast '::text'
SELECT * FROM account_ledgers
WHERE
    (sqlc.narg(owner_type_filter) IS NULL OR owner_type = sqlc.narg(owner_type_filter))
LIMIT ? OFFSET ?;

-- name: ListLedgerEntriesByLedgerID :many
-- Tác dụng: Lấy lịch sử giao dịch của 1 ví cụ thể (API: GET /platform/finance/ledgers/{id}/entries)
SELECT * FROM ledger_entries
WHERE ledger_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetTransactionStatusesByOrderIDs :many
-- Tác dụng: Lấy trạng thái thanh toán cho 1 loạt đơn hàng (Dùng nội bộ)
-- ĐÃ SỬA: Sửa 'ANY(sqlc.slice())' thành 'IN (sqlc.slice())'
SELECT order_id, status FROM transactions
WHERE order_id IN (sqlc.slice(order_ids))
  AND type = 'PAYMENT';