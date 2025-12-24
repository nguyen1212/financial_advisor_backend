
-- +migrate Up
CREATE TABLE IF NOT EXISTS `news`
(
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `publisher_id` BIGINT UNSIGNED NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `thumbnail` VARCHAR(500) NULL DEFAULT NULL,
  `link` VARCHAR(500) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0: unknown; 1: added; 2: synced',
  `category` TINYINT NOT NULL DEFAULT 0 COMMENT '0: unknown; 1: finance; 2: military',
  `published_at` TIMESTAMP NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_news_publisherid` FOREIGN KEY (`publisher_id`) REFERENCES `publishers`(`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS `news`;
