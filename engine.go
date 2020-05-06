package sqlca

import (
	"database/sql"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/redigogo"
	_ "github.com/denisenkom/go-mssqldb" //mssql golang driver
	_ "github.com/go-sql-driver/mysql"   //mysql golang driver
	"github.com/jmoiron/sqlx"            //sqlx package
	_ "github.com/lib/pq"                //postgres golang driver
	_ "github.com/mattn/go-sqlite3"      //sqlite3 golang driver
	"strings"
)

type Engine struct {
	db              *sqlx.DB               // sqlx instance
	tx              *sql.Tx                // sql tx instance
	cache           redigogo.Cache         // redis cache instance
	isCacheBefore   bool                   // is cache update before db or not (default false)
	adapterSqlx     AdapterType            // what's adapter of sqlx
	adapterCache    AdapterType            // what's adapter of cache
	modelType       ModelType              // model type
	operType        OperType               // operation type
	expireTime      int                    // cache expire time of seconds
	bUseCache       bool                   // can update to cache or read from cache? (true=yes false=no)
	bCacheFirst     bool                   // cache first or database first (true=cache first; false=db first)
	model           interface{}            // data model [struct object or struct slice]
	dict            map[string]interface{} // data model db dictionary
	strDatabaseName string                 // database name
	strTableName    string                 // table name
	strPkName       string                 // primary key of table, default 'id'
	strPkValue      string                 // primary key's value
	strWhere        string                 // where condition to query or update
	strLimit        string                 // limit
	strOffset       string                 // offset (only for postgres)
	strAscOrDesc    string                 // order by ... [asc|desc]
	strDistinct     string                 // distinct
	selectColumns   []string               // columns to query: select
	conflictColumns []string               // conflict key on duplicate set (just for postgresql)
	orderByColumns  []string               // order by columns
	groupByColumns  []string               // group by columns
	inConditions    []condition            // in condition
	notConditions   []condition            // not in condition
	andConditions   []string               // and condition
	cacheIndexes    []tableIndex           // index read or write cache
	dbTags          []string               // custom db tag names
	readOnly        []string               // read only column names
}

func init() {
	log.SetLevel(log.LEVEL_INFO)
}

func NewEngine(args ...interface{}) *Engine {

	e := &Engine{
		strPkName:  DEFAULT_PRIMARY_KEY_NAME,
		expireTime: DEFAULT_CAHCE_EXPIRE_SECONDS,
	}
	e.dbTags = append(e.dbTags, TAG_NAME_DB, TAG_NAME_SQLCA)
	return e
}

// get data base driver name and data source name
func (e *Engine) getDriverNameAndDSN(adapterType AdapterType, strUrl string) (strDriverName, strDSN string) {

	strDriverName = adapterType.DriverName()
	switch adapterType {
	case AdapterSqlx_MySQL:
		return strDriverName, e.parseMysqlUrl(strUrl)
	case AdapterSqlx_Postgres:
		return strDriverName, e.parsePostgresUrl(strUrl)
	case AdapterSqlx_Sqlite:
		return strDriverName, e.parseSqliteUrl(strUrl)
	case AdapterSqlx_Mssql:
		return strDriverName, e.parseMssqlUrl(strUrl)
	case AdapterCache_Redis:
		return strDriverName, e.parseRedisUrl(strUrl)
	}
	return strDriverName, strUrl
}

// open a sqlx database or cache connection
// strUrl:
//
//  1. data source name
//
// 	   [mysql]    Open("mysql://root:123456@127.0.0.1:3306/mydb?charset=utf8mb4")
// 	   [postgres] Open("postgres://root:123456@127.0.0.1:5432/mydb?sslmode=disable")
// 	   [sqlite]   Open("sqlite:///var/lib/my.db")
// 	   [mssql]    Open("mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false")
//
//  2. cache config
//     [redis-alone]    Open("redis://123456@127.0.0.1:6379/cluster?db=0")
//     [redis-cluster]  Open("redis://123456@127.0.0.1:6379/cluster?db=0&replicate=127.0.0.1:6380,127.0.0.1:6381")
//
// expireSeconds cache data expire seconds, just for redis
func (e *Engine) Open(strUrl string, expireSeconds ...int) *Engine {

	var err error
	var adapterType AdapterType
	var strDriverName, strDSN string
	us := strings.Split(strUrl, URL_SCHEME_SEP)
	if len(us) != 2 { //default mysql
		adapterType = AdapterSqlx_MySQL
		strDriverName, strDSN = e.parseMysqlDSN(adapterType, strUrl)
	} else {
		adapterType = getAdapterType(us[0])
		strDriverName, strDSN = e.getDriverNameAndDSN(adapterType, strUrl)
	}

	switch adapterType {
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres, AdapterSqlx_Sqlite, AdapterSqlx_Mssql:
		if e.db, err = sqlx.Open(strDriverName, strDSN); err != nil {
			assert(false, "open url [%v] driver name [%v] DSN [%v] error [%v]", strUrl, strDriverName, strDSN, err.Error())
		}
		if err = e.db.Ping(); err != nil {
			assert(false, "ping url [%v] driver name [%v] DSN [%v] error [%v]", strUrl, strDriverName, strDSN, err.Error())
		}
		e.adapterSqlx = adapterType
	case AdapterCache_Redis:
		var err error
		if e.cache, err = newCache(strDriverName, strDSN); err != nil {
			assert(false, "new cache by driver name [%v] DSN [%v] error [%v]", strDriverName, strDSN, err.Error())
		}
		e.adapterCache = adapterType
		if len(expireSeconds) > 0 {
			e.expireTime = expireSeconds[0]
		} else {
			e.expireTime = 3600 //one hour expire
		}
	default:
		assert(false, "adapter instance type [%v] url [%s] not support", adapterType, strUrl)
	}

	//log.Struct(e)
	return e
}

// attach from a exist sqlx db instance
func (e *Engine) Attach(strDatabaseName string, db *sqlx.DB) *Engine {
	e.db = db
	e.setDatabaseName(strDatabaseName)
	return e
}

// get sqlx db instance
func (e *Engine) DB() *sqlx.DB {
	return e.db
}

// set log file
func (e *Engine) SetLogFile(strPath string) {
	log.Open(strPath)
}

// set cache indexes. if null, the primary key (eg. 'id') will be cached to redis
func (e *Engine) Cache(indexes ...string) *Engine {
	e.setUseCache(true)
	for _, v := range indexes {

		if itf := e.getModelValue(v); itf != nil {
			e.setIndexes(v, itf)
		} else {
			log.Warnf("index key=%v value=%v", v, itf)
		}
	}
	return e
}

// debug mode on or off
// if debug on, some method will panic if your condition illegal
func (e *Engine) Debug(ok bool) {
	e.setDebug(ok)
}

// orm model
// use to get result set, support single struct object or slice [pointer type]
// notice: will clone a new engine object for orm operations(query/update/insert/upsert)
func (e *Engine) Model(args ...interface{}) *Engine {
	assert(args, "model is nil")
	assert(e.db, "sqlx instance is nil, please call Open method first")

	return e.clone(args...)
}

// set orm query table name(s)
// when your struct type name is not a table name
func (e *Engine) Table(strNames ...string) *Engine {
	assert(strNames, "table name is nil")
	e.setTableName(strNames...)
	return e
}

// set orm primary key's name, default named 'id'
func (e *Engine) SetPkName(strName string) *Engine {
	assert(strName, "name is nil")
	e.strPkName = strName
	return e
}

func (e *Engine) GetPkName() string {
	return e.strPkName
}

// set orm primary key's value
func (e *Engine) Id(value interface{}) *Engine {

	//TODO @libin sql syntax differences of MySQL/Postgresql/Sqlite/Mssql...
	e.setSelectColumns("*")
	e.setPkValue(value)
	return e
}

// orm select/update columns
func (e *Engine) Select(strColumns ...string) *Engine {
	e.setSelectColumns(strColumns...)
	return e
}

// set distinct when select
func (e *Engine) Distinct() *Engine {
	e.setDistinct()
	return e
}

// orm where condition
func (e *Engine) Where(strWhere string, args ...interface{}) *Engine {
	assert(strWhere, "string is nil")
	strWhere = e.formatString(strWhere, args...)
	e.setWhere(strWhere)
	return e
}

func (e *Engine) And(strFmt string, args ...interface{}) *Engine {
	e.andConditions = append(e.andConditions, e.formatString(strFmt, args...))
	return e
}

// set the conflict columns for upsert
// only for postgresql
func (e *Engine) OnConflict(strColumns ...string) *Engine {

	e.setConflictColumns(strColumns...)
	return e
}

// query limit
// Limit(10) - query records limit 10 (mysql/postgres)
func (e *Engine) Limit(args ...int) *Engine {

	//TODO postgresql/mssql limit statement
	nArgs := len(args)
	if nArgs == 0 {
		return e
	}

	switch e.adapterSqlx {
	case AdapterSqlx_Mssql:
		{
			e.setLimit(fmt.Sprintf("TOP %v", args[0]))
		}
	default:
		{
			if nArgs == 1 {
				e.setLimit(fmt.Sprintf("LIMIT %v", args[0]))
			} else if nArgs == 2 {
				e.setLimit(fmt.Sprintf("LIMIT %v,%v", args[0], args[1]))
			}
		}
	}

	return e
}

// query offset (for mysql/postgres)
func (e *Engine) Offset(offset int) *Engine {
	e.setOffset(fmt.Sprintf("OFFSET %v", offset))
	return e
}

// order by [field1,field2...]
func (e *Engine) OrderBy(strColumns ...string) *Engine {
	e.setOrderBy(strColumns...)
	return e
}

// order by [field1,field2...] asc
func (e *Engine) Asc() *Engine {
	e.setAscOrDesc(ORDER_BY_ASC)
	return e
}

// order by [field1,field2...] desc
func (e *Engine) Desc() *Engine {
	e.setAscOrDesc(ORDER_BY_DESC)
	return e
}

// `field_name` IN ('1','2',...)
func (e *Engine) In(strColumn string, args ...interface{}) *Engine {

	v := condition{
		ColumnName:   strColumn,
		ColumnValues: args,
	}
	e.inConditions = append(e.inConditions, v)
	return e
}

// `field_name` NOT IN ('1','2',...)
func (e *Engine) Not(strColumn string, args ...interface{}) *Engine {

	v := condition{
		ColumnName:   strColumn,
		ColumnValues: args,
	}
	e.notConditions = append(e.notConditions, v)
	return e
}

// group by [field1,field2...]
func (e *Engine) GroupBy(strColumns ...string) *Engine {
	e.setGroupBy(strColumns...)
	return e
}

// orm query
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Query() (rowsAffected int64, err error) {
	assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()

	e.setOperType(OperType_Query)
	if e.getUseCache() {

		var ok bool
		if rowsAffected, ok = e.queryCache(); ok {
			log.Debugf("query from cache ok, rows affected [%v]", rowsAffected)
			return
		}
	}

	strSql := e.makeSqlxString()

	var rows *sql.Rows
	if rows, err = e.db.Query(strSql); err != nil {
		log.Errorf("query [%v] error [%v]", strSql, err.Error())
		return
	}

	defer rows.Close()

	return e.fetchRows(rows)
}

// orm find with customer conditions (map[string]interface{})
func (e *Engine) Find(conditions map[string]interface{}) (rowsAffected int64, err error) {
	assert(len(conditions), "find condition is nil")
	e.setOperType(OperType_Query)
	for k, v := range conditions {
		e.And("%v%v%v=%v%v%v", e.getForwardQuote(), k, e.getBackQuote(), e.getSingleQuote(), v, e.getSingleQuote())
	}
	return e.Query()
}

// orm insert
// return last insert id and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Insert() (lastInsertId int64, err error) {
	assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()

	e.setOperType(OperType_Insert)
	var strSql string
	strSql = e.makeSqlxString()

	switch e.adapterSqlx {
	case AdapterSqlx_Mssql:
		{
			if e.isPkInteger() && e.isPkValueNil() {
				lastInsertId, err = e.mssqlQueryInsert(strSql)
			}
		}
	case AdapterSqlx_Postgres:
		{
			if e.isPkInteger() && e.isPkValueNil() {
				lastInsertId, err = e.postgresQueryInsert(strSql)
			}
		}
	default:
		{
			var r sql.Result
			r, err = e.db.Exec(strSql)
			if err != nil {
				log.Errorf("error %v model %+v", err, e.model)
				return
			}

			lastInsertId, _ = r.LastInsertId() //MSSQL Server not support last insert id
		}
	}

	if lastInsertId > 0 {
		e.upsertCache(lastInsertId)
	}
	return
}

// orm insert or update if key(s) conflict
// return last insert id and error, if err is not nil must be something wrong, if your primary key is not a int/int64 type, maybe id return 0
// NOTE: Model function is must be called before call this function and call OnConflict function when you are on postgresql
func (e *Engine) Upsert() (lastInsertId int64, err error) {

	assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	assert(e.getSelectColumns(), "update columns is not set")
	defer e.cleanWhereCondition()

	e.setOperType(OperType_Upsert)
	var strSql string
	strSql = e.makeSqlxString()

	switch e.adapterSqlx {
	case AdapterSqlx_Mssql:
		{
			lastInsertId, err = e.mssqlUpsert(e.makeSqlxInsert())
		}
	case AdapterSqlx_Postgres:
		{
			lastInsertId, err = e.postgresQueryUpsert(strSql)
		}
	default:
		{
			var r sql.Result
			r, err = e.db.Exec(strSql)
			if err != nil {
				log.Errorf("error %v model %+v", err, e.model)
				return
			}
			lastInsertId, err = r.LastInsertId()
			if err != nil {
				log.Errorf("get last insert id error %v model %+v", err, e.model)
				return
			}
			if lastInsertId > 0 {
				e.upsertCache(lastInsertId)
			}
		}
	}

	return
}

// orm update from model
// strColumns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Update() (rowsAffected int64, err error) {
	assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	assert(e.getSelectColumns(), "update columns is not set, please call Select method")

	e.setOperType(OperType_Update)
	defer e.cleanWhereCondition()

	if e.getCacheBefore() {
		e.updateCache() //update data to cache before database updated
	}
	var strSql string
	strSql = e.makeSqlxString()

	var r sql.Result
	r, err = e.db.Exec(strSql)
	if err != nil {
		log.Errorf("error %v model %+v", err, e.model)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		log.Errorf("get rows affected error [%v] query [%v] model [%+v]", err, strSql, e.model)
		return
	}
	log.Debugf("RowsAffected [%v] query [%v]", rowsAffected, strSql)

	if rowsAffected > 0 && !e.getCacheBefore() {
		e.updateCache() //update data to cache after database updated
	}
	return
}

// orm delete record(s) from db and cache
func (e *Engine) Delete() (rowsAffected int64, err error) {
	e.setOperType(OperType_Delete)
	strSql := e.makeSqlxString()
	defer e.cleanWhereCondition()

	var r sql.Result
	r, err = e.db.Exec(strSql)
	if err != nil {
		log.Errorf("error %v model %+v", err, e.model)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		log.Errorf("get rows affected error [%v] query [%v] model [%+v]", err, strSql, e.model)
		return
	}
	log.Debugf("RowsAffected [%v] query [%v]", rowsAffected, strSql)

	if rowsAffected > 0 {
		e.deleteCache() //delete from cache
	}
	return
}

// use raw sql to query results
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) QueryRaw(strQuery string, args ...interface{}) (rowsAffected int64, err error) {
	assert(e.db, "sqlx db instance is nil")
	assert(strQuery, "query sql string is nil")
	assert(e.model, "model is nil, please call Model method first")

	e.setOperType(OperType_QueryRaw)

	var rows *sqlx.Rows
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	if rows, err = e.db.Queryx(strQuery); err != nil {
		log.Errorf("query [%v] error [%v]", strQuery, err.Error())
		return
	}

	defer rows.Close()
	return e.fetchRows(rows.Rows)
}

// use raw sql to query results into a map slice (model type is []map[string]string)
// return results and error
// NOTE: Model function is must be called before call this function
func (e *Engine) QueryMap(strQuery string, args ...interface{}) (rowsAffected int64, err error) {
	assert(strQuery, "query sql string is nil")
	assert(e.model, "model is nil, please call Model method first")

	e.setOperType(OperType_QueryMap)
	var rows *sqlx.Rows

	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	if rows, err = e.db.Queryx(strQuery); err != nil {
		log.Errorf("SQL [%v] query error [%v]", strQuery, err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		rowsAffected++
		fetcher, _ := e.getFecther(rows.Rows)
		*e.model.(*[]map[string]string) = append(*e.model.(*[]map[string]string), fetcher.mapValues)
	}
	return
}

// use raw sql to insert/update database, results can not be cached to redis/memcached/memory...
// return rows affected and error, if err is not nil must be something wrong
func (e *Engine) ExecRaw(strQuery string, args ...interface{}) (rowsAffected, lastInsertId int64, err error) {
	assert(e.db, "sqlx db instance is nil")
	assert(strQuery, "query sql string is nil")

	e.setOperType(OperType_ExecRaw)

	var r sql.Result
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	if r, err = e.db.Exec(strQuery); err != nil {
		log.Errorf("error [%v] model [%+v]", err, e.model)
		return
	}

	rowsAffected, err = r.RowsAffected()
	if err != nil {
		log.Errorf("get rows affected error [%v] query [%v]", err.Error(), strQuery)
		return
	}
	lastInsertId, _ = r.LastInsertId() //MSSQL Server not support last insert id
	return
}

func (e *Engine) TxBegin() (*Engine, error) {
	return e.newTx()
}

func (e *Engine) TxGet(dest interface{}, strQuery string, args ...interface{}) (count int64, err error) {
	assert(e.tx, "TxGet tx instance is nil, please call TxBegin to create a tx instance")
	var rows *sql.Rows

	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	rows, err = e.tx.Query(strQuery)
	if err != nil {
		log.Errorf("TxGet sql [%v] args %v query error [%v]", strQuery, args, err.Error())
		return
	}
	e.setModel(dest)
	defer rows.Close()
	if count, err = e.fetchRows(rows); err != nil {
		log.Errorf("TxGet sql [%v] args %v fetch row error [%v]", strQuery, args, err.Error())
		return
	}
	return
}

func (e *Engine) TxExec(strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	assert(e.tx, "TxExec tx instance is nil, please call TxBegin to create a tx instance")
	var result sql.Result

	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	result, err = e.tx.Exec(strQuery)

	if err != nil {
		log.Errorf("TxExec exec query [%v] args %+v error [%+v]", strQuery, args, err.Error())
		return
	}
	lastInsertId, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()
	return
}

func (e *Engine) TxRollback() error {
	assert(e.tx, "TxRollback tx instance is nil, please call TxBegin to create a tx instance")
	return e.tx.Rollback()
}

func (e *Engine) TxCommit() error {
	assert(e.tx, "TxCommit tx instance is nil, please call TxBegin to create a tx instance")
	return e.tx.Commit()
}

// set your customer tag for db query/insert/update (eg. go structure generated by protobuf not contain 'db' tag)
// this function must calls before Model()
func (e *Engine) SetCustomTag(tagNames ...string) *Engine {
	if len(tagNames) > 0 {
		e.dbTags = append(e.dbTags, tagNames...)
	}
	return e
}

// set cache update before database
func (e *Engine) SetCacheBefore(ok bool) {
	e.isCacheBefore = ok
}

// get cache update before database
func (e *Engine) getCacheBefore() bool {
	return e.isCacheBefore
}
