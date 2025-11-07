package gon

import "fmt"

type (
	greaterNode struct {
		first     Expression
		second    Expression
		inclusive bool
	}
)

func (node greaterNode) Name() string {
	if node.inclusive {
		return "gte"
	}

	return "gt"
}

func (node greaterNode) Shape() []KeyExpression {
	if node.inclusive {
		return []KeyExpression{
			{"first", node.first},
			{"second", node.second},
		}
	}

	return []KeyExpression{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node greaterNode) Type() NodeType {
	return NodeTypeExpression
}

func Greater(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("greater expression cannot compare unset expressions"))
	}

	return greaterNode{
		first:  first,
		second: second,
	}
}

func GreaterOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("greater or equal expression cannot compare unset expressions"))
	}
	return greaterNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node greaterNode) Eval(scope Scope) Value {
	firstValue := node.first.Eval(scope).Value()
	secondValue := node.second.Eval(scope).Value()

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return propagateErr(nil, "cannot compare different types: %T and %T", firstValue, secondValue)
	}

	if node.inclusive {
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}

var (
	_ Expression = greaterNode{}
)
