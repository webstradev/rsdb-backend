CREATE TABLE `platforms` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`name` VARCHAR(128) NOT NULL COLLATE 'utf8mb4_general_ci',
	`website` VARCHAR(255) NOT NULL,
	`country` VARCHAR(50) NOT NULL,
	`notes` TEXT NOT NULL,
	`comments` TEXT NOT NULL,
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;