package gon_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding/goncoder"
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
			"myName": gon.Literal("friendName"),
			// Support for structs and maps.
			"friend": gon.Object(&Friend{
				Name:     "friendName",
				Birthday: birthday,
			}),
			// Support for callable function definitions.
			"reply": gon.Function(func(ctx context.Context, name string, msg any) string {
				switch msg := msg.(type) {
				case error:
					return fmt.Sprintf("unexpected error: %s", msg.Error())
				case string:
					fmt.Printf("Hello %s, you are %s!\n", name, msg)
				}

				return "surprise!"
			}),
			"whoAreYou": gon.Function(func(ctx context.Context) string {
				return "I don't know you!"
			}),
		})
	// Error on invalid key names.
	// Should start with a-z and only contain alphanumeric characters.
	require.NoError(t, err)

	ruleStr := `if(
	condition: equal(myName, friend.name),
	then: call("reply"
		friend.name
		if(lt(friend.birthday, time("2007-10-31T11:07:39+01:00")), "old", "young")
	),
	else: call("whoAreYou")
)`

	rule, err := goncoder.Decode([]byte(ruleStr), goncoder.DefaultExpressionCodex)
	require.NoError(t, err)

	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Value())

	err = goncoder.Encode(t.Output(), rule)
	require.NoError(t, err)

	t.Fail()
}

func Benchmark_Equal(b *testing.B) {
	scope, _ := gon.NewScope().
		WithContext(b.Context()).
		WithDefinitions(gon.Definitions{
			"var1": gon.Literal(1),
			"var2": gon.Literal(1),
		})

	isEqual := gon.Equal(
		gon.Reference("var1"),
		gon.Reference("var2"),
	)

	for b.Loop() {
		isEqual.Eval(scope)
	}
}
