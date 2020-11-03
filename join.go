package sqlca

type JoinType int

const (
	JoinType_Inner = 0 //inner join
	JoinType_Left  = 1 //left join
	JoinType_Right = 2 //right join
)

func (t JoinType) GoString() string {
	return t.String()
}

func (t JoinType) String() string {
	switch t {
	case JoinType_Inner:
		return "JoinType_Inner"
	case JoinType_Left:
		return "JoinType_Left"
	case JoinType_Right:
		return "JoinType_Right"
	}
	return "JoinType_Unknown"
}

func (t JoinType) ToKeyWord() string {
	switch t {
	case JoinType_Inner:
		return DATABASE_KEY_NAME_INNER_JOIN
	case JoinType_Left:
		return DATABASE_KEY_NAME_LEFT_JOIN
	case JoinType_Right:
		return DATABASE_KEY_NAME_RIGHT_JOIN
	}
	return "<nil>"
}

type Join struct {
	e            *Engine
	jt           JoinType
	strTableName string
	strOn        string
}

func (j *Join) On(strOn string, args ...interface{}) *Engine {
	e := j.e
	j.strOn = e.formatString(strOn, args...)
	e.joins = append(e.joins, j)
	return e
}
