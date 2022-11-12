CREATE TABLE `contacts` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`created_at` DATETIME NOT NULL DEFAULT current_timestamp(),
	`modified_at` DATETIME NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	`deleted_at` DATETIME NULL DEFAULT NULL,
	`platform_id` INT(11) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`title` VARCHAR(512) NOT NULL DEFAULT '',
	`email` VARCHAR(128) NOT NULL,
	`phone` VARCHAR(100) NOT NULL,
	`phone2` VARCHAR(100) NOT NULL,
	`address` VARCHAR(255) NOT NULL,
	`notes` TEXT NOT NULL,
	`source` VARCHAR(100) NOT NULL DEFAULT '',
	`privacy` VARCHAR(100) NOT NULL DEFAULT '',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `contacts_platforms_Fk` (`platform_id`) USING BTREE,
	CONSTRAINT `contacts_platforms_Fk` FOREIGN KEY (`platform_id`) REFERENCES `restorebackup`.`platforms` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
