package encoding

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
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

func Decode(buffer []byte, codex Codex) (gon.Expression, error) {
	tokens := tokenize(buffer)
	parser := newParser(tokens)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode, codex)
}
