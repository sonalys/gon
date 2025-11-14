package nodes

import (
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/sliceutils"
)

type OrNode struct {
	nodes []adapters.Node
}

// Or defines an or node, there must be at least one input.
// It returns the value of the first error, true boolean or non-boolean expression.
func Or(nodes ...adapters.Node) adapters.Node {
	if len(nodes) == 0 {
		return adapters.NodeError{
			NodeScalar: "or",
			Cause:      adapters.ErrMustHaveArguments,
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return adapters.NodeError{
				NodeScalar: "or",
				Cause:      adapters.ErrAllNodesMustBeSet,
			}
		}
	}

	return &OrNode{
		nodes: nodes,
	}
}

func (node *OrNode) Scalar() string {
	return "or"
}

func (node *OrNode) Shape() []adapters.KeyNode {
	return sliceutils.Map(node.nodes, func(from adapters.Node) adapters.KeyNode { return adapters.KeyNode{Node: from} })
}

func (node *OrNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *OrNode) Eval(scope adapters.Scope) adapters.Value {
	for _, expr := range node.nodes {
		value, err := scope.Compute(expr)
		if err != nil {
			return adapters.NewNodeError(node, err)
		}

		switch value := value.(type) {
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

func (node *OrNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		_, argsSlice, _ := argSorter(args)

		return Or(argsSlice...), nil
	})
}
