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
	combine   *strings.Builder
	groupBy   *strings.Builder
	orderBy   *strings.Builder
	ctes      *strings.Builder
	offset    uint
	limit     uint
	err       error
}

// NewBuilder creates a new builder
func NewBuilder() *Builder {
	return &Builder{}
}

// Copy copies an existing builder
func (b *Builder) Copy() *Builder {
	newBuilder := &Builder{}
	if b.selects != nil {
		newBuilder.selects = &strings.Builder{}
		newBuilder.selects.WriteString(b.selects.String())
	}
	if b.from != nil {
		newBuilder.from = &strings.Builder{}
		newBuilder.from.WriteString(b.from.String())
	}
	if b.joins != nil {
		newBuilder.joins = &strings.Builder{}
		newBuilder.joins.WriteString(b.joins.String())
	}
	if b.condition != nil {
		newBuilder.condition = &strings.Builder{}
		newBuilder.condition.WriteString(b.condition.String())
	}
	if b.combine != nil {
		newBuilder.combine = &strings.Builder{}
		newBuilder.combine.WriteString(b.combine.String())
	}
	if b.groupBy != nil {
		newBuilder.groupBy = &strings.Builder{}
		newBuilder.groupBy.WriteString(b.groupBy.String())
	}
	if b.orderBy != nil {
		newBuilder.orderBy = &strings.Builder{}
		newBuilder.orderBy.WriteString(b.orderBy.String())
	}
	if b.ctes != nil {
		newBuilder.ctes = &strings.Builder{}
		newBuilder.ctes.WriteString(b.ctes.String())
	}
	if b.offset > 0 {
		newBuilder.offset = b.offset
	}
	if b.limit > 0 {
		newBuilder.limit = b.limit
	}
	newBuilder.err = b.err
	return newBuilder
}

// Build a new string builder from the constructed query
func (b *Builder) Build() *strings.Builder {
	result := &strings.Builder{}
	if b.ctes != nil {
		result.WriteString(b.ctes.String())
	}
	if b.selects != nil {
		result.WriteString(b.selects.String())
		result.WriteRune(' ')
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
	if b.combine != nil {
		result.WriteString(b.combine.String())
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

// SubQueryString builds this query as a subquery string
func (b *Builder) SubQueryString(alias string) string {
	s := b.Build()
	return fmt.Sprintf("( %s ) AS %s", s.String(), alias)
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

// Union adds another query to union
func (b *Builder) Union(other string) *Builder {
	if b.combine == nil {
		b.combine = &strings.Builder{}
	}
	b.combine.WriteString("UNION ")
	b.combine.WriteString(other)
	return b
}

// Intersect adds another query to intercept
func (b *Builder) Intersect(other string) *Builder {
	if b.combine == nil {
		b.combine = &strings.Builder{}
	}
	b.combine.WriteString("INTERSECT ")
	b.combine.WriteString(other)
	return b
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
	fmt.Fprint(b.orderBy, str)
	b.orderBy.WriteRune(' ')
	fmt.Fprint(b.orderBy, mode)
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
	fmt.Fprint(b.groupBy, str)
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

// With add a CTE to the query
func (b *Builder) With(alias string, query func(builder *Builder)) *Builder {
	builder := &Builder{}
	query(builder)
	if b.ctes == nil {
		b.ctes = &strings.Builder{}
		b.ctes.WriteString("WITH ")
	} else if b.ctes.Len() == 0 {
		b.ctes.WriteString("WITH ")
	} else {
		b.ctes.WriteRune(',')
		b.ctes.WriteRune(' ')
	}
	b.ctes.WriteString(alias)
	b.ctes.WriteString(" AS ( ")
	b.ctes.WriteString(builder.Build().String())
	b.ctes.WriteString(") ")
	return b
}

// With add a CTE to the query
func (b *Builder) WithStr(alias string, queryStr string) *Builder {
	if b.ctes == nil {
		b.ctes = &strings.Builder{}
		b.ctes.WriteString("WITH ")
	} else if b.ctes.Len() == 0 {
		b.ctes.WriteString("WITH ")
	} else {
		b.ctes.WriteRune(',')
		b.ctes.WriteRune(' ')
	}
	b.ctes.WriteString(alias)
	b.ctes.WriteString(" AS ( ")
	b.ctes.WriteString(queryStr)
	b.ctes.WriteString(") ")
	return b
}
