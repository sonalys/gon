package gon

import "fmt"

type equalNode struct {
	first  Expression
	second Expression
}

func (node equalNode) Name() string {
	return "equal"
}

func (node equalNode) Shape() []KeyExpression {
	return []KeyExpression{
		{"first", node.first},
		{"second", node.second},
	}
}

func (node equalNode) Type() NodeType {
	return NodeTypeExpression
}

func Equal(first, second Expression) Expression {
	if first == nil || second == nil {
		return Literal(fmt.Errorf("equal expression cannot compare unset expressions"))
	}

	return equalNode{
		first:  first,
		second: second,
	}
}

func (node equalNode) Eval(scope Scope) Value {
	firstValue := node.first.Eval(scope)
	secondValue := node.second.Eval(scope)

	return Literal(firstValue == secondValue)
}

var (
	_ Expression = equalNode{}
)
