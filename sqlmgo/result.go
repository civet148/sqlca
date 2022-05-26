package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type Result struct {
	subType   subType             `json:"-"`
	SqlType   SqlType             `json:"sql_type"`
	MogType   MgoType             `json:"mog_type"`
	SQL       string              `json:"sql"`
	Stmt      sqlparser.Statement `json:"stmt"`
	TableName string              `json:"table_name"`
	Filter    bson.M              `json:"filter"`
	Options   bson.M              `json:"options"`
}

func newResult(sqltype SqlType, strSQL string, stmt sqlparser.Statement) (r *Result, err error) {
	if !sqltype.IsValid() {
		return nil, log.Errorf("sql type [%s] not support, SQL [%s]", sqltype, strSQL)
	}

	r = &Result{
		SqlType: sqltype,
		SQL:     strSQL,
		Stmt:    stmt,
		Filter:  make(bson.M, 0),
	}
	return r.formatSqlNode()
}

func (r *Result) formatSqlNode() (*Result, error) {
	r.MogType = MgoType_Other
	buf := sqlparser.NewTrackedBuffer(r.formatter)
	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type) {
		case sqlparser.GroupBy:
			groupby := node.(sqlparser.GroupBy)
			if len(groupby) != 0 {
				r.MogType = MgoType_GroupBy
			}
		case *sqlparser.Insert:
			r.MogType = MgoType_Insert
		case *sqlparser.Delete:
			r.MogType = MgoType_Delete
		case *sqlparser.Update:
			r.MogType = MgoType_Update
		case *sqlparser.Select:
			r.MogType = MgoType_Query
		}
		return true, nil
	}, r.Stmt)
	log.Debugf("mongo type [%v]", r.MogType.String())
	_ = buf
	buf.Myprintf("%v", r.Stmt)
	return r, nil
}

func (r *Result) formatter(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	switch node.(type) {
	case sqlparser.Comments:
		r.formatSqlNodeComments(buf, node.(sqlparser.Comments))
	case sqlparser.Columns:
		r.formatSqlNodeColumns(buf, node.(sqlparser.Columns))
	case sqlparser.TableExprs:
		r.formatSqlNodeTableExprs(buf, node.(sqlparser.TableExprs))
	case sqlparser.SelectExprs:
		r.formatSqlNodeSelectExprs(buf, node.(sqlparser.SelectExprs))
	case sqlparser.TableNames:
		r.formatSqlNodeTableNames(buf, node.(sqlparser.TableNames))
	case sqlparser.GroupBy:
		r.formatSqlNodeGroupBy(buf, node.(sqlparser.GroupBy))
	case sqlparser.OrderBy:
		r.formatSqlNodeOrderBy(buf, node.(sqlparser.OrderBy))
	case sqlparser.ColIdent:
		r.formatSqlNodeColIdent(buf, node.(sqlparser.ColIdent))
	case sqlparser.TableName:
		r.formatSqlNodeTableName(buf, node.(sqlparser.TableName))
	case sqlparser.TableIdent:
		r.formatSqlNodeTableIdent(buf, node.(sqlparser.TableIdent))
	case *sqlparser.StarExpr:
		r.formatSqlNodeStarExpr(buf, node.(*sqlparser.StarExpr))
	case *sqlparser.Limit:
		r.formatSqlNodeLimit(buf, node.(*sqlparser.Limit))
	case *sqlparser.Order:
		r.formatSqlNodeOrder(buf, node.(*sqlparser.Order))
	case *sqlparser.Where:
		r.formatSqlNodeWhere(buf, node.(*sqlparser.Where))
	case *sqlparser.Update:
		r.formatSqlNodeUpdate(buf, node.(*sqlparser.Update))
	case *sqlparser.Select:
		r.formatSqlNodeSelect(buf, node.(*sqlparser.Select))
	case *sqlparser.ParenSelect:
		r.formatSqlNodeParenSelect(buf, node.(*sqlparser.ParenSelect))
	case *sqlparser.ParenExpr:
		r.formatSqlNodeParenExpr(buf, node.(*sqlparser.ParenExpr))
	case *sqlparser.ParenTableExpr:
		r.formatSqlNodeParenTableExpr(buf, node.(*sqlparser.ParenTableExpr))
	case *sqlparser.GroupConcatExpr:
		r.formatSqlNodeGroupConcatExpr(buf, node.(*sqlparser.GroupConcatExpr))
	case *sqlparser.SQLVal:
		r.formatSqlNodeSQLVal(buf, node.(*sqlparser.SQLVal))
	case *sqlparser.AliasedExpr:
		r.formatSqlNodeAliasedExpr(buf, node.(*sqlparser.AliasedExpr))
	case *sqlparser.AliasedTableExpr:
		r.formatSqlNodeAliasedTableExpr(buf, node.(*sqlparser.AliasedTableExpr))
	case *sqlparser.AndExpr:
		r.formatSqlNodeAndExpr(buf, node.(*sqlparser.AndExpr))
	case *sqlparser.BinaryExpr:
		r.formatSqlNodeBinaryExpr(buf, node.(*sqlparser.BinaryExpr))
	case *sqlparser.CollateExpr:
		r.formatSqlNodeCollateExpr(buf, node.(*sqlparser.CollateExpr))
	case *sqlparser.ColName:
		r.formatSqlNodeColName(buf, node.(*sqlparser.ColName))
	case *sqlparser.Delete:
		r.formatSqlNodeDelete(buf, node.(*sqlparser.Delete))
	case *sqlparser.Insert:
		r.formatSqlNodeInsert(buf, node.(*sqlparser.Insert))
	case *sqlparser.IndexDefinition:
		r.formatSqlNodeIndexDefinition(buf, node.(*sqlparser.IndexDefinition))
	case *sqlparser.IndexHints:
		r.formatSqlNodeIndexHints(buf, node.(*sqlparser.IndexHints))
	case *sqlparser.IndexInfo:
		r.formatSqlNodeIndexInfo(buf, node.(*sqlparser.IndexInfo))
	case *sqlparser.FuncExpr:
		r.formatSqlNodeFuncExpr(buf, node.(*sqlparser.FuncExpr))
	case *sqlparser.Begin:
		r.formatSqlNodeBegin(buf, node.(*sqlparser.Begin))
	case sqlparser.BoolVal:
		r.formatSqlNodeBoolVal(buf, node.(sqlparser.BoolVal))
	case *sqlparser.ComparisonExpr:
		r.formatSqlNodeComparisonExpr(buf, node.(*sqlparser.ComparisonExpr))
	case *sqlparser.CaseExpr:
		r.formatSqlNodeCaseExpr(buf, node.(*sqlparser.CaseExpr))
	case *sqlparser.When:
		r.formatSqlNodeWhen(buf, node.(*sqlparser.When))
	case *sqlparser.MatchExpr:
		r.formatSqlNodeMatchExpr(buf, node.(*sqlparser.MatchExpr))
	case *sqlparser.ListArg:
		r.formatSqlNodeListArg(buf, node.(*sqlparser.ListArg))
	case *sqlparser.Show:
		r.formatSqlNodeShow(buf, node.(*sqlparser.Show))
	case *sqlparser.ShowFilter:
		r.formatSqlNodeShowFilter(buf, node.(*sqlparser.ShowFilter))
	case *sqlparser.Union:
		r.formatSqlNodeUnion(buf, node.(*sqlparser.Union))
	case *sqlparser.ColumnDefinition:
		r.formatSqlNodeColumnDefinition(buf, node.(*sqlparser.ColumnDefinition))
	case *sqlparser.ColumnType:
		r.formatSqlNodeColumnType(buf, node.(*sqlparser.ColumnType))
	case *sqlparser.Commit:
		r.formatSqlNodeCommit(buf, node.(*sqlparser.Commit))
	case *sqlparser.ConvertExpr:
		r.formatSqlNodeConvertExpr(buf, node.(*sqlparser.ConvertExpr))
	case *sqlparser.ConvertType:
		r.formatSqlNodeConvertType(buf, node.(*sqlparser.ConvertType))
	case *sqlparser.ConvertUsingExpr:
		r.formatSqlNodeConvertUsingExpr(buf, node.(*sqlparser.ConvertUsingExpr))
	case *sqlparser.ExistsExpr:
		r.formatSqlNodeExistsExpr(buf, node.(*sqlparser.ExistsExpr))
	case *sqlparser.DBDDL:
		r.formatSqlNodeDBDDL(buf, node.(*sqlparser.DBDDL))
	case *sqlparser.DDL:
		r.formatSqlNodeDDL(buf, node.(*sqlparser.DDL))
	case *sqlparser.IntervalExpr:
		r.formatSqlNodeIntervalExpr(buf, node.(*sqlparser.IntervalExpr))
	case *sqlparser.JoinCondition:
		r.formatSqlNodeJoinCondition(buf, node.(*sqlparser.JoinCondition))
	case *sqlparser.JoinTableExpr:
		r.formatSqlNodeJoinTableExpr(buf, node.(*sqlparser.JoinTableExpr))
	case *sqlparser.Set:
		r.formatSqlNodeSet(buf, node.(*sqlparser.Set))
	case *sqlparser.SetExpr:
		r.formatSqlNodeSetExpr(buf, node.(*sqlparser.SetExpr))
	case *sqlparser.SetExprs:
		r.formatSqlNodeSetExprs(buf, node.(*sqlparser.SetExprs))
	case *sqlparser.Default:
		r.formatSqlNodeDefault(buf, node.(*sqlparser.Default))
	default:
		log.Warnf("unknown sql node type [%s]", reflect.TypeOf(node).String())
		log.Json(node)
	}
}
