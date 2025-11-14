package nodes

import "github.com/sonalys/gon/adapters"

type SmallerNode struct {
	first     adapters.Node
	second    adapters.Node
	inclusive bool
}

// Smaller defines a smaller node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller to the second.
func Smaller(first, second adapters.Node) adapters.Node {
	if first == nil || second == nil {
		return adapters.NodeError{
			NodeScalar: "lt",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &SmallerNode{
		first:  first,
		second: second,
	}
}

// SmallerOrEqual defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller or equal to the second.
func SmallerOrEqual(first, second adapters.Node) adapters.Node {
	if first == nil || second == nil {
		return adapters.NodeError{
			NodeScalar: "lte",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &SmallerNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node *SmallerNode) Scalar() string {
	if node.inclusive {
		return "lte"
	}

	return "lt"
}

func (node *SmallerNode) Shape() []adapters.KeyNode {
	if node.inclusive {
		return []adapters.KeyNode{
			{Key: "first", Node: node.first},
			{Key: "second", Node: node.second},
		}
	}

	return []adapters.KeyNode{
		{Key: "first", Node: node.first},
		{Key: "second", Node: node.second},
	}
}

func (node *SmallerNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *SmallerNode) Eval(scope adapters.Scope) adapters.Value {
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
		return Literal(comparison <= 0)
	}

	return Literal(comparison < 0)
}

func (node *SmallerNode) Register(codex adapters.Codex) error {
	err := codex.Register("lt", func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	})
	if err != nil {
		return err
	}

	return codex.Register("lte", func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	})
}

var (
	_ adapters.Node = &SmallerNode{}
)
