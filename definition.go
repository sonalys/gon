package gon

type (
	definition struct {
		key string
	}

	assignment struct {
		definition string
		expression Expression
	}
)

func Definition(key string) definition {
	return definition{
		key: key,
	}
}

func (d definition) Eval(scope Scope) Value {
	return scope.Definition(d.key).Eval(scope)
}

func (a assignment) Eval(scope Scope) Value {
	scope.Define(a.definition, a.expression)
	return Static(nil)
}
