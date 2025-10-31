package gon

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	object struct {
		valueOf reflect.Value
	}
)

func (e object) Banner() (string, []KeyExpression) {
	return "object", nil
}

func (e object) Type() NodeType {
	return NodeTypeInvalid
}

func Object(target any) Expression {
	valueOf := reflect.ValueOf(target)

	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Struct && valueOf.Kind() != reflect.Map {
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
		typeOf := o.valueOf.Type()
		fieldValue = o.valueOf.FieldByNameFunc(func(fieldName string) bool {
			field, ok := typeOf.FieldByName(fieldName)
			return ok && field.Tag.Get("gon") == key
		})
	case reflect.Map:
		fieldValue = o.valueOf.MapIndex(reflect.ValueOf(topKey))
	}

	if !fieldValue.IsValid() {
		return Static(fmt.Errorf("definition not found: %s", key)), false
	}

	value := fieldValue.Interface()

	if len(parts) == 1 {
		return Static(value), true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey):])
	}

	return propagateErr(nil, "definition not found: %s", key), false
}

var _ DefinitionResolver = &object{}
