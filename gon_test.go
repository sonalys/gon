package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Expression(t *testing.T) {
	scope := gon.NewScope().
		WithContext(t.Context()).
		WithVariables(map[string]gon.Expression{
			"var1": gon.Static("name"),
			"var2": gon.Static("name"),
		})

	expr := gon.Equal(
		gon.Variable("var1"),
		gon.Variable("var2"),
	)

	switch v := expr.Eval(scope).Any().(type) {
	case bool:
		require.True(t, v)
	default:
		require.Fail(t, "unexpected value", "expected bool got %T", v)
	}
}
