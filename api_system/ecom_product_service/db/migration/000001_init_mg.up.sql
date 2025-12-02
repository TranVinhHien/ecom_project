CREATE DATABASE IF NOT EXISTS ecommerce_product_db;
USE ecommerce_product_db;
SET NAMES utf8mb4;

-- 1. Bảng category 
CREATE TABLE category (
  category_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL UNIQUE,
  `key` VARCHAR(128) NOT NULL UNIQUE,
  `path` VARCHAR(500), -- Kiến trúc sư khuyên: Đổi TEXT thành VARCHAR để Index nhanh hơn
  `parent` VARCHAR(36),
  image TEXT,

  FOREIGN KEY (parent) REFERENCES category(category_id)
);
-- Sửa lỗi 1 & 2: Bỏ dấu nháy đơn, dùng VARCHAR
CREATE INDEX idx_category_path ON category(`path`);

-- 2. Bảng brand
CREATE TABLE brand (
  brand_id VARCHAR(36) PRIMARY KEY,
  name NVARCHAR(128) NOT NULL,
  code VARCHAR(15) NOT NULL UNIQUE, -- Đã tự động có Index Unique ở đây
  image TEXT,
  create_date DATETIME DEFAULT NOW(),
  update_date DATETIME DEFAULT NOW()
);
-- Đã xóa dòng CREATE INDEX idx_brand_code vì bị dư thừa

CREATE TABLE product (
    id VARCHAR(36) PRIMARY KEY,
    name TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `key` VARCHAR(500) NOT NULL UNIQUE, 
    description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    short_description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
    
    -- Các trường quan hệ và phân loại
    brand_id VARCHAR(36),
    category_id VARCHAR(36) NOT NULL,
    shop_id VARCHAR(36) NOT NULL,

    -- Hình ảnh đại diện cho SPU
    image TEXT NOT NULL,
    media TEXT, 

    -- Các trường quản lý
    delete_status ENUM('Pending','Active', 'Deleted') DEFAULT 'Active',
    product_is_permission_return BOOLEAN DEFAULT TRUE,
    product_is_permission_check BOOLEAN DEFAULT TRUE,

    -- Các trường theo dõi
    create_date DATETIME DEFAULT NOW(),
    update_date DATETIME DEFAULT NOW() ON UPDATE NOW(),
    create_by VARCHAR(128),
    update_by VARCHAR(128),
    total_sold BIGINT NOT NULL DEFAULT 0,
    
    -- Khóa ngoại
    FOREIGN KEY (brand_id) REFERENCES brand(brand_id),
    FOREIGN KEY (category_id) REFERENCES category(category_id)
);

CREATE INDEX idx_product_category_optimized 
ON product (category_id, delete_status, min_price, total_sold);

CREATE INDEX idx_product_shop_optimized 
ON product (shop_id, delete_status, min_price, total_sold);

CREATE INDEX idx_product_global_filter 
ON product (delete_status, min_price, create_date);

CREATE INDEX idx_product_createdate 
ON product (create_date);

CREATE INDEX idx_key
ON product (`key`);

CREATE FULLTEXT INDEX idx_product_name_fulltext ON product(name);
-- Bảng 3: option_value
CREATE TABLE option_value (
    id VARCHAR(36) PRIMARY KEY,
    option_name VARCHAR(100) NOT NULL, -- Kiến trúc sư khuyên: Đổi TEXT thành VARCHAR
    `value` VARCHAR(255) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    image TEXT, 
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);
-- Sửa lỗi 1 & 2: Bỏ dấu nháy đơn, đổi TEXT thành VARCHAR để index được
CREATE INDEX idx_option_value_product_id_name ON option_value(product_id, option_name);

-- Bảng 4: product_sku
CREATE TABLE product_sku (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    sku_code VARCHAR(100) NOT NULL UNIQUE, 

    price DOUBLE NOT NULL DEFAULT 0 CHECK (price >= 0),
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    quantity_reserver INT NOT NULL DEFAULT 0 CHECK (quantity_reserver >= 0),
    sku_name TEXT,
    weight DOUBLE NOT NULL DEFAULT 0 CHECK (weight >= 0), 

    create_date DATETIME DEFAULT NOW(),
    update_date DATETIME DEFAULT NOW() ON UPDATE NOW(),
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);
-- Index này OK
CREATE INDEX idx_product_sku_product_id_price ON product_sku(product_id, price);

-- Bảng 5: sku_attr
CREATE TABLE sku_attr (
    sku_id VARCHAR(36) NOT NULL,
    option_value_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (sku_id, option_value_id), 
    FOREIGN KEY (sku_id) REFERENCES product_sku(id) ON DELETE CASCADE,
    FOREIGN KEY (option_value_id) REFERENCES option_value(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);