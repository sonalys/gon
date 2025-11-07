package gon

import (
	"cmp"
	"fmt"
	"time"
)

func propagateErr(value Value, mask string, args ...any) Value {
	if value == nil {
		return Literal(fmt.Errorf(mask, args...))
	}

	if err, ok := value.Value().(error); ok {
		return Literal(err)
	}

	return Literal(fmt.Errorf(mask, args...))
}

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
