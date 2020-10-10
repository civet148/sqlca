package schema

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/sqlca"
	"os"
	"strings"
)

const (
	SCHEME_MYSQL    = "mysql"
	SCHEME_POSTGRES = "postgres"
	SCHEME_MSSQL    = "mssql"
)

const (
	IMPORT_GOGO_PROTO = `import "github.com/gogo/protobuf/gogoproto/gogo.proto";`
	IMPORT_SQLCA      = `import "github.com/civet148/sqlca"`
)

type Commander struct {
	ConnUrl       string
	Database      string
	Tables        []string
	Without       []string
	ReadOnly      []string
	Tags          []string
	Scheme        string
	Host          string
	User          string
	Password      string
	Charset       string
	OutDir        string
	Prefix        string
	Suffix        string
	PackageName   string
	Protobuf      bool
	EnableDecimal bool
	OneFile       bool
	GogoOptions   []string
	Orm           bool
	OmitEmpty     bool
}

type TableSchema struct {
	SchemeName         string        `json:"table_schema" db:"table_schema"`   //database name
	TableName          string        `json:"table_name" db:"table_name"`       //table name
	TableEngine        string        `json:"engine" db:"engine"`               //database engine
	TableComment       string        `json:"table_comment" db:"table_comment"` //comment of table schema
	SchemeDir          string        `json:"schema_dir" db:"schema_dir"`       //output path
	PkName             string        `json:"pk_name" db:"pk_name"`             //primary key column name
	StructName         string        `json:"struct_name" db:"struct_name"`     //struct name
	OutDir             string        `json:"out_dir" db:"out_dir"`             //output directory
	FileName           string        `json:"file_name" db:"file_name"`         //output directory
	Columns            []TableColumn `json:"table_columns" db:"table_columns"` //columns with database and golang
	TableNameCamelCase string        //table name in camel case
}

type TableColumn struct {
	Name         string `json:"column_name" db:"column_name"`
	DataType     string `json:"data_type" db:"data_type"`
	Key          string `json:"column_key" db:"column_key"`
	Extra        string `json:"extra" db:"extra"`
	Comment      string `json:"column_comment" db:"column_comment"`
	IsPrimaryKey bool   // is primary key
	IsDecimal    bool   // is decimal type
	IsReadOnly   bool   // is read only
	GoName       string //column name in golang
	GoType       string //column type in golang
}

type Exporter interface {
	ExportGo() (err error)
	ExportProto() (err error)
}

type Instance func(cmd *Commander, e *sqlca.Engine) Exporter

var instances = make(map[string]Instance, 1)

func Register(strScheme string, inst Instance) {
	instances[strScheme] = inst
}

func NewExporter(cmd *Commander, e *sqlca.Engine) Exporter {
	var ok bool
	var inst Instance
	if inst, ok = instances[cmd.Scheme]; !ok {
		log.Errorf("scheme [%v] instance not registered", cmd.Scheme)
		return nil
	}
	return inst(cmd, e)
}

func IsInSlice(in string, s []string) bool {
	for _, v := range s {
		if v == in {
			return true
		}
	}
	return false
}

func MakeTags(strColName, strColType, strTagValue, strComment string, strAppends string, omitEmpty bool) string {
	strComment = ReplaceCRLF(strComment)
	var strJsonValue string
	strJsonValue = strTagValue
	if omitEmpty {
		strJsonValue += ",omitempty"
	}
	return fmt.Sprintf("	%v %v `json:\"%v\" db:\"%v\" %v` //%v \n",
		strColName, strColType, strJsonValue, strTagValue, strAppends, strComment)
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

func CreateOutputFile(cmd *Commander, table *TableSchema, strFileSuffix string, append bool) (file *os.File, err error) {

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

	var flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC

	if append {
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	}

	if cmd.OneFile { //数据库名称作为文件名
		table.FileName = fmt.Sprintf("%v/%v%v%v.%v", table.SchemeDir, strNamePrefix, table.SchemeName, strNameSuffix, strFileSuffix)
	} else { //数据表名作为文件名
		table.FileName = fmt.Sprintf("%v/%v%v%v.%v", table.SchemeDir, strNamePrefix, table.TableName, strNameSuffix, strFileSuffix)
	}

	file, err = os.OpenFile(table.FileName, flag, 0)
	if err != nil {
		log.Errorf("open file [%v] error (%v)", table.FileName, err.Error())
		return
	}
	log.Infof("open file [%v] ok", table.FileName)
	return
}

func CamelCaseConvert(strIn string) (strOut string) {

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

func TrimSpaceSlice(s []string) (ts []string) {
	for _, v := range s {
		ts = append(ts, strings.TrimSpace(v))
	}
	return
}

func GetDatabaseName(strPath string) (strName string) {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		return
	}
	return strPath[idx+1:]
}

//将数据库字段类型转为go语言对应的数据类型
func GetGoColumnType(strTableName, strColName, strDataType string, enableDecimal bool) (strColType string, isDecimal bool) {

	var ok bool
	if strColType, ok = db2goTypes[strDataType]; !ok {
		strColType = "string"
		log.Warnf("table [%v] column [%v] data type [%v] not support yet, set as string type", strTableName, strColName, strDataType)
		return
	}
	switch strDataType {
	case DB_COLUMN_TYPE_DECIMAL:
		if !enableDecimal {
			strColType = "float64"
		} else {
			strColType = "sqlca.Decimal"
		}
	}
	return
}

//将数据库字段类型转为protobuf对应的数据类型
func GetProtoColumnType(strTableName, strColName, strDataType string) (strColType string) {

	var ok bool
	if strColType, ok = db2protoTypes[strDataType]; !ok {
		strColType = "string"
		log.Warnf("table [%v] column [%v] data type [%v] not support yet, set as string type", strTableName, strColName, strDataType)
		return
	}
	return
}

func HandleCommentCRLF(table *TableSchema) {
	//write table name in camel case naming
	table.TableComment = ReplaceCRLF(table.TableComment)
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}
}
