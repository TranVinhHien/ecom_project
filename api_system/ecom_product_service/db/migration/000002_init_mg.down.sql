DROP TRIGGER IF EXISTS before_insert_category;
DROP TRIGGER IF EXISTS before_update_category;
DROP FUNCTION IF EXISTS get_category_path;


DROP TRIGGER IF EXISTS after_option_value_update;
DROP TRIGGER IF EXISTS after_sku_attr_delete;
DROP TRIGGER IF EXISTS after_sku_attr_update;
DROP TRIGGER IF EXISTS after_sku_attr_insert;
DROP PROCEDURE IF EXISTS generate_sku_name;

ALTER TABLE product_sku DROP COLUMN sku_name;