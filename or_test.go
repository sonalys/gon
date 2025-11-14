package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Or(t *testing.T) {
	scope := gon.NewScope()

	t.Run("should have at least one child", func(t *testing.T) {
		expr := gon.Or()

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should not have unset children", func(t *testing.T) {
		expr := gon.Or(nil)

		_, err := scope.Compute(expr)
		require.Error(t, err)
	})

	t.Run("should return first true expression", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(false),
			gon.Literal(true),
		)

		resp := expr.Eval(scope)
		require.Equal(t, true, resp.Value())
	})

	t.Run("should return first non false expression", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(false),
			gon.Literal(1),
		)

		resp := expr.Eval(scope)
		require.Equal(t, 1, resp.Value())
	})

	t.Run("should return false if none is matched", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(false),
			gon.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, false, resp.Value())
	})

	t.Run("should stop at first matched non-boolean condition", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(false),
			gon.Literal(1),
			gon.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, 1, resp.Value())
	})

	t.Run("should stop at first matched boolean condition", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(false),
			gon.Literal(true),
			gon.Literal(false),
		)

		resp := expr.Eval(scope)
		require.Equal(t, true, resp.Value())
	})

	t.Run("should propagate error", func(t *testing.T) {
		expr := gon.Or(
			gon.Literal(assert.AnError),
			gon.Literal(1),
			gon.Literal(false),
		)

		resp := expr.Eval(scope)
		err, ok := resp.Value().(error)
		require.True(t, ok)
		require.ErrorIs(t, err, assert.AnError)
	})
}
