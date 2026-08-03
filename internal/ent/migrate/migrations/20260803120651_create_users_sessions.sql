-- Create "sessions" table
CREATE TABLE `sessions` (`id` text NOT NULL, `expires_at` datetime NOT NULL, `last_seen` datetime NOT NULL, `ip_address` text NOT NULL, `user_agent_hash` text NOT NULL, `created_at` datetime NOT NULL, `user_id` integer NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `sessions_users_sessions` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE);
-- Create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `username` text NOT NULL, `email` text NOT NULL, `password_hash` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `is_active` bool NOT NULL DEFAULT (true));
-- Create index "users_username_key" to table: "users"
CREATE UNIQUE INDEX `users_username_key` ON `users` (`username`);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX `users_email_key` ON `users` (`email`);
