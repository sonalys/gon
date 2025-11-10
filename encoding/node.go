package encoding

import (
	"fmt"

	"github.com/sonalys/gon"
)

type Node struct {
	Children []*Node
	Key      []byte
	Scalar   []byte
	Value    any
	Type     gon.NodeType
}

func translateNode(rootNode *Node, codex Codex) (gon.Node, error) {
	switch rootNode.Type {
	case gon.NodeTypeReference:
		return gon.Reference(string(rootNode.Scalar)), nil
	case gon.NodeTypeLiteral:
		return gon.Literal(rootNode.Value), nil
	}

	constructor, ok := codex[string(rootNode.Scalar)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	children := rootNode.Children
	nodeChildren := make([]gon.KeyNode, 0, len(children))

	for _, child := range children {
		childNode, err := translateNode(child, codex)
		if err != nil {
			return nil, err
		}
		nodeChildren = append(nodeChildren, gon.KeyNode{
			Key:  string(child.Key),
			Node: childNode,
		})
	}

	return constructor(nodeChildren)
}
