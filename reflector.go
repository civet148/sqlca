package sqlca

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Structure struct {
	value interface{}       //value
	dict  map[string]string //dictionary of structure tag and value
}

type Fetcher struct {
	count     int               //column count
	row       *sql.Row          //first row
	cols      []string          //column names
	types     []*sql.ColumnType //column types
	arrValues [][]byte          //value slice
	mapValues map[string]string //value map
	arrIndex  int               //fetch index
}

func Struct(v interface{}) *Structure {

	return &Structure{
		value: v,
		dict:  make(map[string]string),
	}
}

// parse struct tag and value to map
func (s *Structure) ToMap(tagName string) (m map[string]string) {

	typ := reflect.TypeOf(s.value)
	val := reflect.ValueOf(s.value)

	if typ.Kind() == reflect.Ptr { // pointer type
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct { // struct type

		s.parseStructField(typ, val, tagName)

	} else if typ.Kind() == reflect.Slice {

	} else {
		assert(false, "not a struct or slice object")
	}
	return s.dict
}

// get struct field's tag value
func (s *Structure) getTag(sf reflect.StructField, tagName string) string {

	return sf.Tag.Get(tagName)
}

// parse struct fields
func (s *Structure) parseStructField(typ reflect.Type, val reflect.Value, tagName string) {

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
func (s *Structure) setValueByField(field reflect.StructField, val reflect.Value, tagName string) {

	tagVal := s.getTag(field, tagName)
	if tagVal != "" {
		strVal := fmt.Sprintf("%v", val.Interface())
		s.dict[tagVal] = fmt.Sprintf("%v", strings.TrimSpace(strVal)) //trim the first and last blank character and save to map
	}
}

//fetch row to struct or slice, must call rows.Next() before call this function
func (e *Engine) fetchRow(rows *sql.Rows, arg interface{}) (count int64, err error) {

	fetcher, _ := e.getFecther(rows)

	argment := arg
	typ := reflect.TypeOf(argment)
	val := reflect.ValueOf(argment)

	if typ.Kind() == reflect.Ptr {

		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Map:
		{
			err = e.fetchToMap(fetcher, argment)
			count++
		}
	case reflect.Slice:
		{
			val.Set(reflect.MakeSlice(val.Type(), 0, 0)) //make slice for storage
			for {
				fetcher, _ := e.getFecther(rows)
				elemTyp := val.Type().Elem()
				elemVal := reflect.New(elemTyp).Elem()
				err = e.fetchToStruct(rows, fetcher, elemTyp, elemVal) //assign to struct object
				val.Set(reflect.Append(val, elemVal))
				count++
				if !rows.Next() {
					break
				}
			}
		}
	case reflect.Struct:
		{
			err = e.fetchToStruct(rows, fetcher, typ, val)
			count++
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
func (e *Engine) fetchToStruct(rows *sql.Rows, fetcher *Fetcher, typ reflect.Type, val reflect.Value) (err error) {

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
				fmt.Printf("Filed [%s] tag(%s)  is not valid \n", typField.Type.Name(), e.getTagValue(typField, TAG_NAME_DB))
				return
			}
			switch typField.Type.Kind() {
			case reflect.Struct:
				{
					e.fetchToStruct(rows, fetcher, typField.Type, valField)
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

func (e *Engine) getTagValue(sf reflect.StructField, tagName string) string {

	return sf.Tag.Get(tagName)
}

//按结构体字段标签赋值
func (e *Engine) setValueByField(fetcher *Fetcher, field reflect.StructField, val reflect.Value) (err error) {

	//优先给有db标签的成员变量赋值
	strDbTagVal := e.getTagValue(field, TAG_NAME_DB)

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
		fmt.Printf("can't assign value to this type [%v]\n", typ.Kind())
		return
	}
}
