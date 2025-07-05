package sqlca

import (
	"database/sql"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/types"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type tableIndex struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (e *Engine) setDB(db *sqlx.DB) {
	e.db = db
}

// get db instance for query
func (e *Engine) getDB() (db *sqlx.DB) {
	return e.db
}

func (e *Engine) setModel(models ...any) *Engine {
	var strCamelTableName string
	for _, v := range models {

		if v == nil {
			continue
		}

		var isStructPtrPtr bool
		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)
		if typ.Kind() == reflect.Ptr {

			typ = typ.Elem()
			val = val.Elem()
			switch typ.Kind() {
			case reflect.Ptr:
				{
					if typ.Elem().Kind() == reflect.Struct { //struct pointer address (&*StructType)
						if val.IsNil() {
							var typNew = typ.Elem()
							var valNew = reflect.New(typNew)
							val.Set(valNew)
						}
						isStructPtrPtr = true
					}
				}
			}
		}
		if isStructPtrPtr {
			e.model = val.Interface()
			e.setModelType(types.ModelType_Struct)
			typSt := typ.Elem()
			strCamelTableName = typSt.Name()
		} else {
			switch typ.Kind() {
			case reflect.Struct: // struct
				e.setModelType(types.ModelType_Struct)
			case reflect.Slice: //  slice
				e.setModelType(types.ModelType_Slice)
			case reflect.Map: // map
				e.setModelType(types.ModelType_Map)
			default: //base type
				e.setModelType(types.ModelType_BaseType)
			}
			if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
				e.model = models[0] //map, struct or slice
				if typ.Kind() == reflect.Slice {
					modelVal := reflect.ValueOf(e.model)
					elemTyp := modelVal.Type().Elem()
					elemVal := reflect.New(elemTyp).Elem()
					typSt := elemVal.Type().Elem()
					if typSt.Kind() == reflect.Ptr {
						typSt = typSt.Elem()
					}
					valSt := reflect.New(typSt)
					if tabler, ok := valSt.Interface().(types.Tabler); ok {
						e.setTableName(tabler.TableName())
					} else {
						strCamelTableName = typSt.Name()
					}
					if val.IsNil() {
						val.Set(reflect.MakeSlice(elemVal.Type(), 0, 0))
					}
				} else {
					var typSt = typ
					if typ.Kind() == reflect.Ptr {
						typSt = typ.Elem()
					}
					valSt := reflect.New(typSt)
					if tabler, ok := valSt.Interface().(types.Tabler); ok {
						e.setTableName(tabler.TableName())
					} else {
						strCamelTableName = typSt.Name()
					}
				}
			} else {
				e.model = models //built-in types
			}
		}
		if strCamelTableName != "" {
			name := convertCamelToSnake(strCamelTableName)
			e.setTableName(strings.ToLower(name))
		}
		var selectColumns []string
		ref := newReflector(e, e.model)
		ref = ref.ParseModel(e.dbTags...)
		e.dict = ref.Dict
		for _, col := range ref.Columns {
			selectColumns = append(selectColumns, col)
		}
		if len(selectColumns) > 0 {
			e.setSelectColumns(selectColumns...)
		}
		break //only check first argument
	}
	return e.setHooks()
}

// clone engine
func (e *Engine) clone(models ...any) *Engine {

	engine := &Engine{
		strDSN:          e.strDSN,
		dsn:             e.dsn,
		db:              e.db,
		adapterType:     e.adapterType,
		adapterCache:    e.adapterCache,
		strPkName:       e.strPkName,
		expireTime:      e.expireTime,
		strDatabaseName: e.strDatabaseName,
		dbTags:          e.dbTags,
		bForce:          e.bForce,
		noVerbose:       e.noVerbose,
		bAutoRollback:   e.bAutoRollback,
		slowQueryOn:     e.slowQueryOn,
		slowQueryTime:   e.slowQueryTime,
		tx:              e.tx,
		operType:        e.operType,
		idgen:           e.idgen,
		options:         e.options,
	}
	engine.setModel(models...)
	return engine
}

func (e *Engine) genSnowflakeId() ID {
	if e.idgen == nil {
		return 0
	}
	return e.idgen.Generate()
}

func (e *Engine) newTx() (txEngine *Engine, err error) {
	txEngine = e.clone()
	db := e.getDB()
	if txEngine.tx, err = db.Beginx(); err != nil {
		log.Errorf("newTx error [%+v]", err.Error())
		return nil, err
	}
	txEngine.operType = types.OperType_Tx
	return
}

func (e *Engine) execQuery() (rowsAffected int64, err error) {
	var query string
	var args []any
	var db = e.getDB()
	var queryer sqlx.Queryer
	_ = queryer
	query, args = e.makeSqlxQuery(false)
	if e.operType == types.OperType_Tx {
		queryer = sqlx.Queryer(e.tx)
	} else {
		queryer = sqlx.Queryer(db)
	}
	var rows *sql.Rows

	//log.Debugf("query [%v] args %v", query, args)
	rows, err = queryer.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rowsAffected, err = e.fetchRows(rows)
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (e *Engine) execQueryEx(strCountSql string) (rowsAffected, total int64, err error) {
	var query string
	var args []any
	var db = e.getDB()
	var queryer sqlx.Queryer
	_ = queryer
	query, args = e.makeSqlxQuery(false)
	if e.operType == types.OperType_Tx {
		queryer = sqlx.Queryer(e.tx)
	} else {
		queryer = sqlx.Queryer(db)
	}
	var rows *sql.Rows

	//log.Debugf("query [%v] args %v", query, args)
	rows, err = queryer.Query(query, args...)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	rowsAffected, err = e.fetchRows(rows)
	if err != nil {
		return 0, 0, err
	}
	var rowsCount *sql.Rows
	if rowsCount, err = queryer.Query(strCountSql); err != nil {
		return 0, 0, err
	}

	defer rowsCount.Close()
	for rowsCount.Next() {
		total++
	}
	return rowsAffected, total, nil
}

func (e *Engine) txQuery(dest interface{}, strSql string, args ...any) (count int64, err error) {
	var rows *sql.Rows
	strSql = e.buildSqlExprs(strSql, args...).RawSQL(e.GetAdapter())
	c := e.Counter()
	defer c.Stop(fmt.Sprintf("query tx [%s]", strSql))

	rows, err = e.tx.Query(strSql)
	if err != nil {
		if !e.noVerbose {
			log.Errorf("query tx sql [%v] args %v query error [%v] auto rollback [%v]", strSql, args, err.Error(), e.bAutoRollback)
		}
		e.autoRollback()
		return
	}
	e.setModel(dest)

	defer rows.Close()
	if count, err = e.fetchRows(rows); err != nil {
		if !e.noVerbose {
			log.Errorf("query tx sql [%v] args %v fetch row error [%v] auto rollback [%v]", strSql, args, err.Error(), e.bAutoRollback)
		}
		e.autoRollback()
		return
	}
	return
}

func (e *Engine) txRollback() error {
	return e.tx.Rollback()
}

func (e *Engine) txCommit() error {
	return e.tx.Commit()
}

func (e *Engine) postgresQueryInsert(strSQL string) string {
	strSQL += fmt.Sprintf(" RETURNING \"%v\"", e.GetPkName())
	return strSQL
}

func (e *Engine) mysqlQueryUpsert(strSQL string) (lastInsertId int64, err error) {

	if e.operType == types.OperType_Tx {
		lastInsertId, _, err = e.TxExec(strSQL)
		if err != nil {
			return 0, log.Errorf("upsert error [%s]", err)
		}
	} else {
		var r sql.Result
		db := e.getDB()
		r, err = db.Exec(strSQL)
		if err != nil {
			if !e.noVerbose {
				log.Errorf("error %v model %+v", err, e.model)
			}
			return
		}
		lastInsertId, err = r.LastInsertId()
		if err != nil {
			if !e.noVerbose {
				log.Errorf("get last insert id error %v model %+v", err, e.model)
			}
			return
		}
	}
	return
}

func (e *Engine) postgresQueryUpsert(strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	log.Debugf("[%v]", strSQL)
	db := e.getDB()
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

func (e *Engine) mysqlExec(strSQL string) (lastInsertId, rowsAffected int64, err error) {
	var r sql.Result
	var db *sqlx.DB
	db = e.getDB()
	r, err = db.Exec(strSQL)
	if err != nil {
		return 0, 0, err
	}
	rowsAffected, _ = r.RowsAffected()
	lastInsertId, _ = r.LastInsertId()
	return
}

func (e *Engine) mssqlQueryInsert(strSQL string) string {
	strSQL += " SELECT SCOPE_IDENTITY() AS last_insert_id"
	log.Debugf("[%v]", strSQL)
	return strSQL
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
		strSQL = e.mssqlQueryInsert(strSQL)
		if lastInsertId, _, err = db.TxExec(strSQL); err != nil {
			log.Errorf("mssqlQueryInsert [%v] error [%v]", strSQL, err.Error())
			_ = db.TxRollback()
			return
		}
	} else {
		// UPDATE users SET xxx=yyy WHERE id=nnn
		strUpdates := fmt.Sprintf("%v %v %v %v %v %v=%v",
			types.DATABASE_KEY_NAME_UPDATE, e.getTableName(),
			types.DATABASE_KEY_NAME_SET, e.getOnConflictDo(),
			types.DATABASE_KEY_NAME_WHERE, e.GetPkName(), lastInsertId)
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
	e.strDistinct = types.DATABASE_KEY_NAME_DISTINCT
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
	return e.dict[strCol]
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

func (e *Engine) getJoins() (strJoins string) {
	for _, v := range e.joins {
		strJoins += fmt.Sprintf(" %s %s %s %s ", v.jt.ToKeyWord(), v.strTableName, types.DATABASE_KEY_NAME_ON, v.strOn)
	}
	return
}

func (e *Engine) setPkValue(value any) {

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
		}
	default:
		strValue = fmt.Sprintf("%v", value)
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

func (e *Engine) exist(src []string, s string) bool {
	for _, v := range src {
		if s == v {
			return true
		}
	}
	return false
}

func (e *Engine) appendStrings(src []string, dest ...string) []string {
	//check duplicated elements
	for _, v := range dest {
		if !e.exist(src, v) {
			src = append(src, v)
		}
	}
	return src
}

func (e *Engine) setSelectColumnsForce(strColumns ...string) {
	if len(strColumns) > 0 {
		e.selectColumns = strColumns
	}
}

func (e *Engine) setSelectColumns(strColumns ...string) (ok bool) {
	if len(strColumns) == 0 {
		return false
	}
	if e.selected {
		e.selectColumns = e.appendStrings(e.selectColumns, strColumns...)
	} else {
		e.selectColumns = strColumns
	}
	return true
}

func (e *Engine) setExcludeColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		e.excludeColumns = e.appendStrings(e.excludeColumns, strColumns...)
	}
}

func (e *Engine) setNullableColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		e.nullableColumns = e.appendStrings(e.nullableColumns, strColumns...)
	}
}

func (e *Engine) getSelectColumns() (strColumns []string) {
	return e.selectColumns
}

func (e *Engine) setAscColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		var orders []string
		for _, col := range strColumns {
			orders = append(orders, fmt.Sprintf("%s ASC", col))
		}
		e.orderByColumns = append(e.orderByColumns, strings.Join(orders, ","))
	}
}

func (e *Engine) setDescColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		var orders []string
		for _, col := range strColumns {
			orders = append(orders, fmt.Sprintf("%s DESC", col))
		}
		e.orderByColumns = append(e.orderByColumns, strings.Join(orders, ","))
	}
}

func (e *Engine) setCustomizeUpdates(strUpdates ...string) {
	if len(strUpdates) > 0 {
		e.strUpdates = strUpdates
	}
}

func (e *Engine) getCustomizeUpdates() []string {
	return e.strUpdates
}

// SELECT ... FROM xxx ORDER BY c1, c2 ASC, c3 DESC
func (e *Engine) getAscAndDesc() (strAscDesc string) {

	var ss []string
	//make default order by expression
	if len(ss) == 0 {
		ss = append(ss, strings.Join(e.orderByColumns, ","))
	}
	return strings.Join(ss, ",")
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

func (e *Engine) getQuoteColumnValue(v any) (strValue string) {
	return fmt.Sprintf("%v", quotedValue(v))
}

func (e *Engine) setWhere(query string, args ...any) {
	e.andConditions = append(e.andConditions, e.buildSqlExprs(query, args...))
}

func (e *Engine) getModelType() types.ModelType {
	return e.modelType
}

func (e *Engine) setModelType(modelType types.ModelType) {
	e.modelType = modelType
}

func (e *Engine) getOperType() types.OperType {
	return e.operType
}

func (e *Engine) setOperType(operType types.OperType) {
	e.operType = operType
}

func (e *Engine) setConflictColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		e.conflictColumns = e.appendStrings(e.conflictColumns, strColumns...)
	}
}

func (e *Engine) getConflictColumns() []string {
	return e.conflictColumns
}

func (e *Engine) getSingleQuote() (strQuote string) {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL, types.AdapterSqlx_Sqlite:
		return "'"
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return "'"
	case types.AdapterSqlx_Mssql:
		return "'"
	}
	return "'"
}

func (e *Engine) getForwardQuote() (strQuote string) {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL, types.AdapterSqlx_Sqlite:
		return "`"
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return "\""
	case types.AdapterSqlx_Mssql:
		return "["
	}
	return ""
}

func (e *Engine) getBackQuote() (strQuote string) {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL, types.AdapterSqlx_Sqlite:
		return "`"
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return "\""
	case types.AdapterSqlx_Mssql:
		return "]"
	}
	return ""
}

func (e *Engine) getOnConflictForwardKey() (strKey string) {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL, types.AdapterSqlx_Sqlite:
		return "ON DUPLICATE"
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return "ON CONFLICT ("
	case types.AdapterSqlx_Mssql:
		return ""
	}
	return
}

func (e *Engine) getOnConflictBackKey() (strKey string) {
	switch e.adapterType {
	case types.AdapterSqlx_MySQL, types.AdapterSqlx_Sqlite:
		return "KEY UPDATE"
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return ") DO UPDATE SET"
	case types.AdapterSqlx_Mssql:
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
	if len(strColumns) > 0 {
		e.orderByColumns = e.appendStrings(e.orderByColumns, strColumns...)
	}
}

func (e *Engine) getOrderBy() (strOrderBy string) {
	if isNilOrFalse(e.orderByColumns) {
		return
	}
	return fmt.Sprintf("%v %v", types.DATABASE_KEY_NAME_ORDER_BY, e.getAscAndDesc())
}

func (e *Engine) setGroupBy(strColumns ...string) {
	if len(strColumns) > 0 {
		e.groupByColumns = e.appendStrings(e.groupByColumns, strColumns...)
	}
}

func (e *Engine) setHaving(havingCondition string) {
	e.havingCondition = havingCondition
}

func (e *Engine) getHaving() (strHaving string) {

	if isNilOrFalse(e.havingCondition) {
		return
	}
	return fmt.Sprintf("%v %v", types.DATABASE_KEY_NAME_HAVING, e.havingCondition)
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

func (e *Engine) isEmpty(strCol string) bool {

	for _, v := range e.nullableColumns {
		if v == strCol {
			val := fmt.Sprintf("%v", e.dict[strCol])
			if val == "" {
				return true
			}
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

	if e.adapterType != types.AdapterSqlx_Postgres && e.adapterType != types.AdapterSqlx_OpenGauss {
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
func (e *Engine) getCountColumn() string {
	return "COUNT(*)"
}

func (e *Engine) getRawColumns() (strColumns string) {
	var selectCols []string

	if len(e.selectColumns) == 0 {
		return types.SQLCA_CHAR_ASTERISK
	}

	for _, v := range e.selectColumns {
		if e.isColumnSelected(v) {
			selectCols = append(selectCols, v)
		}
	}
	selectCols = e.makeNearbyColumn(selectCols...)
	if len(selectCols) > 0 {
		if e.strCaseWhen != "" {
			selectCols = append(selectCols, e.strCaseWhen)
		}
		strColumns = strings.Join(selectCols, ",")
	}
	return
}

func (e *Engine) trimNearbySameColumn(strAS string, strColumns ...string) (columns []string) {
	for _, v := range strColumns {
		if v != strAS {
			columns = append(columns, v)
		}
	}
	return
}

func (e *Engine) makeNearbyColumn(strColumns ...string) (columns []string) {

	columns = strColumns
	switch e.adapterType {
	case types.AdapterSqlx_MySQL:
		{
			/* -- MySQL
			SELECT  id,lng,lat,name,(6371 * ACOS(COS(RADIANS(lat)) * COS(RADIANS(28.803909723)) * COS(RADIANS(121.5619236231) - RADIANS(lng))
			+ SIN(RADIANS(lat)) * SIN(RADIANS(28.803909723)))) AS distance FROM t_address WHERE 1=1 HAVING distance <= 113
			*/
			//NEARBY additional column
			if e.nearby != nil {
				nb := e.nearby
				strNearBy := fmt.Sprintf(`(6371 * ACOS(COS(RADIANS(%v)) * COS(RADIANS(%v)) * COS(RADIANS(%v) - RADIANS(%v)) + SIN(RADIANS(%v)) * SIN(RADIANS(%v)))) AS %s`,
					nb.strLatCol, nb.lat, nb.lng, nb.strLngCol, nb.strLatCol, nb.lat, nb.strAS)
				columns = e.trimNearbySameColumn(nb.strAS, columns...)
				columns = append(columns, strNearBy)
				e.setHaving(fmt.Sprintf("%s <= %v", nb.strAS, nb.distance))
			}
		}
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		{
			/* -- Postgres
			SELECT  a.* FROM
			 (
			   SELECT id,lng,lat,name,(6371 * ACOS(COS(RADIANS(lat)) * COS(RADIANS(28.803909723)) * COS(RADIANS(121.5619236231) - RADIANS(lng))
			   + SIN(RADIANS(lat)) * SIN(RADIANS(28.803909723)))) AS distance FROM t_address WHERE 1=1
			) a
			WHERE a.distance <= 113
			*/
		}
	}
	return
}

// handle special characters, prevent SQL inject
func (e *Engine) handleSpecialChars(strIn string) (strOut string) {
	return types.PreventSqlInject(e.GetAdapter(), strIn)
}

func (e *Engine) getQuoteUpdates(strColumns []string, strExcepts ...string) (strUpdates string) {

	var cols []string
	for _, v := range strColumns {

		if e.isColumnSelected(v, strExcepts...) && !e.isReadOnly(v) {
			val := e.getModelValue(v)
			if val == nil {
				val = types.SqlNull{}
			}
			val = convertBool2Int(val)
			c := fmt.Sprintf("%v=%v", e.getQuoteColumnName(v), e.getQuoteColumnValue(val)) // column name format to `date`='1583055138',...
			cols = append(cols, c)
		}
	}

	if len(cols) == 0 {
		//may be model is a base type slice
		args, ok := e.model.([]any)
		if !ok {
			return
		}
		count := len(args)
		for i, k := range strColumns {
			if i < count {
				v := args[i]
				typ := reflect.TypeOf(v)
				val := reflect.ValueOf(v)
				var value interface{}
				kind := typ.Kind()
				if kind != reflect.Interface && kind != reflect.Ptr {
					value = val.Interface()
				} else {
					value = val.Elem().Interface()
				}
				c := fmt.Sprintf("%v=%v", e.getQuoteColumnName(k), e.getQuoteColumnValue(value))
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
	var strUpdates string
	var strCustomizeUpdates = e.getCustomizeUpdates()
	switch e.adapterType {
	case types.AdapterSqlx_MySQL:
		{
			if len(strCustomizeUpdates) != 0 {
				strUpdates = strings.Join(strCustomizeUpdates, ",")
			} else {
				strUpdates = e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
			}

			if !isNilOrFalse(strUpdates) {
				if e.isPkInteger() { // primary key type is a integer
					strDo = fmt.Sprintf("%v", strUpdates)
				} else {
					strDo = strUpdates
				}
			}
		}
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		{
			if len(strCustomizeUpdates) != 0 {
				strUpdates = strings.Join(strCustomizeUpdates, ",")
			} else {
				strUpdates = e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
			}
			if !isNilOrFalse(strUpdates) {
				strDo = fmt.Sprintf("%v RETURNING \"%v\"", strUpdates, e.GetPkName()) // TODO @libin test postgresql ON CONFLICT(...) DO UPDATE SET ... RETURNING id
			}
		}
	case types.AdapterSqlx_Mssql:
		{
			strDo = e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
		}
	case types.AdapterSqlx_Sqlite:
		{
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

	if typ.Kind() == reflect.Slice {
		var cols2 []string
		var values [][]string
		var valueQuoteSlice []string
		var excludeIndexes []int
		cols2, values = e.getStructSliceKeyValues(true)
		for i, v := range cols2 {

			if e.isReadOnly(v) || e.isExcluded(v) || e.isEmpty(v) {
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
				vq := e.getQuoteColumnValue(vv)
				valueQoute = append(valueQoute, vq)
			}
			valueQuoteSlice = append(valueQuoteSlice, fmt.Sprintf("(%v)", strings.Join(valueQoute, ",")))
		}

		if len(valueQuoteSlice) > 0 {
			strColonValues = strings.Join(valueQuoteSlice, ",")
		}
	} else {
		for k, v := range e.dict {
			if e.isReadOnly(k) || e.isExcluded(k) || e.isEmpty(k) {
				continue
			}
			v = convertBool2Int(v)
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

func (e *Engine) buildSqlExprs(query string, args ...any) types.Expr {
	var keepSlice bool
	if !strings.Contains(query, "?") && len(args) > 0 {
		query = fmt.Sprintf("%s = ?", query)
	} else {
		if shouldKeepSlice(query, args...) {
			keepSlice = true
		}
	}
	var vars []any
	for _, arg := range args {
		vars = append(vars, indirectValue(arg, keepSlice))
	}
	return types.Expr{SQL: query, Vars: vars}
}

func (e *Engine) makeSqlxQueryPrimaryKey() (strSql string) {

	strSql = fmt.Sprintf("%v %v%v%v %v %v %v %v%v%v=%v%v%v",
		types.DATABASE_KEY_NAME_SELECT, e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		types.DATABASE_KEY_NAME_FROM, e.getTableName(), types.DATABASE_KEY_NAME_WHERE,
		e.getForwardQuote(), e.GetPkName(), e.getBackQuote(),
		e.getSingleQuote(), e.getPkValue(), e.getSingleQuote())
	return
}

func (e *Engine) getCaller(skip int) (strFunc string) {
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		n := runtime.FuncForPC(pc).Name()
		ns := strings.Split(n, ".")
		strFunc = ns[len(ns)-1]
	}
	return
}

func (e *Engine) makeSQL(operType types.OperType, rawSQL bool) (strSql string, args []any) {

	switch operType {
	case types.OperType_Query:
		strSql, args = e.makeSqlxQuery(rawSQL)
	case types.OperType_Update:
		strSql, args = e.makeSqlxUpdate(rawSQL)
	case types.OperType_Insert:
		strSql = e.makeSqlxInsert()
	case types.OperType_Upsert:
		strSql = e.makeSqlxUpsert()
	case types.OperType_Delete:
		strSql, args = e.makeSqlxDelete(rawSQL)
	default:
		log.Errorf("operation illegal")
	}
	return strings.TrimSpace(strSql), args
}

func (e *Engine) makeInCondition(cond types.Expr) (strCondition string, args []any) {

	var strValues []string
	for _, v := range cond.Vars {

		var typ = reflect.TypeOf(v)
		var val = reflect.ValueOf(v)
		switch typ.Kind() {
		case reflect.Slice, reflect.Array:
			{
				n := val.Len()
				for i := 0; i < n; i++ {
					strValues = append(strValues, fmt.Sprintf("%v%v%v", e.getSingleQuote(), val.Index(i).Interface(), e.getSingleQuote()))
				}
			}
		default:
			strValues = append(strValues, fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote()))
		}
	}
	strCondition = fmt.Sprintf("%v %v (%v)", cond.SQL, types.DATABASE_KEY_NAME_IN, strings.Join(strValues, ","))
	return
}

func (e *Engine) makeNotCondition(cond types.Expr) (strCondition string, args []any) {

	var strValues []string
	for _, v := range cond.Vars {
		strValues = append(strValues, fmt.Sprintf("%v%v%v", e.getSingleQuote(), v, e.getSingleQuote()))
	}
	strCondition = fmt.Sprintf("%v %v (%v)", cond.SQL, types.DATABASE_KEY_NAME_NOT_IN, strings.Join(strValues, ","))
	return
}

func (e *Engine) makeWhereCondition(operType types.OperType, rawSQL bool) (strWhere string, args []any) {

	if !e.isPkValueNil() {
		strWhere += e.getPkWhere()
	}

	if strWhere == "" {
		//where condition required when update or delete
		if operType != types.OperType_Update && operType != types.OperType_Delete && len(e.joins) == 0 {
			strWhere += "1=1"
		} else {
			if len(e.joins) > 0 || len(e.andConditions) != 0 {
				strWhere += "1=1"
			}
		}
	}

	//AND conditions
	for _, v := range e.andConditions {
		var query string
		if rawSQL {
			query = v.RawSQL()
		} else {
			query = v.SQL
			if len(v.Vars) != 0 {
				args = append(args, v.Vars...)
			}
		}
		strWhere += fmt.Sprintf(" %v %v ", types.DATABASE_KEY_NAME_AND, query)
	}
	//IN conditions
	for _, v := range e.inConditions {
		var query string
		var vars []any
		query, vars = e.makeInCondition(v)
		if len(vars) != 0 {
			args = append(args, vars...)
		}
		strWhere += fmt.Sprintf(" %v %v ", types.DATABASE_KEY_NAME_AND, query)
	}
	//NOT IN conditions
	for _, v := range e.notConditions {
		var query string
		var vars []any
		query, vars = e.makeNotCondition(v)
		if len(vars) != 0 {
			args = append(args, vars...)
		}
		strWhere += fmt.Sprintf(" %v %v ", types.DATABASE_KEY_NAME_AND, query)
	}
	//OR conditions
	for _, v := range e.orConditions {
		if strings.Contains(v.SQL, "(") && strings.Contains(v.SQL, ")") {
			strWhere += fmt.Sprintf(" %v %v ", types.DATABASE_KEY_NAME_AND, v.RawSQL()) //multiple OR condition append
		} else {
			strWhere += fmt.Sprintf(" %v %v ", types.DATABASE_KEY_NAME_OR, v.RawSQL()) //single OR condition append
		}
	}

	if strWhere != "" {
		strWhere = types.DATABASE_KEY_NAME_WHERE + " " + strWhere
	} else {
		strWhere = types.DATABASE_KEY_NAME_WHERE
	}
	return
}

func (e *Engine) makeSqlxQuery(rawSQL bool) (strSqlx string, args []any) {
	var strWhere string
	strWhere, args = e.makeWhereCondition(types.OperType_Query, rawSQL)

	switch e.adapterType {
	case types.AdapterSqlx_Mssql:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v",
			types.DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getLimit(), e.getRawColumns(), types.DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy())
	default:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v",
			types.DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), types.DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy(), e.getLimit(), e.getOffset())
	}
	return
}

func (e *Engine) makeSqlxQueryCount() (strSqlx string) {
	strWhere, _ := e.makeWhereCondition(types.OperType_Query, true)

	switch e.adapterType {
	case types.AdapterSqlx_Mssql:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v",
			types.DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), types.DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy())
	default:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v",
			types.DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), types.DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy(), e.getOffset())
	}
	return
}

func (e *Engine) makeSqlxForUpdate(rawSQL bool) (strSql string, args []any) {
	strSql, args = e.makeSqlxQuery(rawSQL)
	strSql += " " + types.DATABASE_KEY_NAME_FOR_UPDATE
	return strSql, args
}

func (e *Engine) makeSqlxUpdate(rawSQL bool) (strSqlx string, args []any) {
	var strWhere string
	strWhere, args = e.makeWhereCondition(types.OperType_Update, rawSQL)
	strSqlx = fmt.Sprintf("%v %v %v %v %v %v",
		types.DATABASE_KEY_NAME_UPDATE, e.getTableName(), types.DATABASE_KEY_NAME_SET,
		e.getQuoteUpdates(e.getSelectColumns(), e.GetPkName()), strWhere, e.getLimit())
	assert(strSqlx, "update sql is nil")
	return
}

func (e *Engine) makeSqlxInsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strSqlx = fmt.Sprintf("%v %v %v %v %v", types.DATABASE_KEY_NAME_INSERT, e.getTableName(), strColumns, types.DATABASE_KEY_NAME_VALUES, strValues)
	return
}

func (e *Engine) makeSqlxUpsert() (strSqlx string) {

	strColumns, strValues := e.getInsertColumnsAndValues()
	strOnConflictUpdates := e.getOnConflictUpdates()
	strSqlx = fmt.Sprintf("%v %v %v %v %v %v", types.DATABASE_KEY_NAME_INSERT, e.getTableName(), strColumns, types.DATABASE_KEY_NAME_VALUES, strValues, strOnConflictUpdates)
	return
}

func (e *Engine) makeSqlxDelete(rawSQL bool) (strSqlx string, args []any) {
	var strWhere string
	strWhere, args = e.makeWhereCondition(types.OperType_Delete, rawSQL)
	if strWhere == "" {
		panic("no condition to delete records") //删除必须加条件,WHERE条件可设置为1=1(确保不是人为疏忽)
	}
	strSqlx = fmt.Sprintf("%v %v %v %v %v", types.DATABASE_KEY_NAME_DELETE, types.DATABASE_KEY_NAME_FROM, e.getTableName(), strWhere, e.getLimit())
	return
}

func (e *Engine) cleanWhereCondition() {
	e.strPkValue = ""
	e.cleanHooks()
}

func (e *Engine) autoRollback() {
	if e.bAutoRollback && e.operType == types.OperType_Tx && e.tx != nil {
		_ = e.tx.Rollback()
		log.Debugf("tx auto rollback successful")
	}
}

func (e *Engine) aggFunc(strKey, strColumn string, strAS ...string) string {
	var strAlias string
	if len(strAS) == 0 {
		strAlias = strColumn
	} else {
		strAlias = strAS[0]
	}
	return fmt.Sprintf("%s(%s) AS %s", strKey, strColumn, strAlias)
}

func (e *Engine) roundFunc(strColumn string, round int, strAS ...string) string {
	var strAlias string
	if len(strAS) == 0 {
		strAlias = strColumn
	} else {
		strAlias = strAS[0]
	}
	return fmt.Sprintf("%s(%s, %d) AS %s", types.DATABASE_KEY_NAME_ROUND, strColumn, round, strAlias)
}

func (e *Engine) setForUpdate() *Engine {
	e.strForUpdate = types.DATABASE_KEY_NAME_FOR_UPDATE
	return e
}

func (e *Engine) setLockShareMode() *Engine {
	e.strLockShareMode = types.DATABASE_KEY_NAME_LOCK_SHARE_MODE
	return e
}

func (e *Engine) setAnd(query string, args ...any) *Engine {
	assert(query, "query statement is empty")
	expr := e.buildSqlExprs(query, args...)
	e.andConditions = append(e.andConditions, expr)
	return e
}

func (e *Engine) setOr(exprs ...types.Expr) *Engine {
	if len(exprs) == 1 {
		e.orConditions = append(e.orConditions, exprs[0])
		return e
	}

	var ors []string
	for _, expr := range exprs {
		ors = append(ors, expr.RawSQL(e.GetAdapter()))
	}
	strCombOrs := " ( " + strings.Join(ors, " OR ") + " ) "
	e.orConditions = append(e.orConditions, types.Expr{SQL: strCombOrs})
	return e
}

func (e *Engine) parseQueryAndMap(query any) {
	where := query.(map[string]any)
	for k, v := range where {
		if strings.Contains(k, "?") {
			e.setAnd(k, v)
		} else {
			cond := fmt.Sprintf("%s = ?", k)
			e.setAnd(cond, v)
		}
	}
}

func (e *Engine) parseQueryOrMap(query any) {
	var qss []types.Expr
	where := query.(map[string]any)
	for k, v := range where {
		if strings.Contains(k, "?") {
			qss = append(qss, types.Expr{
				SQL:  k,
				Vars: []any{v},
			})
		} else {
			k = fmt.Sprintf("%s = ?", k)
			qss = append(qss, types.Expr{
				SQL:  k,
				Vars: []any{v},
			})
		}
	}
	e.setOr(qss...)
}

func (e *Engine) setNormalCondition(query any, args ...any) *Engine {
	var strSql string
	qt := parseQueryInterface(query)
	switch qt {
	case queryInterface_String:
		strSql = query.(string)
		expr := e.buildSqlExprs(strSql, args...)
		e.andConditions = append(e.andConditions, expr)
	case queryInterface_Map:
		e.parseQueryAndMap(query)
	}
	return e
}

func (e *Engine) verbose(msg string, args ...any) {
	var isError bool
	for _, arg := range args {
		if _, ok := arg.(error); ok {
			isError = true
		}
	}
	if !e.noVerbose {
		if isError {
			log.Errorf(msg, args...)
		} else {
			log.Debugf(msg, args...)
		}
	}
}
