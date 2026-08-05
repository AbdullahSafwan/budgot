-- Create "login_attempts" table
CREATE TABLE `login_attempts` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `username` text NULL, `ip_address` text NOT NULL, `success` bool NOT NULL, `attempted_at` datetime NOT NULL);
