package gon

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type (
	static struct {
		value any
	}
)

func (s static) Banner() (string, []KeyExpression) {
	switch v := s.value.(type) {
	case time.Time:
		return "time", []KeyExpression{
			{"", Static(v.Format(time.RFC3339))},
		}
	default:
		return "static", nil
	}
}

func (s static) Type() NodeType {
	switch s.value.(type) {
	case time.Time:
		return NodeTypeExpression
	default:
		return NodeTypeValue
	}
}

func Static(value any) static {
	return static{
		value: value,
	}
}

// Function receives a function in the format f(ctx, arg1, arg2, ...) (res1, res2, ...).
// Example:
//
//	gon.Function(func(ctx context.Context, name string) string)
func Function(f any) Expression {
	return static{
		value: f,
	}
}

func Time(t string) Expression {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return Static(err)
	}

	return static{
		value: parsed,
	}
}

func (s static) Value() any {
	if nested, ok := s.value.(Value); ok {
		return nested.Value()
	}

	return s.value
}

func (s static) Eval(scope Scope) Value {
	return s
}

func (s static) Call(ctx context.Context, args ...Value) Value {
	valueOfFunc := reflect.ValueOf(s.value)
	typeOfFunc := valueOfFunc.Type()

	if valueOfFunc.Kind() != reflect.Func {
		return Static(fmt.Errorf("definition is not callable: %T", valueOfFunc.Interface()))
	}

	if expArgs, gotArgs := typeOfFunc.NumIn(), len(args)+1; gotArgs != expArgs {
		return Static(fmt.Errorf("expected %d args, got %d", expArgs, gotArgs))
	}

	valueOfContext := reflect.ValueOf(ctx)
	if !valueOfContext.Type().AssignableTo(typeOfFunc.In(0)) {
		return Static(fmt.Errorf("definition first argument must be assignable to context"))

	}

	argsValue := make([]reflect.Value, 0, len(args)+1)
	argsValue = append(argsValue, valueOfContext)

	for i := range args {
		valueOfArg := reflect.ValueOf(args[i].Value())
		typeOfArg := valueOfArg.Type()
		expectedTypeOfArg := typeOfFunc.In(i + 1)

		if !typeOfArg.AssignableTo(expectedTypeOfArg) {
			return Static(fmt.Errorf("argument mismatch for function, arg %d expected %s, got %s", i+2, expectedTypeOfArg.String(), typeOfArg.String()))
		}

		argsValue = append(argsValue, valueOfArg)
	}

	resp := valueOfFunc.Call(argsValue)

	expResp := typeOfFunc.NumOut()
	if expResp == 0 {
		return Static(nil)
	}

	if typeOfFunc.NumOut() == 1 {
		return Static(resp[0].Interface())
	}

	respValue := make([]Value, 0, len(resp))
	for i := range resp {
		respValue = append(respValue, Static(resp[i].Interface()))
	}

	return Static(respValue)
}

var (
	_ Value    = static{}
	_ Callable = static{}
)
