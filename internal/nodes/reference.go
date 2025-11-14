package nodes

import "github.com/sonalys/gon/adapters"

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
