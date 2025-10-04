# GON

Gon is your dynamic and flexible rule-engine!


## Usage

### Use Gon for:

* **If**, **then**, **else** requirements
* Enforcing your ever-changing business requirements
* Scriptable customer conditions/actions

### Don't use Gon for:

* Iterations
* Recursions
* Statically validated expressions
  * It won't give you validation ahead of execution for invalid operations


### Example Usage

```go
func Test_Expression(t *testing.T) {
	type Friend struct {
		Name string `gon:"name"`
		Age  int    `gon:"age"`
	}

	scope, err := gon.NewScope().
		// Context cancellation
		WithContext(t.Context()).
		// Dynamic, decoupled scope for your rules.
		WithDefinitions(map[string]gon.Expression{
			// Support for static variables of any type.
			"myName": gon.Static("friendName"),
			// Support for structs and maps.
			"friend": gon.Object(&Friend{
				Name: "friendName",
				Age:  5,
			}),
			// Support for callable function definitions.
			"reply": gon.Static(func(name string, msg any) string {
				fmt.Printf("Hello %s, you are %s!\n", name, msg)

				return "surprise!"
			}),
			"whoAreYou": gon.Static(func() string {
				return "I don't know you!"
			}),
		})
	// Error on invalid key names.
	// Should start with a-z and only contain alphanumeric characters.
	require.NoError(t, err)

	// If-else branch.
	rule := gon.If(
		gon.Equal(
			// Scope variable referencing.
			gon.Definition("myName"),
			gon.Definition("friend.name"),
		),
		// Main branch if condition fulfilled.
		gon.Call("reply",
			gon.Definition("friend.name"),
			gon.If(
				gon.Greater(
					gon.Definition("friend.age"),
					gon.Static(18),
				),
				gon.Static("old"),
				gon.Static("young"),
			),
		),
		gon.Call("whoAreYou"),
	)
	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Value())

	err = gon.Encode(t.Output(), rule)
	require.NoError(t, err)

	t.Fail()
}
// Prints:
// Hello friendName, you are young!
//
// Returns:
// surprise!
//
// if(
// 	condition: equal(
// 		first: myName
// 		second: friend.name
// 	)
// 	then: call("reply"
// 		friend.name
// 		if(
// 			condition: gt(
// 				first: friend.age
// 				second: 18
// 			)
// 			then: "old"
// 			else: "young"
// 		)
// 	)
// 	else: call("whoAreYou")
// )
```

## Roadmap

* Decoder
* Better slice definition and referencing
* More operations:
  * Between
  * Contains
  * Not
  * Zero ( zero-value )
  * Etc...