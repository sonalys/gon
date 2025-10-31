package ast

import (
	"fmt"
	"iter"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/internal/sliceutils"
)

type (
	Node interface {
		Children() iter.Seq[Node]
	}

	KeyNode struct {
		Key  string
		Node Node
	}

	InvalidNode struct {
		Error error
	}

	Expression struct {
		Name    string
		KeyArgs []KeyNode
	}

	StaticValue struct {
		Value any
	}

	Reference struct {
		Name string
	}
)

func Walk(rootNode Node, walkFunc func(Node) bool) {
	recursiveWalk(rootNode, walkFunc)
}

func recursiveWalk(rootNode Node, walkFunc func(Node) bool) bool {
	if !walkFunc(rootNode) {
		return false
	}

	for child := range rootNode.Children() {
		if !recursiveWalk(child, walkFunc) {
			return false
		}
	}

	return true
}

func Parse(rootExpression gon.Typed) Node {
	switch rootExpression.Type() {
	case gon.NodeTypeExpression:
		name, keyExpressions := rootExpression.Banner()
		return Expression{
			Name: name,
			KeyArgs: sliceutils.Map(keyExpressions, func(from gon.KeyExpression) KeyNode {
				return KeyNode{
					Key:  from.Key,
					Node: Parse(from.Expression),
				}
			}),
		}
	case gon.NodeTypeReference:
		name, _ := rootExpression.Banner()
		return Reference{
			Name: name,
		}
	case gon.NodeTypeValue:
		valuer := rootExpression.(gon.Valuer)

		return StaticValue{
			Value: valuer.Value(),
		}
	default:
		return InvalidNode{
			Error: fmt.Errorf("invalid node type: %v", rootExpression.Type()),
		}
	}
}

func (i InvalidNode) Children() iter.Seq[Node] {
	return func(yield func(Node) bool) {}
}

func (e Expression) Children() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for _, child := range e.KeyArgs {
			if !yield(child.Node) {
				return
			}
		}
	}
}

func (r Reference) Children() iter.Seq[Node] {
	return func(yield func(Node) bool) {}
}

func (s StaticValue) Children() iter.Seq[Node] {
	return func(yield func(Node) bool) {}
}
