package gon

type SmallerNode struct {
	first     Node
	second    Node
	inclusive bool
}

// Smaller defines a smaller node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller to the second.
func Smaller(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "lt",
			Cause:      ErrAllNodesMustBeSet,
		}
	}

	return &SmallerNode{
		first:  first,
		second: second,
	}
}

// SmallerOrEqual defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is smaller or equal to the second.
func SmallerOrEqual(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "lte",
			Cause:      ErrAllNodesMustBeSet,
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

func (node *SmallerNode) Shape() []KeyNode {
	if node.inclusive {
		return []KeyNode{
			{"first", node.first},
			{"second", node.second},
		}
	}

	return []KeyNode{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node *SmallerNode) Type() NodeType {
	return NodeTypeExpression
}

func (node *SmallerNode) Eval(scope Scope) Value {
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return NewNodeError(node, err)
	}

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return NewNodeError(node, IncompatiblePairError{
			First:  firstValue,
			Second: secondValue,
		})
	}

	if node.inclusive {
		return Literal(comparison <= 0)
	}

	return Literal(comparison < 0)
}

func (node *SmallerNode) Register(codex Codex) error {
	err := codex.Register("lt", func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	})
	if err != nil {
		return err
	}

	return codex.Register("lte", func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	})
}

var (
	_ Node = &SmallerNode{}
)
