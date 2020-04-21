package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca/test/mssql"
)

func main() {
	log.Infof("------------------------------------------------------MYSQL------------------------------------------------------------")
	//mysql.Benchmark()
	log.Infof("----------------------------------------------------SQLSERVER----------------------------------------------------------")
	mssql.Benchmark()

	//log.Info("%+v", log.Report()) //print function report
	log.Info("program exit...")
}
