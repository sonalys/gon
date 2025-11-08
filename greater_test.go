package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Greater(t *testing.T) {
	t.Run("unset first expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(nil, gon.Literal(1))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("unset second expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(gon.Literal(1), nil)

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("cannot compare different types", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(gon.Literal(1), gon.Literal(1.))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should return true for greater", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(gon.Literal(2), gon.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return false for equal", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(gon.Literal(1), gon.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})

	t.Run("should return false for smaller", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.Greater(gon.Literal(1), gon.Literal(2))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})
}

func Test_GreaterOrEqual(t *testing.T) {
	t.Run("unset first expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(nil, gon.Literal(1))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("unset second expression", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(gon.Literal(1), nil)

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("cannot compare different types", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(gon.Literal(1), gon.Literal(1.))

		_, ok := node.Eval(scope).Value().(error)
		require.True(t, ok)
	})

	t.Run("should return true for greater", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(gon.Literal(2), gon.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return true for equal", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(gon.Literal(1), gon.Literal(1))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.True(t, got)
	})

	t.Run("should return false for smaller", func(t *testing.T) {
		scope := gon.NewScope()
		node := gon.GreaterOrEqual(gon.Literal(1), gon.Literal(2))

		got, ok := node.Eval(scope).Value().(bool)
		require.True(t, ok)
		require.False(t, got)
	})
}
