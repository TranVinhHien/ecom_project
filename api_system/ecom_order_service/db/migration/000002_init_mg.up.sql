

-- -- Script để chạy trong MySQL sau khi tạo bảng category.
-- -- Đã sửa để tránh lỗi "Can't update table in trigger" bằng cách dùng BEFORE trigger và set NEW.path trực tiếp (giả sử path của parent đã được tính đúng trước đó).

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

-- DELIMITER ;
