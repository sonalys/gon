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

// HumanEncode encodes the node in a human-friendly format.
func HumanEncode(w io.Writer, root gon.Node, opts ...HumanEncodeOption) error {
	cfg := &humanEncodeConfig{
		showParamName: true,
	}

	for _, opt := range opts {
		opt.applyHumanEncodeOption(cfg)
	}

	astNode, err := ast.Parse(root)
	if err != nil {
		return fmt.Errorf("encoding root expression: %w", err)
	}

	return encodeBody(w, astNode, 0, cfg)
}

type humanEncodeConfig struct {
	compact       bool
	showParamName bool
}

type HumanEncodeOption interface {
	applyHumanEncodeOption(*humanEncodeConfig)
}

type prettyOpt struct{}

func (p prettyOpt) applyHumanEncodeOption(opt *humanEncodeConfig) {
	opt.compact = true
}

type hideParamName struct{}

func (p hideParamName) applyHumanEncodeOption(opt *humanEncodeConfig) {
	opt.showParamName = false
}

func Compact() *prettyOpt {
	return &prettyOpt{}
}

func Unnamed() *hideParamName {
	return &hideParamName{}
}

func encodeBody(w io.Writer, root ast.Node, indentation int, cfg *humanEncodeConfig) error {
	print := func(indentation int, mask string, args ...any) {
		if cfg.compact {
			indentation = 0
		}

		if indentation > 0 {
			_, _ = fmt.Fprint(w, strings.Repeat("\t", indentation))

		}

		_, _ = fmt.Fprintf(w, mask, args...)
	}

	endLine := func() {
		if cfg.compact {
			return
		}

		print(0, "\n")
	}

	switch node := root.(type) {
	case ast.Expression:
		print(0, "%s(", node.Scalar)

		for i, arg := range node.KeyArgs {
			if i > 0 || len(node.KeyArgs) != 1 || (i == 0 && len(node.KeyArgs) > 1 && arg.Key != "" && node.KeyArgs[i+1].Key != "") {
				endLine()
				print(indentation+1, "")
			}
			if cfg.showParamName && arg.Key != "" {
				print(0, "%s: ", arg.Key)
			}

			if err := encodeBody(w, arg.Node, indentation+1, cfg); err != nil {
				return err
			}

			if i < len(node.KeyArgs)-1 {
				print(0, ",")
			}
		}

		if len(node.KeyArgs) > 1 {
			endLine()
			print(indentation, ")")
			break
		}

		print(0, ")")
	case ast.Reference:
		print(0, "%v", node.Name)
	case ast.Literal:
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
