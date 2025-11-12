

-- OPTION VALUE (option_value) CRUD

-- name: CreateOptionValue :exec
INSERT INTO option_value (
  id, option_name, `value`, product_id, image
) VALUES (
  sqlc.arg('id'),
  sqlc.arg('option_name'),
  sqlc.arg('value'),
  sqlc.arg('product_id'),
  sqlc.arg('image')
  );

-- name: UpdateOptionValue :exec
UPDATE option_value
SET
  option_name = COALESCE(sqlc.narg('option_name'), option_name),
  `value` = COALESCE(sqlc.narg('value'), `value`),
  image = COALESCE(sqlc.narg('image'), image)
WHERE id = sqlc.arg('id');

-- name: DeleteOptionValue :exec
DELETE FROM option_value WHERE id = sqlc.arg('id');

-- name: ListOptionValuesByProductID :many
SELECT * FROM option_value WHERE product_id = sqlc.arg('product_id') ORDER BY option_name;