

-- =================================================================
-- Queries for `order_items` table
-- =================================================================

-- name: CreateOrderItem :exec
INSERT INTO order_items (
  id, shop_order_id, product_id, sku_id, quantity, original_unit_price, final_unit_price, total_price,
  promotions_snapshot, product_name_snapshot, product_image_snapshot, sku_attributes_snapshot
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: ListOrderItemsByShopOrderID :many
SELECT * FROM order_items
WHERE shop_order_id = ?;


-- -- =================================================================
-- -- Queries for `order_status_history` table
-- -- =================================================================

-- -- name: CreateOrderStatusHistory :exec
-- INSERT INTO order_status_history (
--   shop_order_id, status, notes, created_by
-- ) VALUES (
--   ?, ?, ?, ?
-- );

-- -- name: ListOrderStatusHistoryByShopOrderID :many
-- SELECT * FROM order_status_history
-- WHERE shop_order_id = ?
-- ORDER BY created_at ASC;

-- name: GetOrderItemsByShopOrderIDs :many
-- Lấy tất cả items thuộc các shop_orders (để gửi event 'order_cancelled' cho Product Service)
SELECT shop_order_id, sku_id, quantity
FROM order_items
WHERE shop_order_id IN (sqlc.slice(shop_order_ids)); --




-- name: GetProductTotalSold :many
-- lấy tổng số lượng đã bán của các product_ids trong các đơn hàng có trạng thái 'PROCESSING', 'SHIPPED', 'COMPLETED'(đang dùng cho product_service)
SELECT
  oi.product_id,
CAST(COALESCE(SUM(oi.quantity), 0) AS SIGNED) AS total_sold
FROM
  order_items AS oi
JOIN
  shop_orders AS so ON oi.shop_order_id = so.id
WHERE
  oi.product_id IN (sqlc.slice(product_ids))
  AND so.status IN ('PROCESSING', 'SHIPPED', 'COMPLETED')
GROUP BY
  oi.product_id;