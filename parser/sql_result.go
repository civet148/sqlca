package parser

import (
	"github.com/xwb1989/sqlparser"
)

type SqlResult struct {
	Type SqlType             `json:"type"`
	SQL  string              `json:"sql"`
	Stmt sqlparser.Statement `json:"stmt"`
}

func NewSqlResult(typ SqlType, strSQL string, stmt sqlparser.Statement) *SqlResult {
	if !typ.IsValid() {
		panic("sql type not valid")
	}
	return &SqlResult{
		Type: typ,
		SQL:  strSQL,
		Stmt: stmt,
	}
}
