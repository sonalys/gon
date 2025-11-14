package nodes

import (
	"fmt"
	"strings"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/gonutils"
)

type HasSuffixNode struct {
	text   adapters.Node
	suffix adapters.Node
}

// Equal defines a suffix node, all input nodes should evaluate to the same type, and be not nil.
// Returns a boolean value indicating whether the text has the suffix.
func HasSuffix(text, suffix adapters.Node) adapters.Node {
	if text == nil || suffix == nil {
		return adapters.NodeError{
			NodeScalar: "suffix",
			Cause:      adapters.ErrAllNodesMustBeSet,
		}
	}

	return &HasSuffixNode{
		text:   text,
		suffix: suffix,
	}
}

func (node *HasSuffixNode) Scalar() string {
	return "hasSuffix"
}

func (node *HasSuffixNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "text", Node: node.text},
		{Key: "suffix", Node: node.suffix},
	}
}

func (node *HasSuffixNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *HasSuffixNode) Eval(scope adapters.Scope) adapters.Value {
	text, err := scope.Compute(node.text)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	prefix, err := scope.Compute(node.suffix)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	textStr, ok1 := text.(string)
	prefixStr, ok2 := prefix.(string)

	if !ok1 || !ok2 {
		return adapters.NewNodeError(node, fmt.Errorf("text and suffix should be string, got %T and %T", text, prefix))
	}

	return Literal(strings.HasSuffix(textStr, prefixStr))
}

func (node *HasSuffixNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := gonutils.SortArgs(args, "text", "suffix")
		if err != nil {
			return nil, err
		}

		return HasSuffix(orderedArgs["text"], orderedArgs["suffix"]), nil
	})
}

var _ adapters.SerializableNode = &HasSuffixNode{}
