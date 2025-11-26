

-- =================================================================
-- Queries for `transactions` table
-- =================================================================

-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, transaction_code, order_id, payment_method_id, amount, currency, type, status, notes
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);
-- name: GetExpiredPendingTransactions :many
SELECT id, order_id
FROM transactions
WHERE status = 'PENDING'
  AND type = 'PAYMENT'
  AND created_at < ?; -- Tham số $1 (expired_before) sẽ được truyền từ code Go

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = ? LIMIT 1;

-- name: GetPendingPaymentByOrderID :one
SELECT * FROM transactions
WHERE order_id = ? AND type = 'PAYMENT' AND status = 'PENDING'
LIMIT 1;

-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET
  status = ?,
  processed_at = sqlc.narg('processed_at'),
  gateway_transaction_id = sqlc.narg('gateway_transaction_id'),
  notes = sqlc.narg('notes')
WHERE id = ?;

