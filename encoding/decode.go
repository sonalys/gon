package encoding

import "github.com/sonalys/gon"

func Decode(buffer []byte, codex Codex) (gon.Node, error) {
	tokens := tokenize(buffer)
	parser := newParser(tokens)

	rootNode, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return translateNode(rootNode, codex)
}
