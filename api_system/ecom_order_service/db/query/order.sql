-- =================================================================
-- Queries for `orders` table
-- =================================================================

-- name: CreateOrder :exec
INSERT INTO orders (
  id, order_code, user_id, grand_total, subtotal, total_shipping_fee, total_discount,
  site_order_voucher_code, site_order_voucher_discount, site_shipping_voucher_code, site_shipping_voucher_discount,
  shipping_address_snapshot, payment_method_snapshot, note
) VALUES (
  ?, ?, ?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?, ?
);

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = ? LIMIT 1;

-- name: GetOrderByCode :one
SELECT * FROM orders
WHERE order_code = ? LIMIT 1;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListOrdersByUserIDPaged :many
SELECT * FROM orders
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;


-- name: UpdateOrderTotals :exec
UPDATE orders
SET
  grand_total = ?,
  subtotal = ?,
  total_shipping_fee = ?,
  total_discount = ?,
  updated_at = NOW()
WHERE id = ?;

-- name: UpdateOrderShippingAddress :exec
UPDATE orders
SET
  shipping_address_snapshot = ?,
  updated_at = NOW()
WHERE id = ?;

