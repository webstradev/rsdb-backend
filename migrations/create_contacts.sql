CREATE TABLE `contacts` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`email` VARCHAR(128) NOT NULL,
	`phone` VARCHAR(50) NOT NULL,
	`phone2` VARCHAR(50) NOT NULL,
	`address` VARCHAR(255) NOT NULL,
	`notes` TEXT NOT NULL,
	`source` VARCHAR(100) NOT NULL DEFAULT '',
	`privacy` VARCHAR(100) NOT NULL DEFAULT '',
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;