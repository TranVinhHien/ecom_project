
-- -- Tạo function để tính path (nếu cần cho logic khác, nhưng trigger không dùng nữa)
-- DELIMITER //

-- CREATE FUNCTION get_category_path(cat_id VARCHAR(255)) RETURNS VARCHAR(1000)
-- DETERMINISTIC
-- READS SQL DATA
-- BEGIN
--     DECLARE result VARCHAR(1000) DEFAULT '';
--     DECLARE current_id VARCHAR(255) DEFAULT cat_id;
--     DECLARE temp_key VARCHAR(255);
--     DECLARE temp_parent VARCHAR(255);
    
--     WHILE current_id IS NOT NULL DO
--         SELECT `key`, parent INTO temp_key, temp_parent
--         FROM category WHERE category_id = current_id;
        
--         IF result = '' THEN
--             SET result = temp_key;
--         ELSE
--             SET result = CONCAT(temp_key, '/', result);
--         END IF;
        
--         SET current_id = temp_parent;
--     END WHILE;
    
--     RETURN CONCAT('/', result);
-- END //

-- DELIMITER ;

-- -- Trigger BEFORE INSERT: Tự động set path dựa trên parent (path của parent phải đã đúng)
-- DELIMITER //

-- CREATE TRIGGER before_insert_category
-- BEFORE INSERT ON category
-- FOR EACH ROW
-- BEGIN
--     IF NEW.parent IS NULL THEN
--         SET NEW.path = CONCAT('/', NEW.`key`);
--     ELSE
--         SET NEW.path = CONCAT((SELECT path FROM category WHERE category_id = NEW.parent), '/', NEW.`key`);
--     END IF;
-- END //

-- DELIMITER ;

-- -- Trigger BEFORE UPDATE: Update path nếu parent hoặc key thay đổi
-- DELIMITER //

-- CREATE TRIGGER before_update_category
-- BEFORE UPDATE ON category
-- FOR EACH ROW
-- BEGIN
--     -- Chỉ chạy logic nếu `key` hoặc `parent` thực sự thay đổi
--     -- Dùng `IFNULL` để so sánh NULL an toàn
--     IF NOT (NEW.`key` <=> OLD.`key` AND NEW.`parent` <=> OLD.`parent`) THEN
--         -- Nếu không có parent, path là /key
--         IF NEW.parent IS NULL THEN
--             SET NEW.path = CONCAT('/', NEW.`key`);
--         -- Nếu có parent, path là path_cua_cha/key
--         ELSE
--             -- Lấy path của cha và ghép với key mới
--             SET NEW.path = CONCAT((SELECT path FROM category WHERE category_id = NEW.parent), '/', NEW.`key`);
--         END IF;
--     END IF;
-- END //

-- DELIMITER ;\

-- DELIMITER $$

-- CREATE PROCEDURE generate_sku_name(IN p_sku_id VARCHAR(36))
-- BEGIN
--     DECLARE v_sku_name VARCHAR(500);
    
--     -- Tạo sku_name bằng cách nối các option_name: value
--     SELECT GROUP_CONCAT(
--         CONCAT(ov.option_name, ': ', ov.value) 
--         ORDER BY ov.option_name 
--         SEPARATOR ', '
--     ) INTO v_sku_name
--     FROM sku_attr sa
--     INNER JOIN option_value ov ON sa.option_value_id = ov.id
--     WHERE sa.sku_id = p_sku_id;
    
--     -- Cập nhật sku_name vào bảng product_sku
--     UPDATE product_sku 
--     SET sku_name = v_sku_name,
--         update_date = NOW()
--     WHERE id = p_sku_id;
-- END$$

-- DELIMITER ;

-- DELIMITER $$

-- CREATE TRIGGER after_sku_attr_insert
-- AFTER INSERT ON sku_attr
-- FOR EACH ROW
-- BEGIN
--     CALL generate_sku_name(NEW.sku_id);
-- END$$

-- DELIMITER ;



-- DELIMITER $$

-- CREATE TRIGGER after_sku_attr_update
-- AFTER UPDATE ON sku_attr
-- FOR EACH ROW
-- BEGIN
--     CALL generate_sku_name(NEW.sku_id);
    
--     -- Nếu sku_id thay đổi, cập nhật cả sku_id cũ
--     IF OLD.sku_id != NEW.sku_id THEN
--         CALL generate_sku_name(OLD.sku_id);
--     END IF;
-- END$$

-- DELIMITER ;


-- DELIMITER $$

-- CREATE TRIGGER after_sku_attr_delete
-- AFTER DELETE ON sku_attr
-- FOR EACH ROW
-- BEGIN
--     CALL generate_sku_name(OLD.sku_id);
-- END$$

-- DELIMITER ;




-- DELIMITER $$

-- CREATE TRIGGER after_option_value_update
-- AFTER UPDATE ON option_value
-- FOR EACH ROW
-- BEGIN
--     -- Cập nhật tất cả SKU liên quan đến option_value này
--     DECLARE done INT DEFAULT FALSE;
--     DECLARE v_sku_id VARCHAR(36);
--     DECLARE sku_cursor CURSOR FOR 
--         SELECT DISTINCT sku_id 
--         FROM sku_attr 
--         WHERE option_value_id = NEW.id;
--     DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
--     OPEN sku_cursor;
    
--     read_loop: LOOP
--         FETCH sku_cursor INTO v_sku_id;
--         IF done THEN
--             LEAVE read_loop;
--         END IF;
        
--         CALL generate_sku_name(v_sku_id);
--     END LOOP;
    
--     CLOSE sku_cursor;
-- END$$

-- DELIMITER ;

-- UPDATE product_sku ps
-- SET sku_name = (
--     SELECT GROUP_CONCAT(
--         CONCAT(ov.option_name, ': ', ov.value) 
--         ORDER BY ov.option_name 
--         SEPARATOR ', '
--     )
--     FROM sku_attr sa
--     INNER JOIN option_value ov ON sa.option_value_id = ov.id
--     WHERE sa.sku_id = ps.id
-- );
