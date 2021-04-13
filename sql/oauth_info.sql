CREATE TABLE `oauth_info`
(
    `id`          bigint(20) UNSIGNED                     NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_name`    varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '应用名称',
    `homepage`    varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '首页',
    `description` varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述简介',
    `callback`    varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '重定向URL',
    `app_id`      varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端标识',
    `app_secret`  varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端密钥',
    `is_delete`   tinyint(4)                              NOT NULL DEFAULT '0' COMMENT '是否删除(0,正常；1,删除；)',
    `created_by`  varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '创建人',
    `updated_by`  varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '修改人',
    `created_at`  timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `oauth_app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='oauth信息表';