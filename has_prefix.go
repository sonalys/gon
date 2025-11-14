package gon

import (
	"fmt"
	"strings"
)

type HasPrefixNode struct {
	text   Node
	prefix Node
}

// Equal defines a prefix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the prefix.
func HasPrefix(text, prefix Node) Node {
	if text == nil || prefix == nil {
		return NodeError{
			NodeScalar: "prefix",
			Cause:      fmt.Errorf("all inputs should be not-nil"),
		}
	}

	return HasPrefixNode{
		text:   text,
		prefix: prefix,
	}
}

func (node HasPrefixNode) Scalar() string {
	return "hasPrefix"
}

func (node HasPrefixNode) Shape() []KeyNode {
	return []KeyNode{
		{"text", node.text},
		{"prefix", node.prefix},
	}
}

func (node HasPrefixNode) Type() NodeType {
	return NodeTypeExpression
}

func (node HasPrefixNode) Eval(scope Scope) Value {
	text, err := scope.Compute(node.text)
	if err != nil {
		return NewNodeError(node, err)
	}

	prefix, err := scope.Compute(node.prefix)
	if err != nil {
		return NewNodeError(node, err)
	}

	textStr, ok1 := text.(string)
	prefixStr, ok2 := prefix.(string)

	if !ok1 || !ok2 {
		return NewNodeError(node, fmt.Errorf("text and prefix should be string, got %T and %T", text, prefix))
	}

	return Literal(strings.HasPrefix(textStr, prefixStr))
}

func (node HasPrefixNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "text", "prefix")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'not' node: %w", err)
		}

		return HasPrefix(orderedArgs["text"], orderedArgs["prefix"]), nil
	})
}
