package parser

import "testing"

func TestSqlParse(t *testing.T) {
	strSQL := "SELECT miner, date, SUM(to_decimal(win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM miner_reward  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND date >= '2021-12-01' AND date <= '2022-01-31' " +
		" GROUP BY miner, date"

	SqlParse(strSQL)
}