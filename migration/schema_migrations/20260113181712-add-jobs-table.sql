
-- +migrate Up
CREATE TABLE IF NOT EXISTS `jobs` (
    `uuid` VARCHAR(36) PRIMARY KEY,
    `payload` JSON NOT NULL,
    `result` JSON NULL DEFAULT NULL,
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0: new; 1: processing; 2: completed; 3: failed',
    `type` TINYINT NOT NULL COMMENT '0: unknown; 1: web_scrapper',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS `jobs`;
