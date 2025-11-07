package gon

type (
	Valued interface {
		Value() any
	}

	Value interface {
		Expression
		Valued
	}

	Typed interface {
		Type() NodeType
	}

	Named interface {
		Name() string
	}

	Shaped interface {
		Shape() []KeyExpression
	}

	Expression interface {
		Named
		Shaped

		Eval(scope Scope) Value
	}

	KeyExpression struct {
		Key        string
		Expression Expression
	}
)
