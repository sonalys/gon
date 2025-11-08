package gon

import "fmt"

type equalNode struct {
	first  Expression
	second Expression
}

func Equal(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("equal expression cannot compare unset expressions"))
	}

	return equalNode{
		first:  first,
		second: second,
	}
}

func (node equalNode) Name() string {
	return "equal"
}

func (node equalNode) Shape() []KeyExpression {
	return []KeyExpression{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node equalNode) Type() NodeType {
	return NodeTypeExpression
}

func (node equalNode) Eval(scope Scope) Value {
	firstValue := node.first.Eval(scope).Value()
	secondValue := node.second.Eval(scope).Value()

	value, ok := cmpAny(firstValue, secondValue)
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

	return Literal(value == 0)
}
