package gon

type ReferenceNode struct {
	definitionName string
}

func Reference(key string) Node {
	return &ReferenceNode{
		definitionName: key,
	}
}

func (node *ReferenceNode) Scalar() string {
	return node.definitionName
}

func (node *ReferenceNode) Shape() []KeyNode {
	return nil
}

func (node *ReferenceNode) Type() NodeType {
	return NodeTypeReference
}

func (node *ReferenceNode) Eval(scope Scope) Value {
	value, ok := scope.Definition(node.definitionName)
	if !ok {
		return NewNodeError(node, DefinitionNotFoundError{
			DefinitionKey: node.definitionName,
		})
	}

	return value.Eval(scope)
}
