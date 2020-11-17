package sqlca

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"reflect"
	"strconv"
	"strings"
)

type ModelReflector struct {
	value  interface{}            //value
	engine *Engine                // database engine
	dict   map[string]interface{} //dictionary of structure tag and value
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
		dict:   make(map[string]interface{}),
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

// handle special characters, prevent SQL inject
func handleSpecialChars(strIn string) (strOut string) {

	strIn = strings.TrimSpace(strIn) //trim blank characters
	strIn = strings.Replace(strIn, `\`, `\\`, -1)
	strIn = strings.Replace(strIn, `'`, `\'`, -1)
	strIn = strings.Replace(strIn, `"`, `\"`, -1)

	return strIn
}

// parse struct tag and value to map
func (s *ModelReflector) ToMap(tagNames ...string) map[string]interface{} {

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
			s.parseStructField(typ, val, tagNames...)
		}
	case reflect.Slice:
		{
			typ = val.Type().Elem()
			val = reflect.New(typ).Elem()
			s.parseStructField(typ, val, tagNames...)
		}
	case reflect.Map:
		{
			if v, ok := s.value.(*map[string]interface{}); ok {
				s.dict = *v
				break
			}
			if v, ok := s.value.(map[string]interface{}); ok {
				s.dict = v
				break
			}
			if v, ok := s.value.(*map[string]string); ok {
				s.dict = s.convertMapString(*v)
				break
			}
			if v, ok := s.value.(map[string]string); ok {
				s.dict = s.convertMapString(v)
				break
			}
		}
	default:
		log.Warnf("kind [%v] not support yet", typ.Kind())
	}
	return s.dict
}

func (s *ModelReflector) convertMapString(ms map[string]string) (mi map[string]interface{}) {
	mi = make(map[string]interface{}, 10)
	for k, v := range ms {
		mi[k] = v
	}
	return
}

// get struct field's tag value
func (s *ModelReflector) getTag(sf reflect.StructField, tagName string) (strValue string) {

	strValue = sf.Tag.Get(tagName)
	if strValue == SQLCA_TAG_VALUE_IGNORE {
		return ""
	}
	return
}

// parse struct fields
func (s *ModelReflector) parseStructField(typ reflect.Type, val reflect.Value, tagNames ...string) {

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
			if !valField.IsValid() || !valField.CanInterface() {
				continue
			}
			//log.Debugf("reflect.Struct field [%v] kind [%+v]", i, typField.Type.Kind())
			if typField.Type.Kind() == reflect.Struct {

				if _, ok := valField.Interface().(driver.Valuer); ok {
					s.parseValuer(typField, valField, tagNames...)
				} else {
					s.parseStructField(typField.Type, valField, tagNames...) //recurse every field that type is a struct
				}
			} else {
				s.setValueByField(typField, valField, tagNames...) // save field tag value and field value to map
			}
		}
	}
}

//parse decimal
func (s *ModelReflector) parseValuer(field reflect.StructField, val reflect.Value, tagNames ...string) {

	s.setValueByField(field, val, tagNames...)
}

//trim the field value's first and last blank character and save to map
func (s *ModelReflector) setValueByField(field reflect.StructField, val reflect.Value, tagNames ...string) {

	if len(tagNames) == 0 {
		log.Errorf("ModelReflector.setValueByField no tag to set value")
		return
	}

	var tagVal string
	for _, v := range tagNames {

		if v == TAG_NAME_SQLCA {
			continue
		}
		//parse db、json、protobuf tag
		tagVal = handleTagValue(v, s.getTag(field, v))
		if tagVal != "" {
			//log.Debugf("ModelReflector.setValueByField tag [%v] value [%+v]", tagVal, val.Interface())
			if d, ok := val.Interface().(driver.Valuer); ok {
				s.dict[tagVal], _ = d.Value()
			} else {
				s.dict[tagVal] = val.Interface()
			}
			break
		}
	}

	for _, v := range tagNames {

		if v == TAG_NAME_SQLCA { //parse sqlca tag
			strTagValue := s.getTag(field, v)
			vs := strings.Split(strTagValue, ",")
			for _, vv := range vs {
				if vv == SQLCA_TAG_VALUE_READ_ONLY { //column is read only
					s.engine.readOnly = append(s.engine.readOnly, tagVal)
					//log.Debugf("%v [%v]", tagVal, SQLCAL_TAG_VALUE_READ_ONLY)
				}
			}
		}
	}
}

func (e *Engine) fetchRows(r *sql.Rows) (count int64, err error) {

	for r.Next() {
		var c int64

		if e.getModelType() == ModelType_BaseType {
			if c, err = e.fetchRow(r, e.model.([]interface{})...); err != nil {
				log.Errorf("fetchRow error [%v]", err.Error())
				return
			}
		} else {
			if c, err = e.fetchRow(r, e.model); err != nil {
				log.Errorf("fetchRow error [%v]", err.Error())
				return
			}
		}
		count += c
	}
	return
}

//fetch row to struct or slice, must call rows.Next() before call this function
func (e *Engine) fetchRow(rows *sql.Rows, args ...interface{}) (count int64, err error) {

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
	return
}

func (e *Engine) getStructSliceKeyValues(excludeReadOnly bool) (keys []string, values [][]string) {

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
				}
				if elemVal.Kind() == reflect.Ptr {
					elemVal = elemVal.Elem()
				}

				if elemTyp.Kind() == reflect.Struct {
					var vs []string
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

func (e *Engine) getStructFieldValues(typ reflect.Type, val reflect.Value, excludeReadOnly bool) (keys, values []string) {

	if typ.Kind() == reflect.Struct {

		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			if !valField.IsValid() || !valField.CanInterface() {
				//fmt.Printf("Filed [%s] tag(%s)  is not valid \n", typField.Type.Name(), e.getTagValue(typField))
				return
			}
			strTagVal := e.getTagValue(typField)
			strFieldVal := fmt.Sprintf("%v", valField)

			if e.isPkValueNil() && strTagVal == e.GetPkName() {
				continue
			}

			if excludeReadOnly {
				if typField.Tag.Get(TAG_NAME_SQLCA) == SQLCA_TAG_VALUE_READ_ONLY {
					continue
				}
			}

			if strTagVal != "" && strTagVal != SQLCA_TAG_VALUE_IGNORE {
				keys = append(keys, strTagVal)
				values = append(values, strFieldVal)
				//log.Debugf("filed tag name [%v] value [%v]", strTagVal, strFieldVal)
			}
		}
	}
	return
}

//fetch cache data to struct or slice or map
func (e *Engine) fetchCache(fetchers []*Fetcher, args ...interface{}) (count int64, err error) {

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
	scans := make([]interface{}, fetcher.count)

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

//fetch row data to map
func (e *Engine) fetchToMap(fetcher *Fetcher, arg interface{}) (err error) {

	typ := reflect.TypeOf(arg)

	if typ.Kind() == reflect.Ptr {

		for k, v := range fetcher.mapValues {
			m := *arg.(*map[string]string) //just support map[string]string type
			m[k] = v
		}
	}

	return
}

//fetch row data to struct
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

			if !valField.IsValid() {
				log.Warnf("struct field type (%s) is not valid or can't interface{} \n", typField.Type.Name())
				continue
			}
			switch typField.Type.Kind() {
			case reflect.Struct:
				{
					if _, ok := valField.Addr().Interface().(sql.Scanner); ok {
						e.fetchToScanner(fetcher, typField, valField)
					} else {
						if e.getTagValue(typField) != "" {
							_ = e.fetchToJsonObject(fetcher, typField, valField)
						} else {
							_ = e.fetchToStruct(fetcher, typField.Type, valField)
						}
					}
				}
			case reflect.Slice:
				if e.getTagValue(typField) != "" {
					_ = e.fetchToJsonObject(fetcher, typField, valField)
				}
			case reflect.Map: //ignore...
			case reflect.Ptr:
				{
					typElem := typField.Type.Elem()
					if typElem.Kind() == reflect.Struct {
						if valField.IsNil() {
							valNew := reflect.New(typElem)
							valField.Set(valNew)
						}
						if e.getTagValue(typField) != "" {
							_ = e.fetchToJsonObject(fetcher, typField, valField)
						} else {
							_ = e.fetchToStruct(fetcher, typField.Type, valField)
						}
					}
				}
			default:
				{
					_ = e.setValueByField(fetcher, typField, valField) //assign value to struct field
				}
			}
		}
	}

	return
}

//json string unmarshal to struct/slice
func (e *Engine) fetchToJsonObject(fetcher *Fetcher, field reflect.StructField, val reflect.Value) (err error) {
	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == SQLCA_TAG_VALUE_IGNORE {
		return
	}

	if v, ok := fetcher.mapValues[strDbTagVal]; ok {
		vp := val.Addr()
		if err = json.Unmarshal([]byte(v), vp.Interface()); err != nil {
			log.Errorf("json.Unmarshal error [%s]", err)
			return
		}
	}
	return
}

//fetch to struct object by customize scanner
func (e *Engine) fetchToScanner(fetcher *Fetcher, field reflect.StructField, val reflect.Value) {
	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == SQLCA_TAG_VALUE_IGNORE {
		return
	}
	if v, ok := fetcher.mapValues[strDbTagVal]; ok {
		vp := val.Addr()
		d := vp.Interface().(sql.Scanner)
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

	if strTagName == TAG_NAME_JSON && strTagValue != "" {

		vs := strings.Split(strTagValue, ",")
		strTagValue = vs[0]
	} else if strTagName == TAG_NAME_PROTOBUF && strTagValue != "" {
		//parse protobuf tag value
		vs := strings.Split(strTagValue, ",")
		for _, vv := range vs {
			ss := strings.Split(vv, "=")
			if len(ss) <= 1 {
				//log.Warnf("protobuf tag value [%v] is not a invalid format", strTagValue)
				continue
			} else {
				if ss[0] == PROTOBUF_VALUE_NAME {
					strTagValue = ss[1]
					return strTagValue
				}
			}
		}
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

//按结构体字段标签赋值
func (e *Engine) setValueByField(fetcher *Fetcher, field reflect.StructField, val reflect.Value) (err error) {

	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)
	if strDbTagVal == SQLCA_TAG_VALUE_IGNORE {
		return
	}
	if v, ok := fetcher.mapValues[strDbTagVal]; ok {
		e.setValue(field.Type, val, v)
	}
	return
}

//将string存储的值赋值到变量
func (e *Engine) setValue(typ reflect.Type, val reflect.Value, v string) {
	switch typ.Kind() {

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
		i, _ := strconv.ParseUint(v, 10, 64)
		val.SetBool(true)
		if i == 0 {
			val.SetBool(false)
		}
	default:
		panic(fmt.Sprintf("can't assign value [%v] to variant type [%v]\n", v, typ.Kind()))
		return
	}
}
