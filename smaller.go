package gon

import "errors"

type smaller struct {
	first  Expression
	second Expression
	equal  bool
}

func (e smaller) Name() (string, []KeyedExpression) {
	if e.equal {
		return "lte", []KeyedExpression{
			{Key: "first", Value: e.first},
			{Key: "second", Value: e.second},
		}
	}

	return "lt", []KeyedExpression{
		{Key: "first", Value: e.first},
		{Key: "second", Value: e.second},
	}
}

func (e smaller) Type() ExpressionType {
	return ExpressionTypeOperation
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
