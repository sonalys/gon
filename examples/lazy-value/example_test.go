package main_test

import (
	"context"
	"fmt"

	"github.com/sonalys/gon"
)

func Example_lazyValue() {
	var called bool

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"lazy": gon.Literal(func(ctx context.Context) int {
				called = true
				return 5
			}),
		})
	if err != nil {
		panic(err)
	}

	if called {
		panic("not lazy")
	}

	lazyRule := gon.Equal(gon.Reference("lazy"), gon.Literal(5))

	got := lazyRule.Eval(scope).Value()
	fmt.Println(got)

	if !called {
		panic("should have been called")
	}
	//Output:
	// true
}
