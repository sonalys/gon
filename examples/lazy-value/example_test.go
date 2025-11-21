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
		WithValues(gon.Values{
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

	value, err := scope.Compute(lazyRule)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	if !called {
		panic("should have been called")
	}
	//Output:
	// true
}
