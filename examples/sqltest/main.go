package main

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/sqlmgo"
)

const (
	strSelect = "SELECT `age`, sex, balance, data.weight.kg FROM user WHERE sex>=10 and weight != 0 and created_time > '2022-05-22 10:00:00' and ok=true ORDER BY created_time, id DESC"
	strGourpBy = "SELECT `miner`, `date`, SUM(to_decimal(miner_reward.win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM `miner_reward` mr  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND (date >= '2021-12-01' OR date <= '2022-01-31') and deleted=0 and ok=true" +
		" GROUP BY miner, date ORDER by date, created_time DESC LIMIT 2,4"
	strInsert = "INSERT INTO user(user_id, name, sex, age, created_time) VALUES(1005, 'lory', 'male', 28, '2006-01-02 15:04:05'),(1006, 'kitty', 'female', 20, '2006-01-02 15:04:05')"
)

func main() {
	//parse(strSelect)
	parse(strGourpBy)
	//parse(strInsert)
}

func parse(strSQL string) {
	log.Infof("SQL [%s]", strSQL)
	r, err := sqlmgo.ParseSQL(strSQL)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	_ = r
	//log.Json(r)
}

