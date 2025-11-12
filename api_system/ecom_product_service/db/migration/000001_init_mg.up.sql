CREATE DATABASE IF NOT EXISTS ecommerce_product_db;
USE ecommerce_product_db;
-- Thiết lập mã hóa UTF-8
SET NAMES utf8mb4;

-- 1. Bảng category 
CREATE TABLE category (
  category_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL UNIQUE,
  `key` VARCHAR(128) NOT NULL UNIQUE,
  `path` TEXT ,
  `parent` VARCHAR(36),
  image TEXT,

  FOREIGN KEY (parent) REFERENCES category(category_id)
);

-- 2. Bảng brand
CREATE TABLE brand (
  brand_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL,
  code VARCHAR(15) NOT NULL UNIQUE,
  image TEXT,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW()
);


CREATE TABLE product (
    id VARCHAR(36) PRIMARY KEY,
    name TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `key` VARCHAR(500) NOT NULL UNIQUE, -- Dùng cho URL thân thiện (slug)
    description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    short_description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    
    -- Các trường quan hệ và phân loại
    brand_id VARCHAR(36),
    category_id VARCHAR(36) NOT NULL,
    shop_id VARCHAR(36) NOT NULL,

    -- Hình ảnh đại diện cho SPU
    image TEXT NOT NULL,
    media TEXT, -- Lưu JSON hoặc danh sách ảnh khác

    -- Các trường quản lý của 
    delete_status ENUM('Active', 'Deleted') DEFAULT 'Active',
    product_is_permission_return BOOLEAN DEFAULT TRUE,
    product_is_permission_check BOOLEAN DEFAULT TRUE,

    -- Các trường theo dõi
    create_date DATETIME DEFAULT NOW(),
    update_date DATETIME DEFAULT NOW() ON UPDATE NOW(),
    create_by VARCHAR(128),
    update_by VARCHAR(128),

    -- Khóa ngoại
    FOREIGN KEY (brand_id) REFERENCES brand(brand_id),
    FOREIGN KEY (category_id) REFERENCES category(category_id)
);

-- Bảng 3: option_values (Thay thế cho phần `value` trong bảng `attr` ) "ssize quaanf"
-- Dùng để lưu các giá trị cụ thể, ví dụ: 'Đen', 'Trắng', 'M', 'L' 41,42
CREATE TABLE option_value (
    id VARCHAR(36) PRIMARY KEY,
    option_name TEXT NOT NULL,
    `value` VARCHAR(255) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    image TEXT, -- Hình ảnh cho giá trị (ví dụ ảnh swatch màu)
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);

-- Bảng 4: product_variants (Tương đương bảng `sku` , nhưng tối ưu hơn)
-- Lưu thông tin cho từng phiên bản sản phẩm cụ thể (SKU)
CREATE TABLE product_sku (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    sku_code VARCHAR(100) NOT NULL UNIQUE, -- Mã SKU người dùng tự đặt, dễ nhận biết

    -- Giá và tồn kho thuộc về SKU
    price DOUBLE NOT NULL DEFAULT 0 CHECK (price >= 0),
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_reserver INT NOT NULL DEFAULT 0 CHECK (quantity_reserver >= 0),
    sku_name TEXT ,
    -- cân nặng  từng sản phẩm 
    weight DOUBLE NOT NULL DEFAULT 0 CHECK (weight >= 0), -- Cân nặng (kg)

    -- Các trường theo dõi 
    create_date DATETIME DEFAULT NOW(),
    update_date DATETIME DEFAULT NOW() ON UPDATE NOW(),
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);

-- Bảng 5: variant_values (Tương đương bảng `sku_attr` , nhưng hiệu quả hơn)
-- Bảng trung gian để kết nối một SKU với các giá trị thuộc tính của nó
CREATE TABLE sku_attr (
    sku_id VARCHAR(36) NOT NULL,
    option_value_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (sku_id, option_value_id), -- Khóa chính kết hợp để đảm bảo không trùng lặp
    FOREIGN KEY (sku_id) REFERENCES product_sku(id) ON DELETE CASCADE,
    FOREIGN KEY (option_value_id) REFERENCES option_value(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);
