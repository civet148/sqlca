package mysql

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"os"
	"strings"
)

type ExporterMysql struct {
	Cmd     *schema.Commander
	Engine  *sqlca.Engine
	Schemas []*schema.TableSchema
}

func init() {
	schema.Register(schema.SCHEME_MYSQL, NewExporterMysql)
}

func NewExporterMysql(cmd *schema.Commander, e *sqlca.Engine) schema.Exporter {

	return &ExporterMysql{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterMysql) ExportGo() (err error) {
	var cmd = m.Cmd
	var e = m.Engine
	var schemas = m.Schemas
	//var tableNames []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	//var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)
	log.Infof("ready to export tables [%v]", cmd.Tables)

	if schemas, err = m.queryTableSchemas(cmd, e); err != nil {
		log.Errorf("query tables error [%s]", err.Error())
		return
	}
	for _, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			log.Error(err.Error())
			return
		}
	}

	return schema.ExportTableSchema(cmd, e, schemas)
}

func (m *ExporterMysql) ExportProto() (err error) {
	var cmd = m.Cmd
	var e = m.Engine
	var schemas = m.Schemas
	if schemas, err = m.queryTableSchemas(cmd, e); err != nil {
		log.Errorf(err.Error())
		return
	}

	var file *os.File
	strHead := schema.MakeProtoHead(cmd)
	for i, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			log.Error(err.Error())
			return
		}

		var append bool
		if i > 0 && cmd.OneFile {
			append = true
		}

		strBody := schema.MakeProtoBody(cmd, v)

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

func (m *ExporterMysql) queryTableSchemas(cmd *schema.Commander, e *sqlca.Engine) (schemas []*schema.TableSchema, err error) {

	var strQuery string
	var tables []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)

	log.Infof("ready to export tables [%v]", cmd.Tables)
	for _, v := range cmd.Tables {
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	if len(tables) == 0 {
		strQuery = fmt.Sprintf("select `table_schema` as table_schema, `table_name` as table_name, `engine` as engine, `table_comment` as table_comment from `information_schema`.`tables` "+
			"where (`engine`='myisam' or `engine` = 'innodb' or `engine` = 'tokudb') and `table_schema` in (%v) order by table_schema",
			strDatabaseName)
	} else {
		strQuery = fmt.Sprintf("select `table_schema` as table_schema, `table_name` as table_name, `engine` as engine, `table_comment` as table_comment from `information_schema`.`tables` "+
			"where (`engine`='myisam' or `engine` = 'innodb' or `engine` = 'tokudb') and `table_schema` in (%v) and table_name in (%v) order by table_schema",
			strDatabaseName, strings.Join(tables, ","))
	}

	_, err = e.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func (m *ExporterMysql) queryTableColumns(table *schema.TableSchema) (err error) {

	/*
	 SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS`
	 WHERE `TABLE_SCHEMA` = 'accounts' AND `TABLE_NAME` = 'users' ORDER BY ORDINAL_POSITION ASC
	*/
	var e = m.Engine
	_, err = e.Model(&table.Columns).QueryRaw("select `table_name` as table_name, `column_name` as column_name, `data_type` as data_type, `extra` as extra, `column_key` as column_key, `column_comment` as column_comment"+
		" from `information_schema`.`columns` where `table_schema` = '%v' and `table_name` = '%v' order by ordinal_position asc", table.SchemeName, table.TableName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	//write table name in camel case naming
	table.TableComment = schema.ReplaceCRLF(table.TableComment)
	for i, v := range table.Columns {
		table.Columns[i].Comment = schema.ReplaceCRLF(v.Comment)
	}
	return
}
