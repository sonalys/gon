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
		return Literal(NodeError{
			NodeName: "equal",
			Cause:    fmt.Errorf("all inputs should be not-nil"),
		})
	}

	return equalNode{
		first:  first,
		second: second,
	}
}

func (node equalNode) Name() string {
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
	firstValue := node.first.Eval(scope).Value()
	secondValue := node.second.Eval(scope).Value()

	value, ok := cmpAny(firstValue, secondValue)
	if !ok {
		if err, ok := firstValue.(error); ok {
			return Literal(NodeError{
				NodeName: node.Name(),
				Cause: NodeError{
					NodeName: "firstValue",
					Cause:    err,
				},
			})
		}

		if err, ok := secondValue.(error); ok {
			return Literal(NodeError{
				NodeName: node.Name(),
				Cause: NodeError{
					NodeName: "secondValue",
					Cause:    err,
				},
			})
		}

		return Literal(NodeError{
			NodeName: node.Name(),
			Cause:    fmt.Errorf("cannot compare %T and %T", firstValue, secondValue),
		})
	}

	return Literal(value == 0)
}
