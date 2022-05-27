package sqlmgo

import (
	"github.com/civet148/log"
	"github.com/xwb1989/sqlparser"
)

func (r *Result) handleSqlNodeSelectExprs(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_SelectExprs
	buf = sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("{%v}", node.SelectExprs)
	log.Infof("SelectExprs [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeFrom(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_From
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.From)
	log.Infof("From [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeWhere(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Where
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("{%v}", node.Where)
	log.Infof("Where [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeGroupBy(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_GroupBy
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.GroupBy)
	log.Infof("GroupBy [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeOrderBy(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_OrderBy
	buf = sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.OrderBy)
	log.Infof("OrderBy [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeLimit(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Limit
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.Limit)
	log.Infof("Limit [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeHaving(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	r.subType = subType_Having
	buf= sqlparser.NewTrackedBuffer(r.formatter)
	buf.Myprintf("%v", node.Having)
	log.Infof("Having [%s]", buf.String())
	return buf
}

func (r *Result) handleSqlNodeDistinct(node *sqlparser.Select) (buf *sqlparser.TrackedBuffer) {
	log.Warnf("TODO...")
	return nil
}