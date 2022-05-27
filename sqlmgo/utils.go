package sqlmgo

func ConvertOperator(strOperator string) string {
	switch strOperator {
	case "=":
		return "$eq"
	case ">":
		return "$gt"
	case "<":
		return "$lt"
	case ">=":
		return "$gte"
	case "<=":
		return "$lte"
	case "!=":
		return "$ne"
	}
	return strOperator
}