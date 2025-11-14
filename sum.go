package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type sumNode struct {
	nodes []Node
}

func Sum(nodes ...Node) Node {
	if len(nodes) == 0 {
		return NodeError{
			NodeScalar: "sum",
			Cause:      fmt.Errorf("must receive at least one expression"),
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return NodeError{
				NodeScalar: "sum",
				Cause:      fmt.Errorf("all expressions should be not-nil"),
			}
		}
	}

	return sumNode{
		nodes: nodes,
	}
}

func (node sumNode) Scalar() string {
	return "sum"
}

func (node sumNode) Shape() []KeyNode {
	return sliceutils.Map(node.nodes, func(from Node) KeyNode { return KeyNode{Node: from} })
}

func (node sumNode) Type() NodeType {
	return NodeTypeExpression
}

func (node sumNode) Eval(scope Scope) Value {
	values := make([]any, 0, len(node.nodes))

	for i := range node.nodes {
		value, err := scope.Compute(node.nodes[i])
		if err != nil {
			return NewNodeError(node, err)
		}

		values = append(values, value)
	}

	sum, ok := sumAny(values...)
	if !ok {
		return NewNodeError(node, fmt.Errorf("all nodes must be of the same type"))
	}

	return Literal(sum)
}
