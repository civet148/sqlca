package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
)

const (
	TABLE_NAME_USERS = "users"
)

type UserDO struct {
	Id    int32  `db:"id"`    // int(11) NOT NULL AUTO_INCREMENT COMMENT 'auto inc id',
	Name  string `db:"name"`  // varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
	Phone string `db:"phone"` // varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
	Sex   int8   `db:"sex"`   // tinyint(1) NOT NULL DEFAULT 1,
	Email string `db:"email"` // varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
}

func main() {

	e := sqlca.NewEngine(true)
	e.Debug(true) //debug on

	e.Open("redis://127.0.0.1:6379/cluster?db=0", 3600) //redis alone mode
	//e.Open(sqlca.AdapterCache_Redis, "redis://123456@127.0.0.1:6379/cluster?db=0&replicate=127.0.0.1:6380,127.0.0.1:6381") //redis cluster mode

	e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")
	//e.Open("postgres://root:`~!@#$%^&*()-_=+@127.0.0.1:5432/test?sslmode=enable")
	//e.Open("sqlite:///var/lib/test.db")
	//e.Open("mssql://sa:123456@127.0.0.1:1433/test?name=test&windows=false")

	//OrmInsertByModel(e)
	//OrmUpsertByModel(e)
	//OrmUpdateByModel(e)
	//OrmQueryIntoModel(e)
	//OrmQueryIntoModelSlice(e)
	//RawQueryIntoModel(e)
	//RawQueryIntoModelSlice(e)
	//RawQueryIntoMap(e)
	//RawExec(e)
	//TxExec(e)
	RawTxExec(e)

	log.Info("program exit...")
}

func OrmInsertByModel(e *sqlca.Engine) {

	user := UserDO{
		//Id:    0,
		Name:  "admin",
		Phone: "8618600000000",
		Sex:   1,
		Email: "admin@golang.org",
	}
	if lastInsertId, err := e.Model(&user).Table(TABLE_NAME_USERS).Insert(); err != nil {
		log.Errorf("insert data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("insert data model [%+v] ok, last insert id [%v]", user, lastInsertId)
	}
}

func OrmUpsertByModel(e *sqlca.Engine) {
	user := UserDO{
		Id:    1,
		Name:  "lory",
		Phone: "8618688888888",
		Sex:   2,
		Email: "lory@gmail.com",
	}
	if lastInsertId, err := e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "email", "sex").Upsert(); err != nil {
		log.Errorf("upsert data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("upsert data model [%+v] ok, last insert id [%v]", user, lastInsertId)
	}
}

func OrmUpdateByModel(e *sqlca.Engine) {
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
}

func OrmQueryIntoModel(e *sqlca.Engine) {
	user := UserDO{}

	//SQL: select id, name, phone from users where id=1
	//e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Select("id", "name", "phone").Query();

	// select * from users where id=1
	if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmQueryIntoModelSlice(e *sqlca.Engine) {

	var users []UserDO

	//SQL: select id, name, phone from users limit 3
	//e.Model(&user).Table(TABLE_NAME_USERS).Select("id", "name", "phone").Limit(3).Query();

	//SQL: select * from users limit 3
	if rowsAffected, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(3).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func RawQueryIntoModel(e *sqlca.Engine) {
	user := UserDO{}

	//SQL: select * from users where id=1
	if rowsAffected, err := e.Model(&user).QueryRaw("select * from users where id=?", 1); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func RawQueryIntoModelSlice(e *sqlca.Engine) {
	var users []UserDO

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryRaw("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func RawQueryIntoMap(e *sqlca.Engine) {
	var users []map[string]string

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryMap("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into map [%+v] error [%v]", users, err.Error())
	} else {
		log.Debugf("query into map [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func RawExec(e *sqlca.Engine) {

	//e.ExecRaw("UPDATE %v SET name='duck' WHERE id='%v'", TABLE_NAME_USERS, 2) //it will work well as question placeholder
	rowsAffected, lasteInsertId, err := e.ExecRaw("UPDATE users SET name=? WHERE id=?", "duck", 1)
	if err != nil {
		log.Errorf("exec raw sql error [%v]", err.Error())
	} else {
		log.Debugf("exec raw sql ok, rows affected [%v] last insert id [%v]", rowsAffected, lasteInsertId)
	}
}

func TxExec(e *sqlca.Engine) {
	user1 := UserDO{
		//Id:    0,
		Name:  "user1",
		Phone: "8618600000001",
		Sex:   1,
		Email: "user1@hotmail.com",
	}

	user2 := UserDO{
		//Id:    0,
		Name:  "user2",
		Phone: "8618600000002",
		Sex:   1,
		Email: "user2@hotmail.com",
	}
	tx1 := e.Model(&user1).Table(TABLE_NAME_USERS).ToTxInsert()
	tx2 := e.Model(&user2).Table(TABLE_NAME_USERS).ToTxInsert()
	if err := e.Tx(tx1, tx2); err != nil {
		log.Errorf("tx error [%v]", err.Error())
	} else {
		log.Debugf("tx ok")
	}
}

func RawTxExec(e *sqlca.Engine) {

	tx1 := "INSERT INTO users (`name`,`phone`,`sex`,`email`) VALUES ('user3','8618600000003','1','user3@hotmail.com')"
	tx2 := "INSERT INTO users (`name`,`phone`,`sex`,`email`) VALUES ('user4','8618600000004','2','user4@hotmail.com')"

	if err := e.TxRaw(tx1, tx2); err != nil {
		log.Errorf("tx raw error [%v]", err.Error())
	} else {
		log.Debugf("tx raw ok")
	}
}
