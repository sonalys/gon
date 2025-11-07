package gon

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type (
	literalNode struct {
		value any
	}
)

func Literal(value any) Value {
	return &literalNode{
		value: value,
	}
}

func (node *literalNode) Name() string {
	switch node.value.(type) {
	case time.Time:
		return "time"
	default:
		return "literal"
	}
}

func (node *literalNode) Shape() []KeyExpression {
	switch v := node.value.(type) {
	case time.Time:
		return []KeyExpression{
			{"", Literal(v.Format(time.RFC3339))},
		}
	default:
		return nil
	}
}

func (node *literalNode) Type() NodeType {
	switch node.value.(type) {
	case time.Time:
		return NodeTypeExpression
	default:
		return NodeTypeValue
	}
}

func (node *literalNode) Value() any {
	if nested, ok := node.value.(Value); ok {
		return nested.Value()
	}

	return node.value
}

func (node *literalNode) Eval(scope Scope) Value {
	return node
}

func (node *literalNode) Call(ctx context.Context, args ...Value) Value {
	valueOfFunc := reflect.ValueOf(node.value)
	typeOfFunc := valueOfFunc.Type()

	if valueOfFunc.Kind() != reflect.Func {
		return Literal(fmt.Errorf("definition is not callable: %T", valueOfFunc.Interface()))
	}

	expArgs := typeOfFunc.NumIn()
	gotArgs := len(args)
	var gotContext bool
	if expArgs > 0 {
		if reflect.TypeOf(ctx).AssignableTo(typeOfFunc.In(0)) {
			gotArgs += 1
			gotContext = true
		}
	}

	if gotArgs != expArgs {
		return Literal(fmt.Errorf("expected %d args, got %d", expArgs, gotArgs))
	}

	argsValue := make([]reflect.Value, 0, expArgs)
	if gotContext {
		argsValue = append(argsValue, reflect.ValueOf(ctx))
	}

	for i := range args {
		valueOfArg := reflect.ValueOf(args[i].Value())
		typeOfArg := valueOfArg.Type()

		targetParamIndex := i
		if gotContext {
			targetParamIndex += 1
		}
		expectedTypeOfArg := typeOfFunc.In(targetParamIndex)

		if !typeOfArg.AssignableTo(expectedTypeOfArg) {
			return Literal(fmt.Errorf("argument mismatch for function, arg %d expected %s, got %s", targetParamIndex, expectedTypeOfArg.String(), typeOfArg.String()))
		}

		argsValue = append(argsValue, valueOfArg)
	}

	resp := valueOfFunc.Call(argsValue)

	expResp := typeOfFunc.NumOut()
	if expResp == 0 {
		return Literal(nil)
	}

	if typeOfFunc.NumOut() == 1 {
		return Literal(resp[0].Interface())
	}

	respValue := make([]Value, 0, len(resp))
	for i := range resp {
		respValue = append(respValue, Literal(resp[i].Interface()))
	}

	return Literal(respValue)
}

func (node *literalNode) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")

	valueOf := reflect.ValueOf(node.value)

	switch valueOf.Kind() {
	case reflect.Struct, reflect.Map, reflect.Pointer:
	default:
		return Literal(fmt.Errorf("literal of type %T cannot define children attributes", node.value)), false
	}

	curValue := valueOf
	for i, partKey := range parts {
		// Pointer resolver.
		for ; curValue.Kind() == reflect.Pointer; curValue = curValue.Elem() {
		}
		switch curValue.Kind() {
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

var (
	_ Value              = &literalNode{}
	_ Callable           = &literalNode{}
	_ DefinitionResolver = &literalNode{}
)
