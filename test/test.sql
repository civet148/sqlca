/*
SQLyog Ultimate
MySQL - 8.0.18 : Database - test
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `test`;

/*Table structure for table `classes` */

CREATE TABLE `classes` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'incr id',
  `class_no` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'class no',
  `user_id` int(11) NOT NULL COMMENT 'student id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `classes` */

insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (1,'S-01',1,'2020-04-10 10:08:08','2020-05-12 19:39:43');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (2,'S-01',2,'2020-04-10 10:08:08','2020-05-12 19:39:44');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (3,'S-01',3,'2020-04-10 10:08:08','2020-05-12 19:39:45');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (4,'S-01',4,'2020-04-10 10:08:08','2020-05-12 19:39:45');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (5,'S-02',5,'2020-04-10 10:08:08','2020-05-12 19:39:46');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (6,'S-02',8,'2020-04-10 10:08:08','2020-05-12 19:40:00');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (7,'S-02',9,'2020-04-10 10:08:08','2020-04-10 10:08:08');
insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values (8,'S-02',10,'2020-04-10 10:08:08','2020-04-10 10:08:08');

/*Table structure for table `users` */

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'auto inc id',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'user name',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'phone number',
  `sex` tinyint(1) NOT NULL COMMENT 'user sex',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'email',
  `disable` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'disabled(0=false 1=true)',
  `balance` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT 'balance of decimal',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

/*Data for the table `users` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
