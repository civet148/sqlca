package parser

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

type MgoType int

const (
	MgoType_Other       MgoType = 0
	MgoType_QuerySingle MgoType = 1
	MgoType_QueryGroup  MgoType = 2
	MgoType_QueryUnion  MgoType = 3
	MgoType_Insert      MgoType = 4
	MgoType_Update      MgoType = 5
	MgoType_Delete      MgoType = 6
)

var mgotypes = map[MgoType]string{
	MgoType_Other:       "MgoType_Other",
	MgoType_QuerySingle: "MgoType_QuerySingle",
	MgoType_QueryGroup:  "MgoType_QueryGroup",
	MgoType_QueryUnion:  "MgoType_QueryUnion",
	MgoType_Insert:      "MgoType_Insert",
	MgoType_Update:      "MgoType_Update",
	MgoType_Delete:      "MgoType_Delete",
}

func (t MgoType) GoString() string {
	return t.String()
}

func (t MgoType) String() string {
	if strName, ok := mgotypes[t]; ok {
		return strName
	}
	return fmt.Sprintf("MgoType_Unknown<%d>", t)
}

/*----------------------------------------------------------------------------------------------------*/

type sqlType int

const (
	sqlType_Other       sqlType = 0
	sqlType_Select      sqlType = 1
	sqlType_Insert      sqlType = 2
	sqlType_Update      sqlType = 3
	sqlType_Delete      sqlType = 4
	sqlType_Union       sqlType = 5
	sqlType_Begin       sqlType = 6
	sqlType_Rollback    sqlType = 7
	sqlType_Commit      sqlType = 8
	sqlType_Set         sqlType = 9
	sqlType_DDL         sqlType = 10
	sqlType_DBDDL       sqlType = 11
	sqlType_Use         sqlType = 12
	sqlType_Show        sqlType = 13
	sqlType_OtherRead   sqlType = 14
	sqlType_OtherAdmin  sqlType = 15
	sqlType_ParenSelect sqlType = 16
	sqlType_Stream      sqlType = 17
)

var sqltypes = map[sqlType]string{
	sqlType_Other:       "sqlType_Other",
	sqlType_Select:      "sqlType_Select",
	sqlType_Insert:      "sqlType_Insert",
	sqlType_Update:      "sqlType_Update",
	sqlType_Delete:      "sqlType_Delete",
	sqlType_Union:       "sqlType_Union",
	sqlType_Begin:       "sqlType_Begin",
	sqlType_Rollback:    "sqlType_Rollback",
	sqlType_Commit:      "sqlType_Commit",
	sqlType_Set:         "sqlType_Set",
	sqlType_DDL:         "sqlType_DDL",
	sqlType_DBDDL:       "sqlType_DBDDL",
	sqlType_Use:         "sqlType_Use",
	sqlType_Show:        "sqlType_Show",
	sqlType_OtherRead:   "sqlType_OtherRead",
	sqlType_OtherAdmin:  "sqlType_OtherAdmin",
	sqlType_ParenSelect: "sqlType_ParenSelect",
	sqlType_Stream:      "sqlType_Stream",
}

func (t sqlType) IsValid() bool {
	if _, ok := sqltypes[t]; ok {
		return true
	}
	return false
}

func (t sqlType) GoString() string {
	return t.String()
}

func (t sqlType) String() string {
	if strName, ok := sqltypes[t]; ok {
		return strName
	}
	return fmt.Sprintf("sqlType_Unknown<%d>", t)
}

func StatementSqlType(stmt sqlparser.Statement) (typ sqlType) {

	switch stmt.(type) {
	case *sqlparser.Select:
		typ = sqlType_Select
	case *sqlparser.Insert:
		typ = sqlType_Insert
	case *sqlparser.Update:
		typ = sqlType_Update
	case *sqlparser.Delete:
		typ = sqlType_Delete
	case *sqlparser.Union:
		typ = sqlType_Union
	case *sqlparser.Begin:
		typ = sqlType_Begin
	case *sqlparser.Rollback:
		typ = sqlType_Rollback
	case *sqlparser.Commit:
		typ = sqlType_Commit
	case *sqlparser.Set:
		typ = sqlType_Set
	case *sqlparser.DDL:
		typ = sqlType_DDL
	case *sqlparser.DBDDL:
		typ = sqlType_DBDDL
	case *sqlparser.Use:
		typ = sqlType_Use
	case *sqlparser.Show:
		typ = sqlType_Show
	case *sqlparser.OtherRead:
		typ = sqlType_OtherRead
	case *sqlparser.OtherAdmin:
		typ = sqlType_OtherAdmin
	case *sqlparser.ParenSelect:
		typ = sqlType_ParenSelect
	case *sqlparser.Stream:
		typ = sqlType_Stream
	default:
		typ = sqlType_Other
	}
	return
}
