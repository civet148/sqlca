package main

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"github.com/civet148/sqlca/v2/models"
	"github.com/civet148/sqlca/v2/types"
	"time"
	//_ "github.com/mattn/go-sqlite3" //import go sqlite3 if you want
)

const (
	TABLE_NAME_USERS = "users"
)

var urls = []string{
	"root:123456@tcp(192.168.2.9:3306)/test?charset=utf8mb4", //raw mysql DSN (default 'mysql' if scheme is not specified)
	//"mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4",
	//"postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable",
	//"mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false",
}

//func (u *UserData) Scan(src interface{}) (err error) {
//	//log.Debugf("UserData -> scan from string...[%+v]", src)
//	if err =  json.Unmarshal(src.(string), u); err != nil {
//		log.Errorf(err.Error())
//		return
//	}
//	return
//}

func main() {

	//e.Open("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4") // open with raw mysql DSN
	//e.Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4", &sqlca.Options{Max: 20, Idle: 2})             //MySQL master
	//e.Open("mysql://root:123456@127.0.0.1:3307/test?charset=utf8mb4", sqlca.Options{Max: 20, Idle: 5, Slave: true}) //MySQL slave
	//e.Open("postgres://root:`~!@#$%^&*()-_=+@127.0.0.1:5432/test?sslmode=enable") //postgres
	//e.Open("sqlite:///var/lib/test.db") //sqlite3
	//e.Open("mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false") //windows MS SQLSERVER

	//connect database directly
	for _, url := range urls {
		e, err := sqlca.NewEngine(url)
		if err != nil {
			log.Errorf(err.Error())
			continue
		}
		e.Debug(true) //debug on
		//e.SetLogFile("sqlca.log")
		e.SlowQuery(true, 0)
		Direct(e)
		log.Infof("------------------------------------------------------------------------------------------------------------------------------------------------------------")
	}

	//connect database through SSH tunnel
	//for _, url := range urls {
	//	e := sqlca.NewEngine(url, &sqlca.Options{
	//		Debug: true,
	//		SSH: &sqlca.SSH{
	//			User:     "ubuntu",          //SSH server login account name
	//			Host:     "192.168.124.162", //SSH server host (default port 22)
	//			Password: "123456",          //SSH server password
	//			//PrivateKey: "path/to/private/key.pem", //private key of SSH
	//		},
	//	})
	//	//e.SetLogFile("sqlca.log")
	//	SSHTunnel(e)
	//	log.Infof("------------------------------------------------------------------------------------------------------------------------------------------------------------")
	//}
	log.Info("program exit...")
}

func SSHTunnel(e *sqlca.Engine) {
	var users []models.UsersDO
	if _, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(10).Query(); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("SSH tunnel query users [%+v]", users)
}

// connect database directly
func Direct(e *sqlca.Engine) {
	SwitchDatabase(e)
	OrmInsertByModel(e)
	OrmUpsertByModel(e)
	OrmUpdateByModel(e)
	OrmQueryIntoModel(e)
	OrmQueryExcludeIntoModel(e)
	OrmQueryIntoModelSlice(e)
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
	TxWrapper(e)
	CustomTag(e)
	BuiltInTypesUpdate(e)
	DuplicateUpdateGetId(e)
	Count(e)
	CaseWhen(e)
	UpdateByMap(e)
	NearBy(e)
	MySqlJsonQuery(e)
	CustomizeUpsert(e)
	JoinQuery(e)
	NilPointerQuery(e)
	JsonStructQuery(e)
	BuiltInSliceQuery(e)
	QueryJSON(e)
	BoolConvert(e)
	QueryEx(e)
	OrmQueryToDecimal(e)
	OrmQueryLike(e)
}

func SwitchDatabase(e *sqlca.Engine) {
	type Block struct {
		Id  int32  `json:"id" db:"id"`
		Cid string `json:"cid" db:"cid"`
	}
	db, err := e.Use("stos-manager")
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var blocks []*Block
	_, err = db.Model(&blocks).Table("block").Limit(5).Query()
	for _, b := range blocks {
		log.Infof("block [%+v]", b)
	}
}

func OrmQueryToDecimal(e *sqlca.Engine) {
	var d sqlca.Decimal
	c, err := e.Model(&d).Table(TABLE_NAME_USERS).Count("id").Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("query result count %d decimal %s", c, d)
}

func OrmInsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()
	var users = make([]models.UsersDO, 0)
	for i := 0; i < 3; i++ {
		users = append(users, models.UsersDO{
			//Id:    0,
			Name:    "lory",
			Phone:   "+8618682371690",
			Sex:     1,
			Balance: sqlca.NewDecimal("123.456"),
			Email:   "lory@example.com",
			Disable: true,
			ExtraData: &models.UserData{
				Age:    32,
				Height: 183,
				Female: true,
			},
		})
	}

	log.Debugf("users [%+v]", users)

	//insert from model except 'created_at', 'updated_at' column
	if lastInsertId, err := e.Model(&users).Table(TABLE_NAME_USERS).Exclude("created_at", "updated_at").Insert(); err != nil {
		log.Errorf("insert data model [%+v] error [%v]", users, err.Error())
	} else {
		log.Infof("insert data model [%+v] exclude created_at and updated_at ok, last insert id [%v]", users, lastInsertId)
	}
}

func OrmUpsertByModel(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()
	user := models.UsersDO{
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

	user := models.UsersDO{
		Id:        1,
		Name:      "john",
		Phone:     "8618699999999",
		Sex:       1,
		Email:     "john@gmail.com",
		Disable:   true,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	//SQL: update users set name='john', phone='8618699999999', sex='1', email='john@gmail.com' where id='1'... exclude created_at and updated_at
	if rowsAffected, err := e.Model(&user).Table(TABLE_NAME_USERS).Select("phone", "sex", "name", "email").Id(1).Exclude("created_at", "updated_at").Update(); err != nil {
		log.Errorf("update data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("update data model [%+v] exclude created_at and updated_at ok, rows affected [%v]", user, rowsAffected)
	}
}

func OrmQueryIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := &models.UsersDO{}

	// select * from users where id=1
	if rowsAffected, err := e.Model(user).Table(TABLE_NAME_USERS).Id(15).Query(); err != nil {
		log.Errorf("query into data model [%+v] error [%v]", user, err.Error())
	} else {
		log.Infof("query into user model [%+v] ok, rows affected [%v]]", user, rowsAffected)
		log.Json("user.ExtraData", user.ExtraData)
	}

	type UserInfo struct {
		Users []*models.UsersDO `json:"users" db:"users"`
	}
}

func OrmQueryExcludeIntoModel(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	user := &models.UsersDO{}

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

	var users []*models.UsersDO

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

	user := models.UsersDO{}

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

	var users []models.UsersDO

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

func OrmSelectMultiTable(e *sqlca.Engine) {

	log.Enter()
	defer log.Leave()

	type UserClass struct {
		UserId  int32  `db:"user_id"`
		Phone   string `db:"phone"`
		ClassNo string `db:"class_no"`
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

	user := models.UsersDO{
		Id: 1000,
	}
	//delete from data model
	if rows, err := e.Model(&user).Table(TABLE_NAME_USERS).Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}

	//delete from where condition (without data model)
	if rows, err := e.Model(nil).Table(TABLE_NAME_USERS).Where("id > 1001").Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}

	//delete from primary key 'id' and value (without data model)
	if rows, err := e.Model(nil).Table(TABLE_NAME_USERS).Id(1002).Where("disable=1").Delete(); err != nil {
		log.Errorf("delete from table error [%v]", err.Error())
	} else {
		log.Infof("delete from table ok, affected rows [%v]", rows)
	}
}

func OrmInCondition(e *sqlca.Engine) {
	log.Enter()
	defer log.Leave()

	var users []models.UsersDO
	//SQL: select * from users where id > 2 and id in (1,3,6,7) and disable in (0,1)
	if rows, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		//Where("id > 2").
		In("id", []int32{1, 3, 6, 7}).
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

	var users []models.UsersDO
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

	var user = models.UsersDO{
		Disable: true,
	}
	if _, err := e.Model(&user).Table(TABLE_NAME_USERS).Update(); err != nil { // expect return error
		log.Errorf("%v", err.Error())
	}
	if _, err := e.Model(&user).Table(TABLE_NAME_USERS).Delete(); err != nil { // expect return error
		log.Errorf("%v", err.Error())
	}
}

func OrmToSQL(e *sqlca.Engine) {
	user := models.UsersDO{
		Id:    1,
		Name:  "john3",
		Phone: "8615011111114",
		Sex:   1,
		Email: "john3@gmail.com",
	}
	log.Infof("ToSQL insert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(types.OperType_Insert))
	log.Infof("ToSQL upsert [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(types.OperType_Upsert))
	log.Infof("ToSQL query [%v]", e.Model(&user).Table(TABLE_NAME_USERS).Select("name", "phone", "sex", "email").ToSQL(types.OperType_Query))
	log.Infof("ToSQL delete [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(types.OperType_Delete))
	log.Infof("ToSQL for update [%v]", e.Model(&user).Table(TABLE_NAME_USERS).ToSQL(types.OperType_ForUpdate))
}

func OrmGroupByHaving(e *sqlca.Engine) {
	var users []models.UsersDO
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

func TxGetExec(e *sqlca.Engine) {
	var err error
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
	var dos []models.UsersDO
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
}

func TxRollback(e *sqlca.Engine) {
	var err error
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

func TxWrapper(e *sqlca.Engine) {

	//transaction wrapper (auto rollback + auto commit)
	_ = e.TxFunc(func(tx *sqlca.Engine) error {
		//query results into a struct object or slice
		var err error
		var dos models.UsersDO
		_, err = tx.Model(&dos).
			Table("users").
			Equal("disable", 1).
			Limit(1).
			Query()
		//_, err = tx.TxGet(&dos, "SELECT * FROM users WHERE disable=1 LIMIT 1")
		if err != nil {
			log.Errorf("TxGet error %v", err.Error())
			return err
		}
		//UPDATE users SET disable=0 WHERE id='xx'
		_, err = tx.Model(0).
			Table("users").
			Id(dos.Id).
			Select("disable").
			Where("disable=1").
			Limit(1).
			Update()
		if err != nil {

			return log.Errorf("tx update error %v", err.Error())
		}
		log.Infof("user id [%d] update disable to 0 ready", dos.Id)
		//DELETE FROM users WHERE id='xx'
		_, err = tx.Model(nil).
			Table("users").
			Equal("id", dos.Id).
			Delete()
		if err != nil {
			log.Errorf("tx delete error %v", err.Error())
			return err
		}
		return nil
	})
}

func TxForUpdate(e *sqlca.Engine) {

	go func() {
		e.SlowQuery(true, 0) //设置输出慢查询执行时间，超过规定则打印告警信息（毫秒）
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
		IgnoreMe string `db:"-" json:"ignore_me"`                                    //ignore it
	}

	var users []CustomUser
	//add custom tag
	e.SetCustomTag(types.TAG_NAME_PROTOBUF, types.TAG_NAME_JSON)
	if count, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		Where("id < ?", 5).
		Query(); err != nil {
		log.Errorf("custom tag query error [%v]", err.Error())
	} else {
		log.Infof("custom tag query results %+v rows [%v]", users, count)
	}
}

func BuiltInTypesUpdate(e *sqlca.Engine) {
	//UPDATE users SET `sex`='1',`disable`='0' WHERE `id`='1'
	_, err := e.Model(1, 0).
		Table("users").
		Select("sex", "disable").
		Id(1).
		Update()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	//UPDATE users SET `sex`='2',`disable`='1' WHERE `id`='2'
	var sex = 2
	if rows, err := e.Model(&sex, 1).Table(TABLE_NAME_USERS).Id(2).Select("sex", "disable").Update(); err != nil {
		log.Error(err.Error())
	} else {
		log.Debugf("built-in type update ok, affected rows [%v]", rows)
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
	var count int64
	if _, err := e.Model(&count).
		Count("id", "id_count").
		Table(TABLE_NAME_USERS).
		Where("created_at > ?", "2020-06-01 02:03:04").
		And("disable=0").
		Query(); err != nil {

		log.Errorf("error [%v]", err.Error())
	} else {
		log.Infof("count = %v", count)
	}
}

func CaseWhen(e *sqlca.Engine) {

	var users []models.UsersDO
	if _, err := e.Model(&users).
		Table(TABLE_NAME_USERS).
		Select("id", "name", "phone", "sex").
		Case("male", "sex=1").
		Case("female", "sex=2").
		Else("unknown").
		End("sex_name").
		Where("created_at > ?", "2020-06-01 02:03:04").
		And("disable=0").
		Limit(5).
		Query(); err != nil {

		log.Errorf("error [%v]", err.Error())
	} else {
		log.Infof("users %+v", users)
	}
}

func UpdateByMap(e *sqlca.Engine) {

	//only map[string]interface{} and map[string]string
	updates := map[string]interface{}{
		"sex":  "4",
		"name": "name 2",
	}
	//UPDATE users SET `name`='name 2',`sex`='4' WHERE `id`='2'
	if _, err := e.Model(&updates).Table(TABLE_NAME_USERS).Id(2).Update(); err != nil {
		log.Errorf("update by map [%+v] error [%s]", updates, err.Error())
		return
	}
}

func NearBy(e *sqlca.Engine) {
	/*
			-- 创建表语句

			CREATE TABLE `t_address` (
			  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
			  `lng` double(11,7) DEFAULT NULL,
			  `lat` double(11,7) DEFAULT NULL,
			  `name` char(80) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
			  PRIMARY KEY (`id`)
			) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

			-- SQL语句

			SELECT
			    *,
			(6371 * ACOS(COS(RADIANS(lat)) * COS(RADIANS(28.8039097230)) * COS(RADIANS(121.5619236231) - RADIANS(lng)) + SIN(RADIANS(lat)) * SIN(RADIANS(28.8039097230)))) AS distance
			FROM  t_address
			WHERE 1=1
			HAVING distance < 113.100
			ORDER BY distance  LIMIT 10

		    ---------------------------------------------------------------
			    id     lng           lat      name          distance
			------  -----------  ----------  ------    --------------------
			     1  121.5619236  29.8079889  addr1    	111.64851043187684
			     2  121.5719236  29.8179889  addr2  	112.76462800646821
	*/

	type NearByDO struct {
		Id       int     `db:"id"`
		Lng      float64 `db:"lng"`
		Lat      float64 `db:"lat"`
		Name     string  `db:"name"`
		Distance float64 `db:"distance"`
	}
	var dos []NearByDO
	rows, err := e.Model(&dos).
		Table("t_address").
		NearBy("lng", "lat", "distance", 121.5619236231, 28.8039097230, 113.100).
		OrderBy("distance").
		Limit(10).
		Query()
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Infof("nearby rows [%d] dos -> %+v ", rows, dos)
	}
}

type JsonsDo struct {
	Id        int32  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Sex       int32  `json:"sex" db:"sex"`
	Height    int32  `json:"height" db:"height"`
	Weight    int32  `json:"weight" db:"weight"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func MySqlJsonQuery(e *sqlca.Engine) {
	//SELECT id, NAME, user_data->>'$.height' AS height, user_data->>'$.weight' AS weight FROM jsons LIMIT 5
	var dos []JsonsDo

	_, err := e.Model(&dos).
		Table("jsons").
		Select("id", "name", "user_data->>'$.height' AS height", "user_data->>'$.weight' AS weight").
		Limit(5).
		Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("%+v", dos)

	//SELECT id, NAME, user_data->>'$.height' AS height, user_data->>'$.weight' AS weight FROM jsons WHERE user_data->>'$.auth_code'=2 AND user_data->>'$.height'=165 AND user_data->>'$.auth_no'="450232200910230012"
	dos = nil
	_, err = e.Model(&dos).
		Table("jsons").
		Select("id", "name", "user_data->>'$.height' AS height", "user_data->>'$.weight' AS weight").
		Where("user_data->>'$.auth_code'=2 AND user_data->>'$.height'=165 AND user_data->>'$.auth_no'=\"450232200910230012\"").
		//Limit(5).
		Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("%+v", dos)
}

func CustomizeUpsert(e *sqlca.Engine) {

	var do = models.UsersDO{
		Id:      1,
		Name:    "customize upsert",
		Phone:   "8617923930922",
		Sex:     1,
		Email:   "civet148@gmail.com",
		Disable: false,
		Balance: sqlca.NewDecimal(6.66),
	}
	//INSERT INTO users(id, name, phone, sex, email, disable, balance)
	//VALUES('1', "customize upsert", "8617923930921", '1', "civet148@gmail.com", '0', '6.66')
	//ON DUPLICATE KEY UPDATE balance=balance+VALUES(balance)

	var strCustomUpdates string
	adapter := e.GetAdapter()
	switch adapter {
	case types.AdapterSqlx_MySQL:
		strCustomUpdates = "balance=balance+VALUES(balance)"
	case types.AdapterSqlx_Postgres:
		strCustomUpdates = "balance=users.balance+excluded.balance"
	}
	if _, err := e.Model(&do).
		Table(TABLE_NAME_USERS).
		OnConflict("id"). // only postgres required
		Upsert(strCustomUpdates); err != nil {

		log.Error(err.Error())
		return
	}
}

func JoinQuery(e *sqlca.Engine) {
	var users []*models.UsersDO
	_, err := e.Model(&users).Table(TABLE_NAME_USERS).JsonEqual(models.USERS_COLUMN_EXTRA_DATA, "age", 20).Query()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Json(users)
}

func NilPointerQuery(e *sqlca.Engine) {
	var user *models.UsersDO //nil pointer of UserDO （pass pointer address to query)
	c, err := e.Model(&user).
		Select("a.*").
		Table("users a").
		Where("a.id <= 9").
		Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("count [%d] UserClass data [%+v]", c, user)
}

// query result returns json string
func QueryJSON(e *sqlca.Engine) {

	var users []models.UsersDO
	strJsonResults, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(10).QueryJson()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("user results to JSON %s", strJsonResults)
}

func QueryEx(e *sqlca.Engine) {

	var users []models.UsersDO
	rows, total, err := e.Model(&users).Table(TABLE_NAME_USERS).Limit(2).QueryEx()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("user rows [%d] total [%d]", rows, total)
}

func JsonStructQuery(e *sqlca.Engine) {

	type JsonsDO struct {
		Id   int32  `db:"id"`
		Name string `db:"name"`
		//UserData *models.UserData `db:"user_data"` //column 'user_data' is a json string in table, eg. {"age": 18, "female": false, "height": 178}
		UserData []*models.UserData `db:"user_data"` //column 'user_data' is a json string in table, eg. [{"age":18, "female":false, "height":178},{"age":28, "female":true, "height": 162}]
	}
	var dos []*JsonsDO
	/* -- SELECT  id, name, user_data FROM jsons  WHERE id='1' --

	    id  name       sex   user_data                                                                                     created_at           updated_at
	------  ------  ------  ----------------------------------------------------------------------------------------  -------------------  ---------------------
	     1  jhon         1  {"age": 18, "female": false, "height": 178}                                               2020-10-24 11:35:11    2020-11-17 14:35:16
	*/
	_, err := e.Model(&dos).
		Select("id", "name", "user_data").
		Table("jsons").
		Where("id=?", 1).
		Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Json(dos)
	//log.Infof("JsonsDO [%+v] UserData [%+v]", do, do.UserData)
}

func BuiltInSliceQuery(e *sqlca.Engine) {
	var err error
	var idList []int32
	//SELECT id FROM users WHERE id < 10  [fetch to a int32 slice]
	if _, err = e.Model(&idList).Table(TABLE_NAME_USERS).Select("id").Where("id < 10").Query(); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("id list %+v", idList)
}

func BoolConvert(e *sqlca.Engine) {
	user := &models.UsersDO{}
	e.Model(user).Table(TABLE_NAME_USERS).Limit(1).Query()
	log.Info("user query [%+v]", user)

	if user.Disable {
		user.Disable = false
	} else {
		user.Disable = true
	}
	e.Model(user).Table(TABLE_NAME_USERS).Upsert()
	log.Info("user upsert [%+v] ", user)
}

func OrmQueryLike(e *sqlca.Engine) {
	var users []*models.UsersDO
	if _, err := e.Model(&users).Table(TABLE_NAME_USERS).Like("name", "oh").Query(); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("users %+v", users)
}
