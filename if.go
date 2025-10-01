package gon

import "fmt"

type ifExpr struct {
	condition Expression
	expr      []Expression
}

func If(condition Expression, expr ...Expression) Expression {
	return ifExpr{
		condition: condition,
		expr:      expr,
	}
}

func (i ifExpr) Eval(scope Scope) Value {
	if len(i.expr) < 1 {
		return Static(fmt.Errorf("no branches specified for if condition"))
	}

	conditionEval := i.condition.Eval(scope)
	fulfilled, ok := conditionEval.Bool()
	if !ok {
		return Static(fmt.Errorf("condition should be bool, got %T", conditionEval.Any()))
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
