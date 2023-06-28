package types

import "strings"

const (
	TAG_NAME_DB       = "db"
	TAG_NAME_JSON     = "json"
	TAG_NAME_PROTOBUF = "protobuf"
	TAG_NAME_SQLCA    = "sqlca"
)

const (
	SQLCA_TAG_VALUE_AUTO_INCR = "autoincr" //auto increment
	SQLCA_TAG_VALUE_READ_ONLY = "readonly" //read only
	SQLCA_TAG_VALUE_IS_NULL   = "isnull"   //is nullable
	SQLCA_TAG_VALUE_IGNORE    = "-"        //ignore
)

const (
	PROTOBUF_VALUE_NAME = "name"
	SQLCA_CHAR_ASTERISK = "*"
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
	DATABASE_KEY_NAME_ON         = "ON"
	DATABASE_KEY_NAME_INNER_JOIN = "INNER JOIN"
	DATABASE_KEY_NAME_LEFT_JOIN  = "LEFT JOIN"
	DATABASE_KEY_NAME_RIGHT_JOIN = "RIGHT JOIN"
	DATABASE_KEY_NAME_FULL_JOIN  = "FULL OUTER JOIN" //MSSQL-SERVER
	DATABASE_KEY_NAME_SUM        = "SUM"
	DATABASE_KEY_NAME_AVG        = "AVG"
	DATABASE_KEY_NAME_MIN        = "MIN"
	DATABASE_KEY_NAME_MAX        = "MAX"
	DATABASE_KEY_NAME_COUNT      = "COUNT"
	DATABASE_KEY_NAME_ROUND      = "ROUND"
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
	AdapterSqlx_MySQL    AdapterType = 1  //mysql
	AdapterSqlx_Postgres AdapterType = 2  //postgresql
	AdapterSqlx_Sqlite   AdapterType = 3  //sqlite
	AdapterSqlx_Mssql    AdapterType = 4  //mssql server
	AdapterCache_Redis   AdapterType = 11 //redis
)

func (a AdapterType) GoString() string {
	return a.String()
}

func (a AdapterType) String() string {

	switch a {
	case AdapterSqlx_MySQL:
		return "MySQL"
	case AdapterSqlx_Postgres:
		return "Postgres"
	case AdapterSqlx_Sqlite:
		return "Sqlite"
	case AdapterSqlx_Mssql:
		return "Mssql"
	case AdapterCache_Redis:
		return "Redis"
	}
	return "<Unknown>"
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

func GetAdapterType(name string) AdapterType {
	if strings.EqualFold(name, "sqlite") {
		name = DRIVER_NAME_SQLITE
	}
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
