CREATE TABLE `platforms_articles` (
	`platform_id` INT(11) NOT NULL,
	`article_id` INT(11) NOT NULL,
	PRIMARY KEY (`platform_id`, `article_id`) USING BTREE,
	INDEX `platform_fk` (`platform_id`) USING BTREE,
	INDEX `platform_articles_article_fk` (`article_id`) USING BTREE,
	CONSTRAINT `platform_articles_article_fk` FOREIGN KEY (`article_id`) REFERENCES `restorebackup`.`articles` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT `platform_articles_platform_fk` FOREIGN KEY (`platform_id`) REFERENCES `restorebackup`.`platforms` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
