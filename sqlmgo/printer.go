package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func printSqlNodeComments(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeColumns(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeTableExprs(node sqlparser.SQLNode) {
	//exprs := node.(sqlparser.TableExprs)
	//for _, expr := range exprs {
	//	t := expr.(*sqlparser.AliasedTableExpr)
	//	e := t.Expr.(sqlparser.TableName)
	//	r.Tables = append(r.Tables, Table{
	//		Name: e.Name.String(),
	//		As:   t.As.String(),
	//	})
	//}
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSelectExprs(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeTableNames(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeGroupBy(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeOrderBy(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeColIdent(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeTableName(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeTableIdent(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeStarExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeLimit(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeOrder(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeWhere(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeUpdate(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSelect(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeParenSelect(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeParenExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeParenTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeGroupConcatExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSQLVal(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeAliasedExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeAliasedTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeAndExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeBinaryExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeCollateExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeColName(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeDelete(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeInsert(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeIndexDefinition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeIndexHints(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeIndexInfo(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeFuncExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeBegin(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeBoolVal(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeComparisonExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeCaseExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeWhen(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeMatchExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeListArg(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeShow(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeShowFilter(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeUnion(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeColumnDefinition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeColumnType(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeCommit(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeConvertExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeConvertType(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeConvertUsingExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeExistsExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeDBDDL(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeDDL(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeIntervalExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeJoinCondition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeJoinTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSet(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSetExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeSetExprs(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}

func printSqlNodeDefault(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
	log.Json(node)
}
