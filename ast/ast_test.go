package ast_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
}
