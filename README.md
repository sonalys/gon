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

* Recursion
* Long lived execution. Example: infinite loops


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
			"myName": gon.Literal("friendName"),
			// Support for structs and maps for children attributes.
			"friend": gon.Literal(&Friend{
				Name:     "friendName",
				Birthday: birthday,
			}),
			// Support for callable function definitions, with or without implicit context.
			"reply": gon.Literal(func(ctx context.Context, name string, msg any) string {
				switch msg := msg.(type) {
				case error:
					return fmt.Sprintf("unexpected error: %s", msg.Error())
				case string:
					fmt.Printf("Hello %s, you are %s!\n", name, msg)
				}

				return "surprise!"
			}),
			"whoAreYou": gon.Literal(func() string {
				return "I don't know you!"
			}),
		})
	// Error on invalid key names.
	// Should start with a-z and only contain alphanumeric characters.
	require.NoError(t, err)

	ruleStr := `if(
	condition: equal(myName, friend.name),
	then: call("reply"
		friend.name
		if(lt(friend.birthday, time("2016-10-31T11:07:39+01:00")), "old", "young")
	),
	else: call("whoAreYou")
)`

	rule, err := encoding.Decode([]byte(ruleStr), encoding.DefaultExpressionCodex)
	require.NoError(t, err)

	resp := rule.Eval(scope)
	require.Equal(t, "surprise!", resp.Value())

	err = encoding.Encode(t.Output(), rule)
	require.NoError(t, err)

	t.Fail()
	// Prints:
	// 	Hello friendName, you are old!
	// --- FAIL: Test_Expression (0.00s)
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
	//     				second: time("2016-10-31T11:07:39+01:00")
	//     			)
	//     			then: "old"
	//     			else: "young"
	//     		)
	//     	)
	//     	else: call("whoAreYou")
	//     )
}
```

## Roadmap

* Better slice definition and referencing
* Extensive test coverage
* More operations:
  * Contains
  * Zero ( zero-value )
  * Etc...