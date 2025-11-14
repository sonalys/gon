package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/sliceutils"
)

type (
	CallNode struct {
		funcName string
		argNodes []adapters.Node
	}
)

// Call defines a function call.
// It tries to find the provided funcName under it's evaluated scope and call it with the given args.
// It will evaluate all args before providing them to the funcName.
// Returns a NodeError if the funcName is not found, callable or has wrong arguments.
// Context doesn't need to be given as an argument, and is handled automatically by nodes.
func Call(funcName string, argNodes ...adapters.Node) adapters.Node {
	return &CallNode{
		funcName: funcName,
		argNodes: argNodes,
	}
}

func (node *CallNode) Scalar() string {
	return "call"
}

func (node *CallNode) Shape() []adapters.KeyNode {
	kv := make([]adapters.KeyNode, 0, len(node.argNodes)+1)
	kv = append(kv,
		adapters.KeyNode{Key: "", Node: Literal(node.funcName)},
	)

	for i := range node.argNodes {
		kv = append(kv,
			adapters.KeyNode{Key: "", Node: node.argNodes[i]},
		)
	}

	return kv
}

func (node *CallNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *CallNode) Eval(scope adapters.Scope) adapters.Value {
	values := make([]adapters.Value, 0, len(node.argNodes))

	for i := range node.argNodes {
		values = append(values, node.argNodes[i].Eval(scope))
	}

	definition, ok := scope.Definition(node.funcName)
	if !ok {
		return adapters.NewNodeError(node, adapters.DefinitionNotFoundError{
			DefinitionKey: node.funcName,
		})
	}

	callable, ok := definition.(adapters.Callable)
	if !ok {
		return adapters.NewNodeError(node, adapters.DefinitionNotCallableError{
			DefinitionKey: node.funcName,
		})
	}

	return callable.Call(scope, node.funcName, values...)
}

func (node *CallNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		funcName, ok := args[0].Node.(adapters.Valued).Value().(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", funcName)
		}

		expressionTransform := func(from adapters.KeyNode) adapters.Node {
			return from.Node
		}

		if len(args) == 1 {
			return Call(funcName), nil
		}

		transformedArgs := sliceutils.Map(args[1:], expressionTransform)

		return Call(funcName, transformedArgs...), nil
	})
}
