CREATE TABLE IF NOT EXISTS `todo_lists` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`user_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
	`title` VARCHAR(200) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`description` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`doing_at` DATE NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT NULL,
	`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `FK_todo_lists_users` (`user_id`) USING BTREE,
	CONSTRAINT `FK_todo_lists_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=15
;
