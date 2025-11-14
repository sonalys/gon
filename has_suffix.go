package gon

import (
	"fmt"
	"strings"
)

type HasSuffixNode struct {
	text   Node
	suffix Node
}

// Equal defines a suffix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the suffix.
func HasSuffix(text, suffix Node) Node {
	if text == nil || suffix == nil {
		return NodeError{
			NodeScalar: "suffix",
			Cause:      ErrAllNodesMustBeSet,
		}
	}

	return HasSuffixNode{
		text:   text,
		suffix: suffix,
	}
}

func (node HasSuffixNode) Scalar() string {
	return "hasSuffix"
}

func (node HasSuffixNode) Shape() []KeyNode {
	return []KeyNode{
		{"text", node.text},
		{"suffix", node.suffix},
	}
}

func (node HasSuffixNode) Type() NodeType {
	return NodeTypeExpression
}

func (node HasSuffixNode) Eval(scope Scope) Value {
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

func (node HasSuffixNode) Register(codex Codex) error {
	return codex.Register(node.Scalar(), func(args []KeyNode) (Node, error) {
		orderedArgs, _, err := argSorter(args, "text", "prefix")
		if err != nil {
			return nil, err
		}

		return HasSuffix(orderedArgs["text"], orderedArgs["prefix"]), nil
	})
}
