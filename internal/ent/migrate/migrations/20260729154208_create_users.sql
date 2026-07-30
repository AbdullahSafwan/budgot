-- Create "sessions" table
CREATE TABLE `sessions` (`id` bigint NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `username` varchar(255) NOT NULL, `email` varchar(255) NOT NULL, `password_hash` varchar(255) NOT NULL, `created_at` timestamp NOT NULL, `updated_at` timestamp NOT NULL, `is_active` bool NOT NULL DEFAULT 1, PRIMARY KEY (`id`), UNIQUE INDEX `email` (`email`), UNIQUE INDEX `username` (`username`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
