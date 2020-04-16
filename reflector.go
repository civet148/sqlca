package sqlca

import (
	"database/sql"
	"fmt"
	"github.com/civet148/gotools/log"
	"reflect"
	"strconv"
	"strings"
)

type ModelReflector struct {
	value interface{}            //value
	dict  map[string]interface{} //dictionary of structure tag and value
}

type Fetcher struct {
	count     int               //column count in db table
	cols      []string          //column names in db table
	types     []*sql.ColumnType //column types in db table
	arrValues [][]byte          //value slice
	mapValues map[string]string //value map
	arrIndex  int               //fetch index
}

func newReflector(v interface{}) *ModelReflector {

	return &ModelReflector{
		value: v,
		dict:  make(map[string]interface{}),
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
func (s *ModelReflector) ToMap(tagName string) map[string]interface{} {

	typ := reflect.TypeOf(s.value)
	val := reflect.ValueOf(s.value)

	if typ.Kind() == reflect.Ptr { // pointer type
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct { // struct data model
		s.parseStructField(typ, val, tagName)
	} else if typ.Kind() == reflect.Slice { // struct slice data model
		typ = val.Type().Elem()
		val = reflect.New(typ).Elem()
		s.parseStructField(typ, val, tagName)
	}
	return s.dict
}

// get struct field's tag value
func (s *ModelReflector) getTag(sf reflect.StructField, tagName string) string {

	return sf.Tag.Get(tagName)
}

// parse struct fields
func (s *ModelReflector) parseStructField(typ reflect.Type, val reflect.Value, tagName string) {

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
			if typField.Type.Kind() == reflect.Struct {
				s.parseStructField(typField.Type, valField, tagName) //recurse every field that type is a struct
			} else {
				s.setValueByField(typField, valField, tagName) // save field tag value and field value to map
			}
		}
	}
}

//trim the field value's first and last blank character and save to map
func (s *ModelReflector) setValueByField(field reflect.StructField, val reflect.Value, tagName string) {

	dbTags := strings.Split(tagName, ",")
	if len(dbTags) == 0 {
		log.Errorf("ModelReflector.setValueByField no tag to set value")
		return
	}

	for _, v := range dbTags {
		tagVal := s.getTag(field, v)
		if tagVal != "" {
			s.dict[tagVal] = val.Interface()
			return
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

	fetcher, _ := e.getFecther(rows)

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
					fetcher, _ := e.getFecther(rows)
					elemTyp := val.Type().Elem()
					elemVal := reflect.New(elemTyp).Elem()

					if elemTyp.Kind() == reflect.Struct {
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

func (e *Engine) getFecther(rows *sql.Rows) (fetcher *Fetcher, err error) {

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

//fetch row data to struct/slice
func (e *Engine) fetchToStruct(fetcher *Fetcher, typ reflect.Type, val reflect.Value) (err error) {

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
				//fmt.Printf("Filed [%s] tag(%s)  is not valid \n", typField.Type.Name(), e.getTagValue(typField))
				return
			}
			switch typField.Type.Kind() {
			case reflect.Struct:
				{
					e.fetchToStruct(fetcher, typField.Type, valField)
				}
			default:
				{
					e.setValueByField(fetcher, typField, valField) //assign value to struct field
				}
			}
		}
	}

	return
}

func (e *Engine) fetchToBaseType(fetcher *Fetcher, typ reflect.Type, val reflect.Value) (err error) {

	v := fetcher.arrValues[fetcher.arrIndex]
	e.setValue(typ, val, string(v))
	fetcher.arrIndex++
	return
}

func (e *Engine) getTagValue(sf reflect.StructField) (strValue string) {

	var tagName string
	tagName = TAG_NAME_DB
	strValue = sf.Tag.Get(tagName)
	if strValue == "" {
		for _, v := range e.dbTags { //support multiple tag
			strValue = sf.Tag.Get(v)
			if strValue != "" {
				return
			}
		}
	}
	return
}

//按结构体字段标签赋值
func (e *Engine) setValueByField(fetcher *Fetcher, field reflect.StructField, val reflect.Value) (err error) {

	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field)

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
