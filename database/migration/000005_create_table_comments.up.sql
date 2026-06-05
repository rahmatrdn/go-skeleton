CREATE TABLE IF NOT EXISTS `comments` (
	`id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`post_id` BIGINT(20) UNSIGNED NOT NULL,
	`user_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
	`parent_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'NULL for top-level comments, parent comment id for replies',
	`body` TEXT NOT NULL COLLATE 'utf8mb4_general_ci',
	`created_at` TIMESTAMP NULL DEFAULT current_timestamp(),
	`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `FK_comments_posts` (`post_id`) USING BTREE,
	INDEX `FK_comments_users` (`user_id`) USING BTREE,
	INDEX `FK_comments_parent` (`parent_id`) USING BTREE,
	CONSTRAINT `FK_comments_posts` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT `FK_comments_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT `FK_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;
