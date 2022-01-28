package parser

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

type sqlParser struct {
	SQL  string              `json:"sql"`
	Stmt sqlparser.Statement `json:"stmt"`
}

//ParseMongo parse sql to mongodb bson result
func ParseMongo(strSQL string) (r *MgoResult, err error) {
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Panic("SQL parse error [%s]", err.Error())
		return
	}
	p := &sqlParser{
		SQL:  strSQL,
		Stmt: stmt,
	}
	return p.ParseMongo()
}

//ParseSQL parse sql to relational db result
func ParseSQL(strSQL string) (r *SqlResult, err error) {
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Panic("SQL parse error [%s]", err.Error())
		return
	}
	p := &sqlParser{
		SQL:  strSQL,
		Stmt: stmt,
	}
	return p.ParseSQL()
}


func (s *sqlParser) ParseMongo() (r *MgoResult, err error) {
	var typ = StatementSqlType(s.Stmt)
	r = NewMgoResult(typ)
	return
}

func (s *sqlParser) ParseSQL() (r *SqlResult, err error) {
	var typ = StatementSqlType(s.Stmt)
	r = NewSqlResult(typ, s.SQL, s.Stmt)
	return
}
