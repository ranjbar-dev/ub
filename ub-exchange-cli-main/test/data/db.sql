-- MariaDB dump 10.19  Distrib 10.7.3-MariaDB, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: prunibd
-- ------------------------------------------------------
-- Server version	10.2.29-MariaDB-1:10.2.29+maria~bionic

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `action_roles`
--

DROP TABLE IF EXISTS `action_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action_roles` (
  `role_id` int(11) NOT NULL,
  `action_id` int(11) NOT NULL,
  PRIMARY KEY (`role_id`,`action_id`),
  KEY `IDX_51F0E27BD60322AC` (`role_id`),
  KEY `IDX_51F0E27B9D32F035` (`action_id`),
  CONSTRAINT `FK_51F0E27B9D32F035` FOREIGN KEY (`action_id`) REFERENCES `actions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_51F0E27BD60322AC` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `actions`
--

DROP TABLE IF EXISTS `actions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `actions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_548F1EF77153098` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=151 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_menu_roles`
--

DROP TABLE IF EXISTS `admin_menu_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin_menu_roles` (
  `role_id` int(11) NOT NULL,
  `admin_menu_id` int(11) NOT NULL,
  PRIMARY KEY (`role_id`,`admin_menu_id`),
  KEY `IDX_EFACAB25D60322AC` (`role_id`),
  KEY `IDX_EFACAB258EB6B9B` (`admin_menu_id`),
  CONSTRAINT `FK_EFACAB258EB6B9B` FOREIGN KEY (`admin_menu_id`) REFERENCES `admin_menus` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_EFACAB25D60322AC` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admin_menus`
--

DROP TABLE IF EXISTS `admin_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin_menus` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `name` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_D25FF335727ACA70` (`parent_id`),
  CONSTRAINT `FK_D25FF335727ACA70` FOREIGN KEY (`parent_id`) REFERENCES `admin_menus` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `app_version`
--

DROP TABLE IF EXISTS `app_version`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_version` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `version` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `platform` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `force_update` tinyint(1) NOT NULL,
  `key_features` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `release_date` datetime NOT NULL,
  `bug_fixes` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `version_code` double DEFAULT NULL,
  `store_url` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `audit_log`
--

DROP TABLE IF EXISTS `audit_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `audit_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `changed_by_id` int(11) DEFAULT NULL,
  `entity_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `entity_id` int(11) NOT NULL,
  `change_set` longtext COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:json)',
  PRIMARY KEY (`id`),
  KEY `IDX_F6E1C0F5828AD0A0` (`changed_by_id`),
  CONSTRAINT `FK_F6E1C0F5828AD0A0` FOREIGN KEY (`changed_by_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1331 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `communication_templates`
--

DROP TABLE IF EXISTS `communication_templates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `communication_templates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `title` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `template_text` longtext COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `configurations`
--

DROP TABLE IF EXISTS `configurations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `configurations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_code` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `k` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `v` varchar(1023) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `countries`
--

DROP TABLE IF EXISTS `countries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `countries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `full_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `image_path` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `iso31661alpha3` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `iso31661alpha2` varchar(4) COLLATE utf8_unicode_ci DEFAULT NULL,
  `iso4217currency_alphabetic_code` varchar(8) COLLATE utf8_unicode_ci DEFAULT NULL,
  `uniterm_english_formal` varchar(60) COLLATE utf8_unicode_ci DEFAULT NULL,
  `region_name` varchar(8) COLLATE utf8_unicode_ci DEFAULT NULL,
  `languages` varchar(97) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=813 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_balance`
--

DROP TABLE IF EXISTS `crypto_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_balance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `currency_id` int(11) DEFAULT NULL,
  `external_exchange_id` int(11) DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `address` varchar(1024) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tag` varchar(1024) COLLATE utf8_unicode_ci DEFAULT NULL,
  `metaData` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `free_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `locked_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `blockchain_network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_crypto_balance` (`currency_id`,`type`,`external_exchange_id`),
  KEY `IDX_7D88B91138248176` (`currency_id`),
  KEY `IDX_7D88B911689CD7` (`external_exchange_id`),
  CONSTRAINT `FK_7D88B91138248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_7D88B911689CD7` FOREIGN KEY (`external_exchange_id`) REFERENCES `external_exchanges` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=96 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_internal_transfer`
--

DROP TABLE IF EXISTS `crypto_internal_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_internal_transfer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `crypto_balance_from_id` int(11) NOT NULL,
  `crypto_balance_to_id` int(11) DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `tx_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `metadata` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `blockchain_network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `to_custom_address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_288A0FACFADBB17B` (`crypto_balance_from_id`),
  KEY `IDX_288A0FAC1B947875` (`crypto_balance_to_id`),
  CONSTRAINT `FK_288A0FAC1B947875` FOREIGN KEY (`crypto_balance_to_id`) REFERENCES `crypto_balance` (`id`),
  CONSTRAINT `FK_288A0FACFADBB17B` FOREIGN KEY (`crypto_balance_from_id`) REFERENCES `crypto_balance` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_news`
--

DROP TABLE IF EXISTS `crypto_news`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_news` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `link` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `currency_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_C70CA78F38248176` (`currency_id`),
  CONSTRAINT `FK_C70CA78F38248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_payment_admin_comment`
--

DROP TABLE IF EXISTS `crypto_payment_admin_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_payment_admin_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_id` int(11) DEFAULT NULL,
  `crypto_payment_id` int(11) DEFAULT NULL,
  `adminComment` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_A8164724642B8210` (`admin_id`),
  KEY `IDX_A8164724793128F` (`crypto_payment_id`),
  CONSTRAINT `FK_A8164724642B8210` FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_A8164724793128F` FOREIGN KEY (`crypto_payment_id`) REFERENCES `crypto_payments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_payment_extra_info`
--

DROP TABLE IF EXISTS `crypto_payment_extra_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_payment_extra_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `crypto_payment_id` int(11) DEFAULT NULL,
  `last_handled_id` int(11) DEFAULT NULL,
  `tag` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `network_fee` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_message` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `rejection_reason` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `ip` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `auto_transfer` tinyint(1) DEFAULT NULL,
  `auto_exchange_order_id` int(11) DEFAULT NULL,
  `auto_exchange_failure_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `auto_exchange_failure_reason` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `btc_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `external_exchange_id` int(11) DEFAULT NULL,
  `external_exchange_withdraw_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `external_exchange_withdraw_info` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_A6A91828793128F` (`crypto_payment_id`),
  UNIQUE KEY `UNIQ_A6A918284DDB48A2` (`auto_exchange_order_id`),
  KEY `IDX_A6A918283AEE42EA` (`last_handled_id`),
  KEY `IDX_A6A91828689CD7` (`external_exchange_id`),
  CONSTRAINT `FK_A6A918283AEE42EA` FOREIGN KEY (`last_handled_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_A6A918284DDB48A2` FOREIGN KEY (`auto_exchange_order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `FK_A6A91828689CD7` FOREIGN KEY (`external_exchange_id`) REFERENCES `external_exchanges` (`id`),
  CONSTRAINT `FK_A6A91828793128F` FOREIGN KEY (`crypto_payment_id`) REFERENCES `crypto_payments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=491 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `crypto_payments`
--

DROP TABLE IF EXISTS `crypto_payments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crypto_payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `type` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `money_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `money_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `from_address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `to_address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `fee` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tx_id` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  `admin_status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `withdraw_type` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `blockchain_network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `fee_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tx_id_currency_type_idx` (`tx_id`,`currency_id`,`type`),
  KEY `IDX_55B6187DA76ED395` (`user_id`),
  KEY `IDX_55B6187D38248176` (`currency_id`),
  CONSTRAINT `FK_55B6187D38248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_55B6187DA76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=614 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `currencies`
--

DROP TABLE IF EXISTS `currencies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `currencies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `code` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `sub_unit` int(11) NOT NULL,
  `priority` int(11) NOT NULL,
  `conversion_ratio` double NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `image_path` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `is_main` tinyint(1) NOT NULL,
  `deposit_fee` double DEFAULT NULL,
  `withdrawal_fee` double DEFAULT NULL,
  `extra_info_id` int(11) DEFAULT NULL,
  `minimum_withdraw` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `supports_withdraw` tinyint(1) DEFAULT NULL,
  `supports_deposit` tinyint(1) DEFAULT NULL,
  `maximum_withdraw` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `blockchain_network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `show_sub_unit` int(11) DEFAULT NULL,
  `other_blockchain_networks_configs` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `completed_network_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deposit_comments` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `withdraw_comments` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `second_image_path` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_37C446938550C61A` (`extra_info_id`),
  CONSTRAINT `FK_37C446938550C61A` FOREIGN KEY (`extra_info_id`) REFERENCES `currencies_extra_info` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `currencies_extra_info`
--

DROP TABLE IF EXISTS `currencies_extra_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `currencies_extra_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `issue_date` datetime NOT NULL,
  `total_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `circulation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `links` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `description` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `background_image_path` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dates`
--

DROP TABLE IF EXISTS `dates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_AB74C91EAA9E377A` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `exchange_asset`
--

DROP TABLE IF EXISTS `exchange_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `exchange_asset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `type` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `blockchain_network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `tx_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `source` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `date` date NOT NULL,
  `description` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_1C981614FDFB2543` (`tx_id`),
  KEY `IDX_1C981614A76ED395` (`user_id`),
  KEY `IDX_1C98161438248176` (`currency_id`),
  CONSTRAINT `FK_1C98161438248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_1C981614A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `external_exchange_orders`
--

DROP TABLE IF EXISTS `external_exchange_orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `external_exchange_orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pair_currency_id` int(11) DEFAULT NULL,
  `external_exchange_id` int(11) DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `exchange_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `external_exchange_other_info` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_get_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_pay_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_trade_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_fee_percentage` double DEFAULT NULL,
  `last_trade_id` int(11) DEFAULT NULL,
  `fail_reason` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `buy_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `buy_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `sell_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `sell_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `order_ids` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `exception_message` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `meta_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_order_id` int(11) DEFAULT NULL,
  `source` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_D44E2D556D128938` (`user_order_id`),
  KEY `IDX_D44E2D552D772A06` (`pair_currency_id`),
  KEY `IDX_D44E2D55689CD7` (`external_exchange_id`),
  CONSTRAINT `FK_D44E2D552D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`),
  CONSTRAINT `FK_D44E2D55689CD7` FOREIGN KEY (`external_exchange_id`) REFERENCES `external_exchanges` (`id`),
  CONSTRAINT `FK_D44E2D556D128938` FOREIGN KEY (`user_order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=97104 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `external_exchanges`
--

DROP TABLE IF EXISTS `external_exchanges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `external_exchanges` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `meta_data` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `status` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `migration_versions`
--

DROP TABLE IF EXISTS `migration_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migration_versions` (
  `version` varchar(14) COLLATE utf8_unicode_ci NOT NULL,
  `executed_at` datetime NOT NULL COMMENT '(DC2Type:datetime_immutable)',
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ohlc`
--

DROP TABLE IF EXISTS `ohlc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ohlc` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pair_currency_id` int(11) DEFAULT NULL,
  `time_frame` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `open_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `close_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `low_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `high_price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `base_volume` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `quote_volume` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `taker_buy_base_volume` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `taker_buy_quote_volume` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `pair_currency_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `spread` double NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_ohlc` (`pair_currency_id`,`time_frame`,`start_time`),
  KEY `IDX_1936B6982D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_1936B6982D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ohlc_sync`
--

DROP TABLE IF EXISTS `ohlc_sync`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ohlc_sync` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pair_currency_id` int(11) DEFAULT NULL,
  `time_frame` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `with_update` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_1387C8CC2D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_1387C8CC2D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21841 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_from_external`
--

DROP TABLE IF EXISTS `order_from_external`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_from_external` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pair_currency_id` int(11) DEFAULT NULL,
  `external_order_id` bigint(20) NOT NULL,
  `client_order_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `side` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `meta_data` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `timestamp` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_B3DE09F06293DFCB` (`external_order_id`),
  KEY `IDX_B3DE09F02D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_B3DE09F02D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=678 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `creator_user_id` int(11) NOT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `type` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `exchange_type` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `demanded_money_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `demanded_money_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `payed_by_money_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `payed_by_money_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `pair_currency_id` int(11) NOT NULL,
  `extra_info_id` int(11) DEFAULT NULL,
  `final_demanded_money` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `trade_price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_maker` tinyint(1) DEFAULT NULL,
  `fee_percentage` double DEFAULT NULL,
  `external_exchange_fee_percentage` double DEFAULT NULL,
  `final_payed_by_money` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `level` int(11) DEFAULT NULL,
  `path` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_demanded_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `final_payed_by_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `stop_point_price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_submitted` tinyint(1) DEFAULT NULL,
  `is_traded_with_bot` tinyint(1) DEFAULT NULL,
  `current_market_price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_fast_exchange` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_E52FFDEE727ACA70` (`parent_id`),
  UNIQUE KEY `UNIQ_E52FFDEE8550C61A` (`extra_info_id`),
  KEY `IDX_E52FFDEE29FC6AE1` (`creator_user_id`),
  KEY `IDX_E52FFDEE2D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_E52FFDEE29FC6AE1` FOREIGN KEY (`creator_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_E52FFDEE2D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`),
  CONSTRAINT `FK_E52FFDEE727ACA70` FOREIGN KEY (`parent_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `FK_E52FFDEE8550C61A` FOREIGN KEY (`extra_info_id`) REFERENCES `orders_extra_info` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6549 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders_extra_info`
--

DROP TABLE IF EXISTS `orders_extra_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `orders_extra_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_agent_info` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `external_exchange_other_info` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `external_exchange_id` int(11) DEFAULT NULL,
  `is_market_order_in_external_exchange` tinyint(1) DEFAULT NULL,
  `payed_by_diff` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `auto_exchange` tinyint(1) DEFAULT NULL,
  `external_exchange_order_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_3A36EBCE689CD7` (`external_exchange_id`),
  CONSTRAINT `FK_3A36EBCE689CD7` FOREIGN KEY (`external_exchange_id`) REFERENCES `external_exchanges` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6352 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pair_currencies`
--

DROP TABLE IF EXISTS `pair_currencies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pair_currencies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `basis_currency_id` int(11) NOT NULL,
  `dependent_currency_id` int(11) NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `maker_fee` double DEFAULT NULL,
  `taker_fee` double DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT NULL,
  `max_our_exchange_limit` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `minimum_order_amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_main` tinyint(1) DEFAULT NULL,
  `bot_rules` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `bot_orders_aggregation_time` int(11) DEFAULT NULL,
  `ohlc_spread` double NOT NULL,
  `trade_status` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `show_digits` int(11) DEFAULT NULL,
  `aggregation_status` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `basis_dependent_unique` (`basis_currency_id`,`dependent_currency_id`),
  KEY `IDX_BE427FCED6C70F73` (`basis_currency_id`),
  KEY `IDX_BE427FCEC1EA2D36` (`dependent_currency_id`),
  CONSTRAINT `FK_BE427FCEC1EA2D36` FOREIGN KEY (`dependent_currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_BE427FCED6C70F73` FOREIGN KEY (`basis_currency_id`) REFERENCES `currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `role` varchar(64) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_B63E2EC757698A6A` (`role`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trade_from_external`
--

DROP TABLE IF EXISTS `trade_from_external`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trade_from_external` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) DEFAULT NULL,
  `external_trade_id` bigint(20) NOT NULL,
  `price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `commission` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `commission_coin` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `meta_data` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `timestamp` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_8FCE449EE3212593` (`external_trade_id`),
  KEY `IDX_8FCE449E8D9F6D38` (`order_id`),
  CONSTRAINT `FK_8FCE449E8D9F6D38` FOREIGN KEY (`order_id`) REFERENCES `order_from_external` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trades`
--

DROP TABLE IF EXISTS `trades`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trades` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `price` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `pair_currency_id` int(11) NOT NULL,
  `buy_order_id` int(11) DEFAULT NULL,
  `sell_order_id` int(11) DEFAULT NULL,
  `bot_order_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_BFA111257FC358ED` (`buy_order_id`),
  UNIQUE KEY `UNIQ_BFA111256CF89127` (`sell_order_id`),
  KEY `IDX_BFA111252D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_BFA111252D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`),
  CONSTRAINT `FK_BFA111256CF89127` FOREIGN KEY (`sell_order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `FK_BFA111257FC358ED` FOREIGN KEY (`buy_order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=237956 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `order_id` int(11) DEFAULT NULL,
  `type` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `money_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `money_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `crypto_payment_id` int(11) DEFAULT NULL,
  `amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_id_type_idx` (`order_id`,`type`),
  KEY `IDX_EAA81A4CA76ED395` (`user_id`),
  KEY `IDX_EAA81A4C38248176` (`currency_id`),
  KEY `IDX_EAA81A4C8D9F6D38` (`order_id`),
  KEY `IDX_EAA81A4C793128F` (`crypto_payment_id`),
  CONSTRAINT `FK_EAA81A4C38248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_EAA81A4C793128F` FOREIGN KEY (`crypto_payment_id`) REFERENCES `crypto_payments` (`id`),
  CONSTRAINT `FK_EAA81A4C8D9F6D38` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
  CONSTRAINT `FK_EAA81A4CA76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=220737 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_admin_comment`
--

DROP TABLE IF EXISTS `user_admin_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_admin_comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `admin_id` int(11) DEFAULT NULL,
  `comment` longtext COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_BEC67DB5A76ED395` (`user_id`),
  KEY `IDX_BEC67DB5642B8210` (`admin_id`),
  CONSTRAINT `FK_BEC67DB5642B8210` FOREIGN KEY (`admin_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_BEC67DB5A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_balances`
--

DROP TABLE IF EXISTS `user_balances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_balances` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `frozen_balance` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `balance_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `balance_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `frozen_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `other_addresses` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `auto_exchange_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_currency_idx` (`user_id`,`currency_id`),
  KEY `IDX_A12A000FA76ED395` (`user_id`),
  KEY `IDX_A12A000F38248176` (`currency_id`),
  CONSTRAINT `FK_A12A000F38248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_A12A000FA76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3507 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_balances_history`
--

DROP TABLE IF EXISTS `user_balances_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_balances_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_balance_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `btc_amount` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `usdt_amount` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `balance_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `balance_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `user_id` int(11) NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_currency_date_idx` (`user_id`,`currency_id`,`date`),
  KEY `IDX_61A213EC9F66531` (`user_balance_id`),
  KEY `IDX_61A213EC38248176` (`currency_id`),
  KEY `IDX_61A213ECA76ED395` (`user_id`),
  CONSTRAINT `FK_61A213EC38248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_61A213EC9F66531` FOREIGN KEY (`user_balance_id`) REFERENCES `user_balances` (`id`),
  CONSTRAINT `FK_61A213ECA76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=148903 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_configs`
--

DROP TABLE IF EXISTS `user_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_configs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `theme` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `mode` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `custom_theme_data` longtext COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '(DC2Type:json)',
  `is_trade_notification_enabled` tinyint(1) NOT NULL,
  `is_email_verification_for_withdraw_enabled` tinyint(1) NOT NULL,
  `is_two_fa_verification_for_withdraw_enabled` tinyint(1) NOT NULL,
  `is_email_verification_for_login_enabled` tinyint(1) NOT NULL,
  `is_two_fa_verification_for_login_enabled` tinyint(1) NOT NULL,
  `is_white_list_enabled` tinyint(1) NOT NULL,
  `is_read_only` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_1A646639A76ED395` (`user_id`),
  CONSTRAINT `FK_1A646639A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_devices`
--

DROP TABLE IF EXISTS `user_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_devices` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `user_agent` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `is_logged_out` tinyint(1) NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `ip` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `platform` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `browser` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_490A5090A76ED395` (`user_id`),
  CONSTRAINT `FK_490A5090A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_exchange_statistics`
--

DROP TABLE IF EXISTS `user_exchange_statistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_exchange_statistics` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `creator_user_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `trade_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `trade_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `base_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `base_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `currency_id` int(11) NOT NULL,
  `abs_money` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `base_abs_money` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `exchange_number` int(11) NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_currency_date_idx` (`creator_user_id`,`currency_id`,`date`),
  KEY `IDX_8434DC9729FC6AE1` (`creator_user_id`),
  KEY `IDX_8434DC9738248176` (`currency_id`),
  CONSTRAINT `FK_8434DC9729FC6AE1` FOREIGN KEY (`creator_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_8434DC9738248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=645 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_favorite_pair_currency`
--

DROP TABLE IF EXISTS `user_favorite_pair_currency`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_favorite_pair_currency` (
  `user_id` int(11) NOT NULL,
  `pair_currency_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`pair_currency_id`),
  KEY `IDX_14520C35A76ED395` (`user_id`),
  KEY `IDX_14520C352D772A06` (`pair_currency_id`),
  CONSTRAINT `FK_14520C352D772A06` FOREIGN KEY (`pair_currency_id`) REFERENCES `pair_currencies` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_14520C35A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_groups`
--

DROP TABLE IF EXISTS `user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_levels`
--

DROP TABLE IF EXISTS `user_levels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_levels` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `maker_fee_percentage` double NOT NULL,
  `taker_fee_percentage` double NOT NULL,
  `withdraw_fee_percentage` double NOT NULL,
  `deposit_fee_percentage` double NOT NULL,
  `code` int(11) DEFAULT NULL,
  `exchange_number_limit` int(11) NOT NULL,
  `exchange_volume_limit_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `exchange_volume_limit_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `min_exchange_volume` double DEFAULT NULL,
  `max_exchange_volume` double DEFAULT NULL,
  `min_kyc_level` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_login_history`
--

DROP TABLE IF EXISTS `user_login_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_login_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `device` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `ip` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `email` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  `password` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_D919C366A76ED395` (`user_id`),
  CONSTRAINT `FK_D919C366A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=509978 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_permissions`
--

DROP TABLE IF EXISTS `user_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_84F605FA5E237E06` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_profile_image`
--

DROP TABLE IF EXISTS `user_profile_image`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_profile_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_profile_id` int(11) DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `image_path` varchar(512) COLLATE utf8_unicode_ci NOT NULL,
  `confirmation_status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `original_file_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `id_card_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `rejection_reason` longtext COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `sub_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT NULL,
  `main_image_id` int(11) DEFAULT NULL,
  `is_back` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_DC6736DDF96E6191` (`image_path`),
  KEY `IDX_DC6736DD6B9DD454` (`user_profile_id`),
  KEY `IDX_DC6736DDE4873418` (`main_image_id`),
  CONSTRAINT `FK_DC6736DD6B9DD454` FOREIGN KEY (`user_profile_id`) REFERENCES `user_profiles` (`id`),
  CONSTRAINT `FK_DC6736DDE4873418` FOREIGN KEY (`main_image_id`) REFERENCES `user_profile_image` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=487 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_profiles`
--

DROP TABLE IF EXISTS `user_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `first_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `country_id` int(11) DEFAULT NULL,
  `last_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `gender` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date_of_birth` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `region_and_city` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `postal_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `id_card_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `admin_comment` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `identity_confirmation_status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address_confirmation_status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone_confirmation_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `registration_ip` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  `refer_key` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `trust_level` int(11) NOT NULL,
  `avatar_image_path` varchar(512) COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_uploaded_image_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_6BBD6130A76ED395` (`user_id`),
  KEY `IDX_6BBD6130F92F3E70` (`country_id`),
  CONSTRAINT `FK_6BBD6130A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_6BBD6130F92F3E70` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=235 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_role`
--

DROP TABLE IF EXISTS `user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_role` (
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `IDX_2DE8C6A3A76ED395` (`user_id`),
  KEY `IDX_2DE8C6A3D60322AC` (`role_id`),
  CONSTRAINT `FK_2DE8C6A3A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_2DE8C6A3D60322AC` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_wallet_balance`
--

DROP TABLE IF EXISTS `user_wallet_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_wallet_balance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `network_currency_id` int(11) DEFAULT NULL,
  `balance` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_B670E5B1A76ED395` (`user_id`),
  KEY `IDX_B670E5B138248176` (`currency_id`),
  KEY `IDX_B670E5B1FA2A2F2E` (`network_currency_id`),
  CONSTRAINT `FK_B670E5B138248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_B670E5B1A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_B670E5B1FA2A2F2E` FOREIGN KEY (`network_currency_id`) REFERENCES `currencies` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1906 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_withdraw_address`
--

DROP TABLE IF EXISTS `user_withdraw_address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_withdraw_address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `address` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `label` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT NULL,
  `is_favorite` tinyint(1) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `network` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_1B970987A76ED395` (`user_id`),
  KEY `IDX_1B97098738248176` (`currency_id`),
  KEY `label_search_idx` (`label`),
  CONSTRAINT `FK_1B97098738248176` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`),
  CONSTRAINT `FK_1B970987A76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=469 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `kyc` int(11) NOT NULL,
  `status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `user_level_id` int(11) DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `google2fa_secret_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `exchange_volume_amount` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `exchange_volume_currency` char(6) COLLATE utf8_unicode_ci NOT NULL COMMENT '(DC2Type:exchange_currency)',
  `ub_id` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `verification_code` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL,
  `number_of_exchange` int(11) NOT NULL,
  `google2fa_disabled_at` datetime DEFAULT NULL,
  `private_channel_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `manager_id` int(11) DEFAULT NULL,
  `group_id` int(11) DEFAULT NULL,
  `is_level_manually_set` tinyint(1) DEFAULT NULL,
  `account_status` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `is_two_fa_enabled` tinyint(1) NOT NULL,
  `refresh_token` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `refresh_token_expiry` datetime NOT NULL DEFAULT current_timestamp(),
  `password_changed_at` datetime DEFAULT NULL,
  `two_fa_changed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_1483A5E9E7927C74` (`email`),
  UNIQUE KEY `UNIQ_1483A5E9FF3FBB08` (`ub_id`),
  UNIQUE KEY `UNIQ_1483A5E9E821C39F` (`verification_code`),
  UNIQUE KEY `UNIQ_1483A5E9940A757A` (`private_channel_name`),
  UNIQUE KEY `UNIQ_1483A5E9C74F2195` (`refresh_token`),
  KEY `IDX_1483A5E9BF3CAFA7` (`user_level_id`),
  KEY `IDX_1483A5E9783E3463` (`manager_id`),
  KEY `IDX_1483A5E9FE54D947` (`group_id`),
  CONSTRAINT `FK_1483A5E9783E3463` FOREIGN KEY (`manager_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_1483A5E9BF3CAFA7` FOREIGN KEY (`user_level_id`) REFERENCES `user_levels` (`id`),
  CONSTRAINT `FK_1483A5E9FE54D947` FOREIGN KEY (`group_id`) REFERENCES `user_groups` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=259 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users_permissions`
--

DROP TABLE IF EXISTS `users_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users_permissions` (
  `user_id` int(11) NOT NULL,
  `user_permission_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`user_permission_id`),
  KEY `IDX_DA58F09DA76ED395` (`user_id`),
  KEY `IDX_DA58F09D1057A19A` (`user_permission_id`),
  CONSTRAINT `FK_DA58F09D1057A19A` FOREIGN KEY (`user_permission_id`) REFERENCES `user_permissions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_DA58F09DA76ED395` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'prunibd'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-02-28 18:43:56

