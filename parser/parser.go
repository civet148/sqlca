package parser

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func SqlParse(strSQL string) {
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Errorf("SQL parse error [%s]", err.Error())
		return
	}
	log.Json(stmt)
}