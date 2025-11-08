# GON

Gon is your dynamic and flexible rule-engine!

### Experimental

Gon is still experimental and under development, it's not yet stable for production.

## Goals

* Provide dynamic rule/action script evaluation
* Provide a flexible and extendable library that allows you to control 100% of the rules
	* Choose which nodes to allow or disallow
	* Implement your own custom nodes
	* Allow custom encoding/decoding formats

## Usage

### Use Gon for:

* **If**, **then**, **else** requirements
* Enforcing your ever-changing business requirements
* Scriptable and dynamic customer conditions/actions

### Don't use Gon for:

* Recursion
* Long lived execution. Example: infinite loops


### Basic Example

```go
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
```

### Further Examples

* [Age Verification](./examples/age-verification/example_test.go)
* [Dates](./examples/dates/example_test.go)
* [Dynamic Comparison](./examples/dynamic-comparison/example_test.go)
* [Lazy-value](./examples/lazy-value/example_test.go)

## Standard Nodes

* Literal
* Reference
* Call
* Equal
* Greater
* Smaller
* If
* Or
* Not

## Roadmap

* Better slice definition and referencing
* Extensive test coverage
* More operations:
  * Contains
  * Zero ( zero-value )
  * Etc...