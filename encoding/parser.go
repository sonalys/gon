package encoding

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/sonalys/gon"
)

type parser struct {
	tokens []Token
	index  int
}

func newParser(tokens []Token) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) parse() (*Node, error) {
	if p.isNext(Token(":")) {
		return nil, fmt.Errorf("invalid syntax")
	}

	if p.isNext(Token("(")) {
		name := p.consume()
		p.consumeExpected(Token("("))

		node := &Node{
			Scalar:   name,
			Children: []*Node{},
		}

		for !p.done() && !bytes.Equal(p.peek(), Token(")")) {
			if p.isNext(Token(":")) {
				paramName := p.consume()
				p.consume() // skip ':'
				childNode, err := p.parse()
				if err != nil {
					return nil, err
				}
				childNode.Key = paramName
				node.Children = append(node.Children, childNode)
				continue
			}

			val, err := p.parse()
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, val)
		}

		p.consumeExpected(Token(")"))

		return node, nil
	}

	switch token := p.consume(); {
	case bytes.HasPrefix(token, Token("\"")) && bytes.HasSuffix(token, Token("\"")):
		val := bytes.Trim(token, "\"")
		return &Node{Value: string(val), Type: gon.NodeTypeLiteral}, nil
	case isInteger(string(token)):
		integer, err := strconv.ParseInt(string(token), 10, 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  gon.NodeTypeLiteral,
			Value: integer,
		}, nil
	case isFloat(string(token)):
		float, err := strconv.ParseFloat(string(token), 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  gon.NodeTypeLiteral,
			Value: float,
		}, nil
	default:
		switch string(token) {
		case "true", "True":
			return &Node{
				Type:  gon.NodeTypeLiteral,
				Value: true,
			}, nil
		case "false", "False":
			return &Node{
				Type:  gon.NodeTypeLiteral,
				Value: false,
			}, nil
		default:
			return &Node{
				Type:   gon.NodeTypeReference,
				Scalar: token,
			}, nil
		}
	}
}

func (p *parser) peek() Token {
	if p.index >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.index]
}

func (p *parser) isNext(next Token) bool {
	if p.index+1 >= len(p.tokens) {
		return false
	}
	return bytes.Equal(p.tokens[p.index+1], next)
}

func (p *parser) consume() Token {
	if p.index >= len(p.tokens) {
		return nil
	}
	t := p.tokens[p.index]
	p.index++
	return t
}

func (p *parser) consumeExpected(t Token) {
	if bytes.Equal(p.peek(), t) {
		p.consume()
	}
}

func (p *parser) done() bool {
	return p.index >= len(p.tokens)
}
