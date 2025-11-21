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

func Test_Coalesce(t *testing.T) {
	t.Run("definition doesn't exist in scope", func(t *testing.T) {
		scope := gon.NewScope()

		rule := nodes.Coalesce("reference", nodes.Literal(1))

		value, err := scope.Compute(rule)
		require.NoError(t, err)
		require.Equal(t, 1, value)
	})

	t.Run("definition exist in scope", func(t *testing.T) {
		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(1),
			})
		require.NoError(t, err)

		rule := nodes.Coalesce("reference", nodes.Literal(2))

		value, err := scope.Compute(rule)
		require.NoError(t, err)
		require.Equal(t, 1, value)
	})

	t.Run("definition doesn't exist in children attribute", func(t *testing.T) {
		type Person struct{}

		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(Person{}),
			})
		require.NoError(t, err)

		rule := nodes.Coalesce("reference.age", nodes.Literal(1))

		value, err := scope.Compute(rule)
		require.NoError(t, err)
		require.Equal(t, 1, value)
	})

	t.Run("definition Coalesce in children attribute", func(t *testing.T) {
		type Person struct {
			Age int `gon:"age"`
		}

		scope, err := gon.NewScope().
			WithValues(gon.Values{
				"reference": nodes.Literal(Person{}),
			})
		require.NoError(t, err)

		rule := nodes.Coalesce("reference.age", nodes.Literal(1))

		value, err := scope.Compute(rule)
		require.NoError(t, err)
		require.Equal(t, 0, value)
	})
}

func Test_Coalesce_Encoding(t *testing.T) {
	t.Run("should decode with children", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := nodes.Coalesce("reference", nodes.Literal(5))

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
