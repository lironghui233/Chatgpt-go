CREATE TABLE `chat_records` (
 `id` bigint NOT NULL AUTO_INCREMENT,
 `account` varchar(255) NOT NULL DEFAULT '',
 `group_id` varchar(255) NOT NULL DEFAULT '',
 `user_msg` text,
 `user_msg_tokens` int NOT NULL DEFAULT '0',
 `user_msg_keywords` varchar(1024) NOT NULL DEFAULT '',
 `ai_msg` text,
 `ai_msg_tokens` int NOT NULL DEFAULT '0',
 `req_tokens` int NOT NULL DEFAULT '0',
 `create_at` bigint NOT NULL DEFAULT '0',
 `enterprise_id` varchar(255) NOT NULL DEFAULT '',
 `endpoint` int NOT NULL  DEFAULT '0',
 `endpoint_account` varchar(255) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`),
 KEY `index_create_at` (`create_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
