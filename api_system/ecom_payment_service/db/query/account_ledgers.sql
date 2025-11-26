
-- =================================================================
-- Queries for `account_ledgers` table
-- =================================================================

-- name: CreateLedger :exec
INSERT INTO account_ledgers (
  id, owner_id, owner_type
) VALUES (
  ?, ?, ?
);

-- name: GetLedgerByOwnerID :one
SELECT * FROM account_ledgers
WHERE owner_id = ? AND owner_type = ?
LIMIT 1;

-- name: UpdateLedgerBalances :exec
-- Use positive values for CREDIT, negative for DEBIT
UPDATE account_ledgers
SET
  balance = balance + sqlc.arg(balance_change),
  pending_balance = pending_balance + sqlc.arg(pending_balance_change)
WHERE id = ?;

