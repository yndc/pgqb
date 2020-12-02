package postgresqb

import "strings"

// Sub creates a subquery
func Sub(query func(builder *Builder), alias string) string {
	builder := &Builder{}
	query(builder)
	result := &strings.Builder{}
	result.WriteString("(")
	result.WriteString(builder.Build().String())
	result.WriteString(") AS ")
	result.WriteString(alias)
	return result.String()
}
