package nodes

import (
	"fmt"
	"strings"

	"github.com/sonalys/gon/adapters"
)

type HasPrefixNode struct {
	text   adapters.Node
	prefix adapters.Node
}

// Equal defines a prefix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the prefix.
func HasPrefix(text, prefix adapters.Node) adapters.Node {
	if text == nil || prefix == nil {
		return adapters.NodeError{
			NodeScalar: "prefix",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &HasPrefixNode{
		text:   text,
		prefix: prefix,
	}
}

func (node *HasPrefixNode) Scalar() string {
	return "hasPrefix"
}

func (node *HasPrefixNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{"text", node.text},
		{"prefix", node.prefix},
	}
}

func (node *HasPrefixNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *HasPrefixNode) Eval(scope adapters.Scope) adapters.Value {
	text, err := scope.Compute(node.text)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	prefix, err := scope.Compute(node.prefix)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	textStr, ok1 := text.(string)
	prefixStr, ok2 := prefix.(string)

	if !ok1 || !ok2 {
		return adapters.NewNodeError(node, fmt.Errorf("text and prefix should be string, got %T and %T", text, prefix))
	}

	return Literal(strings.HasPrefix(textStr, prefixStr))
}

func (node *HasPrefixNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := argSorter(args, "text", "prefix")
		if err != nil {
			return nil, err
		}

		return HasPrefix(orderedArgs["text"], orderedArgs["prefix"]), nil
	})
}
