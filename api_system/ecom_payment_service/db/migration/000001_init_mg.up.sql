CREATE DATABASE IF NOT EXISTS ecommerce_transaction_db;
USE ecommerce_transaction_db;
-- Thiết lập mã hóa UTF-8
-- =================================================================
-- DATABASE SCHEMA FOR ORDER SERVICE
-- =================================================================
SET NAMES utf8mb4;
SET time_zone = '+07:00';
-- =================================================================
-- 1. Bảng `payment_methods`
-- =================================================================
CREATE TABLE `payment_methods` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính',
  `name` VARCHAR(100) NOT NULL COMMENT 'Tên hiển thị cho người dùng (Ví dụ: "Ví MoMo")',
  `code` VARCHAR(50) NOT NULL COMMENT 'Mã định danh (Ví dụ: "MOMO", "COD")',
  `type` ENUM('ONLINE', 'OFFLINE') NOT NULL COMMENT 'Loại hình thanh toán',
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Cho phép bật/tắt phương thức này',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_code` (`code`)
) ENGINE=InnoDB COMMENT='Quản lý các phương thức thanh toán';

-- =================================================================
-- 2. Bảng `transactions` (Giữ nguyên - Chỉ lưu dòng tiền thực tế)
-- =================================================================
CREATE TABLE `transactions` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính',
  `transaction_code` VARCHAR(50) NOT NULL COMMENT 'Mã giao dịch thân thiện, duy nhất',
  `order_id` CHAR(36) DEFAULT NULL COMMENT 'ID đơn hàng tổng. NULL nếu là Payout/Deposit',
  `payment_method_id` CHAR(36) NOT NULL COMMENT 'FK tới payment_methods',
  `amount` DECIMAL(15, 2) NOT NULL COMMENT 'Số tiền thực tế khách trả HOẶC thực tế chuyển/nhận',
  `currency` VARCHAR(3) NOT NULL DEFAULT 'VND',
  `type` ENUM('PAYMENT', 'REFUND', 'PAYOUT', 'DEPOSIT') NOT NULL,
  `status` ENUM('PENDING', 'SUCCESS', 'FAILED') NOT NULL,
  `gateway_transaction_id` VARCHAR(100) DEFAULT NULL,
  `notes` TEXT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `processed_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_transaction_code` (`transaction_code`),
  CONSTRAINT `fk_transactions_payment_method` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_methods` (`id`)
) ENGINE=InnoDB COMMENT='Ghi lại nhật ký các giao dịch tiền tệ thực tế';

-- =================================================================
-- 3. Bảng `account_ledgers` (Giữ nguyên)
-- =================================================================
CREATE TABLE `account_ledgers` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của ví',
  `owner_id` CHAR(36) NOT NULL COMMENT 'ID của chủ sở hữu (shop_id hoặc Sàn)',
  `owner_type` ENUM('SHOP', 'PLATFORM') NOT NULL,
  `balance` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Số dư khả dụng',
  `pending_balance` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Số dư đang chờ quyết toán',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_owner` (`owner_id`, `owner_type`)
) ENGINE=InnoDB COMMENT='Quản lý ví và số dư của Sàn và các Shop';

-- =================================================================
-- 4. Bảng `ledger_entries` (Giữ nguyên)
-- =================================================================
CREATE TABLE `ledger_entries` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `ledger_id` CHAR(36) NOT NULL COMMENT 'FK tới ví/sổ cái',
  `transaction_id` CHAR(36) NOT NULL COMMENT 'FK tới giao dịch gốc',
  `amount` DECIMAL(15, 2) NOT NULL COMMENT 'Số tiền thay đổi (âm/dương)',
  `type` ENUM('CREDIT', 'DEBIT') NOT NULL COMMENT 'CREDIT (cộng tiền), DEBIT (trừ tiền)',
  `description` VARCHAR(255) NOT NULL COMMENT 'Mô tả bút toán',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_ledger_id` (`ledger_id`),
  KEY `idx_transaction_id` (`transaction_id`),
  CONSTRAINT `fk_entries_ledger` FOREIGN KEY (`ledger_id`) REFERENCES `account_ledgers` (`id`),
  CONSTRAINT `fk_entries_transaction` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`id`)
) ENGINE=InnoDB COMMENT='Sổ cái kế toán, ghi lại mọi biến động số dư';

-- =================================================================
-- 5. Bảng `shop_order_settlements` (CẬP NHẬT - Loại bỏ trường Sàn)
-- =================================================================
CREATE TABLE `shop_order_settlements` (
  `id` CHAR(36) NOT NULL,
  `shop_order_id` CHAR(36) NOT NULL COMMENT 'ID đơn hàng shop từ Order Service',
  `order_transaction_id` CHAR(36) NOT NULL COMMENT 'ID giao dịch PAYMENT gốc',
  `status` ENUM('PENDING_SETTLEMENT', 'FUNDS_HELD', 'SETTLED', 'FAILED') NOT NULL,

  -- === Financial Snapshot (Ảnh hưởng tới Shop) ===
  `order_subtotal` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng giá trị hàng hóa GỐC',
  `shop_funded_product_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Giảm giá SP do Shop chịu',
  `site_funded_product_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Giảm giá SP do Sàn trợ giá (Sàn bù cho Shop)',
  `shop_voucher_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Voucher của Shop (Shop chịu)',
  `shop_shipping_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tiền Shop hỗ trợ ship (Shop chịu)',
  `site_order_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Số tiền giảm từ voucher SÀN (tiền hàng) đã được PHÂN BỔ cho đơn hàng shop này',
  `site_shipping_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tiền Sàn hỗ trợ ship (voucher ship) đã được PHÂN BỔ cho đơn hàng shop này',
  `shipping_fee` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Phí vận chuyển thực tế khách trả cho gói hàng này (Doanh thu của Sàn)',

  -- === Kết quả Quyết toán ===
  `commission_fee` DECIMAL(15, 2) NOT NULL COMMENT 'Phí hoa hồng Sàn thu (trên giá gốc)',
  `net_settled_amount` DECIMAL(15, 2) NOT NULL COMMENT 'Số tiền thực tế Sàn chuyển vào ví Shop',

  `order_completed_at` TIMESTAMP NULL DEFAULT NULL,
  `settled_at` TIMESTAMP NULL DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_shop_order_id` (`shop_order_id`),
  CONSTRAINT `fk_settlements_transaction` FOREIGN KEY (`order_transaction_id`) REFERENCES `transactions` (`id`)
) ENGINE=InnoDB COMMENT='Theo dõi quyết toán và snapshot tài chính cho từng đơn hàng shop';

-- =================================================================
-- 6. Bảng `order_platform_costs` (Bảng MỚI - Lưu chi phí Sàn)
-- =================================================================
CREATE TABLE `order_platform_costs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` CHAR(36) NOT NULL COMMENT 'ID của đơn hàng tổng',
  `payment_transaction_id` CHAR(36) NOT NULL COMMENT 'ID của giao dịch PAYMENT gốc',
  `site_order_voucher_discount_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tổng tiền voucher Sàn giảm giá ĐƠN HÀNG (Sàn chịu)',
  `site_promotion_discount_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tổng tiền giảm giá từ ĐỢT KHUYẾN MÃI của Sàn (Sàn chịu)',
  `site_shipping_discount_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tổng tiền Sàn giảm giá SHIP (bao gồm voucher và KM ship)',
  `total_site_funded_product_discount` DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT 'Tổng tiền Sàn TRỢ GIÁ SẢN PHẨM cho toàn bộ đơn hàng',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_payment_transaction_id` (`payment_transaction_id`) COMMENT 'Mỗi giao dịch PAYMENT chỉ có 1 bản ghi chi phí Sàn',
  KEY `idx_order_id` (`order_id`),
  CONSTRAINT `fk_platform_costs_transaction` FOREIGN KEY (`payment_transaction_id`) REFERENCES `transactions` (`id`)
) ENGINE=InnoDB COMMENT='Ghi lại chi tiết các khoản chi phí Sàn chịu cho mỗi đơn hàng';
-- =================================================================
-- Thêm dữ liệu mẫu
-- =================================================================
INSERT INTO `payment_methods` 
  (`id`, `name`, `code`, `type`, `is_active`) 
VALUES
  ('a1b2c3d4-e5f6-7890-1234-567890abcdef', 'Ví điện tử MoMo', 'MOMO', 'ONLINE', TRUE),
  ('b2c3d4e5-f6a7-8901-2345-67890abcdef1', 'Thanh toán khi nhận hàng (COD)', 'COD', 'OFFLINE', TRUE),
  ('c3d4e5f6-a7b8-9012-3456-7890abcdef2', 'Cổng thanh toán VNPAY', 'VNPAY', 'ONLINE', FALSE),
  ('d4e5f6a7-b8c9-0123-4567-890abcdef3', 'Chuyển khoản ngân hàng', 'BANK_TRANSFER', 'OFFLINE', FALSE);

-- INSERT INTO `account_ledgers` 
--   (`id`, `owner_id`, `owner_type`, `balance`, `pending_balance`, `created_at`, `updated_at`)
-- VALUES
--   ('111111111111111111111111111111111111', '111111111111111111111111111111111111', 'PLATFORM', 0, 0.00, NOW(), NOW());