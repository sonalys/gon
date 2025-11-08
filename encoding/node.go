package encoding

import (
	"fmt"

	"github.com/sonalys/gon"
)

type NodeType uint8

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeExpression
	NodeTypeReference
	NodeTypeLiteral
)

type Node struct {
	Children []*Node
	Key      []byte
	Scalar   []byte
	Value    any
	Type     NodeType
}

func translateNode(rootNode *Node, codex Codex) (gon.Node, error) {
	switch rootNode.Type {
	case NodeTypeReference:
		return gon.Reference(string(rootNode.Scalar)), nil
	case NodeTypeLiteral:
		return gon.Literal(rootNode.Value), nil
	}

	constructor, ok := codex[string(rootNode.Scalar)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	children := rootNode.Children
	nodeChildren := make([]gon.KeyNode, 0, len(children))

	for _, child := range children {
		nodeChild, err := translateNode(child, codex)
		if err != nil {
			return nil, err
		}
		nodeChildren = append(nodeChildren, gon.KeyNode{
			Key:  string(child.Key),
			Node: nodeChild,
		})
	}

	return constructor(nodeChildren)
}
