ALTER TABLE `news`
  DROP COLUMN `author`,
  DROP COLUMN `file_size`,
  DROP COLUMN `file_path`,
  DROP INDEX `uidx_news_urlhash`,
  DROP COLUMN `hashed_url`,
  RENAME COLUMN `url` TO `link`;

