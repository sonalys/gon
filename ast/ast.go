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

	ParseableNode interface {
		gon.Named
		gon.Typed
		gon.Shaped
	}

	KeyNode struct {
		Key  string
		Node Node
	}

	Invalid struct {
		Error error
	}

	Expression struct {
		Scalar  string
		KeyArgs []KeyNode
	}

	Literal struct {
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

func Parse(rootExpression gon.Node) (Node, error) {
	nodeExpression, ok := rootExpression.(ParseableNode)
	if !ok {
		return nil, fmt.Errorf("parsing node to ast: %T", rootExpression)
	}

	switch t := nodeExpression.Type(); t {
	case gon.NodeTypeExpression:
		name := nodeExpression.Scalar()
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
			Scalar:  name,
			KeyArgs: keyArgs,
		}, nil
	case gon.NodeTypeReference:
		name := nodeExpression.Scalar()
		return Reference{
			Name: name,
		}, nil
	case gon.NodeTypeLiteral:
		valuer, ok := rootExpression.(gon.Valued)
		if !ok {
			return Invalid{
				Error: fmt.Errorf("node type %v should implement %T", t, new(gon.Valued)),
			}, nil
		}

		return Literal{
			Value: valuer.Value(),
		}, nil
	default:
		return Invalid{
			Error: fmt.Errorf("invalid node type: %v", nodeExpression.Type()),
		}, nil
	}
}

func (i Invalid) Children() iter.Seq[Node] {
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

func (s Literal) Children() iter.Seq[Node] {
	return func(yield func(Node) bool) {}
}
