package gon

import (
	"fmt"
)

type smaller struct {
	first  Expression
	second Expression
	equal  bool
}

func Smaller(first, second Expression) smaller {
	return smaller{
		first:  first,
		second: second,
	}
}

func SmallerOrEqual(first, second Expression) smaller {
	return smaller{
		first:  first,
		second: second,
		equal:  true,
	}
}

func (e smaller) Eval(scope Scope) Value {
	firstValue := e.first.Eval(scope).Any()
	secondValue := e.second.Eval(scope).Any()

	comparison, ok := cmpAny(firstValue, secondValue)

	if !ok {
		return Static(fmt.Errorf("cannot compare different types: %T and %T", firstValue, secondValue))
	}

	if e.equal {
		return Static(comparison <= 0)
	}

	return Static(comparison < 0)
}

var (
	_ Expression = smaller{}
)
