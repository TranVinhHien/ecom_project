-- name: CreateCategory :exec
INSERT INTO category (
  category_id, name, `key`, `path`, parent, image
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: DeleteCategory :exec
DELETE FROM category
WHERE category_id = ?;

-- name: UpdateCategory :exec
UPDATE category
SET 
  name = COALESCE(sqlc.narg('name'), name),
  `key` = COALESCE(sqlc.narg('key'), `key`),
  `path` = COALESCE(sqlc.narg('path'), `path`),
  parent = COALESCE(sqlc.narg('parent'), parent),
  image = COALESCE(sqlc.narg('image'), image)
WHERE category_id = ?;

-- name: GetCategory :one
SELECT * FROM category
WHERE category_id = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM category
ORDER BY name;

-- name: ListCategoriesPaged :many
SELECT * FROM category
ORDER BY name
LIMIT ? OFFSET ?;

-- name: GetSubCategories :many
SELECT * FROM category
WHERE parent = ?;

-- name: GetRootCategories :many
SELECT * FROM category
WHERE parent IS NULL;

-- name: SearchCategoriesByName :many
SELECT * FROM category
WHERE name LIKE '%' || ? || '%'
ORDER BY name;

-- name: GetCategoryByPath :one
SELECT * FROM category
WHERE path = ? LIMIT 1;
-- name: UpdateCategoryParent :exec
UPDATE category
SET parent = ?
WHERE category_id = ?;

-- name: UpdateCategoryImage :exec
UPDATE category
SET image = ?
WHERE category_id = ?;

-- name: CountCategories :one
SELECT COUNT(*) AS total FROM category;
