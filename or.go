package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type OrNode struct {
	nodes []Node
}

// Or defines an or node, there must be at least one input.
// It returns the value of the first error, true boolean or non-boolean expression.
func Or(nodes ...Node) Node {
	if len(nodes) == 0 {
		return NodeError{
			NodeScalar: "or",
			Cause:      fmt.Errorf("must receive at least one expression"),
		}
	}

	for i := range nodes {
		if nodes[i] == nil {
			return NodeError{
				NodeScalar: "or",
				Cause:      ErrAllNodesMustBeSet,
			}
		}
	}

	return OrNode{
		nodes: nodes,
	}
}

func (node OrNode) Scalar() string {
	return "or"
}

func (node OrNode) Shape() []KeyNode {
	return sliceutils.Map(node.nodes, func(from Node) KeyNode { return KeyNode{Node: from} })
}

func (node OrNode) Type() NodeType {
	return NodeTypeExpression
}

func (node OrNode) Eval(scope Scope) Value {
	for _, expr := range node.nodes {
		value, err := scope.Compute(expr)
		if err != nil {
			return NewNodeError(node, err)
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

func (node OrNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		_, argsSlice, _ := argSorter(args)

		return Or(argsSlice...), nil
	})
}
