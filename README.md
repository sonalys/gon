# GON

Gon is your dynamic and flexible rule-engine!

### Experimental

Gon is still experimental and under development, it's not yet stable for production.

## Goals

* Provide dynamic rule/action script evaluation
* Provide a flexible and extendable library that allows you to control 100% of the rules
	* Choose which rules to allow or disallow
	* Implement your own custom expression behavior
	* Allow custom encoding/decoding formats

## Usage

### Use Gon for:

* **If**, **then**, **else** requirements
* Enforcing your ever-changing business requirements
* Scriptable and dynamic customer conditions/actions

### Don't use Gon for:

* Iterations
* Recursions
* Statically validated expressions
  * It won't give you validation ahead of execution for invalid operations


### Example Usage

```go
func Test_Expression(t *testing.T) {
	type Friend struct {
		Name     string    `gon:"name"`
		Birthday time.Time `gon:"birthday"`
	}

	birthday := time.Now().AddDate(-15, 0, 0)

	scope, err := gon.NewScope().
		// Context cancellation
		WithContext(t.Context()).
		// Dynamic, decoupled scope for your rules.
		WithDefinitions(map[string]gon.Expression{
			// Support for static variables of any type.
			"myName": gon.Static("friendName"),
			// Support for structs and maps.
			"friend": gon.Object(&Friend{
				Name:     "friendName",
				Birthday: birthday,
			}),
			// Support for callable function definitions.
			"reply": gon.Static(func(name string, msg any) string {
				switch msg := msg.(type) {
				case error:
					return fmt.Sprintf("unexpected error: %s", msg.Error())
				case string:
					fmt.Printf("Hello %s, you are %s!\n", name, msg)
				}

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
			gon.Reference("myName"),
			gon.Reference("friend.name"),
		),
		// Main branch if condition fulfilled.
		gon.Call("reply",
			gon.Reference("friend.name"),
			gon.If(
				gon.Smaller(
					gon.Reference("friend.birthday"),
					gon.Static(time.Now().AddDate(-18, 0, 0)),
				),
				gon.Static("old"),
				gon.Static("young"),
			),
		),
		gon.Call("whoAreYou"),
	)
	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Value())

	err = goncoder.Encode(t.Output(), rule)
	require.NoError(t, err)
	// Prints:
	// 	Hello friendName, you are young!
	//     if(
	//     	condition: equal(
	//     		first: myName
	//     		second: friend.name
	//     	)
	//     	then: call("reply"
	//     		friend.name
	//     		if(
	//     			condition: lt(
	//     				first: friend.birthday
	//     				second: time("2007-10-31T11:07:39+01:00")
	//     			)
	//     			then: "old"
	//     			else: "young"
	//     		)
	//     	)
	//     	else: call("whoAreYou")
	//     )

	t.Fail()
}
```

## Roadmap

* Decoder
* Better slice definition and referencing
* Extensive test coverage
* More operations:
  * Contains
  * Zero ( zero-value )
  * Etc...