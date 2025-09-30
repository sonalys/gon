package gon

import "time"

type (
	static struct {
		value any
	}
)

func Static(value any) Value {
	return static{
		value: value,
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

var (
	_ Value = static{}
)
