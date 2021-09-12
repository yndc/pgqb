package pgqb_test

import (
	"fmt"
	"testing"

	"github.com/yndc/pgqb"
)

func TestJoin(t *testing.T) {
	builder := pgqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin("another_table", "another_table.id = other_table.another_id")
	if builder.Build().String() != "SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN another_table ON another_table.id = other_table.another_id " {
		t.Fail()
	}
}

func TestJoinSubQuery(t *testing.T) {
	builder := pgqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin(pgqb.Sub(func(builder *pgqb.Builder) {
			builder.Select("sqcol").
				From("sqtable").
				Where("sqcol2 = 'abc'").
				Limit(1)
		}, "sq"), "sq.sqcol = 'x'").
		Limit(20).
		OrderBy("id", "ASC")
	if builder.Build().String() != "SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN ( SELECT sqcol FROM sqtable WHERE sqcol2 = 'abc' LIMIT 1 ) AS sq ON sq.sqcol = 'x' ORDER BY id ASC LIMIT 20 " {
		t.Fail()
	}
}

func TestCondition(t *testing.T) {
	builder := pgqb.Builder{}
	builder.
		Select("one", "two").
		From("some_table").
		Where(pgqb.Condition(func(builder *pgqb.ConditionBuilder) {
			builder.Or("a = 1")
			builder.Or("b > 10")
			builder.Or(pgqb.Condition(func(builder *pgqb.ConditionBuilder) {
				builder.And("c = 0")
				builder.And("d = 0")
			}))
		})).
		Where("z = 0")
	if builder.Build().String() != "SELECT one, two FROM some_table WHERE ( a = 1 OR b > 10 OR ( c = 0 AND d = 0 ) ) AND z = 0 " {
		t.Fail()
	}
}

func TestUnion(t *testing.T) {
	builder := pgqb.Builder{}
	builder.Select("one", "two")
	builder.From("some_table")
	builder.Where(pgqb.Condition(func(builder *pgqb.ConditionBuilder) {
		builder.Or("a = 1")
		builder.Or("b > 10")
		builder.Or(pgqb.Condition(func(builder *pgqb.ConditionBuilder) {
			builder.And("c = 0")
			builder.And("d = 0")
		}))
	}))
	builder.Where("z = 0")
	builder.OrderBy("one", "ASC")
	builder.GroupBy("one")
	builder.Limit(10)
	builder.Union(pgqb.NewBuilder().Select("one", "two").From("another_one").Build().String())
	if builder.Build().String() != "SELECT one, two FROM some_table WHERE ( a = 1 OR b > 10 OR ( c = 0 AND d = 0 ) ) AND z = 0 UNION SELECT one, two FROM another_one GROUP BY one ORDER BY one ASC LIMIT 10 " {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	original := pgqb.Builder{}
	original.Select("ayy", "lmao").From("lol")
	copy := original.Copy()
	copy.From("yikes")
	copy.Select("uh")
	copy.InnerJoin("something", "something.lel = lol.lmao")
	if original.Build().String() != "SELECT ayy, lmao FROM lol " {
		t.Fail()
	}
	if copy.Build().String() != "SELECT ayy, lmao, uh FROM yikes INNER JOIN something ON something.lel = lol.lmao " {
		t.Fail()
	}
}

func TestCte(t *testing.T) {
	builder := pgqb.Builder{}
	builder.
		With("first", func(builder *pgqb.Builder) {
			builder.Select("first", "stuff").From("first_table")
		}).
		With("second", func(builder *pgqb.Builder) {
			builder.Select("second", "stuff").From("second_table")
		}).
		Select("first", "second").
		From("first").
		InnerJoin("second", "first.stuff = second.stuff").
		InnerJoin(pgqb.Sub(func(builder *pgqb.Builder) {
			builder.Select("third").
				From("third").
				Where("sqcol2 = 'abc'").
				Limit(1)
		}, "sq"), "sq.sqcol = 'x'").
		Limit(20).
		OrderBy("id", "ASC")

	fmt.Println(builder.Build().String())
	if builder.Build().String() != "SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN ( SELECT sqcol FROM sqtable WHERE sqcol2 = 'abc' LIMIT 1 ) AS sq ON sq.sqcol = 'x' ORDER BY id ASC LIMIT 20 " {
		t.Fail()
	}
}
