package pgqb

import (
	"fmt"
	"strings"
)

// ConditionBuilder is a builder for conditional queries
type ConditionBuilder struct {
	str *strings.Builder
}

// And adds an and condition to the builder
func (b *ConditionBuilder) And(str string, args ...interface{}) {
	if b.str == nil {
		b.str = &strings.Builder{}
	} else {
		b.str.WriteString(` AND `)
	}
	fmt.Fprintf(b.str, str, args...)
}

// Or adds an and condition to the builder
func (b *ConditionBuilder) Or(str string, args ...interface{}) {
	if b.str == nil {
		b.str = &strings.Builder{}
	} else {
		b.str.WriteString(` OR `)
	}
	fmt.Fprintf(b.str, str, args...)
}

func (b *ConditionBuilder) String() string {
	return fmt.Sprintf(`( %s )`, b.str.String())
}

// Condition creates a condition query
func Condition(query func(builder *ConditionBuilder)) string {
	builder := &ConditionBuilder{}
	query(builder)
	return builder.String()
}
