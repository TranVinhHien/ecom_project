-- =================================================================
-- Bảng 6: product_comment (Thiết kế nghiêm ngặt - Chỉ người mua)
-- =================================================================
CREATE TABLE `product_comment` (
  `comment_id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của đánh giá',
  
  -- === Liên kết bắt buộc ===
  `order_item_id` CHAR(36) NOT NULL COMMENT 'UUID của order_items.id từ Order Service. Đây là "vé" để review.',
  `product_id` VARCHAR(36) NOT NULL COMMENT 'FK (logic) tới product.id. Dùng để tra cứu nhanh.',
  `sku_id` VARCHAR(36) NOT NULL COMMENT 'FK (logic) tới product_sku.id. Dùng để tra cứu nhanh.',
  `user_id` CHAR(36) NOT NULL COMMENT 'UUID của người dùng (từ Identity Service)',

  -- === Dữ liệu Snapshot (Yêu cầu 3) ===
  `sku_name_snapshot` NVARCHAR(500)  DEFAULT NULL COMMENT 'Snapshot tên chi tiết SKU (ví dụ: "Màu Sắc: Đỏ, Size: L")',
  
  -- === Nội dung đánh giá ===
  `rating` TINYINT NOT NULL COMMENT 'Điểm đánh giá (1-5 sao)',
  `title` NVARCHAR(255) DEFAULT NULL COMMENT 'Tiêu đề của đánh giá',
  `content` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Nội dung chi tiết của đánh giá',
  `media` TEXT  COMMENT 'Mảng JSON chứa URLs hình ảnh/video',
  
  -- === Quản lý & Phản hồi ===
  `parent_id` CHAR(36) DEFAULT NULL COMMENT 'FK tự tham chiếu (product_comment.id) cho phép Shop/Admin trả lời',
  
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`comment_id`),
  -- Ràng buộc CỰC KỲ QUAN TRỌNG: 
  -- Đảm bảo mỗi item trong đơn hàng chỉ được đánh giá 1 LẦN DUY NHẤT.
  UNIQUE KEY `uq_order_item_id` (`order_item_id`),
  

  KEY `idx_user_id` (`user_id`),
  KEY `idx_sku_id` (`sku_id`)
) ENGINE=InnoDB COMMENT='Đánh giá sản phẩm (Chỉ dành cho người mua đã xác thực)';

-- =================================================================
-- Bảng 7: review_likes (Bảng theo dõi lượt "Hữu ích")
-- =================================================================
CREATE TABLE `review_likes` (
  `review_id` CHAR(36) NOT NULL COMMENT 'FK tới product_comment.id',
  `user_id` CHAR(36) NOT NULL COMMENT 'UUID của người dùng nhấn "Hữu ích"',
  
  PRIMARY KEY (`review_id`, `user_id`),
  CONSTRAINT `fk_likes_review` FOREIGN KEY (`review_id`) REFERENCES `product_comment` (`comment_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='Theo dõi lượt "Hữu ích" (Helpful) cho mỗi đánh giá';