CREATE DATABASE IF NOT EXISTS ecommerce_interact_db;
USE ecommerce_interact_db;
-- Thiết lập mã hóa UTF-8
-- =================================================================
-- DATABASE SCHEMA FOR ORDER SERVICE
-- =================================================================
SET NAMES utf8mb4;
SET time_zone = '+07:00';


CREATE TABLE `customer_feedback` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của phiếu phản hồi',
  `user_id` CHAR(36) NULL COMMENT 'ID của người dùng (nếu đã đăng nhập)',
  `email` VARCHAR(50) NULL COMMENT 'Email của người dùng (nếu đã đăng nhập)',
  `phone` VARCHAR(50) NULL COMMENT 'Số điện thoại của người dùng (nếu đã đăng nhập)',

  -- === Dữ liệu từ Form ===
  `category` VARCHAR(50) NOT NULL COMMENT 'Phân loại (ví dụ: BUG, COMPLAINT, SUGGESTION)',
  `content` TEXT NOT NULL COMMENT 'Nội dung chi tiết khách hàng nhập từ form',
  
  -- === Dữ liệu cho CSKH xử lý ====
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`)

) ENGINE=InnoDB COMMENT='Lưu trữ các phiếu phản hồi (tickets) từ khách hàng';


CREATE TABLE `message_ratings` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  
  -- === Dữ liệu Liên kết (Để Debug & Xem Context) ===
  `event_id` VARCHAR(128) NOT NULL COMMENT 'FK (logic) tới "events" ADK. Vẫn RẤT QUAN TRỌNG để đội ngũ dev debug.',
  `session_id` VARCHAR(255) NOT NULL COMMENT 'FK (logic) tới "sessions" ADK. Dùng để xem toàn bộ lịch sử hội thoại.',
  `user_id` CHAR(36) NULL COMMENT 'ID của người dùng (nếu đã đăng nhập)',
  
  -- === Dữ liệu Đánh giá Cốt lõi ===
  `rating` TINYINT NOT NULL COMMENT 'Đánh giá (1 = Like, -1 = Dislike)',
  
  -- === [MỚI] Dữ liệu Snapshot (Để Báo cáo Nhanh) ===
  `user_prompt` TEXT NULL COMMENT 'Snapshot câu hỏi/prompt của người dùng đã dẫn đến câu trả lời này',
  `agent_response` TEXT NULL COMMENT 'Snapshot câu trả lời của Agent đã được đánh giá',
  
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`),
  
  -- Mỗi event chỉ được đánh giá 1 lần trong 1 phiên chat
  UNIQUE KEY `uq_event_session` (`event_id`, `session_id`), 
  
  KEY `idx_event_id` (`event_id`),
  KEY `idx_rating_time` (`created_at`, `rating`) -- Tối ưu cho dashboard thống kê
) ENGINE=InnoDB COMMENT='Lưu trữ đánh giá Like/Dislike (và context) cho từng kết quả (event) của Agent';