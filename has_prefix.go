package gon

import (
	"fmt"
	"strings"
)

type hasPrefixNode struct {
	text   Node
	prefix Node
}

// Equal defines a prefix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the prefix.
func HasPrefix(text, prefix Node) Node {
	if text == nil || prefix == nil {
		return Literal(NodeError{
			Scalar: "prefix",
			Cause:  fmt.Errorf("all inputs should be not-nil"),
		})
	}

	return hasPrefixNode{
		text:   text,
		prefix: prefix,
	}
}

func (node hasPrefixNode) Scalar() string {
	return "hasPrefix"
}

func (node hasPrefixNode) Shape() []KeyNode {
	return []KeyNode{
		{"text", node.text},
		{"prefix", node.prefix},
	}
}

func (node hasPrefixNode) Type() NodeType {
	return NodeTypeExpression
}

func (node hasPrefixNode) Eval(scope Scope) Value {
	text, err := scope.Compute(node.prefix)
	if err != nil {
		return Literal(NodeError{
			Scalar: node.Scalar(),
			Cause:  err,
		})
	}

	prefix, err := scope.Compute(node.prefix)
	if err != nil {
		return Literal(NodeError{
			Scalar: node.Scalar(),
			Cause:  err,
		})
	}

	textStr, ok1 := text.(string)
	prefixStr, ok2 := prefix.(string)

	if !ok1 || !ok2 {
		return Literal(NodeError{
			Scalar: node.Scalar(),
			Cause:  fmt.Errorf("text and prefix should be string, got %T and %T", text, prefix),
		})
	}

	return Literal(strings.HasPrefix(textStr, prefixStr))
}
