package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ ast.ParseableNode = &gon.IfNode{}
)

func Test_If(t *testing.T) {
	scope := gon.NewScope()

	t.Run("if expression should not be unset", func(t *testing.T) {
		expr := gon.If(nil, gon.Literal(true), gon.Literal(false))

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should error on multiple else branches", func(t *testing.T) {
		expr := gon.If(gon.Literal(true), gon.Literal(false), gon.Literal(false), gon.Literal(false))

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should compute main branch if expression is true", func(t *testing.T) {
		expr := gon.If(gon.Literal(true), gon.Literal(true), gon.Literal(false))

		value, err := scope.Compute(expr)
		require.NoError(t, err)
		require.True(t, value.(bool))
	})

	t.Run("should compute else branch if expression is true", func(t *testing.T) {
		expr := gon.If(gon.Literal(false), gon.Literal(false), gon.Literal(true))

		value, err := scope.Compute(expr)
		require.NoError(t, err)
		require.True(t, value.(bool))
	})

	t.Run("should error if expression doesn't return a bool", func(t *testing.T) {
		t.Run("should compute main branch if expression is true", func(t *testing.T) {
			expr := gon.If(gon.Literal(1), gon.Literal(true), gon.Literal(false))

			_, err := scope.Compute(expr)
			require.Error(t, err)
		})
	})

	t.Run("should return false if no else branch", func(t *testing.T) {
		expr := gon.If(gon.Literal(false), gon.Literal(true))

		value, err := scope.Compute(expr)
		require.NoError(t, err)
		require.False(t, value.(bool))
	})

	t.Run("should propagate expression error", func(t *testing.T) {
		expr := gon.If(gon.Literal(assert.AnError), gon.Literal(true))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should propagate then branch error", func(t *testing.T) {
		expr := gon.If(gon.Literal(true), gon.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})

	t.Run("should propagate else branch error", func(t *testing.T) {
		expr := gon.If(gon.Literal(false), gon.Literal(true), gon.Literal(assert.AnError))

		_, err := scope.Compute(expr)
		require.ErrorAs(t, err, &gon.NodeError{})
	})
}
