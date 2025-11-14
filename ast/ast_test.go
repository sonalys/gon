package ast_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unparsableNode struct{}

func (u unparsableNode) Register(codex adapters.Codex) error {
	panic("unimplemented")
}

func (u unparsableNode) Eval(scope adapters.Scope) adapters.Value {
	panic("unimplemented")
}

func (u unparsableNode) Scalar() string {
	panic("unimplemented")
}

func (u unparsableNode) Shape() []adapters.KeyNode {
	panic("unimplemented")
}

func (u unparsableNode) Type() adapters.NodeType {
	return adapters.NodeType(255)
}

var _ adapters.SerializableNode = unparsableNode{}
var _ adapters.Node = unparsableNode{}

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
