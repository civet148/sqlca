package sqlca

import (
	"github.com/astaxie/beego/cache"
	_ "github.com/go-sql-driver/mysql" //mysql golang driver
	"github.com/jmoiron/sqlx"
)

type Engine struct {
	db           *sqlx.DB          // sqlx instance
	cache        cache.Cache       // beego cache instance
	adapterSqlx  AdapterType       // what's adapter of sqlx
	adapterCache AdapterType       // what's adapter of cache
	modeType     ModeType          // mode: orm or raw
	operType     OperType          // operate type
	debug        bool              // debug mode [on/off]
	model        interface{}       // data model [struct object or struct slice]
	dict         map[string]string // data model db dictionary
	strTableName string            // table name
	strPkName    string            // primary key of table, default 'id'
	strWhere     string            // where condition to query or update
}

func NewEngine() *Engine {

	return &Engine{
		strPkName: DEFAULT_PRIMARY_KEY_NAME,
	}
}

func (e *Engine) Open(adapterType AdapterType, strUrl string) *Engine {

	var err error
	switch adapterType {
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres, AdapterSqlx_Sqlite, AdapterSqlx_Mssql:
		// TODO @libin open sqlx database conection
		strSchema, strDSN := e.getConnUrl(adapterType, strUrl)
		if e.db, err = sqlx.Open(strSchema, strDSN); err != nil {
			e.panic("open schema %v DSN %v original url [%v] error [%v]", strSchema, strDSN, strUrl, err.Error())
		}
		e.adapterSqlx = adapterType
	case AdapterCache_Redis, AdapterCache_Memcached, AdapterCache_Memory, AdapterCache_File:
		// TODO @libin open beego cache conection
		//e.cache = v.(cache.Cache)
		e.adapterCache = adapterType
	default:
		assert(nil, "open adapter instance type [%v] url [%s] failed", adapterType, strUrl)
	}
	return e
}

// attach from a exist sqlx or beego cache instance
func (e *Engine) Attach(adapterType AdapterType, v interface{}) *Engine {

	switch adapterType {
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres, AdapterSqlx_Sqlite, AdapterSqlx_Mssql:
		{
			assert(!isNilOrFalse(e.db), "already have a [%v] sqlx instance, attach failed", adapterType)
			e.db = v.(*sqlx.DB)
			e.adapterSqlx = adapterType
		}
	case AdapterCache_Redis, AdapterCache_Memcached, AdapterCache_Memory, AdapterCache_File:
		{
			assert(!isNilOrFalse(e.cache), "already have a [%v] beego cache instance, attach failed", adapterType)
			e.cache = v.(cache.Cache)
			e.adapterCache = adapterType
		}
	default:
		assert(nil, "adapter type [%v] attach instance failed", adapterType)
		return nil
	}
	return e
}

// get internal sqlx instance
// you can use it to do what you want
func (e *Engine) DB() *sqlx.DB {
	return e.db
}

// get internal cache instance
// you can use it to do what you want
func (e *Engine) Cache() cache.Cache {
	return e.cache
}

// debug mode on or off
// if debug on, some method will panic if your condition illegal
func (e *Engine) Debug(ok bool) {
	e.setDebug(ok)
}

// orm model
// use to get result set, support single struct object or slice [pointer type]
// notice: will clone a new engine object for orm operations
func (e *Engine) Model(v interface{}) *Engine {
	assert(v, "model is nil")
	return e.clone(v)
}

// set orm query table name
// when your struct type name is not a table name
func (e *Engine) Table(strName string) *Engine {
	assert(strName, "name is nil")
	e.setTableName(strName)
	return e
}

// set orm primary key, default named 'id'
func (e *Engine) PrimaryKey(strName string) *Engine {
	assert(strName, "name is nil")
	e.setPkName(strName)
	return e
}

// orm query
// return rows affected and error, if err is not nil must be something wrong
// Model function is must be called before call this function
func (e *Engine) Query() (rows int64, err error) {
	assert(e.model, "model is nil, please call Model function first")
	// TODO @libin Query() implement
	return
}

// orm insert
// return last insert id and error, if err is not nil must be something wrong
// Model function is must be called before call this function
func (e *Engine) Insert() (lastInsertId int64, err error) {
	assert(e.model, "model is nil, please call Model function first")
	// TODO @libin Insert() implement
	return
}

// orm insert or update if exist
// strColumns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return last insert id and error, if err is not nil must be something wrong, if your primary key is not a int/int64 type, maybe id return 0
// Model function is must be called before call this function
func (e *Engine) Upsert(strColumns ...string) (id int64, err error) {
	assert(e.model, "model is nil, please call Model function first")
	// TODO @libin Upsert() implement
	return
}

// orm update from model
// strColumns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return rows affected and error, if err is not nil must be something wrong
// Model function is must be called before call this function
func (e *Engine) Update(strColumns ...string) (rows int64, err error) {
	assert(e.model, "model is nil, please call Model function first")
	// TODO @libin Update() implement

	return
}

// use raw sql to query results
// return rows and error, if err is not nil must be something wrong
// Model function is must be called before call this function
func (e *Engine) QuerySQL(strQuery string, args ...interface{}) (rows int64, err error) {
	assert(strQuery, "query sql string is nil")
	assert(e.model, "model is nil, please call Model function first")
	// TODO @libin QuerySQL() implement

	return
}

// use raw sql to insert/update database, results can not be cached to redis/memcached/memory...
// return rows and error, if err is not nil must be something wrong
func (e *Engine) ExecSQL(strQuery string, args ...interface{}) (rows int64, err error) {
	assert(strQuery, "query sql string is nil")
	// TODO @libin ExecSQL() implement
	return
}
