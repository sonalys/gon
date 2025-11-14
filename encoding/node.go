package encoding

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/nodes"
)

type Node struct {
	Children []*Node
	Key      []byte
	Scalar   []byte
	Value    any
	Type     adapters.NodeType
}

func translateNode(rootNode *Node, codex Codex) (adapters.Node, error) {
	switch rootNode.Type {
	case adapters.NodeTypeReference:
		return nodes.Reference(string(rootNode.Scalar)), nil
	case adapters.NodeTypeLiteral:
		return nodes.Literal(rootNode.Value), nil
	}

	constructor, ok := codex[string(rootNode.Scalar)]
	if !ok {
		return nil, fmt.Errorf("codex for '%s' not found", rootNode.Scalar)
	}

	children := rootNode.Children
	nodeChildren := make([]adapters.KeyNode, 0, len(children))

	for _, child := range children {
		childNode, err := translateNode(child, codex)
		if err != nil {
			return nil, err
		}
		nodeChildren = append(nodeChildren, adapters.KeyNode{
			Key:  string(child.Key),
			Node: childNode,
		})
	}

	return constructor(nodeChildren)
}
