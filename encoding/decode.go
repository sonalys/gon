package encoding

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
)

func Decode(buffer []byte, codex Codex) (adapters.Node, error) {
	tokens := tokenize(buffer)
	parser := newParser(tokens)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	node, err := translateNode(rootNode, codex)
	if err != nil {
		return nil, fmt.Errorf("translating ast using codex: %w", err)
	}

	return node, nil
}
