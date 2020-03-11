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
CREATE DATABASE IF NOT EXISTS test DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

USE test;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

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
e.Open(sqlca.AdapterSqlx_MySQL, "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")

// open redis driver for cache (optional)
e.Open(sqlca.AdapterCache_Redis, "redis://127.0.0.1:6379/cluster?db=0", 3600) //redis standalone mode
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
just for orm [insert/upsert/update] and tx [update]
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
e.Open(sqlca.AdapterSqlx_MySQL, "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")

// open redis driver for cache (requred)
e.Open(sqlca.AdapterCache_Redis, "redis://127.0.0.1:6379/cluster?db=0", 3600) //redis standalone mode

user := UserDO{
    Id:    0, 
    Name:  "john",
    Phone: "8618699999999",
    Sex:   1,
    Email: "john@gmail.com",
}
e.Model(&user).Table(TABLE_NAME_USERS).UseCache().Insert()
```

## tx: orm and raw
```golang

```
