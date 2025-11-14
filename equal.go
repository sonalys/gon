package gon

import "fmt"

type equalNode struct {
	first  Node
	second Node
}

// Equal defines an equality node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the inputs are equal.
func Equal(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "equal",
			Cause:      fmt.Errorf("all inputs should be not-nil"),
		}
	}

	return equalNode{
		first:  first,
		second: second,
	}
}

func (node equalNode) Scalar() string {
	return "equal"
}

func (node equalNode) Shape() []KeyNode {
	return []KeyNode{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node equalNode) Type() NodeType {
	return NodeTypeExpression
}

func (node equalNode) Eval(scope Scope) Value {
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return NewNodeError(node, err)
	}

	value, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return NewNodeError(node, fmt.Errorf("cannot compare %T and %T", firstValue, secondValue))
	}

	return Literal(value == 0)
}
