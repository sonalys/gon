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

func Function(target any) Callable {
	valueOf := reflect.ValueOf(target)

	return function{
		target: valueOf,
	}
}

func (f function) Call(args ...Value) Value {
	if f.target.Kind() != reflect.Func {
		return Static(fmt.Errorf("target function is not callable: %T", f.target.Interface()))
	}

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
		return Static(resp[0])
	}

	respValue := make([]Value, 0, len(resp))
	for i := range resp {
		respValue = append(respValue, Static(resp[i]))
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

	definition := scope.Definition(c.callable).Eval(scope)

	callable, ok := definition.Eval(scope).Callable()
	if !ok {
		return Static(fmt.Errorf("definition is not callable: %T", definition.Any()))
	}

	return callable.Call(values...)
}
