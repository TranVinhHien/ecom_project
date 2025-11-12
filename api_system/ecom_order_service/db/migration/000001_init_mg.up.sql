CREATE DATABASE IF NOT EXISTS ecommerce_order_db;
USE ecommerce_order_db;
-- Thiết lập mã hóa UTF-8
-- =================================================================
-- DATABASE SCHEMA FOR ORDER SERVICE
-- =================================================================
SET NAMES utf8mb4;
SET time_zone = '+07:00';

CREATE TABLE `orders` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của đơn hàng tổng',
  `order_code` VARCHAR(50) NOT NULL COMMENT 'Mã đơn hàng duy nhất, thân thiện với người dùng (ví dụ: YAN20251013ABC)',
  `user_id` CHAR(36) NOT NULL COMMENT 'UUID của người dùng đặt hàng',  
  -- Tổng hợp tài chính toàn bộ đơn hàng
  `grand_total` DECIMAL(15, 2) NOT NULL COMMENT 'Số tiền cuối cùng khách hàng thực trả sau khi áp dụng tất cả giảm giá và phí ship',
  `subtotal` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng tiền hàng gốc của tất cả sản phẩm từ tất cả các shop',
  `total_shipping_fee` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng phí ship của tất cả các shop (chưa trừ voucher ship)',
  `total_discount` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng tất cả các khoản giảm giá (từ voucher shop, sàn và ship)',

  -- Chi tiết voucher toàn sàn
  `site_order_voucher_code` VARCHAR(50) DEFAULT NULL COMMENT 'Mã voucher của sàn áp dụng cho tiền hàng',
  `site_order_voucher_discount` DECIMAL(15, 2) DEFAULT 0.00 COMMENT 'Số tiền giảm từ voucher tiền hàng của sàn',
  `site_shipping_voucher_code` VARCHAR(50) DEFAULT NULL COMMENT 'Mã voucher của sàn áp dụng cho phí ship',
  `site_shipping_voucher_discount` DECIMAL(15, 2) DEFAULT 0.00 COMMENT 'Số tiền giảm từ voucher phí ship của sàn',

  -- Dữ liệu "đóng băng" tại thời điểm đặt hàng
  `shipping_address_snapshot` JSON NOT NULL COMMENT 'Bản ghi nhanh (snapshot) địa chỉ giao hàng dạng JSON',
  `payment_method_snapshot` JSON NOT NULL COMMENT 'Bản ghi nhanh (snapshot) phương thức thanh toán dạng JSON',
  `note` TEXT COMMENT 'Ghi chú của khách hàng cho toàn bộ đơn hàng',
  
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Thời gian tạo đơn hàng',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Thời gian cập nhật lần cuối',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_order_code` (`order_code`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB COMMENT='Bảng chứa các đơn hàng tổng của khách hàng (một lần checkout)';

CREATE TABLE `shop_orders` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của đơn hàng shop',
  `shop_order_code` VARCHAR(50) NOT NULL COMMENT 'Mã đơn hàng duy nhất, thân thiện với shop (ví dụ: SHOP-ABC-123)',
  `order_id` CHAR(36) NOT NULL COMMENT 'Khóa ngoại, liên kết tới đơn hàng tổng trong bảng orders',
  `shop_id` CHAR(36) NOT NULL COMMENT 'UUID của shop/nhà bán hàng',
  `status` ENUM('AWAITING_PAYMENT', 'PROCESSING', 'SHIPPED', 'COMPLETED', 'CANCELLED', 'REFUNDED') NOT NULL COMMENT 'Trạng thái xử lý CHI TIẾT của gói hàng này. Đây là "Source of Truth"',

  -- Tài chính chi tiết cho shop này
  `subtotal` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng tiền hàng của các sản phẩm thuộc shop này',
  `total_discount` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng giảm giá của shop này (từ voucher shop + khuyến mãi sản phẩm)',
  `total_amount` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng tiền cuối cùng của shop này (subtotal + shipping_fee - total_discount)',

  -- Chi tiết voucher của shop
  `shop_voucher_code` VARCHAR(50) DEFAULT NULL COMMENT 'Mã voucher của shop đã áp dụng',
  `shop_voucher_discount` DECIMAL(15, 2) DEFAULT 0.00 COMMENT 'Số tiền giảm từ voucher của shop',

  -- Thông tin vận chuyển
  `shipping_fee`DECIMAL(15, 2) NOT NULL COMMENT 'Tổng phí ship của  shop (chưa trừ voucher ship)',
  `shipping_method` VARCHAR(100) DEFAULT NULL COMMENT 'Tên đơn vị vận chuyển',
  `tracking_code` VARCHAR(100) DEFAULT NULL COMMENT 'Mã vận đơn của đơn vị vận chuyển',
  `cancellation_reason` TEXT DEFAULT NULL COMMENT 'Lý do hủy đơn hàng (nếu có)',

  -- Các mốc thời gian quan trọng của gói hàng
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `paid_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Thời điểm được xác nhận đã thanh toán',
  `processing_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Thời điểm shop bắt đầu xử lý hàng',
  `shipped_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Thời điểm shop bàn giao cho ĐVVC',
  `completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Thời điểm giao hàng thành công',
  `cancelled_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Thời điểm đơn hàng bị hủy',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_shop_order_code` (`shop_order_code`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_shop_id` (`shop_id`),
  CONSTRAINT `fk_shop_orders_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB COMMENT='Bảng chứa các đơn hàng chi tiết của từng shop, là đơn vị vận hành chính';

CREATE TABLE `order_items` (
  `id` CHAR(36) NOT NULL COMMENT 'UUID, Khóa chính của item',
  `shop_order_id` CHAR(36) NOT NULL COMMENT 'Khóa ngoại, liên kết tới đơn hàng của shop',
  `product_id` CHAR(36) NOT NULL COMMENT 'UUID của sản phẩm gốc',
  `sku_id` CHAR(36) NOT NULL COMMENT 'UUID của biến thể sản phẩm (SKU) gốc',
  `quantity` INT UNSIGNED NOT NULL COMMENT 'Số lượng sản phẩm được mua',

  -- Chi tiết giá và khuyến mãi tại thời điểm mua
  `original_unit_price` DECIMAL(15, 2) NOT NULL COMMENT 'Giá gốc của 1 sản phẩm (trước khi giảm)',
  `final_unit_price` DECIMAL(15, 2) NOT NULL COMMENT 'Giá bán cuối cùng của 1 sản phẩm (sau khi giảm)',
  `total_price` DECIMAL(15, 2) NOT NULL COMMENT 'Tổng tiền cho item này (final_unit_price * quantity)',
  `promotions_snapshot` JSON DEFAULT NULL COMMENT 'Bản ghi nhanh các chương trình khuyến mãi đã áp dụng dạng JSON',

  -- Dữ liệu "đóng băng" của sản phẩm
  `product_name_snapshot` TEXT NOT NULL COMMENT 'Tên sản phẩm tại thời điểm mua',
  `product_image_snapshot` TEXT DEFAULT NULL COMMENT 'URL hình ảnh sản phẩm tại thời điểm mua',
  `sku_attributes_snapshot` TEXT DEFAULT NULL COMMENT 'Các thuộc tính của SKU (Màu, Size...) tại thời điểm mua',
  PRIMARY KEY (`id`),
  KEY `idx_shop_order_id` (`shop_order_id`),
  CONSTRAINT `fk_order_items_shop_order` FOREIGN KEY (`shop_order_id`) REFERENCES `shop_orders` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB COMMENT='Bảng chứa các sản phẩm chi tiết trong một đơn hàng của shop';
