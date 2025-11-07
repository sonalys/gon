package encoding

import (
	"fmt"
	"slices"

	"github.com/sonalys/gon"
)

func argSorter(from []gon.KeyExpression, keys ...string) (map[string]gon.Expression, []gon.Expression, error) {
	if len(from) < len(keys) {
		return nil, nil, fmt.Errorf("missing arguments")
	}

	expectedMap := make(map[string]gon.Expression, len(keys))
	rest := make([]gon.Expression, 0, len(from))

gotArgLoop:
	for fromIndex := range from {
		for keyIndex := range keys {
			if from[fromIndex].Key == "" || from[fromIndex].Key == keys[keyIndex] {
				expectedMap[keys[keyIndex]] = from[fromIndex].Expression
				keys = slices.Delete(keys, keyIndex, keyIndex+1)
				continue gotArgLoop
			}
		}
		rest = append(rest, from[fromIndex].Expression)
	}

	return expectedMap, rest, nil
}
