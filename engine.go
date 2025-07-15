package sqlca

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "gitee.com/opengauss/openGauss-connector-go-pq" //open gauss golang driver of gitee.com
	"github.com/bwmarrin/snowflake"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/types"
	_ "github.com/denisenkom/go-mssqldb" //mssql golang driver
	"github.com/gansidui/geohash"
	_ "github.com/go-sql-driver/mysql" //mysql golang driver
	"github.com/jmoiron/sqlx"          //sqlx package
	_ "github.com/lib/pq"              //postgres golang driver
	_ "github.com/mattn/go-sqlite3"    //sqlite3 golang driver
	//_ "github.com/opengauss-mirror/openGauss-connector-go-pq" //open gauss golang driver of github.com
	"strings"
)

const (
	DefaultConnMax  = 150
	DefaultConnIdle = 5
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type ID = snowflake.ID
type Id = ID

type Options struct {
	Debug         bool       //enable debug mode
	Max           int        //max active connections
	Idle          int        //max idle connections
	SSH           *SSH       //ssh tunnel server config
	SnowFlake     *SnowFlake //snowflake id config
	DisableOffset bool       //disable page offset for LIMIT (default page no is 1, if true then page no start from 0)
	DefaultLimit  int32      //limit default (0 means no limit)
}

type SnowFlake struct {
	NodeId int64 //node id (0~1023)
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
	strDSN           string                 // database source name
	dsn              dsnDriver              // driver name and parameters
	options          Options                // database options
	db               *sqlx.DB               // DB instance masters
	tx               *sqlx.Tx               // sql tx instance
	adapterType      types.AdapterType      // what's adapter
	adapterCache     types.AdapterType      // what's adapter of cache
	modelType        types.ModelType        // model type
	operType         types.OperType         // operation type
	expireTime       int                    // cache expire time of seconds
	bForce           bool                   // force update/insert read only column(s)
	bAutoRollback    bool                   // auto rollback when tx error occurred
	model            interface{}            // data model [struct object or struct slice]
	dict             map[string]interface{} // data model db dictionary
	strDatabaseName  string                 // database name
	strTableName     string                 // table name
	strPkName        string                 // primary key of table, default 'id'
	strPkValue       string                 // primary key's value
	strLimit         string                 // limit
	strOffset        string                 // offset (only for postgres)
	strDistinct      string                 // distinct
	strForUpdate     string                 // tx query for update
	strLockShareMode string                 // tx query lock share mode
	nullableColumns  []string               // nullable columns for update/insert
	excludeColumns   []string               // exclude columns for query: select xxx not contain exclude some columns
	selectColumns    []string               // columns to query: select
	conflictColumns  []string               // conflict key on duplicate set (just for postgresql)
	orderByColumns   []string               // order by columns
	groupByColumns   []string               // group by columns
	havingCondition  string                 // having condition
	inConditions     []types.Expr           // in condition
	notConditions    []types.Expr           // not in condition
	andConditions    []types.Expr           // and condition
	orConditions     []types.Expr           // or condition
	dbTags           []string               // custom db tag names
	readOnly         []string               // read only column names
	slowQueryTime    int                    // slow query alert time (milliseconds)
	slowQueryOn      bool                   // enable slow query alert (default off)
	strCaseWhen      string                 // case..when...then...else...end
	nearby           *nearby                // nearby
	strUpdates       []string               // customize updates when using Upsert() ON DUPLICATE KEY UPDATE
	joins            []*Join                // inner/left/right/full-outer join(s)
	selected         bool                   // column(s) selected
	noVerbose        bool                   // no more verbose
	idgen            *snowflake.Node        // snowflake id generator
	hookMethods      *hookMethods           // hook methods
	insertIgnore     bool                   // insert ignore when conflict
}

func init() {
}

/*
// [mysql]    "mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4"
// [postgres] "postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public"
// [opengauss] "opengauss://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public"
// [mssql]    "mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false"
// [sqlite]   "sqlite:///var/lib/test.db"
*/
func NewEngine(strUrl string, options ...*Options) (*Engine, error) {
	e := &Engine{
		strPkName:     types.DEFAULT_PRIMARY_KEY_NAME,
		expireTime:    types.DEFAULT_CAHCE_EXPIRE_SECONDS,
		slowQueryTime: types.DEFAULT_SLOW_QUERY_ALERT_TIME,
		adapterType:   types.AdapterSqlx_MySQL,
	}
	e.dbTags = append(e.dbTags, types.TAG_NAME_DB, types.TAG_NAME_SQLCA, types.TAG_NAME_PROTOBUF, types.TAG_NAME_JSON)
	return e.open(strUrl, options...)
}

// get data base driver name and data source name
func (e *Engine) getDriverNameAndDSN(adapterType types.AdapterType, strUrl string) (driver dsnDriver) {

	driver.strDriverName = adapterType.DriverName()
	switch adapterType {
	case types.AdapterSqlx_MySQL:
		driver.parameter = e.parseMysqlUrl(strUrl)
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		driver.parameter = e.parsePostgresUrl(strUrl)
	case types.AdapterSqlx_Sqlite:
		driver.parameter = e.parseSqliteUrl(strUrl)
	case types.AdapterSqlx_Mssql:
		driver.parameter = e.parseMssqlUrl(strUrl)
	default:
		panic(fmt.Sprintf("unknown adapter [%s]", adapterType))
	}
	return
}

// open a database or cache connection pool
// strUrl:
//
//  1. data source name
//
//     [mysql]    open("mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4")
//     [postgres] open("postgres://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public")
//     [opengauss] open("opengauss://root:123456@127.0.0.1:5432/test?sslmode=disable&search_path=public")
//     [mssql]    open("mssql://sa:123456@127.0.0.1:1433/mydb?instance=SQLExpress&windows=false")
//     [sqlite]   open("sqlite:///var/lib/test.db")
//
// options:
//  1. specify master or slave, MySQL/Postgres (optional)
func (e *Engine) open(strUrl string, options ...*Options) (*Engine, error) {

	var err error
	var adapter types.AdapterType
	if len(options) == 0 {
		options = append(options, &Options{
			Debug:        false,
			Max:          DefaultConnMax,
			Idle:         DefaultConnIdle,
			DefaultLimit: 0,
		})
	}
	//var strDriverName, strDSN string
	us := strings.Split(strUrl, urlSchemeSep)
	if len(us) != 2 { //default mysql
		adapter = types.AdapterSqlx_MySQL
		e.dsn = e.parseMysqlDSN(adapter, strUrl)
	} else {
		adapter = types.GetAdapterType(us[0])
		e.dsn = e.getDriverNameAndDSN(adapter, strUrl)
	}
	var dsn = &e.dsn
	var opt *Options
	var param = &dsn.parameter

	e.adapterType = adapter
	var db *sqlx.DB
	if len(options) != 0 {
		opt = options[0]
		if opt.Debug {
			e.Debug(true)
		}
		if opt.SSH != nil { //SSH tunnel enable
			dsn = opt.SSH.openSSHTunnel(dsn)
		}
	}

	if db, err = sqlx.Open(dsn.strDriverName, param.strDSN); err != nil {
		//log.Errorf("open database driver name [%v] DSN [%v] error [%v]", dsn.strDriverName, parameter.strDSN, err.Error())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		//log.Errorf("ping database driver name [%v] DSN [%v] error [%v]", dsn.strDriverName, parameter.strDSN, err.Error())
		return nil, err
	}

	if opt != nil {
		dsn.SetMax(opt.Max)
		dsn.SetIdle(opt.Idle)
	}

	if param.max != 0 {
		db.SetMaxOpenConns(param.max)
	}
	if param.idle != 0 {
		db.SetMaxIdleConns(param.idle)
	}
	e.setDB(db)

	if opt != nil && opt.SnowFlake != nil {
		e.idgen, err = snowflake.NewNode(opt.SnowFlake.NodeId)
		if err != nil {
			return nil, log.Errorf("new snowflake id generator error [%s]", err.Error())
		}
	}
	e.strDSN = strUrl
	if opt != nil {
		e.options = *opt
	}
	return e, nil
}

// Use switch database (returns a new instance)
func (e *Engine) Use(strDatabaseName string) (*Engine, error) {
	var strUrl = e.strDSN
	us := strings.Split(strUrl, urlSchemeSep)
	if len(us) != 2 { //mysql raw database source name
		strUrl = rawMySql2Url(strUrl)
	}
	ui := ParseUrl(strUrl)
	if ui == nil {
		return nil, log.Errorf("url %s invalid", strUrl)
	}
	ui.Path = fmt.Sprintf("/%s", strDatabaseName)
	return NewEngine(ui.Url(), &e.options)
}

// Attach attach from a exist sqlx db instance
func (e *Engine) Attach(strDatabaseName string, db *sqlx.DB) *Engine {
	e.setDB(db)
	e.setDatabaseName(strDatabaseName)
	return e
}

// SetLogFile set log file
func (e *Engine) SetLogFile(strPath string) {
	log.Open(strPath)
}

// Debug log debug mode on or off
func (e *Engine) Debug(ok bool) {
	e.setDebug(ok)
}

// Model orm model
// use to get result set, support single struct object or slice [pointer type]
// notice: will clone a new engine object for orm operations(query/update/insert/upsert)
func (e *Engine) Model(args ...interface{}) *Engine {
	//assert(args, "model is nil")
	return e.clone(args...)
}

// Table set orm query table name(s) expression
// when your struct type name is not a table name
func (e *Engine) Table(exprs ...string) *Engine {
	return e.From(exprs...)
}

// From alias of Table method
func (e *Engine) From(exprs ...string) *Engine {
	assert(exprs, "from express is nil")
	e.setTableName(exprs...)
	return e
}

// SetPkName set orm primary key's name, default named 'id'
func (e *Engine) SetPkName(strName string) *Engine {
	assert(strName, "name is nil")
	e.strPkName = strName
	return e
}

func (e *Engine) GetPkName() string {
	return e.strPkName
}

// Id set orm primary key's value
func (e *Engine) Id(value any) *Engine {
	e.setPkValue(value)
	return e
}

// Select orm select/update columns
func (e *Engine) Select(columns ...string) *Engine {
	if e.setSelectColumns(columns...) {
		e.selected = true
	}
	return e
}

// Exclude exclude orm select/update columns
func (e *Engine) Exclude(columns ...string) *Engine {
	e.setExcludeColumns(columns...)
	return e
}

// Omit same as Exclude
func (e *Engine) Omit(columns ...string) *Engine {
	e.setExcludeColumns(columns...)
	return e
}

// Distinct set distinct when select
func (e *Engine) Distinct() *Engine {
	e.setDistinct()
	return e
}

// Where orm where condition
func (e *Engine) Where(query any, args ...any) *Engine {
	return e.setNormalCondition(query, args...)
}

func (e *Engine) And(query any, args ...any) *Engine {
	return e.setNormalCondition(query, args...)
}

func (e *Engine) Or(query any, args ...any) *Engine {
	var strSql string
	qt := parseQueryInterface(query)
	switch qt {
	case queryInterface_String:
		strSql = query.(string)
		assert(strSql, "query statement is empty")
		expr := e.buildSqlExpr(strSql, args...)
		e.setOr(expr)
	case queryInterface_Map:
		e.parseQueryOrMap(query)
	}
	return e
}

// OnConflict set the conflict columns for upsert
// only for postgresql
func (e *Engine) OnConflict(columns ...string) *Engine {
	e.setConflictColumns(columns...)
	return e
}

// Limit query limit
// Limit(10) - query records limit 10 (mysql/postgres)
func (e *Engine) Limit(args ...int) *Engine {

	//TODO postgresql/mssql limit statement
	nArgs := len(args)
	if nArgs == 0 {
		return e
	}

	switch e.adapterType {
	case types.AdapterSqlx_Mssql:
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

func (e *Engine) pageOffset() bool {
	return !e.options.DisableOffset
}

// Page page query
//
//	SELECT ... FROM ... WHERE ... LIMIT (pageNo*pageSize), pageSize
func (e *Engine) Page(pageNo, pageSize int) *Engine {
	if pageNo < 0 || pageSize <= 0 {
		return e
	}
	if e.pageOffset() && pageNo > 0 {
		pageNo -= 1
	}
	return e.Limit(pageNo*pageSize, pageSize)
}

// Offset query offset (for mysql/postgres)
func (e *Engine) Offset(offset int) *Engine {
	e.setOffset(fmt.Sprintf("OFFSET %v", offset))
	return e
}

// Having having [condition]
func (e *Engine) Having(strFmt string, args ...any) *Engine {
	expr := e.buildSqlExpr(strFmt, args...)
	e.setHaving(expr.RawSQL(e.GetAdapter()))
	return e
}

// OrderBy order by [field1,field2...] [ASC]
func (e *Engine) OrderBy(orders ...string) *Engine {
	e.setOrderBy(orders...)
	return e
}

// Asc order by [field1,field2...] asc
func (e *Engine) Asc(columns ...string) *Engine {

	if len(columns) == 0 {
		e.setAscColumns(e.orderByColumns...) // default order by columns as asc
	} else {
		e.setAscColumns(columns...) //custom order by asc columns
	}
	return e
}

// Desc order by [field1,field2...] desc
func (e *Engine) Desc(columns ...string) *Engine {

	if len(columns) == 0 {
		e.setDescColumns(e.orderByColumns...) // default order by columns as desc
	} else {
		e.setDescColumns(columns...) //custom order by desc columns
	}
	return e
}

// In `field_name` IN ('1','2',...)
func (e *Engine) In(strColumn string, args ...any) *Engine {
	v := e.buildSqlExpr(strColumn, args...)
	e.inConditions = append(e.inConditions, v)
	return e
}

// NotIn `field_name` NOT IN ('1','2',...)
func (e *Engine) NotIn(strColumn string, args ...any) *Engine {
	v := e.buildSqlExpr(strColumn, args...)
	e.notConditions = append(e.notConditions, v)
	return e
}

// GroupBy group by [field1,field2...]
func (e *Engine) GroupBy(columns ...string) *Engine {
	e.setGroupBy(columns...)
	return e
}

// Query orm query
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
// if slave == true, try query from a slave connection, if not exist query from master
func (e *Engine) Query() (rowsAffected int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()
	if e.options.DefaultLimit > 0 && e.strLimit == "" {
		e.setLimit(fmt.Sprintf("LIMIT %v", e.options.DefaultLimit))
	}

	strRawSql, _ := e.makeSQL(types.OperType_Query, true)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strRawSql))
	if err = e.execBeforeQueryHooks(); err != nil {
		return 0, err
	}
	rowsAffected, err = e.execQuery()
	if err != nil {
		return 0, log.Errorf("caller [%v] SQL [%s] error: %s", e.getCaller(2), strRawSql, err.Error())
	}
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strRawSql)
	if err = e.execAfterQueryHooks(); err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// QueryEx orm query with total count
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
// if slave == true, try query from a slave connection, if not exist query from master
func (e *Engine) QueryEx() (rowsAffected, total int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()
	if e.options.DefaultLimit > 0 && e.strLimit == "" {
		e.setLimit(fmt.Sprintf("LIMIT %v", e.options.DefaultLimit))
	}
	strSql, _ := e.makeSQL(types.OperType_Query, true)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	if err = e.execBeforeQueryHooks(); err != nil {
		return 0, 0, err
	}

	strCountSql := e.makeSqlxQueryCount(false)
	rowsAffected, total, err = e.execQueryEx(strCountSql)
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strSql)
	if err != nil {
		return 0, 0, log.Errorf("query count sql [%v] error [%v]", strCountSql, err.Error())
	}
	if err = e.execAfterQueryHooks(); err != nil {
		return 0, 0, err
	}
	return rowsAffected, total, nil
}

// MustFind orm find data records, returns error if not found
func (e *Engine) MustFind() (rowsAffected int64, err error) {
	rowsAffected, err = e.Query()
	if rowsAffected == 0 {
		return 0, ErrRecordNotFound
	}
	return rowsAffected, err
}

// Ignore insert ignore when primary key conflict
func (e *Engine) Ignore() *Engine {
	e.insertIgnore = true
	return e
}

// Insert orm insert
// return last insert id and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Insert() (lastInsertId, rowsAffected int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	defer e.cleanWhereCondition()
	if err = e.execBeforeCreateHooks(); err != nil {
		return 0, 0, log.Errorf(err.Error())
	}
	var strSql string
	strSql, _ = e.makeSQL(types.OperType_Insert, true)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	switch e.adapterType {
	case types.AdapterSqlx_Mssql:
		{
			strSql = e.mssqlQueryInsert(strSql)
		}
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		{
			strSql = e.postgresQueryInsert(strSql)
		}
	}

	if e.operType == types.OperType_Tx {
		lastInsertId, rowsAffected, err = e.TxExec(strSql)
	} else {
		lastInsertId, rowsAffected, err = e.mysqlExec(strSql)
	}
	if err != nil {
		return 0, 0, log.Errorf("SQL [%v] error: %v", strSql, err.Error())
	}
	e.verbose("caller [%v] last id [%v] SQL [%s]", e.getCaller(2), lastInsertId, strSql)
	if err = e.execAfterCreateHooks(); err != nil {
		return 0, 0, log.Errorf(err.Error())
	}
	return
}

// Upsert orm insert or update if key(s) conflict
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

	strSql, _ := e.makeSQL(types.OperType_Upsert, true)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	switch e.adapterType {
	case types.AdapterSqlx_Mssql:
		{
			if e.operType == types.OperType_Tx {
				return 0, log.Errorf("MSSQL can not use upsert on tx mode")
			}
			lastInsertId, err = e.mssqlUpsert(e.makeSqlxInsert())
		}
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		{
			if e.operType == types.OperType_Tx {
				return 0, log.Errorf("Postgres can not use upsert on tx mode")
			}
			lastInsertId, err = e.postgresQueryUpsert(strSql)
		}
	default:
		{
			lastInsertId, err = e.mysqlQueryUpsert(strSql)
		}
	}
	if err != nil {
		return 0, log.Errorf("SQL [%v] error: %v", strSql, err.Error())
	}
	e.verbose("caller [%v] last id [%v] SQL [%s]", e.getCaller(2), lastInsertId, strSql)
	return
}

// Update orm update from model
// columns... if set, columns will be updated, if none all columns in model will be updated except primary key
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) Update() (rowsAffected int64, err error) {
	//assert(e.model, "model is nil, please call Model method first")
	assert(e.strTableName, "table name not found")
	assert(e.getSelectColumns(), "update columns is not set, please call Select method")
	defer e.cleanWhereCondition()

	if err = e.execBeforeUpdateHooks(); err != nil {
		return 0, log.Errorf(err.Error())
	}

	strSql, _ := e.makeSQL(types.OperType_Update, true)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	var r sql.Result
	var execer = e.getExecer()

	query, args := e.makeSQL(types.OperType_Update, false)
	//log.Debugf("query %s args %v", query, args)
	r, err = execer.Exec(query, args...)
	if err != nil {
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		return 0, log.Errorf(err)
	}
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strSql)
	if err = e.execAfterUpdateHooks(); err != nil {
		return 0, log.Errorf(err.Error())
	}
	return rowsAffected, nil
}

// Delete orm delete record(s) from db
func (e *Engine) Delete() (rowsAffected int64, err error) {
	if err = e.execBeforeDeleteHooks(); err != nil {
		return 0, log.Errorf(err.Error())
	}

	strSql, args := e.makeSQL(types.OperType_Delete, true)
	defer e.cleanWhereCondition()
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	var execer = e.getExecer()

	var r sql.Result
	r, err = execer.Exec(strSql, args...)
	if err != nil {
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strSql)
	if err = e.execAfterDeleteHooks(); err != nil {
		return 0, log.Errorf(err.Error())
	}
	return rowsAffected, nil
}

// QueryRaw use raw sql to query results
// return rows affected and error, if err is not nil must be something wrong
// NOTE: Model function is must be called before call this function
func (e *Engine) QueryRaw(query string, args ...any) (rowsAffected int64, err error) {

	assert(query, "query sql string is nil")
	//assert(e.model, "model is nil, please call Model method first")

	var rows *sqlx.Rows
	var strSql string
	var expr = e.buildSqlExpr(query, args...)
	strSql = expr.RawSQL(e.GetAdapter())

	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	var queryer sqlx.Queryer
	if e.operType == types.OperType_Tx {
		queryer = e.tx
	} else {
		queryer = e.getDB()
	}

	if rows, err = queryer.Queryx(expr.SQL, expr.Vars...); err != nil {
		return
	}

	defer rows.Close()
	rowsAffected, err = e.fetchRows(rows.Rows)
	if err != nil {
		log.Errorf(err.Error())
		return 0, err
	}
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strSql)
	return rowsAffected, nil
}

// QueryMap use raw sql to query results into a map slice (model type is []map[string]string)
// return results and error
// NOTE: Model function is must be called before call this function
func (e *Engine) QueryMap(query string, args ...any) (rowsAffected int64, err error) {
	assert(query, "query sql string is nil")
	//assert(e.model, "model is nil, please call Model method first")
	var expr = e.buildSqlExpr(query, args...)
	strSql := expr.RawSQL(e.GetAdapter())

	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	var queryer sqlx.Queryer
	if e.operType == types.OperType_Tx {
		queryer = e.tx
	} else {
		queryer = e.getDB()
	}
	if err = e.execBeforeQueryHooks(); err != nil {
		return 0, err
	}
	var rows *sqlx.Rows
	if rows, err = queryer.Queryx(expr.SQL, expr.Vars...); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		rowsAffected++
		fetcher, _ := e.getFetcher(rows.Rows)
		*e.model.(*[]map[string]string) = append(*e.model.(*[]map[string]string), fetcher.mapValues)
	}
	e.verbose("caller [%v] rows [%v] SQL [%s]", e.getCaller(2), rowsAffected, strSql)
	if err = e.execAfterQueryHooks(); err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// ExecRaw use raw sql to insert/update database, results can not be cached to redis/memcached/memory...
// return rows affected and error, if err is not nil must be something wrong
func (e *Engine) ExecRaw(query string, args ...any) (rowsAffected, lastInsertId int64, err error) {

	assert(query, "query sql string is nil")
	var expr = e.buildSqlExpr(query, args...)
	strSql := expr.RawSQL(e.GetAdapter())
	var execer = e.getExecer()

	var r sql.Result

	log.Debugf("exec [%v]", strSql)
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("SQL [%s]", strSql))

	if r, err = execer.Exec(expr.SQL, expr.Vars...); err != nil {
		return
	}

	rowsAffected, err = r.RowsAffected()
	if err != nil {
		return
	}
	lastInsertId, _ = r.LastInsertId() //MSSQL Server not support last insert id
	e.verbose("caller [%v] rows [%v] last id [%v] SQL [%s]", e.getCaller(2), rowsAffected, lastInsertId, strSql)
	return rowsAffected, lastInsertId, nil
}

// Force force update/insert read only column(s)
func (e *Engine) Force() *Engine {
	e.bForce = true
	return e
}

// Close disconnect all database connections
func (e *Engine) Close() *Engine {
	_ = e.db.Close()
	return e
}

func (e *Engine) AutoRollback() *Engine {
	e.bAutoRollback = true
	return e
}

func (e *Engine) CountRows() (count int64, err error) {
	assert(e.strDatabaseName, "table name requires")
	strCountSql := e.makeSqlxQueryCount(true)
	var queryer = e.getQueryer()

	var rowsCount *sql.Rows
	if rowsCount, err = queryer.Query(strCountSql); err != nil {
		return 0, err
	}
	e.verbose("caller [%v] SQL [%s]", e.getCaller(2), strCountSql)
	defer rowsCount.Close()
	for rowsCount.Next() {
		_ = rowsCount.Scan(&count)
	}
	return count, nil
}

func (e *Engine) TxBegin() (*Engine, error) {
	return e.newTx()
}

func (e *Engine) TxGet(dest any, strSql string, args ...any) (count int64, err error) {
	assert(e.tx, "tx instance is nil")
	return e.txQuery(dest, strSql, args...)
}

func (e *Engine) TxExec(query string, args ...any) (lastInsertId, rowsAffected int64, err error) {
	assert(e.tx, "tx instance is nil")
	var result sql.Result
	var expr = e.buildSqlExpr(query, args...)
	strSql := expr.RawSQL(e.GetAdapter())

	c := e.Counter()
	defer c.Stop(fmt.Sprintf("tx [%s]", strSql))

	result, err = e.tx.Exec(expr.SQL, expr.Vars...)
	if err != nil {
		e.autoRollback()
		return
	}
	lastInsertId, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()
	e.verbose("caller [%v] rows [%v] last id [%v] SQL [%s]", e.getCaller(2), rowsAffected, lastInsertId, strSql)
	return
}

func (e *Engine) TxRollback() error {
	assert(e.tx, "tx instance is nil")
	return e.txRollback()
}

func (e *Engine) TxCommit() error {
	assert(e.tx, "tx instance is nil")
	return e.txCommit()
}

// make SQL from orm model and operation type
func (e *Engine) ToSQL(operType types.OperType) (strSql string) {

	switch operType {
	case types.OperType_Query:
		strSql, _ = e.makeSqlxQuery(true)
	case types.OperType_Update:
		strSql, _ = e.makeSqlxUpdate(true)
	case types.OperType_Insert:
		strSql = e.makeSqlxInsert()
	case types.OperType_Upsert:
		strSql = e.makeSqlxUpsert()
	case types.OperType_Delete:
		strSql, _ = e.makeSqlxDelete(true)
	case types.OperType_ForUpdate:
		strSql, _ = e.makeSqlxForUpdate(true)
	default:
		log.Errorf("operation illegal")
	}
	return strSql
}

// SetCustomTag set your customer tag for db query/insert/update (eg. go structure generated by protobuf not contain 'db' tag)
// this function must calls before Model()
func (e *Engine) SetCustomTag(tagNames ...string) *Engine {
	if len(tagNames) > 0 {
		e.dbTags = append(e.dbTags, tagNames...)
	}
	return e
}

// Ping ping database
func (e *Engine) Ping() (err error) {
	if err = e.db.Ping(); err != nil {
		log.Errorf("ping database error [%v]", err.Error())
		return
	}
	return nil
}

// SetReadOnly set read only columns
func (e *Engine) SetReadOnly(columns ...string) {
	e.readOnly = columns
}

// TxHandle execute transaction by customize handler
// auto rollback when handler return error
func (e *Engine) TxHandle(handler TxHandler) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxHandle")
	if tx, err = e.TxBegin(); err != nil {
		return err
	}
	if err = handler.OnTransaction(tx); err != nil {
		_ = tx.TxRollback()
		return
	}
	return tx.TxCommit()
}

// TxFunc execute transaction by customize function
//
//	auto rollback when function return error
func (e *Engine) TxFunc(fn func(tx *Engine) error) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxFunc")
	if tx, err = e.TxBegin(); err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		_ = tx.TxRollback()
		return
	}
	return tx.TxCommit()
}

// TxFuncContext execute transaction by customize function with context
//
//	auto rollback when function return error
func (e *Engine) TxFuncContext(ctx context.Context, fn func(ctx context.Context, tx *Engine) error) (err error) {
	var tx *Engine
	c := e.Counter()
	defer c.Stop("TxFuncContext")
	if tx, err = e.TxBegin(); err != nil {
		return
	}
	if err = fn(ctx, tx); err != nil {
		_ = tx.TxRollback()
		return
	}
	return tx.TxCommit()
}

// QueryJson query result marshal to json
func (e *Engine) QueryJson() (s string, err error) {
	var count int64
	count, err = e.Query()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	if count != 0 && e.model != nil {
		var data []byte
		if data, err = json.Marshal(e.model); err != nil {
			return s, err
		}
		s = string(data)
	}
	return s, nil
}

// SlowQuery slow query alert on or off
//
//	on -> true/false
//	ms -> milliseconds (can be 0 if on is false)
func (e *Engine) SlowQuery(on bool, ms int) {
	e.slowQueryOn = on
	if on {
		e.slowQueryTime = ms
	}
}

func (e *Engine) Case(strThen string, strWhen string, args ...any) *CaseWhen {
	cw := &CaseWhen{
		e: e,
	}
	cw.whens = append(cw.whens, &when{
		strThen: strThen,
		strWhen: e.buildSqlExpr(strWhen, args...).RawSQL(e.GetAdapter()),
	})
	return cw
}

/*
NearBy
-- select geo point as distance where distance <= n km (float64)
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

// GeoHash encode geo hash string (precision 1~8)
//
//	returns geo hash and neighbors areas
func (e *Engine) GeoHash(lng, lat float64, precision int) (strGeoHash string, strNeighbors []string) {
	strGeoHash, _ = geohash.Encode(lat, lng, precision)
	strNeighbors = geohash.GetNeighbors(lat, lng, precision)
	return
}

func (e *Engine) JsonMarshal(v any) (strJson string) {
	if data, err := json.Marshal(v); err != nil {
		log.Error(err.Error())
	} else {
		strJson = string(data)
	}
	return
}

func (e *Engine) JsonUnmarshal(strJson string, v any) (err error) {
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

func (e *Engine) GetAdapter() types.AdapterType {
	return e.adapterType
}

func (e *Engine) Count(strColumn string, as ...string) *Engine {
	return e.Select(e.aggFunc(types.DATABASE_KEY_NAME_COUNT, strColumn, as...))
}

func (e *Engine) Sum(strColumn string, as ...string) *Engine {
	return e.Select(e.aggFunc(types.DATABASE_KEY_NAME_SUM, strColumn, as...))
}

func (e *Engine) Avg(strColumn string, as ...string) *Engine {
	return e.Select(e.aggFunc(types.DATABASE_KEY_NAME_AVG, strColumn, as...))
}

func (e *Engine) Min(strColumn string, as ...string) *Engine {
	return e.Select(e.aggFunc(types.DATABASE_KEY_NAME_MIN, strColumn, as...))
}

func (e *Engine) Max(strColumn string, as ...string) *Engine {
	return e.Select(e.aggFunc(types.DATABASE_KEY_NAME_MAX, strColumn, as...))
}

func (e *Engine) Round(strColumn string, round int, as ...string) *Engine {
	return e.Select(e.roundFunc(strColumn, round, as...))
}

func (e *Engine) NoVerbose() *Engine {
	e.noVerbose = true
	return e
}

func (e *Engine) Like(strColumn, keyword string) *Engine {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL:
		e.And(fmt.Sprintf("LOCATE('%s', %s)", keyword, strColumn))
	default:
		e.And(fmt.Sprintf("%s LIKE '%%%s%%'", strColumn, keyword))
	}
	return e
}

func (e *Engine) Likes(kvs map[string]any) *Engine {
	var likes []string
	for k, v := range kvs {
		likes = append(likes, fmt.Sprintf(" %s LIKE '%%%v%%' ", k, v))
	}
	strLikes := strings.Join(likes, types.DATABASE_KEY_NAME_OR)
	strLikes = "(" + strLikes + ")"
	return e.And(strLikes)
}

func (e *Engine) IsNull(strColumn string) *Engine {
	e.And(strColumn + " IS NULL")
	return e
}

func (e *Engine) Equal(strColumn string, value any) *Engine {
	e.And(strColumn, indirectValue(value))
	return e
}

// Eq alias of Equal
func (e *Engine) Eq(strColumn string, value any) *Engine {
	return e.Equal(strColumn, indirectValue(value))
}

// Ne not equal
func (e *Engine) Ne(strColumn string, value any) *Engine {
	return e.And(strColumn+" != ?", indirectValue(value))
}

func (e *Engine) GreaterThan(strColumn string, value any) *Engine {
	e.And(strColumn+" > ?", indirectValue(value))
	return e
}

// Gt alias of GreaterThan
func (e *Engine) Gt(strColumn string, value any) *Engine {
	return e.GreaterThan(strColumn, indirectValue(value))
}

func (e *Engine) GreaterEqual(strColumn string, value any) *Engine {
	e.And(strColumn+" >= ?", indirectValue(value))
	return e
}

// Gte alias of GreaterEqual
func (e *Engine) Gte(strColumn string, value any) *Engine {
	return e.GreaterEqual(strColumn, indirectValue(value))
}

func (e *Engine) LessThan(strColumn string, value any) *Engine {
	e.And(strColumn+" < ?", indirectValue(value))
	return e
}

// Lt alias of LessThan
func (e *Engine) Lt(strColumn string, value any) *Engine {
	return e.LessThan(strColumn, indirectValue(value))
}

func (e *Engine) LessEqual(strColumn string, value any) *Engine {
	e.And(strColumn+" <= ?", indirectValue(value))
	return e
}

// Lte alias of LessEqual
func (e *Engine) Lte(strColumn string, value any) *Engine {
	return e.LessEqual(strColumn, indirectValue(value))
}

// GteLte greater than equal and less than equal
func (e *Engine) GteLte(strColumn string, value1, value2 any) *Engine {
	e.Gte(strColumn, value1)
	e.Lte(strColumn, value2)
	return e
}

func (e *Engine) IsNULL(strColumn string) *Engine {
	return e.And(strColumn + " is NULL")
}

func (e *Engine) NotNULL(strColumn string) *Engine {
	return e.And(strColumn + " is not NULL")
}

func (e *Engine) jsonExpr(strColumn, strPath string) string {
	return fmt.Sprintf("`%s`->'$.%s'", strColumn, strPath)
}

func (e *Engine) JsonEqual(strColumn, strPath string, value any) *Engine {
	return e.And("%s = ?", e.jsonExpr(strColumn, strPath), value)
}

func (e *Engine) JsonGreater(strColumn, strPath string, value any) *Engine {
	return e.And("%s > ?", e.jsonExpr(strColumn, strPath), value)
}

func (e *Engine) JsonLess(strColumn, strPath string, value any) *Engine {
	return e.And("%s< ?", e.jsonExpr(strColumn, strPath), value)
}

func (e *Engine) JsonGreaterEqual(strColumn, strPath string, value any) *Engine {
	return e.And("%s >= ?", e.jsonExpr(strColumn, strPath), value)
}

func (e *Engine) JsonLessEqual(strColumn, strPath string, value any) *Engine {
	return e.And("%s <= ?", e.jsonExpr(strColumn, strPath), value)
}

func (e *Engine) NewID() ID {
	if e.idgen == nil {
		log.Panic("snowflake node id not set, please set it by Options when using NewEngine method")
	}
	return e.idgen.Generate()
}

func (e *Engine) ForUpdate() *Engine {
	if e.operType != types.OperType_Tx {
		log.Panic("this method is only for transaction")
	}
	return e.setForUpdate()
}

func (e *Engine) LockShareMode() *Engine {
	if e.operType != types.OperType_Tx {
		log.Panic("this method is only for transaction")
	}
	return e.setLockShareMode()
}

func (e *Engine) NewContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, types.SqlcaContextKey, e)
	return ctx
}

func NewContext(ctx context.Context, e *Engine) context.Context {
	ctx = context.WithValue(ctx, types.SqlcaContextKey, e)
	return ctx
}

func FromContext(ctx context.Context) *Engine {
	v := ctx.Value(types.SqlcaContextKey)
	return v.(*Engine)
}
