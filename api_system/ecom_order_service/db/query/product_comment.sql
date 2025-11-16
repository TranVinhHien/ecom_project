-- name: CreateComment :exec
INSERT INTO product_comment (
  comment_id,
  order_item_id,
  product_id,
  sku_id,
  user_id,
  sku_name_snapshot,
  rating,
  title,
  content,
  media,
  parent_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetCommentByOrderItemID :one
-- Dùng để check xem order_item_id này đã được review hay chưa.
SELECT * FROM product_comment
WHERE order_item_id = ?;

-- name: ListCommentsByProduct :many
-- Lấy danh sách bình luận (gốc, không phải trả lời) cho một sản phẩm, hỗ trợ phân trang.
SELECT * FROM product_comment
WHERE product_id = ? 
ORDER BY created_at DESC
LIMIT ? OFFSET ?;


-- name: GetProductRatingStats :one
-- Lấy điểm đánh giá trung bình và tổng số lượt đánh giá cho một sản phẩm.
-- Chỉ tính các bình luận gốc (parent_id IS NULL).
SELECT
  COUNT(*) AS total_reviews,
  AVG(rating) AS average_rating
FROM product_comment
WHERE product_id = ? AND parent_id IS NULL;

-- name: CreateReviewLike :exec
-- Thêm một lượt "Hữu ích" cho review
INSERT INTO review_likes (
  review_id, user_id
) VALUES (
  ?, ?
);

-- name: DeleteReviewLike :exec
-- Bỏ lượt "Hữu ích"
DELETE FROM review_likes
WHERE review_id = ? AND user_id = ?;

-- name: CountReviewLikes :one
-- Đếm số lượt "Hữu ích" của một review
SELECT COUNT(*) FROM review_likes
WHERE review_id = ?;

-- name: CheckBulkOrderItemsReviewed :many
-- Kiểm tra danh sách order_item_id đã được review chưa
-- Chỉ trả về những order_item_id đã có bình luận
SELECT order_item_id FROM product_comment
WHERE order_item_id IN (sqlc.slice('order_item_ids'));

-- name: GetRepliesByCommentID :many
-- Lấy danh sách các bình luận trả lời (replies) cho một comment gốc
SELECT * FROM product_comment
WHERE parent_id = ?
ORDER BY created_at ASC;

-- name: GetBulkProductRatingStats :many
-- Lấy điểm đánh giá trung bình và tổng số lượt đánh giá cho nhiều sản phẩm
-- Chỉ tính các bình luận gốc (parent_id IS NULL).
SELECT
  product_id,
  COUNT(*) AS total_reviews,
  AVG(rating) AS average_rating
FROM product_comment
WHERE product_id IN (sqlc.slice('product_ids')) AND parent_id IS NULL
GROUP BY product_id;