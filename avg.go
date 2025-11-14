package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type AvgNode struct {
	nodes []Node
}

func Avg(nodes ...Node) Node {
	if len(nodes) == 0 {
		return NodeError{
			NodeScalar: "avg",
			Cause:      fmt.Errorf("must receive at least one expression"),
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return NodeError{
				NodeScalar: "avg",
				Cause:      fmt.Errorf("all expressions should be not-nil"),
			}
		}
	}

	return AvgNode{
		nodes: nodes,
	}
}

func (node AvgNode) Scalar() string {
	return "avg"
}

func (node AvgNode) Shape() []KeyNode {
	return sliceutils.Map(node.nodes, func(from Node) KeyNode { return KeyNode{Node: from} })
}

func (node AvgNode) Type() NodeType {
	return NodeTypeExpression
}

func (node AvgNode) Eval(scope Scope) Value {
	values := make([]any, 0, len(node.nodes))

	for i := range node.nodes {
		curValue, err := scope.Compute(node.nodes[i])
		if err != nil {
			return NewNodeError(node, err)
		}

		values = append(values, curValue)
	}

	sum, ok := avgAny(values...)
	if !ok {
		return NewNodeError(node, fmt.Errorf("all nodes must be of the same type"))
	}

	return Literal(sum)
}

func (node AvgNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		_, rest, err := argSorter(args)
		if err != nil {
			return nil, fmt.Errorf("error decoding 'avg' node: %w", err)
		}

		return Avg(rest...), nil
	})
}
