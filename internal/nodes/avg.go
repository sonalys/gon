package nodes

import (
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/sliceutils"
)

type AvgNode struct {
	nodes []adapters.Node
}

func Avg(nodes ...adapters.Node) adapters.Node {
	if len(nodes) == 0 {
		return adapters.NodeError{
			NodeScalar: "avg",
			Cause:      adapters.ErrMustHaveArguments,
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return adapters.NodeError{
				NodeScalar: "avg",
				Cause:      adapters.ErrAllNodesMustBeSet,
			}
		}
	}

	return &AvgNode{
		nodes: nodes,
	}
}

func (node *AvgNode) Scalar() string {
	return "avg"
}

func (node *AvgNode) Shape() []adapters.KeyNode {
	return sliceutils.Map(node.nodes, func(from adapters.Node) adapters.KeyNode { return adapters.KeyNode{Node: from} })
}

func (node *AvgNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *AvgNode) Eval(scope adapters.Scope) adapters.Value {
	values := make([]any, 0, len(node.nodes))

	for i := range node.nodes {
		curValue, err := scope.Compute(node.nodes[i])
		if err != nil {
			return adapters.NewNodeError(node, err)
		}

		values = append(values, curValue)
	}

	sum, ok := avgAny(values...)
	if !ok {
		return adapters.NewNodeError(node, adapters.ErrAllNodesMustMatch)
	}

	return Literal(sum)
}

func (node *AvgNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		_, rest, err := argSorter(args)
		if err != nil {
			return nil, err
		}

		return Avg(rest...), nil
	})
}
