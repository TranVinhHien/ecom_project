-- =================================================================
-- 1. Bảng `vouchers` (Bảng trung tâm)
-- Lưu trữ mọi thông tin và điều kiện của voucher.
-- =================================================================
CREATE TABLE `vouchers` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính',
  `name` VARCHAR(255) NOT NULL COMMENT 'Tên hiển thị (Voucher 50k, Freeship)',
  `voucher_code` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Mã voucher (SALE50K). Phải là UNIQUE',
  
  -- === Chủ sở hữu ===
  `owner_type` ENUM('PLATFORM', 'SHOP') NOT NULL COMMENT 'PLATFORM (Sàn) hay SHOP (Nhà bán hàng)',
  `owner_id` CHAR(36) NOT NULL COMMENT 'ID của Shop (nếu owner_type=SHOP) hoặc UUID cố định của Sàn',

  -- === Hành động giảm giá (Discount Action) ===
  `discount_type` ENUM('PERCENTAGE', 'FIXED_AMOUNT') NOT NULL COMMENT 'Giảm theo % hay Số tiền cố định',
  `discount_value` DECIMAL(15, 2) NOT NULL COMMENT 'Giá trị giảm (VD: 5 cho 5%, hoặc 50000 cho 50k)',
  `max_discount_amount` DECIMAL(15, 2) DEFAULT NULL COMMENT 'Số tiền giảm TỐI ĐA (cho discount_type=PERCENTAGE)',
  `applies_to_type` ENUM('ORDER_TOTAL', 'SHIPPING_FEE') NOT NULL DEFAULT 'ORDER_TOTAL' COMMENT 'Giảm trên: Tổng đơn hàng hay Phí vận chuyển',

  -- === Điều kiện (Đơn giản hóa) ===
  `min_purchase_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'ĐIỀU KIỆN: Giá trị đơn hàng tối thiểu để áp dụng',

  -- === Đối tượng (Audience) ===
  `audience_type` ENUM('PUBLIC', 'ASSIGNED') NOT NULL DEFAULT 'PUBLIC' COMMENT 'PUBLIC (công khai) hay ASSIGNED (chỉ định riêng)',

  -- === Giới hạn sử dụng ===
  `start_date` TIMESTAMP NOT NULL COMMENT 'Thời gian bắt đầu hiệu lực',
  `end_date` TIMESTAMP NOT NULL COMMENT 'Thời gian kết thúc hiệu lực',
  `total_quantity` INT NOT NULL COMMENT 'Tổng số lượt có thể sử dụng (giới hạn toàn hệ thống)',
  `used_quantity` INT NOT NULL DEFAULT 0 COMMENT 'Số lượt đã sử dụng (dùng để check race condition)',
  `max_usage_per_user` INT NOT NULL DEFAULT 1 COMMENT 'Số lượt tối đa mỗi user được dùng',
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Bật/Tắt voucher',
  
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='Bảng chính, định nghĩa voucher và các điều kiện đơn giản';


-- =================================================================
-- 2. Bảng `user_vouchers` (Ví voucher của người dùng)
-- Chỉ lưu các voucher đã được "gán" (ASSIGNED) cho người dùng cụ thể.
-- =================================================================
CREATE TABLE `user_vouchers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` CHAR(36) NOT NULL COMMENT 'UUID của người dùng',
  `voucher_id` CHAR(36) NOT NULL COMMENT 'Khóa ngoại tới bảng vouchers',
  `status` ENUM('AVAILABLE', 'USED', 'EXPIRED') NOT NULL DEFAULT 'AVAILABLE' COMMENT 'Trạng thái voucher trong ví',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_user_voucher` (`user_id`, `voucher_id`),
  CONSTRAINT `fk_user_vouchers_voucher` FOREIGN KEY (`voucher_id`) REFERENCES `vouchers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='Lưu trữ các voucher được gán riêng cho người dùng';


-- =================================================================
-- 3. Bảng `voucher_usage_history` (Lịch sử sử dụng)
-- Ghi lại mọi lượt sử dụng để kiểm kê và đối soát.
-- =================================================================
CREATE TABLE `voucher_usage_history` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `voucher_id` CHAR(36) NOT NULL COMMENT 'Khóa ngoại tới bảng vouchers',
  `user_id` CHAR(36) NOT NULL COMMENT 'Người dùng đã sử dụng',
  -- `shop_order_id` CHAR(36) NOT NULL COMMENT 'Áp dụng cho đơn hàng shop nào (tham chiếu tới Order Service)',
  `discount_amount` DECIMAL(15, 2) NOT NULL COMMENT 'Số tiền thực tế đã giảm',
  `used_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  KEY `idx_voucher_user` (`voucher_id`, `user_id`) COMMENT 'Tối ưu cho việc check max_usage_per_user'
  -- KEY `idx_shop_order` (`shop_order_id`)
) ENGINE=InnoDB COMMENT='Lịch sử sử dụng voucher để đối soát';





-- =================================================================
-- BẢNG `vouchers`
-- Giả định hôm nay là 2025-10-27
-- =================================================================

INSERT INTO `vouchers` (
  `id`, `name`, `voucher_code`, `owner_type`, `owner_id`, 
  `discount_type`, `discount_value`, `max_discount_amount`, `applies_to_type`, 
  `min_purchase_amount`, `audience_type`, 
  `start_date`, `end_date`, 
  `total_quantity`, `used_quantity`, `max_usage_per_user`, `is_active`
) VALUES

-- === BLOCK 1: 10 VOUCHER KHÔNG KHẢ DỤNG (HẾT HẠN HOẶC CHƯA TỚI NGÀY) ===

-- 5 Voucher đã hết hạn (end_date < 2025-10-27)
('v-expired-01', 'Voucher Sàn 10K (ĐÃ HẾT HẠN)', 'EXPIRED10K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 10000, NULL, 'ORDER_TOTAL', 
50000, 'PUBLIC', 
'2025-09-01 00:00:00', '2025-09-30 23:59:59', 
1000, 500, 1, TRUE),
('v-expired-02', 'Voucher Shop 15% (ĐÃ HẾT HẠN)', 'EXPIREDSHOP15', 'SHOP', 'SHOP_UUID_123', 
'PERCENTAGE', 15, 20000, 'ORDER_TOTAL', 
100000, 'PUBLIC', 
'2025-09-01 00:00:00', '2025-09-30 23:59:59', 
500, 150, 1, TRUE),
('v-expired-03', 'Freeship 25K (ĐÃ HẾT HẠN)', 'EXPIREDSHIP25', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 25000, NULL, 'SHIPPING_FEE', 
0, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-25 23:59:59', 
2000, 1900, 1, TRUE),
('v-expired-04', 'Voucher Shop 50K (ĐÃ HẾT HẠN)', 'EXPIREDSHOP50', 'SHOP', 'SHOP_UUID_456', 
'FIXED_AMOUNT', 50000, NULL, 'ORDER_TOTAL', 
200000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-25 23:59:59', 
100, 80, 1, TRUE),
('v-expired-05', 'Voucher Sàn 5% (ĐÃ HẾT HẠN)', 'EXPIRED5PCT', 'PLATFORM', '111111111111111111111111111111111111', 
'PERCENTAGE', 5, 50000, 'ORDER_TOTAL', 
0, 'PUBLIC', 
'2025-09-15 00:00:00', '2025-10-15 23:59:59', 
1000, 950, 1, TRUE),

-- 5 Voucher chưa tới ngày (start_date > 2025-10-27)
('v-future-01', 'Voucher 11.11 Sắp Tới', 'SALE1111', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 111000, NULL, 'ORDER_TOTAL', 
1111000, 'PUBLIC', 
'2025-11-11 00:00:00', '2025-11-11 23:59:59', 
5000, 0, 1, TRUE),
('v-future-02', 'Voucher Shop Tháng 11', 'SHOPTHANG11', 'SHOP', 'SHOP_UUID_123', 
'PERCENTAGE', 10, 30000, 'ORDER_TOTAL', 
99000, 'PUBLIC', 
'2025-11-01 00:00:00', '2025-11-30 23:59:59', 
300, 0, 1, TRUE),
('v-future-03', 'Freeship Cuối Tuần (Tới)', 'SHIPCUOITUAN1', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 15000, NULL, 'SHIPPING_FEE', 
50000, 'PUBLIC', 
'2025-11-01 00:00:00', '2025-11-02 23:59:59', 
1000, 0, 1, TRUE),
('v-future-04', 'Voucher Shop Mới 30K', 'NEWSHOP30K', 'SHOP', 'SHOP_UUID_456', 
'FIXED_AMOUNT', 30000, NULL, 'ORDER_TOTAL', 
150000, 'PUBLIC', 
'2025-11-01 00:00:00', '2025-11-15 23:59:59', 
100, 0, 1, TRUE),
('v-future-05', 'Voucher Sàn 8% (Sắp tới)', 'FUTURE8PCT', 'PLATFORM', '111111111111111111111111111111111111', 
'PERCENTAGE', 8, 40000, 'ORDER_TOTAL', 
100000, 'PUBLIC', 
'2025-10-28 00:00:00', '2025-10-31 23:59:59', 
1000, 0, 1, TRUE),


-- === BLOCK 2: 20 VOUCHER ĐANG KHẢ DỤNG (start_date <= 2025-10-27 <= end_date) ===

-- 17 Voucher CÔNG KHAI (PUBLIC)
('v-public-01', 'Voucher Sàn 15K Đơn 99K', 'SAN15K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 15000, NULL, 'ORDER_TOTAL', 
99000, 'PUBLIC', 
'2025-10-25 00:00:00', '2025-10-31 23:59:59', 
2000, 100, 1, TRUE),
('v-public-02', 'Freeship MAX 15K', 'FREESHIP15', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 15000, NULL, 'SHIPPING_FEE', 
45000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
10000, 1500, 2, TRUE),
('v-public-03', 'Voucher Sàn 10% Tối đa 50K', 'SAN10PCT', 'PLATFORM', '111111111111111111111111111111111111', 
'PERCENTAGE', 10, 50000, 'ORDER_TOTAL', 
150000, 'PUBLIC', 
'2025-10-20 00:00:00', '2025-10-30 23:59:59', 
1000, 250, 1, TRUE),
('v-public-04', 'Voucher Shop 1 20K Đơn 150K', 'SHOP1_20K', 'SHOP', 'SHOP_UUID_123', 
'FIXED_AMOUNT', 20000, NULL, 'ORDER_TOTAL', 
150000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
200, 30, 1, TRUE),
('v-public-05', 'Voucher Shop 1 5% Đơn 50K', 'SHOP1_5PCT', 'SHOP', 'SHOP_UUID_123', 
'PERCENTAGE', 5, 10000, 'ORDER_TOTAL', 
50000, 'PUBLIC', 
'2025-10-15 00:00:00', '2025-10-31 23:59:59', 
500, 50, 2, TRUE),
('v-public-06', 'Voucher Shop 2 30K Đơn 250K', 'SHOP2_30K', 'SHOP', 'SHOP_UUID_456', 
'FIXED_AMOUNT', 30000, NULL, 'ORDER_TOTAL', 
250000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
150, 10, 1, TRUE),
('v-public-07', 'Voucher Sàn 100K Đơn 500K', 'SAN100K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 100000, NULL, 'ORDER_TOTAL', 
500000, 'PUBLIC', 
'2025-10-25 00:00:00', '2025-10-28 23:59:59', 
500, 10, 1, TRUE),
('v-public-08', 'Freeship 5K Đơn 0Đ', 'SHIP5K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 5000, NULL, 'SHIPPING_FEE', 
0, 'PUBLIC', 
'2025-10-27 00:00:00', '2025-10-27 23:59:59', 
1000, 50, 1, TRUE),
('v-public-09', 'Voucher Sàn 50K Đơn 200K', 'SAN50K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 50000, NULL, 'ORDER_TOTAL', 
200000, 'PUBLIC', 
'2025-10-20 00:00:00', '2025-10-31 23:59:59', 
1000, 500, 1, TRUE),
('v-public-10', 'Voucher Shop 1 10K Đơn 50K', 'SHOP1_10K', 'SHOP', 'SHOP_UUID_123', 
'FIXED_AMOUNT', 10000, NULL, 'ORDER_TOTAL', 
50000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
300, 100, 1, TRUE),
('v-public-11', 'Voucher Shop 2 15% Đơn 100K', 'SHOP2_15PCT', 'SHOP', 'SHOP_UUID_456', 
'PERCENTAGE', 15, 30000, 'ORDER_TOTAL', 
100000, 'PUBLIC', 
'2025-10-20 00:00:00', '2025-10-31 23:59:59', 
200, 15, 1, TRUE),
('v-public-12', 'Voucher Sàn 30K Đơn 150K', 'SAN30K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 30000, NULL, 'ORDER_TOTAL', 
150000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
2000, 1800, 1, TRUE),
('v-public-13', 'Freeship MAX 30K', 'FREESHIP30', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 30000, NULL, 'SHIPPING_FEE', 
150000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
1000, 300, 1, TRUE),
('v-public-14', 'Voucher Sàn 8% Tối đa 80K', 'SAN8PCT', 'PLATFORM', '111111111111111111111111111111111111', 
'PERCENTAGE', 8, 80000, 'ORDER_TOTAL', 
300000, 'PUBLIC', 
'2025-10-20 00:00:00', '2025-10-30 23:59:59', 
500, 50, 1, TRUE),
('v-public-15', 'Voucher Shop 1 50K Đơn 300K', 'SHOP1_50K', 'SHOP', 'SHOP_UUID_123', 
'FIXED_AMOUNT', 50000, NULL, 'ORDER_TOTAL', 
300000, 'PUBLIC', 
'2025-10-15 00:00:00', '2025-10-31 23:59:59', 
100, 10, 1, TRUE),
('v-public-16', 'Voucher Sàn 20K Đơn 100K', 'SAN20K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 20000, NULL, 'ORDER_TOTAL', 
100000, 'PUBLIC', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
3000, 100, 1, TRUE),
('v-public-17', 'Freeship 10K Đơn 25K', 'SHIP10K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 10000, NULL, 'SHIPPING_FEE', 
25000, 'PUBLIC', 
'2025-10-25 00:00:00', '2025-10-31 23:59:59', 
1000, 100, 1, TRUE),

-- 3 Voucher RIÊNG (ASSIGNED) cho user123
('v-user123-01', 'Voucher Sàn 50K (RIÊNG)', 'RIENG50K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 50000, NULL, 'ORDER_TOTAL', 
100000, 'ASSIGNED', 
'2025-10-20 00:00:00', '2025-10-31 23:59:59', 
100, 10, 1, TRUE),
('v-user123-02', 'Voucher Shop 1 30K (RIÊNG)', 'RIENGSHOP1', 'SHOP', 'SHOP_UUID_123', 
'FIXED_AMOUNT', 30000, NULL, 'ORDER_TOTAL', 
50000, 'ASSIGNED', 
'2025-10-01 00:00:00', '2025-10-31 23:59:59', 
50, 5, 1, TRUE),
('v-user123-03', 'Freeship 20K (RIÊNG)', 'RIENGSHIP20K', 'PLATFORM', '111111111111111111111111111111111111', 
'FIXED_AMOUNT', 20000, NULL, 'SHIPPING_FEE', 
0, 'ASSIGNED', 
'2025-10-25 00:00:00', '2025-10-30 23:59:59', 
100, 0, 1, TRUE);


-- =================================================================
-- BẢNG `user_vouchers`
-- Gán 3 voucher ASSIGNED ở trên cho 'user123'
-- =================================================================

INSERT INTO `user_vouchers` (`user_id`, `voucher_id`, `status`)
VALUES
('user123', 'v-user123-01', 'AVAILABLE'),
('user123', 'v-user123-02', 'AVAILABLE'),
('user123', 'v-user123-03', 'AVAILABLE');