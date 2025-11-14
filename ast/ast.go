package ast

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

type (
	AstNode interface {
	}

	KeyNode struct {
		Key  string
		Node AstNode
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

func Parse(rootExpression adapters.Node) (AstNode, error) {
	nodeExpression, ok := rootExpression.(adapters.SerializableNode)
	if !ok {
		return nil, fmt.Errorf("parsing node to ast: %T", rootExpression)
	}

	switch t := nodeExpression.Type(); t {
	case adapters.NodeTypeExpression:
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
	case adapters.NodeTypeReference:
		name := nodeExpression.Scalar()
		return Reference{
			Name: name,
		}, nil
	case adapters.NodeTypeLiteral:
		valuer, ok := rootExpression.(adapters.Valued)
		if !ok {
			return Invalid{
				Error: fmt.Errorf("node type %v should implement %T", t, new(adapters.Valued)),
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
