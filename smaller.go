package gon

import (
	"errors"
	"fmt"
)

type smaller struct {
	first  Expression
	second Expression
	equal  bool
}

func (e smaller) Banner() (string, []KeyExpression) {
	if e.equal {
		return "lte", []KeyExpression{
			{"first", e.first},
			{"second", e.second},
		}
	}

	return "lt", []KeyExpression{
		{"first", e.first},
		{"second", e.second},
	}
}

func (e smaller) Type() NodeType {
	return NodeTypeExpression
}

func Smaller(first, second Expression) Expression {
	if first == nil || second == nil {
		return Static(fmt.Errorf("smaller expression cannot compare unset expressions"))
	}

	return smaller{
		first:  first,
		second: second,
	}
}

func SmallerOrEqual(first, second Expression) Expression {
	if first == nil || second == nil {
		return Static(fmt.Errorf("smaller or equal expression cannot compare unset expressions"))
	}

	return smaller{
		first:  first,
		second: second,
		equal:  true,
	}
}

func (e smaller) Eval(scope Scope) Value {
	firstValue := e.first.Eval(scope).Value()
	secondValue := e.second.Eval(scope).Value()

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		errs := make([]error, 0, 2)
		if err, ok := firstValue.(error); ok {
			errs = append(errs, err)
		}
		if err, ok := secondValue.(error); ok {
			errs = append(errs, err)
		}

		return propagateErr(Static(errors.Join(errs...)), "cannot compare different types: %T and %T", firstValue, secondValue)
	}

	if e.equal {
		return Static(comparison <= 0)
	}

	return Static(comparison < 0)
}

var (
	_ Expression = smaller{}
)
