package parser

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

type SqlParser struct {
	SQL  string              `json:"sql"`
	Stmt sqlparser.Statement `json:"stmt"`
}

//ParseMongo parse sql to mongodb bson result
func ParseMongo(strSQL string) (r *Result, err error) {
	p := newSqlParser(strSQL)
	return p.ParseMongo()
}

func newSqlParser(strSQL string) *SqlParser {
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Panic("SQL parse error [%s]", err.Error())
		return nil
	}
	return &SqlParser{
		SQL:  strSQL,
		Stmt: stmt,
	}
}

func (s *SqlParser) ParseMongo() (r *Result, err error) {
	var typ = StatementSqlType(s.Stmt)
	return newResult(typ, s.SQL, s.Stmt)
}
