package gon

import "fmt"

type notNode struct {
	expression Node
}

// Not defines a not node, the input node should evaluate to boolean and be not-nil.
// It inverts the result of the evaluated boolean.
func Not(expression Node) Node {
	if expression == nil {
		return Literal(fmt.Errorf("not expression cannot be unset"))
	}

	return notNode{
		expression: expression,
	}
}

func (node notNode) Scalar() string {
	return "not"
}

func (node notNode) Shape() []KeyNode {
	return []KeyNode{
		{"expression", node.expression},
	}
}

func (node notNode) Type() NodeType {
	return NodeTypeExpression
}

func (node notNode) Eval(scope Scope) Value {
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
