CREATE TABLE `user_info`
(
    `id`         bigint(20) UNSIGNED                     NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `account`    varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户账号',
    `password`   varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',
    `resource1`  varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '受保护资源1',
    `resource2`  varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '受保护资源2',
    `resource3`  varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '受保护资源3',
    `is_delete`  tinyint(4)                              NOT NULL DEFAULT '0' COMMENT '是否删除(0,正常；1,删除；)',
    `created_by` varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '创建人',
    `updated_by` varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '修改人',
    `created_at` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_account` (`account`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='账号信息表';