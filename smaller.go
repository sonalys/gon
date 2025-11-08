package gon

import (
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
		return Literal(NodeError{
			Scalar: "lt",
			Cause:  fmt.Errorf("cannot compare unset expressions"),
		})
	}

	return smallerNode{
		first:  first,
		second: second,
	}
}

func SmallerOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(NodeError{
			Scalar: "lte",
			Cause:  fmt.Errorf("cannot compare unset expressions"),
		})
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
		if err, ok := firstValue.(error); ok {
			return Literal(NodeError{
				Scalar: node.Name(),
				Cause: NodeError{
					Scalar: "firstValue",
					Cause:  err,
				},
			})
		}

		if err, ok := secondValue.(error); ok {
			return Literal(NodeError{
				Scalar: node.Name(),
				Cause: NodeError{
					Scalar: "secondValue",
					Cause:  err,
				},
			})
		}

		return Literal(NodeError{
			Scalar: node.Name(),
			Cause:  fmt.Errorf("cannot compare %T and %T", firstValue, secondValue),
		})
	}

	if node.inclusive {
		return Literal(comparison <= 0)
	}

	return Literal(comparison < 0)
}

var (
	_ Expression = smallerNode{}
)
