package main

import (
	"flag"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/mysql"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"strings"
)

var argvUrl = flag.String("url", "", "mysql://root:123456@127.0.0.1:3306/test?charset=utf8")
var argvOutput = flag.String("out", ".", "output directory, default .")
var argvDatabase = flag.String("db", "", "export databases, like 'test,chat_db'")
var argvTables = flag.String("table", "", "export tables, eg. 'users,devices'")
var argvTags = flag.String("tag", "", "golang struct tag name, default json,db")
var argvPrefix = flag.String("prefix", "", "export file prefix")
var argvSuffix = flag.String("suffix", "", "export file suffix")
var argvPackage = flag.String("package", "", "export package name")

func main() {

	//var err error
	var si schema.SchemaInfo
	log.Infof("argument: url [%v]", *argvUrl)
	log.Infof("argument: databases [%v]", *argvDatabase)
	log.Infof("argument: output [%v]", *argvOutput)
	log.Infof("argument: tag [%v]", *argvTags)
	log.Infof("argument: tables [%v]", *argvTables)
	log.Infof("argument: prefix [%v]", *argvPrefix)
	log.Infof("argument: suffix [%v]", *argvSuffix)
	log.Infof("argument: package name [%v]", *argvPackage)

	if *argvUrl == "" {
		fmt.Println("need --url parameter")
		flag.Usage()
		return
	}

	si.Tags = *argvTags
	si.Prefix = *argvPackage
	si.Prefix = *argvPrefix
	si.Suffix = *argvSuffix
	si.OutDir = *argvOutput
	si.ConnUrl = *argvUrl
	si.PackageName = *argvPackage

	ui := sqlca.ParseUrl(*argvUrl)

	if *argvDatabase == "" {
		//use default database
		si.Databases = append(si.Databases, getDatabaseName(ui.Path))
	} else {
		//use input databases
		si.Databases = strings.Split(*argvDatabase, ",")
	}

	if *argvTables != "" {
		si.Tables = strings.Split(*argvTables, ",")
	}

	si.Scheme = ui.Scheme
	si.Host = ui.Host
	si.User = ui.User
	si.Password = ui.Password

	switch si.Scheme {
	case "mysql":
		exportMysql(&si)
	case "postgres":
		exportPostgres(&si)
	}
}

func init() {
	flag.Parse()
}

func getDatabaseName(strPath string) (strName string) {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		return
	}
	return strPath[idx+1:]
}

func exportMysql(si *schema.SchemaInfo) {
	if err := mysql.Export(si); err != nil {
		log.Errorf("export mysql schema error [%v]", err.Error())
	}
}

func exportPostgres(si *schema.SchemaInfo) {

}
