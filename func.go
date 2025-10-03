package gon

import (
	"fmt"
	"reflect"
)

type (
	function struct {
		target reflect.Value
	}

	call struct {
		callable string
		args     []Expression
	}
)

func (c call) Name() (string, []KeyedExpression) {
	kv := make([]KeyedExpression, 0, len(c.args)+1)
	kv = append(kv, KeyedExpression{Key: "", Value: Static(c.callable)})

	for i := range c.args {
		kv = append(kv, KeyedExpression{Key: "", Value: c.args[i]})
	}

	return "call", kv
}

func (c call) Type() ExpressionType {
	return ExpressionTypeOperation
}

func (f function) Name() (string, []KeyedExpression) {
	return "function", nil
}

func (f function) Type() ExpressionType {
	return ExpressionTypeInvalid
}

func Function(target any) Expression {
	valueOf := reflect.ValueOf(target)

	if valueOf.Kind() != reflect.Func {
		return Static(fmt.Errorf("target function is not callable: %T", valueOf.Interface()))
	}

	return function{
		target: valueOf,
	}
}

func (f function) Eval(scope Scope) Value {
	return Static(f)
}

func (f function) Call(args ...Value) Value {
	typeOf := f.target.Type()

	if expArgs, gotArgs := typeOf.NumIn(), len(args); gotArgs != expArgs {
		return Static(fmt.Errorf("expected %d args, got %d", expArgs, gotArgs))
	}

	argsValue := make([]reflect.Value, 0, len(args))
	for i := range args {
		argsValue = append(argsValue, reflect.ValueOf(args[i].Any()))
	}

	resp := f.target.Call(argsValue)

	expResp := typeOf.NumOut()
	if expResp == 0 {
		return Static(nil)
	}

	if typeOf.NumOut() == 1 {
		return Static(resp[0].Interface())
	}

	respValue := make([]Value, 0, len(resp))
	for i := range resp {
		respValue = append(respValue, Static(resp[i].Interface()))
	}

	return Static(respValue)
}

func Call(callable string, args ...Expression) Expression {
	return call{
		callable: callable,
		args:     args,
	}
}

func (c call) Eval(scope Scope) Value {
	values := make([]Value, 0, len(c.args))

	for i := range c.args {
		values = append(values, c.args[i].Eval(scope))
	}

	definition, ok := scope.Definition(c.callable)
	if !ok {
		return Static(fmt.Errorf("no callable definition found for %s", c.callable))
	}

	resp := definition.Eval(scope)

	callable, ok := definition.Eval(scope).Callable()
	if !ok {
		return Static(fmt.Errorf("definition is not callable: %T", resp.Any()))
	}

	return callable.Call(values...)
}
