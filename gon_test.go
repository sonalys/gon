package gon_test

import (
	"fmt"
	"testing"

	"github.com/sonalys/gon"
)

func Test_Expression(t *testing.T) {
	type Friend struct {
		Name string
		Age  int
	}

	scope := gon.NewScope().
		WithContext(t.Context()).
		WithDefinitions(map[string]gon.Expression{
			"myName": gon.Static("friendName"),
			"friend": gon.Object(&Friend{
				Name: "friendName",
				Age:  50,
			}),
			"reply": gon.Function(func(name string, age int) string {
				fmt.Printf("Hello %s, you are %d years old!\n", name, age)

				return "surprise!"
			}),
			"theFinger": gon.Function(func() string {
				return "fuck off stranger!"
			}),
		})

	rule := gon.If(
		gon.Equal(
			gon.Definition("myName"),
			gon.Definition("friend.Name"),
		),
		gon.Call("reply", gon.Definition("friend.Name"), gon.Definition("friend.Age")),
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
