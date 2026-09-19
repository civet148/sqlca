package types

import (
	"fmt"
	"reflect"
	"strings"
)

type Expr struct {
	SQL  string
	Vars []any
}

func (expr Expr) RawSQL(adapters ...AdapterType) string {
	var adapter AdapterType
	if len(adapters) > 0 {
		adapter = adapters[0]
	}
	query := strings.Replace(expr.SQL, "?", "%v", -1)
	vars := expr.quoteValues(adapter, expr.Vars...)
	return fmt.Sprintf(query, vars...)
}

func (expr Expr) quoteValues(adapter AdapterType, values ...any) (vars []any) {
	for _, v := range values {
		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.String:
			s := PreventSqlInject(adapter, v.(string))
			vars = append(vars, fmt.Sprintf("'%v'", s))
		case reflect.Struct:
			if sn, ok := v.(SqlNull); ok {
				// 判断类型名称和包路径是否一致
				if typ.Name() == "SqlNull" && typ.PkgPath() == reflect.TypeOf(SqlNull{}).PkgPath() {
					vars = append(vars, sn.String())
				}
			} else {
				vars = append(vars, v)
			}
		default:
			vars = append(vars, v)
		}
	}
	return vars
}

// PreventSqlInject handle special characters, prevent SQL inject
func PreventSqlInject(adapter AdapterType, strIn string) (strOut string) {

	strIn = strings.TrimSpace(strIn) //trim blank characters
	switch adapter {
	case AdapterSqlx_MySQL:
		strIn = strings.Replace(strIn, `\`, `\\`, -1)
		strIn = strings.Replace(strIn, `'`, `\'`, -1)
		strIn = strings.Replace(strIn, `"`, `\"`, -1)
	case AdapterSqlx_Postgres, AdapterSqlx_OpenGauss:
		strIn = strings.Replace(strIn, `'`, `''`, -1)
	case AdapterSqlx_Mssql:
		strIn = strings.Replace(strIn, `'`, `''`, -1)
	case AdapterSqlx_Sqlite:
	}

	return strIn
}
