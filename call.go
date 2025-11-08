package gon

import (
	"context"
)

type (
	callNode struct {
		funcName string
		args     []Expression
	}

	Callable interface {
		Expression
		Call(ctx context.Context, name string, values ...Value) Value
	}
)

func Call(callable string, args ...Expression) Expression {
	return callNode{
		funcName: callable,
		args:     args,
	}
}

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

func (node callNode) Eval(scope Scope) Value {
	values := make([]Value, 0, len(node.args))

	for i := range node.args {
		values = append(values, node.args[i].Eval(scope))
	}

	definition, ok := scope.Definition(node.funcName)
	if !ok {
		return Literal(NodeError{
			Scalar: node.Name(),
			Cause: DefinitionNotFoundError{
				DefinitionName: node.funcName,
			},
		})
	}

	callable, ok := definition.(Callable)
	if !ok {
		return Literal(NodeError{
			Scalar: node.Name(),
			Cause: DefinitionNotCallable{
				DefinitionName: node.funcName,
			},
		})
	}

	return callable.Call(scope, node.funcName, values...)
}
