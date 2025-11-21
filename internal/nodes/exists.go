package nodes

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

type ExistsNode struct {
	definition string
}

func Exists(definition string) adapters.Node {
	return &ExistsNode{
		definition: definition,
	}
}

func (node *ExistsNode) Eval(scope adapters.Scope) adapters.Value {
	_, ok := scope.Definition(node.definition)
	return Literal(ok)
}

func (node *ExistsNode) Scalar() string {
	return "exists"
}

func (node *ExistsNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "", Node: Literal(node.definition)},
	}
}

func (node *ExistsNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(kn []adapters.KeyNode) (adapters.Node, error) {
		definitionName, ok := kn[0].Node.(adapters.Valued).Value().(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", definitionName)
		}

		return Exists(definitionName), nil
	})
}

func (node *ExistsNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

var (
	_ adapters.SerializableNode = &ExistsNode{}
)
