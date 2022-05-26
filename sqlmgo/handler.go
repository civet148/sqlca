package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func (r *Result) handleSqlNodeSelect(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Select
	buf = sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("{%v}", node.SelectExprs)
	log.Infof("select [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeFrom(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_From
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.From)
	log.Infof("table [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeWhere(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Where
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.Where)
	log.Infof("where [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeGroupBy(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_GroupBy
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.GroupBy)
	log.Infof("group by [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeOrderBy(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_OrderBy
	buf = sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.OrderBy)
	log.Infof("order by [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeLimit(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Limit
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.Limit)
	log.Infof("limit [%s]", buf.String())
	return buf
}