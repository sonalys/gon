package gon

import "fmt"

type notNode struct {
	expression Expression
}

func Not(expression Expression) Expression {
	if expression == nil {
		return Literal(fmt.Errorf("not expression cannot be unset"))
	}

	return notNode{
		expression: expression,
	}
}

func (node notNode) Name() string {
	return "not"
}

func (node notNode) Shape() []KeyExpression {
	return []KeyExpression{
		{"expression", node.expression},
	}
}

func (node notNode) Type() NodeType {
	return NodeTypeExpression
}

func (node notNode) Eval(scope Scope) Value {
	value := node.expression.Eval(scope)
	resp, ok := value.Value().(bool)
	if !ok {
		return propagateErr(value, "cannot negate non-boolean expression")
	}

	return Literal(!resp)
}
