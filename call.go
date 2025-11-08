package gon

import (
	"context"
)

type (
	callNode struct {
		funcName string
		argNodes []Node
	}

	// Callable defines a node that can be called.
	// It represents a function as a node.
	Callable interface {
		Node
		Call(ctx context.Context, funcName string, argValues ...Value) Value
	}
)

// Call defines a function call.
// It tries to find the provided funcName under it's evaluated scope and call it with the given args.
// It will evaluate all args before providing them to the funcName.
// Returns a NodeError if the funcName is not found, callable or has wrong arguments.
// Context doesn't need to be given as an argument, and is handled automatically by gon.
func Call(funcName string, argNodes ...Node) Node {
	return callNode{
		funcName: funcName,
		argNodes: argNodes,
	}
}

func (node callNode) Name() string {
	return "call"
}

func (node callNode) Shape() []KeyNode {
	kv := make([]KeyNode, 0, len(node.argNodes)+1)
	kv = append(kv,
		KeyNode{"", Literal(node.funcName)},
	)

	for i := range node.argNodes {
		kv = append(kv,
			KeyNode{"", node.argNodes[i]},
		)
	}

	return kv
}

func (node callNode) Type() NodeType {
	return NodeTypeExpression
}

func (node callNode) Eval(scope Scope) Value {
	values := make([]Value, 0, len(node.argNodes))

	for i := range node.argNodes {
		values = append(values, node.argNodes[i].Eval(scope))
	}

	definition, ok := scope.Definition(node.funcName)
	if !ok {
		return Literal(NodeError{
			NodeName: node.Name(),
			Cause: DefinitionNotFoundError{
				DefinitionName: node.funcName,
			},
		})
	}

	callable, ok := definition.(Callable)
	if !ok {
		return Literal(NodeError{
			NodeName: node.Name(),
			Cause: DefinitionNotCallable{
				DefinitionName: node.funcName,
			},
		})
	}

	return callable.Call(scope, node.funcName, values...)
}
