CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `account` varchar(100) NOT NULL COMMENT '账号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `user_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1-普通用户,2管理员,100超级管理员',
  `user_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1-可用,2-不可用,3-未激活',
  `mobile` char(11) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像地址',
  `bio` varchar(200) DEFAULT NULL COMMENT '个人说明',
  `token` varchar(100) DEFAULT NULL COMMENT '登陆token',
  `token_expire` datetime DEFAULT NULL COMMENT 'token超时时间',
  `last_login_ip` varchar(100) DEFAULT NULL COMMENT '最后登陆ip地址',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登陆时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `doc` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '文档内容',
  `doc_status` tinyint NOT NULL DEFAULT '1' COMMENT '1正常2审核中3禁用',
  `group_id` int NOT NULL DEFAULT '0' COMMENT '分组id',
  `view_count` int NOT NULL DEFAULT '0' COMMENT '查看次数',
  `like_count` int NOT NULL DEFAULT '0' COMMENT '点赞次数',
  `is_top` tinyint NOT NULL DEFAULT '2' COMMENT '1置顶2不置顶',
  `priority` int NOT NULL DEFAULT '0' COMMENT '优先级',
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `doc_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户Id',
  `title` varchar(100) NOT NULL COMMENT '组名',
  `icon` varchar(100) NOT NULL DEFAULT '' COMMENT '图标',
  `p_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `priority` int NOT NULL DEFAULT '0' COMMENT '优先级',
  `deleted_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='分组表';

CREATE TABLE `global_conf` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL,
  `value` varchar(2500) NOT NULL DEFAULT '',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `version` int unsigned NOT NULL DEFAULT '0',
  `created_by` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='全局配置';

CREATE TABLE `team_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `team_id` bigint unsigned NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='团队成员表';

CREATE TABLE `team` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '名字',
  `number` int unsigned NOT NULL DEFAULT '0' COMMENT '人数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='团队';

