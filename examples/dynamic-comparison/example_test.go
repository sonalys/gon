package main_test

import (
	"fmt"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_dynamicComparison() {
	type Person struct {
		Age int64 `gon:"age"`
	}

	type AgeConfig struct {
		MinAge int64 `gon:"min_age"`
	}

	person := &Person{Age: 19}
	config := &AgeConfig{MinAge: 18}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"person":     gon.Literal(person),
			"age_config": gon.Literal(config),
		})
	if err != nil {
		panic(err)
	}

	ageRuleStr := `if(gte(person.age, age_config.min_age), "pass", "fail")`

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
