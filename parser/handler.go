package parser

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func (r *Result) handleComments(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleColumns(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleTableExprs(node sqlparser.SQLNode) {
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
}

func (r *Result) handleSelectExprs(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleTableNames(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleGroupBy(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleOrderBy(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleColIdent(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleTableName(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleTableIdent(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleStarExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleLimit(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleOrder(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleWhere(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleUpdate(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleSelect(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleParenSelect(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleParenExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleParenTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleGroupConcatExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleSQLVal(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleAliasedExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleAliasedTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleAndExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleBinaryExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleCollateExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleColName(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleDelete(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleInsert(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleIndexDefinition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleIndexHints(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleIndexInfo(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleFuncExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleBegin(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleBoolVal(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleComparisonExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleCaseExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleWhen(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleMatchExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleListArg(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleShow(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleShowFilter(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleUnion(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleColumnDefinition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleColumnType(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleCommit(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleConvertExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleConvertType(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleConvertUsingExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleExistsExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleDBDDL(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleDDL(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleIntervalExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleJoinCondition(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleJoinTableExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleSet(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleSetExpr(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleSetExprs(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}

func (r *Result) handleDefault(node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v] unreachable", node)
}
