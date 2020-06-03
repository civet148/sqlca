package mysql

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"time"
)

const (
	TABLE_NAME_USERS = "users"
)

type UserDO struct {
	Id        int32         `db:"id"`
	Name      string        `db:"name"`
	Phone     string        `db:"phone"`
	Sex       int8          `db:"sex"`
	Email     string        `db:"email"`
	Disable   int8          `db:"disable"`
	Balance   sqlca.Decimal `db:"balance"`
	CreatedAt string        `db:"created_at" sqlca:"readonly"`
	IgnoreMe  string        `db:"-"`
}

type ClassDo struct {
	Id        int32  `db:"id"`
	UserId    int32  `db:"user_id"`
	ClassNo   string `db:"class_no"`
	CreatedAt string `db:"created_at" sqlca:"readonly"`
	IgnoreMe  string `db:"-"`
}

func Benchmark() {

	e := sqlca.NewEngine(true)
	e.Debug(true) //debug on

	e.Open("redis://127.0.0.1:6379", 3600) //redis alone mode
	//e.Open("redis://123456@127.0.0.1:6379/cluster?db=0&replicate=127.0.0.1:6380,127.0.0.1:6381") //redis cluster mode

	e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4") //MySQL
	//e.Open("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4") // open with raw mysql DSN
	//e.Open("postgres://root:`~!@#$%^&*()-_=+@127.0.0.1:5432/test?sslmode=enable") //postgres
	//e.Open("sqlite:///var/lib/test.db") //sqlite3
	//e.Open("mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false") //windows MS SQLSERVER

	MYSQL_OrmInsertByModel(e)
	MYSQL_OrmUpsertByModel(e)
	MYSQL_OrmUpdateByModel(e)
	MYSQL_OrmQueryIntoModel(e)
	MYSQL_OrmQueryIntoModelSlice(e)
	MYSQL_OrmUpdateIndexToCache(e)
	MYSQL_OrmSelectMultiTable(e)
	MYSQL_OrmDeleteFromTable(e)
	MYSQL_OrmInCondition(e)
	MYSQL_OrmFind(e)
	MYSQL_OrmWhereRequire(e)
	MYSQL_OrmToSQL(e)
	MYSQL_OrmGroupByHaving(e)
	MYSQL_RawQueryIntoModel(e)
	MYSQL_RawQueryIntoModelSlice(e)
	MYSQL_RawQueryIntoMap(e)
	MYSQL_RawExec(e)
	MYSQL_TxGetExec(e)
	MYSQL_TxRollback(e)
	MYSQL_TxForUpdate(e)
	MYSQL_CustomTag(e)
	MYSQL_BaseTypesUpdate(e)
	MYSQL_DuplicateUpdateGetId(e)
}

func MYSQL_OrmInsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	user := UserDO{
		//Id:    0,
		Name:    "admin",
		Phone:   "8618600000000",
		Sex:     1,
		Balance: sqlca.NewDecimal("123.45"),
		Email:   "admin@golang.org",
	}
	log.Debugf("user [%+v]", user)
	if lastInsertId, err := e.Model(&user).Table(TABLE_NAME_USERS).Insert(); err != nil {
		log.Errorf("insert data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("insert data model [%+v] ok, last insert id [%v]", user, lastInsertId)
	}
}

func MYSQL_OrmUpsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()
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

func MYSQL_OrmUpdateByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	user := UserDO{
		Id:      1,
		Name:    "john",
		Phone:   "8618699999999",
		Sex:     1,
		Email:   "john@gmail.com",
		Disable: 1,
	}

	//SQL: update users set name='john', phone='8618699999999', sex='1', email='john@gmail.com' where id='1'
	if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Update(); err != nil {
		log.Errorf("update data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("update data model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func MYSQL_OrmQueryIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := &UserDO{}

	//SQL: select id, name, phone from users where id=1
	//e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Select("id", "name", "phone").Query();

	// select * from users where id=1
	if rowsAffected, err := e.Model(user).Table(TABLE_NAME_USERS).Id(1).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func MYSQL_OrmQueryIntoModelSlice(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	var users []*UserDO

	//SQL: select id, name, phone from users limit 3
	//e.Model(&user).Table(TABLE_NAME_USERS).Select("id", "name", "phone").Limit(3).Query();

	//SQL: select * from users limit 3
	if rowsAffected, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(3).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {

		if len(users) == 0 {
			log.Errorf("query into model failed, rows affected [%v]", rowsAffected)
		} else {
			for i, v := range users {
				log.Debugf("query into model slice of [%v]*User [%+v] ", i, v)
			}
		}
	}
}

func MYSQL_RawQueryIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := UserDO{}

	//SQL: select * from users where id=1
	if rowsAffected, err := e.Model(&user).QueryRaw("select * from users where id=?", 1); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func MYSQL_RawQueryIntoModelSlice(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	var users []UserDO

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryRaw("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Debugf("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func MYSQL_RawQueryIntoMap(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	var users []map[string]string

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryMap("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into map [%+v] error [%v]", users, err.Error())
	} else {
		log.Debugf("query into map [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func MYSQL_RawExec(e *sqlca.Engine) {

	//e.ExecRaw("UPDATE %v SET name='duck' WHERE id='%v'", TABLE_NAME_USERS, 2) //it will work well as question placeholder
	rowsAffected, lasteInsertId, err := e.ExecRaw("UPDATE users SET name=? WHERE id=?", "duck", 1)
	if err != nil {
		log.Errorf("exec raw sql error [%v]", err.Error())
	} else {
		log.Debugf("exec raw sql ok, rows affected [%v] last insert id [%v]", rowsAffected, lasteInsertId)
	}
}

func MYSQL_OrmUpdateIndexToCache(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

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
		Distinct().
		Select("name", "phone", "email", "sex").
		Cache("name", "phone").
		Update(); err != nil {
		log.Errorf("update data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Debugf("update data model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func MYSQL_OrmSelectMultiTable(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	type UserClass struct {
		UserId   int32  `db:"user_id"`
		UserName string `db:"user_name"`
		Phone    string `db:"phone"`
		ClassNo  string `db:"class_no"`
	}
	var ucs []UserClass
	//SQL: SELECT a.*, b.class_no FROM users a, classes b WHERE a.id=b.user_id AND a.id=3
	_, err := e.Model(&ucs).
		Distinct().
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
}

func MYSQL_OrmDeleteFromTable(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

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
	if rows, err := e.Table(TABLE_NAME_USERS).Where("id > 1001").Delete(); err != nil {
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
}

func MYSQL_OrmInCondition(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	var users []UserDO
	//SQL: select * from users where id > 2 and id in (1,3,6,7) and disable in (0,1)
	if rows, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		Select("*").
		Where("id > 2").
		In("id", 1, 3, 6, 7).
		In("disable", 0, 1).
		Query(); err != nil {
		log.Errorf("select from table by in condition error [%v]", err.Error())
	} else {
		log.Debugf("select from table by in condition ok, affected rows [%v]", rows)
	}
}

func MYSQL_OrmFind(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	var users []UserDO
	if rows, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		Find(map[string]interface{}{
			"id":      1,
			"disable": 1,
		}); err != nil {
		log.Errorf("select from table by find condition error [%v]", err.Error())
	} else {
		log.Debugf("select from table by find condition ok, affected rows [%v] users %+v", rows, users)
	}
}

func MYSQL_OrmWhereRequire(e *sqlca.Engine) {

	var user = UserDO{
		Disable: 2,
	}
	if _, err := e.Model(&user).Table(TABLE_NAME_USERS).Update(); err != nil { // expect return error
		log.Errorf("%v", err.Error())
	}
	if _, err := e.Model(&user).Table(TABLE_NAME_USERS).Delete(); err != nil { // expect return error
		log.Errorf("%v", err.Error())
	}
}

func MYSQL_OrmToSQL(e *sqlca.Engine) {
	user := UserDO{
		Id:    1,
		Name:  "john3",
		Phone: "8615011111114",
		Sex:   1,
		Email: "john3@gmail.com",
	}
	log.Debugf("ToSQL insert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(sqlca.OperType_Insert))
	log.Debugf("ToSQL upsert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(sqlca.OperType_Upsert))
	log.Debugf("ToSQL query [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(sqlca.OperType_Query))
	log.Debugf("ToSQL delete [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(sqlca.OperType_Delete))
}

func MYSQL_OrmGroupByHaving(e *sqlca.Engine) {
	var users []UserDO
	rows, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		GroupBy("id", "name").
		Having("id>?", 1).
		OrderBy().
		Asc("name").
		Desc("created_at").
		Query()
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Infof("rows [%v] users [%+v]", rows, users)
	}
}

func MYSQL_TxGetExec(e *sqlca.Engine) (err error) {
	log.Enter()
	defer log.Leave()

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
	_, err = tx.TxGet(&dos, "SELECT * FROM users WHERE disable=1 LIMIT 5")
	if err != nil {
		log.Errorf("TxGet error %v", err.Error())
		_ = tx.TxRollback()
		return
	}
	for _, do := range dos {
		log.Debugf("struct user data object [%+v]", do)
	}

	if err = tx.TxCommit(); err != nil {
		log.Errorf("TxCommit error [%v]", err.Error())
		return
	}
	return
}

func MYSQL_TxRollback(e *sqlca.Engine) (err error) {

	log.Enter()
	defer log.Leave()

	var tx *sqlca.Engine
	//transaction: insert and rollback
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("TxBegin error [%v]", err.Error())
		return
	}

	_, _, err = tx.TxExec("INSERT INTO users(id, name, phone, sex, email) VALUES(1, 'john3', '8618600000000', 2, 'john3@gmail.com')")
	if err != nil {
		log.Errorf("TxExec error %v, rollback", err.Error())
		_ = tx.TxRollback()
		return
	}

	if err = tx.TxCommit(); err != nil {
		log.Errorf("TxCommit error [%v]", err.Error())
		return
	}
	return
}

func MYSQL_TxForUpdate(e *sqlca.Engine) {

	go func() {

		if tx, err := e.TxBegin(); err != nil {
			log.Errorf("[TX1] tx begin error [%v]", err.Error())
			return
		} else {
			var id int32
			if _, err = tx.TxGet(&id, "SELECT id FROM users WHERE id=1 FOR UPDATE"); err != nil {
				log.Errorf("[TX1] tx get error [%v]", err.Error())
				tx.TxRollback()
				return
			}

			if _, _, err = tx.TxExec("UPDATE users SET name='i am tx 1' WHERE id=1"); err != nil {
				log.Errorf("[TX1] tx exec error [%v]", err.Error())
				tx.TxRollback()
				return
			}

			time.Sleep(2 * time.Second) //sleep for lock the record where id=1

			log.Infof("[TX1] id [%v] update ok", id)
			if err = tx.TxCommit(); err != nil {
				log.Errorf("[TX1] tx commit error [%v]", err.Error())
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)

	go func() {
		if tx, err := e.TxBegin(); err != nil {
			log.Errorf("[TX2] tx begin error [%v]", err.Error())
			return
		} else {
			var id int32
			if _, err = tx.TxGet(&id, "SELECT id FROM users WHERE id=1 FOR UPDATE"); err != nil {
				log.Errorf("[TX2] tx get error [%v]", err.Error())
				tx.TxRollback()
				return
			}
			if _, _, err = tx.TxExec("UPDATE users SET name='i am tx 2' WHERE id=1"); err != nil {
				log.Errorf("[TX2] tx exec error [%v]", err.Error())
				tx.TxRollback()
				return
			}
			log.Infof("[TX2] id [%v] update ok", id)
			if err = tx.TxCommit(); err != nil {
				log.Errorf("[TX2] tx commit error [%v]", err.Error())
				return
			}
		}
	}()

	time.Sleep(3 * time.Second)
}

func MYSQL_CustomTag(e *sqlca.Engine) {
	type CustomUser struct {
		Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // protobuf tag
		Name     string `json:"name"`                                                // json tag
		Phone    string `db:"phone"`                                                 // db tag
		IgnoreMe string `db:"-" json:"-"`
	}

	var users []CustomUser
	//add custom tag
	e.SetCustomTag(sqlca.TAG_NAME_PROTOBUF, sqlca.TAG_NAME_JSON)
	if count, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		Where("id < ?", 5).
		Query(); err != nil {
		log.Errorf("custom tag query error [%v]", err.Error())
	} else {
		log.Debugf("custom tag query results %+v rows [%v]", users, count)
	}
}

func MYSQL_BaseTypesUpdate(e *sqlca.Engine) {

	var sex = 3
	//var disable=4
	if rows, err := e.Model(&sex).Table(TABLE_NAME_USERS).Id(2).Select("sex", "disable").Update(); err != nil {
		log.Error(err.Error())
	} else {
		log.Debugf("base type update ok, affected rows [%v]", rows)
	}
}

func MYSQL_DuplicateUpdateGetId(e *sqlca.Engine) {
	strSQL := "INSERT INTO users(id, NAME, phone, sex) VALUE(1, 'li2','', 1) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)"
	if rowsAffected, lastInsertId, err := e.ExecRaw(strSQL); err != nil {
		log.Errorf(err.Error())
	} else {
		log.Infof("rows affected [%v] last insert id [%v] ", rowsAffected, lastInsertId)
	}
}
