package query

import "strings"

type SelectQuery struct {
	Table    string
	Distinct bool
	Columns  []string
	Where    WhereBuilder
	GroupBy  []string
	Having   HavingBuilder
	OrderBy  []string
	Limit    int
	Offset   *int
	builder  strings.Builder
}

func (q *SelectQuery) Build() (string, []any) {
	q.builder.Grow(256)

	args := make([]any, 0, 10)

	q.builder.WriteString("SELECT ")
	q.builder.WriteString(strings.Join(q.Columns, ", "))
	q.builder.WriteString(" FROM ")
	q.builder.WriteString(q.Table)
	q.builder.WriteString(" ")

	// TODO: build join here

	// add WHERE to the query
	cond, arg := q.Where.Build()
	if cond != "" {
		q.builder.WriteString(cond)
		q.builder.WriteString(" ")

		args = append(args, arg...)
	}

	// add GROUP BY to the query
	if len(q.GroupBy) > 0 {
		q.builder.WriteString("GROUP BY ")
		q.builder.WriteString(strings.Join(q.GroupBy, ", "))
		q.builder.WriteString(" ")
	}

	// add HAVING to the query
	cond, arg = q.Having.Build()
	if cond != "" {
		q.builder.WriteString(cond)
		q.builder.WriteString(" ")

		args = append(args, arg...)
	}

	// add ORDER BY to the query
	if len(q.OrderBy) > 0 {
		q.builder.WriteString("ORDER BY ")
		q.builder.WriteString(strings.Join(q.OrderBy, ", "))
		q.builder.WriteString(" ")
	}

	// add LIMIT to the query
	if q.Limit > 0 {
		q.builder.WriteString("LIMIT ?")
		q.builder.WriteString(" ")
		args = append(args, q.Limit)
	}

	// add OFFSET to the query
	if q.Offset != nil {
		q.builder.WriteString("OFFSET ?")
		args = append(args, *q.Offset)
	}

	return q.builder.String(), args
}
