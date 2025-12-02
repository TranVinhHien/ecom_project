-- =================================================================
-- Queries for `order_platform_costs` table (New Table)
-- =================================================================

-- name: CreateOrderPlatformCost :exec
-- Records the total costs incurred by the platform for a specific order payment.
INSERT INTO order_platform_costs (
  order_id,
  payment_transaction_id,
  site_order_voucher_discount_amount,
  site_promotion_discount_amount,
  site_shipping_discount_amount,
  total_site_funded_product_discount
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: GetOrderPlatformCostByPaymentTxID :one
-- Retrieves the platform cost details associated with a specific payment transaction.
SELECT * FROM order_platform_costs
WHERE payment_transaction_id = ? LIMIT 1;

-- name: GetOrderPlatformCostByOrderID :one
-- Retrieves the platform cost details associated with a specific order ID.
SELECT * FROM order_platform_costs
WHERE order_id = ? LIMIT 1;