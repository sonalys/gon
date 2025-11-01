package gon

import "fmt"

type (
	greater struct {
		first  Expression
		second Expression
		equal  bool
	}
)

func (e greater) Banner() (string, []KeyExpression) {
	if e.equal {
		return "gte", []KeyExpression{
			{"first", e.first},
			{"second", e.second},
		}
	}

	return "gt", []KeyExpression{
		{"first", e.first},
		{"second", e.second},
	}
}

func (e greater) Type() NodeType {
	return NodeTypeExpression
}

func Greater(first, second Expression) Expression {
	if first == nil || second == nil {
		return Static(fmt.Errorf("greater expression cannot compare unset expressions"))
	}

	return greater{
		first:  first,
		second: second,
	}
}

func GreaterOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Static(fmt.Errorf("greater or equal expression cannot compare unset expressions"))
	}
	return greater{
		first:  first,
		second: second,
		equal:  true,
	}
}

func (e greater) Eval(scope Scope) Value {
	firstValue := e.first.Eval(scope).Value()
	secondValue := e.second.Eval(scope).Value()

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return propagateErr(nil, "cannot compare different types: %T and %T", firstValue, secondValue)
	}

	if e.equal {
		return Static(comparison >= 0)
	}

	return Static(comparison > 0)
}

var (
	_ Expression = greater{}
)
