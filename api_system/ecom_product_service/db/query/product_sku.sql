
-- PRODUCT SKU (product_sku) CRUD

-- name: CreateProductSKU :exec
INSERT INTO product_sku (
  id, product_id, sku_code, price, quantity, weight
) VALUES (
  sqlc.arg('id'),
  sqlc.arg('product_id'),
  sqlc.arg('sku_code'),
  sqlc.arg('price'), 
  sqlc.arg('quantity'), 
  sqlc.arg('weight')
);

-- name: UpdateProductSKU :exec
UPDATE product_sku
SET
  sku_code = COALESCE(sqlc.narg('sku_code'), sku_code),
  price = COALESCE(sqlc.narg('price'), price),
  quantity = COALESCE(sqlc.narg('quantity'), quantity),
  quantity_reserver = COALESCE(sqlc.narg('quantity_reserver'), quantity_reserver),
  weight = COALESCE(sqlc.narg('weight'), weight),
  update_date = NOW()
WHERE id = sqlc.arg('id');

-- name: DeleteProductSKU :exec
DELETE FROM product_sku WHERE id = sqlc.arg('id');

-- name: GetProductSKU :one
SELECT * FROM product_sku WHERE id = sqlc.arg('id') LIMIT 1;

-- name: ListSKUsByProduct :many
SELECT * FROM product_sku
WHERE product_id = sqlc.arg('product_id')
ORDER BY price;
