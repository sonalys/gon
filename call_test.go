package gon_test

import (
	"context"
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Call(t *testing.T) {
	t.Run("definition not found", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Call("not_found")

		got := node.Eval(scope)

		err, ok := got.Value().(error)
		require.True(t, ok)

		var target gon.DefinitionNotFoundError
		require.ErrorAs(t, err, &target)
	})

	t.Run("definition not callable", func(t *testing.T) {
		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(1),
			})
		require.NoError(t, err)

		node := gon.Call("var")

		got := node.Eval(scope)

		err, ok := got.Value().(error)
		require.True(t, ok)

		var target gon.DefinitionNotCallableError
		require.ErrorAs(t, err, &target)
	})

	t.Run("should call func without context", func(t *testing.T) {
		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(func() int {
					return 5
				}),
			})
		require.NoError(t, err)

		node := gon.Call("var")

		got := node.Eval(scope).Value()
		require.Equal(t, 5, got)
	})

	t.Run("should call func with context", func(t *testing.T) {
		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(func(ctx context.Context) int {
					require.NotNil(t, ctx)
					return 5
				}),
			})
		require.NoError(t, err)

		node := gon.Call("var")

		got := node.Eval(scope).Value()
		require.Equal(t, 5, got)
	})
}

func Test_Call_Encoding(t *testing.T) {
	t.Run("should decode without children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := gon.Call("funcName")

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

	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := gon.Call("funcName", gon.Literal(1))

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
