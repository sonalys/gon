package gon

import (
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

func Time(t string) static {
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

func (s static) Call(args ...Value) Value {
	valueOf := reflect.ValueOf(s.value)
	typeOf := valueOf.Type()

	if valueOf.Kind() != reflect.Func {
		return Static(fmt.Errorf("definition is not callable: %T", valueOf.Interface()))
	}

	if expArgs, gotArgs := typeOf.NumIn(), len(args); gotArgs != expArgs {
		return Static(fmt.Errorf("expected %d args, got %d", expArgs, gotArgs))
	}

	argsValue := make([]reflect.Value, 0, len(args))
	for i := range args {
		argsValue = append(argsValue, reflect.ValueOf(args[i].Value()))
	}

	resp := valueOf.Call(argsValue)

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

var (
	_ Value    = static{}
	_ Callable = static{}
)
