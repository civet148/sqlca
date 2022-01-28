package parser

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

type SqlType int

const (
	SqlType_Other       SqlType = 0
	SqlType_Select      SqlType = 1
	SqlType_Insert      SqlType = 2
	SqlType_Update      SqlType = 3
	SqlType_Delete      SqlType = 4
	SqlType_Union       SqlType = 5
	SqlType_Begin       SqlType = 6
	SqlType_Rollback    SqlType = 7
	SqlType_Commit      SqlType = 8
	SqlType_Set         SqlType = 9
	SqlType_DDL         SqlType = 10
	SqlType_DBDDL       SqlType = 11
	SqlType_Use         SqlType = 12
	SqlType_Show        SqlType = 13
	SqlType_OtherRead   SqlType = 14
	SqlType_OtherAdmin  SqlType = 15
	SqlType_ParenSelect SqlType = 16
	SqlType_Stream      SqlType = 17
)

var sqltypes = map[SqlType]string{
	SqlType_Other:       "SqlType_Other",
	SqlType_Select:      "SqlType_Select",
	SqlType_Insert:      "SqlType_Insert",
	SqlType_Update:      "SqlType_Update",
	SqlType_Delete:      "SqlType_Delete",
	SqlType_Union:       "SqlType_Union",
	SqlType_Begin:       "SqlType_Begin",
	SqlType_Rollback:    "SqlType_Rollback",
	SqlType_Commit:      "SqlType_Commit",
	SqlType_Set:         "SqlType_Set",
	SqlType_DDL:         "SqlType_DDL",
	SqlType_DBDDL:       "SqlType_DBDDL",
	SqlType_Use:         "SqlType_Use",
	SqlType_Show:        "SqlType_Show",
	SqlType_OtherRead:   "SqlType_OtherRead",
	SqlType_OtherAdmin:  "SqlType_OtherAdmin",
	SqlType_ParenSelect: "SqlType_ParenSelect",
	SqlType_Stream:      "SqlType_Stream",
}

func (t SqlType) IsValid() bool {
	if _, ok := sqltypes[t]; ok {
		return true
	}
	return false
}

func (t SqlType) GoString() string {
	return t.String()
}

func (t SqlType) String() string {
	if strName, ok := sqltypes[t]; ok {
		return strName
	}
	return fmt.Sprintf("SqlType_Unknown<%d>", t)
}

func StatementSqlType(stmt sqlparser.Statement) (typ SqlType) {

	switch stmt.(type) {
	case *sqlparser.Select:
		typ = SqlType_Select
	case *sqlparser.Insert:
		typ = SqlType_Insert
	case *sqlparser.Update:
		typ = SqlType_Update
	case *sqlparser.Delete:
		typ = SqlType_Delete
	case *sqlparser.Union:
		typ = SqlType_Union
	case *sqlparser.Begin:
		typ = SqlType_Begin
	case *sqlparser.Rollback:
		typ = SqlType_Rollback
	case *sqlparser.Commit:
		typ = SqlType_Commit
	case *sqlparser.Set:
		typ = SqlType_Set
	case *sqlparser.DDL:
		typ = SqlType_DDL
	case *sqlparser.DBDDL:
		typ = SqlType_DBDDL
	case *sqlparser.Use:
		typ = SqlType_Use
	case *sqlparser.Show:
		typ = SqlType_Show
	case *sqlparser.OtherRead:
		typ = SqlType_OtherRead
	case *sqlparser.OtherAdmin:
		typ = SqlType_OtherAdmin
	case *sqlparser.ParenSelect:
		typ = SqlType_ParenSelect
	case *sqlparser.Stream:
		typ = SqlType_Stream
	default:
		typ = SqlType_Other
	}
	return
}