package gon

type equal struct {
	first  Expression
	second Expression
}

func Equal(first, second Expression) Expression {
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
