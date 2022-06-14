package sqlmgo

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func (r *Result) formatSqlNodeSelect(buf *sqlparser.TrackedBuffer, node *sqlparser.Select) {
	log.Json(node)
	//_ = r.handleSqlNodeSelectExprs(node)
	//_ = r.handleSqlNodeFrom(node)
	_ = r.handleSqlNodeWhere(node)
	//_ = r.handleSqlNodeGroupBy(node)
	//_ = r.handleSqlNodeOrderBy(node)
	//_ = r.handleSqlNodeLimit(node)
	//_ = r.handleSqlNodeHaving(node)
	//_ = r.handleSqlNodeDistinct(node)
}

func (r *Result) formatSqlNodeUpdate(buf *sqlparser.TrackedBuffer, node *sqlparser.Update) {
	log.Json(node)

}

func (r *Result) formatSqlNodeDelete(buf *sqlparser.TrackedBuffer, node *sqlparser.Delete) {
	log.Json(node)

}

func (r *Result) formatSqlNodeInsert(buf *sqlparser.TrackedBuffer, node *sqlparser.Insert) {
	log.Json(node)

}

func (r *Result) formatSqlNodeComments(buf *sqlparser.TrackedBuffer, node sqlparser.Comments) {
	log.Json(node)

}

func (r *Result) formatSqlNodeColumns(buf *sqlparser.TrackedBuffer, node sqlparser.Columns) {
	log.Json(node)

}

func (r *Result) formatSqlNodeTableExprs(buf *sqlparser.TrackedBuffer, node sqlparser.TableExprs) {
	log.Json(node)
	for _, expr := range node {
		buf.Myprintf("%v", expr)
	}
}

func (r *Result) formatSqlNodeSelectExprs(buf *sqlparser.TrackedBuffer, node sqlparser.SelectExprs) {
	log.Json(node)
	var prefix string
	for _, n := range node {
		buf.Myprintf("%s%v", prefix, n)
		prefix = ", "
	}
}

func (r *Result) formatSqlNodeTableNames(buf *sqlparser.TrackedBuffer, node sqlparser.TableNames) {
	log.Json(node)
	var prefix string
	for _, n := range node {
		buf.Myprintf("%v", prefix, n.Name)
		prefix = ", "
	}
}

func (r *Result) formatSqlNodeGroupBy(buf *sqlparser.TrackedBuffer, node sqlparser.GroupBy) {
	log.Json(node)

}

func (r *Result) formatSqlNodeOrderBy(buf *sqlparser.TrackedBuffer, node sqlparser.OrderBy) {
	log.Json(node)

}

func (r *Result) formatSqlNodeColIdent(buf *sqlparser.TrackedBuffer, node sqlparser.ColIdent) {
	log.Json(node)

}

func (r *Result) formatSqlNodeTableName(buf *sqlparser.TrackedBuffer, node sqlparser.TableName) {
	log.Json(node)
	buf.Myprintf("%v", node.Name)
}

func (r *Result) formatSqlNodeTableIdent(buf *sqlparser.TrackedBuffer, node sqlparser.TableIdent) {
	log.Json(node)
	buf.Myprintf("%s", node.String())
}

func (r *Result) formatSqlNodeStarExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.StarExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeLimit(buf *sqlparser.TrackedBuffer, node *sqlparser.Limit) {
	log.Json(node)

}

func (r *Result) formatSqlNodeOrder(buf *sqlparser.TrackedBuffer, node *sqlparser.Order) {
	log.Json(node)

}

func (r *Result) formatSqlNodeWhere(buf *sqlparser.TrackedBuffer, node *sqlparser.Where) {
	log.Json(node)
	buf.Myprintf("%v", node.Expr)
}

func (r *Result) formatSqlNodeParenSelect(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenSelect) {
	log.Json(node)

}

func (r *Result) formatSqlNodeParenExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenExpr) {
	log.Json(node)
	buf.Myprintf("%v", node.Expr)
}

func (r *Result) formatSqlNodeParenTableExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ParenTableExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeGroupConcatExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.GroupConcatExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeSQLVal(buf *sqlparser.TrackedBuffer, node *sqlparser.SQLVal) {
	log.Json(node)
	pv, err := sqlparser.NewPlanValue(node)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	v := pv.Value.ToString()
	switch node.Type {
	case sqlparser.StrVal:
		v = fmt.Sprintf("\"%s\"", v)
	//case sqlparser.IntVal:
	//case sqlparser.FloatVal:
	//case sqlparser.HexNum:
	//case sqlparser.HexVal:
	//case sqlparser.ValArg:
	//case sqlparser.BitVal:
	}
	buf.Myprintf("%s", v)
}

func (r *Result) formatSqlNodeAliasedExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AliasedExpr) {
	log.Json(node)
	buf.Myprintf("%v", node.Expr)
}

func (r *Result) formatSqlNodeAliasedTableExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AliasedTableExpr) {
	log.Json(node)
	buf.Myprintf("%v", node.Expr)
}

func (r *Result) formatSqlNodeAndExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.AndExpr) {
	log.Json(node)
	buf.Myprintf("%v", node.Left)
	buf.Myprintf(",")
	buf.Myprintf("%v", node.Right)
}

func (r *Result) formatSqlNodeOrExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.OrExpr) {
	log.Json(node)
	buf.Myprintf("\"$or\":[")
	buf.Myprintf("{%v}", node.Left)
	buf.Myprintf(",")
	buf.Myprintf("{%v}", node.Right)
	buf.Myprintf("]")
}

func (r *Result) formatSqlNodeBinaryExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.BinaryExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeCollateExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.CollateExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeColName(buf *sqlparser.TrackedBuffer, node *sqlparser.ColName) {
	log.Json(node)
	switch r.subType {
	case subType_SelectExprs:
		{
			qualifier := node.Qualifier.Qualifier.String()
			name := node.Qualifier.Name.String()
			if qualifier != "" || name != "" {
				if qualifier != "" {
					name = qualifier + "." + name
				}
				buf.Myprintf("\"%s.%s\":1", name, node.Name.String())
			} else {
				buf.Myprintf("\"%s\":1", node.Name.String())
			}
		}
	case subType_From:
	case subType_Where:
		{
			buf.Myprintf("\"%s\"", node.Name.String())
		}
	case subType_GroupBy:
		{
			buf.Myprintf("\"$%s\":\"$%s\"", node.Name.String(), node.Name.String())
		}

	case subType_OrderBy:
	case subType_Update:
	case subType_Delete:
	case subType_Insert:
	case subType_Limit:
	case subType_Having:
	}
}

func (r *Result) formatSqlNodeIndexDefinition(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexDefinition) {
	log.Json(node)

}

func (r *Result) formatSqlNodeIndexHints(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexHints) {
	log.Json(node)

}

func (r *Result) formatSqlNodeIndexInfo(buf *sqlparser.TrackedBuffer, node *sqlparser.IndexInfo) {
	log.Json(node)
}

func (r *Result) formatSqlNodeFuncExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.FuncExpr) {
	log.Json(node)
}

func (r *Result) formatSqlNodeBegin(buf *sqlparser.TrackedBuffer, node *sqlparser.Begin) {
	log.Json(node)

}

func (r *Result) formatSqlNodeBoolVal(buf *sqlparser.TrackedBuffer, node sqlparser.BoolVal) {
	log.Json(node)
	buf.Myprintf("%s", fmt.Sprintf("%v", node))
}

func (r *Result) formatSqlNodeComparisonExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ComparisonExpr) {
	log.Json(node)
	buf.Myprintf("%v", node.Left)
	buf.Myprintf(":{\"%s\"", ConvertOperator(node.Operator))
	buf.Myprintf(":%v}", node.Right)
}

func (r *Result) formatSqlNodeCaseExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.CaseExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeWhen(buf *sqlparser.TrackedBuffer, node *sqlparser.When) {
	log.Json(node)

}

func (r *Result) formatSqlNodeMatchExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.MatchExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeListArg(buf *sqlparser.TrackedBuffer, node *sqlparser.ListArg) {
	log.Json(node)

}

func (r *Result) formatSqlNodeShow(buf *sqlparser.TrackedBuffer, node *sqlparser.Show) {
	log.Json(node)

}

func (r *Result) formatSqlNodeShowFilter(buf *sqlparser.TrackedBuffer, node *sqlparser.ShowFilter) {
	log.Json(node)

}

func (r *Result) formatSqlNodeUnion(buf *sqlparser.TrackedBuffer, node *sqlparser.Union) {
	log.Json(node)

}

func (r *Result) formatSqlNodeColumnDefinition(buf *sqlparser.TrackedBuffer, node *sqlparser.ColumnDefinition) {
	log.Json(node)

}

func (r *Result) formatSqlNodeColumnType(buf *sqlparser.TrackedBuffer, node *sqlparser.ColumnType) {
	log.Json(node)

}

func (r *Result) formatSqlNodeCommit(buf *sqlparser.TrackedBuffer, node *sqlparser.Commit) {
	log.Json(node)

}

func (r *Result) formatSqlNodeConvertExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeConvertType(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertType) {
	log.Json(node)

}

func (r *Result) formatSqlNodeConvertUsingExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ConvertUsingExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeExistsExpr(buf *sqlparser.TrackedBuffer, node *sqlparser.ExistsExpr) {
	log.Json(node)

}

func (r *Result) formatSqlNodeDBDDL(buf *sqlparser.TrackedBuffer, node *sqlparser.DBDDL) {
	log.Json(node)

}

func (r *Result) formatSqlNodeDDL(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeIntervalExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeJoinCondition(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeJoinTableExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeSet(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeSetExpr(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeSetExprs(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}

func (r *Result) formatSqlNodeDefault(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	log.Json(node)

}
