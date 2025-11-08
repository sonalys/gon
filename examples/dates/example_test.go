package main_test

import (
	"fmt"
	"time"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_dates() {
	type Person struct {
		Birthday time.Time `gon:"age"`
	}

	person := &Person{Birthday: time.Date(2000, 1, 1, 1, 1, 1, 0, time.Local)}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"person": gon.Literal(person),
		})
	if err != nil {
		panic(err)
	}

	ageRuleStr := `if(lt(person.age, time("2016-10-31T11:07:39+01:00")), "pass", "fail")`

	rule, err := encoding.Decode([]byte(ageRuleStr), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	value := rule.Eval(scope).Value()
	fmt.Println(value)
	// Output:
	// pass
}
