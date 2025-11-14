package nodes_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/encoding"
	"github.com/sonalys/gon/internal/nodes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HasPrefix(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should error if text is nil", func(t *testing.T) {
		expr := nodes.HasPrefix(nil, nodes.Literal(""))

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should error if prefix is nil", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal(""), nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should propagate text error", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal(assert.AnError), nodes.Literal(""))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should propagate prefix error", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal(""), nodes.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should error if text is not string", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal(1), nodes.Literal(""))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should error if prefix is not string", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal(""), nodes.Literal(1))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should return true for prefix match", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal("important"), nodes.Literal("im"))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.True(t, got.(bool))
	})

	t.Run("should return false for no prefix match", func(t *testing.T) {
		expr := nodes.HasPrefix(nodes.Literal("important"), nodes.Literal("tant"))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.False(t, got.(bool))
	})
}

func Test_HasPrefix_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.HasPrefix(nodes.Literal(true), nodes.Literal(1))

			shaped, ok := node.(adapters.Shaped)
			require.True(t, ok)

			kns := shaped.Shape()

			registerer, ok := node.(encoding.AutoRegisterer)
			require.True(t, ok)

			codex := make(encoding.Codex)

			err := registerer.Register(&codex)
			require.NoError(t, err)

			named, ok := node.(adapters.Named)
			require.True(t, ok)
			assert.NotEmpty(t, named.Scalar())

			got, err := codex[named.Scalar()](kns)
			require.NoError(t, err)
			require.Equal(t, node, got)
		})
	})
}
