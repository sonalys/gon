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

func Test_Or(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should have at least one child", func(t *testing.T) {
		expr := nodes.Or()

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should not have unset children", func(t *testing.T) {
		expr := nodes.Or(nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should return first true expression", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(false),
			nodes.Literal(true),
		)

		resp := expr.Eval(scope)
		require.Equal(t, true, resp.Value())
	})

	t.Run("should return first non false expression", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(false),
			nodes.Literal(1),
		)

		resp := expr.Eval(scope)
		require.Equal(t, 1, resp.Value())
	})

	t.Run("should return false if none is matched", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(false),
			nodes.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, false, resp.Value())
	})

	t.Run("should stop at first matched non-boolean condition", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(false),
			nodes.Literal(1),
			nodes.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, 1, resp.Value())
	})

	t.Run("should stop at first matched boolean condition", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(false),
			nodes.Literal(true),
			nodes.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, true, resp.Value())
	})

	t.Run("should propagate error", func(t *testing.T) {
		expr := nodes.Or(
			nodes.Literal(assert.AnError),
			nodes.Literal(1),
			nodes.Literal(false),
		)

		resp := expr.Eval(scope)
		err, ok := resp.Value().(error)
		require.True(t, ok)
		require.ErrorIs(t, err, assert.AnError)
	})
}

func Test_Or_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.Or(nodes.Literal(true), nodes.Literal(1))

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
