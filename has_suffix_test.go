package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HasSuffix(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should error if text is nil", func(t *testing.T) {
		expr := gon.HasSuffix(nil, gon.Literal(""))

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should error if prefix is nil", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal(""), nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should propagate text error", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal(assert.AnError), gon.Literal(""))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should propagate prefix error", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal(""), gon.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should error if text is not string", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal(1), gon.Literal(""))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should error if prefix is not string", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal(""), gon.Literal(1))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should return true for prefix match", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal("important"), gon.Literal("tant"))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.True(t, got.(bool))
	})

	t.Run("should return false for no prefix match", func(t *testing.T) {
		expr := gon.HasSuffix(gon.Literal("important"), gon.Literal("im"))

		got, err := scope.Compute(expr)
		require.NoError(t, err)
		require.False(t, got.(bool))
	})
}
