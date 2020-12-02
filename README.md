# postgresqb

`postgresqb` (Postgres Query Builder) is a small utility library to build PostgreSQL query strings in Go.

Features:

- Supported commands: `SELECT`, `FROM`, `WHERE`, `JOIN` (INNER, LEFT, and RIGHT), `ORDER BY`, `GROUP BY`, `OFFSET`, `LIMIT`
- SubQueries

Please be aware that the inputs are not sanitized by default. There's a helper function to do it, `Sanitize`.

## Examples

### Basic `SELECT` query

```golang
builder := postgresqb.Builder{}
builder.Select("one", "two").From("some_table")
```

```sql
SELECT one,two FROM some_table
```

### Multiple Joins

```golang
builder := postgresqb.Builder{}
builder.
    Select("one", "two").
    From("some_table").
    InnerJoin("other_table", "other_table.id = some_table.other_id").
	InnerJoin("another_table", "another_table.id = other_table.another_id")
```

```sql
SELECT one,two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN another_table ON another_table.id = other_table.another_id 
```

### Joins with SubQuery

```golang
builder := postgresqb.Builder{}
builder.
    Select("one", "two").
    From("some_table").
    InnerJoin("other_table", "other_table.id = some_table.other_id").
    InnerJoin(postgresqb.Sub(func(builder *postgresqb.Builder) {
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
 SELECT one,two FROM some_table INNER JOIN other_table ON other_table.id = some_table.other_id INNER JOIN ( SELECT sqcol FROM sqtable WHERE sqcol2 = 'abc' LIMIT 1 ) AS sq ON sq.sqcol = 'x' ORDER BY id ASC LIMIT 20  
```