package nodes

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sonalys/gon/adapters"
)

type (
	LiteralNode struct {
		value  reflect.Value
		isLazy bool
	}
)

var typeOfContext = reflect.TypeOf(context.Background())

// Literal represents a value/node.
// Use Literal with functions to define callable definitions.
// Use Literal with structs or maps to define definitions with children attributes.
// time.Time is serialized as time(RFC3339) by default.
func Literal(value any) *LiteralNode {
	valueOf := reflect.ValueOf(value)

	var isLazy bool
	if valueOf.IsValid() {
		typeOf := valueOf.Type()
		isLazy = typeOf.Kind() == reflect.Func && (typeOf.NumIn() == 0 || typeOfContext.AssignableTo(typeOf.In(0)))
	}

	return &LiteralNode{
		value:  valueOf,
		isLazy: isLazy,
	}
}

func (node *LiteralNode) Scalar() string {
	if !node.value.IsValid() || !node.value.CanInterface() {
		return "literal"
	}

	switch node.value.Interface().(type) {
	case time.Time:
		return "time"
	default:
		return "literal"
	}
}

func (node *LiteralNode) Shape() []adapters.KeyNode {
	if !node.value.IsValid() || !node.value.CanInterface() {
		return []adapters.KeyNode{
			{Node: node},
		}
	}

	switch v := node.value.Interface().(type) {
	case time.Time:
		return []adapters.KeyNode{
			{Key: "", Node: Literal(v.Format(time.RFC3339))},
		}
	default:
		return []adapters.KeyNode{
			{Node: node},
		}
	}
}

func (node *LiteralNode) Type() adapters.NodeType {
	switch node.value.Interface().(type) {
	case time.Time:
		return adapters.NodeTypeExpression
	default:
		return adapters.NodeTypeLiteral
	}
}

func (node *LiteralNode) Value() any {
	if node.value.IsValid() {
		if nested, ok := node.value.Interface().(adapters.Value); ok {
			return nested.Value()
		}
		return node.value.Interface()
	}

	return nil
}

func (node *LiteralNode) Eval(scope adapters.Scope) adapters.Value {
	if node.isLazy {
		return node.Call(scope, "")
	}
	return node
}

func (node *LiteralNode) Call(ctx context.Context, key string, args ...adapters.Value) adapters.Value {
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
			return adapters.NewNodeError(node, adapters.DefinitionNotFoundError{
				DefinitionKey: strings.Join(parts[:i+1], "."),
			})
		}
	}

	typeOfFunc := curValue.Type()

	if curValue.Kind() != reflect.Func {
		return adapters.NewNodeError(node, adapters.DefinitionNotCallableError{
			DefinitionKey: key,
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
		return adapters.NewNodeError(node, fmt.Errorf("expected %d args, got %d", expArgs, gotArgs))
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
			return adapters.NewNodeError(node, fmt.Errorf("argument mismatch for function, arg %d expected %s, got %s", targetParamIndex, expectedTypeOfArg.String(), typeOfArg.String()))
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

	respValue := make([]adapters.Value, 0, len(resp))
	for i := range resp {
		respValue = append(respValue, Literal(resp[i].Interface()))
	}

	return Literal(respValue)
}

func (node *LiteralNode) Definition(key string) (adapters.Value, bool) {
	parts := strings.Split(key, ".")

	curValue := node.value

	if !curValue.IsValid() {
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

		if !curValue.IsValid() {
			partPath := strings.Join(parts[:i+1], ".")

			return Literal(adapters.DefinitionNotFoundError{
				DefinitionKey: partPath,
			}), false
		}
	}

	value := curValue.Interface()
	return Literal(value), true
}

func (node *LiteralNode) Register(codex adapters.Codex) error {
	err := codex.Register("time", func(args []adapters.KeyNode) (adapters.Node, error) {
		valuer, ok := args[0].Node.(adapters.Valued)
		if !ok {
			return nil, fmt.Errorf("invalid value received")
		}

		rawTime, ok := valuer.Value().(string)
		if !ok {
			return nil, fmt.Errorf("time should be parsed only from string")
		}

		t, err := time.Parse(time.RFC3339, rawTime)
		if err != nil {
			return nil, fmt.Errorf("time is invalid: %w", err)
		}

		return Literal(t), nil
	})
	if err != nil {
		return err
	}

	err = codex.Register("bool", func(args []adapters.KeyNode) (adapters.Node, error) {
		valuer, ok := args[0].Node.(adapters.Valued)
		if !ok {
			return nil, fmt.Errorf("invalid value received")
		}

		raw, ok := valuer.Value().(string)
		if !ok {
			return nil, fmt.Errorf("bool should be parsed only from string")
		}

		t, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("time is invalid: %w", err)
		}

		return Literal(t), nil
	})
	if err != nil {
		return err
	}

	err = codex.Register("literal", func(args []adapters.KeyNode) (adapters.Node, error) {
		valuer, ok := args[0].Node.(adapters.Valued)
		if !ok {
			return nil, fmt.Errorf("invalid value received")
		}

		return Literal(valuer.Value()), nil
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	_ adapters.Value            = &LiteralNode{}
	_ adapters.Callable         = &LiteralNode{}
	_ adapters.DefinitionReader = &LiteralNode{}
	_ adapters.SerializableNode = &LiteralNode{}
)
