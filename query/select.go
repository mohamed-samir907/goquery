package query

import "strings"

type SelectQuery struct {
	Table    string
	Distinct bool
	Columns  []string
	Where    *WhereBuilder
	GroupBy  []string
	Having   *HavingBuilder
	OrderBy  []string
	Limit    int
	Offset   *int
}

func (q *SelectQuery) Build() (string, []any) {
	var b strings.Builder
	args := []any{}

	b.WriteString("SELECT ")
	b.WriteString(strings.Join(q.Columns, ", "))
	b.WriteString(" FROM ")
	b.WriteString(q.Table)
	b.WriteString(" ")

	// TODO: build join here

	// add WHERE to the query
	if q.Where != nil {
		cond, arg := q.Where.Build()

		b.WriteString(cond)
		b.WriteString(" ")

		args = append(args, arg...)
	}

	// add GROUP BY to the query
	if len(q.GroupBy) > 0 {
		b.WriteString("GROUP BY ")
		b.WriteString(strings.Join(q.GroupBy, ", "))
		b.WriteString(" ")
	}

	// add HAVING to the query
	if q.Having != nil {
		cond, arg := q.Having.Build()

		b.WriteString(cond)
		b.WriteString(" ")

		args = append(args, arg...)
	}

	// add ORDER BY to the query
	if len(q.OrderBy) > 0 {
		b.WriteString("ORDER BY ")
		b.WriteString(strings.Join(q.OrderBy, ", "))
		b.WriteString(" ")
	}

	// add LIMIT to the query
	if q.Limit > 0 {
		b.WriteString("LIMIT ?")
		b.WriteString(" ")
		args = append(args, q.Limit)
	}

	// add OFFSET to the query
	if q.Offset != nil {
		b.WriteString("OFFSET ?")
		args = append(args, *q.Offset)
	}

	return b.String(), args
}
