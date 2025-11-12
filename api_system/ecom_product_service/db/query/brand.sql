-- name: CreateBrand :exec
INSERT INTO brand (
  brand_id, name, code, image, create_date, update_date
) VALUES (
  ?, ?, ?, ?, NOW(), NOW()
);

-- name: DeleteBrand :exec
DELETE FROM brand
WHERE brand_id = ?;

-- name: UpdateBrand :exec
UPDATE brand
SET 
  name = COALESCE(sqlc.narg('name'), name),
  code = COALESCE(sqlc.narg('code'), code),
  image = COALESCE(sqlc.narg('image'), image),
  update_date = NOW()
WHERE brand_id = ?;

-- name: GetBrand :one
SELECT * FROM brand
WHERE brand_id = ? LIMIT 1;

-- name: GetBrandByCode :one
SELECT * FROM brand
WHERE code = ? LIMIT 1;

-- name: ListBrands :many
SELECT * FROM brand
ORDER BY create_date DESC;

-- name: ListBrandsPaged :many
SELECT * FROM brand
ORDER BY create_date DESC
LIMIT ? OFFSET ?;

-- name: SearchBrandsByName :many
SELECT * FROM brand
WHERE name LIKE '%' || ? || '%'
ORDER BY name;

-- name: CountBrands :one
SELECT COUNT(*) AS total FROM brand;

-- name: UpdateBrandImage :exec
UPDATE brand
SET image = ?, update_date = NOW()
WHERE brand_id = ?;
