-- Create index "transaction_account_id" to table: "transactions"
CREATE INDEX `transaction_account_id` ON `transactions` (`account_id`);
-- Create index "transaction_category_id" to table: "transactions"
CREATE INDEX `transaction_category_id` ON `transactions` (`category_id`);
