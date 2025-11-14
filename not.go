package gon

import "fmt"

type NotNode struct {
	expression Node
}

// Not defines a not node, the input node should evaluate to boolean and be not-nil.
// It inverts the result of the evaluated boolean.
func Not(expression Node) Node {
	if expression == nil {
		return Literal(fmt.Errorf("not expression cannot be unset"))
	}

	return NotNode{
		expression: expression,
	}
}

func (node NotNode) Scalar() string {
	return "not"
}

func (node NotNode) Shape() []KeyNode {
	return []KeyNode{
		{"expression", node.expression},
	}
}

func (node NotNode) Type() NodeType {
	return NodeTypeExpression
}

func (node NotNode) Eval(scope Scope) Value {
	value, err := scope.Compute(node.expression)
	if err != nil {
		return NewNodeError(node, err)
	}

	resp, ok := value.(bool)
	if !ok {
		return NewNodeError(node, fmt.Errorf("expected bool got %T", value))
	}

	return Literal(!resp)
}

func (node NotNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "expression")
		if err != nil {
			return nil, err
		}
		return Not(orderedArgs["expression"]), nil
	})
}
