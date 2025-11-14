package main_test

import (
	"fmt"
	"os"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_ageVerification() {
	type Person struct {
		Age int64 `gon:"age"`
	}

	person := &Person{Age: 19}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"person": gon.Literal(person),
		})
	if err != nil {
		panic(err)
	}

	// Write rules as code, and encode them to text:
	exampleRule := gon.If(
		gon.GreaterOrEqual(
			gon.Reference("person.age"),
			gon.Literal(18),
		),
		gon.Literal("pass"),
		gon.Literal("fail"),
	)

	err = encoding.HumanEncode(os.Stdout, exampleRule, encoding.Compact(), encoding.Unnamed())
	if err != nil {
		panic(err)
	}

	// Or write rules as text, and parse them to code.
	ageRuleStr := `if(gte(person.age, 18), "pass", "fail")`

	rule, err := encoding.Decode([]byte(ageRuleStr), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	value, err := scope.Compute(rule)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nBefore: %s", value)

	person.Age = 5

	value, err = scope.Compute(rule)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nAfter: %s", value)
	// Output:
	// if(gte(person.age,18),"pass","fail")
	//
	// Before: pass
	// After: fail
}
