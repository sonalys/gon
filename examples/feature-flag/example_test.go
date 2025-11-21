package main_test

import (
	"fmt"

	"github.com/sonalys/gon"
)

func Example_featureFlag() {
	scope, err := gon.
		NewScope().
		WithValues(gon.Values{})
	if err != nil {
		panic(err)
	}

	rateLimitPolicy := gon.If(gon.Exists("optionalConnector"), gon.Reference("optionalConnector.premiumLimit"), gon.Literal(10))

	rateLimit, err := scope.Compute(rateLimitPolicy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nBefore: %d", rateLimit)

	type premiumMembership struct {
		RateLimit int `gon:"premiumLimit"`
	}

	scope, err = scope.WithValues(gon.Values{
		"optionalConnector": gon.Literal(premiumMembership{RateLimit: 500}),
	})
	if err != nil {
		panic(err)
	}

	rateLimit, err = scope.Compute(rateLimitPolicy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nAfter: %d", rateLimit)

	//Output:
	// Before: 10
	// After: 500
}
