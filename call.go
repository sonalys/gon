package gon

import (
	"context"
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type (
	CallNode struct {
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
	return &CallNode{
		funcName: funcName,
		argNodes: argNodes,
	}
}

func (node *CallNode) Scalar() string {
	return "call"
}

func (node *CallNode) Shape() []KeyNode {
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

func (node *CallNode) Type() NodeType {
	return NodeTypeExpression
}

func (node *CallNode) Eval(scope Scope) Value {
	values := make([]Value, 0, len(node.argNodes))

	for i := range node.argNodes {
		values = append(values, node.argNodes[i].Eval(scope))
	}

	definition, ok := scope.Definition(node.funcName)
	if !ok {
		return NewNodeError(node, DefinitionNotFoundError{
			DefinitionKey: node.funcName,
		})
	}

	callable, ok := definition.(Callable)
	if !ok {
		return NewNodeError(node, DefinitionNotCallableError{
			DefinitionKey: node.funcName,
		})
	}

	return callable.Call(scope, node.funcName, values...)
}

func (node *CallNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		funcName, ok := args[0].Node.(Valued).Value().(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", funcName)
		}

		expressionTransform := func(from KeyNode) Node {
			return from.Node
		}

		if len(args) == 1 {
			return Call(funcName), nil
		}

		transformedArgs := sliceutils.Map(args[1:], expressionTransform)

		return Call(funcName, transformedArgs...), nil
	})
}
