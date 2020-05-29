package mysql

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"strings"
)

func queryTables(cmd *schema.Commander, e *sqlca.Engine) (schemas []*schema.TableSchema, err error) {

	var strQuery string

	var dbs, tables []string

	for _, v := range cmd.Databases {
		dbs = append(dbs, fmt.Sprintf("'%v'", v))
	}

	if len(dbs) == 0 {
		return nil, fmt.Errorf("no database selected")
	}

	log.Infof("ready to export tables [%v]", cmd.Tables)
	for _, v := range cmd.Tables {
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	if len(tables) == 0 {
		strQuery = fmt.Sprintf("SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `ENGINE`, `TABLE_COMMENT` FROM `INFORMATION_SCHEMA`.`TABLES` "+
			"WHERE (`ENGINE`='MyISAM' OR `ENGINE` = 'InnoDB' OR `ENGINE` = 'TokuDB') AND `TABLE_SCHEMA` IN (%v) ORDER BY TABLE_SCHEMA",
			strings.Join(dbs, ","))
	} else {
		strQuery = fmt.Sprintf("SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `ENGINE`, `TABLE_COMMENT` FROM `INFORMATION_SCHEMA`.`TABLES` "+
			"WHERE (`ENGINE`='MyISAM' OR `ENGINE` = 'InnoDB' OR `ENGINE` = 'TokuDB') AND `TABLE_SCHEMA` IN (%v) AND TABLE_NAME IN (%v) ORDER BY TABLE_SCHEMA",
			strings.Join(dbs, ","), strings.Join(tables, ","))
	}

	_, err = e.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func queryTableColumns(cmd *schema.Commander, e *sqlca.Engine, table *schema.TableSchema) (err error) {

	/*
	 SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS`
	 WHERE `TABLE_SCHEMA` = 'accounts' AND `TABLE_NAME` = 'users' ORDER BY ORDINAL_POSITION ASC
	*/
	_, err = e.Model(&table.Columns).QueryRaw("SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS` "+
		"WHERE `TABLE_SCHEMA` = '%v' AND `TABLE_NAME` = '%v' ORDER BY ORDINAL_POSITION ASC", table.SchemeName, table.TableName)
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

//将数据库字段类型转为go语言对应的数据类型
func getGoColumnType(strTableName, strColName, strDataType, strColKey, strExtra string, disableDecimal bool) (strColType string, isDecimal bool) {

	switch strDataType {
	case "bigint":
		strColType = "int64"
	case "int", "integer", "mediumint":
		strColType = "int32"
	case "smallint":
		strColType = "int16"
	case "tinyint", "bit":
		strColType = "int8"
	case "bool", "boolean":
		strColType = "bool"
	case "decimal":
		if disableDecimal {
			strColType = "float64"
		} else {
			strColType = "sqlca.Decimal"
		}
		isDecimal = true
	case "real", "double", "float", "numeric":
		strColType = "float64"
	case "datetime", "year", "date", "time", "timestamp":
		strColType = "string"
	case "enum", "set", "varchar", "char", "text", "tinytext", "mediumtext", "longtext":
		strColType = "string"
	case "blob", "tinyblob", "mediumblob", "longblob", "binary", "varbinary", "json":
		strColType = "string"
	default:
		{
			err := fmt.Errorf("table [%v] column [%v] data type [%v] unsupport", strTableName, strColName, strDataType)
			log.Errorf("%v", err.Error())
			panic(err.Error())
		}
	}
	return
}

//将数据库字段类型转为protobuf对应的数据类型
func getProtoColumnType(strTableName, strColName, strDataType, strColKey, strExtra string, disableDecimal bool) (strColType string) {

	switch strDataType {
	case "bigint":
		strColType = "int64"
	case "int", "integer", "mediumint":
		strColType = "int32"
	case "smallint":
		strColType = "int32"
	case "tinyint", "bit":
		strColType = "int32"
	case "bool", "boolean":
		strColType = "bool"
	case "double", "decimal":
		if disableDecimal {
			strColType = "double"
		}
	case "real", "float", "numeric":
		strColType = "float"
	case "datetime", "year", "date", "time", "timestamp":
		strColType = "string"
	case "enum", "set", "varchar", "char", "text", "tinytext", "mediumtext", "longtext":
		strColType = "string"
	case "blob", "tinyblob", "mediumblob", "longblob", "binary", "varbinary", "json":
		strColType = "string"
	default:
		{
			err := fmt.Errorf("table [%v] column [%v] data type [%v] unsupport", strTableName, strColName, strDataType)
			log.Errorf("%v", err.Error())
			panic(err.Error())
		}
	}
	return
}
