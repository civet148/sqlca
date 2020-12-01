package sqlca

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/redigogo"
	_ "github.com/denisenkom/go-mssqldb" //mssql golang driver
	"github.com/gansidui/geohash"
	_ "github.com/go-sql-driver/mysql" //mysql golang driver
	"github.com/jmoiron/sqlx"          //sqlx package
	_ "github.com/lib/pq"              //postgres golang driver
	_ "github.com/mattn/go-sqlite3"    //sqlite3 golang driver
	"strconv"
	"strings"
)

type Options struct {
	Debug bool //enable debug mode
	Max   int  //max active connections
	Idle  int  //max idle connections
	Slave bool //is a slave DSN ?
	SSH   *SSH //SSH tunnel server config
}

type TxHandler interface {
	OnTransaction(tx *Engine) error
}

type nearby struct {
	strLngCol string  //lng name of table column
	strLatCol string  //lat name of table column
	strAS     string  //select xxx as yyy (alias)
	lng       float64 //marked lng
	lat       float64 //marked lat
	distance  float64 //distance of km
}

type Engine struct {
	dsn             dsnDriver              // driver name and parameters
	slave           bool                   // use slave to query ?
	dbMasters       []*sqlx.DB             // DB instance masters
	dbSlaves        []*sqlx.DB             // DB instance slaves
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
	bForce          bool                   // force update/insert read only column(s)
	bAutoRollback   bool                   // auto rollback when tx error occurred
	model           interface{}            // data model [struct object or struct slice]
	dict            map[string]interface{} // data model db dictionary
	strDatabaseName string                 // database name
	strTableName    string                 // table name
	strPkName       string                 // primary key of table, default 'id'
	strPkValue      string                 // primary key's value
	strWhere        string                 // where condition to query or update
	strLimit        string                 // limit
	strOffset       string                 // offset (only for postgres)
	strDistinct     string                 // distinct
	excludeColumns  []string               // exclude columns for query: select xxx not contain exclude some columns
	selectColumns   []string               // columns to query: select
	conflictColumns []string               // conflict key on duplicate set (just for postgresql)
	orderByColumns  []string               // order by columns
	groupByColumns  []string               // group by columns
	ascColumns      []string               // order by xxx ASC
	descColumns     []string               // order by xxx DESC
	havingCondition string                 // having condition
	inConditions    []condition            // in condition
	notConditions   []condition            // not in condition
	andConditions   []string               // and condition
	orConditions    []string               // or condition
	cacheIndexes    []tableIndex           // index read or write cache
	dbTags          []string               // custom db tag names
	readOnly        []string               // read only column names
	slowQueryTime   int                    // slow query alert time (milliseconds)
	slowQueryOn     bool                   // enable slow query alert (default off)
	strCaseWhen     string                 // case..when...then...else...end
	nearby          *nearby                // nearby
	strUpdates      []string               // customize updates when using Upsert() ON DUPLICATE KEY UPDATE
	joins           []*Join                //inner/left/right/full-outer join(s)
}

func init() {
	log.SetLevel(log.LEVEL_INFO)
}

// args[0] data source name url
// args[1] options
// if length of args is 0, must call Open method manual
func NewEngine(args ...interface{}) *Engine {

	e := &Engine{
		strPkName:     DEFAULT_PRIMARY_KEY_NAME,
		expireTime:    DEFAULT_CAHCE_EXPIRE_SECONDS,
		slowQueryTime: DEFAULT_SLOW_QUERY_ALERT_TIME,
		adapterSqlx:   AdapterSqlx_MySQL,
	}
	e.dbTags = append(e.dbTags, TAG_NAME_DB, TAG_NAME_SQLCA, TAG_NAME_PROTOBUF, TAG_NAME_JSON)

	var ok bool
	var strOpenUrl string
	var argc = len(args)
	if argc == 0 {
		return e
	} else if argc > 0 {
		if argc == 1 {
			if strOpenUrl, ok = args[0].(string); ok {
				e.Open(strOpenUrl)
			}
		} else {
			var v1 Options
			var v2 *Options
			var options *Options
			strOpenUrl = args[0].(string)
			if v1, ok = args[1].(Options); ok {
				options = &v1
			} else {
				if v2, ok = args[1].(*Options); ok {
					options = v2
				} else {
					strOpenUrl = fmt.Sprintf(args[0].(string), args[1:]...) //legacy version compatible
					return e.Open(strOpenUrl)
				}
			}
			e.Open(strOpenUrl, options)
		}
	}

	return e
}

// get data base driver name and data source name
func (e *Engine) getDriverNameAndDSN(adapterType AdapterType, strUrl string) (driver dsnDriver) {

	driver.strDriverName = adapterType.DriverName()
	switch adapterType {
	case AdapterSqlx_MySQL:
		driver.parameter = e.parseMysqlUrl(strUrl)
		return
	case AdapterSqlx_Postgres:
		driver.parameter = e.parsePostgresUrl(strUrl)
		return
	case AdapterSqlx_Sqlite:
		driver.parameter = e.parseSqliteUrl(strUrl)
		return
	case AdapterSqlx_Mssql:
		driver.parameter = e.parseMssqlUrl(strUrl)
		return
	case AdapterCache_Redis:
		driver.parameter = e.parseRedisUrl(strUrl)
		return
	}
	return
}

// open a database or cache connection pool
// strUrl:
//
//  1. data source name
//
// 	   [mysql]    Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&slave=false&max=100&idle=1")
// 	   [postgres] Open("postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&slave=false&max=100&idle=1")
// 	   [mssql]    Open("mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false&max=100&idle=1")
// 	   [sqlite]   Open("sqlite:///var/lib/test.db")
//
//  2. cache config
//     [redis-alone]    Open("redis://123456@127.0.0.1:6379/cluster?db=0")
//     [redis-cluster]  Open("redis://123456@127.0.0.1:6379/cluster?db=0&replicate=127.0.0.1:6380,127.0.0.1:6381")
//
//  3. SSH tunnel
// 	   [mysql]    Open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&slave=false&max=100&idle=1&ssh=root:123456")
// 	   [postgres] Open("postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&slave=false&max=100&idle=1")
// 	   [mssql]    Open("mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false&max=100&idle=1")
// options:
//        1. specify master or slave, MySQL/Postgres (Options)
//        2. cache data expire seconds, just for redis (Integer)
func (e *Engine) Open(strUrl string, options ...interface{}) *Engine {

	var err error
	var adapter AdapterType
	//var strDriverName, strDSN string
	us := strings.Split(strUrl, URL_SCHEME_SEP)
	if len(us) != 2 { //default mysql
		adapter = AdapterSqlx_MySQL
		e.dsn = e.parseMysqlDSN(adapter, strUrl)
	} else {
		adapter = getAdapterType(us[0])
		e.dsn = e.getDriverNameAndDSN(adapter, strUrl)
	}
	var dsn = &e.dsn
	var opt *Options
	var parameter = &dsn.parameter
	switch adapter {
	case AdapterSqlx_MySQL, AdapterSqlx_Postgres, AdapterSqlx_Sqlite, AdapterSqlx_Mssql:
		e.adapterSqlx = adapter
		var db *sqlx.DB
		if len(options) != 0 {

			if v, ok := options[0].(*Options); ok {
				opt = v
			} else {
				o := options[0].(Options)
				opt = &o
			}
			e.Debug(opt.Debug)
			if opt.SSH != nil { //SSH tunnel enable
				dsn = opt.SSH.openSSHTunnel(dsn)
			}
		}

		if db, err = sqlx.Open(dsn.strDriverName, parameter.strDSN); err != nil {
			log.Errorf("open url [%v] driver name [%v] DSN [%v] error [%v]", strUrl, dsn.strDriverName, parameter.strDSN, err.Error())
			return nil
		}
		if err = db.Ping(); err != nil {
			log.Errorf("ping url [%v] driver name [%v] DSN [%v] error [%v]", strUrl, dsn.strDriverName, parameter.strDSN, err.Error())
			panic(err.Error())
			return nil
		}

		if opt != nil {
			dsn.SetMax(opt.Max)
			dsn.SetIdle(opt.Idle)
			dsn.SetSlave(opt.Slave)
		}
		log.Debugf("dsn parameter [%+v]", dsn.parameter)
		if parameter.max != 0 {
			db.SetMaxOpenConns(parameter.max)
		}
		if parameter.idle != 0 {
			db.SetMaxIdleConns(parameter.idle)
		}

		if parameter.slave {
			e.appendSlave(db)
		} else {
			e.appendMaster(db)
		}
	case AdapterCache_Redis:
		if e.cache, err = newCache(dsn.strDriverName, parameter.strDSN); err != nil {
			log.Errorf("new cache by driver name [%v] DSN [%v] error [%v]", dsn.strDriverName, parameter.strDSN, err.Error())
		}
		e.adapterCache = adapter
		if len(options) > 0 {
			expSec := fmt.Sprintf("%v", options[0])
			e.expireTime, _ = strconv.Atoi(expSec)
		} else {
			e.expireTime = 3600 //one hour expire
		}
	default:
		log.Errorf("adapter instance type [%v] url [%s] not support", adapter, strUrl)
		return nil
	}
	log.Infof("[%s] open url [%s] with options [%+v] ok", adapter.String(), parameter.strDSN, opt)
	return e
}

// attach from a exist sqlx db instance
func (e *Engine) Attach(strDatabaseName string, db *sqlx.DB) *Engine {
	e.appendMaster(db)
	e.setDatabaseName(strDatabaseName)
	return e
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

// log debug mode on or off
func (e *Engine) Debug(ok bool) {
	e.setDebug(ok)
}

// orm model
// use to get result set, support single struct object or slice [pointer type]
// notice: will clone a new engine object for orm operations(query/update/insert/upsert)
func (e *Engine) Model(args ...interface{}) *Engine {
	//assert(args, "model is nil")
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
	e.setPkValue(value)
	return e
}

// orm select/update columns
func (e *Engine) Select(strColumns ...string) *Engine {
	e.setSelectColumns(strColumns...)
	return e
}

// orm select/update columns
func (e *Engine) Exclude(strColumns ...string) *Engine {
	e.setExcludeColumns(strColumns...)
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

func (e *Engine) Or(strFmt string, args ...interface{}) *Engine {
	e.orConditions = append(e.orConditions, e.formatString(strFmt, args...))
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

//page query
//SELECT ... FROM ... WHERE ... LIMIT (pageNo*pageSize), pageSize
func (e *Engine) Page(pageNo, pageSize int) *Engine {
	return e.Limit(pageNo*pageSize, pageSize)
}

// query offset (for mysql/postgres)
func (e *Engine) Offset(offset int) *Engine {
	e.setOffset(fmt.Sprintf("OFFSET %v", offset))
	return e
}

// having [condition]
func (e *Engine) Having(strFmt string, args ...interface{}) *Engine {
	strCondition := e.formatString(strFmt, args...)
	e.setHaving(strCondition)
	return e
}

// order by [field1,field2...] [ASC]
func (e *Engine) OrderBy(strColumns ...string) *Engine {
	e.setOrderBy(strColumns...)
	return e
}

// order by [field1,field2...] asc
func (e *Engine) Asc(strColumns ...string) *Engine {

	if len(strColumns) == 0 {
		e.setAscColumns(e.orderByColumns...) // default order by columns as asc
	} else {
		e.setAscColumns(strColumns...) //custom order by asc columns
	}
	return e
}

// order by [field1,field2...] desc
func (e *Engine) Desc(strColumns ...string) *Engine {

	if len(strColumns) == 0 {
		e.setDescColumns(e.orderByColumns...) // default order by columns as desc
	} else {
		e.setDescColumns(strColumns...) //custom order by desc columns
	}
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

// orm query from a slave db
func (e *Engine) Slave() *Engine {
	e.slave = true
	return e
}

// orm count records
// SELECT COUNT(*) FROM table WHERE ...
// count, err := e.Model(nil).Table("users").Where("delete=1").Count()
func (e *Engine) Count() (count int64, err error) {
	e.setModel(&count)
	e.setSelectColumns("COUNT(*)")
	_, err = e.Query()
	return
}

// orm query
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
// if slave == true, try query from a slave connection, if not exist query from master
func (e *Engine) Query() (rowsAffected int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
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
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("Query [%s]", strSql))

	var rows *sql.Rows

	db := e.getQueryDB()
	if rows, err = db.Query(strSql); err != nil {
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
		e.And("%v=%v", e.getQuoteColumnName(k), e.getQuoteColumnValue(v))
	}
	return e.Query()
}

// orm insert
// return last insert id and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Insert() (lastInsertId int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()

	e.setOperType(OperType_Insert)
	var strSql string
	strSql = e.makeSqlxString()
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("Insert [%s]", strSql))

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
			var db *sqlx.DB

			db = e.getMaster()
			r, err = db.Exec(strSql)
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
// updates -> customize updates condition when key(s) conflict
// [MySQL]
// INSERT INTO messages(id, message_type, unread_count) VALUES('10000', '2', '1', '3')
// ON DUPLICATE KEY UPDATE message_type=values(message_type), unread_count=unread_count+values(unread_count)
// ---------------------------------------------------------------------------------------------------------------------------------------
// e.Model(&do).Table("messages").Upsert("message_type=values(message_type)", "unread_count=unread_count+values(unread_count)")
// ---------------------------------------------------------------------------------------------------------------------------------------
func (e *Engine) Upsert(strCustomizeUpdates ...string) (lastInsertId int64, err error) {

	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	e.setCustomizeUpdates(strCustomizeUpdates...)

	defer e.cleanWhereCondition()

	e.setOperType(OperType_Upsert)
	var strSql string
	strSql = e.makeSqlxString()
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("Upsert [%s]", strSql))
	db := e.getMaster()

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
			r, err = db.Exec(strSql)
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
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	assert(e.getSelectColumns(), "update columns is not set, please call Select method")

	e.setOperType(OperType_Update)
	defer e.cleanWhereCondition()

	if e.getCacheBefore() {
		e.updateCache() //update data to cache before database updated
	}
	var r sql.Result
	var strSql string
	strSql = e.makeSqlxString()

	c := e.Counter()
	defer c.Stop(fmt.Sprintf("Update [%s]", strSql))

	db := e.getMaster()
	r, err = db.Exec(strSql)
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
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("Delete [%s]", strSql))

	var r sql.Result
	db := e.getMaster()
	r, err = db.Exec(strSql)
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

	assert(strQuery, "query sql string is nil")
	//assert(e.model, "model is nil, please call Model method first")
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("QueryRaw [%s]", strQuery))
	e.setOperType(OperType_QueryRaw)

	var rows *sqlx.Rows
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	db := e.getQueryDB()
	if rows, err = db.Queryx(strQuery); err != nil {
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
	//assert(e.model, "model is nil, please call Model method first")
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("QueryMap [%s]", strQuery))

	e.setOperType(OperType_QueryMap)
	var rows *sqlx.Rows

	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)
	db := e.getQueryDB()
	if rows, err = db.Queryx(strQuery); err != nil {
		log.Errorf("SQL [%v] query error [%v]", strQuery, err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		rowsAffected++
		fetcher, _ := e.getFetcher(rows.Rows)
		*e.model.(*[]map[string]string) = append(*e.model.(*[]map[string]string), fetcher.mapValues)
	}
	return
}

// use raw sql to insert/update database, results can not be cached to redis/memcached/memory...
// return rows affected and error, if err is not nil must be something wrong
func (e *Engine) ExecRaw(strQuery string, args ...interface{}) (rowsAffected, lastInsertId int64, err error) {

	assert(strQuery, "query sql string is nil")
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("ExecRaw [%s]", strQuery))

	e.setOperType(OperType_ExecRaw)

	var r sql.Result
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)
	db := e.getMaster()
	if r, err = db.Exec(strQuery); err != nil {
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

// force update/insert read only column(s)
func (e *Engine) Force() *Engine {
	e.bForce = true
	return e
}

func (e *Engine) AutoRollback() *Engine {
	e.bAutoRollback = true
	return e
}

func (e *Engine) TxBegin() (*Engine, error) {
	return e.newTx()
}

func (e *Engine) TxGet(dest interface{}, strQuery string, args ...interface{}) (count int64, err error) {
	assert(e.tx, "TxGet tx instance is nil, please call TxBegin to create a tx instance")
	var rows *sql.Rows
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("TxGet [%s]", strQuery))
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("TxGet query [%v]", strQuery)

	rows, err = e.tx.Query(strQuery)
	if err != nil {
		log.Errorf("TxGet sql [%v] args %v query error [%v] auto rollback [%v]", strQuery, args, err.Error(), e.bAutoRollback)
		e.autoRollback()
		return
	}
	//log.Debugf("TxGet query [%v] rows ok", strQuery)
	e.setModel(dest)

	defer rows.Close()
	if count, err = e.fetchRows(rows); err != nil {
		log.Errorf("TxGet sql [%v] args %v fetch row error [%v] auto rollback [%v]", strQuery, args, err.Error(), e.bAutoRollback)
		e.autoRollback()
		return
	}
	//log.Debugf("TxGet query [%v] ok", strQuery)
	return
}

func (e *Engine) TxExec(strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	assert(e.tx, "TxExec tx instance is nil, please call TxBegin to create a tx instance")
	var result sql.Result
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("TxExec [%s]", strQuery))
	strQuery = e.formatString(strQuery, args...)
	log.Debugf("query [%v]", strQuery)

	result, err = e.tx.Exec(strQuery)

	if err != nil {
		log.Errorf("TxExec exec query [%v] args %+v error [%+v] auto rollback [%v]", strQuery, args, err.Error(), e.bAutoRollback)
		e.autoRollback()
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

// make SQL from orm model and operation type
func (e *Engine) ToSQL(operType OperType) (strSql string) {

	switch operType {
	case OperType_Query:
		strSql = e.makeSqlxQuery()
	case OperType_Update:
		strSql = e.makeSqlxUpdate()
	case OperType_Insert:
		strSql = e.makeSqlxInsert()
	case OperType_Upsert:
		strSql = e.makeSqlxUpsert()
	case OperType_Delete:
		strSql = e.makeSqlxDelete()
	case OperType_ForUpdate:
		strSql = e.makeSqlxForUpdate()
	default:
		log.Errorf("operation illegal")
	}
	return
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

// ping database
func (e *Engine) Ping() (err error) {

	for _, m := range e.dbMasters {
		if err = m.Ping(); err != nil {
			log.Errorf("ping master database error [%v]", err.Error())
			return
		}
	}

	for _, s := range e.dbSlaves {
		if err = s.Ping(); err != nil {
			log.Errorf("ping slave database error [%v]", err.Error())
			return
		}
	}
	return
}

// set read only columns
func (e *Engine) SetReadOnly(columns ...string) {
	e.readOnly = columns
}

//execute transaction by customize handler
//auto rollback when handler return error
func (e *Engine) TxHandle(handler TxHandler) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxHandle")
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("transaction begin error [%v]", err.Error())
		return
	}
	if err = handler.OnTransaction(tx); err != nil {
		_ = tx.TxRollback()
		log.Warnf("transaction rollback by handler error [%v]", err.Error())
		return
	}
	return tx.TxCommit()
}

//execute transaction by customize function
//auto rollback when function return error
func (e *Engine) TxFunc(fn func(tx *Engine) error) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxFunc")
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("transaction begin error [%v]", err.Error())
		return
	}
	if err = fn(tx); err != nil {
		_ = tx.TxRollback()
		log.Errorf("transaction rollback by handler error [%v]", err.Error())
		return
	}
	return tx.TxCommit()
}

//execute transaction by customize function with context
//auto rollback when function return error
func (e *Engine) TxFuncContext(ctx context.Context, fn func(ctx context.Context, tx *Engine) error) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxFuncContext")
	if tx, err = e.TxBegin(); err != nil {
		log.Errorf("transaction begin error [%v]", err.Error())
		return
	}
	if err = fn(ctx, tx); err != nil {
		_ = tx.TxRollback()
		log.Warnf("transaction rollback by handler error [%v]", err.Error())
		return
	}
	return tx.TxCommit()
}

//slow query alert on or off
//on -> true/false
//ms -> milliseconds (can be 0 if on is false)
func (e *Engine) SlowQuery(on bool, ms int) {
	e.slowQueryOn = on
	if on {
		e.slowQueryTime = ms
	}
}

func (e *Engine) Case(strThen string, strWhen string, args ...interface{}) *CaseWhen {
	cw := &CaseWhen{
		e: e,
	}
	cw.whens = append(cw.whens, &when{
		strThen: strThen,
		strWhen: e.formatString(strWhen, args...),
	})
	return cw
}

/* -- select geo point as distance where distance <= n km (float64)
SELECT
    a.*,
    (
    6371 * ACOS (
    COS( RADIANS( a.lat ) ) * COS( RADIANS( 28.8039097230 ) ) * COS(
      RADIANS( 121.5619236231 ) - RADIANS( a.lng )
     ) + SIN( RADIANS( a.lat ) ) * SIN( RADIANS( 28.8039097230 ) )
    )
    ) AS distance
FROM
    t_address a
HAVING distance <= 200 -- less than or equal 200km
ORDER BY
    distance
    LIMIT 10
*/
func (e *Engine) NearBy(strLngCol, strLatCol, strAS string, lng, lat, distance float64) *Engine {
	e.nearby = &nearby{
		strLngCol: strLngCol,
		strLatCol: strLatCol,
		strAS:     strAS,
		lng:       lng,
		lat:       lat,
		distance:  distance,
	}
	return e
}

//encode geo hash string (precision 1~8)
//returns geo hash and neighbors areas
func (e *Engine) GeoHash(lng, lat float64, precision int) (strGeoHash string, strNeighbors []string) {
	strGeoHash, _ = geohash.Encode(lat, lng, precision)
	strNeighbors = geohash.GetNeighbors(lat, lng, precision)
	return
}

func (e *Engine) JsonMarshal(v interface{}) (strJson string) {
	if data, err := json.Marshal(v); err != nil {
		log.Error(err.Error())
		return
	} else {
		strJson = string(data)
	}
	return
}

func (e *Engine) JsonUnmarshal(strJson string, v interface{}) (err error) {
	err = json.Unmarshal([]byte(strJson), v)
	return
}

func (e *Engine) InnerJoin(strTableName string) *Join {
	return &Join{
		e:            e,
		jt:           JoinType_Inner,
		strTableName: strTableName,
	}
}

func (e *Engine) LeftJoin(strTableName string) *Join {
	return &Join{
		e:            e,
		jt:           JoinType_Left,
		strTableName: strTableName,
	}
}

func (e *Engine) RightJoin(strTableName string) *Join {
	return &Join{
		e:            e,
		jt:           JoinType_Right,
		strTableName: strTableName,
	}
}
