package gon

import (
	"errors"
	"reflect"
	"strings"
)

type (
	object struct {
		valueOf reflect.Value
	}
)

func Object(target any) Expression {
	valueOf := reflect.ValueOf(target)

	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Struct {
		return Static(errors.New("object can only be defined as struct or pointer of struct"))
	}

	return &object{
		valueOf: valueOf,
	}
}

func (o *object) Eval(scope Scope) Value {
	return Static(o.valueOf.Interface())
}

func (o *object) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	fieldValue := o.valueOf.FieldByName(topKey)
	if !fieldValue.IsValid() {
		return Static(errors.New("definition not found")), false
	}

	value := fieldValue.Interface()

	if len(parts) == 1 {
		return Static(value), true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey):])
	}

	return Static(errors.New("definition not found")), false
}
