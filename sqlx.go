package sqlca

import (
	"database/sql"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const (
	TAG_NAME_DB       = "db"
	TAG_NAME_JSON     = "json"
	TAG_NAME_PROTOBUF = "protobuf"
	TAG_NAME_SQLCA    = "sqlca"
)

const (
	SQLCA_TAG_VALUE_AUTO_INCR = "autoincr" //auto increment
	SQLCA_TAG_VALUE_READ_ONLY = "readonly" //read only (eg. created_at)
	SQLCA_TAG_VALUE_IGNORE    = "-"        //ignore
)

const (
	PROTOBUF_VALUE_NAME = "name"
)

const (
	DRIVER_NAME_MYSQL    = "mysql"
	DRIVER_NAME_POSTGRES = "postgres"
	DRIVER_NAME_SQLITE   = "sqlite3"
	DRIVER_NAME_MSSQL    = "mssql"
	DRIVER_NAME_REDIS    = "redis"
)

const (
	DATABASE_KEY_NAME_WHERE      = "WHERE"
	DATABASE_KEY_NAME_UPDATE     = "UPDATE"
	DATABASE_KEY_NAME_SET        = "SET"
	DATABASE_KEY_NAME_FROM       = "FROM"
	DATABASE_KEY_NAME_DELETE     = "DELETE"
	DATABASE_KEY_NAME_SELECT     = "SELECT"
	DATABASE_KEY_NAME_DISTINCT   = "DISTINCT"
	DATABASE_KEY_NAME_IN         = "IN"
	DATABASE_KEY_NAME_NOT_IN     = "NOT IN"
	DATABASE_KEY_NAME_OR         = "OR"
	DATABASE_KEY_NAME_AND        = "AND"
	DATABASE_KEY_NAME_INSERT     = "INSERT INTO"
	DATABASE_KEY_NAME_VALUE      = "VALUE"
	DATABASE_KEY_NAME_VALUES     = "VALUES"
	DATABASE_KEY_NAME_FOR_UPDATE = "FOR UPDATE"
	DATABASE_KEY_NAME_ORDER_BY   = "ORDER BY"
	DATABASE_KEY_NAME_ASC        = "ASC"
	DATABASE_KEY_NAME_DESC       = "DESC"
	DATABASE_KEY_NAME_HAVING     = "HAVING"
	DATABASE_KEY_NAME_CASE       = "CASE"
	DATABASE_KEY_NAME_WHEN       = "WHEN"
	DATABASE_KEY_NAME_THEN       = "THEN"
	DATABASE_KEY_NAME_ELSE       = "ELSE"
	DATABASE_KEY_NAME_END        = "END"
)

type AdapterType int

const (
	ORDER_BY_ASC                  = "asc"
	ORDER_BY_DESC                 = "desc"
	DEFAULT_CAHCE_EXPIRE_SECONDS  = 24 * 60 * 60
	DEFAULT_PRIMARY_KEY_NAME      = "id"
	DEFAULT_SLOW_QUERY_ALERT_TIME = 500 //milliseconds
)

const (
	AdapterSqlx_MySQL    AdapterType = 1  //sqlx: mysql
	AdapterSqlx_Postgres AdapterType = 2  //sqlx: postgresql
	AdapterSqlx_Sqlite   AdapterType = 3  //sqlx: sqlite
	AdapterSqlx_Mssql    AdapterType = 4  //sqlx: mssql server
	AdapterCache_Redis   AdapterType = 11 //cache: redis
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
	case AdapterSqlx_Sqlite:
		return "AdapterSqlx_Sqlite"
	case AdapterSqlx_Mssql:
		return "AdapterSqlx_Mssql"
	case AdapterCache_Redis:
		return "AdapterCache_Redis"
	default:
	}
	return "Adapter_Unknown"
}

func (a AdapterType) DriverName() string {
	switch a {
	case AdapterSqlx_MySQL:
		return DRIVER_NAME_MYSQL
	case AdapterSqlx_Postgres:
		return DRIVER_NAME_POSTGRES
	case AdapterSqlx_Sqlite:
		return DRIVER_NAME_SQLITE
	case AdapterSqlx_Mssql:
		return DRIVER_NAME_MSSQL
	case AdapterCache_Redis:
		return DRIVER_NAME_REDIS
	default:
	}
	return "unknown"
}

var adapterNames = map[string]AdapterType{
	DRIVER_NAME_MYSQL:    AdapterSqlx_MySQL,
	DRIVER_NAME_POSTGRES: AdapterSqlx_Postgres,
	DRIVER_NAME_SQLITE:   AdapterSqlx_Sqlite,
	DRIVER_NAME_MSSQL:    AdapterSqlx_Mssql,
	DRIVER_NAME_REDIS:    AdapterCache_Redis,
}

func getAdapterType(name string) AdapterType {

	return adapterNames[name]
}

//------------------------------------------------------------------------------------------------------

type OperType int

const (
	OperType_Query     OperType = 1  // orm: query sql
	OperType_Update    OperType = 2  // orm: update sql
	OperType_Insert    OperType = 3  // orm: insert sql
	OperType_Upsert    OperType = 4  // orm: insert or update sql
	OperType_Tx        OperType = 5  // orm: tx sql
	OperType_QueryRaw  OperType = 6  // raw: query sql into model
	OperType_ExecRaw   OperType = 7  // raw: insert/update sql
	OperType_QueryMap  OperType = 8  // raw: query sql into map
	OperType_Delete    OperType = 9  // orm: delete sql
	OperType_ForUpdate OperType = 10 // orm: select ... for update sql
)

func (o OperType) GoString() string {
	return o.String()
}

func (o OperType) String() string {
	switch o {
	case OperType_Query:
		return "OperType_Query"
	case OperType_Update:
		return "OperType_Update"
	case OperType_Insert:
		return "OperType_Insert"
	case OperType_Upsert:
		return "OperType_Upsert"
	case OperType_QueryRaw:
		return "OperType_QueryRaw"
	case OperType_ExecRaw:
		return "OperType_ExecRaw"
	case OperType_Tx:
		return "OperType_Tx"
	case OperType_QueryMap:
		return "OperType_QueryMap"
	case OperType_Delete:
		return "OperType_Delete"
	}
	return "OperType_Unknown"
}

type ModelType int

const (
	ModelType_Struct   = 1
	ModelType_Slice    = 2
	ModelType_Map      = 3
	ModelType_BaseType = 4
)

func (m ModelType) GoString() string {
	return m.String()
}

func (m ModelType) String() string {
	switch m {
	case ModelType_Struct:
		return "ModelType_Struct"
	case ModelType_Slice:
		return "ModelType_Slice"
	case ModelType_Map:
		return "ModelType_Map"
	case ModelType_BaseType:
		return "ModelType_BaseType"
	}
	return "ModelType_Unknown"
}

type tableIndex struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type condition struct {
	ColumnName   string
	ColumnValues []interface{}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (e *Engine) appendMaster(db *sqlx.DB) {
	e.dbMasters = append(e.dbMasters, db)
	log.Debugf("db masters [%v]", len(e.dbMasters))
}

func (e *Engine) appendSlave(db *sqlx.DB) {
	e.dbSlaves = append(e.dbSlaves, db)
	log.Debugf("db slaves [%v]", len(e.dbSlaves))
}

// get slave db instance if use Slave() method to query, if not exist return a master db instance
func (e *Engine) getQueryDB() (db *sqlx.DB) {
	if e.slave {
		db = e.getSlave()
		if db != nil {
			return
		}
	}
	return e.getMaster()
}

// get a master db instance
func (e *Engine) getMaster() *sqlx.DB {

	n := len(e.dbMasters)
	if n > 0 {
		return e.dbMasters[rand.Intn(n)]
	}
	return nil
}

// get a slave db instance
func (e *Engine) getSlave() *sqlx.DB {
	n := len(e.dbSlaves)
	if n > 0 {
		return e.dbSlaves[rand.Intn(n)]
	}
	return e.getMaster()
}

func (e *Engine) setModel(models ...interface{}) *Engine {

	for _, v := range models {

		if v == nil {
			continue
		}
		typ := reflect.TypeOf(v)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		switch typ.Kind() {
		case reflect.Struct: // struct
			e.setModelType(ModelType_Struct)
		case reflect.Slice: //  slice
			e.setModelType(ModelType_Slice)
		case reflect.Map: // map
			e.setModelType(ModelType_Map)
		default: //base type
			e.setModelType(ModelType_BaseType)
		}
		if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
			e.model = models[0] //map, struct or slice
		} else {
			e.model = models //base type argument like int/string/float32...
		}
		var selectColumns []string
		e.dict = newReflector(e, e.model).ToMap(e.dbTags...)
		for k, _ := range e.dict {
			selectColumns = append(selectColumns, k)
		}
		if len(selectColumns) == 0 {
			e.setSelectColumns("*")
		} else {
			e.setSelectColumns(selectColumns...)
		}
		//log.Debugf("dict [%+v] select columns %+v model type [%+v]", e.dict, selectColumns, e.modelType)
		break //only check first argument
	}
	return e
}

// clone engine
func (e *Engine) clone(models ...interface{}) *Engine {

	engine := &Engine{
		dsn:             e.dsn,
		dbMasters:       e.dbMasters,
		dbSlaves:        e.dbSlaves,
		cache:           e.cache,
		adapterSqlx:     e.adapterSqlx,
		adapterCache:    e.adapterCache,
		strPkName:       e.strPkName,
		expireTime:      e.expireTime,
		strDatabaseName: e.strDatabaseName,
		dbTags:          e.dbTags,
		bForce:          e.bForce,
		bAutoRollback:   e.bAutoRollback,
		slowQueryOn:     e.slowQueryOn,
		slowQueryTime:   e.slowQueryTime,
	}

	engine.setModel(models...)
	return engine
}

func (e *Engine) newTx() (txEngine *Engine, err error) {

	txEngine = e.clone()
	db := e.getMaster()
	if txEngine.tx, err = db.Begin(); err != nil {
		log.Errorf("newTx error [%+v]", err.Error())
		return nil, err
	}
	txEngine.operType = OperType_Tx
	return
}

func (e *Engine) postgresQueryInsert(strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	strSQL += fmt.Sprintf(" RETURNING \"%v\"", e.GetPkName())
	log.Debugf("[%v]", strSQL)
	db := e.getMaster()
	if rows, err = db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (e *Engine) postgresQueryUpsert(strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	log.Debugf("[%v]", strSQL)
	db := e.getMaster()
	if rows, err = db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (e *Engine) mssqlQueryInsert(strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	strSQL += " SELECT SCOPE_IDENTITY() AS last_insert_id"
	log.Debugf("[%v]", strSQL)
	db := e.getMaster()
	if rows, err = db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (e *Engine) mssqlUpsert(strSQL string) (lastInsertId int64, err error) {

	var db *Engine
	var query = e.makeSqlxQueryPrimaryKey()
	if db, err = e.TxBegin(); err != nil {
		log.Errorf("TxBegin error [%v]", err.Error())
		return
	}
	var count int64
	if count, err = db.TxGet(&lastInsertId, query); err != nil {
		log.Errorf("TxGet [%v] error [%v]", query, err.Error())
		_ = db.TxRollback()
		return
	}
	if count == 0 {
		// INSERT INTO users(...) values(...)  SELECT SCOPE_IDENTITY() AS last_insert_id
		//if _, _, err = db.TxExec(strSQL); err != nil
		if lastInsertId, err = e.mssqlQueryInsert(strSQL); err != nil {
			log.Errorf("mssqlQueryInsert [%v] error [%v]", strSQL, err.Error())
			_ = db.TxRollback()
			return
		}
	} else {
		// UPDATE users SET xxx=yyy WHERE id=nnn
		strUpdates := fmt.Sprintf("%v %v %v %v %v %v=%v",
			DATABASE_KEY_NAME_UPDATE, e.getTableName(),
			DATABASE_KEY_NAME_SET, e.getOnConflictDo(),
			DATABASE_KEY_NAME_WHERE, e.GetPkName(), lastInsertId)
		log.Debugf("%v", strUpdates)
		if _, _, err = db.TxExec(strUpdates); err != nil {
			log.Errorf("TxExec [%v] error [%v]", strSQL, err.Error())
			_ = db.TxRollback()
			return
		}
	}

	if err = db.TxCommit(); err != nil {
		log.Errorf("TxCommit [%v] error [%v]", strSQL, err.Error())
		return
	}
	return
}

func (e *Engine) getDistinct() string {
	return e.strDistinct
}

func (e *Engine) setDistinct() {
	e.strDistinct = DATABASE_KEY_NAME_DISTINCT
}

func (e *Engine) setUseCache(enable bool) {
	e.bUseCache = enable
}

func (e *Engine) getUseCache() bool {
	return e.bUseCache
}

func (e *Engine) sepStrByDot(strIn string) (strPrefix, strSuffix string) {
	strSuffix = strIn
	nIndex := strings.LastIndex(strIn, ".")
	if nIndex == -1 || nIndex == 0 {
		return
	} else {
		strPrefix = strIn[:nIndex+1] //contain . character
		if nIndex < len(strIn) {
			strSuffix = strIn[nIndex+1:]
		}
	}

	return
}

func (e *Engine) getModelValue(strKey string) interface{} {

	_, strCol := e.sepStrByDot(strKey)
	//log.Debugf("getModelValue raw=%v col=%v dict=%+v", strKey, strCol, e.dict)
	return e.dict[strCol]
}

func (e *Engine) setIndexes(name string, value interface{}) {

	if name == e.GetPkName() {
		return
	}
	typ := reflect.TypeOf(value)
	switch typ.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Func, reflect.Ptr, reflect.Chan, reflect.UnsafePointer:
		log.Errorf("index value type [%v] illegal", typ.Kind())
	}
	e.cacheIndexes = append(e.cacheIndexes, tableIndex{
		Name:  name,
		Value: value,
	})
}

func (e *Engine) getIndexes() []tableIndex {
	return e.cacheIndexes
}

func (e *Engine) setDatabaseName(strName string) {
	e.strDatabaseName = strName
}

func (e *Engine) getDatabaseName() string {
	return e.strDatabaseName
}

func (e *Engine) getTableName() string {
	return e.strTableName
}

func (e *Engine) setTableName(strNames ...string) {
	e.strTableName = strings.Join(strNames, ",")
}

func (e *Engine) setPkValue(value interface{}) {

	var strValue string
	switch value.(type) {
	case string:
		strValue = value.(string)
		if strValue == "" {
			log.Debugf("primary key's value is nil")
		}
	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
		{
			strValue = fmt.Sprintf("%v", value)
			if strValue == "0" {
				log.Debugf("primary key's value is 0")
			}
		}
	default:
		log.Errorf("primary key's value type illegal")
	}
	e.strPkValue = strValue
}

func (e *Engine) getPkValue() string {
	if e.strPkValue == "" {

		modelValue := e.getModelValue(e.GetPkName())
		if modelValue != nil {
			e.strPkValue = fmt.Sprintf("%v", modelValue) //get primary value from model dictionary
		}
	}
	return e.strPkValue
}

func (e *Engine) setSelectColumns(strColumns ...string) {
	e.selectColumns = strColumns
}

func (e *Engine) setExcludeColumns(strColumns ...string) {
	e.excludeColumns = strColumns
}

func (e *Engine) getSelectColumns() (strColumns []string) {
	return e.selectColumns
}

func (e *Engine) setAscColumns(strColumns ...string) {
	e.ascColumns = strColumns
}

func (e *Engine) setDescColumns(strColumns ...string) {
	e.descColumns = strColumns
}

// SELECT ... FROM xxx ORDER BY c1, c2 ASC, c3 DESC
func (e *Engine) getAscAndDesc() (strAscDesc string) {

	var ss []string
	//make ASC expression
	if len(e.ascColumns) > 0 {
		ss = append(ss, fmt.Sprintf("%v %v", strings.Join(e.ascColumns, ","), DATABASE_KEY_NAME_ASC))
	}
	//make DESC expression
	if len(e.descColumns) > 0 {
		ss = append(ss, fmt.Sprintf("%v %v", strings.Join(e.descColumns, ","), DATABASE_KEY_NAME_DESC))
	}
	//make default order by expression
	if len(ss) == 0 {
		ss = append(ss, fmt.Sprintf("%v %v", strings.Join(e.orderByColumns, ","), DATABASE_KEY_NAME_ASC))
	}
	return strings.Join(ss, ",")
}

// use Where function to set custom where condition
func (e *Engine) getCustomWhere() string {
	return e.strWhere
}

func (e *Engine) getIndexWhere() (strCondition string) {
	if e.getOperType() == OperType_Query && len(e.getIndexes()) > 0 {

		var conditions []string
		for _, v := range e.getIndexes() {
			cond := fmt.Sprintf("%v=%v",
				e.getQuoteColumnName(v.Name), e.getQuoteColumnValue(v.Value))
			conditions = append(conditions, cond)
		}
		strCondition = strings.Join(conditions, " AND ")
	}
	return
}

// primary key value like 'id'=xxx condition
func (e *Engine) getPkWhere() (strCondition string) {

	if e.isPkValueNil() {
		log.Debugf("query condition primary key or index is nil")
		return
	}
	strCondition = fmt.Sprintf("%v=%v",
		e.getQuoteColumnName(e.GetPkName()), e.getQuoteColumnValue(e.getPkValue()))
	return
}

func (e *Engine) getQuoteColumnName(v string) (strColumn string) {
	return fmt.Sprintf("%v%v%v", e.getForwardQuote(), v, e.getBackQuote())
}

func (e *Engine) getQuoteColumnValue(v interface{}) (strValue string) {
	return fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote())
}

func (e *Engine) setWhere(strWhere string) {
	e.strWhere = strWhere
}

func (e *Engine) getModelType() ModelType {
	return e.modelType
}

func (e *Engine) setModelType(modelType ModelType) {
	e.modelType = modelType
}

func (e *Engine) getOperType() OperType {
	return e.operType
}

func (e *Engine) setOperType(operType OperType) {
	e.operType = operType
}

func (e *Engine) setConflictColumns(strColumns ...string) {
	e.conflictColumns = strColumns
}

func (e *Engine) getConflictColumns() []string {
	return e.conflictColumns
}

func (e *Engine) getSingleQuote() (strQuote string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		return "'"
	case AdapterSqlx_Postgres:
		return "'"
	case AdapterSqlx_Mssql:
		return "'"
	}
	return "'"
}

func (e *Engine) getForwardQuote() (strQuote string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		return "`"
	case AdapterSqlx_Postgres:
		return "\""
	case AdapterSqlx_Mssql:
		return "["
	}
	return "`"
}

func (e *Engine) getBackQuote() (strQuote string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		return "`"
	case AdapterSqlx_Postgres:
		return "\""
	case AdapterSqlx_Mssql:
		return "]"
	}
	return "`"
}

func (e *Engine) getOnConflictForwardKey() (strKey string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		return "ON DUPLICATE"
	case AdapterSqlx_Postgres:
		return "ON CONFLICT ("
	case AdapterSqlx_Mssql:
		return ""
	}
	return
}

func (e *Engine) getOnConflictBackKey() (strKey string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		return "KEY UPDATE"
	case AdapterSqlx_Postgres:
		return ") DO UPDATE SET"
	case AdapterSqlx_Mssql:
		return ""
	}
	return
}

func (e *Engine) setLimit(strLimit string) {
	e.strLimit = strLimit
}

func (e *Engine) getLimit() string {
	return e.strLimit
}

func (e *Engine) setOffset(strOffset string) {
	e.strOffset = strOffset
}

func (e *Engine) getOffset() string {
	return e.strOffset
}

func (e *Engine) setOrderBy(strColumns ...string) {
	e.orderByColumns = strColumns
}

func (e *Engine) getOrderBy() (strOrderBy string) {

	if isNilOrFalse(e.orderByColumns) && isNilOrFalse(e.ascColumns) && isNilOrFalse(e.descColumns) {
		return
	}
	return fmt.Sprintf("%v %v", DATABASE_KEY_NAME_ORDER_BY, e.getAscAndDesc())
}

func (e *Engine) setGroupBy(strColumns ...string) {
	e.groupByColumns = strColumns
}

func (e *Engine) setHaving(havingCondition string) {
	e.havingCondition = havingCondition
}

func (e *Engine) getHaving() (strHaving string) {

	if isNilOrFalse(e.havingCondition) {
		return
	}
	return fmt.Sprintf("%v %v", DATABASE_KEY_NAME_HAVING, e.havingCondition)
}

func (e *Engine) getGroupBy() (strGroupBy string) {
	if isNilOrFalse(e.groupByColumns) {
		return
	}
	return fmt.Sprintf(" GROUP BY %v", strings.Join(e.groupByColumns, ","))
}

func (e *Engine) isPkValueNil() bool {

	if e.getPkValue() == "" || e.getPkValue() == "0" {
		return true
	}

	return false
}

func (e *Engine) isPkInteger() bool {

	id := e.getModelValue(e.GetPkName())
	switch id.(type) {
	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
		return true
	}
	return false
}

func (e *Engine) isReadOnly(strIn string) bool {

	if e.bForce {
		return false
	}
	for _, v := range e.readOnly {
		if v == strIn {
			return true
		}
	}
	return false
}

func (e *Engine) isExcluded(strCol string) bool {

	for _, v := range e.excludeColumns {
		if v == strCol {
			return true
		}
	}
	return false
}

func (e *Engine) isSelected(strCol string) bool {

	for _, v := range e.selectColumns {

		if v == strCol {
			return true
		}
	}
	return false
}

func (e *Engine) isExcepted(strCol string, strExcepts ...string) bool {

	for _, v := range strExcepts {
		if v == strCol {
			return true
		}
	}
	return false
}

func (e *Engine) isColumnSelected(strCol string, strExcepts ...string) bool {

	if e.isExcepted(strCol, strExcepts...) {
		return false
	}

	if len(e.selectColumns) == 0 {
		return true
	}

	if e.isExcluded(strCol) {
		return false
	}

	if e.isSelected(strCol) {
		return true
	}
	return false
}

func (e *Engine) getQuoteConflicts() (strQuoteConflicts string) {

	if e.adapterSqlx != AdapterSqlx_Postgres {
		return //only postgres need conflicts fields
	}

	assert(e.conflictColumns, "on conflict columns is nil")

	var cols []string

	for _, v := range e.getConflictColumns() {

		c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), v, e.getBackQuote()) // postgresql conflict column name format to `id`,...
		cols = append(cols, c)
	}

	if len(cols) > 0 {
		strQuoteConflicts = strings.Join(cols, ",")
	}
	return
}

func (e *Engine) getRawColumns() (strColumns string) {
	var selectCols []string
	if len(e.selectColumns) == 0 {
		return "*"
	}
	for _, v := range e.selectColumns {
		if e.isColumnSelected(v) {
			selectCols = append(selectCols, v)
		}
	}
	if len(selectCols) > 0 {
		if e.strCaseWhen != "" {
			selectCols = append(selectCols, e.strCaseWhen)
		}
		strColumns = strings.Join(selectCols, ",")
	}
	return
}

func (e *Engine) getQuoteUpdates(strColumns []string, strExcepts ...string) (strUpdates string) {

	var cols []string
	for _, v := range strColumns {

		if e.isColumnSelected(v, strExcepts...) && !e.isReadOnly(v) {
			val := e.getModelValue(v)
			if val == nil {
				//log.Warnf("column [%v] selected but have no value", v)
				continue
			}
			strVal := handleSpecialChars(fmt.Sprintf("%v", val))
			c := fmt.Sprintf("%v=%v", e.getQuoteColumnName(v), e.getQuoteColumnValue(strVal)) // column name format to `date`='1583055138',...
			cols = append(cols, c)
		}
	}

	if len(cols) == 0 {
		//may be model is a base type slice
		args := e.model.([]interface{})
		count := len(args)
		//log.Debugf("args count [%v] values [%+v]", count, args)
		for i, k := range strColumns {
			if i < count {
				v := args[i]
				val := reflect.ValueOf(v)
				//log.Debugf("columns[%v] name [%v] value [%v]", i, k, val.Elem().Interface())
				c := fmt.Sprintf("%v=%v", e.getQuoteColumnName(k), e.getQuoteColumnValue(val.Elem().Interface())) // column name format to `date`='1583055138',...
				cols = append(cols, c)
			}
		}
	}

	if len(cols) > 0 {
		strUpdates = strings.Join(cols, ",")
	}
	return
}

func (e *Engine) getOnConflictDo() (strDo string) {
	switch e.adapterSqlx {
	case AdapterSqlx_MySQL, AdapterSqlx_Sqlite:
		{
			strUpdates := e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
			if !isNilOrFalse(strUpdates) {
				if e.isPkInteger() { // primary key type is a integer
					strDo = fmt.Sprintf("%v", strUpdates)
				} else {
					strDo = strUpdates
				}
			}
		}
	case AdapterSqlx_Postgres:
		{
			strUpdates := e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
			if !isNilOrFalse(strUpdates) {
				strDo = fmt.Sprintf("%v RETURNING \"%v\"", strUpdates, e.GetPkName()) // TODO @libin test postgresql ON CONFLICT(...) DO UPDATE SET ... RETURNING id
			}
		}
	case AdapterSqlx_Mssql:
		{
			strDo = e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
		}
	}
	return
}

func (e *Engine) isContainInts(i int, values []int) bool {
	for _, v := range values {
		if v == i {
			return true
		}
	}
	return false
}

func (e *Engine) getInsertColumnsAndValues() (strQuoteColumns, strColonValues string) {
	var cols, vals []string

	typ := reflect.TypeOf(e.model)
	kind := typ.Kind()

	if kind == reflect.Ptr {
		typ = typ.Elem()
	}
	//log.Debugf("reflect.TypeOf(e.model) = %v", typ.Kind())
	if typ.Kind() == reflect.Slice {
		var cols2 []string
		var values [][]string
		var valueQuoteSlice []string
		var excludeIndexes []int
		cols2, values = e.getStructSliceKeyValues(true)
		for i, v := range cols2 {

			if e.isReadOnly(v) || e.isExcluded(v) {
				excludeIndexes = append(excludeIndexes, i) //index of exclude columns slice
				continue
			}

			k := fmt.Sprintf("%v%v%v", e.getForwardQuote(), v, e.getBackQuote()) // column name format to `id`,...
			cols = append(cols, k)
		}

		for _, v := range values {
			var valueQoute []string
			for ii, vv := range v {
				if e.isContainInts(ii, excludeIndexes) {
					continue
				}
				vq := fmt.Sprintf("%v%v%v", e.getSingleQuote(), vv, e.getSingleQuote()) // column value format
				valueQoute = append(valueQoute, vq)
			}
			valueQuoteSlice = append(valueQuoteSlice, fmt.Sprintf("(%v)", strings.Join(valueQoute, ",")))
		}

		if len(valueQuoteSlice) > 0 {
			strColonValues = strings.Join(valueQuoteSlice, ",")
		}
	} else {
		for k, v := range e.dict {

			if e.isReadOnly(k) || e.isExcluded(k) {
				continue
			}

			c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), k, e.getBackQuote()) // column name format to `id`,...

			if k == e.GetPkName() && e.isPkValueNil() {
				continue
			}
			vq := fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote()) // column value format
			cols = append(cols, c)
			vals = append(vals, vq)
		}
		strColonValues = fmt.Sprintf("(%v)", strings.Join(vals, ","))
	}

	strQuoteColumns = fmt.Sprintf("(%v)", strings.Join(cols, ","))
	//log.Debugf("strQuoteColumns [%v] strColonValues [%v]", strQuoteColumns, strColonValues)
	return
}

func (e *Engine) getOnConflictUpdates(strExcepts ...string) (strUpdates string) {

	//mysql/sqlite: ON DUPLICATE KEY UPDATE id=last_insert_id(id), date='1582980944'...
	//postgres: ON CONFLICT (id) DO UPDATE SET date='1582980944'...
	//mssql: nothing...
	strUpdates = fmt.Sprintf("%v %v %v %v",
		e.getOnConflictForwardKey(), e.getQuoteConflicts(), e.getOnConflictBackKey(), e.getOnConflictDo())
	return
}

func (e *Engine) formatString(strIn string, args ...interface{}) (strFmt string) {
	strFmt = strIn
	if e.isQuestionPlaceHolder(strIn, args...) { //question placeholder exist
		strFmt = strings.Replace(strFmt, "?", "'%v'", -1)
	}
	return fmt.Sprintf(strFmt, args...)
}

func (e *Engine) makeSqlxQueryPrimaryKey() (strSql string) {

	strSql = fmt.Sprintf("%v %v%v%v %v %v %v %v%v%v=%v%v%v",
		DATABASE_KEY_NAME_SELECT, e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		DATABASE_KEY_NAME_FROM, e.getTableName(), DATABASE_KEY_NAME_WHERE,
		e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		e.getSingleQuote(), e.getPkValue(), e.getSingleQuote())
	return
}

func (e *Engine) makeSqlxString() (strSql string) {

	switch e.operType {
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
	default:
		log.Errorf("operation illegal")
	}
	strSql = strings.TrimSpace(strSql)
	log.Debugf("[%v] SQL [%s]", e.operType, strSql)

	return
}

func (e *Engine) makeInCondition(cond condition) (strCondition string) {

	var strValues []string
	for _, v := range cond.ColumnValues {
		strValues = append(strValues, fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote()))
	}
	strCondition = fmt.Sprintf("%v %v (%v)", cond.ColumnName, DATABASE_KEY_NAME_IN, strings.Join(strValues, ","))
	return
}

func (e *Engine) makeNotCondition(cond condition) (strCondition string) {

	var strValues []string
	for _, v := range cond.ColumnValues {
		strValues = append(strValues, fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote()))
	}
	strCondition = fmt.Sprintf("%v %v (%v)", cond.ColumnName, DATABASE_KEY_NAME_NOT_IN, strings.Join(strValues, ","))
	return
}

func (e *Engine) makeWhereCondition() (strWhere string) {

	if e.isPkValueNil() {
		strIndexCond := e.getIndexWhere()
		if strIndexCond != "" {
			strWhere += e.getIndexWhere()
		}
	} else {
		strWhere += e.getPkWhere()
	}

	if strWhere == "" {
		strCustomer := e.getCustomWhere()
		if strCustomer == "" {
			//where condition required when update or delete
			if e.operType != OperType_Update && e.operType != OperType_Delete {
				strWhere += "1=1"
			} else {
				log.Warnf("where condition required when use orm update or delete")
			}
		} else {
			strWhere += strCustomer
		}
	}

	for _, v := range e.andConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, v)
	}
	for _, v := range e.inConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, e.makeInCondition(v))
	}
	for _, v := range e.notConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, e.makeNotCondition(v))
	}
	for _, v := range e.orConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_OR, v)
	}
	strWhere = DATABASE_KEY_NAME_WHERE + " " + strWhere
	return
}

func (e *Engine) makeSqlxQuery() (strSqlx string) {
	strWhere := e.makeWhereCondition()

	switch e.adapterSqlx {
	case AdapterSqlx_Mssql:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getLimit(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy())
	default:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy(), e.getLimit(), e.getOffset())
	}

	return
}

func (e *Engine) makeSqlxForUpdate() (strSqlx string) {
	return e.makeSqlxQuery() + " " + DATABASE_KEY_NAME_FOR_UPDATE
}

func (e *Engine) makeSqlxUpdate() (strSqlx string) {

	strWhere := e.makeWhereCondition()
	strSqlx = fmt.Sprintf("%v %v %v %v %v %v",
		DATABASE_KEY_NAME_UPDATE, e.getTableName(), DATABASE_KEY_NAME_SET,
		e.getQuoteUpdates(e.getSelectColumns(), e.GetPkName()), strWhere, e.getLimit())
	assert(strSqlx, "update sql is nil")
	return
}

func (e *Engine) makeSqlxInsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strSqlx = fmt.Sprintf("%v %v %v %v %v", DATABASE_KEY_NAME_INSERT, e.getTableName(), strColumns, DATABASE_KEY_NAME_VALUES, strValues)
	return
}

func (e *Engine) makeSqlxUpsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strOnConflictUpdates := e.getOnConflictUpdates()
	strSqlx = fmt.Sprintf("%v %v %v %v %v %v", DATABASE_KEY_NAME_INSERT, e.getTableName(), strColumns, DATABASE_KEY_NAME_VALUES, strValues, strOnConflictUpdates)
	return
}

func (e *Engine) makeSqlxDelete() (strSqlx string) {
	strWhere := e.makeWhereCondition()
	if strWhere == "" {
		panic("no condition to delete records")
	}
	strSqlx = fmt.Sprintf("%v %v %v %v", DATABASE_KEY_NAME_DELETE, DATABASE_KEY_NAME_FROM, e.getTableName(), strWhere)
	return
}

func (e *Engine) isQuestionPlaceHolder(query string, args ...interface{}) bool {
	count := strings.Count(query, "?")
	if count > 0 && count == len(args) {
		return true
	}
	return false
}

func (e *Engine) cleanWhereCondition() {
	e.strWhere = ""
	e.strPkValue = ""
	e.cacheIndexes = nil
}

func (e *Engine) autoRollback() {
	if e.bAutoRollback && e.operType == OperType_Tx && e.tx != nil {
		_ = e.tx.Rollback()
		log.Debugf("tx auto rollback successful")
	}
}
