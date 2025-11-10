package gon

import (
	"fmt"
)

type IfNode struct {
	condition  Node
	thenBranch Node
	elseBranch Node
}

func If(condition, thenBranch Node, elseBranch ...Node) Node {
	if condition == nil {
		return Literal(fmt.Errorf("if condition cannot be unset"))
	}

	return IfNode{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: safeGet(elseBranch, 0),
	}
}

func (node IfNode) Scalar() string {
	return "if"
}

func (node IfNode) Shape() []KeyNode {
	kv := []KeyNode{
		{"condition", node.condition},
		{"then", node.thenBranch},
	}
	if node.elseBranch != nil {
		kv = append(kv,
			KeyNode{"else", node.elseBranch},
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
				Scalar: node.Scalar(),
				Cause:  err,
			})
		}
		return Literal(NodeError{
			Scalar: node.Scalar(),
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
