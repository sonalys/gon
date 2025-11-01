package gon

import "fmt"

type equal struct {
	first  Expression
	second Expression
}

func (e equal) Banner() (string, []KeyExpression) {
	return "equal", []KeyExpression{
		{"first", e.first},
		{"second", e.second},
	}
}

func (e equal) Type() NodeType {
	return NodeTypeExpression
}

func Equal(first, second Expression) Expression {
	if first == nil || second == nil {
		return Static(fmt.Errorf("equal expression cannot compare unset expressions"))
	}

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
