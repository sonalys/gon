package gon

import "fmt"

type orNode struct {
	expressions []Expression
}

func (node orNode) Name() string {
	return "or"
}

func (node orNode) Shape() []KeyExpression {
	return nil
}

func (node orNode) Type() NodeType {
	return NodeTypeExpression
}

func Or(expressions ...Expression) Expression {
	if len(expressions) == 0 {
		return Literal(fmt.Errorf("if condition cannot be unset"))
	}

	return orNode{
		expressions: expressions,
	}
}

func (node orNode) Eval(scope Scope) Value {
	for _, expr := range node.expressions {
		switch value := expr.Eval(scope).Value().(type) {
		case error:
			return Literal(fmt.Errorf("evaluating or node: %w", value))
		case bool:
			if value {
				return Literal(true)
			}
		default:
			return Literal(value)
		}
	}

	return Literal(false)
}
