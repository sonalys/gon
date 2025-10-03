package gon

type equal struct {
	first  Expression
	second Expression
}

func (e equal) Name() (string, []KeyedExpression) {
	return "equal", []KeyedExpression{
		{Key: "first", Value: e.first},
		{Key: "second", Value: e.second},
	}
}

func (e equal) Type() ExpressionType {
	return ExpressionTypeOperation
}

func Equal(first, second Expression) equal {
	return equal{
		first:  first,
		second: second,
	}
}

func (e equal) Eval(scope Scope) Value {
	firstValue := e.first.Eval(scope)
	secondValue := e.second.Eval(scope)

	return Static(firstValue == secondValue)
}

var (
	_ Expression = equal{}
)
