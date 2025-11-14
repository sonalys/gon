package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

type IfNode struct {
	condition  adapters.Node
	thenBranch adapters.Node
	elseBranch adapters.Node
}

func If(condition, thenBranch adapters.Node, elseBranch ...adapters.Node) adapters.Node {
	if condition == nil {
		return adapters.NodeError{
			NodeScalar: "if",
			Cause:      fmt.Errorf("condition cannot be unset"),
		}
	}

	if len(elseBranch) > 1 {
		return adapters.NodeError{
			NodeScalar: "if",
			Cause:      fmt.Errorf("only one else branch can be set"),
		}
	}

	return &IfNode{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: safeGet(elseBranch, 0),
	}
}

func (node *IfNode) Scalar() string {
	return "if"
}

func (node *IfNode) Shape() []adapters.KeyNode {
	kv := []adapters.KeyNode{
		{"condition", node.condition},
		{"then", node.thenBranch},
	}
	if node.elseBranch != nil {
		kv = append(kv,
			adapters.KeyNode{"else", node.elseBranch},
		)
	}
	return kv
}

func (node *IfNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *IfNode) Eval(scope adapters.Scope) adapters.Value {
	value, err := scope.Compute(node.condition)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	fulfilled, ok := value.(bool)
	if !ok {
		return adapters.NewNodeError(node, fmt.Errorf("expected bool got %T", value))
	}

	if fulfilled {
		value, err := scope.Compute(node.thenBranch)
		if err != nil {
			return adapters.NewNodeError(node, err)
		}

		return Literal(value)
	}

	if node.elseBranch != nil {
		value, err := scope.Compute(node.elseBranch)
		if err != nil {
			return adapters.NewNodeError(node, err)
		}

		return Literal(value)
	}

	return Literal(false)
}

func (node *IfNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, rest, err := argSorter(args, "condition", "then")
		if err != nil {
			return nil, err
		}

		return If(orderedArgs["condition"], orderedArgs["then"], rest...), nil
	})
}
