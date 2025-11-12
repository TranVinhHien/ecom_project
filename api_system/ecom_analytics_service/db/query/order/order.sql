-- =================================================================
-- I. API CHO SHOP (Shop-facing APIs)
-- =================================================================

-- name: GetShopOrderOverview :one
-- Tác dụng: Lấy chỉ số tổng quan cho Dashboard của Shop (API: GET /shop/overview)
SELECT
    COUNT(*) AS total_orders,
    COALESCE(SUM(subtotal), 0.00) AS total_gmv,
    COUNT(CASE WHEN status = 'PROCESSING' THEN 1 END) AS processing_orders
FROM shop_orders
WHERE
    shop_id = ?
    AND status NOT IN ('CANCELLED', 'AWAITING_PAYMENT')
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    );

-- name: ListShopOrders :many
-- Tác dụng: Lấy danh sách đơn hàng (phân trang) cho Shop (API: GET /shop/orders)
SELECT * FROM shop_orders
WHERE
    shop_id = ?
    AND (sqlc.narg(status_filter) IS NULL OR status = sqlc.narg(status_filter))
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetShopOrderByID :one
-- Tác dụng: Lấy chi tiết 1 đơn hàng của Shop (API: GET /shop/orders/{id}/enriched)
SELECT * FROM shop_orders
WHERE id = ? AND shop_id = ?;

-- name: GetOrderItemsByShopOrderID :many
-- Tác dụng: Lấy các sản phẩm của 1 đơn hàng Shop (API: GET /shop/orders/{id}/enriched)
SELECT * FROM order_items
WHERE shop_order_id = ?;

-- name: GetShopOrderIDs :many
-- Tác dụng: Lấy danh sách ID đơn hàng của Shop (Dùng nội bộ để truy vấn CSDL khác)
SELECT id FROM shop_orders
WHERE shop_id = ?;

-- name: GetShopRevenueTimeSeries :many
-- Tác dụng: Lấy dữ liệu doanh thu GMV theo ngày để vẽ biểu đồ (API: GET /shop/revenue/timeseries)
SELECT
    DATE(completed_at) AS report_date,
    COALESCE(SUM(subtotal), 0.00) AS gmv
FROM shop_orders
WHERE
    shop_id = ?
    AND (
        sqlc.narg(from_completed_at) IS NULL 
        OR sqlc.narg(to_completed_at) IS NULL 
        OR completed_at BETWEEN sqlc.narg(from_completed_at) AND sqlc.narg(to_completed_at)
    )
GROUP BY report_date
ORDER BY report_date ASC;

-- name: GetShopTopProductsByRevenue :many
-- Tác dụng: Xếp hạng sản phẩm theo Doanh thu (API: GET /shop/ranking/products/by-revenue)
SELECT
    oi.product_id,
    oi.sku_id,
    oi.product_name_snapshot,
    oi.sku_attributes_snapshot,
    COALESCE(SUM(oi.total_price), 0.00) AS total_revenue
FROM order_items oi
JOIN shop_orders so ON oi.shop_order_id = so.id
WHERE
    so.shop_id = ?
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR so.completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
GROUP BY oi.product_id, oi.sku_id, oi.product_name_snapshot, oi.sku_attributes_snapshot
ORDER BY total_revenue DESC
LIMIT ?;

-- name: GetShopTopProductsByQuantity :many
-- Tác dụng: Xếp hạng sản phẩm theo Số lượng bán (API: GET /shop/ranking/products/by-quantity)
SELECT
    oi.product_id,
    oi.sku_id,
    oi.product_name_snapshot,
    oi.sku_attributes_snapshot,
    COALESCE(SUM(oi.quantity), 0) AS total_quantity
FROM order_items oi
JOIN shop_orders so ON oi.shop_order_id = so.id
WHERE
    so.shop_id = ?
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR so.completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
GROUP BY oi.product_id, oi.sku_id, oi.product_name_snapshot, oi.sku_attributes_snapshot
ORDER BY total_quantity DESC
LIMIT ?;

-- =================================================================
-- II. API CHO SÀN (Platform-facing APIs)
-- =================================================================

-- name: GetPlatformOrderOverview :one
-- Tác dụng: Lấy tổng quan GMV và Đơn hàng toàn Sàn (API: GET /platform/overview)
SELECT
    COUNT(*) AS total_orders,
    COALESCE(SUM(subtotal), 0.00) AS total_gmv,
    COUNT(DISTINCT shop_id) AS total_shops
FROM shop_orders
WHERE
    status = 'COMPLETED'
    AND (
        sqlc.narg(from_completed_at) IS NULL 
        OR sqlc.narg(to_completed_at) IS NULL 
        OR completed_at BETWEEN sqlc.narg(from_completed_at) AND sqlc.narg(to_completed_at)
    );

-- name: ListPlatformOrders :many
-- Tác dụng: Lấy danh sách TẤT CẢ đơn hàng trên Sàn (API: GET /platform/orders)
SELECT * FROM shop_orders
WHERE
    (sqlc.narg(shop_id_filter) IS NULL OR shop_id = sqlc.narg(shop_id_filter))
    AND (sqlc.narg(status_filter) IS NULL OR status = sqlc.narg(status_filter))
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetOrderByID :one
-- Tác dụng: Lấy thông tin đơn hàng TỔNG (API: GET /platform/orders/{id}/detail)
SELECT * FROM orders
WHERE id = ?;

-- name: GetShopOrdersByOrderID :many
-- Tác dụng: Lấy các đơn hàng SHOP con của đơn hàng TỔNG (API: GET /platform/orders/{id}/detail)
SELECT * FROM shop_orders
WHERE order_id = ?;

-- name: GetOrderIDsFromShopOrders :many
-- Tác dụng: Lấy danh sách ID đơn hàng TỔNG (Dùng nội bộ để truy vấn CSDL khác)
-- ĐÃ SỬA: Sửa 'ANY(sqlc.slice())' thành 'IN (sqlc.slice())'
SELECT DISTINCT order_id FROM shop_orders
WHERE id IN (sqlc.slice(shop_order_ids));

-- name: GetPlatformTopShopsByGMV :many
-- Tác dụng: Xếp hạng Shop theo Doanh thu GMV (API: GET /platform/ranking/shops)
SELECT
    shop_id,
    COALESCE(SUM(subtotal), 0.00) AS total_gmv,
    COUNT(*) AS total_orders
FROM shop_orders
WHERE
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    
GROUP BY shop_id
ORDER BY total_gmv DESC
LIMIT ?;

-- name: GetPlatformTopProductsByQuantity :many
-- Tác dụng: Xếp hạng sản phẩm hot nhất toàn Sàn (API: GET /platform/ranking/products)
SELECT
    oi.product_id,
    oi.product_name_snapshot,
    COALESCE(SUM(oi.quantity), 0) AS total_quantity
FROM order_items oi
JOIN shop_orders so ON oi.shop_order_id = so.id
WHERE
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR so.completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    
GROUP BY oi.product_id, oi.product_name_snapshot
ORDER BY total_quantity DESC
LIMIT ?;

-- name: GetPlatformTopUsersBySpend :many
-- Tác dụng: Xếp hạng khách hàng chi tiêu nhiều nhất (API: GET /platform/ranking/users)
SELECT
    user_id,
    COALESCE(SUM(grand_total), 0.00) AS total_spent,
    COUNT(*) AS total_orders
FROM orders
WHERE
    sqlc.narg(start_date) IS NULL 
    OR sqlc.narg(end_date) IS NULL 
    OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
GROUP BY user_id
ORDER BY total_spent DESC
LIMIT ?;

-- name: GetPlatformGMVTimeSeries :many
SELECT
    DATE(completed_at) AS report_date,
    COALESCE(SUM(subtotal), 0.00) AS gmv
FROM shop_orders
WHERE
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR completed_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    
GROUP BY report_date
ORDER BY report_date ASC;