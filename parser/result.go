package parser

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type Result struct {
	SqlType    sqlType             `json:"sql_type"`
	MogType    MgoType             `json:"mog_type"`
	SQL        string              `json:"sql"`
	Stmt       sqlparser.Statement `json:"stmt"`
	Match      bson.M              `json:"match"`
	Sort       bson.M              `json:"sort"`
	Set        bson.M              `json:"set"`
	Group      bson.M              `json:"group"`
	Projection bson.M              `json:"projection"`
}

func newResult(sqltype sqlType, strSQL string, stmt sqlparser.Statement) (r *Result, err error) {
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
		log.Infof("--------------------------------------------------------------------------------------------")
		var strNodeType string
		switch node.(type) {
		case sqlparser.Comments:
			r.handleComments(node)
		case sqlparser.Columns:
			r.handleColumns(node)
		case sqlparser.TableExprs:
			r.handleTableExprs(node)
		case sqlparser.SelectExprs:
			r.handleSelectExprs(node)
		case sqlparser.TableNames:
			r.handleTableNames(node)
		case sqlparser.GroupBy:
			r.handleGroupBy(node)
		case sqlparser.OrderBy:
			r.handleOrderBy(node)
		case sqlparser.ColIdent:
			r.handleColIdent(node)
		case sqlparser.TableName:
			r.handleTableName(node)
		case sqlparser.TableIdent:
			r.handleTableIdent(node)
		case *sqlparser.StarExpr:
			r.handleStarExpr(node)
		case *sqlparser.Limit:
			r.handleLimit(node)
		case *sqlparser.Order:
			r.handleOrder(node)
		case *sqlparser.Where:
			r.handleWhere(node)
		case *sqlparser.Update:
			r.handleUpdate(node)
		case *sqlparser.Select:
			r.handleSelect(node)
		case *sqlparser.ParenSelect:
			r.handleParenSelect(node)
		case *sqlparser.ParenExpr:
			r.handleParenExpr(node)
		case *sqlparser.ParenTableExpr:
			r.handleParenTableExpr(node)
		case *sqlparser.GroupConcatExpr:
			r.handleGroupConcatExpr(node)
		case *sqlparser.SQLVal:
			r.handleSQLVal(node)
		case *sqlparser.AliasedExpr:
			r.handleAliasedExpr(node)
		case *sqlparser.AliasedTableExpr:
			r.handleAliasedTableExpr(node)
		case *sqlparser.AndExpr:
			r.handleAndExpr(node)
		case *sqlparser.BinaryExpr:
			r.handleBinaryExpr(node)
		case *sqlparser.CollateExpr:
			r.handleCollateExpr(node)
		case *sqlparser.ColName:
			r.handleColName(node)
		case *sqlparser.Delete:
			r.handleDelete(node)
		case *sqlparser.Insert:
			r.handleInsert(node)
		case *sqlparser.IndexDefinition:
			r.handleIndexDefinition(node)
		case *sqlparser.IndexHints:
			r.handleIndexHints(node)
		case *sqlparser.IndexInfo:
			r.handleIndexInfo(node)
		case *sqlparser.FuncExpr:
			r.handleFuncExpr(node)
		case *sqlparser.Begin:
			r.handleBegin(node)
		case sqlparser.BoolVal:
			r.handleBoolVal(node)
		case *sqlparser.ComparisonExpr:
			r.handleComparisonExpr(node)
		case *sqlparser.CaseExpr:
			r.handleCaseExpr(node)
		case *sqlparser.When:
			r.handleWhen(node)
		case *sqlparser.MatchExpr:
			r.handleMatchExpr(node)
		case *sqlparser.ListArg:
			r.handleListArg(node)
		case *sqlparser.Show:
			r.handleShow(node)
		case *sqlparser.ShowFilter:
			r.handleShowFilter(node)
		case *sqlparser.Union:
			r.handleUnion(node)
		case *sqlparser.ColumnDefinition:
			r.handleColumnDefinition(node)
		case *sqlparser.ColumnType:
			r.handleColumnType(node)
		case *sqlparser.Commit:
			r.handleCommit(node)
		case *sqlparser.ConvertExpr:
			r.handleConvertExpr(node)
		case *sqlparser.ConvertType:
			r.handleConvertType(node)
		case *sqlparser.ConvertUsingExpr:
			r.handleConvertUsingExpr(node)
		case *sqlparser.ExistsExpr:
			r.handleExistsExpr(node)
		case *sqlparser.DBDDL:
			r.handleDBDDL(node)
		case *sqlparser.DDL:
			r.handleDDL(node)
		case *sqlparser.IntervalExpr:
			r.handleIntervalExpr(node)
		case *sqlparser.JoinCondition:
			r.handleJoinCondition(node)
		case *sqlparser.JoinTableExpr:
			r.handleJoinTableExpr(node)
		case *sqlparser.Set:
			r.handleSet(node)
		case *sqlparser.SetExpr:
			r.handleSetExpr(node)
		case *sqlparser.SetExprs:
			r.handleSetExprs(node)
		case *sqlparser.Default:
			r.handleDefault(node)
		default:
			strNodeType = reflect.TypeOf(node).Name()
			log.Errorf("unknown sql node type [%s]", strNodeType)
		}
		log.Json(node)
		return true, nil
	}, r.Stmt)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	return r, nil
}
