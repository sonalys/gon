package gon

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type (
	literalNode struct {
		value any
	}
)

func (node literalNode) Name() string {
	switch node.value.(type) {
	case time.Time:
		return "time"
	default:
		return "literal"
	}
}

func (node literalNode) Shape() []KeyExpression {
	switch v := node.value.(type) {
	case time.Time:
		return []KeyExpression{
			{"", Literal(v.Format(time.RFC3339))},
		}
	default:
		return nil
	}
}

func (node literalNode) Type() NodeType {
	switch node.value.(type) {
	case time.Time:
		return NodeTypeExpression
	default:
		return NodeTypeValue
	}
}

func Literal(value any) literalNode {
	return literalNode{
		value: value,
	}
}

// Function receives a function in the format f(ctx, arg1, arg2, ...) (res1, res2, ...).
// Example:
//
//	gon.Function(func(ctx context.Context, name string) string)
func Function(f any) Expression {
	return literalNode{
		value: f,
	}
}

func Time(t string) Expression {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return Literal(err)
	}

	return literalNode{
		value: parsed,
	}
}

func (node literalNode) Value() any {
	if nested, ok := node.value.(Value); ok {
		return nested.Value()
	}

	return node.value
}

func (node literalNode) Eval(scope Scope) Value {
	return node
}

func (node literalNode) Call(ctx context.Context, args ...Value) Value {
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

var (
	_ Value    = literalNode{}
	_ Callable = literalNode{}
)
