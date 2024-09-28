package query

type HavingBuilder struct {
	g *ConditionGroup
}

func NewHavingBuilder() *HavingBuilder {
	return &HavingBuilder{
		g: NewConditionGroup(),
	}
}

func (b *HavingBuilder) Having(column string, op Operator, value any) *HavingBuilder {
	if len(b.g.Conditions) == 0 {
		b.g.AddCondition(column, op, value, TypeInit)
	} else {
		b.g.AddCondition(column, op, value, TypeAnd)
	}

	return b
}

func (b *HavingBuilder) OrHaving(column string, op Operator, value any) *HavingBuilder {
	if len(b.g.Conditions) == 0 {
		b.g.AddCondition(column, op, value, TypeInit)
	} else {
		b.g.AddCondition(column, op, value, TypeOr)
	}

	return b
}

func (b *HavingBuilder) HavingFunc(f func(query *HavingBuilder)) *HavingBuilder {
	var newBuilder *HavingBuilder

	if len(b.g.SubGroups) == 0 && len(b.g.Conditions) == 0 {
		newBuilder = convertGroupToHavingBuilder(
			b.g.AddSubGroup(TypeInit),
		)
	} else {
		newBuilder = convertGroupToHavingBuilder(
			b.g.AddSubGroup(TypeAnd),
		)
	}

	f(newBuilder)

	return b
}

func (b *HavingBuilder) OrHavingFunc(f func(query *HavingBuilder)) *HavingBuilder {
	var newBuilder *HavingBuilder

	if len(b.g.SubGroups) == 0 && len(b.g.Conditions) == 0 {
		newBuilder = convertGroupToHavingBuilder(
			b.g.AddSubGroup(TypeInit),
		)
	} else {
		newBuilder = convertGroupToHavingBuilder(
			b.g.AddSubGroup(TypeOr),
		)
	}

	f(newBuilder)

	return b
}

func (b *HavingBuilder) Build() (string, []any) {
	q, arg := b.g.Build()

	if q == "" {
		return "", nil
	}

	return "HAVING " + q, arg
}

func convertGroupToHavingBuilder(cg *ConditionGroup) *HavingBuilder {
	return &HavingBuilder{
		g: cg,
	}
}
