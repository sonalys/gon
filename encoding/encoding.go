package encoding

import (
	"fmt"

	"github.com/sonalys/gon"
)

type (
	NodeConstructor func(keyedNodes []gon.KeyNode) (gon.Node, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)

var DefaultExpressionCodex = Codex{}

func (c *Codex) Register(name string, constructor func([]gon.KeyNode) (gon.Node, error)) error {
	if _, conflicts := (*c)[name]; conflicts {
		return fmt.Errorf("node with name '%s' is already registered", name)
	}

	(*c)[name] = constructor

	return nil
}

type AutoRegisterer interface {
	Register(codex gon.Codex) error
}

func (c *Codex) AutoRegister(nodes ...AutoRegisterer) error {
	for _, registerer := range nodes {
		if err := registerer.Register(c); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	err := DefaultExpressionCodex.AutoRegister(
		&gon.AvgNode{},
		&gon.CallNode{},
		&gon.EqualNode{},
		&gon.GreaterNode{},
		&gon.HasPrefixNode{},
		&gon.HasSuffixNode{},
		&gon.IfNode{},
		&gon.LiteralNode{},
		&gon.NotNode{},
		&gon.OrNode{},
		&gon.SmallerNode{},
		&gon.SumNode{},
	)
	if err != nil {
		panic(fmt.Errorf("unexpected error registering default nodes: %s", err))
	}
}
