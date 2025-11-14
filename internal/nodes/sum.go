package nodes

import (
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/gonutils"
	"github.com/sonalys/gon/internal/sliceutils"
)

type SumNode struct {
	nodes []adapters.Node
}

func Sum(nodes ...adapters.Node) adapters.Node {
	if len(nodes) == 0 {
		return adapters.NodeError{
			NodeScalar: "sum",
			Cause:      adapters.ErrMustHaveArguments,
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return adapters.NodeError{
				NodeScalar: "sum",
				Cause:      adapters.ErrAllNodesMustBeSet,
			}
		}
	}

	return &SumNode{
		nodes: nodes,
	}
}

func (node *SumNode) Scalar() string {
	return "sum"
}

func (node *SumNode) Shape() []adapters.KeyNode {
	return sliceutils.Map(node.nodes, func(from adapters.Node) adapters.KeyNode { return adapters.KeyNode{Node: from} })
}

func (node *SumNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *SumNode) Eval(scope adapters.Scope) adapters.Value {
	values := make([]any, 0, len(node.nodes))

	for i := range node.nodes {
		value, err := scope.Compute(node.nodes[i])
		if err != nil {
			return adapters.NewNodeError(node, err)
		}

		values = append(values, value)
	}

	sum, ok := sumAny(values...)
	if !ok {
		return adapters.NewNodeError(node, adapters.ErrAllNodesMustMatch)
	}

	return Literal(sum)
}

func (node *SumNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		_, rest, err := gonutils.SortArgs(args)
		if err != nil {
			return nil, err
		}

		return Sum(rest...), nil
	})
}

var _ adapters.SerializableNode = &SumNode{}
