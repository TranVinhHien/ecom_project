-- PRODUCT CRUD & UTILS

-- name: CreateProduct :exec
INSERT INTO product (
  id, name, `key`, description, short_description,
  brand_id, category_id, shop_id,
  image, media,
   product_is_permission_return, product_is_permission_check,
  create_by
) VALUES (
  sqlc.arg('id'),
  sqlc.arg('name'),
  sqlc.arg('key'),
  sqlc.arg('description'), 
  sqlc.arg('short_description'), 
  sqlc.arg('brand_id'), 
  sqlc.arg('category_id'),
  sqlc.arg('shop_id'),
  sqlc.arg('image'), 
  sqlc.arg('media'), 
  COALESCE(sqlc.narg('product_is_permission_return'), TRUE),
  COALESCE(sqlc.narg('product_is_permission_check'), TRUE),
  sqlc.arg('create_by')
);

-- name: UpdateProduct :exec
UPDATE product
SET
  name = COALESCE(sqlc.narg('name'), name),
  `key` = COALESCE(sqlc.narg('key'), `key`),
  description = COALESCE(sqlc.narg('description'), description),
  short_description = COALESCE(sqlc.narg('short_description'), short_description),
  brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
  category_id = COALESCE(sqlc.narg('category_id'), category_id),
  shop_id = COALESCE(sqlc.narg('shop_id'), shop_id),
  image = COALESCE(sqlc.narg('image'), image),
  media = COALESCE(sqlc.narg('media'), media),
  delete_status = COALESCE(sqlc.narg('delete_status'), delete_status),
  product_is_permission_return = COALESCE(sqlc.narg('product_is_permission_return'), product_is_permission_return),
  product_is_permission_check = COALESCE(sqlc.narg('product_is_permission_check'), product_is_permission_check),
  update_by = COALESCE(sqlc.narg('update_by'), update_by),
  update_date = NOW()
WHERE id = sqlc.arg('id');

-- name: IncrementProductTotalSold :exec
UPDATE product
SET 
  total_sold = total_sold + sqlc.arg('quantity'),
  update_date = NOW()
WHERE id = sqlc.arg('id');

-- name: DeleteProduct :exec
DELETE FROM product WHERE id = sqlc.arg('id');

-- name: GetProduct :one
SELECT
  p.id, p.name, p.`key`, p.description, p.short_description,
  p.brand_id, p.category_id, p.shop_id,
  p.image, p.media, p.delete_status,
  p.product_is_permission_return, p.product_is_permission_check,
  p.create_date, p.update_date, p.create_by, p.update_by,
  (SELECT MIN(price) FROM product_sku WHERE product_id = p.id) AS min_price,
  (SELECT MAX(price) FROM product_sku WHERE product_id = p.id) AS max_price,
  (SELECT id FROM product_sku WHERE product_id = p.id ORDER BY price ASC LIMIT 1) AS min_price_sku_id,
  (SELECT id FROM product_sku WHERE product_id = p.id ORDER BY price DESC LIMIT 1) AS max_price_sku_id
FROM product p
WHERE p.id = sqlc.arg('id')
LIMIT 1;

-- name: GetProductByKey :one
SELECT
  p.id, p.name, p.`key`, p.description, p.short_description,
  p.brand_id, p.category_id, p.shop_id,
  p.image, p.media, p.delete_status,
  p.product_is_permission_return, p.product_is_permission_check,
  p.create_date, p.update_date, p.create_by, p.update_by,
  (SELECT MIN(price) FROM product_sku WHERE product_id = p.id) AS min_price,
  (SELECT MAX(price) FROM product_sku WHERE product_id = p.id) AS max_price,
  (SELECT id FROM product_sku WHERE product_id = p.id ORDER BY price ASC LIMIT 1) AS min_price_sku_id,
  (SELECT id FROM product_sku WHERE product_id = p.id ORDER BY price DESC LIMIT 1) AS max_price_sku_id
FROM product p
WHERE p.`key` = sqlc.arg('key')
LIMIT 1;

-- name: GetProductStockTotal :one
SELECT COALESCE(SUM(quantity), 0) AS total_stock
FROM product_sku
WHERE product_id = sqlc.arg('product_id');


-- name: ListProductsAdvanced :many
SELECT 
    p.id,
    p.name,
    p.`key`,
    p.description,
    p.short_description,
    p.brand_id,
    p.category_id,
    p.shop_id,
    p.image,
    p.media,
    p.delete_status,
    p.product_is_permission_return,
    p.product_is_permission_check,
    p.create_date,
    p.update_date,
    p.create_by,
    p.update_by,
    p.total_sold,
    p.min_price,
    p.max_price,
    -- Subquery lấy SKU ID đại diện (giữ nguyên logic tối ưu)
    (SELECT ps.id FROM product_sku ps WHERE ps.product_id = p.id ORDER BY ps.price ASC LIMIT 1) AS min_price_sku_id,
    (SELECT ps.id FROM product_sku ps WHERE ps.product_id = p.id ORDER BY ps.price DESC LIMIT 1) AS max_price_sku_id
FROM product p
WHERE 
    p.delete_status = 'Active'
    -- Bộ lọc động (Dynamic Filtering)
    AND (sqlc.narg('shop_id') IS NULL OR p.shop_id = sqlc.narg('shop_id'))
    AND (sqlc.narg('category_id') IS NULL OR p.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('brand_id') IS NULL OR p.brand_id = sqlc.narg('brand_id'))
    AND (sqlc.narg('price_min') IS NULL OR p.min_price >= sqlc.narg('price_min'))
    AND (sqlc.narg('price_max') IS NULL OR p.min_price <= sqlc.narg('price_max'))
    -- Tìm kiếm (Lưu ý: Nếu cột name là VARCHAR(255) có index, hiệu năng sẽ tốt hơn TEXT)
    AND (sqlc.narg('keyword') IS NULL OR p.name LIKE CONCAT('%', sqlc.narg('keyword'), '%'))
ORDER BY 
    -- Sắp xếp động (Dynamic Sorting)
    CASE WHEN sqlc.narg('sort') = 'best_sell'  THEN p.total_sold END DESC,
    CASE WHEN sqlc.narg('sort') = 'price_asc'  THEN p.min_price END ASC,
    CASE WHEN sqlc.narg('sort') = 'price_desc' THEN p.min_price END DESC,
    CASE WHEN sqlc.narg('sort') = 'name_asc'   THEN p.name END ASC,
    CASE WHEN sqlc.narg('sort') = 'name_desc'  THEN p.name END DESC,
    p.create_date DESC -- Mặc định sắp xếp theo mới nhất nếu không chọn gì
LIMIT ? OFFSET ?;


-- name: CountProductsAdvanced :one
SELECT COUNT(*)
FROM product p
WHERE 
    p.delete_status = 'Active'
    AND (sqlc.narg('shop_id') IS NULL OR p.shop_id = sqlc.narg('shop_id'))
    AND (sqlc.narg('category_id') IS NULL OR p.category_id = sqlc.narg('category_id'))
    AND (sqlc.narg('brand_id') IS NULL OR p.brand_id = sqlc.narg('brand_id'))
    AND (sqlc.narg('price_min') IS NULL OR p.min_price >= sqlc.narg('price_min'))
    AND (sqlc.narg('price_max') IS NULL OR p.min_price <= sqlc.narg('price_max'))
    AND (sqlc.narg('keyword') IS NULL OR p.name LIKE CONCAT('%', sqlc.narg('keyword'), '%'));


-- name: GetAllProductID :many
SELECT id FROM product;

-- name: GetProductIDs :many
SELECT * FROM product 
WHERE id IN (sqlc.slice(product_ids));
