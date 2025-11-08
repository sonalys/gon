package gon_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/stretchr/testify/require"
)

func Test_Literal_Definition(t *testing.T) {
	t.Run("structs/should resolve pointer structs", func(t *testing.T) {
		type str struct {
			Value int `gon:"value"`
		}

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(&str{
					Value: 5,
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.value")
		require.True(t, ok)

		require.Equal(t, 5, gotValue.Value())
	})

	t.Run("structs/should resolve pointer fields", func(t *testing.T) {
		type str struct {
			Value *int `gon:"value"`
		}

		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(str{
					Value: &expected,
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.value")
		require.True(t, ok)

		require.Equal(t, &expected, gotValue.Value())
	})

	t.Run("structs/should resolve nested fields", func(t *testing.T) {
		type attribute struct {
			Value int `gon:"value"`
		}

		type str struct {
			Attribute attribute `gon:"attribute"`
		}

		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(str{
					Attribute: attribute{
						Value: expected,
					},
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.attribute.value")
		require.True(t, ok)

		require.Equal(t, expected, gotValue.Value())
	})

	t.Run("structs/should resolve nested pointer fields", func(t *testing.T) {
		type attribute struct {
			Value int `gon:"value"`
		}

		type str struct {
			Attribute *attribute `gon:"attribute"`
		}

		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(str{
					Attribute: &attribute{
						Value: expected,
					},
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.attribute.value")
		require.True(t, ok)

		require.Equal(t, expected, gotValue.Value())
	})

	t.Run("maps/should resolve pointer structs", func(t *testing.T) {
		type str map[string]int

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(&str{
					"value": 5,
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.value")
		require.True(t, ok)

		require.Equal(t, 5, gotValue.Value())
	})

	t.Run("maps/should resolve pointer fields", func(t *testing.T) {
		type str map[string]*int

		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(str{
					"value": &expected,
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.value")
		require.True(t, ok)

		require.Equal(t, &expected, gotValue.Value())
	})

	t.Run("maps/should resolve nested fields", func(t *testing.T) {
		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(map[string]map[string]int{
					"attribute": {
						"value": expected,
					},
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.attribute.value")
		require.True(t, ok)

		require.Equal(t, expected, gotValue.Value())
	})

	t.Run("maps/should resolve nested pointer fields", func(t *testing.T) {
		expected := 5

		scope, err := gon.
			NewScope().
			WithDefinitions(gon.Definitions{
				"var": gon.Literal(map[string]*map[string]int{
					"attribute": {
						"value": expected,
					},
				}),
			})
		require.NoError(t, err)

		gotValue, ok := scope.Definition("var.attribute.value")
		require.True(t, ok)

		require.Equal(t, expected, gotValue.Value())
	})
}

func Test_Literal_Call(t *testing.T) {
	t.Run("should resolve callable attribute", func(t *testing.T) {
		expected := 5

		node := gon.Literal(map[string]map[string]func() int{
			"attribute": {
				"value": func() int { return expected },
			},
		})

		valued := node.Call(t.Context(), "attribute.value")
		require.Equal(t, expected, valued.Value())

		scope, err := gon.NewScope().WithDefinitions(gon.Definitions{
			"var": node,
		})
		require.NoError(t, err)

		scopeValued := gon.Call("var.attribute.value").Eval(scope)
		require.Equal(t, expected, scopeValued.Value())

	})
}

func Test_Literal_Value(t *testing.T) {
	t.Run("nil should not panic", func(t *testing.T) {
		node := gon.Literal(nil)

		require.NotPanics(t, func() {
			got := node.Value()
			require.Nil(t, got)
		})
	})

	t.Run("nested value should unwrap", func(t *testing.T) {
		inner := gon.Literal(1)
		outer := gon.Literal(inner)

		got := outer.Value()
		require.Equal(t, 1, got)
	})
}
