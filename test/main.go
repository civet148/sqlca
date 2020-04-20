package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca/test/mssql"
	"github.com/civet148/sqlca/test/mysql"
)

func main() {

	mysql.Benchmark()
	log.Infof("------------------------------------------------------------------------------------------------------------------")
	mssql.Benchmark()

	//log.Info("%+v", log.Report()) //print function report
	log.Info("program exit...")
}
