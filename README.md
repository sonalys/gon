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


### Example Usage

```go
func TestExample(t *testing.T) {
	type Friend struct {
		Name string
		Age  int
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
			"reply": gon.Function(func(name string, msg any) string {
				fmt.Printf("Hello %s, you are %s!\n", name, msg)

				return "surprise!"
			}),
			"whoAreYou": gon.Function(func() string {
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
			gon.Definition("friend.Name"),
		),
        // Main branch if condition fulfilled.
		gon.Call("reply",
			gon.Definition("friend.Name"),
			gon.If(
				gon.Greater(
					gon.Definition("friend.Age"),
					gon.Static(18),
				),
				gon.Static("old"),
				gon.Static("young"),
			),
		),
		gon.Call("whoAreYou"),
	)

	resp := rule.Eval(scope)
    // Prints:
    // Hello friendName, you are young!
    //
    // Returns:
    // surprise!
}
```

## Roadmap

* Encoder/Decoder
* Support GON tags