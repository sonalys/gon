package sliceutils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("should convert a slice of integers to a slice of strings", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []string{"1", "2", "3"}
		actual := Map(input, func(n int) string {
			return strconv.Itoa(n)
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should extract a field from a slice of structs", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		input := []User{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		expected := []string{"Alice", "Bob"}
		actual := Map(input, func(u User) string {
			return u.Name
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should return an empty slice when given an empty slice", func(t *testing.T) {
		input := []int{}
		expected := []string{}
		actual := Map(input, func(n int) string {
			return strconv.Itoa(n)
		})
		assert.Equal(t, expected, actual)
	})
}

func TestToMap(t *testing.T) {
	t.Run("should convert a slice of structs to a map", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		input := []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		expected := map[int]string{
			1: "Alice",
			2: "Bob",
		}
		actual := ToMap(input, func(u User) (int, string) {
			return u.ID, u.Name
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should handle the last value for duplicate keys", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		input := []User{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 1, Name: "Alicia"},
		}
		expected := map[int]string{
			1: "Alicia",
			2: "Bob",
		}
		actual := ToMap(input, func(u User) (int, string) {
			return u.ID, u.Name
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should return an empty map when given an empty slice", func(t *testing.T) {
		input := []struct{}{}
		expected := map[int]string{}
		actual := ToMap(input, func(s struct{}) (int, string) {
			return 0, ""
		})
		assert.Equal(t, expected, actual)
	})
}

func TestFilter(t *testing.T) {
	t.Run("should filter a slice of integers based on a condition", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []int{2, 4}
		actual := Filter(input, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should return the original slice when all elements match the filter", func(t *testing.T) {
		input := []int{2, 4, 6}
		expected := []int{2, 4, 6}
		actual := Filter(input, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should return an empty slice when no elements match the filter", func(t *testing.T) {
		input := []int{1, 3, 5}
		expected := []int{}
		actual := Filter(input, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("should return an empty slice when given an empty slice", func(t *testing.T) {
		input := []int{}
		expected := []int{}
		actual := Filter(input, func(n int) bool {
			return n%2 == 0
		})
		assert.Equal(t, expected, actual)
	})
}
