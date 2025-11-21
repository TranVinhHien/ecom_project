-- ==============================================================
-- 1. THỐNG KÊ TỔNG QUAN (DASHBOARD)
-- Tác dụng: Lấy tổng số session, tin nhắn user, tin nhắn agent.
-- Linh hoạt: Có thể lọc theo khoảng thời gian hoặc lấy toàn bộ.
-- ==============================================================

-- name: GetDashboardStats :one
SELECT 
    COUNT(DISTINCT session_id) AS total_sessions,
    
    -- Ép kiểu về SIGNED (số nguyên) và mặc định là 0 nếu NULL
    CAST(COALESCE(SUM(CASE WHEN author = 'user' THEN 1 ELSE 0 END), 0) AS SIGNED) AS total_user_messages,
    
    -- Tương tự cho agent
    CAST(COALESCE(SUM(CASE WHEN author = 'Host_Agent' THEN 1 ELSE 0 END), 0) AS SIGNED) AS total_agent_messages
FROM events
WHERE 
    (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'));
-- ==============================================================
-- 2. BIỂU ĐỒ HEATMAP (THEO GIỜ TRONG NGÀY)
-- Tác dụng: Xem mật độ tin nhắn theo giờ.
-- ==============================================================

-- name: GetMessageVolumeByHour :many
SELECT 
    HOUR(timestamp) AS hour_of_day,
    COUNT(*) AS message_count
FROM events
WHERE 
    author = 'user'
    AND (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'))
GROUP BY hour_of_day
ORDER BY hour_of_day ASC;

-- ==============================================================
-- 3. TOP NGƯỜI DÙNG TÍCH CỰC
-- Tác dụng: Tìm User nhắn nhiều nhất.
-- Tham số limit_count bắt buộc để giới hạn số lượng trả về (ví dụ: Top 10).
-- ==============================================================

-- name: GetTopActiveUsers :many
SELECT 
    user_id, 
    COUNT(*) AS message_count
FROM events 
WHERE 
    author = 'user'
    AND (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'))
GROUP BY user_id 
ORDER BY message_count DESC 
LIMIT ?;

-- ==============================================================
-- 4. PHÂN TÍCH CHỦ ĐỀ (TOPIC ANALYSIS)
-- Tác dụng: Thống kê các Topic từ cột JSON custom_metadata.
-- ==============================================================

-- name: GetTopicStats :many
SELECT 
    custom_metadata->>'$.topic' AS topic,
    COUNT(*) AS count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) AS percentage
FROM events
WHERE 
    author = 'user' 
    AND custom_metadata IS NOT NULL
    AND (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'))
GROUP BY topic
ORDER BY count DESC;

-- ==============================================================
-- 5. PHÂN TÍCH Ý ĐỊNH MUA HÀNG (PURCHASE INTENT)
-- Tác dụng: Đếm số lượng theo mức độ Intent (High, Medium, Low).
-- ==============================================================

-- name: GetPurchaseIntentStats :many
SELECT 
    custom_metadata->>'$.purchase_intent' AS purchase_intent,
    COUNT(*) AS count
FROM events
WHERE 
    author = 'user' 
    AND custom_metadata IS NOT NULL
    AND (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'))
GROUP BY purchase_intent
ORDER BY count DESC;

-- ==============================================================
-- 6. TOP SẢN PHẨM ĐƯỢC NHẮC ĐẾN (ENTITIES)
-- Tác dụng: Trích xuất tên sản phẩm/danh mục từ mảng entities trong JSON.
-- ==============================================================

-- name: GetTopMentionedCategories :many
SELECT 
    custom_metadata
FROM events
WHERE 
    author = 'user' 
    AND custom_metadata IS NOT NULL
    AND (sqlc.narg('start_time') IS NULL OR timestamp >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time') IS NULL OR timestamp <= sqlc.narg('end_time'));