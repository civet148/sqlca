package main

import (
	"fmt"
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
	UpdatedAt string        `db:"updated_at" sqlca:"readonly"`
	IgnoreMe  string        `db:"-"`
}

type ClassDo struct {
	Id        int32  `db:"id"`
	UserId    int32  `db:"user_id"`
	ClassNo   string `db:"class_no"`
	CreatedAt string `db:"created_at" sqlca:"readonly"`
	UpdatedAt string `db:"updated_at" sqlca:"readonly"`
	IgnoreMe  string `db:"-"`
}

var urls = []string{
	"mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4",
	//"mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false",
	//"postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable",
}

func main() {

	//e.Open("redis://123456@127.0.0.1:6379/cluster?db=0&replicate=127.0.0.1:6380,127.0.0.1:6381") //redis cluster mode
	//e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4", &sqlca.Options{Max: 20, Idle: 2})             //MySQL master
	//e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4", sqlca.Options{Max: 20, Idle: 5, Slave: true}) //MySQL slave
	//e.Open("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4") // open with raw mysql DSN
	//e.Open("postgres://root:`~!@#$%^&*()-_=+@127.0.0.1:5432/test?sslmode=enable") //postgres
	//e.Open("sqlite:///var/lib/test.db") //sqlite3
	//e.Open("mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false") //windows MS SQLSERVER

	for _, v := range urls {
		_ = v
		e := sqlca.NewEngine(v)
		e.Debug(true) //debug on
		Benchmark(e)
		log.Infof("")
		log.Infof("------------------------------------------------------------------------------------------------------------------------------------------------------------")
		log.Infof("")
	}

	log.Info("program exit...")
}

func Benchmark(e *sqlca.Engine) {

	e.Open("redis://127.0.0.1:6379", 3600) //redis alone mode

	OrmInsertByModel(e)
	OrmUpsertByModel(e)
	OrmUpdateByModel(e)
	OrmQueryIntoModel(e)
	OrmQueryExcludeIntoModel(e)
	OrmQueryIntoModelSlice(e)
	OrmUpdateIndexToCache(e)
	OrmSelectMultiTable(e)
	OrmDeleteFromTable(e)
	OrmInCondition(e)
	OrmFind(e)
	OrmWhereRequire(e)
	OrmToSQL(e)
	OrmGroupByHaving(e)
	RawQueryIntoModel(e)
	RawQueryIntoModelSlice(e)
	RawQueryIntoMap(e)
	RawExec(e)
	TxGetExec(e)
	TxRollback(e)
	TxForUpdate(e)
	CustomTag(e)
	BaseTypesUpdate(e)
	DuplicateUpdateGetId(e)
	Count(e)
}

func OrmInsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	user := UserDO{
		//Id:    0,
		Name:      "admin",
		Phone:     "8618600000000",
		Sex:       1,
		Balance:   sqlca.NewDecimal("123.45"),
		Email:     "admin@golang.org",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	log.Debugf("user [%+v]", user)
	if lastInsertId, err := e.Model(&user).Table(TABLE_NAME_USERS).Exclude("created_at", "updated_at").Insert(); err != nil {
		log.Errorf("insert data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("insert data model [%+v] exclude created_at and updated_at ok, last insert id [%v]", user, lastInsertId)
	}

	//bulk insert
	var users []UserDO
	for i := 0; i < 3; i++ {
		users = append(users, UserDO{
			Id:        0,
			Name:      fmt.Sprintf("name(%v)", i),
			Phone:     fmt.Sprintf("phone(%v)", i),
			Sex:       0,
			Email:     "xxx@gmail.com",
			Disable:   0,
			Balance:   sqlca.NewDecimal(i),
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	if lastInsertId, err := e.Model(&users).Table(TABLE_NAME_USERS).Exclude("email", "created_at", "updated_at").Insert(); err != nil {
		log.Errorf("bulk insert data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Infof("bulk insert data model [%+v] exclude email, created_at and updated_at ok, last insert id [%v]", users, lastInsertId)
	}
}

func OrmUpsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()
	user := UserDO{
		Id:    1,
		Name:  "lory",
		Phone: "8618688888888",
		Sex:   2,
		Email: "lory@gmail.com",
	}
	if lastInsertId, err := e.Model(&user).
		Table(TABLE_NAME_USERS).
		Select("name", "phone", "sex").
		OnConflict("id"). // only for postgres
		Exclude("email", "created_at", "updated_at").
		Upsert(); err != nil {
		log.Errorf("upsert data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("upsert data model [%+v] exclude email, created_at and updated_at ok, last insert id [%v]", user, lastInsertId)
	}
}

func OrmUpdateByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	user := UserDO{
		Id:        1,
		Name:      "john",
		Phone:     "8618699999999",
		Sex:       1,
		Email:     "john@gmail.com",
		Disable:   1,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	//SQL: update users set name='john', phone='8618699999999', sex='1', email='john@gmail.com' where id='1'... exclude created_at and updated_at
	if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Exclude("created_at", "updated_at").Update(); err != nil {
		log.Errorf("update data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("update data model [%+v] exclude created_at and updated_at ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmQueryIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := &UserDO{}

	//SQL: select id, name, phone from users where id=1
	//e.Model(&user).Table(TABLE_NAME_USERS).Id(1).Select("id", "name", "phone").Query();

	// select * from users where id=1
	if rowsAffected, err := e.Model(user).Table(TABLE_NAME_USERS).Id(1).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmQueryExcludeIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := &UserDO{}

	// select * from users where id=1 ..exclude email and disable
	if rowsAffected, err := e.Model(user).Table(TABLE_NAME_USERS).Id(1).Exclude("email", "disable").Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmQueryIntoModelSlice(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	var users []*UserDO

	//SQL: select id, name, phone from users limit 3
	//e.Model(&user).Table(TABLE_NAME_USERS).Select("id", "name", "phone").Limit(3).Query();

	//SQL: select * from users limit 3
	if rowsAffected, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(3).Slave().Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {

		if len(users) == 0 {
			log.Errorf("query into model failed, rows affected [%v]", rowsAffected)
		} else {
			for i, v := range users {
				log.Infof("query into model slice of [%v]*User [%+v] ", i, v)
			}
		}
	}
}

func RawQueryIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := UserDO{}

	//SQL: select * from users where id=1
	if rowsAffected, err := e.Model(&user).QueryRaw("select * from users where id=?", 1); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("query into model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func RawQueryIntoModelSlice(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	var users []UserDO

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryRaw("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Infof("query into model [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func RawQueryIntoMap(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	var users []map[string]string

	//SQL: select * from users where id < 5
	if rowsAffected, err := e.Model(&users).QueryMap("select * from %v where id < %v", TABLE_NAME_USERS, 5); err != nil {
		log.Errorf("query into map [%+v] error [%v]", users, err.Error())
	} else {
		log.Infof("query into map [%+v] ok, rows affected [%v]", users, rowsAffected)
	}
}

func RawExec(e *sqlca.Engine) {

	//e.ExecRaw("UPDATE %v SET name='duck' WHERE id='%v'", TABLE_NAME_USERS, 2) //it will work well as question placeholder
	rowsAffected, lasteInsertId, err := e.ExecRaw("UPDATE users SET name=? WHERE id=?", "duck", 1)
	if err != nil {
		log.Errorf("exec raw sql error [%v]", err.Error())
	} else {
		log.Infof("exec raw sql ok, rows affected [%v] last insert id [%v]", rowsAffected, lasteInsertId)
	}
}

func OrmUpdateIndexToCache(e *sqlca.Engine) {

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
		log.Infof("update data model [%+v] ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmSelectMultiTable(e *sqlca.Engine) {

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
		log.Infof("user class info [%+v]", ucs)
	}
}

func OrmDeleteFromTable(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	user := UserDO{
		Id: 1000,
	}
	//delete from data model
	if rows, err := e.Model(&user).Table(TABLE_NAME_USERS).Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}

	//delete from where condition (without data model)
	if rows, err := e.Table(TABLE_NAME_USERS).Where("id > 1001").Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}

	//delete from primary key 'id' and value (without data model)
	if rows, err := e.Table(TABLE_NAME_USERS).Id(1002).Where("disable=1").Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}
}

func OrmInCondition(e *sqlca.Engine) {
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
		Or("created_at > ?", "2020-06-01 00:00:00").
		Query(); err != nil {
		log.Errorf("select from table by in condition error [%v]", err.Error())
	} else {
		log.Infof("select from table by in condition ok, affected rows [%v]", rows)
	}
}

func OrmFind(e *sqlca.Engine) {

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
		log.Infof("select from table by find condition ok, affected rows [%v] users %+v", rows, users)
	}
}

func OrmWhereRequire(e *sqlca.Engine) {

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

func OrmToSQL(e *sqlca.Engine) {
	user := UserDO{
		Id:    1,
		Name:  "john3",
		Phone: "8615011111114",
		Sex:   1,
		Email: "john3@gmail.com",
	}
	log.Infof("ToSQL insert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(sqlca.OperType_Insert))
	log.Infof("ToSQL upsert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(sqlca.OperType_Upsert))
	log.Infof("ToSQL query [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(sqlca.OperType_Query))
	log.Infof("ToSQL delete [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(sqlca.OperType_Delete))
	log.Infof("ToSQL for update [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(sqlca.OperType_ForUpdate))
}

func OrmGroupByHaving(e *sqlca.Engine) {
	var users []UserDO
	rows, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		GroupBy("id", "name").
		Having("id>?", 1).
		OrderBy().
		Asc("name").
		Desc("created_at").
		Limit(10).
		Query()
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Infof("rows [%v] users [%+v]", rows, users)
	}
}

func TxGetExec(e *sqlca.Engine) (err error) {
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
	log.Infof("user id [%v] disabled, last insert id [%v] rows affected [%v]", UserId, lastInsertId, rowsAffected)

	//query results into a struct object or slice
	var dos []UserDO
	_, err = tx.TxGet(&dos, "SELECT * FROM users WHERE disable=1 LIMIT 5")
	if err != nil {
		log.Errorf("TxGet error %v", err.Error())
		_ = tx.TxRollback()
		return
	}
	for _, do := range dos {
		log.Infof("struct user data object [%+v]", do)
	}

	if err = tx.TxCommit(); err != nil {
		log.Errorf("TxCommit error [%v]", err.Error())
		return
	}
	return
}

func TxRollback(e *sqlca.Engine) (err error) {

	log.Enter()
	defer log.Leave()

	var tx *sqlca.Engine
	//transaction: insert and rollback
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("TxBegin error [%v]", err.Error())
		return
	}
	// tx auto rollback
	_, _, err = tx.AutoRollback().TxExec("INSERT INTO users(id, name, phone, sex, email) VALUES(1, 'john3', '8618600000000', 2, 'john3@gmail.com')")
	if err != nil {
		log.Errorf("TxExec error %v, rollback", err.Error())
		return
	}

	if err = tx.TxCommit(); err != nil {
		log.Errorf("TxCommit error [%v]", err.Error())
		return
	}
	return
}

func TxForUpdate(e *sqlca.Engine) {

	go func() {
		e.SetSlowQueryAlertTime(2000) //设置输出慢查询执行时间，超过规定则打印告警信息（毫秒）
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

func CustomTag(e *sqlca.Engine) {
	type CustomUser struct {
		Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // protobuf tag
		Name     string `json:"name,omitempty"`                                      // json tag
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
		log.Infof("custom tag query results %+v rows [%v]", users, count)
	}
}

func BaseTypesUpdate(e *sqlca.Engine) {

	var sex = 3
	//var disable=4
	if rows, err := e.Model(&sex).Table(TABLE_NAME_USERS).Id(2).Select("sex", "disable").Update(); err != nil {
		log.Error(err.Error())
	} else {
		log.Debugf("base type update ok, affected rows [%v]", rows)
	}
}

func DuplicateUpdateGetId(e *sqlca.Engine) {
	strSQL := "INSERT INTO users(id, NAME, phone, sex) VALUE(1, 'li2','', 1) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)"
	if rowsAffected, lastInsertId, err := e.ExecRaw(strSQL); err != nil {
		log.Errorf(err.Error())
	} else {
		log.Infof("rows affected [%v] last insert id [%v] ", rowsAffected, lastInsertId)
	}
}

func Count(e *sqlca.Engine) {

	if count, err := e.Model(nil).
		Table(TABLE_NAME_USERS).
		Where("created_at > ?", "2020-06-01 02:03:04").
		And("disable=0").
		Count(); err != nil {

		log.Errorf("error [%v]", err.Error())
	} else {
		log.Infof("count = %v", count)
	}
}
