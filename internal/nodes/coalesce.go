package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

type CoalesceNode struct {
	definition string
	or         adapters.Node
}

func Coalesce(definition string, or adapters.Node) adapters.Node {
	return &CoalesceNode{
		definition: definition,
		or:         or,
	}
}

func (node *CoalesceNode) Eval(scope adapters.Scope) adapters.Value {
	value, ok := scope.Definition(node.definition)
	if ok {
		return value
	}

	return node.or.Eval(scope)
}

func (node *CoalesceNode) Scalar() string {
	return "coalesce"
}

func (node *CoalesceNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "", Node: Literal(node.definition)},
		{Key: "", Node: node.or},
	}
}

func (node *CoalesceNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(kn []adapters.KeyNode) (adapters.Node, error) {
		definitionName, ok := kn[0].Node.(adapters.Valued).Value().(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", definitionName)
		}

		orNode := kn[1].Node

		return Coalesce(definitionName, orNode), nil
	})
}

func (node *CoalesceNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

var (
	_ adapters.SerializableNode = &CoalesceNode{}
)
