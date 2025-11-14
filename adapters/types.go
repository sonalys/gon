package adapters

type NodeType uint8

const (
	NodeTypeInvalid NodeType = iota
	// NodeTypeExpression represents an expression() node type. Example: if()
	NodeTypeExpression
	// NodeTypeReference represents a variable reference. Example: friend.name.
	NodeTypeReference
	// NodeTypeLiteral represents a direct value. Example: "string", 5.
	NodeTypeLiteral
	_nodeTypeCeiling
)

func (t NodeType) IsValid() bool {
	return t > NodeTypeInvalid && t < _nodeTypeCeiling
}
