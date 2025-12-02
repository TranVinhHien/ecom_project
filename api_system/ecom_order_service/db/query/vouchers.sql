-- name: CreateVoucher :exec
INSERT INTO vouchers (
    id,
    name,
    voucher_code,
    owner_type,
    owner_id,
    discount_type,
    discount_value,
    max_discount_amount,
    applies_to_type,
    min_purchase_amount,
    audience_type,
    start_date,
    end_date,
    total_quantity,
    max_usage_per_user,
    is_active
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateVoucher :exec
-- Cập nhật từng phần (Partial Update)
UPDATE vouchers
SET
    name = COALESCE(sqlc.narg(name), name),
    voucher_code = COALESCE(sqlc.narg(voucher_code), voucher_code),
    discount_type = COALESCE(sqlc.narg(discount_type), discount_type),
    discount_value = COALESCE(sqlc.narg(discount_value), discount_value),
    max_discount_amount = COALESCE(sqlc.narg(max_discount_amount), max_discount_amount),
    applies_to_type = COALESCE(sqlc.narg(applies_to_type), applies_to_type),
    min_purchase_amount = COALESCE(sqlc.narg(min_purchase_amount), min_purchase_amount),
    audience_type = COALESCE(sqlc.narg(audience_type), audience_type),
    start_date = COALESCE(sqlc.narg(start_date), start_date),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    total_quantity = COALESCE(sqlc.narg(total_quantity), total_quantity),
    max_usage_per_user = COALESCE(sqlc.narg(max_usage_per_user), max_usage_per_user),
    is_active = COALESCE(sqlc.narg(is_active), is_active)
WHERE
    id = sqlc.arg(id);

-- name: GetPublicVouchers :many
-- Lấy danh sách voucher CÔNG KHAI (cho toàn bộ người dùng)
-- Chỉ lấy voucher còn hiệu lực và còn số lượng
SELECT * FROM vouchers
WHERE
    audience_type = 'PUBLIC'
    AND is_active = TRUE
    AND NOW() BETWEEN start_date AND end_date
    AND used_quantity < total_quantity
ORDER BY
    created_at DESC;

-- name: GetPublicVouchersWithFilter :many
-- Lấy danh sách voucher CÔNG KHAI với bộ lọc
SELECT * FROM vouchers
WHERE
    audience_type = 'PUBLIC'
    AND is_active = TRUE
    AND NOW() BETWEEN start_date AND end_date
    AND used_quantity < total_quantity
    AND (sqlc.narg('owner_type') IS NULL OR owner_type = sqlc.narg('owner_type'))
    AND (sqlc.narg('shop_id') IS NULL OR owner_id = sqlc.narg('shop_id'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
ORDER BY
    CASE 
        WHEN sqlc.arg('sort_by') = 'discount_asc' THEN discount_value
    END ASC,
    CASE 
        WHEN sqlc.arg('sort_by') = 'discount_desc' THEN discount_value
    END DESC,
    CASE 
        WHEN sqlc.arg('sort_by') = 'created_at' OR sqlc.arg('sort_by') = '' THEN created_at
    END DESC;

-- name: GetAssignedVouchersByUser :many
-- Lấy danh sách voucher ĐƯỢC GÁN RIÊNG (cho 1 user)
-- Chỉ lấy voucher còn hiệu lực, còn trạng thái AVAILABLE
SELECT v.*
FROM vouchers v
JOIN user_vouchers uv ON v.id = uv.voucher_id
WHERE
    uv.user_id = sqlc.arg(user_id)
    AND uv.status = 'AVAILABLE'
    AND v.is_active = TRUE
    AND NOW() BETWEEN v.start_date AND v.end_date
ORDER BY
    v.created_at DESC;

-- name: GetAssignedVouchersByUserWithFilter :many
-- Lấy danh sách voucher ĐƯỢC GÁN RIÊNG với bộ lọc
SELECT v.*
FROM vouchers v
JOIN user_vouchers uv ON v.id = uv.voucher_id
WHERE
    uv.user_id = sqlc.arg(user_id)
    AND uv.status = 'AVAILABLE'
    AND v.is_active = TRUE
    AND NOW() BETWEEN v.start_date AND v.end_date
    AND (sqlc.narg('owner_type') IS NULL OR v.owner_type = sqlc.narg('owner_type'))
    AND (sqlc.narg('shop_id') IS NULL OR v.owner_id = sqlc.narg('shop_id'))
    AND (sqlc.narg('applies_to_type') IS NULL OR v.applies_to_type = sqlc.narg('applies_to_type'))
ORDER BY
    CASE 
        WHEN sqlc.arg('sort_by') = 'discount_asc' THEN v.discount_value
    END ASC,
    CASE 
        WHEN sqlc.arg('sort_by') = 'discount_desc' THEN v.discount_value
    END DESC,
    CASE 
        WHEN sqlc.arg('sort_by') = 'created_at' OR sqlc.arg('sort_by') = '' THEN v.created_at
    END DESC;

-- =============================================
-- CÁC HÀM LẤY DỮ LIỆU ĐỂ "CHECK" TRONG CODE LOGIC
-- =============================================

-- name: GetVoucherForValidation :one
-- Lấy thông tin voucher bằng MÃ (code) để kiểm tra
-- Chỉ trả về voucher nếu nó CƠ BẢN hợp lệ (còn hạn, còn lượt)
SELECT * FROM vouchers
WHERE
    voucher_code = ?
    AND is_active = TRUE
    AND NOW() BETWEEN start_date AND end_date
    AND used_quantity < total_quantity
LIMIT 1;

-- name: CountVoucherUsageByUser :one
-- Đếm số lần user đã sử dụng 1 voucher (cho check max_usage_per_user)
SELECT COUNT(*) FROM voucher_usage_history
WHERE
    voucher_id = ? AND user_id = ?;

-- name: GetUserVoucherStatus :one
-- Lấy trạng thái voucher trong ví của user (cho check voucher ĐƯỢC GÁN)
SELECT * FROM user_vouchers
WHERE
    voucher_id = ? AND user_id = ?
LIMIT 1;

-- =============================================
-- CÁC HÀM "COMMAND" KHI TẠO ĐƠN HÀNG
-- =============================================

-- name: IncrementVoucherUsage :execrows
-- Tăng số lượng đã dùng. Dùng :execrows để check race condition
-- (Logic code phải kiểm tra RowsAffected() == 1)
UPDATE vouchers
SET
    used_quantity = used_quantity + 1
WHERE
    id = ? AND used_quantity < total_quantity;

-- name: SetUserVoucherStatus :execrows
-- Cập nhật trạng thái voucher trong ví user (từ AVAILABLE -> USED)
-- (Logic code nên kiểm tra RowsAffected() == 1)
UPDATE user_vouchers
SET
    status = sqlc.arg(status)
WHERE
    voucher_id = sqlc.arg(voucher_id)
    AND user_id = sqlc.arg(user_id)
    AND status = 'AVAILABLE';

-- name: CreateVoucherUsageHistory :exec
-- Ghi lại lịch sử sử dụng voucher
INSERT INTO voucher_usage_history (
    voucher_id,
    user_id,
    discount_amount
) VALUES (
    ?, ?, ?
);
-- name: GetVoucherByIDForValidation :one
-- Lấy thông tin voucher bằng ID để kiểm tra (THÊM MỚI)
-- Chỉ trả về voucher nếu nó CƠ BẢN hợp lệ (còn hạn, còn lượt)
SELECT * FROM vouchers
WHERE
    id = ?
    AND is_active = TRUE
    AND NOW() BETWEEN start_date AND end_date
    AND used_quantity < total_quantity
LIMIT 1;



-- name: GetVoucherByID :one
-- Lấy voucher bằng ID (không check điều kiện)
SELECT * FROM vouchers WHERE id = ? LIMIT 1;

-- name: DecrementVoucherUsage :execrows
-- Giảm số lượng đã dùng (khi hủy đơn)
UPDATE vouchers
SET
    used_quantity = used_quantity - 1
WHERE
    id = ? AND used_quantity > 0;

-- name: GetVoucherUsageHistory :one
-- Lấy 1 dòng lịch sử sử dụng cụ thể
SELECT * FROM voucher_usage_history
WHERE
    voucher_id = ? 
    AND user_id = ? 
LIMIT 1;

-- name: DeleteVoucherUsageHistory :execrows
-- Xóa 1 dòng lịch sử cụ thể (khi hủy đơn)
DELETE FROM voucher_usage_history
WHERE
    id = ?; -- Xóa bằng ID của bảng history

-- name: ResetUserVoucherStatus :execrows
-- Reset trạng thái ví voucher (từ USED về AVAILABLE)
UPDATE user_vouchers
SET
    status = 'AVAILABLE'
WHERE
    voucher_id = ?
    AND user_id = ?
    AND status = 'USED';

-- =============================================
-- CÁC HÀM QUẢN LÝ CHO ADMIN/SELLER
-- =============================================

-- name: CountVouchersForManagement :one
-- Đếm tổng số voucher theo owner với filters
SELECT COUNT(*) as total
FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    );

-- name: ListVouchersForManagementBySortCreatedAtDesc :many
-- Lấy danh sách voucher cho admin/seller - Sắp xếp theo created_at DESC
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListVouchersForManagementBySortCreatedAtAsc :many
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListVouchersForManagementBySortStartDateDesc :many
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY start_date DESC
LIMIT ? OFFSET ?;

-- name: ListVouchersForManagementBySortStartDateAsc :many
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY start_date ASC
LIMIT ? OFFSET ?;

-- name: ListVouchersForManagementBySortEndDateDesc :many
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY end_date DESC
LIMIT ? OFFSET ?;

-- name: ListVouchersForManagementBySortEndDateAsc :many
SELECT * FROM vouchers
WHERE
    owner_id = sqlc.arg(owner_id)
    AND owner_type = sqlc.arg(owner_type)
    AND (sqlc.narg('voucher_code') IS NULL OR voucher_code LIKE sqlc.narg('voucher_code'))
    AND (sqlc.narg('name') IS NULL OR name LIKE sqlc.narg('name'))
    AND (sqlc.narg('discount_type') IS NULL OR discount_type = sqlc.narg('discount_type'))
    AND (sqlc.narg('applies_to_type') IS NULL OR applies_to_type = sqlc.narg('applies_to_type'))
    AND (sqlc.narg('audience_type') IS NULL OR audience_type = sqlc.narg('audience_type'))
    AND (sqlc.narg('is_active') IS NULL OR is_active = sqlc.narg('is_active'))
    AND (
        sqlc.narg('status') IS NULL
        OR (sqlc.narg('status') = 'ACTIVE' AND is_active = 1 AND start_date <= NOW() AND end_date >= NOW() AND used_quantity < total_quantity)
        OR (sqlc.narg('status') = 'EXPIRED' AND end_date < NOW())
        OR (sqlc.narg('status') = 'UPCOMING' AND start_date > NOW())
        OR (sqlc.narg('status') = 'DEPLETED' AND used_quantity >= total_quantity)
    )
ORDER BY end_date ASC
LIMIT ? OFFSET ?;
