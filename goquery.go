package goquery

import (
	"fmt"

	"github.com/mohamed-samir907/goquery/drivers"
	"github.com/mohamed-samir907/goquery/query"
)

type Database struct {
	driver drivers.Driver
}

func New(driver drivers.Driver) *Database {
	return &Database{
		driver: driver,
	}
}

type Query struct {
	driver   drivers.Driver
	table    string
	distinct bool
	columns  []string
	where    query.WhereBuilder
	groupBy  []string
	having   query.HavingBuilder
	orderBy  []string
	limit    int
	offset   *int
}

func (d *Database) Table(table string) *Query {
	return &Query{
		driver:   d.driver,
		table:    table,
		distinct: false,
		columns:  []string{"*"},
		groupBy:  make([]string, 0, 5),
		orderBy:  make([]string, 0, 2),
		where:    query.NewWhereBuilder(),
		having:   query.NewHavingBuilder(),
	}
}

func (q *Query) Select(columns ...string) *Query {
	q.columns = columns
	return q
}

func (q *Query) Distinct() *Query {
	q.distinct = true
	return q
}

// func (q *Query) Join(table string, joinBuilder func(b *JoinBuilder)) *Query {
// 	return q
// }

func (q *Query) Where(column string, op query.Operator, value any) *Query {
	q.where.Where(column, op, value)
	return q
}

func (q *Query) OrWhere(column string, op query.Operator, value any) *Query {
	q.where.OrWhere(column, op, value)
	return q
}

func (q *Query) WhereFunc(f func(builder query.WhereBuilder)) *Query {
	q.where.WhereFunc(f)
	return q
}

func (q *Query) OrWhereFunc(f func(builder query.WhereBuilder)) *Query {
	q.where.OrWhereFunc(f)
	return q
}

func (q *Query) GroupBy(columns ...string) *Query {
	q.groupBy = columns
	return q
}

func (q *Query) Having(column string, op query.Operator, value any) *Query {
	q.having.Having(column, op, value)
	return q
}

func (q *Query) OrHaving(column string, op query.Operator, value any) *Query {
	q.having.OrHaving(column, op, value)
	return q
}

func (q *Query) HavingFunc(f func(builder query.HavingBuilder)) *Query {
	q.having.HavingFunc(f)
	return q
}

func (q *Query) OrHavingFunc(f func(builder query.HavingBuilder)) *Query {
	q.having.OrHavingFunc(f)
	return q
}

func (q *Query) OrderBy(column, orderType string) *Query {
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s %s", column, orderType))
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = &offset
	return q
}

func (q *Query) Get() ([]map[string]any, error) {
	query := query.SelectQuery{
		Table:    q.table,
		Distinct: q.distinct,
		Columns:  q.columns,
		Where:    q.where,
		GroupBy:  q.groupBy,
		Having:   q.having,
		OrderBy:  q.orderBy,
		Limit:    q.limit,
		Offset:   q.offset,
	}

	return q.driver.Get(query)
}

func (q *Query) First() (map[string]any, error) {
	rows, err := q.Limit(1).Get()
	if err != nil {
		return nil, err
	}

	return rows[0], nil
}

func (q *Query) Insert() {
	q.driver.Insert()
}

func (q *Query) Update() {
	q.driver.Update()
}

func (q *Query) Delete() {
	q.driver.Delete()
}
