-- MySQL dump 10.13  Distrib 8.0.31, for Win64 (x86_64)
--
-- Host: localhost    Database: treasure_doc
-- ------------------------------------------------------
-- Server version	8.0.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `td_doc`
--

DROP TABLE IF EXISTS `td_doc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_doc` (
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
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_doc`
--

LOCK TABLES `td_doc` WRITE;
/*!40000 ALTER TABLE `td_doc` DISABLE KEYS */;
INSERT INTO `td_doc` VALUES (6,2,'术多记无起东政','农应值造线集思干类劳先又联万。志按无都正开际问事次工线了七教革。影是较现高工相四把表该习具。议复热再近习备极总或始革飞。志容志广越通向马院色义养根业。农性将业政美社气非美与儿公样样相革。较每农法我到方越指报我认认子。',1,0,0,0,1,0,NULL,'2023-05-15 17:47:55','2023-05-15 17:47:55'),(7,2,'还龙厂节除','格资号委他影属都反切目照水般没你选。学得除这知认观术三及所目图外指。可都行际治如群据想质东即工带得。资群列群向元包真族党列增观山题量。',1,0,0,0,2,0,NULL,'2023-05-15 17:48:05','2023-05-15 17:48:05'),(8,2,'义根装再性主相','段部记题者共力共二心次产。细思感深其品心二积带律矿亲。声大候接习组下事革素深率下白油二西。酸便手场八着上状九强与查示人题信省。重图月并装报法切去间权正号备克关。是容马些业报拉单满法理省直进。',1,0,0,0,2,0,NULL,'2023-05-15 17:48:08','2023-05-15 17:48:08'),(9,2,'义点记连记那响','采话主持放写整广等达劳划易。会命关气级须进此府科院领中。西节先志无王难约济便他书转式。离位须动及走名易式权属采流没命影。特劳动名并铁特毛量直会条和数精。',1,0,0,0,1,0,NULL,'2023-05-15 17:48:11','2023-05-15 17:48:11'),(10,2,'测试写一条数据1','测试写一条数据1\n',1,0,0,0,0,0,NULL,'2023-05-15 18:00:17','2023-05-15 18:00:17'),(11,2,'test','sdfdsfs\n',1,0,0,0,0,0,NULL,'2023-05-16 18:10:51','2023-05-16 18:10:51'),(12,2,'test1','test2\n',1,0,0,0,0,0,NULL,'2023-05-16 18:13:49','2023-05-16 18:13:49'),(13,2,'test2','test3\n',1,0,0,0,0,0,NULL,'2023-05-16 18:14:58','2023-05-16 18:14:58'),(14,2,'哈哈哈哈标题','sdfsf\n',1,0,0,0,0,0,NULL,'2023-05-16 18:20:39','2023-05-16 18:20:39'),(15,2,'青无头认节素理','十叫声何好代据前林从布见指八争共近。林风土响参除传数部劳识响队广。林研接量本确每实之条者太关论据技。间第中切工王相观数工动满根。感点加现山拉斯布认在老争即共标。',1,0,0,0,1,0,NULL,'2023-05-17 10:58:52','2023-05-17 10:58:52'),(16,2,'今日记录','<p>今天的天气还可以,但是今天还是要上班！</p><p><br></p><p>今天是周五，没有办法。 那明天是周六，结果还是要上班？</p>',1,0,0,0,0,0,NULL,'2023-05-27 10:52:16','2023-05-19 16:10:28');
/*!40000 ALTER TABLE `td_doc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_doc_group`
--

DROP TABLE IF EXISTS `td_doc_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_doc_group` (
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
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_doc_group`
--

LOCK TABLES `td_doc_group` WRITE;
/*!40000 ALTER TABLE `td_doc_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_doc_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_global_conf`
--

DROP TABLE IF EXISTS `td_global_conf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_global_conf` (
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
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_global_conf`
--

LOCK TABLES `td_global_conf` WRITE;
/*!40000 ALTER TABLE `td_global_conf` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_global_conf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_goods`
--

DROP TABLE IF EXISTS `td_goods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_goods` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `img` varchar(2000) DEFAULT NULL,
  `enabled` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1-可用，2-禁用',
  `goods_name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_goods`
--

LOCK TABLES `td_goods` WRITE;
/*!40000 ALTER TABLE `td_goods` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_goods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_goods_sku`
--

DROP TABLE IF EXISTS `td_goods_sku`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_goods_sku` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `enabled` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1-可用，2-禁用',
  `goods_id` int unsigned NOT NULL DEFAULT '0' COMMENT '商品id',
  `goods_spec_ids` varchar(10) NOT NULL COMMENT '规格id',
  `price` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '价格',
  `stock` int NOT NULL DEFAULT '0' COMMENT '库存',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_goods_sku`
--

LOCK TABLES `td_goods_sku` WRITE;
/*!40000 ALTER TABLE `td_goods_sku` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_goods_sku` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_goods_spec`
--

DROP TABLE IF EXISTS `td_goods_spec`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_goods_spec` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `good_id` int unsigned NOT NULL DEFAULT '0',
  `spec` varchar(100) NOT NULL COMMENT '规格',
  `units` varchar(100) DEFAULT NULL COMMENT '单位',
  `spec_val` varchar(100) NOT NULL COMMENT '规格值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_goods_spec`
--

LOCK TABLES `td_goods_spec` WRITE;
/*!40000 ALTER TABLE `td_goods_spec` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_goods_spec` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_order`
--

DROP TABLE IF EXISTS `td_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_order` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(100) NOT NULL DEFAULT '' COMMENT '订单号',
  `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `amount` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '金额',
  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态,0-异常,1-待支付,2-已支付,3-支付失败,4-用户取消,5-系统取消,6-订单异常',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_order`
--

LOCK TABLES `td_order` WRITE;
/*!40000 ALTER TABLE `td_order` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_order_detail`
--

DROP TABLE IF EXISTS `td_order_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_order_detail` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int unsigned NOT NULL DEFAULT '0' COMMENT '订单id',
  `good_id` int unsigned NOT NULL DEFAULT '0',
  `sku_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'sku id',
  `price` decimal(10,4) NOT NULL DEFAULT '0.0000' COMMENT '单价',
  `quantity` int unsigned NOT NULL DEFAULT '0' COMMENT '数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_order_detail`
--

LOCK TABLES `td_order_detail` WRITE;
/*!40000 ALTER TABLE `td_order_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_order_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_team`
--

DROP TABLE IF EXISTS `td_team`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_team` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '名字',
  `number` int unsigned NOT NULL DEFAULT '0' COMMENT '人数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='团队';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_team`
--

LOCK TABLES `td_team` WRITE;
/*!40000 ALTER TABLE `td_team` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_team` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_team_user`
--

DROP TABLE IF EXISTS `td_team_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_team_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `team_id` bigint unsigned NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='团队成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_team_user`
--

LOCK TABLES `td_team_user` WRITE;
/*!40000 ALTER TABLE `td_team_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `td_team_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `td_user`
--

DROP TABLE IF EXISTS `td_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `td_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(50) NOT NULL COMMENT '''昵称''',
  `account` varchar(100) NOT NULL COMMENT '''账号''',
  `email` varchar(100) DEFAULT NULL COMMENT '''邮箱''',
  `password` varchar(100) NOT NULL COMMENT '''密码''',
  `user_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '''1-普通用户,2管理员,100超级管理员''',
  `user_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '''1-可用,2-不可用,3-未激活''',
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `td_user`
--

LOCK TABLES `td_user` WRITE;
/*!40000 ALTER TABLE `td_user` DISABLE KEYS */;
INSERT INTO `td_user` VALUES (2,'tangzhiqiang','tangzhiqiang','tangzhiqiang@test.com','$2a$10$Qer3Rxy8fJd7GGGLjB76Q.uAq6m4o/R9vME86yK.tjwUqLRUXuc9W',1,1,'','','','3b5b3d702a9637860ac351550859cd19','2023-06-16 05:55:45','::1','2023-06-14 17:55:45','2023-05-15 17:10:08','2023-06-14 17:55:45',NULL);
/*!40000 ALTER TABLE `td_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'treasure_doc'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-06-15 10:27:29
