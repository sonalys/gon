package encoding

import (
	"fmt"
	"slices"

	"github.com/sonalys/gon"
)

func argSorter(from []gon.KeyNode, keys ...string) (map[string]gon.Node, []gon.Node, error) {
	if len(from) < len(keys) {
		return nil, nil, fmt.Errorf("missing arguments")
	}

	expectedMap := make(map[string]gon.Node, len(keys))
	rest := make([]gon.Node, 0, len(from))

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
