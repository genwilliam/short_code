-- ============================================================
-- Short URL Service
-- Database: MySQL 8.0+
-- ============================================================

CREATE DATABASE IF NOT EXISTS short_url
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_0900_ai_ci;

USE short_url;


-- ============================================================
-- Table: short_links
--
-- 用于保存：
--   short_code -> original_url
--
-- 示例：
--   abc123 -> https://www.example.com/very/long/url
-- ============================================================

CREATE TABLE IF NOT EXISTS short_links (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT
        COMMENT 'internal id',

    short_code VARCHAR(32) NOT NULL
        COMMENT 'short code',

    original_url TEXT NOT NULL
        COMMENT 'original url',

    status TINYINT UNSIGNED NOT NULL DEFAULT 1
        COMMENT 'link status',
        -- 0=disabled,1=enabled

    expires_at DATETIME NULL
        COMMENT 'expiration time',
        -- NULL means never expires
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        COMMENT 'creation time',

    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP
        COMMENT 'last update time',

    deleted_at DATETIME NULL
        COMMENT 'soft delete time',
        -- NULL means not deleted

    PRIMARY KEY (id),

    UNIQUE KEY uk_short_code (short_code),

    KEY idx_status (status),

    KEY idx_expires_at (expires_at),

    KEY idx_created_at (created_at),

    KEY idx_deleted_at (deleted_at)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci
  COMMENT='短链接映射表';
