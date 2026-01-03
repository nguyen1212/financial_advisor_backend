
-- +migrate Up
ALTER TABLE `news`
  RENAME COLUMN `file_id` TO `news_with_fulltext_id`;
  
-- +migrate Down
ALTER TABLE `news`
  RENAME COLUMN `news_with_fulltext_id` TO `file_id`;

