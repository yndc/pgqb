package postgresqb

import (
	"fmt"
	"strings"
)

// Builder is the query builder
type Builder struct {
	selects *strings.Builder
	from    *strings.Builder
	joins   *strings.Builder
	where   *strings.Builder
	groupBy *strings.Builder
	orderBy *strings.Builder
	offset  uint
	limit   uint
	err     error
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
	if b.where != nil {
		result.WriteString(b.where.String())
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
	}
	for i, v := range cols {
		b.selects.WriteString(DoubleQuote(v))
		if i < len(cols)-1 {
			b.selects.WriteRune(',')
		}
	}
	b.selects.WriteRune(' ')
	return b
}

// From adds a from query
func (b *Builder) From(column string) *Builder {
	if b.from == nil {
		b.from = &strings.Builder{}
		b.from.WriteString("FROM ")
	} else {
		b.from.Reset()
	}
	b.from.WriteString(DoubleQuote(column))
	b.from.WriteRune(' ')
	return b
}

// Where adds a where query
func (b *Builder) Where(str string, values ...interface{}) *Builder {
	if b.where == nil {
		b.where = &strings.Builder{}
	}
	if b.where.Len() == 0 {
		b.where.WriteString("WHERE ")
	} else {
		b.where.WriteString("AND ")
	}
	fmt.Fprintf(b.where, str, values...)
	b.where.WriteRune(' ')
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

// Or adds an or query
func (b *Builder) Or(str string, values ...interface{}) *Builder {
	if b.where == nil {
		b.where = &strings.Builder{}
	}
	if b.where.Len() == 0 {
		b.where.WriteString("WHERE ")
	} else {
		b.where.WriteString("OR ")
	}
	fmt.Fprintf(b.where, str, values...)
	b.where.WriteRune(' ')
	return b
}

// OrderBy adds an order by query
func (b *Builder) OrderBy(str string, mode string) *Builder {
	if b.orderBy == nil {
		b.orderBy = &strings.Builder{}
		b.orderBy.WriteString("ORDER BY ")
	} else {
		b.orderBy.WriteString(", ")
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
