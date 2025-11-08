package gon

import (
	"fmt"
)

type IfNode struct {
	condition  Expression
	thenBranch Expression
	elseBranch Expression
}

func If(condition, thenBranch Expression, elseBranch ...Expression) Expression {
	if condition == nil {
		return Literal(fmt.Errorf("if condition cannot be unset"))
	}

	return IfNode{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: safeGet(elseBranch, 0),
	}
}

func (node IfNode) Name() string {
	return "if"
}

func (node IfNode) Shape() []KeyExpression {
	kv := []KeyExpression{
		{"condition", node.condition},
		{"then", node.thenBranch},
	}
	if node.elseBranch != nil {
		kv = append(kv,
			KeyExpression{"else", node.elseBranch},
		)
	}
	return kv
}

func (node IfNode) Type() NodeType {
	return NodeTypeExpression
}

func (node IfNode) Eval(scope Scope) Value {
	value := node.condition.Eval(scope)
	fulfilled, ok := value.Value().(bool)
	if !ok {
		if err, ok := value.Value().(error); ok {
			return Literal(NodeError{
				Scalar: node.Name(),
				Cause:  err,
			})
		}
		return Literal(NodeError{
			Scalar: node.Name(),
			Cause:  fmt.Errorf("expected a boolean value"),
		})
	}

	if fulfilled {
		return node.thenBranch.Eval(scope)
	}

	if node.elseBranch != nil {
		return node.elseBranch.Eval(scope)
	}

	return Literal(false)
}
