-- =================================================================
-- Queries for `shop_order_settlements` table
-- =================================================================

-- name: CreateShopOrderSettlement :exec
-- Creates the financial snapshot and settlement record for a shop order.
INSERT INTO shop_order_settlements (
  id,
  shop_order_id,
  order_transaction_id,
  status,
  order_subtotal,
  shop_funded_product_discount,
  site_funded_product_discount,
  shop_voucher_discount,
  shop_shipping_discount,
  shipping_fee,
  commission_fee,
  net_settled_amount
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSettlementByID :one
SELECT * FROM shop_order_settlements
WHERE id = ? LIMIT 1;

-- name: GetSettlementByShopOrderID :one
-- Finds the settlement record for a specific shop order ID.
SELECT * FROM shop_order_settlements
WHERE shop_order_id = ? LIMIT 1;

-- name: UpdateSettlementStatusToFundsHeld :exec
-- Marks the settlement as awaiting the holding period after successful delivery.
UPDATE shop_order_settlements
SET
  status = 'FUNDS_HELD',
  order_completed_at = ? -- Pass the completion timestamp from Order Service
WHERE shop_order_id = ?;

-- name: ListEligibleSettlementsForProcessing :many
-- Finds settlements ready for automatic payout processing (status is FUNDS_HELD and holding period passed).
-- The '?' parameter should be calculated as NOW() - holding_period (e.g., NOW() - INTERVAL 7 DAY)
SELECT * FROM shop_order_settlements
WHERE status = 'FUNDS_HELD' AND order_completed_at IS NOT NULL AND order_completed_at <= ?;

-- name: UpdateSettlementStatusToSettled :exec
-- Marks the settlement as completed after funds are moved to the shop's available balance.
UPDATE shop_order_settlements
SET
  status = 'SETTLED',
  settled_at = NOW()
WHERE id = ?;

-- name: UpdateSettlementStatusToFailed :exec
-- Marks the settlement processing as failed (e.g., if accounting entries fail).
UPDATE shop_order_settlements
SET
  status = 'FAILED'
WHERE id = ?;