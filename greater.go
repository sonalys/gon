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
		return NodeError{
			NodeScalar: "gt",
			Cause:      fmt.Errorf("cannot compare unset expressions"),
		}
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
		return NodeError{
			NodeScalar: "gte",
			Cause:      fmt.Errorf("cannot compare unset expressions"),
		}
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
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return NewNodeError(node, err)
	}

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return NewNodeError(node, fmt.Errorf("cannot compare %T and %T", firstValue, secondValue))
	}

	if node.inclusive {
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}
