CREATE TABLE `users_tokens` (
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`hashed_token` VARCHAR(255) NOT NULL,
	`type` VARCHAR(50) NOT NULL DEFAULT '',
	`created_by` INT(11) NULL DEFAULT NULL,
	`user_id` INT(11) NULL DEFAULT NULL,
	`used` TINYINT(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`hashed_token`) USING BTREE,
	INDEX `FK_users_tokens_users` (`created_by`) USING BTREE,
	INDEX `FK_userid_users` (`user_id`) USING BTREE,
	CONSTRAINT `FK_userid_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
	CONSTRAINT `FK_users_tokens_users` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
)
COLLATE='utf8mb4_bin'
ENGINE=InnoDB
;
