/*
SQLyog Ultimate v13.1.1 (64 bit)
MySQL - 8.0.18 : Database - test
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`test` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `test`;

/*Table structure for table `classes` */

DROP TABLE IF EXISTS `classes`;

CREATE TABLE `classes` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'incr id',
  `class_no` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'class no',
  `user_id` int(11) NOT NULL COMMENT 'student id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `classes` */

insert  into `classes`(`id`,`class_no`,`user_id`,`created_at`,`updated_at`) values 
(1,'S-01',3,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(2,'S-01',4,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(3,'S-01',5,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(4,'S-01',6,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(5,'S-02',7,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(6,'S-02',8,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(7,'S-02',9,'2020-04-10 10:08:08','2020-04-10 10:08:08'),
(8,'S-02',10,'2020-04-10 10:08:08','2020-04-10 10:08:08');

/*Table structure for table `users` */

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'auto inc id',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'user name',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'phone number',
  `sex` tinyint(1) NOT NULL DEFAULT '1' COMMENT 'user sex',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'email',
  `disable` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'disabled(0=false 1=true)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

/*Data for the table `users` */

insert  into `users`(`id`,`name`,`phone`,`sex`,`email`,`disable`,`created_at`,`updated_at`) values 
(1,'lory','8618688888888',2,'2822922@qq.com',0,'2020-03-11 08:46:57','2020-04-13 12:23:34'),
(2,'lucas','8618699999999',1,'john@gmail.com',0,'2020-03-11 08:46:57','2020-03-11 16:02:51'),
(3,'std00','8618600000000',1,'admin@golang.org',1,'2020-03-11 14:42:53','2020-04-13 12:09:50'),
(4,'std01','8618600000001',1,'user1@hotmail.com',0,'2020-03-11 16:58:45','2020-04-10 10:07:15'),
(5,'std02','8618600000002',1,'user2@hotmail.com',0,'2020-03-11 16:58:45','2020-04-10 10:07:20'),
(6,'std03','8618600000003',1,'user1@hotmail.com',0,'2020-03-11 16:59:58','2020-04-10 10:07:22'),
(7,'std04','8618600000004',1,'user2@hotmail.com',0,'2020-03-11 16:59:58','2020-04-10 10:07:26'),
(9,'std05','8618600000005',1,'user1@hotmail.com',0,'2020-03-11 17:03:51','2020-04-10 10:07:28'),
(10,'std06','8618600000006',1,'user2@hotmail.com',0,'2020-03-11 17:03:51','2020-04-10 10:07:29'),
(11,'std07','8618600000007',1,'user1@hotmail.com',0,'2020-03-11 17:04:17','2020-04-10 10:07:31'),
(12,'std08','8618600000008',2,'user2@hotmail.com',0,'2020-03-11 17:04:17','2020-04-10 10:07:33'),
(13,'std09','8618600000009',1,'user3@hotmail.com',0,'2020-03-11 17:04:49','2020-04-10 10:07:35'),
(14,'std10','8618600000010',2,'user4@hotmail.com',0,'2020-03-11 17:04:49','2020-04-10 10:07:38');

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
