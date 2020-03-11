/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1-3306
 Source Server Type    : MySQL
 Source Server Version : 80018
 Source Host           : localhost:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 80018
 File Encoding         : 65001

 Date: 11/03/2020 08:49:40
*/

CREATE DATABASE IF NOT EXISTS test DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

USE test;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'auto inc id',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `sex` tinyint(1) NOT NULL DEFAULT 1,
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'rose', '8613011112222', 2, 'rose123@hotmail.com', '2020-03-11 08:46:57', '2020-03-11 08:48:14');
INSERT INTO `users` VALUES (2, 'john', '8613100003333', 1, 'john333@hotmail.com', '2020-03-11 08:46:57', '2020-03-11 08:47:59');

SET FOREIGN_KEY_CHECKS = 1;
