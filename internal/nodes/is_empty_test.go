package nodes_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IsEmpty(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should error on unset node", func(t *testing.T) {
		node := gon.IsEmpty(nil)

		_, err := scope.Compute(node)
		require.ErrorAs(t, err, &adapters.NodeError{})
	})

	t.Run("should work for chan", func(t *testing.T) {
		ch := make(chan struct{}, 1)
		ch <- struct{}{}

		node := gon.IsEmpty(gon.Literal(ch))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		<-ch

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})

	t.Run("should work for map", func(t *testing.T) {
		value := map[string]struct{}{
			"a": {},
		}

		node := gon.IsEmpty(gon.Literal(value))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		delete(value, "a")

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})

	t.Run("should work for slice", func(t *testing.T) {
		value := []int{1}

		node := gon.IsEmpty(gon.Literal(value))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		value = value[:0]
		node = gon.IsEmpty(gon.Literal(value))

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})

	t.Run("should work for array", func(t *testing.T) {
		value := [1]int{1}

		node := gon.IsEmpty(gon.Literal(value))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		value2 := [0]int{}
		node = gon.IsEmpty(gon.Literal(value2))

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})

	t.Run("should work for string", func(t *testing.T) {
		value := "a"

		node := gon.IsEmpty(gon.Literal(value))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		value = ""
		node = gon.IsEmpty(gon.Literal(value))

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})

	t.Run("should work for pointer", func(t *testing.T) {
		value := &[]int{1}

		node := gon.IsEmpty(gon.Literal(value))

		got, err := scope.Compute(node)
		require.NoError(t, err)
		assert.False(t, got.(bool))

		*value = (*value)[:0]

		got, err = scope.Compute(node)
		require.NoError(t, err)
		assert.True(t, got.(bool))
	})
}

func Test_IsEmpty_Encoding(t *testing.T) {
	t.Run("should decode without else", func(t *testing.T) {
		require.NotPanics(t, func() {
			node := gon.IsEmpty(gon.Literal(true))

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
