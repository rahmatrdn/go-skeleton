CREATE TABLE IF NOT EXISTS `posts` (
	`id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
	`title` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`slug` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`body` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`status` TINYINT(4) NOT NULL DEFAULT 1 COMMENT 'Post status (ex. 1: Draft, 2: Published, 3: Archived)',
	`published_at` TIMESTAMP NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT current_timestamp(),
	`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `slug_k` (`slug`) USING BTREE,
	INDEX `FK_posts_users` (`user_id`) USING BTREE,
	INDEX `status_k` (`status`) USING BTREE,
	CONSTRAINT `FK_posts_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;
