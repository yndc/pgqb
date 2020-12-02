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
		From("some_table").
		InnerJoin("other_table", "other_table.id = some_table.other_id").
		InnerJoin("another_table", "another_table.id = other_table.another_id").Limit(20).Offset(10)
	fmt.Println(builder.Build().String())
}
