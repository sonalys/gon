package gon

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	objectNode struct {
		valueOf reflect.Value
	}
)

func (node objectNode) Name() string {
	return "object"
}

func (node objectNode) Shape() []KeyExpression {
	return nil
}

func (node objectNode) Type() NodeType {
	return NodeTypeInvalid
}

func Object(target any) Expression {
	valueOf := reflect.ValueOf(target)

	for valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}

	if valueOf.Kind() != reflect.Struct && valueOf.Kind() != reflect.Map {
		return Literal(errors.New("object can only be defined as pointer or value of struct or map"))
	}

	return &objectNode{
		valueOf: valueOf,
	}
}

func (node *objectNode) Eval(scope Scope) Value {
	return Literal(node.valueOf.Interface())
}

func (node *objectNode) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	var fieldValue reflect.Value

	switch node.valueOf.Kind() {
	case reflect.Struct:
		typeOf := node.valueOf.Type()
		fieldValue = node.valueOf.FieldByNameFunc(func(fieldName string) bool {
			field, ok := typeOf.FieldByName(fieldName)
			return ok && field.Tag.Get("gon") == key
		})
	case reflect.Map:
		fieldValue = node.valueOf.MapIndex(reflect.ValueOf(topKey))
	}

	if !fieldValue.IsValid() {
		return Literal(fmt.Errorf("definition not found: %s", key)), false
	}

	value := fieldValue.Interface()

	if len(parts) == 1 {
		return Literal(value), true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey):])
	}

	return propagateErr(nil, "definition not found: %s", key), false
}

var _ DefinitionResolver = &objectNode{}
