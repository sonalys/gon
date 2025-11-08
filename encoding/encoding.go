package encoding

import "github.com/sonalys/gon"

type (
	NodeConstructor func(args []gon.KeyNode) (gon.Node, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)
