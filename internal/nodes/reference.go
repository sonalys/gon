package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

type ReferenceNode struct {
	definitionName string
}

func Reference(key string) adapters.Node {
	return &ReferenceNode{
		definitionName: key,
	}
}

func (node *ReferenceNode) Scalar() string {
	return node.definitionName
}

func (node *ReferenceNode) Shape() []adapters.KeyNode {
	return nil
}

func (node *ReferenceNode) Type() adapters.NodeType {
	return adapters.NodeTypeReference
}

func (node *ReferenceNode) Eval(scope adapters.Scope) adapters.Value {
	value, ok := scope.Definition(node.definitionName)
	if !ok {
		return adapters.NewNodeError(node, adapters.DefinitionNotFoundError{
			DefinitionKey: node.definitionName,
		})
	}

	return value.Eval(scope)
}

func (node *ReferenceNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
		}

		valued, ok := args[0].Node.(adapters.Valued)
		if !ok {
			return nil, fmt.Errorf("expected string literal")
		}

		referenceKey, ok := valued.Value().(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", referenceKey)
		}

		return Reference(referenceKey), nil
	})
}

var _ adapters.SerializableNode = &ReferenceNode{}
