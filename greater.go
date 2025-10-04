package gon

import (
	"fmt"
)

type (
	greater struct {
		first  Expression
		second Expression
		equal  bool
	}
)

func (e greater) Name() (string, []KeyedExpression) {
	if e.equal {
		return "gte", []KeyedExpression{
			{Key: "first", Value: e.first},
			{Key: "second", Value: e.second},
		}
	}

	return "gt", []KeyedExpression{
		{Key: "first", Value: e.first},
		{Key: "second", Value: e.second},
	}
}

func (e greater) Type() ExpressionType {
	return ExpressionTypeOperation
}

func Greater(first, second Expression) greater {
	return greater{
		first:  first,
		second: second,
	}
}

func GreaterOrEqual(first, second Expression) greater {
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
		return Static(fmt.Errorf("cannot compare different types: %T and %T", firstValue, secondValue))
	}

	if e.equal {
		return Static(comparison >= 0)
	}

	return Static(comparison > 0)
}

var (
	_ Expression = greater{}
)
