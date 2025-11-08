package main_test

import (
	"fmt"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_gonAgeVerification() {
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

	ageRuleStr := `if(gte(person.age, 18), "pass", "fail")`

	rule, err := encoding.Decode([]byte(ageRuleStr), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	value := rule.Eval(scope).Value()
	fmt.Println(value)

	person.Age = 5

	value = rule.Eval(scope).Value()
	fmt.Println(value)
	// Output:
	// pass
	// fail
}
