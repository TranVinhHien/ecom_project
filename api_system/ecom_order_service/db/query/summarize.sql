-- name: CheckReviewPermission :one
-- Xác thực quyền review:
-- 1. order_item_id có tồn tại (JOIN order_items)
-- 2. order_item_id đó có thuộc về user_id này (WHERE o.user_id = ?)
-- 3. Đơn hàng shop (shop_order) chứa item đó PHẢI ở trạng thái 'COMPLETED' (WHERE so.status = 'COMPLETED')
SELECT
    oi.product_id,
    oi.sku_id,
    oi.sku_attributes_snapshot
FROM 
    shop_orders so
JOIN 
    orders o ON so.order_id = o.id
JOIN 
    order_items oi ON oi.shop_order_id = so.id
WHERE 
    oi.id = ?           -- (sqlc.arg: order_item_id)
    AND o.user_id = ?   -- (sqlc.arg: user_id)
    AND so.status = 'COMPLETED'
LIMIT 1;