/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50639
 Source Host           : localhost:3306
 Source Schema         : db_loveworld_online

 Target Server Type    : MySQL
 Target Server Version : 50639
 File Encoding         : 65001

 Date: 30/01/2019 10:21:24
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for online_cnt
-- ----------------------------
DROP TABLE IF EXISTS `online_cnt`;
CREATE TABLE `online_cnt` (
  `channel` varchar(100) NOT NULL DEFAULT '',
  `timekey` int(11) NOT NULL DEFAULT '0',
  `gsid` varchar(32) NOT NULL DEFAULT '',
  `onlinecntios` int(11) NOT NULL DEFAULT '0',
  `onlinecntandroid` int(11) NOT NULL DEFAULT '0',
  KEY `timekey` (`timekey`,`gsid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for player_payment
-- ----------------------------
DROP TABLE IF EXISTS `player_payment`;
CREATE TABLE `player_payment` (
  `timekey` int(11) DEFAULT NULL,
  `global_id` bigint(20) DEFAULT NULL,
  `device_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount` double DEFAULT NULL,
  `channel` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `platform` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sequence` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_svr_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  KEY `idx_player_payment_platform` (`platform`),
  KEY `idx_player_payment_game_svr_id` (`game_svr_id`),
  KEY `idx_player_payment_timekey` (`timekey`),
  KEY `idx_player_payment_global_id` (`global_id`),
  KEY `idx_player_payment_device_id` (`device_id`),
  KEY `idx_player_payment_channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for player_register
-- ----------------------------
DROP TABLE IF EXISTS `player_register`;
CREATE TABLE `player_register` (
  `timekey` int(11) DEFAULT NULL,
  `global_id` bigint(20) DEFAULT NULL,
  `device_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_svr_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  KEY `idx_player_register_timekey` (`timekey`),
  KEY `idx_player_register_global_id` (`global_id`),
  KEY `idx_player_register_device_id` (`device_id`),
  KEY `idx_player_register_channel` (`channel`),
  KEY `idx_player_register_game_svr_id` (`game_svr_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
