
-- +migrate Up
ALTER TABLE `news`
  ADD COLUMN `tmp` VARCHAR(255) NOT NULL;

-- +migrate Down
ALTER TABLE `news`
  DROP COLUMN `tmp`;
