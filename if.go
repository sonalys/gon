package gon

import (
	"fmt"

	"github.com/sonalys/gon/internal/sliceutils"
)

type ifExpr struct {
	condition Expression
	expr      []Expression
}

func (e ifExpr) Banner() (string, []KeyExpression) {
	kv := []KeyExpression{
		{"condition", e.condition},
		{"then", e.expr[0]},
	}
	if len(e.expr) > 1 {
		kv = append(kv,
			KeyExpression{"else", e.expr[1]},
		)
	}
	return "if", kv
}

func (e ifExpr) Type() NodeType {
	return NodeTypeExpression
}

func If(condition Expression, expr ...Expression) Expression {
	if condition == nil {
		return Static(fmt.Errorf("if condition cannot be unset"))
	}

	expr = sliceutils.Filter(expr, func(e Expression) bool { return e != nil })

	if len(expr) < 1 {
		return Static(fmt.Errorf("no set branches specified for if condition"))
	}

	if len(expr) > 2 {
		return Static(fmt.Errorf("if expression only accepts up to 2 expressions: main and alternative branches"))
	}

	return ifExpr{
		condition: condition,
		expr:      expr,
	}
}

func (i ifExpr) Eval(scope Scope) Value {
	value := i.condition.Eval(scope)
	fulfilled, ok := value.Value().(bool)
	if !ok {
		return propagateErr(value, "if expected bool, got %T", value.Value())
	}

	exprLen := len(i.expr)

	if fulfilled {
		return i.expr[0].Eval(scope)
	}

	if exprLen == 1 {
		return Static(nil)
	}

	return i.expr[1].Eval(scope)
}
