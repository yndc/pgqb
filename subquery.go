package postgresqb

import "strings"

// SubQuery creates a subquery
func SubQuery(alias string, subQueryBuilder func(builder *Builder)) string {
	builder := &Builder{}
	subQueryBuilder(builder)
	result := &strings.Builder{}
	result.WriteString("( ")
	result.WriteString(builder.Build().String())
	result.WriteString(") AS ")
	result.WriteString(alias)
	result.WriteRune(' ')
	return result.String()
}
