-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_accounts" table
CREATE TABLE `new_accounts` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `account_type` text NOT NULL DEFAULT ('checking'), `is_active` bool NOT NULL DEFAULT (true), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `country_id` integer NOT NULL, `currency_id` integer NOT NULL, `user_id` integer NOT NULL, CONSTRAINT `accounts_countries_accounts` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE NO ACTION, CONSTRAINT `accounts_currencies_accounts` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`) ON DELETE NO ACTION, CONSTRAINT `accounts_users_accounts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "accounts" to new temporary table "new_accounts"
INSERT INTO `new_accounts` (`id`, `name`, `account_type`, `created_at`, `updated_at`, `country_id`, `currency_id`, `user_id`) SELECT `id`, `name`, `account_type`, `created_at`, `updated_at`, `country_id`, `currency_id`, `user_id` FROM `accounts`;
-- Drop "accounts" table after copying rows
DROP TABLE `accounts`;
-- Rename temporary table "new_accounts" to "accounts"
ALTER TABLE `new_accounts` RENAME TO `accounts`;
-- Add column "created_at" to table: "budgets"
ALTER TABLE `budgets` ADD COLUMN `created_at` datetime NOT NULL;
-- Add column "updated_at" to table: "budgets"
ALTER TABLE `budgets` ADD COLUMN `updated_at` datetime NOT NULL;
-- Create index "budget_month_year_user_id_category_id_country_id_currency_id" to table: "budgets"
CREATE UNIQUE INDEX `budget_month_year_user_id_category_id_country_id_currency_id` ON `budgets` (`month`, `year`, `user_id`, `category_id`, `country_id`, `currency_id`);
-- Create "new_categories" table
CREATE TABLE `new_categories` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `type` text NOT NULL, `color` text NOT NULL, `is_active` bool NOT NULL DEFAULT (true), `user_id` integer NOT NULL, CONSTRAINT `categories_users_categories` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "categories" to new temporary table "new_categories"
INSERT INTO `new_categories` (`id`, `name`, `type`, `color`, `user_id`) SELECT `id`, `name`, `type`, `color`, `user_id` FROM `categories`;
-- Drop "categories" table after copying rows
DROP TABLE `categories`;
-- Rename temporary table "new_categories" to "categories"
ALTER TABLE `new_categories` RENAME TO `categories`;
-- Create index "category_name_user_id" to table: "categories"
CREATE UNIQUE INDEX `category_name_user_id` ON `categories` (`name`, `user_id`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
