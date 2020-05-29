package schema

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"os"
	"strings"
)

type Commander struct {
	ConnUrl        string
	Databases      []string
	Tables         []string
	Without        []string
	ReadOnly       []string
	Tags           []string
	Scheme         string
	Host           string
	User           string
	Password       string
	Charset        string
	OutDir         string
	Prefix         string
	Suffix         string
	PackageName    string
	Protobuf       bool
	DisableDecimal bool
}

type TableSchema struct {
	SchemeName   string        `json:"TABLE_SCHEMA" db:"TABLE_SCHEMA"`   //database name
	TableName    string        `json:"TABLE_NAME" db:"TABLE_NAME"`       //table name
	TableEngine  string        `json:"ENGINE" db:"ENGINE"`               //database engine
	TableComment string        `json:"TABLE_COMMENT" db:"TABLE_COMMENT"` //comment of table schema
	SchemeDir    string        `json:"SCHEMA_DIR" db:"SCHEMA_DIR"`       //output path
	PkName       string        `json:"PK_NAME" db:"PK_NAME"`             //primary key column name
	StructName   string        `json:"STRUCT_NAME" db:"STRUCT_NAME"`     //struct name
	OutDir       string        `json:"OUT_DIR" db:"OUT_DIR"`             //output directory
	FileName     string        `json:"FILE_NAME" db:"FILE_NAME"`         //output directory
	Columns      []TableColumn `json:"TABLE_COLUMNS" db:"TABLE_COLUMNS"` //columns with database and golang
}

type TableColumn struct {
	Name         string `json:"COLUMN_NAME" db:"COLUMN_NAME"`
	DataType     string `json:"DATA_TYPE" db:"DATA_TYPE"`
	Key          string `json:"COLUMN_KEY" db:"COLUMN_KEY"`
	Extra        string `json:"EXTRA" db:"EXTRA"`
	Comment      string `json:"COLUMN_COMMENT" db:"COLUMN_COMMENT"`
	IsPrimaryKey bool   // is primary key
	IsDecimal    bool   // is decimal type
	IsReadOnly   bool   // is read only
	GoName       string //column name in golang
	GoType       string //column type in golang
}

func IsInSlice(in string, s []string) bool {
	for _, v := range s {
		if v == in {
			return true
		}
	}
	return false
}

func MakeTags(strColName, strColType, strTagValue, strComment string, strAppends string) string {
	strComment = ReplaceCRLF(strComment)
	return fmt.Sprintf("	%v %v `json:\"%v\" db:\"%v\" %v` //%v \n",
		strColName, strColType, strTagValue, strTagValue, strAppends, strComment)
}

func MakeGetter(strStructName, strColName, strColType string) (strGetter string) {

	return fmt.Sprintf("func (do *%v) Get%v() %v { return do.%v } \n", strStructName, strColName, strColType, strColName)
}

func MakeSetter(strStructName, strColName, strColType string) (strSetter string) {

	return fmt.Sprintf("func (do *%v) Set%v(v %v) { do.%v = v } \n", strStructName, strColName, strColType, strColName)
}

func ReplaceCRLF(strIn string) (strOut string) {
	strOut = strings.ReplaceAll(strIn, "\r", "")
	strOut = strings.ReplaceAll(strOut, "\n", "")
	return
}

func CreateOutputFile(cmd *Commander, table *TableSchema, strFileSuffix string) (file *os.File, err error) {

	var strOutDir = cmd.OutDir
	var strPackageName = cmd.PackageName
	var strNamePrefix = cmd.Prefix
	var strNameSuffix = cmd.Suffix

	_, errStat := os.Stat(strOutDir)
	if errStat != nil && os.IsNotExist(errStat) {

		log.Info("mkdir [%v]", strOutDir)
		if err = os.Mkdir(strOutDir, os.ModeDir); err != nil {
			log.Error("mkdir [%v] error (%v)", strOutDir, err.Error())
			return
		}
	}

	table.OutDir = strOutDir

	if strPackageName == "" {
		//mkdir by output dir + scheme name
		strPackageName = table.SchemeName
		if strings.LastIndex(strOutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
			table.SchemeDir = fmt.Sprintf("%v/%v", strOutDir, strPackageName)
		} else {
			table.SchemeDir = fmt.Sprintf("%v%v", strOutDir, strPackageName)
		}
	} else {
		table.SchemeDir = fmt.Sprintf("%v/%v", strOutDir, strPackageName) //mkdir by package name
	}

	_, errStat = os.Stat(table.SchemeDir)

	if errStat != nil && os.IsNotExist(errStat) {

		log.Info("mkdir [%v]", table.SchemeDir)
		if err = os.Mkdir(table.SchemeDir, os.ModeDir); err != nil {
			log.Errorf("mkdir path name [%v] error (%v)", table.SchemeDir, err.Error())
			return
		}
	}

	if strNamePrefix != "" {
		strNamePrefix = fmt.Sprintf("%v_", strNamePrefix)
	}
	if strNameSuffix != "" {
		strNameSuffix = fmt.Sprintf("_%v", strNameSuffix)
	}

	table.FileName = fmt.Sprintf("%v/%v%v%v.%v", table.SchemeDir, strNamePrefix, table.TableName, strNameSuffix, strFileSuffix)

	file, err = os.OpenFile(table.FileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		log.Errorf("open file [%v] error (%v)", table.FileName, err.Error())
		return
	}
	log.Infof("create file [%v] ok", table.FileName)
	return
}
