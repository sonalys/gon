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

func Test_Exists(t *testing.T) {
	t.Run("definition doesn't exist in scope", func(t *testing.T) {
		scope := gon.NewScope()

		rule := nodes.Exists("reference")

		exists, err := scope.Compute(rule)
		require.NoError(t, err)
		assert.False(t, exists.(bool))
	})

	t.Run("definition exist in scope", func(t *testing.T) {
		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(1),
			})
		require.NoError(t, err)

		rule := nodes.Exists("reference")

		exists, err := scope.Compute(rule)
		require.NoError(t, err)
		assert.True(t, exists.(bool))
	})

	t.Run("definition doesn't exist in children attribute", func(t *testing.T) {
		type Person struct{}

		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(Person{}),
			})
		require.NoError(t, err)

		rule := nodes.Exists("reference.age")

		exists, err := scope.Compute(rule)
		require.NoError(t, err)
		assert.False(t, exists.(bool))
	})

	t.Run("definition exists in children attribute", func(t *testing.T) {
		type Person struct {
			Age int `gon:"age"`
		}

		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(Person{}),
			})
		require.NoError(t, err)

		rule := nodes.Exists("reference.age")

		exists, err := scope.Compute(rule)
		require.NoError(t, err)
		assert.True(t, exists.(bool))
	})
}

func Test_Exists_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.Exists("reference")

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
