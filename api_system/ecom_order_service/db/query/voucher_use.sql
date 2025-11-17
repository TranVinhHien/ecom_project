-- name: CreateVoucherUser :exec
-- tạo ra voucher cho riêng người dùng
INSERT INTO user_vouchers (
    user_id,
    voucher_id,
    `status`
) VALUES (
    ?, ?, ?
);