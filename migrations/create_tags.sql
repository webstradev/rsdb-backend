CREATE TABLE `tags` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`tag` VARCHAR(50) NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `tag` (`tag`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;