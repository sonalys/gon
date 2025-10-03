package gon

import (
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
		Slice() (value []Value, ok bool)
		Callable() (value Callable, ok bool)
	}

	Expression interface {
		Eval(scope Scope) Value
	}

	Definitions map[string]Expression

	Callable interface {
		Expression
		Call(...Value) Value
	}
)
