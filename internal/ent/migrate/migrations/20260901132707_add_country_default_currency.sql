-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_countries" table
CREATE TABLE `new_countries` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `code` text NOT NULL,
  `name` text NOT NULL,
  `default_currency_id` integer NULL,
  CONSTRAINT `countries_currencies_default_currency` FOREIGN KEY (`default_currency_id`) REFERENCES `currencies` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Copy rows from old table "countries" to new temporary table "new_countries"
INSERT INTO `new_countries` (`id`, `code`, `name`) SELECT `id`, `code`, `name` FROM `countries`;
-- Drop "countries" table after copying rows
DROP TABLE `countries`;
-- Rename temporary table "new_countries" to "countries"
ALTER TABLE `new_countries` RENAME TO `countries`;
-- Create index "countries_code_key" to table: "countries"
CREATE UNIQUE INDEX `countries_code_key` ON `countries` (`code`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
