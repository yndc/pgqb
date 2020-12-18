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

// InvertOperator flips comparison operator
func InvertOperator(op string) string {
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
	case "IS NULL":
		return "IS NOT NULL"
	case "IS NOT NULL":
		return "IS NULL"
	case "TRUE":
		return "FALSE"
	case "FALSE":
		return "TRUE"
	default:
		return op
	}
}
