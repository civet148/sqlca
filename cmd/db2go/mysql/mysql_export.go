package mysql

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"os"
	"strings"
)

/*
-- 查询数据库表名、引擎及注释
SELECT `TABLE_SCHEMA`, `TABLE_NAME`, `ENGINE`, `TABLE_COMMENT` FROM `INFORMATION_SCHEMA`.`TABLES`
WHERE `TABLE_SCHEMA`='accounts' AND (`ENGINE`='MyISAM' OR `ENGINE` = 'InnoDB' OR `ENGINE` = 'TokuDB')

-- 查询数据表字段名、字段类型及注释
SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`,  `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS`
WHERE `TABLE_SCHEMA` = 'accounts' AND `TABLE_NAME` = 'acc_3pl'
*/

type TableSchema struct {
	SchemeName   string `json:"TABLE_SCHEMA" db:"TABLE_SCHEMA"`
	TableName    string `json:"TABLE_NAME" db:"TABLE_NAME"`
	TableEngine  string `json:"ENGINE" db:"ENGINE"`
	TableComment string `json:"TABLE_COMMENT" db:"TABLE_COMMENT"`
	SchemeDir    string `json:"SCHEMA_DIR" db:"SCHEMA_DIR"`
}

type TableColumn struct {
	TableName     string `json:"TABLE_NAME" db:"TABLE_NAME"`
	ColumnName    string `json:"COLUMN_NAME" db:"COLUMN_NAME"`
	DataType      string `json:"DATA_TYPE" db:"DATA_TYPE"`
	ColumnKey     string `json:"COLUMN_KEY" db:"COLUMN_KEY"`
	Extra         string `json:"EXTRA" db:"EXTRA"`
	ColumnComment string `json:"COLUMN_COMMENT" db:"COLUMN_COMMENT"`
}

type TableColumnGo struct {
	SchemeName string
	TableName  string
	ColumnName string
	ColumnType string
}

func Export(si *schema.SchemaInfo) (err error) {

	e := sqlca.NewEngine(false)
	e.Debug(true)
	e.Open(si.ConnUrl)
	var strQuery string
	var tableSchemas []TableSchema

	var dbs, tables []string

	for _, v := range si.Databases {
		dbs = append(dbs, fmt.Sprintf("'%v'", v))
	}

	if len(dbs) == 0 {
		return fmt.Errorf("no database selected")
	}
	log.Warnf("tables [%v]", si.Tables)
	for _, v := range si.Tables {
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

	_, err = e.Model(&tableSchemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	return exportTableSchema(si, e, tableSchemas)
}

func exportTableSchema(si *schema.SchemaInfo, e *sqlca.Engine, tables []TableSchema) (err error) {

	for _, v := range tables {

		_, errStat := os.Stat(si.OutDir)
		if errStat != nil && os.IsNotExist(errStat) {

			log.Info("mkdir [%v]", si.OutDir)
			if err = os.Mkdir(si.OutDir, os.ModeDir); err != nil {
				log.Error("mkdir [%v] error (%v)", si.OutDir, err.Error())
				return
			}
		}

		if si.PackageName == "" {
			//mkdir by output dir + scheme name
			si.PackageName = v.SchemeName
			if strings.LastIndex(si.OutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
				v.SchemeDir = fmt.Sprintf("%v/%v", si.OutDir, si.PackageName)
			} else {
				v.SchemeDir = fmt.Sprintf("%v%v", si.OutDir, si.PackageName)
			}
		} else {
			v.SchemeDir = fmt.Sprintf("%v/%v", si.OutDir, si.PackageName) //mkdir by package name
		}

		_, errStat = os.Stat(v.SchemeDir)

		if errStat != nil && os.IsNotExist(errStat) {

			log.Info("mkdir [%v]", v.SchemeDir)
			if err = os.Mkdir(v.SchemeDir, os.ModeDir); err != nil {
				log.Errorf("mkdir path name [%v] error (%v)", v.SchemeDir, err.Error())
				return
			}
		}

		var strPrefix, strSuffix string
		if si.Prefix != "" {
			strPrefix = fmt.Sprintf("%v_", si.Prefix)
		}
		if si.Suffix != "" {
			strSuffix = fmt.Sprintf("_%v", si.Suffix)
		}
		strFileName := fmt.Sprintf("%v/%v%v%v.go", v.SchemeDir, strPrefix, v.TableName, strSuffix)
		if err = exportTableColumns(si, e, v, strFileName); err != nil {
			return
		}
	}

	return
}

func isInSlice(in string, s []string) bool {
	for _, v := range s {
		if v == in {
			return true
		}
	}
	return false
}

func exportTableColumns(si *schema.SchemaInfo, e *sqlca.Engine, table TableSchema, strFileName string) (err error) {

	File, err := os.OpenFile(strFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		log.Errorf("open file [%v] error (%v)", strFileName, err.Error())
		return
	}
	log.Infof("exporting table schema [%v] to file [%v]", table.TableName, strFileName)

	var strContent string

	//write package name
	strContent += fmt.Sprintf("package %v\n\n", si.PackageName)

	var TableCols []TableColumn
	var TableColsGo []TableColumnGo

	/*
	 SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS`
	 WHERE `TABLE_SCHEMA` = 'accounts' AND `TABLE_NAME` = 'users'
	*/
	e.Model(&TableCols).QueryRaw("SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT` FROM `INFORMATION_SCHEMA`.`COLUMNS` "+
		"WHERE `TABLE_SCHEMA` = '%v' AND `TABLE_NAME` = '%v'", table.SchemeName, table.TableName)

	//write table name in camel case naming
	strTableName := camelCaseConvert(table.TableName)
	strContent += fmt.Sprintf("var TableName%v = \"%v\" //%v \n\n", strTableName, table.TableName, table.TableComment)

	strContent += fmt.Sprintf("type %vDO struct { \n", strTableName)
	for _, v := range TableCols {

		if isInSlice(v.ColumnName, si.Without) {
			continue
		}

		var tagValues []string
		strColName := camelCaseConvert(v.ColumnName)
		strColType := getColumnType(v.TableName, v.ColumnName, v.DataType, v.ColumnKey, v.Extra)
		if si.Tags != "" {
			tags := strings.Split(si.Tags, ",")
			for _, t := range tags {
				tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t, v.ColumnName))
			}
		}
		strContent += fmt.Sprintf("	%v %v `json:\"%v\" db:\"%v\" %v` //%v \n",
			strColName, strColType, v.ColumnName, v.ColumnName, strings.Join(tagValues, " "), v.ColumnComment)

		var colGo TableColumnGo
		colGo.SchemeName = table.SchemeName
		colGo.TableName = v.TableName
		colGo.ColumnName = strColName
		colGo.ColumnType = strColType
		TableColsGo = append(TableColsGo, colGo)
	}

	strContent += "}\n"
	_, _ = File.WriteString(strContent)
	return
}

//将数据库字段类型转为go语言对应的数据类型
func getColumnType(strTableName, strColName, strDataType, strColKey, strExtra string) (strColType string) {

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
	case "real", "decimal", "double", "float":
		strColType = "float64"
	case "datetime", "year", "date", "time", "timestamp":
		strColType = "string"
	case "enum", "set", "varchar", "char", "text", "tinytext", "mediumtext", "longtext":
		strColType = "string"
	case "blob", "tinyblob", "mediumblob", "longblob", "binary", "varbinary":
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

func camelCaseConvert(strIn string) (strOut string) {

	var idxUnderLine = int(-1)
	for i, v := range strIn {
		strChr := string(v)

		if i == 0 {

			strOut += strings.ToUpper(strChr)
		} else {
			if v == '_' {
				idxUnderLine = i //ignore
			} else {

				if i == idxUnderLine+1 {

					strOut += strings.ToUpper(strChr)
				} else {
					strOut += strChr
				}
			}
		}
	}

	return
}
