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

func Test_Smaller(t *testing.T) {
	t.Run("unset first expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nil, nodes.Literal(1))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("unset second expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nodes.Literal(1), nil)

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("cannot compare different types", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nodes.Literal(1), nodes.Literal(1.))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should return true for smaller", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nodes.Literal(1), nodes.Literal(2))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return false for equal", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nodes.Literal(1), nodes.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})

	t.Run("should return false for greater", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.Smaller(nodes.Literal(2), nodes.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})
}

func Test_SmallerOrEqual(t *testing.T) {
	t.Run("unset first expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nil, nodes.Literal(1))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("unset second expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nodes.Literal(1), nil)

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("cannot compare different types", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nodes.Literal(1), nodes.Literal(1.))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should return true for smaller", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nodes.Literal(1), nodes.Literal(2))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return true for equal", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nodes.Literal(1), nodes.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return false for greater", func(t *testing.T) {
		scope := gon.NewScope()
		node := nodes.SmallerOrEqual(nodes.Literal(2), nodes.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})
}

func Test_Smaller_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.Smaller(nodes.Literal(true), nodes.Literal(1))

			shaped, ok := node.(adapters.Shaped)
			require.True(t, ok)

			kns := shaped.Shape()

			registerer, ok := node.(adapters.AutoRegisterer)
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
