package gon

type referenceNode struct {
	definitionName string
}

func Reference(key string) Node {
	return referenceNode{
		definitionName: key,
	}
}

func (node referenceNode) Scalar() string {
	return node.definitionName
}

func (node referenceNode) Shape() []KeyNode {
	return nil
}

func (node referenceNode) Type() NodeType {
	return NodeTypeReference
}

func (node referenceNode) Eval(scope Scope) Value {
	value, ok := scope.Definition(node.definitionName)
	if !ok {
		return NewNodeError(node, DefinitionNotFoundError{
			DefinitionKey: node.definitionName,
		})
	}

	return value.Eval(scope)
}
