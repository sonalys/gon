package gon

type (
	variable struct {
		key string
	}
)

func Variable(key string) Expression {
	return variable{
		key: key,
	}
}

func (v variable) Eval(scope Scope) Value {
	return scope.Variable(v.key).Eval(scope)
}
