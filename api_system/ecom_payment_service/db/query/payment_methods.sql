-- =================================================================
-- Queries for `payment_methods` table
-- =================================================================

-- name: GetPaymentMethodByID :one
SELECT * FROM payment_methods
WHERE id = ? AND is_active = TRUE;

-- name: ListActivePaymentMethods :many
SELECT * FROM payment_methods
WHERE is_active = TRUE
ORDER BY type, name;
