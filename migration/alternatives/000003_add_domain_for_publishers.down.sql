-- +migrate Down
ALTER TABLE `publishers` DROP COLUMN `domain`;
