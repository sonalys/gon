package gon

type once struct {
	eval  func(scope Scope) Value
	value Value
}

func Once(expression Expression) once {
	return once{
		eval: expression.Eval,
	}
}

func (e once) Eval(scope Scope) Value {
	if e.value == nil {
		e.value = e.eval(scope)
	}

	return e.value
}
