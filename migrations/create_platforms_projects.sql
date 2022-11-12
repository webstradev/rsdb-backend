CREATE TABLE `platforms_projects` (
	`platform_id` INT(11) NOT NULL,
	`project_id` INT(11) NOT NULL,
	PRIMARY KEY (`platform_id`, `project_id`) USING BTREE,
	INDEX `platform_projects_project_fk` (`project_id`) USING BTREE,
	INDEX `platform_id` (`platform_id`) USING BTREE,
	CONSTRAINT `platform_projects_platform_fk` FOREIGN KEY (`platform_id`) REFERENCES `platforms` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT,
	CONSTRAINT `platform_projects_project_fk` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
