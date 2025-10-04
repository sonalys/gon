package gon

type not struct {
	expression Expression
}

func (e not) Name() (string, []KeyedExpression) {
	return "not", []KeyedExpression{
		{Key: "expression", Value: e.expression},
	}
}

func (e not) Type() ExpressionType {
	return ExpressionTypeOperation
}

func Not(expression Expression) not {
	return not{
		expression: expression,
	}
}

func (e not) Eval(scope Scope) Value {
	value, ok := e.expression.Eval(scope).Value().(bool)
	if !ok {
		return propagateErr(value, "cannot negate non-boolean expression")
	}

	return Static(!value)
}

var (
	_ Expression = not{}
)
