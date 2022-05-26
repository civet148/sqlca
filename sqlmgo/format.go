package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func (r *Result) formatSqlNodeComments(buf *sqlparser.TrackedBuffer, node sqlparser.Comments) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeColumns(buf *sqlparser.TrackedBuffer, node sqlparser.Columns) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeTableExprs(buf *sqlparser.TrackedBuffer, node sqlparser.TableExprs) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSelectExprs(buf *sqlparser.TrackedBuffer, node sqlparser.SelectExprs) {
	log.Debugf("sql node [%#v]", node)
	buf.Myprintf("%v", node)
}

func (r *Result) formatSqlNodeTableNames(buf *sqlparser.TrackedBuffer, node sqlparser.TableNames) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeGroupBy(buf *sqlparser.TrackedBuffer, node sqlparser.GroupBy) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeOrderBy(buf *sqlparser.TrackedBuffer, node sqlparser.OrderBy) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeColIdent(buf *sqlparser.TrackedBuffer, node sqlparser.ColIdent) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeTableName(buf *sqlparser.TrackedBuffer, node sqlparser.TableName) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeTableIdent(buf *sqlparser.TrackedBuffer, node sqlparser.TableIdent) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeStarExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.StarExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeLimit(buf *sqlparser.TrackedBuffer, node *sqlparser.Limit) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeOrder(buf *sqlparser.TrackedBuffer, node *sqlparser.Order) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeWhere(buf *sqlparser.TrackedBuffer, node *sqlparser.Where) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeUpdate(buf *sqlparser.TrackedBuffer, node *sqlparser.Update) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSelect(buf *sqlparser.TrackedBuffer, node *sqlparser.Select) {
	log.Debugf("sql node [%#v]", node)
}

func (r *Result) formatSqlNodeParenSelect(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenSelect) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeParenExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeParenTableExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenTableExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeGroupConcatExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.GroupConcatExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSQLVal(buf *sqlparser.TrackedBuffer, node *sqlparser.SQLVal) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeAliasedExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AliasedExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeAliasedTableExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AliasedTableExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeAndExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AndExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeBinaryExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.BinaryExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeCollateExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.CollateExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeColName(buf *sqlparser.TrackedBuffer, node *sqlparser.ColName) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeDelete(buf *sqlparser.TrackedBuffer, node *sqlparser.Delete) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeInsert(buf *sqlparser.TrackedBuffer, node *sqlparser.Insert) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeIndexDefinition(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexDefinition) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeIndexHints(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexHints) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeIndexInfo(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexInfo) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeFuncExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.FuncExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeBegin(buf *sqlparser.TrackedBuffer, node *sqlparser.Begin) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeBoolVal(buf *sqlparser.TrackedBuffer, node sqlparser.BoolVal) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeComparisonExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ComparisonExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeCaseExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.CaseExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeWhen(buf *sqlparser.TrackedBuffer, node *sqlparser.When) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeMatchExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.MatchExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeListArg(buf *sqlparser.TrackedBuffer, node *sqlparser.ListArg) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeShow(buf *sqlparser.TrackedBuffer, node *sqlparser.Show) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeShowFilter(buf *sqlparser.TrackedBuffer, node *sqlparser.ShowFilter) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeUnion(buf *sqlparser.TrackedBuffer, node *sqlparser.Union) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeColumnDefinition(buf *sqlparser.TrackedBuffer, node *sqlparser.ColumnDefinition) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeColumnType(buf *sqlparser.TrackedBuffer, node *sqlparser.ColumnType) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeCommit(buf *sqlparser.TrackedBuffer, node *sqlparser.Commit) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeConvertExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeConvertType(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertType) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeConvertUsingExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertUsingExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeExistsExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ExistsExpr) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeDBDDL(buf *sqlparser.TrackedBuffer, node *sqlparser.DBDDL) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeDDL(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeIntervalExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeJoinCondition(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeJoinTableExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSet(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSetExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeSetExprs(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}

func (r *Result) formatSqlNodeDefault(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Debugf("sql node [%#v]", node)
	
}
