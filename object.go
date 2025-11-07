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

	curValue := node.valueOf
	for i, partKey := range parts {
		switch curValue.Kind() {
		case reflect.Pointer:
			curValue = curValue.Elem()
		case reflect.Struct:
			typeOf := curValue.Type()
			curValue = curValue.FieldByNameFunc(func(fieldName string) bool {
				field, ok := typeOf.FieldByName(fieldName)
				return ok && field.Tag.Get("gon") == partKey
			})
		case reflect.Map:
			curValue = curValue.MapIndex(reflect.ValueOf(partKey))
		}

		if curValue.IsZero() {
			return Literal(fmt.Errorf("definition not found: %s", strings.Join(parts[:i+1], "."))), false
		}
	}

	value := curValue.Interface()
	return Literal(value), true
}

var _ DefinitionResolver = &objectNode{}
