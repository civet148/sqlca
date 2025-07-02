package types

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// Builder builder interface
type Builder interface {
	Writer
	WriteQuoted(field interface{})
	AddVar(Writer, ...interface{})
}

type Expr struct {
	SQL  string
	Vars []any
}

// Build build raw expression
func (expr Expr) Build(builder Builder) {
	var (
		idx int
	)

	for _, v := range []byte(expr.SQL) {
		if v == '?' && len(expr.Vars) > idx {
			builder.AddVar(builder, expr.Vars[idx])
			idx++
		}
	}

	if idx < len(expr.Vars) {
		for _, v := range expr.Vars[idx:] {
			builder.AddVar(builder, sql.NamedArg{Value: v})
		}
	}
}

func (expr Expr) String() string {
	query := strings.Replace(expr.SQL, "?", "%v", -1)
	vars := expr.quoteValues(expr.Vars...)
	return fmt.Sprintf(query, vars...)
}

func (expr Expr) quoteValues(values ...any) (vars []any) {
	for _, v := range values {
		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.String:
			vars = append(vars, fmt.Sprintf("'%v'", v))
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
