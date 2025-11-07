package gon

import "fmt"

type (
	greater struct {
		first  Expression
		second Expression
		equal  bool
	}
)

func (g greater) Name() string {
	if g.equal {
		return "gte"
	}

	return "gt"
}

func (g greater) Shape() []KeyExpression {
	if g.equal {
		return []KeyExpression{
			{"first", g.first},
			{"second", g.second},
		}
	}

	return []KeyExpression{
		{"first", g.first},
		{"second", g.second},
	}
}

func (g greater) Type() NodeType {
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

func (g greater) Eval(scope Scope) Value {
	firstValue := g.first.Eval(scope).Value()
	secondValue := g.second.Eval(scope).Value()

	comparison, ok := cmpAny(firstValue, secondValue)
	if !ok {
		return propagateErr(nil, "cannot compare different types: %T and %T", firstValue, secondValue)
	}

	if g.equal {
		return Static(comparison >= 0)
	}

	return Static(comparison > 0)
}

var (
	_ Expression = greater{}
)
