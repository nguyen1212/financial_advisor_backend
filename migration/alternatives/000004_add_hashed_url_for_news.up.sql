ALTER TABLE `news`
  RENAME COLUMN `link` TO `url`,
  ADD COLUMN `hashed_url` BINARY(16) NOT NULL AFTER `url`,
  ADD UNIQUE INDEX `uidx_news_urlhash` (`hashed_url`),
  ADD COLUMN `file_path` VARCHAR(255) NULL AFTER `hashed_url`,
  ADD COLUMN `file_size` INT UNSIGNED NULL AFTER `file_path`,
  ADD COLUMN `author` VARCHAR(100) NULL AFTER `title`;
