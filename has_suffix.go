package gon

import (
	"fmt"
	"strings"
)

type hasSuffixNode struct {
	text   Node
	suffix Node
}

// Equal defines a suffix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the suffix.
func HasSuffix(text, suffix Node) Node {
	if text == nil || suffix == nil {
		return NodeError{
			NodeScalar: "suffix",
			Cause:      fmt.Errorf("all inputs should be not-nil"),
		}
	}

	return hasSuffixNode{
		text:   text,
		suffix: suffix,
	}
}

func (node hasSuffixNode) Scalar() string {
	return "hasSuffix"
}

func (node hasSuffixNode) Shape() []KeyNode {
	return []KeyNode{
		{"text", node.text},
		{"suffix", node.suffix},
	}
}

func (node hasSuffixNode) Type() NodeType {
	return NodeTypeExpression
}

func (node hasSuffixNode) Eval(scope Scope) Value {
	text, err := scope.Compute(node.text)
	if err != nil {
		return NewNodeError(node, err)
	}

	prefix, err := scope.Compute(node.suffix)
	if err != nil {
		return NewNodeError(node, err)
	}

	textStr, ok1 := text.(string)
	prefixStr, ok2 := prefix.(string)

	if !ok1 || !ok2 {
		return NewNodeError(node, fmt.Errorf("text and suffix should be string, got %T and %T", text, prefix))
	}

	return Literal(strings.HasSuffix(textStr, prefixStr))
}
