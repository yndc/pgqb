# pgqb

`pgqb` (Postgres Query Builder) is a small utility library to build PostgreSQL query strings in Go.

Features:

- Supported commands: `SELECT`, `FROM`, `WHERE`, `JOIN` (INNER, LEFT, and RIGHT), `ORDER BY`, `GROUP BY`, `OFFSET`, `LIMIT`
- SubQueries builder
- Nested conditionals builder

Please be aware that the inputs are not sanitized by default. There's a helper function to do it, `Sanitize`.

## Examples

### Basic `SELECT` query

```golang
builder := pgqb.Builder{}
builder.Select("one", "two").From("some_table")
```

```sql
SELECT one,two FROM some_table
```

### Multiple Joins

```golang
builder := pgqb.Builder{}
builder.
    Select("one", "two").
    From("some_table").
    InnerJoin("other_table", "other_table.id = some_table.other_id").
	InnerJoin("another_table", "another_table.id = other_table.another_id")
```

```sql
SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN another_table ON another_table.id = other_table.another_id
```

### Joins with SubQuery

```golang
builder := pgqb.Builder{}
builder.
    Select("one", "two").
    From("some_table").
    InnerJoin("other_table", "other_table.id = some_table.other_id").
    InnerJoin(pgqb.Sub(func(builder *pgqb.Builder) {
        builder.
            Select("sqcol").
            From("sqtable").
            Where("sqcol2 = 'abc'").
            Limit(1)
    }, "sq"), "sq.sqcol = 'x'").
    Limit(20).
    OrderBy("id", "ASC")
```

```sql
SELECT one, two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN ( SELECT sqcol FROM sqtable WHERE sqcol2 = 'abc' LIMIT 1 ) AS sq ON sq.sqcol = 'x' ORDER BY id ASC LIMIT 20  
```

### Nested conditionals

```golang
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
```

```sql
SELECT one, two FROM some_table WHERE ( a = 1 OR b > 10 OR ( c = 0 AND d = 0 ) ) AND z = 0 
```