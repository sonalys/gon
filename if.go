package gon

import "fmt"

type ifExpr struct {
	condition Expression
	expr      []Expression
}

func (e ifExpr) Name() (string, []KeyedExpression) {
	kv := []KeyedExpression{
		{Key: "condition", Value: e.condition},
		{Key: "then", Value: e.expr[0]},
	}
	if len(e.expr) > 1 {
		kv = append(kv, KeyedExpression{
			Key: "else", Value: e.expr[1],
		})
	}
	return "if", kv
}

func (e ifExpr) Type() ExpressionType {
	return ExpressionTypeOperation
}

func If(condition Expression, expr ...Expression) Expression {
	if len(expr) < 1 {
		return Static(fmt.Errorf("no branches specified for if condition"))
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
	conditionEval := i.condition.Eval(scope)
	fulfilled, ok := conditionEval.Value().(bool)
	if !ok {
		return Static(fmt.Errorf("condition should be bool, got %T", conditionEval.Value()))
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
