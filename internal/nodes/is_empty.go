package nodes

import (
	"fmt"
	"reflect"

	"github.com/sonalys/gon/adapters"
)

type IsEmptyNode struct {
	node adapters.Node
}

func IsEmpty(node adapters.Node) adapters.Node {
	if node == nil {
		return adapters.NodeError{
			NodeScalar: "isEmpty",
			Cause:      fmt.Errorf("node cannot be unset"),
		}
	}

	return &IsEmptyNode{
		node: node,
	}
}

func (node *IsEmptyNode) Eval(scope adapters.Scope) adapters.Value {
	value, err := scope.Compute(node.node)
	if err != nil {
		return adapters.NewNodeError(node, err)
	}

	valueOf := reflect.ValueOf(value)

	// Dereferenciation.
	for ; valueOf.Kind() == reflect.Pointer; valueOf = valueOf.Elem() {
	}

	switch valueOf.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		return Literal(valueOf.Len() == 0)
	default:
		return adapters.NewNodeError(node, fmt.Errorf("cannot calculate emptiness for %T", value))
	}
}

func (node *IsEmptyNode) Scalar() string {
	return "isEmpty"
}

func (node *IsEmptyNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "", Node: node.node},
	}
}

func (node *IsEmptyNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(kn []adapters.KeyNode) (adapters.Node, error) {
		if len(kn) != 1 {
			return nil, fmt.Errorf("expected 1 argument, got %d", len(kn))
		}

		return IsEmpty(kn[0].Node), nil
	})
}

var (
	_ adapters.SerializableNode = &IsEmptyNode{}
)
