# author 
lory.li
# email
civet148@126.com
# QQ 
93864947
# sqlca
a enhancement database and cache tool based on sqlx and redigogo which based on redigo and go-redis-cluster

# database schema
```mysql
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

```

# data model
```golang
type UserDO struct {
	Id    int32  `db:"id"`  
	Name  string `db:"name"`  
	Phone string `db:"phone"` 
	Sex   int8   `db:"sex"`   
	Email string `db:"email"` 
}
```

# open database/redis

```golang

e := sqlca.NewEngine()
e.Debug(true) //debug mode on

// open database driver (requred)
e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")

// open redis driver for cache (optional)
e.Open("redis://127.0.0.1:6379/cluster?db=0", 3600) //redis standalone mode
``` 

# global variants 
```golang
const (
	TABLE_NAME_USERS = "users"
)
```

## orm: insert/upsert from data model

```golang
user := UserDO{
        Name:  "admin",
        Phone: "8618600000000",
        Sex:   1,
        Email: "admin@golang.org",
}

e.Model(&user).Table(TABLE_NAME_USERS).Insert()
```
```golang
user := UserDO{
    Id:    1,
    Name:  "lory",
    Phone: "8618688888888",
    Sex:   2,
    Email: "lory@gmail.com",
}

e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "email", "sex").Upsert()
```
## orm: update from data model
```golang
user := UserDO{
    Name:  "john",
    Phone: "8618699999999",
    Sex:   1,
    Email: "john@gmail.com",
}

e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Select("name", "phone", "email", "sex").Update()
```

## orm: query results into data model
```golang
user := UserDO{}

// default 'select * from ...'
e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Query()

// just select 'name' and 'phone' 
e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Select("name", "phone").Query()
```

## orm: query results into data model slice
```golang
var users []UserDO

// select id, name, phone from users limit 3
//e.Model(&user).Table(TABLE_NAME_USERS).Select("id", name", "phone").Limit(3).Query();

// select * from users limit 3
if rowsAffected, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(3).Query(); err != nil {
    log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
} else {
    log.Debugf("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
}
```

## orm: update from data model
```golang
user := UserDO{
    Id:    1, 
    Name:  "john",
    Phone: "8618699999999",
    Sex:   1,
    Email: "john@gmail.com",
}
//SQL: update users set name='john', phone='8618699999999', sex='1', email='john@gmail.com' where id='1'
if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "email", "sex").Update(); err != nil {
    log.Errorf("update data model [%+v] error [%v]", user, err.Error())
} else {
    log.Debugf("update data model [%+v] ok, rows affected [%v]", user, rowsAffected)
}
```

## raw: query results into data model 
```golang
user := UserDO{}

//SQL: select * from users where id=1
if rowsAffected, err := e.Model(&user).QueryRaw("select * from users where id=?", 1); err != nil {
    log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
} else {
    log.Debugf("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
}
```

## raw: query results into data model slice 
```golang
var users []UserDO

//SQL: select * from users where id < 5
if rowsAffected, err := e.Model(&users).QueryRaw("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
    log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
} else {
    log.Debugf("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
}
```

## raw: query results into data model map[string]string slice  
```golang
var users []map[string]string

//SQL: select * from users where id < 5
if rowsAffected, err := e.Model(&users).QueryMap("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
    log.Errorf("query into map [%+v] error [%v]", users, err.Error())
} else {
    log.Debugf("query into map [%+v] ok, rows affected [%v]", users, rowsAffected)
}
```

## raw: exec without data model
```golang
//e.ExecRaw("UPDATE %v SET name='duck' WHERE id='%v'", TABLE_NAME_USERS, 2) //it will work well as question placeholder
rowsAffected, lasteInsertId, err := e.ExecRaw("UPDATE users SET name=? WHERE id=?", "duck", 1)
if err != nil {
    log.Errorf("exec raw sql error [%v]", err.Error())
} else {
    log.Debugf("exec raw sql ok, rows affected [%v] last insert id [%v]", rowsAffected, lasteInsertId)
}
```

## save data to cache by id or index 
just for orm [insert/upsert/update]
```golang

```

## change primary key name
```golang
e.SetPkName("uuid")
```

## use cache when orm query/update/insert/upsert
```golang

e := sqlca.NewEngine()
e.Debug(true) //debug mode on

// open database driver (requred)
e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")

// open redis driver for cache (requred)
e.Open("redis://127.0.0.1:6379/cluster?db=0", 3600) //redis standalone mode

user := UserDO{
    Id:    0, 
    Name:  "john",
    Phone: "8618699999999",
    Sex:   1,
    Email: "john@gmail.com",
}
e.Model(&user).Table(TABLE_NAME_USERS).Cache("phone").Insert()
```

## tx (TxBegin...TxGet...TxExec...TxCommit)

```golang

func SqlcaTxGetAndExec(e *sqlca.Engine) (err error) {

	var tx *sqlca.Engine
	//transaction: select user id form users where phone is '8618600000000' and update users disable to 1 by user id
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("TxBegin error [%v]", err.Error())
		return
	}

	var UserId int32

	//query results into base variants
	_, err = tx.TxGet(&UserId, "SELECT id FROM users WHERE phone='%v'", "8618600000000")
	if err != nil {
		log.Errorf("TxGet error %v", err.Error())
		_ = tx.TxRollback()
		return
	}
	var lastInsertId, rowsAffected int64
	if UserId == 0 {
		log.Warnf("select id users by phone number but user not exist")
		_ = tx.TxRollback()
		return
	}
	log.Debugf("base variant of user id [%+v]", UserId)
	lastInsertId, rowsAffected, err = tx.TxExec("UPDATE users SET disable=? WHERE id=?", 1, UserId)
	if err != nil {
		log.Errorf("TxExec error %v", err.Error())
		_ = tx.TxRollback()
		return
	}
	log.Debugf("user id [%v] disabled, last insert id [%v] rows affected [%v]", UserId, lastInsertId, rowsAffected)

	//query results into a struct object or slice
	var dos []UserDO
	_, err = tx.TxGet(&dos, "SELECT * FROM users WHERE disable=1")
	if err != nil {
		log.Errorf("TxGet error %v", err.Error())
		_ = tx.TxRollback()
		return
	}
	for _, do := range dos {
		log.Debugf("struct user data object [%+v]", do)
	}

	err = tx.TxCommit()
	return
}
```

## index record to cache
```golang
user := UserDO{
    Id:    1,
    Name:  "john3",
    Phone: "8615011111114",
    Sex:   1,
    Email: "john3@gmail.com",
}

//SQL: update users set name='john3', phone='8615011111114', sex='1', email='john3@gmail.com' where id='1'
//index: name, phone
//redis key:  sqlx:cache:[table]:[column]:[column value]
if rowsAffected, err := e.Model(&user).
    Table(TABLE_NAME_USERS).
    Select("name", "phone", "email", "sex").
    Cache("name", "phone").
    Update(); err != nil {
    log.Errorf("update data model [%+v] error [%v]", user, err.Error())
} else {
    log.Debugf("update data model [%+v] ok, rows affected [%v]", user, rowsAffected)
}
``` 

## attach exist sqlx db instance
```golang
// db is a exist sqlx.DB instance
e := sqlca.NewEngine().Attach(db)
```

## set cache update before db update
```golang
e.SetCacheBefore(true)
```