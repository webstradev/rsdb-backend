CREATE TABLE `articles_tags` (
	`article_id` INT(11) NOT NULL,
	`tag_id` INT(11) NOT NULL,
	PRIMARY KEY (`article_id`, `tag_id`) USING BTREE,
	INDEX `tag_id` (`tag_id`) USING BTREE,
	INDEX `article_id` (`article_id`) USING BTREE,
	CONSTRAINT `articles_tags_tag_fk` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT `articles_tags_article_fk` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
