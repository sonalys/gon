package gon

import "fmt"

type not struct {
	expression Expression
}

func (n not) Name() string {
	return "not"
}

func (n not) Shape() []KeyExpression {
	return []KeyExpression{
		{"expression", n.expression},
	}
}

func (n not) Type() NodeType {
	return NodeTypeExpression
}

func Not(expression Expression) Expression {
	if expression == nil {
		return Static(fmt.Errorf("not expression cannot be unset"))
	}

	return not{
		expression: expression,
	}
}

func (n not) Eval(scope Scope) Value {
	value := n.expression.Eval(scope)
	resp, ok := value.Value().(bool)
	if !ok {
		return propagateErr(value, "cannot negate non-boolean expression")
	}

	return Static(!resp)
}

var (
	_ Expression = not{}
)
