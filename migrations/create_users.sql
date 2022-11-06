CREATE TABLE `users` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`email` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci',
	`password` VARCHAR(512) NOT NULL,
	`role` VARCHAR(50) NOT NULL DEFAULT 'user',
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;