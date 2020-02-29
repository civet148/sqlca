package sqlca

import (
	"fmt"
	"reflect"
	"strings"
)

func Struct(v interface{}) *Structure {

	return &Structure{
		value: v,
		dict:  make(map[string]string),
	}
}

type Structure struct {
	value interface{}       //value
	dict  map[string]string //dictionary of structure tag and value
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
