package gon

import "errors"

type assignment struct {
	definition string
	expression Expression
}

func (a assignment) Eval(scope Scope) Value {
	definer, ok := scope.(Definer)
	if !ok {
		return Literal(errors.New("scope is read-only"))
	}

	return Literal(definer.Define(a.definition, a.expression))
}
