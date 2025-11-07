package goncoder

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
	"github.com/sonalys/gon/internal/sliceutils"
)

type (
	NodeConstructor func(args []gon.KeyExpression) (gon.Expression, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)

func Encode(w io.Writer, root gon.Expression) error {
	astNode, err := ast.Parse(root)
	if err != nil {
		return fmt.Errorf("encoding root expression: %w", err)
	}
	return encodeBody(w, astNode, 0)
}

func encodeBody(w io.Writer, root ast.Node, indentation int) error {
	print := func(indentation int, mask string, args ...any) {
		if indentation > 0 {
			fmt.Fprint(w, strings.Repeat("\t", indentation))
		}
		fmt.Fprintf(w, mask, args...)
	}

	switch node := root.(type) {
	case ast.Expression:
		print(0, "%s(", node.Name)

		for i, arg := range node.KeyArgs {
			if i == 0 && arg.Key != "" || i > 0 {
				print(0, "\n")
				print(indentation+1, "")
			}
			if arg.Key != "" {
				print(0, "%s: ", arg.Key)
			}

			if err := encodeBody(w, arg.Node, indentation+1); err != nil {
				return err
			}
		}

		if len(node.KeyArgs) > 1 {
			print(0, "\n")
			print(indentation, ")")
			break
		}

		print(0, ")")
	case ast.Reference:
		print(0, "%v", node.Name)
	case ast.StaticValue:
		value := node.Value
		if str, ok := value.(string); ok {
			value = strconv.Quote(str)
		}
		print(0, "%v", value)
	default:
		return errors.New("cannot encode invalid expression type")
	}
	return nil
}

func argSorter(from []gon.KeyExpression, keys ...string) (map[string]gon.Expression, []gon.Expression, error) {
	if len(from) < len(keys) {
		return nil, nil, fmt.Errorf("missing arguments")
	}

	expectedMap := make(map[string]gon.Expression, len(keys))
	rest := make([]gon.Expression, 0, len(from))

gotArgLoop:
	for fromIndex := range from {
		for keyIndex := range keys {
			if from[fromIndex].Key == "" || from[fromIndex].Key == keys[keyIndex] {
				expectedMap[keys[keyIndex]] = from[fromIndex].Expression
				keys = slices.Delete(keys, keyIndex, keyIndex+1)
				continue gotArgLoop
			}
		}
		rest = append(rest, from[fromIndex].Expression)
	}

	return expectedMap, rest, nil
}

var DefaultExpressionCodex = Codex{
	"if": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, rest, err := argSorter(args, "condition", "then")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'if' node: %w", err)
		}

		return gon.If(orderedArgs["condition"], orderedArgs["then"], rest...), nil
	},
	"equal": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'equal' node: %w", err)
		}
		return gon.Equal(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"lt": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'lt' node: %w", err)
		}
		return gon.Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"lte": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'lte' node: %w", err)
		}
		return gon.SmallerOrEqual(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"gt": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'gt' node: %w", err)
		}
		return gon.Greater(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"gte": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'gte' node: %w", err)
		}
		return gon.GreaterOrEqual(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"not": func(args []gon.KeyExpression) (gon.Expression, error) {
		orderedArgs, _, err := argSorter(args, "expression")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'not' node: %w", err)
		}
		return gon.Not(orderedArgs["expression"]), nil
	},
	"call": func(args []gon.KeyExpression) (gon.Expression, error) {
		valuer := args[0].Expression.(gon.Valued)

		expressionTransform := func(from gon.KeyExpression) gon.Expression {
			return from.Expression
		}

		transformedArgs := sliceutils.Map(args[1:], expressionTransform)

		return gon.Call(valuer.Value().(string), transformedArgs...), nil
	},
	"time": func(args []gon.KeyExpression) (gon.Expression, error) {
		valuer := args[0].Expression.(gon.Valued)

		rawTime, ok := valuer.Value().(string)
		if !ok {
			return nil, fmt.Errorf("time should be parsed only from string")
		}

		t, err := time.Parse(time.RFC3339, rawTime)
		if err != nil {
			return nil, fmt.Errorf("time is invalid: %w", err)
		}

		return gon.Static(t), nil
	},
}

func Decode(buffer []byte, codex Codex) (gon.Expression, error) {
	tokens := tokenize(buffer)
	parser := newParser(tokens)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode, codex)
}
