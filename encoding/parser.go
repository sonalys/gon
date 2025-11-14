package encoding

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/sonalys/gon/adapters"
)

type parser struct {
	tokens []Token
	index  int
}

func newParser(tokens []Token) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) parse() (*Node, error) {
	if p.isNext([]byte(":")) {
		return nil, fmt.Errorf("invalid syntax")
	}

	if p.isNext([]byte("(")) {
		name := p.consume()
		p.consumeExpected([]byte("("))

		node := &Node{
			Scalar:   name.content,
			Children: []*Node{},
		}

		for !p.done() && !bytes.Equal(p.peek().content, []byte(")")) {
			if p.isNext([]byte(":")) {
				paramName := p.consume()
				p.consume() // skip ':'
				childNode, err := p.parse()
				if err != nil {
					return nil, err
				}
				childNode.Key = paramName.content
				node.Children = append(node.Children, childNode)
				continue
			}

			val, err := p.parse()
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, val)
		}

		p.consumeExpected([]byte(")"))

		return node, nil
	}

	switch token := p.consume(); {
	case bytes.HasPrefix(token.content, []byte("\"")) && bytes.HasSuffix(token.content, []byte("\"")):
		val := bytes.Trim(token.content, "\"")
		return &Node{Value: string(val), Type: adapters.NodeTypeLiteral}, nil
	case isInteger(string(token.content)):
		integer, err := strconv.ParseInt(string(token.content), 10, 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  adapters.NodeTypeLiteral,
			Value: integer,
		}, nil
	case isFloat(string(token.content)):
		float, err := strconv.ParseFloat(string(token.content), 64)
		if err != nil {
			return nil, err
		}

		return &Node{
			Type:  adapters.NodeTypeLiteral,
			Value: float,
		}, nil
	default:
		switch string(token.content) {
		case "true", "True":
			return &Node{
				Type:  adapters.NodeTypeLiteral,
				Value: true,
			}, nil
		case "false", "False":
			return &Node{
				Type:  adapters.NodeTypeLiteral,
				Value: false,
			}, nil
		default:
			return &Node{
				Type:   adapters.NodeTypeReference,
				Scalar: token.content,
			}, nil
		}
	}
}

func (p *parser) peek() Token {
	if p.index >= len(p.tokens) {
		return Token{}
	}
	return p.tokens[p.index]
}

func (p *parser) isNext(next []byte) bool {
	if p.index+1 >= len(p.tokens) {
		return false
	}
	return bytes.Equal(p.tokens[p.index+1].content, next)
}

func (p *parser) consume() Token {
	if p.index >= len(p.tokens) {
		return Token{}
	}
	t := p.tokens[p.index]
	p.index++
	return t
}

func (p *parser) consumeExpected(t []byte) {
	if bytes.Equal(p.peek().content, t) {
		p.consume()
	}
}

func (p *parser) done() bool {
	return p.index >= len(p.tokens)
}
