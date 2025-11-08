package gon

import (
	"cmp"
	"time"

	"golang.org/x/exp/constraints"
)

func safeGet[T any](slice []T, index int) T {
	if len(slice) <= index {
		var zero T
		return zero
	}
	return slice[index]
}

func cmpAny(firstValue, secondValue any) (int, bool) {
	switch c1 := firstValue.(type) {
	case int:
		c2, ok := secondValue.(int)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case int8:
		c2, ok := secondValue.(int8)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case int16:
		c2, ok := secondValue.(int16)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case int32:
		c2, ok := secondValue.(int32)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case int64:
		c2, ok := secondValue.(int64)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uint:
		c2, ok := secondValue.(uint)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uint8:
		c2, ok := secondValue.(uint8)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uint16:
		c2, ok := secondValue.(uint16)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uint32:
		c2, ok := secondValue.(uint32)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uint64:
		c2, ok := secondValue.(uint64)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case uintptr:
		c2, ok := secondValue.(uintptr)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case float32:
		c2, ok := secondValue.(float32)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case float64:
		c2, ok := secondValue.(float64)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case string:
		c2, ok := secondValue.(string)
		if !ok {
			return 0, false
		}
		return cmp.Compare(c1, c2), true
	case time.Time:
		c2, ok := secondValue.(time.Time)
		if !ok {
			return 0, false
		}
		if c1.Equal(c2) {
			return 0, true
		}

		if c1.Before(c2) {
			return -1, true
		}

		return 1, true
	default:
		return 0, false
	}
}

func sumAny(values ...any) (any, bool) {
	if len(values) == 0 {
		return 0, false
	}

	switch values[0].(type) {
	case int:
		values, ok := castAll[int](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case int8:
		values, ok := castAll[int8](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case int16:
		values, ok := castAll[int16](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case int32:
		values, ok := castAll[int32](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case int64:
		values, ok := castAll[int64](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uint:
		values, ok := castAll[uint](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uint8:
		values, ok := castAll[uint8](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uint16:
		values, ok := castAll[uint16](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uint32:
		values, ok := castAll[uint32](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uint64:
		values, ok := castAll[uint64](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case uintptr:
		values, ok := castAll[uintptr](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case float32:
		values, ok := castAll[float32](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	case float64:
		values, ok := castAll[float64](values...)
		if !ok {
			return 0, false
		}
		return sum(values...), true
	default:
		return 0, false
	}
}

func avgAny(values ...any) (any, bool) {
	if len(values) == 0 {
		return 0, false
	}

	switch values[0].(type) {
	case int:
		values, ok := castAll[int](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case int8:
		values, ok := castAll[int8](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case int16:
		values, ok := castAll[int16](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case int32:
		values, ok := castAll[int32](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case int64:
		values, ok := castAll[int64](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uint:
		values, ok := castAll[uint](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uint8:
		values, ok := castAll[uint8](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uint16:
		values, ok := castAll[uint16](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uint32:
		values, ok := castAll[uint32](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uint64:
		values, ok := castAll[uint64](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case uintptr:
		values, ok := castAll[uintptr](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case float32:
		values, ok := castAll[float32](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	case float64:
		values, ok := castAll[float64](values...)
		if !ok {
			return 0, false
		}
		return avg(values...), true
	default:
		return 0, false
	}
}

func sum[T constraints.Float | constraints.Integer](values ...T) T {
	var total T

	for i := range values {
		total += values[i]
	}

	return total
}

func avg[T constraints.Float | constraints.Integer](values ...T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}

	count := T(len(values))
	sum := sum(values...)

	return sum / count
}

func castAll[T any](values ...any) ([]T, bool) {
	output := make([]T, 0, len(values))

	for i := range values {
		value, ok := values[i].(T)
		if !ok {
			return nil, false
		}

		output = append(output, value)
	}

	return output, true
}
