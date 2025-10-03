package gon_test

import (
	"fmt"
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Expression(t *testing.T) {
	type Friend struct {
		Name string `gon:"name"`
		Age  int    `gon:"age"`
	}

	scope, err := gon.NewScope().
		// Context cancellation
		WithContext(t.Context()).
		// Dynamic, decoupled scope for your rules.
		WithDefinitions(map[string]gon.Expression{
			// Support for static variables of any type.
			"myName": gon.Static("friendName"),
			// Support for structs and maps.
			"friend": gon.Object(&Friend{
				Name: "friendName",
				Age:  5,
			}),
			// Support for callable function definitions.
			"reply": gon.Function(func(name string, msg any) string {
				fmt.Printf("Hello %s, you are %s!\n", name, msg)

				return "surprise!"
			}),
			"whoAreYou": gon.Function(func() string {
				return "I don't know you!"
			}),
		})
		// Error on invalid key names.
		// Should start with a-z and only contain alphanumeric characters.
	require.NoError(t, err)

	// If-else branch.
	rule := gon.If(
		gon.Equal(
			// Scope variable referencing.
			gon.Definition("myName"),
			gon.Definition("friend.name"),
		),
		// Main branch if condition fulfilled.
		gon.Call("reply",
			gon.Definition("friend.name"),
			gon.If(
				gon.Greater(
					gon.Definition("friend.age"),
					gon.Static(18),
				),
				gon.Static("old"),
				gon.Static("young"),
			),
		),
		gon.Call("whoAreYou"),
	)
	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Any())

	err = gon.Encode(t.Output(), rule)
	require.NoError(t, err)

	t.Fail()
}

func Benchmark_Equal(b *testing.B) {
	scope, _ := gon.NewScope().
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
