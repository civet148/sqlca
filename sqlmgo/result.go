package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type Result struct {
	SqlType    SqlType             `json:"sql_type"`
	MogType    MgoType             `json:"mog_type"`
	SQL        string              `json:"sql"`
	Stmt       sqlparser.Statement `json:"stmt"`
	Match      bson.M              `json:"match"`
	Sort       bson.M              `json:"sort"`
	Set        bson.M              `json:"set"`
	Group      bson.M              `json:"group"`
	Projection bson.M              `json:"projection"`
}

func newResult(sqltype SqlType, strSQL string, stmt sqlparser.Statement) (r *Result, err error) {
	if !sqltype.IsValid() {
		log.Panic("sql type [%s] not valid, SQL [%s]", sqltype, strSQL)
	}

	r = &Result{
		SqlType:    sqltype,
		SQL:        strSQL,
		Stmt:       stmt,
		Match:      make(bson.M, 0),
		Sort:       make(bson.M, 0),
		Set:        make(bson.M, 0),
		Group:      make(bson.M, 0),
		Projection: make(bson.M, 0),
	}
	return r.walkSqlNode()
}

func (r *Result) walkSqlNode() (*Result, error) {

	err := sqlparser.Walk(func(node sqlparser.SQLNode) (ok bool, err error) {
		log.Json(node)
		return true, nil
	})
	return r, err
}

func (r *Result) printSqlNode() (*Result, error) {
	err := sqlparser.Walk(func(node sqlparser.SQLNode) (ok bool, err error) {
		log.Infof("--------------------------------------------------------------------------------------------")
		var strNodeType string
		switch node.(type) {
		case sqlparser.Comments:
			printSqlNodeComments(node)
		case sqlparser.Columns:
			printSqlNodeColumns(node)
		case sqlparser.TableExprs:
			printSqlNodeTableExprs(node)
		case sqlparser.SelectExprs:
			printSqlNodeSelectExprs(node)
		case sqlparser.TableNames:
			printSqlNodeTableNames(node)
		case sqlparser.GroupBy:
			printSqlNodeGroupBy(node)
		case sqlparser.OrderBy:
			printSqlNodeOrderBy(node)
		case sqlparser.ColIdent:
			printSqlNodeColIdent(node)
		case sqlparser.TableName:
			printSqlNodeTableName(node)
		case sqlparser.TableIdent:
			printSqlNodeTableIdent(node)
		case *sqlparser.StarExpr:
			printSqlNodeStarExpr(node)
		case *sqlparser.Limit:
			printSqlNodeLimit(node)
		case *sqlparser.Order:
			printSqlNodeOrder(node)
		case *sqlparser.Where:
			printSqlNodeWhere(node)
		case *sqlparser.Update:
			printSqlNodeUpdate(node)
		case *sqlparser.Select:
			printSqlNodeSelect(node)
		case *sqlparser.ParenSelect:
			printSqlNodeParenSelect(node)
		case *sqlparser.ParenExpr:
			printSqlNodeParenExpr(node)
		case *sqlparser.ParenTableExpr:
			printSqlNodeParenTableExpr(node)
		case *sqlparser.GroupConcatExpr:
			printSqlNodeGroupConcatExpr(node)
		case *sqlparser.SQLVal:
			printSqlNodeSQLVal(node)
		case *sqlparser.AliasedExpr:
			printSqlNodeAliasedExpr(node)
		case *sqlparser.AliasedTableExpr:
			printSqlNodeAliasedTableExpr(node)
		case *sqlparser.AndExpr:
			printSqlNodeAndExpr(node)
		case *sqlparser.BinaryExpr:
			printSqlNodeBinaryExpr(node)
		case *sqlparser.CollateExpr:
			printSqlNodeCollateExpr(node)
		case *sqlparser.ColName:
			printSqlNodeColName(node)
		case *sqlparser.Delete:
			printSqlNodeDelete(node)
		case *sqlparser.Insert:
			printSqlNodeInsert(node)
		case *sqlparser.IndexDefinition:
			printSqlNodeIndexDefinition(node)
		case *sqlparser.IndexHints:
			printSqlNodeIndexHints(node)
		case *sqlparser.IndexInfo:
			printSqlNodeIndexInfo(node)
		case *sqlparser.FuncExpr:
			printSqlNodeFuncExpr(node)
		case *sqlparser.Begin:
			printSqlNodeBegin(node)
		case sqlparser.BoolVal:
			printSqlNodeBoolVal(node)
		case *sqlparser.ComparisonExpr:
			printSqlNodeComparisonExpr(node)
		case *sqlparser.CaseExpr:
			printSqlNodeCaseExpr(node)
		case *sqlparser.When:
			printSqlNodeWhen(node)
		case *sqlparser.MatchExpr:
			printSqlNodeMatchExpr(node)
		case *sqlparser.ListArg:
			printSqlNodeListArg(node)
		case *sqlparser.Show:
			printSqlNodeShow(node)
		case *sqlparser.ShowFilter:
			printSqlNodeShowFilter(node)
		case *sqlparser.Union:
			printSqlNodeUnion(node)
		case *sqlparser.ColumnDefinition:
			printSqlNodeColumnDefinition(node)
		case *sqlparser.ColumnType:
			printSqlNodeColumnType(node)
		case *sqlparser.Commit:
			printSqlNodeCommit(node)
		case *sqlparser.ConvertExpr:
			printSqlNodeConvertExpr(node)
		case *sqlparser.ConvertType:
			printSqlNodeConvertType(node)
		case *sqlparser.ConvertUsingExpr:
			printSqlNodeConvertUsingExpr(node)
		case *sqlparser.ExistsExpr:
			printSqlNodeExistsExpr(node)
		case *sqlparser.DBDDL:
			printSqlNodeDBDDL(node)
		case *sqlparser.DDL:
			printSqlNodeDDL(node)
		case *sqlparser.IntervalExpr:
			printSqlNodeIntervalExpr(node)
		case *sqlparser.JoinCondition:
			printSqlNodeJoinCondition(node)
		case *sqlparser.JoinTableExpr:
			printSqlNodeJoinTableExpr(node)
		case *sqlparser.Set:
			printSqlNodeSet(node)
		case *sqlparser.SetExpr:
			printSqlNodeSetExpr(node)
		case *sqlparser.SetExprs:
			printSqlNodeSetExprs(node)
		case *sqlparser.Default:
			printSqlNodeDefault(node)
		default:
			strNodeType = reflect.TypeOf(node).Name()
			log.Errorf("unknown sql node type [%s]", strNodeType)
			log.Json(node)
		}
		return true, nil
	}, r.Stmt)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	return r, nil
}
