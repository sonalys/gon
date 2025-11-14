package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Avg(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should have at least one child", func(t *testing.T) {
		expr := gon.Avg()

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should not have unset children", func(t *testing.T) {
		expr := gon.Avg(nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should propagate error", func(t *testing.T) {
		expr := gon.Avg(gon.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("all children should be of the same type", func(t *testing.T) {
		expr := gon.Avg(gon.Literal(1), gon.Literal(1.))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should average values", func(t *testing.T) {
		expr := gon.Avg(gon.Literal(1), gon.Literal(3))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.Equal(t, 2, got)
	})
}

func Test_Avg_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := gon.Avg(gon.Literal(true), gon.Literal(1))

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
