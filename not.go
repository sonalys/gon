package gon

type not struct {
	expression Expression
}

func (e not) Banner() (string, []KeyExpression) {
	return "not", []KeyExpression{
		KeyExpression{"expression", e.expression},
	}
}

func (e not) Type() NodeType {
	return NodeTypeExpression
}

func Not(expression Expression) not {
	return not{
		expression: expression,
	}
}

func (e not) Eval(scope Scope) Value {
	value := e.expression.Eval(scope)
	resp, ok := value.Value().(bool)
	if !ok {
		return propagateErr(value, "cannot negate non-boolean expression")
	}

	return Static(!resp)
}

var (
	_ Expression = not{}
)
