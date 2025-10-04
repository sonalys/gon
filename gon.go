package gon

import (
	"time"
)

type (
	ExpressionType uint8

	KeyedExpression struct {
		Key   string
		Value Expression
	}

	KeyedValue struct {
		Key   string
		Value Value
	}

	Typed interface {
		Type() ExpressionType
	}

	Named interface {
		Name() (string, []KeyedExpression)
	}

	Expression interface {
		Typed
		Named

		Eval(scope Scope) Value
	}

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
	}

	Definitions map[string]Expression

	Callable interface {
		Expression
		Call(...Value) Value
	}
)

const (
	ExpressionTypeInvalid ExpressionType = iota
	// ExpressionTypeOperation represents an operation() node type.
	ExpressionTypeOperation
	// ExpressionTypeReference represents a variable reference. Example: friend.name.
	ExpressionTypeReference
	// ExpressionTypeValue represents a direct value. Example: "string", 5.
	ExpressionTypeValue
)

/*

if(
	expression: equal(
		myName,
		friend.name
	),
	then: call(
		name: "reply",
		args: if(
			expression: greater(
				friend.age,
				18
			),
			then: "old"
			else: "young"
		),
	),
	else: call(
		"whoAreYou"
	)
)

*/
