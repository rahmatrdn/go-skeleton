BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255),
    phone VARCHAR(16),
    name VARCHAR(255),
    role TINYINT NOT NULL COMMENT 'Role user (ex. 1: Admin, 2: User, etc.)',
    password VARCHAR(200) COMMENT 'User password hash with Bcrypt',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY email_k (email),
    KEY phone_k (phone),
    KEY source_k (role)
);

COMMIT;