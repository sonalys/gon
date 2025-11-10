package gon

import "fmt"

type (
	greaterNode struct {
		first     Node
		second    Node
		inclusive bool
	}
)

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater than the second.
func Greater(first, second Node) Node {
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

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater or equal than the second.
func GreaterOrEqual(first, second Node) Node {
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

func (node greaterNode) Scalar() string {
	if node.inclusive {
		return "gte"
	}

	return "gt"
}

func (node greaterNode) Shape() []KeyNode {
	if node.inclusive {
		return []KeyNode{
			{"first", node.first},
			{"second", node.second},
		}
	}

	return []KeyNode{
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
				Scalar: node.Scalar(),
				Cause: NodeError{
					Scalar: "firstValue",
					Cause:  err,
				},
			})
		}

		if err, ok := secondValue.(error); ok {
			return Literal(NodeError{
				Scalar: node.Scalar(),
				Cause: NodeError{
					Scalar: "secondValue",
					Cause:  err,
				},
			})
		}

		return Literal(NodeError{
			Scalar: node.Scalar(),
			Cause:  fmt.Errorf("cannot compare %T and %T", firstValue, secondValue),
		})
	}

	if node.inclusive {
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}
