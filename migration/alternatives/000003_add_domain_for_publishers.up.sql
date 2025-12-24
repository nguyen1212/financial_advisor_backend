ALTER TABLE `publishers` 
  ADD COLUMN `domain` VARCHAR(255) NOT NULL AFTER `description`,
  ADD UNIQUE INDEX `uidx_publishers_domain` (`domain`);
