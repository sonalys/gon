package gon

import (
	"fmt"
)

type (
	call struct {
		callable string
		args     []Expression
	}
)

func (c call) Name() (string, []KeyedExpression) {
	kv := make([]KeyedExpression, 0, len(c.args)+1)
	kv = append(kv, KeyedExpression{Key: "", Value: Static(c.callable)})

	for i := range c.args {
		kv = append(kv, KeyedExpression{Key: "", Value: c.args[i]})
	}

	return "call", kv
}

func (c call) Type() ExpressionType {
	return ExpressionTypeOperation
}

func Call(callable string, args ...Expression) Expression {
	return call{
		callable: callable,
		args:     args,
	}
}

func (c call) Eval(scope Scope) Value {
	values := make([]Value, 0, len(c.args))

	for i := range c.args {
		values = append(values, c.args[i].Eval(scope))
	}

	definition, ok := scope.Definition(c.callable)
	if !ok {
		return Static(fmt.Errorf("no callable definition found for %s", c.callable))
	}

	callable, ok := definition.(Callable)
	if !ok {
		return Static(fmt.Errorf("definition is not callable: %s", c.callable))
	}

	return callable.Call(values...)
}
