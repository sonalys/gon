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

	if valueOf.Kind() != reflect.Struct || valueOf.Kind() != reflect.Map {
		return Static(errors.New("object can only be defined as pointer or value of struct or map"))
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

	var fieldValue reflect.Value

	switch o.valueOf.Kind() {
	case reflect.Struct:
		fieldValue = o.valueOf.FieldByName(topKey)
	case reflect.Map:
		fieldValue = o.valueOf.MapIndex(reflect.ValueOf(topKey))
	}

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
