package sqlca

import (
	"github.com/astaxie/beego/cache"
	"github.com/jmoiron/sqlx"
)

type AdapterType int

const (
	DEFAULT_PRIMARY_KEY_NAME = "id"
)

const (
	AdapterSqlx_MySQL      AdapterType = 1 //sqlx: mysql
	AdapterSqlx_Postgres   AdapterType = 2 //sqlx: postgresql
	AdapterCache_Redis     AdapterType = 3 //cache: redis
	AdapterCache_Memcached AdapterType = 4 //cache: memcached
	AdapterCache_Memory    AdapterType = 5 //cache: memory
	AdapterCache_File      AdapterType = 6 //cache: file
)

func (a AdapterType) GoString() string {
	return a.String()
}

func (a AdapterType) String() string {

	switch a {
	case AdapterSqlx_MySQL:
		return "AdapterSqlx_MySQL"
	case AdapterSqlx_Postgres:
		return "AdapterSqlx_Postgres"
	case AdapterCache_Redis:
		return "AdapterCache_Redis"
	case AdapterCache_Memcached:
		return "AdapterCache_Memcached"
	case AdapterCache_Memory:
		return "AdapterCache_Memory"
	case AdapterCache_File:
		return "AdapterCache_File"
	default:
	}
	return "Adapter_Unknown"
}

type SqlxWhere map[string]interface{}

type Engine struct {
	db           *sqlx.DB    // sqlx instance
	cache        cache.Cache // beego cache instance
	adapterSqlx  AdapterType // what's adapter of sqlx
	adapterCache AdapterType // what's adapter of cache
	debug        bool        // debug on/off
	model        interface{} // data model of record
	strTableName string      // table name
	strPkName    string      // primary key of table, default 'id'
	strWhere     string      // where condition to query or update
}

func NewEngine() *Engine {

	return &Engine{
		strPkName: DEFAULT_PRIMARY_KEY_NAME,
	}
}

func (e *Engine) Open(adapterType AdapterType, strUrl string) *Engine {
	switch adapterType {
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres:
		// TODO @libin open sqlx database conection
		//e.db = v.(*sqlx.DB)
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
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres:
		{
			assert(e.db, "already have a [%v] instance, attach failed", adapterType)
			e.db = v.(*sqlx.DB)
			e.adapterSqlx = adapterType
		}
	case AdapterCache_Redis, AdapterCache_Memcached, AdapterCache_Memory:
		{
			assert(e.cache, "already have a [%v] instance, attach failed", adapterType)
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
	return e.clone(v)
}

// set orm query table name
// when your struct type name is not a table name
func (e *Engine) Table(strName string) *Engine {
	e.setTableName(strName)
	return e
}

// set orm primary key, default named 'id'
func (e *Engine) PrimaryKey(strName string) *Engine {
	e.setPkName(strName)
	return e
}

// orm query
// return rows affected and error, if err is not nil must be something wrong
func (e *Engine) Query() (rows int64, err error) {
	// TODO @libin Query() implement
	return
}

// orm insert
// return last insert id and error, if err is not nil must be something wrong
func (e *Engine) Insert() (lastInsertId int64, err error) {
	// TODO @libin Insert() implement
	return
}

// orm insert or update if exist
// strColumns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return last insert id and error, if err is not nil must be something wrong, if your primary key is not a int/int64 type, maybe lastInsertId return 0
func (e *Engine) Upsert(strColumns ...string) (lastInsertId int64, err error) {
	// TODO @libin Upsert() implement
	return
}

// orm update from model
// strColumns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return rows affected and error, if err is not nil must be something wrong
func (e *Engine) Update(strColumns ...string) (rows int64, err error) {
	// TODO @libin Update() implement
	return
}

// use raw sql to query results
// return rows and error, if err is not nil must be something wrong
func (e *Engine) QuerySQL(strQuery string, args ...interface{}) (rows int64, err error) {
	assert(strQuery, "query sql string is nil")
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
