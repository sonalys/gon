package ast

import (
	"fmt"
	"iter"

	"github.com/sonalys/gon"
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

type NodeExpression interface {
	gon.Named
	gon.Typed
	gon.Shaped
}

func Parse(rootExpression gon.Node) (Node, error) {
	nodeExpression, ok := rootExpression.(NodeExpression)
	if !ok {
		return nil, fmt.Errorf("could not parse node to ast: %T", rootExpression)
	}

	switch t := nodeExpression.Type(); t {
	case gon.NodeTypeExpression:
		name := nodeExpression.Name()
		keyExpressions := nodeExpression.Shape()

		keyArgs := make([]KeyNode, 0, len(keyExpressions))

		for i := range keyExpressions {
			parsed, err := Parse(keyExpressions[i].Node)
			if err != nil {
				return nil, fmt.Errorf("parsing keyed expression: %w", err)
			}

			keyArgs = append(keyArgs, KeyNode{
				Key:  keyExpressions[i].Key,
				Node: parsed,
			})
		}

		return Expression{
			Name:    name,
			KeyArgs: keyArgs,
		}, nil
	case gon.NodeTypeReference:
		name := nodeExpression.Name()
		return Reference{
			Name: name,
		}, nil
	case gon.NodeTypeValue:
		valuer, ok := rootExpression.(gon.Valued)
		if !ok {
			return InvalidNode{
				Error: fmt.Errorf("node type %v should implement %T", t, new(gon.Valued)),
			}, nil
		}

		return StaticValue{
			Value: valuer.Value(),
		}, nil
	default:
		return InvalidNode{
			Error: fmt.Errorf("invalid node type: %v", nodeExpression.Type()),
		}, nil
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
