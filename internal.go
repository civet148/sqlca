package sqlca

import (
	"fmt"
	"github.com/civet148/log"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (e *Engine) appendMaster(db executor) {
	e.dbMasters = append(e.dbMasters, db)
	log.Debugf("db masters [%v]", len(e.dbMasters))
}

func (e *Engine) appendSlave(db executor) {
	e.dbSlaves = append(e.dbSlaves, db)
	log.Debugf("db slaves [%v]", len(e.dbSlaves))
}

// get slave db instance if use Slave() method to query, if not exist return a master db instance
func (e *Engine) getQueryDB() (db executor) {
	if e.slave {
		db = e.getSlave()
		if db != nil {
			return
		}
	}
	return e.getMaster()
}

// get a master db instance
func (e *Engine) getMaster() executor {

	n := len(e.dbMasters)
	if n > 0 {
		return e.dbMasters[rand.Intn(n)]
	}
	return nil
}

// get a slave db instance
func (e *Engine) getSlave() executor {
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
							//log.Warnf("[%+v] -> pointer is nil", typ.Elem().Name())
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
			e.setModelType(ModelType_Struct)
		} else {
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
				if typ.Kind() == reflect.Slice && val.IsNil() {
					modelVal := reflect.ValueOf(e.model)
					elemTyp := modelVal.Type().Elem()
					elemVal := reflect.New(elemTyp).Elem()
					val.Set(reflect.MakeSlice(elemVal.Type(), 0, 0))
				}
			} else {
				e.model = models //built-in types, eg int/string/float32...
			}
		}

		var selectColumns []string
		e.dict = newReflector(e, e.model).ToMap(e.dbTags...)
		for k, _ := range e.dict {
			selectColumns = append(selectColumns, k)
		}
		if len(selectColumns) > 0 {
			e.setSelectColumns(selectColumns...)
		}
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
		adapterType:     e.adapterType,
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

func (e *Engine) txBegin() (txe *Engine, err error) {

	txe = e.clone()
	db := e.getMaster()
	if txe.tx, err = db.txBegin(); err != nil {
		log.Errorf("newTx error [%+v]", err.Error())
		return nil, err
	}
	txe.operType = OperType_Tx
	return
}

func (e *Engine) defaultInsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Insert(e, strSQL)
}

func (e *Engine) postgresQueryInsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Insert(e, strSQL)
}

func (e *Engine) defaultUpsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Upsert(e, strSQL)
}

func (e *Engine) postgresQueryUpsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Upsert(e, strSQL)
}

func (e *Engine) mssqlQueryInsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Insert(e, strSQL)
}

func (e *Engine) mssqlUpsert(strSQL string) (lastInsertId int64, err error) {
	db := e.getMaster()
	return db.Upsert(e, strSQL)
}

func (e *Engine) getDistinct() string {
	return e.strDistinct
}

func (e *Engine) setDistinct() {
	e.strDistinct = DATABASE_KEY_NAME_DISTINCT
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
		strJoins += fmt.Sprintf(" %s %s %s %s ", v.jt.ToKeyWord(), v.strTableName, DATABASE_KEY_NAME_ON, v.strOn)
	}
	return
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

func (e *Engine) getSelectColumns() (strColumns []string) {
	return e.selectColumns
}

func (e *Engine) setAscColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		e.ascColumns = e.appendStrings(e.ascColumns, strColumns...)
	}
}

func (e *Engine) setDescColumns(strColumns ...string) {
	if len(strColumns) > 0 {
		e.descColumns = e.appendStrings(e.descColumns, strColumns...)
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
	v = e.handleSpecialChars(fmt.Sprintf("%v", v))
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
	if len(strColumns) > 0 {
		e.conflictColumns = e.appendStrings(e.conflictColumns, strColumns...)
	}
}

func (e *Engine) getConflictColumns() []string {
	return e.conflictColumns
}

func (e *Engine) getSingleQuote() (strQuote string) {
	switch e.adapterType {
	case AdapterType_MySQL, AdapterType_Sqlite:
		return "'"
	case AdapterType_Postgres:
		return "'"
	case AdapterType_Mssql:
		return "'"
	}
	return "'"
}

func (e *Engine) getForwardQuote() (strQuote string) {
	switch e.adapterType {
	case AdapterType_MySQL, AdapterType_Sqlite:
		return "`"
	case AdapterType_Postgres:
		return "\""
	case AdapterType_Mssql:
		return "["
	}
	return ""
}

func (e *Engine) getBackQuote() (strQuote string) {
	switch e.adapterType {
	case AdapterType_MySQL, AdapterType_Sqlite:
		return "`"
	case AdapterType_Postgres:
		return "\""
	case AdapterType_Mssql:
		return "]"
	}
	return ""
}

func (e *Engine) getOnConflictForwardKey() (strKey string) {
	switch e.adapterType {
	case AdapterType_MySQL, AdapterType_Sqlite:
		return "ON DUPLICATE"
	case AdapterType_Postgres:
		return "ON CONFLICT ("
	case AdapterType_Mssql:
		return ""
	}
	return
}

func (e *Engine) getOnConflictBackKey() (strKey string) {
	switch e.adapterType {
	case AdapterType_MySQL, AdapterType_Sqlite:
		return "KEY UPDATE"
	case AdapterType_Postgres:
		return ") DO UPDATE SET"
	case AdapterType_Mssql:
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

	if isNilOrFalse(e.orderByColumns) && isNilOrFalse(e.ascColumns) && isNilOrFalse(e.descColumns) {
		return
	}
	return fmt.Sprintf("%v %v", DATABASE_KEY_NAME_ORDER_BY, e.getAscAndDesc())
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

	if e.adapterType != AdapterType_Postgres {
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
		return SQLCA_CHAR_ASTERISK
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
	case AdapterType_MySQL:
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
	case AdapterType_Postgres:
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

	strIn = strings.TrimSpace(strIn) //trim blank characters
	switch e.adapterType {
	case AdapterType_MySQL:
		strIn = strings.Replace(strIn, `\`, `\\`, -1)
		strIn = strings.Replace(strIn, `'`, `\'`, -1)
		strIn = strings.Replace(strIn, `"`, `\"`, -1)
	case AdapterType_Postgres:
		strIn = strings.Replace(strIn, `'`, `''`, -1)
	case AdapterType_Mssql:
		strIn = strings.Replace(strIn, `'`, `''`, -1)
	case AdapterType_Sqlite:
	case AdapterCache_Redis:
	}

	return strIn
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
			val = convertBool2Int(val)
			strVal := fmt.Sprintf("%v", val)
			c := fmt.Sprintf("%v=%v", e.getQuoteColumnName(v), e.getQuoteColumnValue(strVal)) // column name format to `date`='1583055138',...
			cols = append(cols, c)
		}
	}

	if len(cols) == 0 {
		//may be model is a base type slice
		args, ok := e.model.([]interface{})
		if !ok {
			return
		}
		count := len(args)
		//log.Debugf("args count [%v] values [%+v]", count, args)
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
	case AdapterType_MySQL:
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
	case AdapterType_Postgres:
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
	case AdapterType_Mssql:
		{
			strDo = e.getQuoteUpdates(e.getSelectColumns(), e.strPkName)
		}
	case AdapterType_Sqlite:
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

			if e.isReadOnly(k) || e.isExcluded(k) {
				continue
			}
			//log.Warnf("dict key %+v value %#v", k, v)
			v = convertBool2Int(v)
			c := fmt.Sprintf("%v%v%v", e.getForwardQuote(), k, e.getBackQuote()) // column name format to `id`,...

			if k == e.GetPkName() && e.isPkValueNil() {
				continue
			}
			v = e.handleSpecialChars(fmt.Sprintf("%v", v))
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

func (e *Engine) getCaller(skip int) (strFunc string) {
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		n := runtime.FuncForPC(pc).Name()
		ns := strings.Split(n, ".")
		strFunc = ns[len(ns)-1]
	}
	return
}

func (e *Engine) makeSQL() (strSql string) {

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
	if !e.noVerbose {
		log.Debugf("[%v] SQL [%s]", e.getCaller(3), strSql)
	}
	return
}

func (e *Engine) makeInCondition(cond condition) (strCondition string) {

	var strValues []string
	for _, v := range cond.ColumnValues {

		var typ = reflect.TypeOf(v)
		var val = reflect.ValueOf(v)
		switch typ.Kind() {
		case reflect.Slice:
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

	if !e.isPkValueNil() {
		strWhere += e.getPkWhere()
	}

	if strWhere == "" {
		strCustomer := e.getCustomWhere()
		if strCustomer == "" {
			//where condition required when update or delete
			if e.operType != OperType_Update && e.operType != OperType_Delete && len(e.joins) == 0 {
				strWhere += "1=1"
			} else {
				if len(e.joins) > 0 {
					strWhere += "1=1"
				}
			}
		} else {
			strWhere += strCustomer
		}
	}

	//AND conditions
	for _, v := range e.andConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, v)
	}
	//IN conditions
	for _, v := range e.inConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, e.makeInCondition(v))
	}
	//NOT IN conditions
	for _, v := range e.notConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_AND, e.makeNotCondition(v))
	}
	//OR conditions
	for _, v := range e.orConditions {
		strWhere += fmt.Sprintf(" %v %v ", DATABASE_KEY_NAME_OR, v)
	}

	if strWhere != "" {
		strWhere = DATABASE_KEY_NAME_WHERE + " " + strWhere
	} else {
		strWhere = DATABASE_KEY_NAME_WHERE
	}
	return
}

func (e *Engine) makeSqlxQuery() (strSqlx string) {
	strWhere := e.makeWhereCondition()

	switch e.adapterType {
	case AdapterType_Mssql:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getLimit(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy())
	default:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy(), e.getLimit(), e.getOffset())
	}

	return
}

func (e *Engine) makeSqlxQueryCount() (strSqlx string) {
	strWhere := e.makeWhereCondition()

	switch e.adapterType {
	case AdapterType_Mssql:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy())
	default:
		strSqlx = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v",
			DATABASE_KEY_NAME_SELECT, e.getDistinct(), e.getRawColumns(), DATABASE_KEY_NAME_FROM, e.getTableName(), e.getJoins(),
			strWhere, e.getGroupBy(), e.getHaving(), e.getOrderBy(), e.getOffset())
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
		panic("no condition to delete records") //删除必须加条件,WHERE条件可设置为1=1(确保不是人为疏忽)
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
}

func (e *Engine) autoRollback() bool {
	if e.bAutoRollback && e.operType == OperType_Tx && e.tx != nil {
		return true
	}
	return false
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
	return fmt.Sprintf("%s(%s, %d) AS %s", DATABASE_KEY_NAME_ROUND, strColumn, round, strAlias)
}
