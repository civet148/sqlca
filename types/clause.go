package types

import (
	"fmt"
	"reflect"
	"strings"
)

// Builder builder interface
type Builder interface {
	AddSQL(string)
	AddVar(...any)
}

type Expr struct {
	SQL  string
	Vars []any
}

// Build build raw expression
func (expr Expr) Build(builder Builder) {
	var idx int
	builder.AddSQL(expr.SQL)
	for _, v := range []byte(expr.SQL) {
		if v == '?' && len(expr.Vars) > idx {
			builder.AddVar(expr.Vars[idx])
			idx++
		}
	}
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
			var sn = v.(SqlNull)
			// 判断类型名称和包路径是否一致
			if typ.Name() == "SqlNull" && typ.PkgPath() == reflect.TypeOf(SqlNull{}).PkgPath() {
				vars = append(vars, sn.String())
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
