-- =================================================================
-- SQLC QUERIES FOR MESSAGE_RATINGS
-- =================================================================

-- name: CreateMessageRating :exec
INSERT INTO message_ratings (
    event_id,
    session_id,
    user_id,
    rating,
    user_prompt,
    agent_response
) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetMessageRatingStats :one
-- Thống kê tổng quan đánh giá message ratings
SELECT
    COUNT(*) AS total_ratings,
    COUNT(CASE WHEN rating = 1 THEN 1 END) AS like_count,
    COUNT(CASE WHEN rating = -1 THEN 1 END) AS dislike_count,
    -- Ép kiểu kết quả cuối cùng sang DOUBLE để sqlc hiểu là float64
    CAST(
        COALESCE(
            ROUND(
                (COUNT(CASE WHEN rating = 1 THEN 1 END) * 100.0) / NULLIF(COUNT(*), 0),
                2
            ), 
            0.0
        ) AS DOUBLE
    ) AS satisfaction_rate
FROM message_ratings
WHERE
    (sqlc.narg('start_date') IS NULL OR created_at >= sqlc.narg('start_date'))
    AND (sqlc.narg('end_date') IS NULL OR created_at <= sqlc.narg('end_date'));

-- name: GetMessageRatingsTimeSeries :many
-- Thống kê theo thời gian (theo ngày)
SELECT
    DATE(created_at) AS report_date,
    COUNT(*) AS total_ratings,
    COUNT(CASE WHEN rating = 1 THEN 1 END) AS like_count,
    COUNT(CASE WHEN rating = -1 THEN 1 END) AS dislike_count,
    -- Tương tự ở đây
    CAST(
        COALESCE(
            ROUND(
                (COUNT(CASE WHEN rating = 1 THEN 1 END) * 100.0) / NULLIF(COUNT(*), 0),
                2
            ), 
            0.0
        ) AS DOUBLE
    ) AS satisfaction_rate
FROM message_ratings
WHERE
    (sqlc.narg('start_date') IS NULL OR created_at >= sqlc.narg('start_date'))
    AND (sqlc.narg('end_date') IS NULL OR created_at <= sqlc.narg('end_date'))
GROUP BY report_date
ORDER BY report_date ASC;
-- name: GetMessageRatingsBySession :many
-- Lấy danh sách ratings theo session để xem chi tiết
SELECT 
    id,
    event_id,
    session_id,
    user_id,
    rating,
    user_prompt,
    agent_response,
    created_at
FROM message_ratings
WHERE
    (sqlc.narg(session_id_filter) IS NULL OR session_id = sqlc.narg(session_id_filter))
    AND (sqlc.narg(user_id_filter) IS NULL OR user_id = sqlc.narg(user_id_filter))
    AND (sqlc.narg(rating_filter) IS NULL OR rating = sqlc.narg(rating_filter))
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- =================================================================
-- SQLC QUERIES FOR CUSTOMER_FEEDBACK
-- =================================================================

-- name: CreateCustomerFeedback :exec
INSERT INTO customer_feedback (
    id,
    user_id,
    email,
    phone,
    category,
    content
) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListCustomerFeedbacks :many
-- Lấy danh sách feedback cho Admin xem
SELECT * FROM customer_feedback
WHERE
    (sqlc.narg(category_filter) IS NULL OR category = sqlc.narg(category_filter))
    AND (sqlc.narg(user_id_filter) IS NULL OR user_id = sqlc.narg(user_id_filter))
    AND (
        sqlc.narg(start_date) IS NULL 
        OR sqlc.narg(end_date) IS NULL 
        OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
    )
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetCustomerFeedbackByID :one
-- Lấy chi tiết 1 feedback
SELECT * FROM customer_feedback
WHERE id = ?;

-- name: GetCustomerFeedbackStats :one
-- Thống kê tổng quan feedback
SELECT
    COUNT(*) AS total_feedbacks,
    COUNT(CASE WHEN category = 'BUG' THEN 1 END) AS bug_count,
    COUNT(CASE WHEN category = 'COMPLAINT' THEN 1 END) AS complaint_count,
    COUNT(CASE WHEN category = 'SUGGESTION' THEN 1 END) AS suggestion_count,
    COUNT(DISTINCT user_id) AS unique_users
FROM customer_feedback
WHERE
    sqlc.narg(start_date) IS NULL 
    OR sqlc.narg(end_date) IS NULL 
    OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date);

-- name: GetCustomerFeedbacksByCategory :many
-- Thống kê feedback theo category
SELECT
    category,
    COUNT(*) AS feedback_count
FROM customer_feedback
WHERE
    sqlc.narg(start_date) IS NULL 
    OR sqlc.narg(end_date) IS NULL 
    OR created_at BETWEEN sqlc.narg(start_date) AND sqlc.narg(end_date)
GROUP BY category
ORDER BY feedback_count DESC;
