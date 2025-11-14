package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
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
