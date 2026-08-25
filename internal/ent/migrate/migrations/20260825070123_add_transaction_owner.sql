-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_transactions" table
CREATE TABLE `new_transactions` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `amount` integer NOT NULL,
  `description` text NULL,
  `transaction_date` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `account_id` integer NOT NULL,
  `category_id` integer NOT NULL,
  `transaction_linked_transaction` integer NULL,
  `user_id` integer NOT NULL,
  CONSTRAINT `transactions_users_transactions` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `transactions_transactions_linked_transaction` FOREIGN KEY (`transaction_linked_transaction`) REFERENCES `transactions` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT `transactions_categories_transactions` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `transactions_accounts_transactions` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Copy rows from old table "transactions" to new temporary table "new_transactions"
INSERT INTO `new_transactions` (`id`, `amount`, `description`, `transaction_date`, `created_at`, `account_id`, `category_id`, `transaction_linked_transaction`) SELECT `id`, `amount`, `description`, `transaction_date`, `created_at`, `account_id`, `category_id`, `transaction_linked_transaction` FROM `transactions`;
-- Drop "transactions" table after copying rows
DROP TABLE `transactions`;
-- Rename temporary table "new_transactions" to "transactions"
ALTER TABLE `new_transactions` RENAME TO `transactions`;
-- Create index "transactions_transaction_linked_transaction_key" to table: "transactions"
CREATE UNIQUE INDEX `transactions_transaction_linked_transaction_key` ON `transactions` (`transaction_linked_transaction`);
-- Create index "transaction_user_id" to table: "transactions"
CREATE INDEX `transaction_user_id` ON `transactions` (`user_id`);
-- Create index "transaction_account_id" to table: "transactions"
CREATE INDEX `transaction_account_id` ON `transactions` (`account_id`);
-- Create index "transaction_category_id" to table: "transactions"
CREATE INDEX `transaction_category_id` ON `transactions` (`category_id`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
