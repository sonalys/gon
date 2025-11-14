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

func Test_Sum(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should have at least one child", func(t *testing.T) {
		expr := nodes.Sum()

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should not have unset children", func(t *testing.T) {
		expr := nodes.Sum(nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should propagate error", func(t *testing.T) {
		expr := nodes.Sum(nodes.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("all children should be of the same type", func(t *testing.T) {
		expr := nodes.Sum(nodes.Literal(1), nodes.Literal(1.))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should average values", func(t *testing.T) {
		expr := nodes.Sum(nodes.Literal(1), nodes.Literal(3))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.Equal(t, 4, got)
	})
}

func Test_Sum_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.Sum(nodes.Literal(true), nodes.Literal(1))

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
