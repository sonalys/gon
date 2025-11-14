package gon_test

import (
	"context"
	"testing"

	"github.com/sonalys/gon"
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
