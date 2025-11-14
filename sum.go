package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type SumNode struct {
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

	return SumNode{
		nodes: nodes,
	}
}

func (node SumNode) Scalar() string {
	return "sum"
}

func (node SumNode) Shape() []KeyNode {
	return sliceutils.Map(node.nodes, func(from Node) KeyNode { return KeyNode{Node: from} })
}

func (node SumNode) Type() NodeType {
	return NodeTypeExpression
}

func (node SumNode) Eval(scope Scope) Value {
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

func (node SumNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		_, rest, err := argSorter(args)
		if err != nil {
			return nil, err
		}

		return Sum(rest...), nil
	})
}
