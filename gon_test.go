package gon_test

import (
	"fmt"
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Expression(t *testing.T) {
	scope := gon.NewScope().
		WithContext(t.Context()).
		WithDefinitions(map[string]gon.Expression{
			"var1": gon.Static("name"),
			"var2": gon.Static("name2"),
			"callable": gon.Static(gon.Function(func(name string, age int) string {
				fmt.Printf("Hello %s, you are %d years old!\n", name, age)

				return "surprise!"
			})),
		})

	isEqual := gon.Equal(
		gon.Equal(
			gon.Definition("var1"),
			gon.Definition("var2"),
		),
		gon.Equal(
			gon.Definition("var1"),
			gon.Definition("var2"),
		),
	)

	got, ok := isEqual.Eval(scope).Bool()
	require.True(t, ok)
	require.True(t, got)

	resp := gon.Call("callable", gon.Definition("var1"), gon.Static(5)).Eval(scope)
	t.Errorf("got: %v", resp.Any())
	t.Fail()
}

func Benchmark_Equal(b *testing.B) {
	scope := gon.NewScope().
		WithContext(b.Context()).
		WithDefinitions(gon.Definitions{
			"var1": gon.Static(1),
			"var2": gon.Static(1),
		})

	isEqual := gon.Equal(
		gon.Definition("var1"),
		gon.Definition("var2"),
	)

	for b.Loop() {
		isEqual.Eval(scope)
	}
}
