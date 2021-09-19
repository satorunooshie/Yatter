CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(255),
  `avatar` text,
  `header` text,
  `note` text,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `delete_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `delete_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_account_id` (`account_id`),
  CONSTRAINT `fk_status_account_id` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`)
);

CREATE TABLE `follow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `follower_id` bigint(20) NOT NULL,
  `followee_id` bigint(20) NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `delete_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_follower_id` (`follower_id`),
  INDEX `idx_followee_id` (`followee_id`)
);

CREATE TABLE `media_attachment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `status_id` bigint(20) NOT NULL,
  `type` tinyint(3) DEFAULT 1 COMMENT '1->画像',
  `url` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `delete_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_status_id` (`status_id`),
  CONSTRAINT `fk_media_attachments_status_id` FOREIGN KEY (`status_id`) REFERENCES `status` (`id`)
);
