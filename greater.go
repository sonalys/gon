package gon

type (
	GreaterNode struct {
		first     Node
		second    Node
		inclusive bool
	}
)

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater than the second.
func Greater(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "gt",
			Cause:      ErrAllNodesMustBeSet,
		}
	}

	return GreaterNode{
		first:  first,
		second: second,
	}
}

// Greater defines a greater node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the first node is greater or equal than the second.
func GreaterOrEqual(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "gte",
			Cause:      ErrAllNodesMustBeSet,
		}
	}
	return GreaterNode{
		first:     first,
		second:    second,
		inclusive: true,
	}
}

func (node GreaterNode) Scalar() string {
	if node.inclusive {
		return "gte"
	}

	return "gt"
}

func (node GreaterNode) Shape() []KeyNode {
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

func (node GreaterNode) Type() NodeType {
	return NodeTypeExpression
}

func (node GreaterNode) Eval(scope Scope) Value {
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
		return Literal(comparison >= 0)
	}

	return Literal(comparison > 0)
}

func (node GreaterNode) Register(codex Codex) error {
	err := codex.Register("gt", func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Greater(orderedArgs["first"], orderedArgs["second"]), nil
	})
	if err != nil {
		return err
	}

	err = codex.Register("gte", func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Greater(orderedArgs["first"], orderedArgs["second"]), nil
	})

	return err
}
