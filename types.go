package gon

type NodeType uint8

const (
	NodeTypeInvalid NodeType = iota
	// NodeTypeExpression represents an expression() node type. Example: if()
	NodeTypeExpression
	// NodeTypeReference represents a variable reference. Example: friend.name.
	NodeTypeReference
	// NodeTypeValue represents a direct value. Example: "string", 5.
	NodeTypeValue
)
