
-- SKU_ATTR (sku_attr) CRUD

-- name: CreateSKUAttr :exec
INSERT INTO sku_attr (sku_id, option_value_id, product_id)
VALUES (sqlc.arg('sku_id'), sqlc.arg('option_value_id'), sqlc.arg('product_id'));

-- name: DeleteSKUAttr :exec
DELETE FROM sku_attr WHERE sku_id = sqlc.arg('sku_id') AND option_value_id = sqlc.arg('option_value_id');

-- name: ListSKUOptionValuesByProductID :many
SELECT *
FROM sku_attr sa
WHERE sa.product_id = sqlc.arg('product_id');
