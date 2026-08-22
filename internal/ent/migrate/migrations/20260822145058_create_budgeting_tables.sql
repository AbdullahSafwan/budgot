-- Create "accounts" table
CREATE TABLE `accounts` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `account_type` text NOT NULL DEFAULT ('checking'), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `country_id` integer NOT NULL, `currency_id` integer NOT NULL, `user_id` integer NOT NULL, CONSTRAINT `accounts_countries_accounts` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE NO ACTION, CONSTRAINT `accounts_currencies_accounts` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`) ON DELETE NO ACTION, CONSTRAINT `accounts_users_accounts` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Create "budgets" table
CREATE TABLE `budgets` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `month` integer NOT NULL, `year` integer NOT NULL, `amount` integer NOT NULL, `category_id` integer NOT NULL, `country_id` integer NOT NULL, `currency_id` integer NOT NULL, `user_id` integer NOT NULL, CONSTRAINT `budgets_categories_budgets` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE NO ACTION, CONSTRAINT `budgets_countries_budgets` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE NO ACTION, CONSTRAINT `budgets_currencies_budgets` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`) ON DELETE NO ACTION, CONSTRAINT `budgets_users_budgets` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Create "categories" table
CREATE TABLE `categories` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `type` text NOT NULL, `color` text NOT NULL, `user_id` integer NOT NULL, CONSTRAINT `categories_users_categories` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Create "countries" table
CREATE TABLE `countries` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `name` text NOT NULL);
-- Create index "countries_code_key" to table: "countries"
CREATE UNIQUE INDEX `countries_code_key` ON `countries` (`code`);
-- Create "currencies" table
CREATE TABLE `currencies` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `name` text NOT NULL, `symbol` text NULL, `decimal_places` integer NOT NULL DEFAULT (2));
-- Create index "currencies_code_key" to table: "currencies"
CREATE UNIQUE INDEX `currencies_code_key` ON `currencies` (`code`);
-- Create "transactions" table
CREATE TABLE `transactions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `transaction_type` text NOT NULL, `amount` integer NOT NULL, `description` text NULL, `transaction_date` datetime NOT NULL, `created_at` datetime NOT NULL, `account_id` integer NOT NULL, `category_id` integer NOT NULL, `transaction_linked_transaction` integer NULL, CONSTRAINT `transactions_accounts_transactions` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE, CONSTRAINT `transactions_categories_transactions` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE NO ACTION, CONSTRAINT `transactions_transactions_linked_transaction` FOREIGN KEY (`transaction_linked_transaction`) REFERENCES `transactions` (`id`) ON DELETE SET NULL);
-- Create index "transactions_transaction_linked_transaction_key" to table: "transactions"
CREATE UNIQUE INDEX `transactions_transaction_linked_transaction_key` ON `transactions` (`transaction_linked_transaction`);
