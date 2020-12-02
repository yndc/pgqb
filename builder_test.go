package postgresqb_test

import (
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
	if builder.Build().String() != " SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN another_table ON another_table.id = other_table.another_id " {
		t.Fail()
	}
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
	if builder.Build().String() != " SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN ( SELECT sqcol FROM sqtable WHERE sqcol2 = 'abc' LIMIT 1 ) AS sq ON sq.sqcol = 'x' ORDER BY id ASC LIMIT 20 " {
		t.Fail()
	}
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
	if builder.Build().String() != " SELECT one, two FROM some_table WHERE ( a = 1 OR b > 10 OR ( c = 0 AND d = 0 ) ) AND z = 0 " {
		t.Fail()
	}
}
