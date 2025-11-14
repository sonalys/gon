package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/gonutils"
)

type NotNode struct {
	expression adapters.Node
}

// Not defines a not node, the input node should evaluate to boolean and be not-nil.
// It inverts the result of the evaluated boolean.
func Not(expression adapters.Node) adapters.Node {
	if expression == nil {
		return Literal(adapters.ErrAllNodesMustBeSet)
	}

	return &NotNode{
		expression: expression,
	}
}

func (node *NotNode) Scalar() string {
	return "not"
}

func (node *NotNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "expression", Node: node.expression},
	}
}

func (node *NotNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *NotNode) Eval(scope adapters.Scope) adapters.Value {
	value, err := scope.Compute(node.expression)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	resp, ok := value.(bool)
	if !ok {
		return adapters.NewNodeError(node, fmt.Errorf("expected bool got %T", value))
	}

	return Literal(!resp)
}

func (node *NotNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := gonutils.SortArgs(args, "expression")
		if err != nil {
			return nil, err
		}
		return Not(orderedArgs["expression"]), nil
	})
}

var _ adapters.SerializableNode = &NotNode{}
