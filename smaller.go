package gon

import (
	"fmt"
)

type smallerNode struct {
	first     Node
	second    Node
	inclusive bool
}

// Smaller defines a smaller node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller to the second.
func Smaller(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "lt",
			Cause:      fmt.Errorf("cannot compare unset expressions"),
		}
	}

	return smallerNode{
		first:  first,
		second: second,
	}
}

// SmallerOrEqual defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller or equal to the second.
func SmallerOrEqual(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "lte",
			Cause:      fmt.Errorf("cannot compare unset expressions"),
		}
	}

	return smallerNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node smallerNode) Scalar() string {
	if node.inclusive {
		return "lte"
	}

	return "lt"
}

func (node smallerNode) Shape() []KeyNode {
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

func (node smallerNode) Type() NodeType {
	return NodeTypeExpression
}

func (node smallerNode) Eval(scope Scope) Value {
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
		return Literal(comparison <= 0)
	}

	return Literal(comparison < 0)
}

var (
	_ Node = smallerNode{}
)
