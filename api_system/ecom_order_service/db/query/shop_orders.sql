
-- =================================================================
-- Queries for `shop_orders` table
-- =================================================================

-- name: CreateShopOrder :exec
INSERT INTO shop_orders (
  id, shop_order_code, order_id, shop_id, status, subtotal, total_discount, total_amount, shipping_fee,
  shop_voucher_code, shop_voucher_discount, processing_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,sqlc.narg('processing_at')
);

-- name: GetShopOrderByID :one
SELECT * FROM shop_orders
WHERE id = ? LIMIT 1;

-- name: ListShopOrdersByOrderID :many
SELECT * FROM shop_orders
WHERE (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'))
  AND order_id = ? 
ORDER BY created_at ASC;


-- -- name: ListShopOrdersByStatus :many
-- SELECT shop_orders.* FROM shop_orders
-- JOIN orders ON shop_orders.order_id = orders.id
-- WHERE 
--     shop_orders.status = ? AND orders.user_id = ?
-- ORDER BY 
--     shop_orders.created_at DESC
-- LIMIT ? OFFSET ?;

-- name: ListShopOrdersByStatus :many
SELECT shop_orders.* FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
    -- 1. Nếu @status là NULL, vế (1) -> TRUE, bỏ qua điều kiện status
    -- 2. Nếu @status có giá trị (vd: 'PROCESSING'), vế (1) -> FALSE, 
    --    và vế (2) -> shop_orders.status = 'PROCESSING'
    (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
    
    AND orders.user_id = sqlc.arg('user_id')
ORDER BY 
    shop_orders.created_at DESC
LIMIT ? OFFSET ?;

-- name: ListShopOrdersSHOP :many
SELECT shop_orders.* FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
    -- 1. Nếu @status là NULL, vế (1) -> TRUE, bỏ qua điều kiện status
    -- 2. Nếu @status có giá trị (vd: 'PROCESSING'), vế (1) -> FALSE, 
    --    và vế (2) -> shop_orders.status = 'PROCESSING'
    (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
    
    AND shop_orders.shop_id = sqlc.arg('shop_id')
ORDER BY 
    shop_orders.created_at DESC
LIMIT ? OFFSET ?;



-- name: CancelShopOrdersByIDs :exec
-- Cập nhật trạng thái một loạt shop_orders thành CANCELLED
UPDATE shop_orders
SET status = 'CANCELLED',
    cancellation_reason = ?, -- Tham số $1 (cancellation_reason)
    cancelled_at = NOW()
WHERE id IN (sqlc.slice(shop_order_ids)); --

-- name: ListShopOrdersByStatusCount :one
SELECT COUNT(*) FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
    (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
   AND orders.user_id = ?;

-- name: ListShopOrdersSHOPCount :one
SELECT COUNT(*) FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
    (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
   AND shop_orders.shop_id = ?;


-- name: ListShopOrdersByShopIDPaged :many
SELECT * FROM shop_orders
WHERE shop_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateShopOrderGeneralInfo :exec
UPDATE shop_orders
SET
  shipping_method = COALESCE(sqlc.narg('shipping_method'), shipping_method),
  tracking_code = COALESCE(sqlc.narg('tracking_code'), tracking_code),
  updated_at = NOW()
WHERE id = ?;


-- name: UpdateShopOrderStatusToRefunded :exec
UPDATE shop_orders
SET
  status = 'REFUNDED',
  processing_at = NOW(),
  updated_at = NOW()
WHERE id = ?;
-- name: UpdateShopOrderStatusToProcessing :exec
UPDATE shop_orders
SET
  status = 'PROCESSING',
  processing_at = NOW(),
  paid_at = NOW(),
  updated_at = NOW()
WHERE id = ?;

-- name: UpdateShopOrderStatusToShipped :exec
UPDATE shop_orders
SET
  status = 'SHIPPED',
  shipped_at = NOW(),
  tracking_code = ?,
  shipping_method = ?,
  updated_at = NOW()
WHERE id = ?;
  
-- name: UpdateShopOrderStatusToCompleted :exec
UPDATE shop_orders
SET
  status = 'COMPLETED',
  completed_at = NOW(),
  updated_at = NOW()
WHERE id = ?;

-- name: UpdateShopOrderStatusToCancelled :exec
UPDATE shop_orders
SET
  status = 'CANCELLED',
  cancelled_at = NOW(),
  cancellation_reason = ?,
  updated_at = NOW()
WHERE id = ?;

-- name: SearchShopOrders :many
SELECT shop_orders.* FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
  -- Lọc theo user_id từ bảng orders
  (sqlc.narg('user_id') IS NULL OR orders.user_id = sqlc.narg('user_id'))
  
  -- Lọc theo status (nếu có)
  AND (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
  
  -- Lọc theo shop_id (nếu có)
  AND (sqlc.narg('shop_id') IS NULL OR shop_orders.shop_id = sqlc.narg('shop_id'))
  
  -- Lọc theo khoảng tổng tiền (min_amount và max_amount)
  AND (sqlc.narg('min_amount') IS NULL OR shop_orders.total_amount >= sqlc.narg('min_amount'))
  AND (sqlc.narg('max_amount') IS NULL OR shop_orders.total_amount <= sqlc.narg('max_amount'))
  
  -- Lọc theo khoảng thời gian created_at
  AND (sqlc.narg('created_from') IS NULL OR shop_orders.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to') IS NULL OR shop_orders.created_at <= sqlc.narg('created_to'))
  
  -- Lọc theo khoảng thời gian paid_at
  AND (sqlc.narg('paid_from') IS NULL OR shop_orders.paid_at >= sqlc.narg('paid_from'))
  AND (sqlc.narg('paid_to') IS NULL OR shop_orders.paid_at <= sqlc.narg('paid_to'))
  
  -- Lọc theo khoảng thời gian processing_at
  AND (sqlc.narg('processing_from') IS NULL OR shop_orders.processing_at >= sqlc.narg('processing_from'))
  AND (sqlc.narg('processing_to') IS NULL OR shop_orders.processing_at <= sqlc.narg('processing_to'))
  
  -- Lọc theo khoảng thời gian shipped_at
  AND (sqlc.narg('shipped_from') IS NULL OR shop_orders.shipped_at >= sqlc.narg('shipped_from'))
  AND (sqlc.narg('shipped_to') IS NULL OR shop_orders.shipped_at <= sqlc.narg('shipped_to'))
  
  -- Lọc theo khoảng thời gian completed_at
  AND (sqlc.narg('completed_from') IS NULL OR shop_orders.completed_at >= sqlc.narg('completed_from'))
  AND (sqlc.narg('completed_to') IS NULL OR shop_orders.completed_at <= sqlc.narg('completed_to'))
  
  -- Lọc theo khoảng thời gian cancelled_at
  AND (sqlc.narg('cancelled_from') IS NULL OR shop_orders.cancelled_at >= sqlc.narg('cancelled_from'))
  AND (sqlc.narg('cancelled_to') IS NULL OR shop_orders.cancelled_at <= sqlc.narg('cancelled_to'))
  
ORDER BY 
  CASE sqlc.narg('sort_by')
    WHEN 'created_at' THEN shop_orders.created_at
    WHEN 'total_amount' THEN shop_orders.total_amount
    WHEN 'paid_at' THEN shop_orders.paid_at
    WHEN 'processing_at' THEN shop_orders.processing_at
    WHEN 'shipped_at' THEN shop_orders.shipped_at
    WHEN 'completed_at' THEN shop_orders.completed_at
    ELSE shop_orders.created_at
  END DESC
LIMIT ? OFFSET ?;

-- name: SearchShopOrdersCount :one
SELECT COUNT(*) FROM shop_orders
JOIN orders ON shop_orders.order_id = orders.id
WHERE 
  -- Lọc theo user_id từ bảng orders
  (sqlc.narg('user_id') IS NULL OR orders.user_id = sqlc.narg('user_id'))
  
  -- Lọc theo status (nếu có)
  AND (sqlc.narg('status') IS NULL OR shop_orders.status = sqlc.narg('status'))
  
  -- Lọc theo shop_id (nếu có)
  AND (sqlc.narg('shop_id') IS NULL OR shop_orders.shop_id = sqlc.narg('shop_id'))
  
  -- Lọc theo khoảng tổng tiền (min_amount và max_amount)
  AND (sqlc.narg('min_amount') IS NULL OR shop_orders.total_amount >= sqlc.narg('min_amount'))
  AND (sqlc.narg('max_amount') IS NULL OR shop_orders.total_amount <= sqlc.narg('max_amount'))
  
  -- Lọc theo khoảng thời gian created_at
  AND (sqlc.narg('created_from') IS NULL OR shop_orders.created_at >= sqlc.narg('created_from'))
  AND (sqlc.narg('created_to') IS NULL OR shop_orders.created_at <= sqlc.narg('created_to'))
  
  -- Lọc theo khoảng thời gian paid_at
  AND (sqlc.narg('paid_from') IS NULL OR shop_orders.paid_at >= sqlc.narg('paid_from'))
  AND (sqlc.narg('paid_to') IS NULL OR shop_orders.paid_at <= sqlc.narg('paid_to'))
  
  -- Lọc theo khoảng thời gian processing_at
  AND (sqlc.narg('processing_from') IS NULL OR shop_orders.processing_at >= sqlc.narg('processing_from'))
  AND (sqlc.narg('processing_to') IS NULL OR shop_orders.processing_at <= sqlc.narg('processing_to'))
  
  -- Lọc theo khoảng thời gian shipped_at
  AND (sqlc.narg('shipped_from') IS NULL OR shop_orders.shipped_at >= sqlc.narg('shipped_from'))
  AND (sqlc.narg('shipped_to') IS NULL OR shop_orders.shipped_at <= sqlc.narg('shipped_to'))
  
  -- Lọc theo khoảng thời gian completed_at
  AND (sqlc.narg('completed_from') IS NULL OR shop_orders.completed_at >= sqlc.narg('completed_from'))
  AND (sqlc.narg('completed_to') IS NULL OR shop_orders.completed_at <= sqlc.narg('completed_to'))
  
  -- Lọc theo khoảng thời gian cancelled_at
  AND (sqlc.narg('cancelled_from') IS NULL OR shop_orders.cancelled_at >= sqlc.narg('cancelled_from'))
  AND (sqlc.narg('cancelled_to') IS NULL OR shop_orders.cancelled_at <= sqlc.narg('cancelled_to'));
