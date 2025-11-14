package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Not(t *testing.T) {
	t.Run("unset expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Not(nil)

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should propagate error for non boolean value", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Not(gon.Literal(1))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should return true for false", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Not(gon.Literal(false))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return false for true", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Not(gon.Literal(true))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})
}

func Test_Not_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := gon.Not(gon.Literal(1))

			shaped, ok := node.(gon.Shaped)
			require.True(t, ok)

			kns := shaped.Shape()

			registerer, ok := node.(encoding.AutoRegisterer)
			require.True(t, ok)

			codex := make(encoding.Codex)

			err := registerer.Register(&codex)
			require.NoError(t, err)

			named, ok := node.(gon.Named)
			require.True(t, ok)
			assert.NotEmpty(t, named.Scalar())

			got, err := codex[named.Scalar()](kns)
			require.NoError(t, err)
			require.Equal(t, node, got)
		})
	})
}
