CREATE TABLE IF NOT EXISTS `users` (
	`id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	`email` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`phone` VARCHAR(16) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`role` TINYINT(4) NOT NULL COMMENT 'Role user (ex. 1: Admin, 2: User Guest, etc.)',
	`password` VARCHAR(200) NULL DEFAULT NULL COMMENT 'User password hash with Bcrypt' COLLATE 'utf8mb4_general_ci',
	`created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `email_k` (`email`) USING BTREE,
	INDEX `phone_k` (`phone`) USING BTREE,
	INDEX `source_k` (`role`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=9
;
