package sqlca

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/types"
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

var (
	inBlank    = " IN "
	inQuestion = "?"
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
	args    []any
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

func (s *StringBuilder) Append(query string, args ...any) *StringBuilder {
	var strQuery string
	if isQuestionPlaceHolder(query, args...) { //question placeholder exist
		query = strings.Replace(query, "?", "%v", -1) + " "
		s.builder.WriteString(query)
		s.args = append(s.args, args...)
	} else {
		strQuery = " " + fmt.Sprintf(query, args...) + " "
	}
	s.builder.WriteString(strQuery)
	return s
}

func (s *StringBuilder) String() string {
	return fmt.Sprintf(s.builder.String(), s.args...)
}

func (s *StringBuilder) Args() []interface{} {
	return s.args
}

func indirectValue(v any, isClauses ...bool) any {
	var isClauseVal bool
	if len(isClauses) > 0 {
		isClauseVal = isClauses[0]
	}

	if v == nil {
		return types.SqlNull{}
	}

	value := reflect.ValueOf(v)
	// 循环处理指针，直到获取到非指针的值
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return types.SqlNull{}
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
		} else {
			data, err := json.Marshal(value.Interface())
			if err == nil {
				return string(data)
			}
		}
	case reflect.Slice, reflect.Array:
		if !isClauseVal {
			data, err := json.Marshal(value.Interface())
			if err == nil {
				return string(data)
			}
		} else {
			var strValues []string
			var val = reflect.ValueOf(v)
			n := val.Len()
			for i := 0; i < n; i++ {
				sv := fmt.Sprintf("'%v'", indirectValue(val.Index(i).Interface()))
				strValues = append(strValues, sv)
			}
			return types.SqlClauseValue{Val: strings.Join(strValues, ",")}
		}
	case reflect.Map:
		data, err := json.Marshal(value.Interface())
		if err == nil {
			return string(data)
		}
	default:
		return fmt.Sprintf("%v", v)
	}
	return v
}

func quotedValue(v any) (sv string) {
	val := reflect.ValueOf(v)
	val = reflect.Indirect(val)

	switch val.Kind() {
	case reflect.String:
		sv = fmt.Sprintf("'%v'", v.(string))
	case reflect.Struct:
		sn, ok := isSqlNull(v)
		if ok {
			sv = sn.String()
		} else {
			var scv types.SqlClauseValue
			if scv, ok = isSqlClauseValue(v); ok {
				return scv.String()
			} else {
				sv = fmt.Sprintf("'%v'", indirectValue(v))
			}
		}
	default:
		sv = fmt.Sprintf("'%v'", indirectValue(v))
	}
	return sv
}

func isSqlNull(v any) (types.SqlNull, bool) {
	val := reflect.ValueOf(v)
	val = reflect.Indirect(val)
	typ := val.Type()
	// 判断类型名称和包路径是否一致
	if typ.Name() == "SqlNull" && typ.PkgPath() == reflect.TypeOf(types.SqlNull{}).PkgPath() {
		return val.Interface().(types.SqlNull), true
	}
	return types.SqlNull{}, false
}

func isSqlClauseValue(v any) (types.SqlClauseValue, bool) {
	val := reflect.ValueOf(v)
	val = reflect.Indirect(val)
	typ := val.Type()
	// 判断类型名称和包路径是否一致
	if typ.Name() == "SqlClauseValue" && typ.PkgPath() == reflect.TypeOf(types.SqlClauseValue{}).PkgPath() {
		return val.Interface().(types.SqlClauseValue), true
	}
	return types.SqlClauseValue{}, false
}

// 判断是否为IN、NOT IN条件
func hasClauseSlice(query string, args ...any) bool {
	var upper = strings.ToUpper(query)
	if strings.Contains(upper, inBlank) && strings.Contains(upper, inQuestion) {
		return true
	}
	return false
}
