package main

import (
	"flag"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/consts"
	_ "github.com/civet148/sqlca/cmd/db2go/mssql"
	_ "github.com/civet148/sqlca/cmd/db2go/mysql"
	_ "github.com/civet148/sqlca/cmd/db2go/postgres"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"github.com/civet148/sqlca/cmd/db2go/structs"
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
var argvWithout = flag.String("without", "", "exclude columns")
var argvReadOnly = flag.String("readonly", "", "read only columns")
var argvProtobuf = flag.Bool("proto", false, "output proto buffer file")
var argvDisableDecimal = flag.Bool("disable-decimal", true, "compatible with legacy version [deprecated]")
var argvEnableDecimal = flag.Bool("enable-decimal", false, "decimal as sqlca.Decimal type")
var argvGogoOptions = flag.String("gogo-options", "", "gogo proto options")
var argvOneFile = flag.Bool("one-file", false, "output go/proto file into one file which named by database name")
var argvOrm = flag.Bool("orm", false, "generate ORM code inner data object")
var argvOmitEmpty = flag.Bool("omitempty", false, "omit empty for json tag")
var argvStruct = flag.Bool("struct", false, "generate struct getter and setter")
var argvConst = flag.Bool("const", false, "generate const variants format methods")

func init() {
	flag.Parse()
	log.SetLevel("info")
}

func main() {

	//var err error
	var cmd = schema.Commander{}
	cmd.Prefix = *argvPackage
	cmd.Prefix = *argvPrefix
	cmd.Suffix = *argvSuffix
	cmd.OutDir = *argvOutput
	cmd.ConnUrl = *argvUrl
	cmd.PackageName = *argvPackage
	cmd.Protobuf = *argvProtobuf
	cmd.EnableDecimal = *argvEnableDecimal
	cmd.Orm = *argvOrm
	cmd.OmitEmpty = *argvOmitEmpty

	if *argvStruct {
		structs.ExportStruct(&cmd)
	} else if *argvConst {
		consts.ExportConst(&cmd)
	} else {

		if *argvUrl == "" {
			log.Infof("")
			fmt.Println("need --url parameter")
			flag.Usage()
			return
		}

		ui := sqlca.ParseUrl(*argvUrl)

		log.Infof("%+v", cmd.String())

		if *argvDatabase == "" {
			//use default database
			cmd.Database = schema.GetDatabaseName(ui.Path)
		} else {
			//use input database
			cmd.Database = strings.TrimSpace(*argvDatabase)
		}

		if *argvTables != "" {
			cmd.Tables = schema.TrimSpaceSlice(strings.Split(*argvTables, ","))
		}

		if *argvWithout != "" {
			cmd.Without = schema.TrimSpaceSlice(strings.Split(*argvWithout, ","))
		}

		if *argvProtobuf {
			if *argvGogoOptions != "" {
				cmd.GogoOptions = schema.TrimSpaceSlice(strings.Split(*argvGogoOptions, ","))
				if len(cmd.GogoOptions) == 0 {
					cmd.GogoOptions = schema.TrimSpaceSlice(strings.Split(*argvGogoOptions, ";"))
				}
			}
		}

		if *argvOneFile {
			cmd.OneFile = true
		}

		if *argvTags != "" {
			cmd.Tags = schema.TrimSpaceSlice(strings.Split(*argvTags, ","))
		}
		if *argvReadOnly != "" {
			cmd.ReadOnly = schema.TrimSpaceSlice(strings.Split(*argvReadOnly, ","))
		}

		cmd.Scheme = ui.Scheme
		cmd.Host = ui.Host
		cmd.User = ui.User
		cmd.Password = ui.Password
		e := sqlca.NewEngine(false)
		e.Open(cmd.ConnUrl)
		export(&cmd, e)
	}
}

func export(cmd *schema.Commander, e *sqlca.Engine) {
	e.Debug(true)
	exporter := schema.NewExporter(cmd, e)
	if exporter == nil {
		log.Errorf("new exporter error, nil object")
		return
	}
	if cmd.Protobuf {
		if err := exporter.ExportProto(); err != nil {
			log.Errorf("export [%v] to protobuf file error [%v]", cmd.Scheme, err.Error())
		}
	} else {
		if err := exporter.ExportGo(); err != nil {
			log.Errorf("export [%v] to go file error [%v]", cmd.Scheme, err.Error())
		}
	}
}
