package gon

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sonalys/gon/internal/sliceutils"
	"github.com/stretchr/testify/assert"
)

type baseNumericTestCase struct {
	name           string
	values         []int
	expectedResult any
	expectedOK     bool
}

func runNumericTests(
	t *testing.T,
	numericOnly bool,
	baseCases []baseNumericTestCase,
	run func(tc baseNumericTestCase, args ...any),
) {
	t.Helper()

	typeCasters := map[string]func(...int) []any{
		"int":     func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return int(from) }) },
		"int8":    func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return int8(from) }) },
		"int16":   func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return int16(from) }) },
		"int32":   func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return int32(from) }) },
		"int64":   func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return int64(from) }) },
		"uint":    func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uint(from) }) },
		"uint8":   func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uint8(from) }) },
		"uint16":  func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uint16(from) }) },
		"uint32":  func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uint32(from) }) },
		"uint64":  func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uint64(from) }) },
		"uintptr": func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return uintptr(from) }) },
		"float32": func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return float32(from) }) },
		"float64": func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return float64(from) }) },
		"string":  func(a ...int) []any { return sliceutils.Map(a, func(from int) any { return fmt.Sprint(from) }) },
		"time": func(a ...int) []any {
			return sliceutils.Map(a, func(from int) any { return time.Unix(int64(from), 0) })
		},
	}

	if numericOnly {
		delete(typeCasters, "string")
		delete(typeCasters, "time")
	}

	for typeName, caster := range typeCasters {
		for _, tc := range baseCases {
			t.Run(fmt.Sprintf("%s/%s", tc.name, typeName), func(t *testing.T) {
				t.Parallel()

				values := caster(tc.values...)
				run(tc, values...)
			})
		}
	}
}

func Test_sumAny(t *testing.T) {
	t.Parallel()

	baseCases := []baseNumericTestCase{
		{name: "sum of positive numbers", values: []int{1, 2, 3}, expectedResult: 6, expectedOK: true},
		{name: "sum with a single value", values: []int{100}, expectedResult: 100, expectedOK: true},
		{name: "sum with zero", values: []int{10, 0, -5}, expectedResult: 5, expectedOK: true},
	}

	runNumericTests(t, true, baseCases, func(tc baseNumericTestCase, args ...any) {
		got, ok := sumAny(args...)
		assert.Equal(t, tc.expectedOK, ok)
		assert.EqualValues(t, tc.expectedResult, got)

		_, ok = sumAny(append(args, nil)...)
		assert.False(t, ok)
	})

	t.Run("general cases", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name       string
			values     []any
			expected   any
			expectedOK bool
		}{
			{name: "empty slice", values: []any{}, expected: 0, expectedOK: false},
			{name: "type mismatch", values: []any{1, "2", 3}, expected: 0, expectedOK: false},
			{name: "unsupported type for sum", values: []any{"a", "b"}, expected: 0, expectedOK: false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				actual, ok := sumAny(tc.values...)
				assert.Equal(t, tc.expectedOK, ok)
				assert.Equal(t, tc.expected, actual)
			})
		}
	})
}

func Test_avgAny(t *testing.T) {
	t.Parallel()

	baseCases := []baseNumericTestCase{
		{name: "avg of positive numbers", values: []int{1, 2, 3, 4, 5}, expectedResult: 3, expectedOK: true},
		{name: "avg with a single value", values: []int{10}, expectedResult: 10, expectedOK: true},
	}

	runNumericTests(t, true, baseCases, func(tc baseNumericTestCase, args ...any) {
		got, ok := avgAny(args...)
		assert.Equal(t, tc.expectedOK, ok)
		assert.EqualValues(t, tc.expectedResult, got)

		_, ok = avgAny(append(args, nil)...)
		assert.False(t, ok)
	})

	t.Run("general cases", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name       string
			values     []any
			expected   any
			expectedOK bool
		}{
			{name: "empty slice", values: []any{}, expected: 0, expectedOK: false},
			{name: "type mismatch", values: []any{1, 2.5, 3}, expected: 0, expectedOK: false},
			{name: "unsupported type for avg", values: []any{time.Now()}, expected: 0, expectedOK: false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				actual, ok := avgAny(tc.values...)
				assert.Equal(t, tc.expectedOK, ok)
				assert.Equal(t, tc.expected, actual)
			})
		}
	})
}

func Test_safeGet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		slice    []any
		index    int
		expected any
	}{
		{name: "valid index", slice: []any{1, "hello", true}, index: 1, expected: "hello"},
		{name: "index out of bounds", slice: []any{1, 2, 3}, index: 5, expected: nil},
		{name: "empty slice", slice: []any{}, index: 0, expected: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := safeGet(tc.slice, tc.index)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_cmpAny(t *testing.T) {
	t.Parallel()

	baseCases := []baseNumericTestCase{
		{name: "equal", values: []int{1, 1}, expectedResult: 0, expectedOK: true},
		{name: "bigger", values: []int{1, 0}, expectedResult: 1, expectedOK: true},
		{name: "smaller", values: []int{0, 1}, expectedResult: -1, expectedOK: true},
	}

	runNumericTests(t, false, baseCases, func(tc baseNumericTestCase, args ...any) {
		got, ok := cmpAny(args[0], args[1])
		assert.Equal(t, tc.expectedOK, ok)
		assert.EqualValues(t, tc.expectedResult, got)

		_, ok = cmpAny(args[0], nil)
		assert.False(t, ok)
	})

}

func Test_castAll(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []struct {
		name       string
		values     []any
		castFunc   func() (any, bool)
		assertFunc func(t *testing.T, actual any, ok bool)
	}{
		{
			name: "successful cast to int",
			castFunc: func() (any, bool) {
				return castAll[int](1, 2, 3)
			},
			assertFunc: func(t *testing.T, actual any, ok bool) {
				assert.True(t, ok)
				assert.Equal(t, []int{1, 2, 3}, actual)
			},
		},
		{
			name: "failed cast due to type mismatch",
			castFunc: func() (any, bool) {
				return castAll[int](1, "two", 3)
			},
			assertFunc: func(t *testing.T, actual any, ok bool) {
				assert.False(t, ok)
				assert.Nil(t, actual)
			},
		},
	}

	for _, tc := range testCases {
		if ctx.Err() != nil {
			t.Skip("Context cancelled")
		}
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, ok := tc.castFunc()
			tc.assertFunc(t, actual, ok)
		})
	}
}
