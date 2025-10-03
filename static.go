package gon

import "time"

type (
	static struct {
		value any
	}
)

func (s static) Name() (string, []KeyedExpression) {
	switch s.value.(type) {
	case time.Time:
		return "time", nil
	default:
		return "static", nil
	}
}

func (s static) Type() ExpressionType {
	switch s.value.(type) {
	case time.Time:
		return ExpressionTypeOperation
	default:
		return ExpressionTypeValue
	}
}

func Static(value any) static {
	return static{
		value: value,
	}
}

func Time(t string) static {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return Static(err)
	}

	return static{
		value: parsed,
	}
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

func (s static) Slice() (value []Value, ok bool) {
	values, ok := s.value.([]Value)
	return values, ok
}

func (s static) Callable() (value Callable, ok bool) {
	value, ok = s.value.(Callable)
	return
}

var (
	_ Value = static{}
)
