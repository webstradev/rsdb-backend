CREATE TABLE `projects_tags` (
	`project_id` INT(11) NOT NULL,
	`tag_id` INT(11) NOT NULL,
	PRIMARY KEY (`project_id`, `tag_id`) USING BTREE,
	INDEX `tag_id` (`tag_id`) USING BTREE,
	INDEX `project_id` (`project_id`) USING BTREE,
	CONSTRAINT `projects_tags_tag_fk` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT `projects_tags_project_fk` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
