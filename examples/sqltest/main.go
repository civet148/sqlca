package main

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func main() {
	strSQL0 := "SELECT `miner`, `date`, SUM(to_decimal(miner_reward.win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM `miner_reward` mr  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND (date >= '2021-12-01' AND date <= '2022-01-31') and ok=true" +
		" GROUP BY miner, date ORDER by date DESC"

	parse(strSQL0)
}

func parse(strSQL string) {

	log.Infof("SQL [%s]", strSQL)
	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Errorf("SQL parse error [%s]", err.Error())
		return
	}
	log.Json(stmt)
	//sqlparser.Walk(func(node sqlparser.SQLNode) (ok bool, err error) {
	//	log.Infof("--------------------------------------------------------------------------------------------")
	//	if node != nil {
	//		log.Infof("sql node [%#v]", node)
	//		log.Json(node)
	//	}
	//	return true, nil
	//}, stmt)
	buf := sqlparser.NewTrackedBuffer(Formatter)
	buf.Myprintf("%v", stmt)
	//log.Infof("SQL [%s]", buf.String())
	//parser.ParseMongo(strSQL)
}

func Formatter(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Infof("buffer [%s] node [%#v]", buf.String(), node)
}
