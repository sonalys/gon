package nodes

import "github.com/sonalys/gon/adapters"

type EqualNode struct {
	first  adapters.Node
	second adapters.Node
}

// Equal defines an equality node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the inputs are equal.
func Equal(first, second adapters.Node) adapters.Node {
	if first == nil || second == nil {
		return adapters.NodeError{
			NodeScalar: "equal",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &EqualNode{
		first:  first,
		second: second,
	}
}

func (node *EqualNode) Scalar() string {
	return "equal"
}

func (node *EqualNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "first", Node: node.first},
		{Key: "second", Node: node.second},
	}
}

func (node *EqualNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *EqualNode) Eval(scope adapters.Scope) adapters.Value {
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	value, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return adapters.NewNodeError(node, adapters.IncompatiblePairError{
			First:  firstValue,
			Second: secondValue,
		})
	}

	return Literal(value == 0)
}

func (node *EqualNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Equal(orderedArgs["first"], orderedArgs["second"]), nil
	})
}
