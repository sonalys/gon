package ast_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unparsableNode struct{}

func (u unparsableNode) Eval(scope gon.Scope) gon.Value {
	panic("unimplemented")
}

func (u unparsableNode) Scalar() string {
	panic("unimplemented")
}

func (u unparsableNode) Shape() []gon.KeyNode {
	panic("unimplemented")
}

func (u unparsableNode) Type() gon.NodeType {
	return gon.NodeType(255)
}

var _ ast.ParseableNode = unparsableNode{}
var _ gon.Node = unparsableNode{}

func Test_Parse(t *testing.T) {
	t.Run("should parse entire tree", func(t *testing.T) {
		rootNode := gon.If(gon.Literal(true), gon.Reference("key"))

		astNode, err := ast.Parse(rootNode)
		require.NoError(t, err)

		expected := ast.Expression{
			Scalar: "if",
			KeyArgs: []ast.KeyNode{
				{
					Key: "condition",
					Node: ast.Literal{
						Value: true,
					},
				},
				{
					Key: "then",
					Node: ast.Reference{
						Name: "key",
					},
				},
			},
		}

		assert.Equal(t, expected, astNode)
	})

	t.Run("should return invalid node for unknown type", func(t *testing.T) {
		rootNode := &unparsableNode{}

		astNode, err := ast.Parse(rootNode)
		require.NoError(t, err)

		_, ok := astNode.(ast.Invalid)
		require.True(t, ok)
	})
}
