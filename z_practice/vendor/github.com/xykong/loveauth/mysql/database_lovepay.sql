/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50639
 Source Host           : localhost:3306
 Source Schema         : lovepay

 Target Server Type    : MySQL
 Target Server Version : 50639
 File Encoding         : 65001

 Date: 04/03/2019 16:23:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for activities
-- ----------------------------
DROP TABLE IF EXISTS `activities`;
CREATE TABLE `activities` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `global_id` bigint(20) DEFAULT NULL,
  `open_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` bigint(20) DEFAULT NULL,
  `activity_id` bigint(20) DEFAULT NULL,
  `activity_group_id` bigint(20) DEFAULT NULL,
  `activity_count` bigint(20) DEFAULT NULL,
  `activity_occur` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_activities_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for bilibili_pay_callbacks
-- ----------------------------
DROP TABLE IF EXISTS `bilibili_pay_callbacks`;
CREATE TABLE `bilibili_pay_callbacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `bili_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_no` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `money` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_money` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_money` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `merchant_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `product_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `product_desc` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `client_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `extension_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_status` int(11) DEFAULT NULL,
  `sign` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_bilibili_pay_callbacks_deleted_at` (`deleted_at`),
  KEY `idx_bilibili_pay_callbacks_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for do_alipay_callback_reqs
-- ----------------------------
DROP TABLE IF EXISTS `do_alipay_callback_reqs`;
CREATE TABLE `do_alipay_callback_reqs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `auth_app_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `app_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `charset` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_biz_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `buyer_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `buyer_logon_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `seller_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `seller_email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `receipt_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `invoice_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `buyer_pay_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `point_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `refund_fee` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `subject` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `body` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gmt_create` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gmt_payment` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gmt_refund` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gmt_close` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fund_bill_list` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `passback_params` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `voucher_detail_list` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_do_alipay_callback_reqs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for do_midas_callback_reqs
-- ----------------------------
DROP TABLE IF EXISTS `do_midas_callback_reqs`;
CREATE TABLE `do_midas_callback_reqs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `open_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `app_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time_stamp` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_item` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bill_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `provide_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `appmeta` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cft_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `clientver` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `payamt_coins` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pubacct_payamt_coins` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bazinga` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sig` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_do_midas_callback_reqs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for douyin_call_back_requests
-- ----------------------------
DROP TABLE IF EXISTS `douyin_call_back_requests`;
CREATE TABLE `douyin_call_back_requests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `notify_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `way` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `client_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_no` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total_fee` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `buyer_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tt_sign` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tt_sign_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_douyin_call_back_requests_deleted_at` (`deleted_at`),
  KEY `idx_douyin_call_back_requests_trade_no` (`trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for huawei_pay_callbacks
-- ----------------------------
DROP TABLE IF EXISTS `huawei_pay_callbacks`;
CREATE TABLE `huawei_pay_callbacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `result` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `product_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_type` int(11) DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `currency` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notify_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `request_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bank_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `access_mode` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `spending` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ext_reserved` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sys_reserved` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_huawei_pay_callbacks_deleted_at` (`deleted_at`),
  KEY `idx_huawei_pay_callbacks_order_id` (`order_id`),
  KEY `idx_huawei_pay_callbacks_ext_reserved` (`ext_reserved`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for inventories
-- ----------------------------
DROP TABLE IF EXISTS `inventories`;
CREATE TABLE `inventories` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `global_id` bigint(20) DEFAULT NULL,
  `open_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` bigint(20) DEFAULT NULL,
  `item_id` bigint(20) DEFAULT NULL,
  `item_count` bigint(20) DEFAULT NULL,
  `item_history` bigint(20) DEFAULT NULL,
  `item_purchased` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for mgtv_call_back_requests
-- ----------------------------
DROP TABLE IF EXISTS `mgtv_call_back_requests`;
CREATE TABLE `mgtv_call_back_requests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `sign` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `uuid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `business_order_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_create` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total_fee` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ext_data` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_mgtv_call_back_requests_deleted_at` (`deleted_at`),
  KEY `idx_mgtv_call_back_requests_business_order_id` (`business_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `global_id` bigint(20) DEFAULT NULL,
  `shop_id` int(11) DEFAULT NULL,
  `sequence` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `timestamp` timestamp NULL DEFAULT NULL,
  `state` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `num` int(11) DEFAULT '1',
  `amount` int(11) DEFAULT NULL,
  `sns_order_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_method` int(11) DEFAULT NULL,
  `vendor` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_sns_order_id` (`sns_order_id`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_global_id` (`global_id`),
  KEY `idx_orders_sequence` (`sequence`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for q_pay_callbacks
-- ----------------------------
DROP TABLE IF EXISTS `q_pay_callbacks`;
CREATE TABLE `q_pay_callbacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `appid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mch_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nonce_str` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_state` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bank_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fee_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total_fee` int(11) DEFAULT NULL,
  `cash_fee` int(11) DEFAULT NULL,
  `coupon_fee` int(11) DEFAULT NULL,
  `transaction_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `attach` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time_end` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `openid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_q_pay_callbacks_deleted_at` (`deleted_at`),
  KEY `idx_q_pay_callbacks_transaction_id` (`transaction_id`),
  KEY `idx_q_pay_callbacks_attach` (`attach`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for quick_pay_callbacks
-- ----------------------------
DROP TABLE IF EXISTS `quick_pay_callbacks`;
CREATE TABLE `quick_pay_callbacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `is_test` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel_uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `game_order` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `extras_params` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_quick_pay_callbacks_deleted_at` (`deleted_at`),
  KEY `idx_quick_pay_callbacks_game_order` (`game_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for repire_order_requests
-- ----------------------------
DROP TABLE IF EXISTS `repire_order_requests`;
CREATE TABLE `repire_order_requests` (
  `global_id` bigint(20) DEFAULT NULL,
  `sequence` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mail_title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mail_content` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `award_list` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `operator` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `repire_time` timestamp NULL DEFAULT NULL,
  KEY `idx_repire_order_requests_global_id` (`global_id`),
  KEY `idx_repire_order_requests_operator` (`operator`),
  KEY `idx_repire_order_requests_repire_time` (`repire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for vivo_query_order_responses
-- ----------------------------
DROP TABLE IF EXISTS `vivo_query_order_responses`;
CREATE TABLE `vivo_query_order_responses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `ret` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `message` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign_method` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `signature` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cp_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `app_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cp_order_number` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_number` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `order_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ext_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_vivo_query_order_responses_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for wechat_pay_callbacks
-- ----------------------------
DROP TABLE IF EXISTS `wechat_pay_callbacks`;
CREATE TABLE `wechat_pay_callbacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `return_code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `return_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `appid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mch_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_info` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nonce_str` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sign` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `result_code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `err_code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `err_code_des` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `openid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_subscribe` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `trade_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bank_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total_fee` int(11) DEFAULT NULL,
  `fee_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cash_fee` int(11) DEFAULT NULL,
  `cash_fee_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `transaction_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `out_trade_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `attach` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time_end` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_wechat_pay_callbacks_deleted_at` (`deleted_at`),
  KEY `idx_wechat_pay_callbacks_transaction_id` (`transaction_id`),
  KEY `idx_wechat_pay_callbacks_attach` (`attach`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
