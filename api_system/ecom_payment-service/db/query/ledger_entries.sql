
-- =================================================================
-- Queries for `ledger_entries` table
-- =================================================================

-- name: CreateLedgerEntry :exec
INSERT INTO ledger_entries (
  ledger_id, transaction_id, amount, type, description
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: ListLedgerEntriesByLedgerIDPaged :many
SELECT * FROM ledger_entries
WHERE ledger_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
