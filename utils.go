package pgqb

// FlipSortingMode flips.. sorting mode
func FlipSortingMode(str string) string {
	switch str {
	case "ASC":
		return "DESC"
	case "DESC":
		return "ASC"
	default:
		return "DESC"
	}
}

// FlipComparisonOperator flips comparison operator
func FlipComparisonOperator(op string) string {
	switch op {
	case ">":
		return "<"
	case "<":
		return ">"
	case ">=":
		return "<="
	case "<=":
		return ">="
	case "=":
		return "!="
	case "!=":
		return "="
	default:
		return op
	}
}
