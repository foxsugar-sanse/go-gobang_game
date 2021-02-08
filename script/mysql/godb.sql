/*
 Navicat Premium Data Transfer

 Source Server         : GO_WEB_MYSQL
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : 192.168.31.161:3306
 Source Schema         : godb

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 08/02/2021 17:05:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_friend
-- ----------------------------
DROP TABLE IF EXISTS `user_friend`;
CREATE TABLE `user_friend` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `main_uid` int(11) DEFAULT NULL,
  `friend_uid` int(11) DEFAULT NULL,
  `friend_note` varchar(255) DEFAULT NULL,
  `user_group` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `user_nick_name` varchar(255) DEFAULT NULL,
  `user_age` int(3) DEFAULT NULL,
  `user_sex` varchar(1) DEFAULT NULL,
  `user_brief` varchar(255) DEFAULT NULL,
  `user_contact` varchar(255) DEFAULT NULL,
  `pass_word` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
