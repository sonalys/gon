package gon_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Expression(t *testing.T) {
	type Friend struct {
		Name     string    `gon:"name"`
		Birthday time.Time `gon:"birthday"`
	}

	birthday := time.Now().AddDate(-15, 0, 0)

	scope, err := gon.NewScope().
		// Context cancellation
		WithContext(t.Context()).
		// Dynamic, decoupled scope for your rules.
		WithDefinitions(map[string]gon.Expression{
			// Support for static variables of any type.
			"myName": gon.Static("friendName"),
			// Support for structs and maps.
			"friend": gon.Object(&Friend{
				Name:     "friendName",
				Birthday: birthday,
			}),
			// Support for callable function definitions.
			"reply": gon.Static(func(name string, msg any) string {
				switch msg := msg.(type) {
				case error:
					return fmt.Sprintf("unexpected error: %s", msg.Error())
				case string:
					fmt.Printf("Hello %s, you are %s!\n", name, msg)
				}

				return "surprise!"
			}),
			"whoAreYou": gon.Static(func() string {
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
				gon.Smaller(
					gon.Definition("friend.birthday"),
					gon.Static(time.Now().AddDate(-18, 0, 0)),
				),
				gon.Static("old"),
				gon.Static("young"),
			),
		),
		gon.Call("whoAreYou"),
	)
	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Value())

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
