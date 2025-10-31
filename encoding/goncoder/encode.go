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

var DefaultFactory = map[string]func(args []gon.KeyExpression) gon.Expression{
	"if": func(args []gon.KeyExpression) gon.Expression {
		return gon.If(args[0].Expression, sliceutils.Map(args[1:], func(from gon.KeyExpression) gon.Expression { return from.Expression })...)
	},
	"equal": func(args []gon.KeyExpression) gon.Expression {
		return gon.Equal(args[0].Expression, args[1].Expression)
	},
	"definition": func(args []gon.KeyExpression) gon.Expression {
		return gon.Reference(args[0].Key)
	},
}
