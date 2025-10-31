package gon

type (
	NodeType uint8

	Typed interface {
		Type() NodeType
		Banner() (string, []KeyExpression)
	}

	Expression interface {
		Typed

		Eval(scope Scope) Value
	}

	KeyExpression struct {
		Key        string
		Expression Expression
	}

	Valuer interface {
		Value() any
	}

	Value interface {
		Expression
		Valuer
	}
)

const (
	NodeTypeInvalid NodeType = iota
	// NodeTypeExpression represents an expression() node type. Example: if()
	NodeTypeExpression
	// NodeTypeReference represents a variable reference. Example: friend.name.
	NodeTypeReference
	// NodeTypeValue represents a direct value. Example: "string", 5.
	NodeTypeValue
)
