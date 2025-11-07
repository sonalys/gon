package gon

import (
	"errors"
	"fmt"
)

type smallerNode struct {
	first     Expression
	second    Expression
	inclusive bool
}

func (node smallerNode) Name() string {
	if node.inclusive {
		return "lte"
	}

	return "lt"
}

func (node smallerNode) Shape() []KeyExpression {
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

func (node smallerNode) Type() NodeType {
	return NodeTypeExpression
}

func Smaller(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("smaller expression cannot compare unset expressions"))
	}

	return smallerNode{
		first:  first,
		second: second,
	}
}

func SmallerOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("smaller or equal expression cannot compare unset expressions"))
	}

	return smallerNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node smallerNode) Eval(scope Scope) Value {
	firstValue := node.first.Eval(scope).Value()
	secondValue := node.second.Eval(scope).Value()

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		errs := make([]error, 0, 2)
		if err, ok := firstValue.(error); ok {
			errs = append(errs, err)
		}
		if err, ok := secondValue.(error); ok {
			errs = append(errs, err)
		}

		return propagateErr(Literal(errors.Join(errs...)), "cannot compare different types: %T and %T", firstValue, secondValue)
	}

	if node.inclusive {
		return Literal(comparison <= 0)
	}

	return Literal(comparison < 0)
}

var (
	_ Expression = smallerNode{}
)
