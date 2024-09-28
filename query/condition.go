package query

import (
	"strings"
)

// Operator represent a condition operator.
type Operator string

const (
	Eq      Operator = "="
	Neq     Operator = "!="
	Gt      Operator = ">"
	Lt      Operator = "<"
	Gte     Operator = ">="
	Lte     Operator = "<="
	Like    Operator = "LIKE"
	NotLike Operator = "NOT LIKE"
	Is      Operator = "IS"
	IsNot   Operator = "IS NOT"
	In      Operator = "IN"
	NotIn   Operator = "NOT IN"
)

type Type string

const (
	TypeAnd  Type = "AND"
	TypeOr   Type = "OR"
	TypeInit Type = ""
)

type Condition struct {
	Column   string
	Operator Operator
	Value    any
	Type     Type
}

type ConditionGroup struct {
	Conditions []Condition
	Type       Type
	SubGroups  []ConditionGroup
}

func NewConditionGroup() *ConditionGroup {
	return &ConditionGroup{
		Type: TypeInit,
	}
}

func (cg *ConditionGroup) AddCondition(column string, op Operator, value any, _type Type) *ConditionGroup {
	cg.Conditions = append(cg.Conditions, Condition{
		Column:   column,
		Operator: op,
		Value:    value,
		Type:     _type,
	})
	return cg
}

func (cg *ConditionGroup) AddSubGroup(_type Type) *ConditionGroup {
	newGroup := ConditionGroup{Type: _type}
	cg.SubGroups = append(cg.SubGroups, newGroup)
	return &cg.SubGroups[len(cg.SubGroups)-1]
}

func (cg *ConditionGroup) Build() (string, []any) {
	return cg.buildGroup(*cg, true)
}

func (cb *ConditionGroup) buildGroup(group ConditionGroup, root bool) (string, []any) {
	var conditions []string
	var args []any

	for _, cond := range group.Conditions {
		condition, arg := cb.buildCondition(cond)
		conditions = append(conditions, condition)
		args = append(args, arg...)
	}

	for _, subGroup := range group.SubGroups {
		subConditions, subArgs := cb.buildGroup(subGroup, false)

		if subConditions == "" {
			continue
		}

		// this will detect the Logical operator for the sub group conditions
		// because first WhereFunc make TypeInit and if the query has conditions
		// then you need to convert TypeInit to TypeAnd
		condType := string(subGroup.Type)
		if subGroup.Type == TypeInit && len(conditions) > 0 {
			condType = string(TypeAnd)
		}

		// this adds a whitespace between logical operator and condition
		if condType != "" {
			condType = condType + " "
		}

		conditions = append(conditions, condType+subConditions)
		args = append(args, subArgs...)
	}

	if len(conditions) == 0 {
		return "", nil
	}

	if root {
		return strings.Join(conditions, " "), args
	}

	return "(" + strings.Join(conditions, " ") + ")", args
}

func (cb *ConditionGroup) buildCondition(cond Condition) (string, []any) {
	var stmt strings.Builder

	// this will build the beginning of the condition like: "AND " the space will be added
	// if the cond.Type is not TypeInit
	stmt.WriteString(string(cond.Type))
	if cond.Type != TypeInit {
		stmt.WriteString(" ")
	}

	// contine build the beginning of the condition like: "id = " and the value will
	// be deterimined later based on the cond.Operator
	stmt.WriteString(cond.Column)
	stmt.WriteString(" ")
	stmt.WriteString(string(cond.Operator))

	switch cond.Operator {
	// this will build stmt: AND id IN (?, ?, ?, ...)
	case In, NotIn:
		if values, ok := cond.Value.([]any); ok {
			placeholders := make([]string, len(values))
			for i := range placeholders {
				placeholders[i] = "?"
			}

			stmt.WriteString(" (")
			stmt.WriteString(strings.Join(placeholders, ","))
			stmt.WriteString(")")

			return stmt.String(), values
		}

	// this will build stmt: AND id IS NOT NULL
	case Is, IsNot:
		stmt.WriteString(" NULL")

		return stmt.String(), []any{}

	// this will build stmt: AND id = ?
	default:
		stmt.WriteString(" ?")
	}

	return stmt.String(), []any{cond.Value}
}
