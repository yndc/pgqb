package pgqb

import (
	"fmt"
	"strings"
)

// Builder is the query builder
type Builder struct {
	selects   *strings.Builder
	from      *strings.Builder
	joins     *strings.Builder
	condition *strings.Builder
	groupBy   *strings.Builder
	orderBy   *strings.Builder
	offset    uint
	limit     uint
	err       error
}

// Build a new string builder from the constructed query
func (b *Builder) Build() *strings.Builder {
	result := &strings.Builder{}
	result.WriteRune(' ')
	if b.selects != nil {
		result.WriteString(b.selects.String())
	}
	if b.from != nil {
		result.WriteString(b.from.String())
	}
	if b.joins != nil {
		result.WriteString(b.joins.String())
	}
	if b.condition != nil {
		result.WriteString(b.condition.String())
	}
	if b.groupBy != nil {
		result.WriteString(b.groupBy.String())
	}
	if b.orderBy != nil {
		result.WriteString(b.orderBy.String())
	}
	if b.offset > 0 {
		fmt.Fprintf(result, "OFFSET %d ", b.offset)
	}
	if b.limit > 0 {
		fmt.Fprintf(result, "LIMIT %d ", b.limit)
	}
	return result
}

// Select add a select to the query
func (b *Builder) Select(cols ...string) *Builder {
	if b.selects == nil {
		b.selects = &strings.Builder{}
	}
	if b.selects.Len() == 0 {
		b.selects.WriteString("SELECT ")
	} else {
		b.selects.WriteRune(',')
		b.selects.WriteRune(' ')
	}
	for i, v := range cols {
		b.selects.WriteString(v)
		if i < len(cols)-1 {
			b.selects.WriteRune(',')
			b.selects.WriteRune(' ')
		}
	}
	b.selects.WriteRune(' ')
	return b
}

// From adds a from query
func (b *Builder) From(column string) *Builder {
	if b.from == nil {
		b.from = &strings.Builder{}
	} else {
		b.from.Reset()
	}
	b.from.WriteString("FROM ")
	b.from.WriteString(column)
	b.from.WriteRune(' ')
	return b
}

// Where adds a where query, if a previous where query already exists, then it will be AND-ed
// Use WhereOr to use OR on multiple conditions
func (b *Builder) Where(str string, values ...interface{}) *Builder {
	if b.condition == nil {
		b.condition = &strings.Builder{}
	}
	if b.condition.Len() == 0 {
		b.condition.WriteString("WHERE ")
	} else {
		b.condition.WriteString("AND ")
	}
	fmt.Fprintf(b.condition, str, values...)
	b.condition.WriteRune(' ')
	return b
}

// Or adds an or query
func (b *Builder) Or(str string, values ...interface{}) *Builder {
	if b.condition == nil {
		b.condition = &strings.Builder{}
	}
	if b.condition.Len() == 0 {
		b.condition.WriteString("WHERE ")
	} else {
		b.condition.WriteString("OR ")
	}
	fmt.Fprintf(b.condition, str, values...)
	b.condition.WriteRune(' ')
	return b
}

// Join adds a join query
func (b *Builder) Join(joinType string, targetTable string, on string) *Builder {
	if b.joins == nil {
		b.joins = &strings.Builder{}
	}
	b.joins.WriteString(joinType)
	b.joins.WriteString(" JOIN ")
	b.joins.WriteString(targetTable)
	b.joins.WriteString(" ON ")
	b.joins.WriteString(on)
	b.joins.WriteRune(' ')
	return b
}

// InnerJoin adds an inner join query
func (b *Builder) InnerJoin(targetTable string, on string) *Builder {
	return b.Join("INNER", targetTable, on)
}

// LeftJoin adds a left join query
func (b *Builder) LeftJoin(targetTable string, on string) *Builder {
	return b.Join("LEFT", targetTable, on)
}

// RightJoin adds an inner join query
func (b *Builder) RightJoin(targetTable string, on string) *Builder {
	return b.Join("RIGHT", targetTable, on)
}

// OrderBy adds an order by query
func (b *Builder) OrderBy(str string, mode string) *Builder {
	if b.orderBy == nil {
		b.orderBy = &strings.Builder{}
		b.orderBy.WriteString("ORDER BY ")
	} else {
		b.orderBy.WriteRune(',')
		b.orderBy.WriteRune(' ')
	}
	fmt.Fprintf(b.orderBy, str)
	b.orderBy.WriteRune(' ')
	fmt.Fprintf(b.orderBy, mode)
	b.orderBy.WriteRune(' ')
	return b
}

// GroupBy adds an group by query
func (b *Builder) GroupBy(str string) *Builder {
	if b.groupBy == nil {
		b.groupBy = &strings.Builder{}
		b.groupBy.WriteString("GROUP BY ")
	} else {
		b.groupBy.WriteString(", ")
	}
	fmt.Fprintf(b.groupBy, str)
	b.groupBy.WriteRune(' ')
	return b
}

// Offset adds an offset to the query
// If an offset is already set before, the value will be overwritten
func (b *Builder) Offset(n uint) *Builder {
	b.offset = n
	return b
}

// Limit adds a limit to the query
// If a limit is already set before, the value will be overwritten
func (b *Builder) Limit(n uint) *Builder {
	b.limit = n
	return b
}

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
