# GON

[![Go Reference](https://pkg.go.dev/badge/github.com/sonalys/gon.svg)](https://pkg.go.dev/github.com/sonalys/gon)
[![Tests](https://github.com/sonalys/gon/actions/workflows/test.yml/badge.svg)](https://github.com/sonalys/gon/actions/workflows/test.yml)
[![Linter](https://github.com/sonalys/gon/actions/workflows/lint.yml/badge.svg)](https://github.com/sonalys/gon/actions/workflows/lint.yml)
[![codecov](https://codecov.io/github/sonalys/gon/graph/badge.svg?token=N0XL7NLXIL)](https://codecov.io/github/sonalys/gon)

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
```

### Further Examples

* [Age Verification](./examples/age-verification/example_test.go)
* [Dates](./examples/dates/example_test.go)
* [Dynamic Comparison](./examples/dynamic-comparison/example_test.go)
* [Lazy value](./examples/lazy-value/example_test.go)
* [Custom Node](./examples/custom-node/example_test.go)
* [Functions](./examples/functions/example_test.go)
* [Object access rule](./examples/object-access-rule/example_test.go)

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
* HasPrefix
* HasSuffix
* Sum
* Avg

## Limitations

* Uses reflect package
* Some node types like Literal, Bool and Time are still not fully customizable

## Roadmap

I want to extend the project in the direction of having further:

* Better slice definition and referencing
* Benchmarks
* More operations:
  * Contains
  * Zero ( zero-value )
  * Etc...

### Contributing

If you want to contribute, feel free to open discussions and issues.