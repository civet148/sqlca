package sqlmgo

import (
	"github.com/civet148/log"
	"testing"
)

func TestParseMongo(t *testing.T) {
	strSQL := "SELECT `miner`, `date`, SUM(to_decimal(miner_reward.win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM `miner_reward` mr  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND (date >= '2021-12-01' AND date <= '2022-01-31') and ok=true" +
		" GROUP BY miner, date ORDER by date DESC"

	_, err := ParseSQL(strSQL)
	if err != nil {
		log.Errorf(err.Error())
	}
}
