package goncoder

import (
	"bytes"
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
	parser := newParser(buffer)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode)
}

func translateNode(rootNode *Node) (gon.Expression, error) {
	switch rootNode.Type {
	case NodeTypeReference:
		return gon.Reference(string(rootNode.Scalar)), nil
	case NodeTypeLiteral:
		return gon.Static(rootNode.Value), nil
	}

	constructor, ok := DefaultExpressionCodex[string(rootNode.Scalar)]
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
			Key:        string(child.Key),
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
	Key      []byte
	Scalar   []byte
	Value    any
	Type     NodeType
}

type parser struct {
	tokens [][]byte
	index  int
}

func newParser(input []byte) *parser {
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
	if p.isNext([]byte(":")) {
		return nil, fmt.Errorf("invalid syntax")
	}

	// Function / object-style node
	if p.isNext([]byte("(")) {
		name := p.consume()
		p.consumeExpected([]byte("("))

		obj := &Node{
			Scalar:   name,
			Children: []*Node{},
		}

		for !p.done() && !bytes.Equal(p.peek(), []byte(")")) {
			if p.isNext([]byte(":")) {
				key := p.consume()
				p.consume() // skip ':'
				val, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				val.Key = key
				obj.Children = append(obj.Children, val)
				continue
			}

			// Unnamed child (e.g., call("reply"))
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			obj.Children = append(obj.Children, val)
		}

		p.consumeExpected([]byte(")"))

		return obj, nil
	}

	switch token := p.consume(); {
	case bytes.HasPrefix(token, []byte("\"")) && bytes.HasSuffix(token, []byte("\"")):
		val := bytes.Trim(token, "\"")
		return &Node{Value: string(val), Type: NodeTypeLiteral}, nil
	case isInteger(string(token)):
		integer, err := strconv.ParseInt(string(token), 10, 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  NodeTypeLiteral,
			Value: integer,
		}, nil
	case isFloat(string(token)):
		float, err := strconv.ParseFloat(string(token), 64)
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

func tokenize(input []byte) [][]byte {
	var tokens [][]byte

	var startPos int

	var inString bool
	var inComment bool

	for i, r := range input {
		resetCursor := func() {
			startPos = i + 1
		}

		getCurrent := func(inclusive bool) []byte {
			defer resetCursor()
			if inclusive {
				return input[startPos : i+1]
			}
			return input[startPos:i]
		}

		curLength := i - startPos

		switch {
		case r == '\n':
			if !inComment {
				if curLength > 0 {
					tokens = append(tokens, getCurrent(true))
				}
			}
			inComment = false
			resetCursor()
		case inComment:
		case bytes.Equal(input[i:i+2], []byte("//")):
			inComment = true
		case r == '"' && !inString:
			inString = true
		case r == '"' && inString:
			tokens = append(tokens, bytes.TrimSpace(getCurrent(true)))
			inString = false
		case inString:
		case unicode.IsSpace(rune(r)):
			if !inComment && curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			resetCursor()
		case bytes.Contains([]byte("():"), input[i:i+1]):
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			tokens = append(tokens, input[i:i+1])
			resetCursor()
		case r == ',':
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			resetCursor()
		default:
		}
	}
	if startPos < len(input) {
		tokens = append(tokens, input[startPos:])
	}
	return tokens
}

func (p *parser) peek() []byte {
	if p.index >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.index]
}

func (p *parser) isNext(next []byte) bool {
	if p.index+1 >= len(p.tokens) {
		return false
	}
	return bytes.Equal(p.tokens[p.index+1], next)
}

func (p *parser) consume() []byte {
	if p.index >= len(p.tokens) {
		return nil
	}
	t := p.tokens[p.index]
	p.index++
	return t
}

func (p *parser) consumeExpected(t []byte) {
	if bytes.Equal(p.peek(), t) {
		p.consume()
	}
}

func (p *parser) done() bool {
	return p.index >= len(p.tokens)
}
