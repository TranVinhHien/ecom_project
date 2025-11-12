-- =================================================================
-- I. API CHO SHOP (Shop-facing APIs)
-- =================================================================

-- name: ListVouchersByOwner :many
-- Tác dụng: Lấy danh sách voucher Shop đã tạo (API: GET /shop/vouchers)
SELECT * FROM vouchers
WHERE
    owner_id = ? AND owner_type = 'SHOP'
    AND (sqlc.narg(is_active_filter) IS NULL OR is_active = sqlc.narg(is_active_filter))
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetVoucherUsagePerformanceByOwner :one
-- Tác dụng: Thống kê hiệu suất voucher của Shop (API: GET /shop/vouchers/performance)
-- Lưu ý: CSDL voucher_usage_history thiếu shop_order_id, nên chúng ta join bằng voucher_id
SELECT
    COUNT(vuh.id) AS total_usage_count,
    COALESCE(SUM(vuh.discount_amount), 0.00) AS total_discount_value
FROM voucher_usage_history vuh
JOIN vouchers v ON vuh.voucher_id = v.id
WHERE
    v.owner_id = ? AND v.owner_type = 'SHOP'
    AND (vuh.used_at BETWEEN ? AND ?);

-- name: GetVoucherUsageHistory :many
-- Tác dụng: Lấy lịch sử sử dụng của 1 voucher (API: GET /shop/vouchers/{id}/detail)
SELECT * FROM voucher_usage_history
WHERE
    voucher_id = ?
ORDER BY used_at DESC
LIMIT ? OFFSET ?;


-- =================================================================
-- II. API CHO SÀN (Platform-facing APIs)
-- =================================================================

-- name: ListPlatformVouchers :many
-- Tác dụng: Lấy tất cả voucher trên Sàn (API: GET /platform/vouchers)
-- ĐÃ SỬA: Bỏ các cast '::text'
SELECT * FROM vouchers
WHERE
    (sqlc.narg(owner_type_filter) IS NULL OR owner_type = sqlc.narg(owner_type_filter))
    AND (sqlc.narg(is_active_filter) IS NULL OR is_active = sqlc.narg(is_active_filter))
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetPlatformVoucherPerformance :one
-- TácG: Thống kê hiệu suất voucher Sàn (API: GET /platform/vouchers/performance/platform)
SELECT
    COUNT(vuh.id) AS total_usage_count,
    COALESCE(SUM(vuh.discount_amount), 0.00) AS total_discount_value
FROM voucher_usage_history vuh
JOIN vouchers v ON vuh.voucher_id = v.id
WHERE
    v.owner_type = 'PLATFORM'
    AND (vuh.used_at BETWEEN ? AND ?);