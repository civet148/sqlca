package main

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/sqlmgo"
	"github.com/xwb1989/sqlparser"
)

const (
	strSQL1 = "SELECT age, sex FROM user WHERE id=1 and name='tom'"
	strSQL2 = "SELECT `miner`, `date`, SUM(to_decimal(miner_reward.win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM `miner_reward` mr  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND (date >= '2021-12-01' AND date <= '2022-01-31') and ok=true" +
		" GROUP BY miner, date ORDER by date DESC LIMIT 2,4"
)

func main() {
	parse(strSQL1)
}

func parse(strSQL string) {
	log.Infof("SQL [%s]", strSQL)
	r, err := sqlmgo.ParseSQL(strSQL)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Json(r)
}

func Formatter(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	//log.Infof("buffer [%s] node [%#v]", buf.String(), node)
	//log.Json("SQL NODE", node)
}
