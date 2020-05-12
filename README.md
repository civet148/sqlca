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


```

# data model
```golang
type UserDO struct {
    Id          int32           `db:"id"`  
    Name        string          `db:"name"`  
    Phone       string          `db:"phone"` 
    Sex         int8            `db:"sex"`   
    Email       string          `db:"email"` 
    Disable     int8            `db:"disable"`
    Balance     sqlca.Decimal   `db:"balance"`
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
if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "email", "sex").Cache("phone").Update(); err != nil {
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
just for orm [insert/upsert/update] see the example of orm update

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
//redis key:  sqlca:cache:[db]:[table]:[column]:[column value]
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

## delete from table
```golang

user := UserDO{
		Id: 1000,
}
//delete from data model
if rows, err := e.Model(&user).Table(TABLE_NAME_USERS).Delete(); err != nil {
    log.Errorf("delete from table error [%v]", err.Error())
} else {
    log.Debugf("delete from table ok, affected rows [%v]", rows)
}

//delete from where condition (without data model)
if rows, err := e.Table(TABLE_NAME_USERS).Where("id=1001").Delete(); err != nil {
    log.Errorf("delete from table error [%v]", err.Error())
} else {
    log.Debugf("delete from table ok, affected rows [%v]", rows)
}

//delete from primary key 'id' and value (without data model)
if rows, err := e.Table(TABLE_NAME_USERS).Id(1002).Where("disable=1").Delete(); err != nil {
    log.Errorf("delete from table error [%v]", err.Error())
} else {
    log.Debugf("delete from table ok, affected rows [%v]", rows)
}
```

## select from multiple tables
```golang
type UserClass struct {
    UserId   int32  `db:"user_id"`
    UserName string `db:"user_name"`
    Phone    string `db:"phone"`
    ClassNo  string `db:"class_no"`
}
var ucs []UserClass
//SQL: SELECT a.*, b.class_no FROM users a, classes b WHERE a.id=b.user_id AND a.id=3
_, err := e.Model(&ucs).
    Select("a.id as user_id", "a.name", "a.phone", "b.class_no").
    Table("users a", "classes b").
    Where("a.id=b.user_id").
    And("a.id=?", 3).
    Query()
if err != nil {
    log.Errorf("query error [%v]", err.Error())
} else {
    log.Debugf("user class info [%+v]", ucs)
}
```

## custom tag
```golang
type CustomUser struct {
    Id    int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // protobuf tag
    Name  string `json:"name"`   // json tag
    Phone string `db:"phone"`    // db tag
}

var users []CustomUser
//add custom tag
e.SetCustomTag("protobuf", "json")
if count, err := e.Model(&users).
    Table(TABLE_NAME_USERS).
    Where("id < ?", 5).
    Query(); err != nil {
    log.Errorf("custom tag query error [%v]", err.Error())
} else {
    log.Debugf("custom tag query results %+v rows [%v]", users, count)
}
```

## sqlca properties tag: readonly 
```golang
type UserDO struct {
    Id    int32  `db:"id"` 
    Name  string `db:"name"`   
    Phone string `db:"phone"`    
    CreatedAt string `db:"created_at" sqlca:"readonly"` //sqlca tag: readonly
}
```

