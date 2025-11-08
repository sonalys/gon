package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type orNode struct {
	nodes []Node
}

// Or defines an or node, there must be at least one input.
// It returns the value of the first error, true boolean or non-boolean expression.
func Or(nodes ...Node) Node {
	if len(nodes) == 0 {
		return Literal(NodeError{
			NodeName: "or",
			Cause:    fmt.Errorf("must receive at least one expression"),
		})
	}

	for i := range nodes {
		if nodes[i] == nil {
			return Literal(NodeError{
				NodeName: "or",
				Cause:    fmt.Errorf("all expressions should be not-nil"),
			})
		}
	}

	return orNode{
		nodes: nodes,
	}
}

func (node orNode) Name() string {
	return "or"
}

func (node orNode) Shape() []KeyNode {
	return sliceutils.Map(node.nodes, func(from Node) KeyNode { return KeyNode{Node: from} })
}

func (node orNode) Type() NodeType {
	return NodeTypeExpression
}

func (node orNode) Eval(scope Scope) Value {
	for _, expr := range node.nodes {
		switch value := expr.Eval(scope).Value().(type) {
		case error:
			return Literal(NodeError{
				NodeName: "or",
				Cause:    value,
			})
		case bool:
			if value {
				return Literal(true)
			}
		default:
			return Literal(value)
		}
	}

	return Literal(false)
}
