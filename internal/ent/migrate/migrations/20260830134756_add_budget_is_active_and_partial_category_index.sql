-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Add column "is_active" to table: "budgets"
ALTER TABLE `budgets` ADD COLUMN `is_active` bool NOT NULL DEFAULT true;
-- Create "new_categories" table
CREATE TABLE `new_categories` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `type` text NOT NULL,
  `color` text NOT NULL,
  `is_active` bool NOT NULL DEFAULT true,
  `user_id` integer NOT NULL,
  CONSTRAINT `categories_users_categories` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Copy rows from old table "categories" to new temporary table "new_categories"
INSERT INTO `new_categories` (`id`, `name`, `type`, `color`, `is_active`, `user_id`) SELECT `id`, `name`, `type`, `color`, `is_active`, `user_id` FROM `categories`;
-- Drop "categories" table after copying rows
DROP TABLE `categories`;
-- Rename temporary table "new_categories" to "categories"
ALTER TABLE `new_categories` RENAME TO `categories`;
-- Create index "category_name_user_id" to table: "categories"
CREATE UNIQUE INDEX `category_name_user_id` ON `categories` (`name`, `user_id`) WHERE is_active;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
