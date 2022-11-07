-- my_treasure_doc.doc definition

CREATE TABLE `doc` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '文档内容',
  `doc_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1正常2审核中3禁用',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组id',
  `view_count` int(11) NOT NULL DEFAULT '0' COMMENT '查看次数',
  `like_count` int(11) NOT NULL DEFAULT '0' COMMENT '点赞次数',
  `is_top` tinyint(4) NOT NULL DEFAULT '2' COMMENT '1置顶2不置顶',
  `priority` int(255) NOT NULL DEFAULT '0' COMMENT '优先级',
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.doc_group definition

CREATE TABLE `doc_group` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户Id',
  `title` varchar(100) NOT NULL COMMENT '组名',
  `icon` varchar(100) NOT NULL DEFAULT '' COMMENT '图标',
  `p_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `priority` int(11) NOT NULL DEFAULT '0' COMMENT '优先级',
  `deleted_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分组表';


-- my_treasure_doc.global_conf definition

CREATE TABLE `global_conf` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL,
  `value` varchar(2500) NOT NULL DEFAULT '',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `version` int(10) unsigned NOT NULL DEFAULT '0',
  `created_by` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='全局配置';


-- my_treasure_doc.goods definition

CREATE TABLE `goods` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `img` varchar(2000) DEFAULT NULL,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1-可用，2-禁用',
  `goods_name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.goods_sku definition

CREATE TABLE `goods_sku` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1-可用，2-禁用',
  `goods_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品id',
  `goods_spec_ids` varchar(10) NOT NULL COMMENT '规格id',
  `price` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '价格',
  `stock` int(11) NOT NULL DEFAULT '0' COMMENT '库存',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.goods_spec definition

CREATE TABLE `goods_spec` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `good_id` int(10) unsigned NOT NULL DEFAULT '0',
  `spec` varchar(100) NOT NULL COMMENT '规格',
  `units` varchar(100) DEFAULT NULL COMMENT '单位',
  `spec_val` varchar(100) NOT NULL COMMENT '规格值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.`order` definition

CREATE TABLE `order` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(100) NOT NULL DEFAULT '' COMMENT '订单号',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `amount` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '金额',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态,0-异常,1-待支付,2-已支付,3-支付失败,4-用户取消,5-系统取消,6-订单异常',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.order_detail definition

CREATE TABLE `order_detail` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单id',
  `good_id` int(10) unsigned NOT NULL DEFAULT '0',
  `sku_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'sku id',
  `price` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '单价',
  `quantity` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4;


-- my_treasure_doc.team definition

CREATE TABLE `team` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '名字',
  `number` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '人数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队';


-- my_treasure_doc.team_user definition

CREATE TABLE `team_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `team_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='团队成员表';


-- my_treasure_doc.`user` definition

CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) NOT NULL COMMENT '''昵称''',
  `account` varchar(100) NOT NULL COMMENT '''账号''',
  `email` varchar(100) DEFAULT NULL COMMENT '''邮箱''',
  `password` varchar(100) NOT NULL COMMENT '''密码''',
  `user_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '''1-普通用户,2管理员,100超级管理员''',
  `user_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '''1-可用,2-不可用,3-未激活''',
  `mobile` char(11) DEFAULT NULL COMMENT '''手机号''',
  `avatar` varchar(500) DEFAULT NULL COMMENT '''头像地址''',
  `bio` varchar(200) DEFAULT NULL COMMENT '''个人说明''',
  `token` varchar(100) DEFAULT NULL COMMENT '''登陆token''',
  `token_expire` datetime DEFAULT NULL COMMENT '''token超时时间''',
  `last_login_ip` varchar(100) DEFAULT NULL COMMENT '''最后登陆ip地址''',
  `last_login_time` datetime DEFAULT NULL COMMENT '''最后登陆时间''',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;