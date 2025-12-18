package sqlca

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/types"
	"reflect"
	"strconv"
	"strings"
)

type ModelReflector struct {
	value   interface{}            //value
	engine  *Engine                //database engine
	Dict    map[string]interface{} //dictionary of structure tag and value
	Columns []string               //column names
}

type Fetcher struct {
	count     int               //column count in db table
	cols      []string          //column names in db table
	types     []*sql.ColumnType //column types in db table
	arrValues [][]byte          //value slice
	mapValues map[string]string //value map
	arrIndex  int               //fetch index
}

func newReflector(e *Engine, v interface{}) *ModelReflector {

	return &ModelReflector{
		value:  v,
		engine: e,
		Dict:   make(map[string]interface{}),
	}
}

// map[string]string to [][]byte
func mapToBytesSlice(m map[string]string) (arrays [][]byte) {
	for _, v := range m {
		arr := []byte(v)
		arrays = append(arrays, arr)
	}
	return
}

// parse struct tag and value to map
func (s *ModelReflector) ParseModel(tagNames ...string) *ModelReflector {

	typ := reflect.TypeOf(s.value)
	val := reflect.ValueOf(s.value)

	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	switch kind {
	case reflect.Struct:
		{
			s.parseStructFields(typ, val, tagNames...)
		}
	case reflect.Slice:
		{
			typ = val.Type().Elem()
			val = reflect.New(typ).Elem()
			s.parseStructFields(typ, val, tagNames...)
		}
	case reflect.Map:
		{
			var dict map[string]any
			if v, ok := s.value.(*map[string]any); ok {
				dict = *v
			}
			if v, ok := s.value.(map[string]any); ok {
				dict = v
			}
			for k, v := range dict {
				s.Columns = append(s.Columns, k)
				_ = v
			}
			s.Dict = dict
		}
	default:
		log.Warnf("kind [%v] not support yet", typ.Kind())
	}
	if len(s.Columns) == 0 {
		s.Columns = []string{"*"}
	}
	return s
}

func (s *ModelReflector) convertMapString(ms map[string]string) (mi map[string]any) {
	mi = make(map[string]any, 10)
	for k, v := range ms {
		mi[k] = v
	}
	return
}

// get struct field's tag value
func (s *ModelReflector) getTag(sf reflect.StructField, tagNames ...string) (strValue string, ignore bool) {
	for _, tagName := range tagNames {
		strValue = handleTagValue(tagName, sf.Tag.Get(tagName))
		if strValue == types.SQLCA_TAG_VALUE_IGNORE {
			return "", true
		}
		if strValue != "" {
			return strValue, false
		}
	}

	return strValue, false
}

// parse struct fields
func (s *ModelReflector) parseStructFields(typ reflect.Type, val reflect.Value, tagNames ...string) {
	kind := typ.Kind()
	if kind == reflect.Struct {
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			tagVal, ignore := s.getTag(typField, types.TAG_NAME_DB, types.TAG_NAME_GORM, types.TAG_NAME_XORM, types.TAG_NAME_PROTOBUF)
			if ignore {
				continue
			}
			if !valField.IsValid() || !valField.CanInterface() {
				s.Dict[tagVal] = indirectValue(nil)
				continue
			}
			var tagSqlca string
			tagSqlca, ignore = s.getTag(typField, types.TAG_NAME_SQLCA)
			if !ignore {
				if tagSqlca == types.SQLCA_TAG_VALUE_IS_NULL {
					s.engine.setNullableColumns(tagVal)
				}
			}
			if typField.Type.Kind() == reflect.Struct {
				if tagVal == "" {
					s.parseStructFields(typField.Type, valField, tagNames...) //recurse every field that type is a struct
				} else {
					s.Dict[tagVal] = indirectValue(valField.Interface())
					s.Columns = append(s.Columns, tagVal)
				}
			} else if typField.Type.Kind() == reflect.Slice || typField.Type.Kind() == reflect.Map {
				if tagVal != "" {
					s.Dict[tagVal] = indirectValue(valField.Interface())
					s.Columns = append(s.Columns, tagVal)
				}
			} else {
				s.setValueByField(typField, valField, tagNames...) // save field tag value and field value to map
			}
		}
	}
}

// trim the field value's first and last blank character and save to map
func (s *ModelReflector) setValueByField(field reflect.StructField, val reflect.Value, tagNames ...string) {

	if len(tagNames) == 0 {
		log.Errorf("ModelReflector.setValueByField no tag to set value")
		return
	}

	var tagVal string
	for _, v := range tagNames {

		if v == types.TAG_NAME_SQLCA {
			continue
		}

		strTagValue, ignore := s.getTag(field, v)
		//parse db、json、protobuf tag
		tagVal = handleTagValue(v, strTagValue)
		if ignore {
			break
		}
		if tagVal != "" {
			s.Dict[tagVal] = indirectValue(val.Interface())
			s.Columns = append(s.Columns, tagVal)
			break
		}
	}

	for _, v := range tagNames {
		if v == types.TAG_NAME_SQLCA { //parse sqlca tag
			strTagValue, ignore := s.getTag(field, v)
			if !ignore && strTagValue != "" {
				vs := strings.Split(strTagValue, ",")
				for _, vv := range vs {
					if vv == types.SQLCA_TAG_VALUE_READ_ONLY { //column is read only
						s.engine.readOnly = append(s.engine.readOnly, tagVal)
					}
				}
			}
		}
	}
}

func (e *Engine) fetchRows(r *sql.Rows) (count int64, err error) {

	var i int
	for r.Next() {
		var c int64
		i++
		if e.getModelType() == types.ModelType_BaseType {
			if c, err = e.fetchRow(r, e.model.([]any)...); err != nil {
				log.Errorf("fetch row error [%v]", err.Error())
				return
			}
		} else {
			if c, err = e.fetchRow(r, e.model); err != nil {
				log.Errorf("fetch row error [%v]", err.Error())
				return
			}
		}
		count += c
	}
	return
}

// fetch row to struct or slice, must call rows.Next() before call this function
func (e *Engine) fetchRow(rows *sql.Rows, args ...any) (count int64, err error) {

	fetcher, _ := e.getFetcher(rows)

	for _, arg := range args {

		typ := reflect.TypeOf(arg)
		val := reflect.ValueOf(arg)

		if typ.Kind() == reflect.Ptr {

			typ = typ.Elem()
			val = val.Elem()
		}

		switch typ.Kind() {
		case reflect.Map:
			{
				err = e.fetchToMap(fetcher, arg)
				count++
			}
		case reflect.Slice:
			{
				if val.IsNil() {
					val.Set(reflect.MakeSlice(val.Type(), 0, 0)) //make slice for storage
				}
				for {
					fetcher, _ := e.getFetcher(rows)
					var elemVal reflect.Value
					var elemTyp reflect.Type
					elemTyp = val.Type().Elem()

					if elemTyp.Kind() == reflect.Ptr {
						elemVal = reflect.New(elemTyp.Elem())
					} else {
						elemVal = reflect.New(elemTyp).Elem()
					}

					if elemTyp.Kind() == reflect.Struct || elemTyp.Kind() == reflect.Ptr {
						err = e.fetchToStruct(fetcher, elemTyp, elemVal) // assign to struct type variant
					} else {
						err = e.fetchToBaseType(fetcher, elemTyp, elemVal) // assign to base type variant
					}

					val.Set(reflect.Append(val, elemVal))
					count++
					if !rows.Next() {
						break
					}
				}
			}
		case reflect.Struct:
			{
				if _, ok := val.Addr().Interface().(sql.Scanner); !ok {
					err = e.fetchToStruct(fetcher, typ, val)
				} else {
					err = e.fetchToBaseType(fetcher, typ, val)
				}
				count++
			}
		default:
			{
				err = e.fetchToBaseType(fetcher, typ, val)
				count++
			}
		}
	}
	return
}

func (e *Engine) getStructSliceKeyValues(excludeReadOnly bool) (keys []string, values [][]any) {

	typ := reflect.TypeOf(e.model)
	val := reflect.ValueOf(e.model)

	if typ.Kind() == reflect.Ptr {

		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Slice:
		{
			elemTyp := val.Type().Elem()
			for i := 0; i < val.Len(); i++ {
				elemVal := val.Index(i)
				if elemTyp.Kind() == reflect.Ptr {
					elemTyp = elemTyp.Elem()
					elemVal = elemVal.Elem()
				}

				if elemTyp.Kind() == reflect.Struct {
					var vs []any
					keys, vs = e.getStructFieldValues(elemTyp, elemVal, excludeReadOnly)
					values = append(values, vs)
				}
			}
		}
	default:
		{
			panic(fmt.Sprintf("expect struct got %v", typ.Name()))
		}
	}

	return
}

func (e *Engine) getStructFieldValues(typ reflect.Type, val reflect.Value, excludeReadOnly bool) (keys []string, values []any) {

	if typ.Kind() == reflect.Struct {

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}

			var tagName string
			var fieldVal any
			tagName = e.getTagValue(typField)
			if !valField.IsValid() || !valField.CanInterface() {
				fieldVal = indirectValue(nil)
			} else {
				fieldVal = indirectValue(valField.Interface())
			}
			if (fieldVal == "" || fieldVal == "0") && tagName == e.GetPkName() {
				continue
			}
			if excludeReadOnly {
				if typField.Tag.Get(types.TAG_NAME_SQLCA) == types.SQLCA_TAG_VALUE_READ_ONLY {
					continue
				}
			}

			if tagName != "" && tagName != types.SQLCA_TAG_VALUE_IGNORE {
				keys = append(keys, tagName)
				values = append(values, fieldVal)
			}
		}
	}
	return
}

// fetch cache data to struct or slice or map
func (e *Engine) fetchCache(fetchers []*Fetcher, args ...any) (count int64, err error) {

	for i, fetcher := range fetchers {

		for _, arg := range args {

			typ := reflect.TypeOf(arg)
			val := reflect.ValueOf(arg)

			if typ.Kind() == reflect.Ptr {

				typ = typ.Elem()
				val = val.Elem()
			}

			switch typ.Kind() {
			case reflect.Map:
				{
					err = e.fetchToMap(fetcher, arg)
					count++
				}
			case reflect.Slice:
				{
					if val.IsNil() {
						val.Set(reflect.MakeSlice(val.Type(), 0, 0)) //make slice for storage
					}
					for {
						fetcher = fetchers[i]
						elemTyp := val.Type().Elem()
						elemVal := reflect.New(elemTyp).Elem()

						if elemTyp.Kind() == reflect.Struct {
							err = e.fetchToStruct(fetcher, elemTyp, elemVal) // assign to struct type variant
						} else {
							err = e.fetchToBaseType(fetcher, elemTyp, elemVal) // assign to base type variant
						}

						val.Set(reflect.Append(val, elemVal))
						count++
						if len(fetchers) == i+1 {
							break
						}
						i++
					}
				}
			case reflect.Struct:
				{
					err = e.fetchToStruct(fetcher, typ, val)
					count++
				}
			default:
				{
					e.fetchToBaseType(fetcher, typ, val)
					count++
				}
			}
		}
	}
	return
}

func (e *Engine) getFetcher(rows *sql.Rows) (fetcher *Fetcher, err error) {

	fetcher = &Fetcher{}
	fetcher.cols, _ = rows.Columns()
	fetcher.count = len(fetcher.cols)
	fetcher.types, _ = rows.ColumnTypes()
	fetcher.arrValues = make([][]byte, fetcher.count)
	fetcher.mapValues = make(map[string]string)
	scans := make([]any, fetcher.count)

	for i := range fetcher.arrValues {
		scans[i] = &fetcher.arrValues[i]
	}

	if err = rows.Scan(scans...); err != nil {

		return
	}
	for i, v := range fetcher.arrValues {
		fetcher.mapValues[fetcher.cols[i]] = string(v)
	}
	return
}

// fetch row data to map
func (e *Engine) fetchToMap(fetcher *Fetcher, arg any) (err error) {

	typ := reflect.TypeOf(arg)

	if typ.Kind() == reflect.Ptr {

		for k, v := range fetcher.mapValues {
			m := *arg.(*map[string]string) //just support map[string]string type
			m[k] = v
		}
	}

	return
}

// fetch row data to struct
func (e *Engine) fetchToStruct(fetcher *Fetcher, typ reflect.Type, val reflect.Value) (err error) {

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)
			e.fetchToStructField(fetcher, typField.Type, typField, valField)
		}
	}
	return
}

func (e *Engine) fetchToStructField(fetcher *Fetcher, typ reflect.Type, field reflect.StructField, val reflect.Value, ptr ...reflect.Value) {

	switch typ.Kind() {
	case reflect.Struct:
		{
			e.fetchToStructAny(fetcher, field, val, ptr...)
		}
	case reflect.Slice:
		if e.getTagValue(field) != "" {
			_ = e.fetchToJsonObject(fetcher, field, val, ptr...)
		}
	case reflect.Map: //ignore...
	case reflect.Ptr:
		{
			typElem := field.Type.Elem()
			if val.IsNil() {
				valNew := reflect.New(typElem)
				val.Set(valNew)
			}
			e.fetchToStructField(fetcher, typElem, field, val.Elem(), val)
		}
	default:
		{
			_ = e.setValueByField(fetcher, field, val, ptr...) //assign value to struct field
		}
	}
}

func (e *Engine) fetchToStructAny(fetcher *Fetcher, field reflect.StructField, val reflect.Value, ptr ...reflect.Value) {
	if _, ok := val.Addr().Interface().(sql.Scanner); ok {
		e.fetchToScanner(fetcher, field, val)
	} else {
		if e.getTagValue(field) != "" {
			_ = e.fetchToJsonObject(fetcher, field, val, ptr...)
		} else {
			_ = e.fetchToStruct(fetcher, field.Type, val)
		}
	}
}

func setNilPtr(vals ...reflect.Value) {
	for _, val := range vals {
		if !val.CanAddr() {
			return
		}
		val.Set(reflect.Zero(val.Type()))
	}
}

// json string unmarshal to struct/slice
func (e *Engine) fetchToJsonObject(fetcher *Fetcher, field reflect.StructField, val reflect.Value, ptr ...reflect.Value) (err error) {
	var assigned bool
	defer func() {
		if !assigned && len(ptr) != 0 {
			setNilPtr(ptr...)
		}
	}()
	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == types.SQLCA_TAG_VALUE_IGNORE {
		return nil
	}

	if v, ok := fetcher.mapValues[strDbTagVal]; ok {
		vp := val.Addr()
		if strings.TrimSpace(v) != "" && canUnmarshalJson(v) {
			if err = json.Unmarshal([]byte(v), vp.Interface()); err != nil {
				return log.Errorf("json.Unmarshal [%s] error [%s]", v, err)
			}
			assigned = true
		} else {
			//if struct field is a slice type and content is nil make space for it
			if field.Type.Kind() == reflect.Slice {
				val.Set(reflect.MakeSlice(field.Type, 0, 0))
			}
		}
	}
	return nil
}

// fetch to struct object by customize scanner
func (e *Engine) fetchToScanner(fetcher *Fetcher, field reflect.StructField, val reflect.Value) {
	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == types.SQLCA_TAG_VALUE_IGNORE {
		return
	}
	if v, ok := fetcher.mapValues[strDbTagVal]; ok {
		vp := val.Addr()
		d := vp.Interface().(sql.Scanner)
		if v == "" {
			return
		}
		if err := d.Scan(v); err != nil {
			log.Errorf("scan '%v' to scanner [%+v] error [%+v]", v, vp.Interface(), err.Error())
		}
	}
}

func (e *Engine) fetchToBaseType(fetcher *Fetcher, typ reflect.Type, val reflect.Value) (err error) {

	v := fetcher.arrValues[fetcher.arrIndex]
	e.setValue(typ, val, string(v))
	fetcher.arrIndex++
	return
}

func handleTagValue(strTagName, strTagValue string) string {

	if strTagValue == "" {
		return ""
	}

	if strTagName == types.TAG_NAME_JSON {
		vs := strings.Split(strTagValue, ",")
		strTagValue = vs[0]
	} else if strTagName == types.TAG_NAME_PROTOBUF {
		//parse protobuf tag value
		vs := strings.Split(strTagValue, ",")
		for _, vv := range vs {
			ss := strings.Split(vv, "=")
			if len(ss) <= 1 {
				continue
			} else {
				if ss[0] == types.PROTOBUF_VALUE_NAME {
					strTagValue = ss[1]
					return strTagValue
				}
			}
		}
	} else {
		vs := strings.Split(strTagValue, ";")
		col := vs[0]
		if strings.Contains(col, ":") {
			vs = strings.Split(col, ":")
			col = vs[1]
		}
		strTagValue = col
	}
	return strTagValue
}

func (e *Engine) getTagValue(sf reflect.StructField) (strValue string) {

	for _, v := range e.dbTags { //support multiple tag
		strValue = handleTagValue(v, sf.Tag.Get(v))
		if strValue != "" {
			return
		}
	}
	return
}

// 按结构体字段标签赋值
func (e *Engine) setValueByField(fetcher *Fetcher, field reflect.StructField, val reflect.Value, ptr ...reflect.Value) (err error) {

	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == types.SQLCA_TAG_VALUE_IGNORE {
		return nil
	}
	var assigned bool
	v, ok := fetcher.mapValues[strDbTagVal]
	if ok {
		assigned = e.setValue(field.Type, val, v)
	}
	if !assigned && len(ptr) != 0 {
		setNilPtr(ptr...)
	}
	return nil
}

// 将string存储的值赋值到变量
func (e *Engine) setValue(typ reflect.Type, val reflect.Value, v string) bool {
	if strings.TrimSpace(v) == "" {
		return false
	}
	switch typ.Kind() {
	case reflect.Struct:
		s, ok := val.Addr().Interface().(sql.Scanner)
		if !ok {
			log.Warnf("struct type %s not implement sql.Scanner interface", typ.Name())
			return false
		}
		if err := s.Scan(v); err != nil {
			panic(fmt.Sprintf("scan value %s to sql.Scanner implement object error [%s]", v, err))
		}
	case reflect.String:
		val.SetString(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(v, 10, 64)
		val.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(v, 10, 64)
		val.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, _ := strconv.ParseFloat(v, 64)
		val.SetFloat(i)
	case reflect.Bool:
		if v == "true" { //postgresql boolean value
			val.SetBool(true)
		} else { //other database integer value
			i, _ := strconv.ParseUint(v, 10, 64)
			if i != 0 {
				val.SetBool(true)
			}
		}
	case reflect.Ptr:
		typ = typ.Elem()
		e.setValue(typ, val, v)
	default:
		panic(fmt.Sprintf("can't assign value [%v] to variant type [%v]\n", v, typ.Kind()))
		return false
	}
	return true
}

func convertBool2Int(v any) any {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	switch typ.Kind() {
	case reflect.Bool:
		{
			if val.Interface() == false {
				return 0
			} else {
				return 1
			}
		}
	}
	return v
}
