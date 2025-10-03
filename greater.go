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
	firstValue := e.first.Eval(scope).Any()
	secondValue := e.second.Eval(scope).Any()

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
