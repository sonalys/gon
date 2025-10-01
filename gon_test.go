package gon_test

import (
	"fmt"
	"testing"

	"github.com/sonalys/gon"
)

func Test_Expression(t *testing.T) {
	scope := gon.NewScope().
		WithContext(t.Context()).
		WithDefinitions(map[string]gon.Expression{
			"myName":     gon.Static("name"),
			"friendName": gon.Static("name2"),
			"reply": gon.Static(gon.Function(func(name string, age int) string {
				fmt.Printf("Hello %s, you are %d years old!\n", name, age)

				return "surprise!"
			})),
			"theFinger": gon.Static(gon.Function(func() string {
				return "fuck off stranger!"
			})),
		})

	rule := gon.If(
		gon.Equal(
			gon.Definition("myName"),
			gon.Definition("friendName"),
		),
		gon.Call("reply", gon.Definition("myName"), gon.Static(5)),
		gon.Call("theFinger"),
	)

	resp := rule.Eval(scope)
	t.Errorf("got resp: %v", resp.Any())
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
