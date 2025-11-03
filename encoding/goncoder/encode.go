package goncoder

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

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
	"call": func(args []gon.KeyExpression) gon.Expression {
		valuer := args[0].Expression.(gon.Valuer)

		expressionTransform := func(from gon.KeyExpression) gon.Expression {
			return from.Expression
		}

		transformedArgs := sliceutils.Map(args[1:], expressionTransform)

		return gon.Call(valuer.Value().(string), transformedArgs...)
	},
}

func Decode(buffer []byte) (gon.Expression, error) {
	parser := newParser(string(buffer))

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode)
}

func translateNode(rootNode *Node) (gon.Expression, error) {
	switch rootNode.Type {
	case NodeTypeReference:
		return gon.Reference(rootNode.Scalar), nil
	case NodeTypeLiteral:
		return gon.Static(rootNode.Value), nil
	}

	constructor, ok := DefaultExpressionCodex[rootNode.Scalar]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	children := rootNode.Children
	nodeChildren := make([]gon.KeyExpression, 0, len(children))

	for _, child := range children {
		nodeChild, err := translateNode(child)
		if err != nil {
			return nil, err
		}
		nodeChildren = append(nodeChildren, gon.KeyExpression{
			Key:        child.Key,
			Expression: nodeChild,
		})
	}

	return constructor(nodeChildren), nil
}

type NodeType uint8

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeExpression
	NodeTypeReference
	NodeTypeLiteral
)

type Node struct {
	Children []*Node
	Key      string
	Scalar   string
	Value    any
	Type     NodeType
}

type parser struct {
	tokens []string
	index  int
}

func newParser(input string) *parser {
	// Strip comments
	lines := []string{}
	for l := range strings.SplitSeq(input, "\n") {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "//") {
			continue
		}
		lines = append(lines, l)
	}
	input = strings.Join(lines, " ")

	tokens := tokenize(input)
	return &parser{tokens: tokens}
}

func (p *parser) parse() (*Node, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no input")
	}

	return p.parseExpr()
}

func (p *parser) parseExpr() (*Node, error) {
	// Named field
	if p.peekNext() == ":" {
		return nil, fmt.Errorf("invalid syntax")
	}

	// Function / object-style node
	if p.peekNext() == "(" {
		name := p.consume()
		p.consumeExpected("(")

		obj := &Node{
			Scalar:   name,
			Children: []*Node{},
		}

		for !p.done() && p.peek() != ")" {
			if p.peekNext() == ":" {
				key := p.consume()
				p.consume() // skip ':'
				val, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				val.Key = key
				obj.Children = append(obj.Children, val)
			} else {
				// Unnamed child (e.g., call("reply"))
				val, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				obj.Children = append(obj.Children, val)
			}
		}

		p.consumeExpected(")")

		return obj, nil
	}

	switch token := p.consume(); {
	case strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\""):
		val := strings.Trim(token, "\"")
		return &Node{Value: val, Type: NodeTypeLiteral}, nil
	case isInteger(token):
		integer, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  NodeTypeLiteral,
			Value: integer,
		}, nil
	case isFloat(token):
		float, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  NodeTypeLiteral,
			Value: float,
		}, nil
	default:
		return &Node{
			Type:   NodeTypeReference,
			Scalar: token,
		}, nil
	}
}

func isInteger(s string) bool {
	if s == "" {
		return false
	}
	// Allow leading + or -
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	if s == "" {
		return false
	}
	// Allow leading + or -
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}

	dotSeen := false
	digitSeen := false
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			digitSeen = true
		case r == '.':
			if dotSeen {
				return false
			}
			dotSeen = true
		default:
			return false
		}
	}
	// Must contain at least one digit
	return digitSeen && dotSeen
}

func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inString := false

	for _, r := range input {
		switch {
		case r == '"' && !inString:
			inString = true
			current.WriteRune(r)
		case r == '"' && inString:
			current.WriteRune(r)
			tokens = append(tokens, strings.TrimSpace(current.String()))
			current.Reset()
			inString = false
		case inString:
			current.WriteRune(r)
		case unicode.IsSpace(r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case strings.ContainsRune("():", r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))
		// Ignore commas
		case r == ',':
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func (p *parser) peek() string {
	if p.index >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.index]
}

func (p *parser) peekNext() string {
	if p.index+1 >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.index+1]
}

func (p *parser) consume() string {
	if p.index >= len(p.tokens) {
		return ""
	}
	t := p.tokens[p.index]
	p.index++
	return t
}

func (p *parser) consumeExpected(t string) {
	if p.peek() == t {
		p.consume()
	}
}

func (p *parser) done() bool {
	return p.index >= len(p.tokens)
}
