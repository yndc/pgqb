package postgresqb_test

import (
	"fmt"
	"testing"

	"github.com/yndc/postgresqb"
)

func TestBuilder(t *testing.T) {
	builder := postgresqb.Builder{}
	builder.
		Select("one", "two", "three").
		From("wrong_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin("another_table", "another_table.id = other_table.another_id").
		InnerJoin(postgresqb.Sub(func(builder *postgresqb.Builder) {
			builder.Select("sqcol").
				From("sqtable").
				Where("sqcol2 = 'abc'").
				Limit(1)
		}, "sq"), "sq.sqcol = 'x'").
		Limit(20).
		Offset(10).
		From("some_table").
		Limit(100).
		Offset(5).
		OrderBy("id", postgresqb.FlipSortingMode("ASC"))
	fmt.Println(builder.Build().String())
}
