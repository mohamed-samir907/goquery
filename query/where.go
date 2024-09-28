package query

type WhereBuilder struct {
	g *ConditionGroup
}

func NewWhereBuilder() *WhereBuilder {
	return &WhereBuilder{
		g: NewConditionGroup(),
	}
}

func (b *WhereBuilder) Where(column string, op Operator, value any) *WhereBuilder {
	if len(b.g.Conditions) == 0 {
		b.g.AddCondition(column, op, value, TypeInit)
	} else {
		b.g.AddCondition(column, op, value, TypeAnd)
	}

	return b
}

func (b *WhereBuilder) OrWhere(column string, op Operator, value any) *WhereBuilder {
	if len(b.g.Conditions) == 0 {
		b.g.AddCondition(column, op, value, TypeInit)
	} else {
		b.g.AddCondition(column, op, value, TypeOr)
	}

	return b
}

func (b *WhereBuilder) WhereFunc(f func(query *WhereBuilder)) *WhereBuilder {
	var newBuilder *WhereBuilder

	if len(b.g.SubGroups) == 0 && len(b.g.Conditions) == 0 {
		newBuilder = convertGroupToWhereBuilder(
			b.g.AddSubGroup(TypeInit),
		)
	} else {
		newBuilder = convertGroupToWhereBuilder(
			b.g.AddSubGroup(TypeAnd),
		)
	}

	f(newBuilder)

	return b
}

func (b *WhereBuilder) OrWhereFunc(f func(query *WhereBuilder)) *WhereBuilder {
	var newBuilder *WhereBuilder

	if len(b.g.SubGroups) == 0 && len(b.g.Conditions) == 0 {
		newBuilder = convertGroupToWhereBuilder(
			b.g.AddSubGroup(TypeInit),
		)
	} else {
		newBuilder = convertGroupToWhereBuilder(
			b.g.AddSubGroup(TypeOr),
		)
	}

	f(newBuilder)

	return b
}

func (b *WhereBuilder) Build() (string, []any) {
	q, arg := b.g.Build()

	if q == "" {
		return "", nil
	}

	return "WHERE " + q, arg
}

func convertGroupToWhereBuilder(cg *ConditionGroup) *WhereBuilder {
	return &WhereBuilder{
		g: cg,
	}
}
