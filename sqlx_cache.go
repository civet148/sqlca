package sqlca

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"reflect"
	"strings"
)

const (
	TAG_NAME_DB    = "db"
	TAG_NAME_JSON  = "json"
	TAG_NAME_REDIS = "redis"
)

type AdapterType int

const (
	ORDER_BY_ASC                 = "asc"
	ORDER_BY_DESC                = "desc"
	DEFAULT_CAHCE_EXPIRE_SECONDS = 60 * 60
	DEFAULT_PRIMARY_KEY_NAME     = "id"
	SQLX_IGNORE_CREATED_AT       = "created_at"
	SQLX_IGNORE_UPDATED_AT       = "updated_at"
)

const (
	AdapterSqlx_MySQL      AdapterType = 1  //sqlx: mysql
	AdapterSqlx_Postgres   AdapterType = 2  //sqlx: postgresql
	AdapterSqlx_Sqlite     AdapterType = 3  //sqlx: sqlite
	AdapterSqlx_Mssql      AdapterType = 4  //sqlx: mssql server
	AdapterCache_Redis     AdapterType = 11 //cache: redis
	AdapterCache_Memcached AdapterType = 12 //cache: memcached
	AdapterCache_Memory    AdapterType = 13 //cache: memory
	AdapterCache_File      AdapterType = 14 //cache: file
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

func (a AdapterType) Schema() string {
	switch a {
	case AdapterSqlx_MySQL:
		return "mysql"
	case AdapterSqlx_Postgres:
		return "postges"
	case AdapterSqlx_Sqlite:
		return "sqlite"
	case AdapterSqlx_Mssql:
		return "mssql"
	case AdapterCache_Redis:
		return "redis"
	case AdapterCache_Memcached:
		return "memcached"
	case AdapterCache_Memory:
		return "memory"
	case AdapterCache_File:
		return "file"
	default:
	}
	return "unknown"
}

type OperType int

const (
	OperType_Query    OperType = 1 // orm: query sql
	OperType_Update   OperType = 2 // orm: update sql
	OperType_Insert   OperType = 3 // orm: insert sql
	OperType_Upsert   OperType = 4 // orm: insert or update sql
	OperType_QueryRaw OperType = 5 // raw: query sql
	OperType_ExecRaw  OperType = 6 // raw: insert/update sql
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
		return "ModelType_StructSlice"
	case ModelType_Map:
		return "ModelType_Map"
	case ModelType_BaseType:
		return "ModelType_BaseType"
	}
	return "ModelType_Unknown"
}

// clone engine
func (e *Engine) clone(models ...interface{}) *Engine {

	engine := &Engine{
		db:           e.db,
		cache:        e.cache,
		debug:        e.debug,
		adapterSqlx:  e.adapterSqlx,
		adapterCache: e.adapterCache,
		strPkName:    e.strPkName,
		expireTime:   e.expireTime,
	}

	for _, v := range models {

		typ := reflect.TypeOf(v)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		switch typ.Kind() {
		case reflect.Struct: // struct
			engine.setModelType(ModelType_Struct)
		case reflect.Slice: //  slice
			engine.setModelType(ModelType_Slice)
		case reflect.Map: // map
			//TODO @libin support map type model
			engine.setModelType(ModelType_Map)
			assert(false, "map type model not support yet")
		default: //base type
			engine.setModelType(ModelType_BaseType)
		}
		if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Slice {
			engine.model = models[0] //map, struct or slice
		} else {
			engine.model = models //base type argument like int/string/float32...
		}
		break //only check first argument
	}

	engine.dict = newReflector(engine.model).ToMap(TAG_NAME_DB)
	for k, _ := range engine.dict {
		engine.selectColumns = append(engine.selectColumns, k)
	}

	return engine
}

func (e *Engine) clean() *Engine {

	return e
}

func (e *Engine) getTableName() string {
	return e.strTableName
}

func (e *Engine) setTableName(strName string) {
	e.strTableName = strName
}

func (e *Engine) getPkValue() string {
	return e.strPkValue
}

func (e *Engine) setPkValue(strValue string) {
	e.strPkValue = strValue
}

func (e *Engine) setSelectColumns(strColumns ...string) {
	e.selectColumns = strColumns
}

func (e *Engine) getSelectColumns() (strColumns []string) {
	return e.selectColumns
}

func (e *Engine) setAscOrDesc(strSort string) {
	e.strAscOrDesc = strSort
}

func (e *Engine) getAscOrDesc() string {
	return e.strAscOrDesc
}

// use Where function to set custom where condition
func (e *Engine) getCustomWhere() string {
	return e.strWhere
}

// primary key value like 'id'=xxx condition
func (e *Engine) getPkWhere() (strPkCondition string) {

	strPkName := e.GetPkName()
	strPkValue := e.getPkValue()
	if isNilOrFalse(strPkValue) {
		//use model primary value
		strPkCondition = fmt.Sprintf("%v%v%v=%v%v%v", e.getForwardQuote(), strPkName, e.getBackQuote(), e.getSingleQuote(), e.dict[strPkName], e.getSingleQuote())
	} else {
		//use custom primary value
		strPkCondition = fmt.Sprintf("%v%v%v=%v%v%v", e.getForwardQuote(), strPkName, e.getBackQuote(), e.getSingleQuote(), strPkValue, e.getSingleQuote())
	}
	return
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

// get data base driver name and data source name
func (e *Engine) getConnUrl(adapterType AdapterType, strUrl string) (strScheme, strDSN string) {
	//TODO @libin parse connect url for database
	strScheme = adapterType.Schema()
	// TODO parse url ....
	//
	//switch e.adapterSqlx {
	//case AdapterSqlx_MySQL:
	//case AdapterSqlx_Postgres:
	//case AdapterSqlx_Sqlite:
	//case AdapterSqlx_Mssql:
	//}
	return strScheme, strUrl
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
	return
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
	return
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
	return
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
		return " "
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

	if isNilOrFalse(e.orderByColumns) {
		return
	}
	return fmt.Sprintf("ORDER BY %v %v", strings.Join(e.orderByColumns, ","), e.getAscOrDesc())
}

func (e *Engine) setGroupBy(strColumns ...string) {
	e.groupByColumns = strColumns
}

func (e *Engine) getGroupBy() (strGroupBy string) {
	if isNilOrFalse(e.groupByColumns) {
		return
	}
	return fmt.Sprintf(" GROUP BY %v", strings.Join(e.groupByColumns, ","))
}

func (e *Engine) isColumnSelected(strCol string, strExcepts ...string) bool {

	for _, v := range strExcepts {
		if v == strCol {
			return false
		}
	}

	if len(e.selectColumns) == 0 {
		return true
	}

	for _, v := range e.selectColumns {

		if v == strCol {
			return true
		}
	}
	return false
}

func (e *Engine) getQuoteConflicts(strExcepts ...string) (strQuoteConflicts string) {

	if e.adapterSqlx != AdapterSqlx_Postgres {
		return //only postgres need conflicts fields
	}

	assert(e.conflictColumns, "on conflict columns is nil")

	var cols []string

	for _, v := range e.conflictColumns {

		if e.isColumnSelected(v, strExcepts...) {
			c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), v, e.getBackQuote()) // postgresql conflict column name format to `id`,...
			cols = append(cols, c)
		}
	}

	if len(cols) > 0 {
		strQuoteConflicts = strings.Join(cols, ",")
	}
	return
}

func (e *Engine) getQuoteColumns() (strColumns string) {
	var cols []string

	selectCols := e.getSelectColumns()
	if len(selectCols) == 0 {
		return "*"
	}

	for _, v := range selectCols {

		if v == "*" {
			return "*"
		}

		if e.isColumnSelected(v) {
			c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), v, e.getBackQuote()) // column name format to `id`,`date`,...
			cols = append(cols, c)
		}
	}

	if len(cols) > 0 {
		strColumns = strings.Join(cols, ",")
	}
	return
}

func (e *Engine) getQuoteUpdates(strColumns []string, strExcepts ...string) (strUpdates string) {

	var cols []string
	for _, v := range strColumns {

		if e.isColumnSelected(v, strExcepts...) {
			c := fmt.Sprintf("%v%v%v=%v%v%v", e.getForwardQuote(), v, e.getBackQuote(), e.getSingleQuote(), e.dict[v], e.getSingleQuote()) // column name format to `date`='1583055138',...
			cols = append(cols, c)
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
			strDo = fmt.Sprintf("`%v`=LAST_INSERT_ID(`%v`)", e.strPkName, e.strPkName)
			strUpdates := e.getQuoteUpdates(e.getSelectColumns(), e.strPkName, SQLX_IGNORE_CREATED_AT, SQLX_IGNORE_UPDATED_AT)
			if !isNilOrFalse(strUpdates) {
				strDo = fmt.Sprintf("%v, %v", strDo, strUpdates)
			}
		}
	case AdapterSqlx_Postgres:
		{
			strUpdates := e.getQuoteUpdates(e.getSelectColumns(), e.strPkName, SQLX_IGNORE_CREATED_AT, SQLX_IGNORE_UPDATED_AT)
			if !isNilOrFalse(strUpdates) {
				strDo = fmt.Sprintf("%v RETURNING %v", strUpdates, e.GetPkName()) // TODO @libin test postgresql ON CONFLICT(...) DO UPDATE SET ... RETURNING id
			}
		}
	case AdapterSqlx_Mssql:
		{
			// TODO @libin MSSQL Server upsert...do...
		}
	}
	return
}

func (e *Engine) getInsertColumnsAndValues() (strQuoteColumns, strColonValues string) {
	var cols, vals []string

	for k, _ := range e.dict {

		if k == SQLX_IGNORE_CREATED_AT || k == SQLX_IGNORE_UPDATED_AT {
			continue
		}
		c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), k, e.getBackQuote()) // column name format to `id`,...
		v := fmt.Sprintf(":%v", k)                                           // column value format to :id,...
		cols = append(cols, c)
		vals = append(vals, v)
	}
	strQuoteColumns = strings.Join(cols, ",")
	strColonValues = strings.Join(vals, ",")
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

func (e *Engine) makeSqlxString() (strSqlx string) {

	switch e.operType {
	case OperType_Query:
		strSqlx = e.makeSqlxQuery()
	case OperType_Update:
		strSqlx = e.makeSqlxUpdate()
	case OperType_Insert:
		strSqlx = e.makeSqlxInsert()
	case OperType_Upsert:
		strSqlx = e.makeSqlxUpsert()
	default:
		assert(false, "operation illegal")
	}

	log.Debugf("sqlx query [%s]", strSqlx)

	return
}

func (e *Engine) makeSqlxQuery() (strSqlx string) {

	if isNilOrFalse(e.getCustomWhere()) {
		strSqlx = fmt.Sprintf("SELECT %v FROM %v WHERE %v %v %v %v %v",
			e.getQuoteColumns(), e.getTableName(), e.getPkWhere(), e.getOrderBy(), e.getGroupBy(), e.getLimit(), e.getOffset()) //where condition by model primary key value
	} else {
		strSqlx = fmt.Sprintf("SELECT %v FROM %v WHERE %v %v %v %v %v",
			e.getQuoteColumns(), e.getTableName(), e.getCustomWhere(), e.getOrderBy(), e.getGroupBy(), e.getLimit(), e.getOffset()) //where condition by custom where condition from Where()
	}
	assert(strSqlx, "query sql is nil")
	return
}

func (e *Engine) makeSqlxUpdate() (strSqlx string) {

	if isNilOrFalse(e.getCustomWhere()) {
		//where condition by model primary key value (not include primary key `id` and created_at/updated_at)
		strSqlx = fmt.Sprintf("UPDATE %v SET %v WHERE %v %v",
			e.getTableName(),
			e.getQuoteUpdates(e.getSelectColumns(), e.GetPkName(), SQLX_IGNORE_CREATED_AT, SQLX_IGNORE_UPDATED_AT),
			e.getPkWhere(),
			e.getLimit())
	} else {
		//where condition by custom condition (not include primary key like `id` and created_at/updated_at)
		strSqlx = fmt.Sprintf("UPDATE %v SET %v WHERE %v %v",
			e.getTableName(),
			e.getQuoteUpdates(e.getSelectColumns(), e.GetPkName(), SQLX_IGNORE_CREATED_AT, SQLX_IGNORE_UPDATED_AT),
			e.getCustomWhere(),
			e.getLimit())
	}
	assert(strSqlx, "update sql is nil")
	return
}

func (e *Engine) makeSqlxInsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strSqlx = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", e.strTableName, strColumns, strValues)
	return
}

func (e *Engine) makeSqlxUpsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strOnConflictUpdates := e.getOnConflictUpdates()
	strSqlx = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v) %v",
		e.strTableName, strColumns, strValues, strOnConflictUpdates)
	return
}
