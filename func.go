package gon

import (
	"context"
	"fmt"
)

type (
	callNode struct {
		funcName string
		args     []Expression
	}

	Callable interface {
		Expression
		Call(ctx context.Context, values ...Value) Value
	}
)

func (node callNode) Name() string {
	return "call"
}

func (node callNode) Shape() []KeyExpression {
	kv := make([]KeyExpression, 0, len(node.args)+1)
	kv = append(kv,
		KeyExpression{"", Literal(node.funcName)},
	)

	for i := range node.args {
		kv = append(kv,
			KeyExpression{"", node.args[i]},
		)
	}

	return kv
}

func (node callNode) Type() NodeType {
	return NodeTypeExpression
}

func Call(callable string, args ...Expression) Expression {
	return callNode{
		funcName: callable,
		args:     args,
	}
}

func (node callNode) Eval(scope Scope) Value {
	values := make([]Value, 0, len(node.args))

	for i := range node.args {
		values = append(values, node.args[i].Eval(scope))
	}

	definition, ok := scope.Definition(node.funcName)
	if !ok {
		return Literal(fmt.Errorf("no callable definition found for %s", node.funcName))
	}

	callable, ok := definition.(Callable)
	if !ok {
		return Literal(fmt.Errorf("definition is not callable: %s", node.funcName))
	}

	return callable.Call(scope, values...)
}
