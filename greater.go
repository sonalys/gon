package gon

import "fmt"

type (
	greaterNode struct {
		first     Expression
		second    Expression
		inclusive bool
	}
)

func Greater(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(NodeError{
			Scalar: "gt",
			Cause:  fmt.Errorf("cannot compare unset expressions"),
		})
	}

	return greaterNode{
		first:  first,
		second: second,
	}
}

func GreaterOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(NodeError{
			Scalar: "gte",
			Cause:  fmt.Errorf("cannot compare unset expressions"),
		})
	}
	return greaterNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

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

func (node greaterNode) Eval(scope Scope) Value {
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
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}
