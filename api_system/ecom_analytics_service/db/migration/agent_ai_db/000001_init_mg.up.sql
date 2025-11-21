CREATE DATABASE IF NOT EXISTS agent_ai_db;
USE agent_ai_db;
-- Thiết lập mã hóa UTF-8
-- =================================================================
-- DATABASE SCHEMA FOR ORDER SERVICE
-- =================================================================
SET NAMES utf8mb4;
SET time_zone = '+07:00';

CREATE TABLE `app_states` (
  `app_name` varchar(128) NOT NULL,
  `state` longtext NOT NULL,
  `update_time` datetime(6) NOT NULL,
  PRIMARY KEY (`app_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_states` (
  `app_name` varchar(128) NOT NULL,
  `user_id` varchar(128) NOT NULL,
  `state` longtext NOT NULL,
  `update_time` datetime(6) NOT NULL,
  PRIMARY KEY (`app_name`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `sessions` (
  `app_name` varchar(128) NOT NULL,
  `user_id` varchar(128) NOT NULL,
  `id` varchar(128) NOT NULL,
  `state` longtext NOT NULL,
  `create_time` datetime(6) NOT NULL,
  `update_time` datetime(6) NOT NULL,
  PRIMARY KEY (`app_name`, `user_id`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `events` (
  `id` varchar(128) NOT NULL,
  `app_name` varchar(128) NOT NULL,
  `user_id` varchar(128) NOT NULL,
  `session_id` varchar(128) NOT NULL,
  `invocation_id` varchar(256) NOT NULL,
  `author` varchar(256) NOT NULL,
  `actions` blob DEFAULT NULL,
  `long_running_tool_ids_json` text DEFAULT NULL,
  `branch` varchar(256) DEFAULT NULL,
  `timestamp` datetime(6) NOT NULL,
  `content` longtext DEFAULT NULL,
  `grounding_metadata` longtext DEFAULT NULL,
  `custom_metadata` longtext DEFAULT NULL,
  `partial` tinyint(1) DEFAULT NULL,
  `turn_complete` tinyint(1) DEFAULT NULL,
  `error_code` varchar(256) DEFAULT NULL,
  `error_message` varchar(1024) DEFAULT NULL,
  `interrupted` tinyint(1) DEFAULT NULL,
  
  -- Định nghĩa khóa chính (Primary Key) gồm 4 trường
  PRIMARY KEY (`id`, `app_name`, `user_id`, `session_id`),
  
  -- Định nghĩa Index phụ (dựa trên phần index_name 'app_name' trong ảnh)
  KEY `app_name` (`app_name`, `user_id`, `session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
