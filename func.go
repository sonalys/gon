package gon

import (
	"context"
	"fmt"
)

type (
	call struct {
		callable string
		args     []Expression
	}

	Callable interface {
		Expression
		Call(ctx context.Context, values ...Value) Value
	}
)

func (c call) Banner() (string, []KeyExpression) {
	kv := make([]KeyExpression, 0, len(c.args)+1)
	kv = append(kv,
		KeyExpression{"", Static(c.callable)},
	)

	for i := range c.args {
		kv = append(kv,
			KeyExpression{"", c.args[i]},
		)
	}

	return "call", kv
}

func (c call) Type() NodeType {
	return NodeTypeExpression
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

	return callable.Call(scope, values...)
}
