package gon

import "fmt"

type referenceNode struct {
	definitionName string
}

func Reference(key string) Expression {
	return referenceNode{
		definitionName: key,
	}
}

func (node referenceNode) Name() string {
	return node.definitionName
}

func (node referenceNode) Shape() []KeyExpression {
	return nil
}

func (node referenceNode) Type() NodeType {
	return NodeTypeReference
}

func (node referenceNode) Eval(scope Scope) Value {
	expression, ok := scope.Definition(node.definitionName)
	if !ok {
		return Literal(fmt.Errorf("definition not found: %s", node.definitionName))
	}

	return expression.Eval(scope)
}
