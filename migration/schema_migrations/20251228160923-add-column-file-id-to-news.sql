
-- +migrate Up
ALTER TABLE news
  ADD COLUMN file_id BIGINT UNSIGNED NULL DEFAULT NULL AFTER `file_path`;

-- +migrate Down
ALTER TABLE news
  DROP COLUMN file_id;
