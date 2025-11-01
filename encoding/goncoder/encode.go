package goncoder

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
	"github.com/sonalys/gon/internal/sliceutils"
)

type (
	Codex map[string]func(args []gon.KeyExpression) gon.Expression
)

func Encode(w io.Writer, root gon.Expression) error {
	astNode := ast.Parse(root)
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

func keyExpressionSplitter(arg gon.KeyExpression) (string, gon.Expression) {
	return arg.Key, arg.Expression
}

var DefaultExpressionCodex = Codex{
	"if": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.If(argMap["condition"], argMap["then"], argMap["else"])
	},
	"equal": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.Equal(argMap["first"], argMap["second"])
	},
	"lt": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.Smaller(argMap["first"], argMap["second"])
	},
	"lte": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.SmallerOrEqual(argMap["first"], argMap["second"])
	},
	"gt": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.Greater(argMap["first"], argMap["second"])
	},
	"gte": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.GreaterOrEqual(argMap["first"], argMap["second"])
	},
	"not": func(args []gon.KeyExpression) gon.Expression {
		argMap := sliceutils.ToMap(args, keyExpressionSplitter)
		return gon.Not(argMap["expression"])
	},
}

func Decode(r io.Reader) (gon.Expression, error) {
	return nil, nil
}
