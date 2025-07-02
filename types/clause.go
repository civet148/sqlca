package types

import (
	"database/sql"
	"fmt"
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
	query := strings.Replace(expr.SQL, "?", "'%v'", -1)
	return fmt.Sprintf(query, expr.Vars...)
}
