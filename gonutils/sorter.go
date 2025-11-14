package gonutils

import (
	"fmt"
	"slices"

	"github.com/sonalys/gon/adapters"
)

// SortArgs will parse any given keys as required.
// The required args will be put into the map, and error if any is missing.
// The rest of the keys found are appended to the slice.
func SortArgs(from []adapters.KeyNode, keys ...string) (map[string]adapters.Node, []adapters.Node, error) {
	if len(from) < len(keys) {
		return nil, nil, fmt.Errorf("missing arguments")
	}

	expectedMap := make(map[string]adapters.Node, len(keys))
	rest := make([]adapters.Node, 0, len(from))

gotArgLoop:
	for fromIndex := range from {
		for keyIndex := range keys {
			if from[fromIndex].Key == "" || from[fromIndex].Key == keys[keyIndex] {
				expectedMap[keys[keyIndex]] = from[fromIndex].Node
				keys = slices.Delete(keys, keyIndex, keyIndex+1)
				continue gotArgLoop
			}
		}
		rest = append(rest, from[fromIndex].Node)
	}

	return expectedMap, rest, nil
}
