package sqlca

const (
	TAG_NAME_DB    = "db"
	TAG_NAME_JSON  = "json"
	TAG_NAME_REDIS = "redis"
)

type AdapterType int

const (
	DEFAULT_PRIMARY_KEY_NAME = "id"
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
	OperType_Query  OperType = 1 // query sql
	OperType_Update OperType = 2 // update sql
	OperType_Insert OperType = 3 // insert sql
	OperType_Upsert OperType = 4 // insert or update sql
	OperType_Tx     OperType = 5 // transaction sql
	OperType_Alter  OperType = 6 // alter sql
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
	case OperType_Tx:
		return "OperType_Tx"
	case OperType_Alter:
		return "OperType_Alter"
	}
	return "OperType_Unknown"
}

type ModeType int

const (
	ModeType_ORM = 1
	ModeType_Raw = 2
)

func (m ModeType) GoString() string {
	return m.String()
}

func (m ModeType) String() string {
	switch m {
	case ModeType_ORM:
		return "ModeType_ORM"
	case ModeType_Raw:
		return "ModeType_Raw"
	}
	return "ModeType_Unknown"
}

// clone engine
func (e *Engine) clone(model interface{}) *Engine {

	dict := Struct(model).ToMap(TAG_NAME_DB)
	return &Engine{
		db:        e.db,
		cache:     e.cache,
		debug:     e.debug,
		model:     model,
		dict:      dict,
		strPkName: e.strPkName,
	}
}

func (e *Engine) checkModel() bool {

	if e.model == nil {
		e.panic("orm model is nil, please call Model() method before query or update")
		return false
	}
	return true
}

func (e *Engine) getTableName() string {
	return e.strTableName
}

func (e *Engine) setTableName(strName string) {
	e.strTableName = strName
}

func (e *Engine) getPkName() string {
	return e.strPkName
}

func (e *Engine) setPkName(strName string) {
	e.strPkName = strName
}

func (e *Engine) getWhere() string {
	return e.strWhere
}

func (e *Engine) setWhere(strWhere string) {
	e.strWhere = strWhere
}

func (e *Engine) getModeType() ModeType {
	return e.modeType
}

func (e *Engine) setModeType(modeType ModeType) {
	e.modeType = modeType
}

func (e *Engine) getOperType() OperType {
	return e.operType
}

func (e *Engine) setOperType(operType OperType) {
	e.operType = operType
}

// get data base driver name and data source name
func (e *Engine) getConnUrl(adapterType AdapterType, strUrl string) (strScheme, strDSN string) {
	//TODO @libin connect to database
	strScheme = adapterType.Schema()
	return strScheme, strUrl
}

func (e *Engine) makeOrmQueryMysql() (strSQL string) {
	assert(e.getModeType() == ModeType_ORM, "not a orm mode")
	//TODO: @libin make SQL query string (mysql)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}

func (e *Engine) makeOrmQuerySqlite() (strSQL string) {
	assert(e.getModeType() == ModeType_ORM, "not a orm mode")
	//TODO: @libin make SQL query string (sqlite)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}

func (e *Engine) makeOrmQueryPostgresql() (strSQL string) {
	assert(e.getModeType() == ModeType_ORM, "not a orm mode")
	//TODO: @libin make SQL query string (postgresql)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}

func (e *Engine) makeOrmQueryMssql() (strSQL string) {
	assert(e.getModeType() == ModeType_ORM, "not a orm mode")
	//TODO: @libin make SQL query string (mssql)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}
