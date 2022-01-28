package main

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func main() {
	TestSqlParse()
}

func TestSqlParse() {
	strSQL := "SELECT miner, date, SUM(to_decimal(win_reward)) as total_reward, COUNT(1) AS total_count " +
		" FROM miner_reward  WHERE miner='0x45a36a8e118c37e4c47ef4ab827a7c9e579e11e2' AND date >= '2021-12-01' AND date <= '2022-01-31' " +
		" GROUP BY miner, date"

	stmt, err := sqlparser.Parse(strSQL)
	if err != nil {
		log.Errorf("SQL parse error [%s]", err.Error())
		return
	}
	log.Json(stmt)

	switch stmt.(type) {
	case *sqlparser.Select:
		log.Infof("Select")
	case *sqlparser.Insert:
		log.Infof("Insert")
	case *sqlparser.Update:
		log.Infof("Update")
	case *sqlparser.Delete:
		log.Infof("Delete")
	case *sqlparser.Union:
		log.Infof("Union")
	case *sqlparser.Begin:
		log.Infof("Begin")
	case *sqlparser.Rollback:
		log.Infof("Rollback")
	case *sqlparser.Commit:
		log.Infof("Commit")
	case *sqlparser.Set:
		log.Infof("Set")
	case *sqlparser.DDL:
		log.Infof("DDL")
	case *sqlparser.DBDDL:
		log.Infof("DBDDL")
	case *sqlparser.Use:
		log.Infof("Use")
	case *sqlparser.Show:
		log.Infof("Show")
	case *sqlparser.OtherRead:
		log.Infof("OtherRead")
	case *sqlparser.OtherAdmin:
		log.Infof("OtherAdmin")
	case *sqlparser.ParenSelect:
		log.Infof("ParenSelect")
	case *sqlparser.Stream:
		log.Infof("Stream")
	default:
		log.Infof("Unknown")
	}
}
