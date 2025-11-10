package encoding

import "github.com/sonalys/gon"

type (
	NodeConstructor func(keyedNodes []gon.KeyNode) (gon.Node, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)
