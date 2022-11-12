CREATE TABLE `platforms_categories` (
	`platform_id` INT(11) NOT NULL,
	`category_id` INT(11) NOT NULL,
	PRIMARY KEY (`platform_id`, `category_id`) USING BTREE,
	INDEX `category_id` (`category_id`) USING BTREE,
	INDEX `platform_id` (`platform_id`) USING BTREE,
	CONSTRAINT `platforms_categories_category_id` FOREIGN KEY (`category_id`) REFERENCES `restorebackup`.`categories` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT `platforms_categories_platform_fk` FOREIGN KEY (`platform_id`) REFERENCES `restorebackup`.`platforms` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
