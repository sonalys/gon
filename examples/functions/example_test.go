package main_test

import (
	"context"
	"fmt"
	"time"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_functions() {
	type Person struct {
		Name     string    `gon:"name"`
		Birthday time.Time `gon:"birthday"`
	}

	person := &Person{
		Name:     "Bob",
		Birthday: time.Date(1992, 3, 24, 0, 0, 0, 0, time.UTC),
	}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"friend": gon.Literal(person),
			"greet": gon.Literal(func(name string, birthday time.Time) string {

				return fmt.Sprintf("Hello %s, your birthday is at %s", name, birthday)
			}),
			// Context is handled automatic by gon, if specified by the function literal.
			"print": gon.Literal(func(ctx context.Context, message string) {
				fmt.Println(message)
			}),
		})
	if err != nil {
		panic(err)
	}

	ruleStr := `call("print", call("greet", friend.name, friend.birthday))`

	rule, err := encoding.Decode([]byte(ruleStr), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	_, err = scope.Compute(rule)
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello Bob, your birthday is at 1992-03-24 00:00:00 +0000 UTC
}
