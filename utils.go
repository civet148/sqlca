package sqlca

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type queryInterfaceType int

const (
	queryInterface_Unknown queryInterfaceType = 0
	queryInterface_String  queryInterfaceType = 1
	queryInterface_Map     queryInterfaceType = 2
)

// convertCamelToSnake converts a CamelCase string to snake_case
func convertCamelToSnake(s string) string {
	// 使用正则表达式匹配大写字母前的所有字符
	re := regexp.MustCompile("([A-Z])")
	// 将匹配到的大写字母替换为下划线加小写字母
	snake := re.ReplaceAllString(s, "_$1")
	// 将字符串转换为小写并去掉开头的下划线
	return strings.ToLower(strings.TrimPrefix(snake, "_"))
}

// checkTruth check string true or not
func checkTruth(vals ...string) bool {
	for _, val := range vals {
		if val != "" && !strings.EqualFold(val, "false") {
			return true
		}
	}
	return false
}

func parseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return ""
}

func parseQueryInterface(query interface{}) queryInterfaceType {
	switch query.(type) {
	case map[string]interface{}:
		return queryInterface_Map
	case string:
		return queryInterface_String
	}
	return queryInterface_Unknown
}

func isQuestionPlaceHolder(query string, args ...interface{}) bool {
	count := strings.Count(query, "?")
	if count > 0 {
		if count != len(args) {
			log.Warnf("question placeholder count %d not equal args count %d", count, len(args))
			return true
		}
		return true
	}
	return false
}

type StringBuilder struct {
	builder strings.Builder
	args    []interface{}
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

func (s *StringBuilder) Append(query string, args ...interface{}) *StringBuilder {
	var strQuery string
	if isQuestionPlaceHolder(query, args...) { //question placeholder exist
		s.builder.WriteString(query)
		s.args = append(s.args, args...)
	} else {
		strQuery = " " + fmt.Sprintf(query, args...) + " "
	}
	s.builder.WriteString(strQuery)
	return s
}

func (s *StringBuilder) String() string {
	return s.builder.String()
}

func (s *StringBuilder) Args() []interface{} {
	return s.args
}

func indirectValue(v any) any {
	if v == nil {
		return "NULL"
	}

	value := reflect.ValueOf(v)
	// 循环处理指针，直到获取到非指针的值
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return "NULL"
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(value.Uint())
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.String:
		return value.String()
	case reflect.Struct:
		if valuer, ok := value.Interface().(driver.Valuer); ok {
			result, err := valuer.Value()
			if err == nil {
				return result
			}
		}
		data, err := json.Marshal(value.Interface())
		if err == nil {
			return string(data)
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		data, err := json.Marshal(value.Interface())
		if err == nil {
			return string(data)
		}
	default:
		return fmt.Sprintf("%v", v)
	}
	return v
}

// ... existing code ...
