package nodes

import "github.com/sonalys/gon/adapters"

type (
	GreaterNode struct {
		first     adapters.Node
		second    adapters.Node
		inclusive bool
	}
)

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater than the second.
func Greater(first, second adapters.Node) adapters.Node {
	if first == nil || second == nil {
		return adapters.NodeError{
			NodeScalar: "gt",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &GreaterNode{
		first:  first,
		second: second,
	}
}

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater or equal than the second.
func GreaterOrEqual(first, second adapters.Node) adapters.Node {
	if first == nil || second == nil {
		return adapters.NodeError{
			NodeScalar: "gte",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}
	return &GreaterNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node *GreaterNode) Scalar() string {
	if node.inclusive {
		return "gte"
	}

	return "gt"
}

func (node *GreaterNode) Shape() []adapters.KeyNode {
	if node.inclusive {
		return []adapters.KeyNode{
			{"first", node.first},
			{"second", node.second},
		}
	}

	return []adapters.KeyNode{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node *GreaterNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *GreaterNode) Eval(scope adapters.Scope) adapters.Value {
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return adapters.NewNodeError(node, adapters.IncompatiblePairError{
			First:  firstValue,
			Second: secondValue,
		})
	}

	if node.inclusive {
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}

func (node *GreaterNode) Register(codex adapters.Codex) error {
	err := codex.Register("gt", func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Greater(orderedArgs["first"], orderedArgs["second"]), nil
	})
	if err != nil {
		return err
	}

	err = codex.Register("gte", func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Greater(orderedArgs["first"], orderedArgs["second"]), nil
	})

	return err
}
