package gon

import (
	"context"
	"errors"
	"time"
)

type (
	Value interface {
		Expression

		Any() any
		Bool() (value bool, ok bool)
		Duration() (value time.Duration, ok bool)
		Error() error
		Float() (value float64, ok bool)
		Int() (value int, ok bool)
		String() (value string, ok bool)
		Time() (value time.Time, ok bool)
	}

	Scope interface {
		context.Context
		Variable(key string) Expression
	}

	Expression interface {
		Eval(scope Scope) Value
	}

	static struct {
		value any
	}

	variable struct {
		key string
	}

	scope struct {
		context.Context
		store map[string]Expression
	}

	equal struct {
		first  Expression
		second Expression
	}
)

func (v variable) Eval(scope Scope) Value {
	return scope.Variable(v.key).Eval(scope)
}

func Variable(key string) Expression {
	return variable{
		key: key,
	}
}

func Equal(first, second Expression) Expression {
	return equal{
		first:  first,
		second: second,
	}
}

func (e equal) Eval(scope Scope) Value {
	firstValue := e.first.Eval(scope)
	secondValue := e.second.Eval(scope)

	return Static(firstValue == secondValue)
}

func Static(value any) Value {
	return static{
		value: value,
	}
}

func Context(ctx context.Context, fromMap map[string]Expression) Scope {
	return scope{
		Context: ctx,
		store:   fromMap,
	}
}

func (k scope) Variable(key string) Expression {
	value, ok := k.store[key]
	if !ok {
		return Static(errors.New("fact not found"))
	}

	return value
}

func (s static) Any() any {
	return s.value
}

func (s static) Bool() (value bool, ok bool) {
	value, ok = s.value.(bool)
	return
}

func (s static) Duration() (value time.Duration, ok bool) {
	value, ok = s.value.(time.Duration)
	return
}

func (s static) Error() error {
	err, _ := s.value.(error)
	return err
}

func (s static) Eval(scope Scope) Value {
	return s
}

func (s static) Float() (value float64, ok bool) {
	value, ok = s.value.(float64)
	return
}

func (s static) Int() (value int, ok bool) {
	value, ok = s.value.(int)
	return
}

func (s static) String() (value string, ok bool) {
	value, ok = s.value.(string)
	return
}

func (s static) Time() (value time.Time, ok bool) {
	value, ok = s.value.(time.Time)
	return
}

var (
	_ Value      = static{}
	_ Expression = equal{}
)
