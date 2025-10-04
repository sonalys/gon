package gon

type (
	ExpressionType uint8

	Typed interface {
		Type() ExpressionType
	}

	KeyedExpression struct {
		Key   string
		Value Expression
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

		Value() any
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
