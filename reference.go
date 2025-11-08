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
	value, ok := scope.Definition(node.definitionName)
	if !ok {
		if err, ok := value.Value().(error); ok {
			return Literal(NodeError{
				Scalar: "reference",
				Cause:  err,
			})

		}
		return Literal(NodeError{
			Scalar: "reference",
			Cause:  fmt.Errorf("definition not found: %s", node.definitionName),
		})
	}

	return value
}
