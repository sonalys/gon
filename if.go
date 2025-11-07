package gon

import (
	"fmt"
)

type IfNode struct {
	condition  Expression
	thenBranch Expression
	elseBranch Expression
}

func (node IfNode) Banner() (string, []KeyExpression) {
	kv := []KeyExpression{
		{"condition", node.condition},
		{"then", node.thenBranch},
	}
	if node.elseBranch != nil {
		kv = append(kv,
			KeyExpression{"else", node.elseBranch},
		)
	}
	return "if", kv
}

func (node IfNode) Type() NodeType {
	return NodeTypeExpression
}

func If(condition, thenBranch Expression, elseBranch ...Expression) Expression {
	if condition == nil {
		return Static(fmt.Errorf("if condition cannot be unset"))
	}

	return IfNode{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: tryGet(elseBranch, 0),
	}
}

func (node IfNode) Eval(scope Scope) Value {
	value := node.condition.Eval(scope)
	fulfilled, ok := value.Value().(bool)
	if !ok {
		return propagateErr(value, "if expected bool, got %T", value.Value())
	}

	if fulfilled {
		return node.thenBranch.Eval(scope)
	}

	if node.elseBranch != nil {
		return node.elseBranch.Eval(scope)
	}

	return Static(false)
}
