package gon

type EqualNode struct {
	first  Node
	second Node
}

// Equal defines an equality node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the inputs are equal.
func Equal(first, second Node) Node {
	if first == nil || second == nil {
		return NodeError{
			NodeScalar: "equal",
			Cause:      ErrAllNodesMustBeSet,
		}
	}

	return EqualNode{
		first:  first,
		second: second,
	}
}

func (node EqualNode) Scalar() string {
	return "equal"
}

func (node EqualNode) Shape() []KeyNode {
	return []KeyNode{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node EqualNode) Type() NodeType {
	return NodeTypeExpression
}

func (node EqualNode) Eval(scope Scope) Value {
	firstValue, err := scope.Compute(node.first)
	if err != nil {
		return NewNodeError(node, err)
	}

	secondValue, err := scope.Compute(node.second)
	if err != nil {
		return NewNodeError(node, err)
	}

	value, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return NewNodeError(node, IncompatiblePairError{
			First:  firstValue,
			Second: secondValue,
		})
	}

	return Literal(value == 0)
}

func (node EqualNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, err
		}
		return Equal(orderedArgs["first"], orderedArgs["second"]), nil
	})
}
