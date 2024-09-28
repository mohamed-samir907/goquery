package query

import (
	"strings"
)

type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
)

type JoinBuilder struct {
	joinType JoinType
	table    string
	left     string
	op       string
	right    string
}

func NewJoinBuilder(table string, joinType JoinType) *JoinBuilder {
	return &JoinBuilder{
		joinType: JoinTypeInner,
		table:    table,
	}
}

// b.On("users.id", "=", "orders.user_id")
func (b *JoinBuilder) On(left, op, right string) *JoinBuilder {
	b.left = left
	b.right = right
	b.op = op
	return b
}

func (b *JoinBuilder) Using(column string) *JoinBuilder {
	b.left = column
	return b
}

func (b *JoinBuilder) Build() string {
	var query strings.Builder

	query.WriteString(string(b.joinType))
	query.WriteString(" JOIN ")
	query.WriteString(b.table)
	query.WriteString(" ON ")
	query.WriteString(b.left)
	query.WriteString(" ")
	query.WriteString(b.op)
	query.WriteString(" ")
	query.WriteString(b.right)

	return query.String()
}
