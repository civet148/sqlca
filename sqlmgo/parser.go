package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

//ParseSQL parse sql to mongodb bson result
func ParseSQL(strSQL string) (r *Result, err error) {
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Panic("SQL parse error [%s]", err.Error())
		return nil, err
	}
	typ := getSqlType(stmt)
	return newResult(typ, strSQL, stmt)
}
