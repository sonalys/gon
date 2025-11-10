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
		value  reflect.Value
		isLazy bool
	}
)

var typeOfContext = reflect.TypeOf(context.Background())

// Literal represents a value/node.
// Use Literal with functions to define callable definitions.
// Use Literal with structs or maps to define definitions with children attributes.
// time.Time is serialized as time(RFC3339) by default.
func Literal(value any) *literalNode {
	valueOf := reflect.ValueOf(value)

	var isLazy bool
	if valueOf.IsValid() {
		typeOf := valueOf.Type()
		isLazy = typeOf.Kind() == reflect.Func && (typeOf.NumIn() == 0 || typeOfContext.AssignableTo(typeOf.In(0)))
	}

	return &literalNode{
		value:  valueOf,
		isLazy: isLazy,
	}
}

func (node *literalNode) Scalar() string {
	switch node.value.Interface().(type) {
	case time.Time:
		return "time"
	default:
		return "literal"
	}
}

func (node *literalNode) Shape() []KeyNode {
	switch v := node.value.Interface().(type) {
	case time.Time:
		return []KeyNode{
			{"", Literal(v.Format(time.RFC3339))},
		}
	default:
		return nil
	}
}

func (node *literalNode) Type() NodeType {
	switch node.value.Interface().(type) {
	case time.Time:
		return NodeTypeExpression
	default:
		return NodeTypeLiteral
	}
}

func (node *literalNode) Value() any {
	if node.value.IsValid() {
		if nested, ok := node.value.Interface().(Value); ok {
			return nested.Value()
		}
		return node.value.Interface()
	}

	return nil
}

func (node *literalNode) Eval(scope Scope) Value {
	if node.isLazy {
		return node.Call(scope, "")
	}
	return node
}

func (node *literalNode) Call(ctx context.Context, key string, args ...Value) Value {
	parts := strings.Split(key, ".")

	curValue := node.value
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

		if !curValue.IsValid() || curValue.IsZero() {
			return Literal(NodeError{
				Scalar: "literal",
				Cause:  fmt.Errorf("definition '%s' not found", strings.Join(parts[:i+1], ".")),
			})
		}
	}

	typeOfFunc := curValue.Type()

	if curValue.Kind() != reflect.Func {
		return Literal(NodeError{
			Scalar: "literal",
			Cause: DefinitionNotCallable{
				DefinitionName: key,
			},
		})
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
		return Literal(NodeError{
			Scalar: "literal",
			Cause:  fmt.Errorf("expected %d args, got %d", expArgs, gotArgs),
		})
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
			return Literal(NodeError{
				Scalar: "literal",
				Cause:  fmt.Errorf("argument mismatch for function, arg %d expected %s, got %s", targetParamIndex, expectedTypeOfArg.String(), typeOfArg.String()),
			})
		}

		argsValue = append(argsValue, valueOfArg)
	}

	resp := curValue.Call(argsValue)

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

func (node *literalNode) Definition(key string) (Value, bool) {
	parts := strings.Split(key, ".")

	curValue := node.value

	if !curValue.IsValid() || curValue.IsZero() {
		return Literal(nil), false
	}

	for i, partKey := range parts {
		// Pointer reference resolver, necessary to resolve pointer fields.
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

		if !curValue.IsValid() || curValue.IsZero() {
			partPath := strings.Join(parts[:i+1], ".")

			return Literal(DefinitionNotFoundError{
				DefinitionName: partPath,
			}), false
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
