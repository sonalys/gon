package gon

import (
	"errors"
	"reflect"
	"strings"
)

type (
	definition struct {
		key string
	}

	DefinitionResolver interface {
		Definition(key string) (Expression, bool)
	}

	definitionResolver map[string]Expression

	object struct {
		valueOf reflect.Value
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

type errorString string

const (
	DefinitionNotFoundErr errorString = "definition not found"
)

func (s errorString) Error() string {
	return string(s)
}

func (s definitionResolver) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := s[topKey]
	if !ok {
		return Static(DefinitionNotFoundErr), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return Static(errors.New("definition doesn't have children")), false
}

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
		return Static(DefinitionNotFoundErr), false
	}

	value := fieldValue.Interface()

	if len(parts) == 1 {
		return Static(value), true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey):])
	}

	return Static(DefinitionNotFoundErr), false
}
