package mysql

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"os"
)

func ExportProtobuf(cmd *schema.Commander, e *sqlca.Engine) (err error) {

	var schemas []*schema.TableSchema
	if schemas, err = queryTables(cmd, e); err != nil {
		log.Errorf(err.Error())
		return
	}

	var file *os.File
	strHead := makeProtoHead(cmd)
	for i, v := range schemas {
		if err = queryTableColumns(cmd, e, v); err != nil {
			log.Error(err.Error())
			return
		}

		var append bool
		if i > 0 && cmd.OneFile {
			append = true
		}

		strBody := makeProtoBody(cmd, v)

		if file, err = schema.CreateOutputFile(cmd, v, "proto", append); err != nil {
			log.Error(err.Error())
			return
		}

		if i == 0 {
			file.WriteString(strHead)
		} else if !cmd.OneFile {
			file.WriteString(strHead)
		}
		file.WriteString(strBody)
	}
	file.Close()
	return
}

func makeProtoHead(cmd *schema.Commander) (strContent string) {

	strContent += "syntax = \"proto3\";\n"
	strContent += fmt.Sprintf("package %v;\n\n", cmd.PackageName)

	if len(cmd.GogoOptions) > 0 {
		strContent += schema.IMPORT_GOGO_PROTO + "\n"
		//strContent += schema.IMPORT_GOOGOLE_PROTOBUF + "\n"
	}
	strContent += "\n"
	for _, v := range cmd.GogoOptions {
		strContent += fmt.Sprintf("option %v;\n", v)
	}
	strContent += "\n"
	return
}

func makeProtoBody(cmd *schema.Commander, table *schema.TableSchema) (strContent string) {

	strTableName := schema.CamelCaseConvert(table.TableName)
	strContent += fmt.Sprintf("message %vDO {\n", strTableName)
	for i, v := range table.Columns {

		if schema.IsInSlice(v.Name, cmd.Without) {
			continue
		}
		no := i + 1
		strColName := v.Name
		strColType := getProtoColumnType(table.TableName, v.Name, v.DataType, v.Key, v.Extra, true)
		strContent += fmt.Sprintf("	%-10s %-22s = %-2d; //%v\n", strColType, strColName, no, v.Comment)
	}
	strContent += "}\n\n"
	return
}
