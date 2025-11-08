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
	NodeConstructor func(args []gon.KeyNode) (gon.Node, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)

func Encode(w io.Writer, root gon.Node) error {
	astNode, err := ast.Parse(root)
	if err != nil {
		return fmt.Errorf("encoding root expression: %w", err)
	}
	return encodeBody(w, astNode, 0)
}

func encodeBody(w io.Writer, root ast.Node, indentation int) error {
	print := func(indentation int, mask string, args ...any) error {
		if indentation > 0 {
			_, err := fmt.Fprint(w, strings.Repeat("\t", indentation))
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprintf(w, mask, args...)
		if err != nil {
			return err
		}

		return nil
	}

	switch node := root.(type) {
	case ast.Expression:
		if err := print(0, "%s(", node.Name); err != nil {
			return err
		}

		for i, arg := range node.KeyArgs {
			if len(node.KeyArgs) != 1 && (i == 0 && arg.Key != "" || i > 0) {
				if err := print(0, "\n"); err != nil {
					return err
				}
				if err := print(indentation+1, ""); err != nil {
					return err
				}
			}
			if arg.Key != "" {
				if err := print(0, "%s: ", arg.Key); err != nil {
					return err
				}
			}

			if err := encodeBody(w, arg.Node, indentation+1); err != nil {
				return err
			}
		}

		if len(node.KeyArgs) > 1 {
			if err := print(0, "\n"); err != nil {
				return err
			}
			if err := print(indentation, ")"); err != nil {
				return err
			}
			break
		}

		if err := print(0, ")"); err != nil {
			return err
		}
	case ast.Reference:
		if err := print(0, "%v", node.Name); err != nil {
			return err
		}
	case ast.StaticValue:
		value := node.Value
		if str, ok := value.(string); ok {
			value = strconv.Quote(str)
		}
		if err := print(0, "%v", value); err != nil {
			return err
		}
	default:
		return errors.New("cannot encode invalid expression type")
	}
	return nil
}

func Decode(buffer []byte, codex Codex) (gon.Node, error) {
	tokens := tokenize(buffer)
	parser := newParser(tokens)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode, codex)
}
