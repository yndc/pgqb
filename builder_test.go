package postgresqb_test

import (
	"fmt"
	"testing"

	"github.com/yndc/postgresqb"
)

func TestJoin(t *testing.T) {
	builder := postgresqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin("another_table", "another_table.id = other_table.another_id")
	fmt.Println(builder.Build().String())
}

func TestJoinSubQuery(t *testing.T) {
	builder := postgresqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin(postgresqb.Sub(func(builder *postgresqb.Builder) {
			builder.Select("sqcol").
				From("sqtable").
				Where("sqcol2 = 'abc'").
				Limit(1)
		}, "sq"), "sq.sqcol = 'x'").
		Limit(20).
		OrderBy("id", "ASC")
	fmt.Println(builder.Build().String())

}

func TestCondition(t *testing.T) {
	builder := postgresqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		Where(postgresqb.Condition(func(builder *postgresqb.ConditionBuilder) {
			builder.Or("a = 1")
			builder.Or("b > 10")
			builder.Or(postgresqb.Condition(func(builder *postgresqb.ConditionBuilder) {
				builder.And("c = 0")
				builder.And("d = 0")
			}))
		})).
		Where("z = 0")
	fmt.Println(builder.Build().String())
}
